package main

import (
	"context"
	"fmt"
	service "grpc_demo/service"
	"io"
	"strconv"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/grpclog"
)

const (
	address = "127.0.0.1:50051"
)

func EasyModeFunc(client service.CyberManagerClient, info *service.Request) {
	fmt.Println("###############简单模式###############")
	//简单模式下直接调用方法就可以传递参数和获得返回值
	re, err := client.EasyMode(context.Background(), info)
	if err != nil {
		fmt.Println("Error in EasyModeFunc:", err)
	}
	fmt.Println("Recieve one data:", re)
}

func ClientStreamFunc(client service.CyberManagerClient) {
	fmt.Println("###############客户端流模式###############")
	re, err := client.ClientStream(context.Background())
	if err != nil {
		fmt.Println("Error in ClientStreamFunc", err)
	}
	for i := int32(1); i < 10; i++ {
		sq := &service.Request{CyberId: i}
		re.Send(sq) //客户端要先发送数据再接收
	}

	//使用while循环接收数据
	for {
		//先关闭send然后再接收数据。该方法内部调用了CloseSend和RecvMsg方法。但是后者协程不安全
		r, err2 := re.CloseAndRecv()
		if err2 == io.EOF {
			fmt.Println("Recieve done.")
			return
		}
		if err2 != nil {
			fmt.Println("Error in ClientStreamFunc:", err2)
			return
		}
		fmt.Println(r)
	}
}

func ServerStreamFunc(client service.CyberManagerClient) {
	fmt.Println("###############服务端流模式###############")
	s := &service.Request{CyberId: 2}
	//向服务器端以参数形式发送数据s，并创建了一个通道，通过通道来获取服务器端返回的流
	re, err := client.ServerStream(context.Background(), s)
	if err != nil {
		fmt.Println("Error in ServerStreamFunc:", err)
	}
	for {
		//上面已经以参数的形式发送了数据，此处直接使用Recv接收数据。因是参数发送，所有不用使用CloseAndRecv
		v, err2 := re.Recv()
		if err2 == io.EOF {
			fmt.Println("Recieve done.")
			return
		}
		if err2 != nil {
			fmt.Println("Error in ServerStreamFunc:", err2)
			return
		}
		fmt.Println("Recieve one data:", v)
	}
}

func BidirectionalStreamFunc(client service.CyberManagerClient) {
	fmt.Println("###############双端模式###############")
	re, err := client.BidirectionalStream(context.Background())
	if err != nil {
		fmt.Println("Error in BidirectionalStreamFunc:", err)
	}
	//创建一个通道作为控制协程协作
	var wg sync.WaitGroup
	wg.Add(1)
	//开启协程进行接收信息
	go func() {
		for {
			in, err2 := re.Recv()
			if err2 == io.EOF {
				fmt.Println("Receive done.")
				wg.Done()
				return
			}
			if err2 != nil {
				fmt.Println("Error in BidirectionalStreamFunc:", err2)
			}
			fmt.Printf("Got message:%s, no:%d, time:%s\r\n", in.Msg, in.No, strconv.FormatInt(time.Now().UTC().Unix(), 10))
			time.Sleep(500 * time.Microsecond)
		}
	}()
	//主进程继续执行发送
	for i := int32(1); i < 10; i++ {
		s1 := &service.Request{CyberId: i}
		err3 := re.Send(s1)
		if err3 != nil {
			fmt.Println("Error in BidirectionalStreamFunc:", err3)
		}
	}

	re.CloseSend()
	wg.Wait()
}

func RemoteControlFunc(client service.CyberManagerClient, info *service.Request) {
	fmt.Println("###############简单模式###############")
	//简单模式下直接调用方法就可以传递参数和获得返回值
	re, err := client.RemoteControl(context.Background(), info)
	if err != nil {
		fmt.Println("Error in RemoteControlFunc:", err)
	}
	fmt.Println("Recieve one data:", re)
}

func CyberControlFunc(client service.CyberManagerClient) {
	fmt.Println("###############服务端流模式###############")
	s := &service.Request{CyberId: 2}
	//向服务器端以参数形式发送数据s，并创建了一个通道，通过通道来获取服务器端返回的流
	re, err := client.CyberControl(context.Background(), s)
	if err != nil {
		fmt.Println("Error in CyberControlFunc:", err)
	}
	for {
		//上面已经以参数的形式发送了数据，此处直接使用Recv接收数据。因是参数发送，所有不用使用CloseAndRecv
		v, err2 := re.Recv()
		if err2 == io.EOF {
			fmt.Println("Recieve done.")
			return
		}
		if err2 != nil {
			fmt.Println("Error in ServerStreamFunc:", err2)
			return
		}
		fmt.Println("Recieve one data:", v)
	}
}

func main() {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		grpclog.Fatal("Error in dial.\r\n")
	}
	defer conn.Close()

	c := service.NewCyberManagerClient(conn)

	// 调用服务端函数
	// res, err := c.SayHello(context.Background(), &service.HelloRequest{Name: "CyberName"})
	// if err != nil {
	// 	fmt.Printf("调用服务端代码失败: %s", err)
	// 	return
	// }
	// fmt.Printf("调用成功: %s", res.Message)

	req := &service.Request{CyberId: int32(2)}
	// EasyModeFunc(c, req)
	// ClientStreamFunc(c)
	// ServerStreamFunc(c)
	// BidirectionalStreamFunc(c)
	RemoteControlFunc(c, req)
	// CyberControlFunc(c)
}
