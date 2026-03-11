package service

import (
	"context"
	"fmt"
	"heroServer/proto/rpcPb"
	"heroServer/proto/rpcPb/hero"

	//"github.com/golang/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

type BaseMessageServer struct {
	rpcPb.UnimplementedBaseMessageServer
	userHeroService *GetUserHeroServer
}

func NewBaseMessageServer() *BaseMessageServer {
	return &BaseMessageServer{
		userHeroService: &GetUserHeroServer{},
	}
}

func (s *BaseMessageServer) GetBaseMessage(ctx context.Context, req *rpcPb.BaseRequest) (*rpcPb.BaseResponse, error) {
	// 检查是否有payload
	if req.Payload == nil {
		return &rpcPb.BaseResponse{
			Code: 400,
			Msg:  "Payload is required",
		}, nil
	}

	// 根据service和method进行路由
	switch req.Service {
	case "hero":
		return s.handleHeroService(ctx, req.Method, req.Payload)
	default:
		return &rpcPb.BaseResponse{
			Code: 404,
			Msg:  fmt.Sprintf("Service %s not found", req.Service),
		}, nil
	}
}

func (s *BaseMessageServer) handleHeroService(ctx context.Context, method string, payload *anypb.Any) (*rpcPb.BaseResponse, error) {
	switch method {
	case "GetUserHero":
		return s.handleGetUserHero(ctx, payload)
	default:
		return &rpcPb.BaseResponse{
			Code: 404,
			Msg:  fmt.Sprintf("Method %s not found in hero service", method),
		}, nil
	}
}

func (s *BaseMessageServer) handleGetUserHero(ctx context.Context, payload *anypb.Any) (*rpcPb.BaseResponse, error) {
	// 将Any转换为GetUserHeroRequest
	var request hero.GetUserHeroRequest
	if err := payload.UnmarshalTo(&request); err != nil {
		return &rpcPb.BaseResponse{
			Code: 500,
			Msg:  fmt.Sprintf("Failed to unmarshal payload: %v", err),
		}, nil
	}

	// 调用GetUserHero服务
	response, err := s.userHeroService.GetUserHero(ctx, &request)
	if err != nil {
		return &rpcPb.BaseResponse{
			Code: 500,
			Msg:  fmt.Sprintf("GetUserHero service error: %v", err),
		}, nil
	}

	// 将response转换为Any
	anyResponse, err := anypb.New(response)
	if err != nil {
		return &rpcPb.BaseResponse{
			Code: 500,
			Msg:  fmt.Sprintf("Failed to marshal response: %v", err),
		}, nil
	}

	return &rpcPb.BaseResponse{
		Code:    200,
		Msg:     "Success",
		Payload: anyResponse,
	}, nil
}
