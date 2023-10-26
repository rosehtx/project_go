package main

import (
	"fmt"
	"google.golang.org/grpc"
	"heroServer/register"
	"net"
)

func main() {
	// 监听本地端口
	lis, err := net.Listen("tcp", ":8050")
	if err != nil {
		fmt.Printf("监听端口失败: %s", err)
		return
	}

	// 创建gRPC服务器
	server := grpc.NewServer()
	// 注册服务
	register.RegisterServer(server)

	fmt.Printf("服务启动port: %s", "8050")
	err = server.Serve(lis)
	if err != nil {
		fmt.Printf("开启服务失败: %s", err)
		return
	}
}
