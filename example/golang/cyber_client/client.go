package main

import (
	"context"
	"fmt"
	service "grpc_demo/service"
	"io"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/keepalive"
)

const (
	address = "127.0.0.1:50051"
)

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

	kacp := keepalive.ClientParameters{
		Time:                10 * time.Second, // send pings every 10 seconds if there is no activity
		Timeout:             time.Second,      // wait 1 second for ping ack before considering the connection dead
		PermitWithoutStream: true,             // send pings even without active streams
	}

	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithKeepaliveParams(kacp))
	if err != nil {
		grpclog.Fatal("Error in dial.\r\n")
	}
	defer conn.Close()

	c := service.NewCyberManagerClient(conn)

	CyberControlFunc(c)
}
