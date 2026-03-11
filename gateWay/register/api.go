package register

import (
	"gateWay/proto/rpcPb"
	"gateWay/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// 注册服务
func RegisterServer(s *grpc.Server)  {
	//这边统一注册
	rpcPb.RegisterBaseMessageServer(s,service.NewBaseMessageServer())

	reflection.Register(s)
}