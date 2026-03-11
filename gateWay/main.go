package main

import (
	"fmt"
	"gateWay/register"
	"gateWay/service"
	"google.golang.org/grpc"
	"net"
)

func main() {
	//连接其他服务
	conError := service.ConnectServer()
	if conError != nil {
		fmt.Printf("连接其他服务失败: %s", conError)
		return
	}
	//开启gateWay的服务
	lis, err := net.Listen("tcp", ":8010")
	if err != nil {
		fmt.Printf("监听端口失败: %s", err)
		return
	}
	// 创建gRPC服务器
	server := grpc.NewServer()
	// 注册服务
	register.RegisterServer(server)

	fmt.Printf("服务启动port: %s", "8010")
	err = server.Serve(lis)
	if err != nil {
		fmt.Printf("开启服务失败: %s", err)
		return
	}

}
