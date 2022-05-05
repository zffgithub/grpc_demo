package main

import (
	"fmt"
	service "grpc_demo/service"
	"io"
	"net"
	"strconv"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
)

const (
	PORT = ":50051"
)

type Server struct{}

func (s *Server) SayHello(ctx context.Context, in *service.HelloRequest) (*service.HelloReply, error) {
	return &service.HelloReply{Message: "hello " + in.Name + "\r\n"}, nil
}

//简单模式
func (s *Server) EasyMode(ctx context.Context, in *service.Request) (*service.Response, error) {
	fmt.Println("###############简单模式###############")
	//简单模式不需要调用recv和send。数据直接通过方法的参数和返回值来接收发送。
	re := &service.Response{No: in.CyberId, Msg: "Receive one data"}
	return re, nil
}

//客户端流模式
func (s *Server) ClientStream(stream service.CyberManager_ClientStreamServer) error {
	fmt.Println("###############客户端流模式###############")
	//客户端端的代码stream是CyberManager_ClientStreamServer类型，这种类型只有recv和sendandclose两种方法接收和发送并关闭
	//而服务器跟着相反是send和closeandrecv两种方法。
	//客户端流模式要先接收然后再发送。
	//使用for循环接收数据，因为接收的是流模式的Request类型，所以每次循环得到的是一个studentRequest类型
	for {
		in, err := stream.Recv() //通过for循环接收客户端传来的流，该方法协程不安全。
		//必须把判断EOF的条件放在nil前面，因为如果获取完数据err会返回EOF类型
		if err == io.EOF { //在读取完所有数据之后会返回EOF来执行该条件。
			fmt.Println("Recieve done.")
			//发送（只在读取完数据后发送了一次）
			rp := &service.Response{No: int32(0), Msg: "Recieve all data."}
			//注意：该方法内部调用的SendMsg方法，但是该方法协程不安全。不同协程可以同时调用会出现问题。
			return stream.SendAndClose(rp)
		}
		if err != nil {
			fmt.Println("Error in ClientStream:", err)
			return err
		}
		fmt.Println("Recieve data:", in)
	}
}

//服务端流模式
func (s *Server) ServerStream(in *service.Request, stream service.CyberManager_ServerStreamServer) error {
	fmt.Println("###############服务端流模式###############")
	fmt.Println("Recieve one data:", in)
	for i := int32(1); i < 10; i++ {
		p := &service.Response{No: in.CyberId + i, Msg: fmt.Sprintf("Recieve client data:%d", in.CyberId)}
		stream.Send(p)
	}
	return nil
}

//双向模式
func (s *Server) BidirectionalStream(stream service.CyberManager_BidirectionalStreamServer) error {
	fmt.Println("###############双端模式###############")
	i := int32(100000)
	for {
		select {
		case <-stream.Context().Done():
			fmt.Println("Client disconnect by context done.")
			return stream.Context().Err()
		default:
			//接收
			in, err := stream.Recv()
			if err == io.EOF {
				fmt.Println("Recieve done.")
				return nil
			}
			if err != nil {
				fmt.Println("Error in BidirectionalStream:", err)
			}
			fmt.Println(in, strconv.FormatInt(time.Now().UTC().UnixNano(), 10))

			//发送
			i += 1
			rp := &service.Response{No: i, Msg: "test"}
			stream.Send(rp)
			time.Sleep(1 * time.Second)
		}

	}
	// return nil
}

var q *Client

func init() {
	fmt.Println("Init")
	q = NewClient()
	q.SetConditions(100)
}

//简单模式
func (s *Server) RemoteControl(ctx context.Context, in *service.Request) (*service.Response, error) {
	fmt.Println("RemoteControl================简单模式================")
	re := &service.Response{No: in.CyberId, Msg: in.Danmu}
	topic := string(in.CyberId)
	msg := fmt.Sprintf("RemoteControl Server send time:%s, Danmu:%s", time.Now().Format("2006/01/02 15:04:05"), in.Danmu)
	if err := q.Publish(topic, msg); err != nil {
		fmt.Println(err)
	}
	return re, nil
}

//服务端流模式
func (s *Server) CyberControl(in *service.Request, stream service.CyberManager_CyberControlServer) error {
	fmt.Println("CyberControl###############服务端流模式###############")
	fmt.Println("Receive client request:", in)
	topic := string(in.CyberId)
	ch, err := q.Subscribe(topic)
	if err != nil {
		fmt.Println(err)
	}
	for {
		select {
		case <-stream.Context().Done():
			fmt.Println("Client disconnect by context done.")
			return stream.Context().Err()
		default:
			msg := q.GetPayLoad(ch)
			fmt.Println("CyberControl Receive queue data:", msg)
			p := &service.Response{No: in.CyberId, Msg: fmt.Sprintf("CyberControl Server push data:%v", msg)}
			send_err := stream.Send(p)
			if send_err != nil {
				fmt.Println("Send msg error:", send_err)
			}
		}
	}
	// var wg sync.WaitGroup
	// wg.Add(1)
	// go func() {
	// 	msg := q.GetPayLoad(ch)
	// 	p := &service.Response{No: in.CyberId, Msg: fmt.Sprintf("Server push data:%v", msg)}
	// 	send_err := stream.Send(p)
	// 	if send_err != nil {
	// 		fmt.Println("Send msg error:", send_err)
	// 	}
	// 	fmt.Println("###############")
	// 	if err := q.Unsubscribe(topic, ch); err != nil {
	// 		fmt.Println(err)
	// 	}
	// 	wg.Done()
	// }()
	// wg.Wait()
	// return nil
}

func main() {
	lis, err := net.Listen("tcp", PORT)
	if err != nil {
		grpclog.Fatal("Error start tcp listen.\r\n")
	}

	kaep := keepalive.EnforcementPolicy{
		MinTime:             5 * time.Second, // If a client pings more than once every 5 seconds, terminate the connection
		PermitWithoutStream: true,            // Allow pings even when there are no active streams
	}

	kasp := keepalive.ServerParameters{
		MaxConnectionIdle:     15 * time.Second,   // If a client is idle for 15 seconds, send a GOAWAY
		MaxConnectionAge:      3600 * time.Second, // If any connection is alive for more than 30 seconds, send a GOAWAY
		MaxConnectionAgeGrace: 5 * time.Second,    // Allow 5 seconds for pending RPCs to complete before forcibly closing connections
		Time:                  5 * time.Second,    // Ping the client if it is idle for 5 seconds to ensure the connection is still active
		Timeout:               3 * time.Second,    // Wait 1 second for the ping ack before assuming the connection is dead
	}

	opts := []grpc.ServerOption{
		grpc.KeepaliveEnforcementPolicy(kaep),
		grpc.KeepaliveParams(kasp),
	}
	s := grpc.NewServer(opts...)
	// 注册服务
	service.RegisterCyberManagerServer(s, &Server{})
	reflection.Register(s)
	s.Serve(lis)
	if err != nil {
		fmt.Printf("开启服务失败: %s \r\n", err)
		return
	}
}
