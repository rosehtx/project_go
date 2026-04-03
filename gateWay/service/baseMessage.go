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

// sendResp 向调用方回写唯一一次结果；若 ctx 已取消则不再阻塞 worker
func sendResp(task *DispatchTask, resp *rpcPb.BaseResponse) {
	if task == nil || resp == nil {
		return
	}
	select {
	case task.respChan <- resp:
	case <-task.ctx.Done():
	}
}

func (server *BaseMessageServer) worker() {
	for task := range server.taskChan {
		requestData := task.request
		serverName := requestData.GetServer()

		conn, ok := GetServerCon(serverName)
		if !ok {
			sendResp(task, &rpcPb.BaseResponse{
				Code: 404,
				Msg:  "后端服务未注册: " + serverName,
			})
			continue
		}

		clint := rpcPb.NewBaseMessageClient(conn)
		resp, err := clint.GetBaseMessage(task.ctx, requestData)
		if err != nil {
			sendResp(task, &rpcPb.BaseResponse{
				Code: 500,
				Msg:  err.Error(),
			})
			continue
		}
		sendResp(task, resp)
	}
}

func (server *BaseMessageServer) GetBaseMessage(ctx context.Context, requestData *rpcPb.BaseRequest) (*rpcPb.BaseResponse, error) {
	respChan := make(chan *rpcPb.BaseResponse, 1)
	task := &DispatchTask{
		ctx:      ctx,
		request:  requestData,
		respChan: respChan,
	}
	server.taskChan <- task

	select {
	case resp := <-respChan:
		return resp, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}
