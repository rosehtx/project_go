package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"heroServer/proto/rpcPb/hero"
	"testing"
)

func TestMainTest(t *testing.T)  {
	// 连接服务器
	conn, err := grpc.Dial(":8050", grpc.WithInsecure())
	if err != nil {
		fmt.Printf("连接服务端失败: %s", err)
		return
	}
	defer conn.Close()

	// 新建一个客户端
	c := hero.NewGetUserHeroClient(conn)

	// 调用服务端函数
	r, err := c.GetUserHero(context.Background(), &hero.GetUserHeroRequest{Id: 444})
	if err != nil {
		fmt.Printf("调用服务端代码失败: %s", err)
		return
	}

	fmt.Println(r)
	//fmt.Printf("调用成功: %s", r.Name)
	//fmt.Println(r.Age)
}
