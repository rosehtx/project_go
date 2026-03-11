package service

import (
	"context"
	"gateWay/proto/rpcPb"
)

type DispatchTask struct {
	ctx      context.Context
	request  *rpcPb.BaseRequest
	respChan chan *rpcPb.BaseResponse
}

type BaseMessageServer struct {
	taskChan chan *DispatchTask
}

func NewBaseMessageServer() *BaseMessageServer {
	server := &BaseMessageServer{
		//有缓冲，避免高并发阻塞
		taskChan: make(chan *DispatchTask, 1000),
	}
	// 启动worker pool
	for i := 0; i < 20; i++ {
		go server.worker()
	}
	return server
}

func (server *BaseMessageServer) worker() {
	//监控taskChan
	for task := range server.taskChan {
		// 这里做实际的分发逻辑（如根据BaseRequest.Service/Method路由到后端）
		requestData := task.request
		serverName := requestData.GetServer()

		//获取对应的连接
		conn, exists := GetServerCon(serverName)
		if exists == false {
			task.respChan <- &rpcPb.BaseResponse{
				Code: 404,
				Msg:  "后端服务未注册: " + serverName,
			}
		}
		//分发
		clint := rpcPb.NewBaseMessageClient(conn)
		resp, err := clint.GetBaseMessage(task.ctx, requestData)
		if err != nil {
			task.respChan <- &rpcPb.BaseResponse{
				Code: 500,
				Msg:  err.Error(),
			}
		}
		// 回传响应
		task.respChan <- resp
	}
}

func (server *BaseMessageServer) GetBaseMessage(ctx context.Context, requestData *rpcPb.BaseRequest) (*rpcPb.BaseResponse, error) {
	respChan := make(chan *rpcPb.BaseResponse, 1)
	task := &DispatchTask{
		ctx:      ctx,
		request:  requestData,
		respChan: respChan,
	}
	// 投递到异步分发队列
	server.taskChan <- task

	// 阻塞等待worker处理并返回结果（支持超时）
	select {
	case resp := <-respChan:
		return resp, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}
