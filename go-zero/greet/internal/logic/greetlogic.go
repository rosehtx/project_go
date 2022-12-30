package logic

import (
	"context"
	"errors"

	"goZero/greet/internal/svc"
	"goZero/greet/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GreetLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGreetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GreetLogic {
	return &GreetLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GreetLogic) Greet(req *types.Request) (resp *types.Response, err error) {
	// todo: add your logic here and delete this line
	l.Logger.Info("sdfdsfds")
	return &types.Response{
		Message: "Hello go-zero",
	}, errors.New("errrrrrr")
}

func (l *GreetLogic) GreetTwo(req *types.Request) (resp *types.Response, err error) {
	l.Logger.Info("GreetTwo get")
	return &types.Response{
		Message: req.Name,
	},nil
}

