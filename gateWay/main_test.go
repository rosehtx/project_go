package main

import (
	"context"
	"fmt"
	"gateWay/config"
	"gateWay/proto/rpcPb"
	"gateWay/proto/rpcPb/hero"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/anypb"
	"testing"
)

func TestMainTest(t *testing.T) {
	// 连接服务器
	conn, err := grpc.Dial(":8010", grpc.WithInsecure())
	if err != nil {
		fmt.Printf("连接服务端失败: %s", err)
		return
	}
	defer conn.Close()

	// 新建一个客户端
	c := rpcPb.NewBaseMessageClient(conn)

	// 调用服务端函数
	request := &hero.GetUserHeroRequest{
		Id: 4433,
	}
	anyRequest, err := anypb.New(request)
	r, err := c.GetBaseMessage(context.Background(), &rpcPb.BaseRequest{
		Server:  config.AGENT_NAME,
		Service: "hero",
		Method:  "GetUserHero",
		Payload: anyRequest,
	})
	if err != nil {
		fmt.Printf("调用服务端代码失败: %s", err)
		return
	}

	fmt.Println(r)
	//fmt.Printf("调用成功: %s", r.Name)
	//fmt.Println(r.Age)
}
