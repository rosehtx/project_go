package register

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"heroServer/proto/rpcPb/hero"
	"heroServer/service"
)

// 注册服务
func RegisterServer(s *grpc.Server)  {
	//这边统一注册
	hero.RegisterGetUserHeroServer(s,&service.GetUserHeroServer{})


	reflection.Register(s)
}