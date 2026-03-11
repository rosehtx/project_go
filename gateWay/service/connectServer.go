package service

import (
	"context"
	"fmt"
	"gateWay/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
	"sync"
	"time"
)

type backendEntry struct {
	name string
	addr string

	mu    sync.Mutex
	conn  *grpc.ClientConn
	alive bool // 连接是否认为健康

	cancel context.CancelFunc // 停止管理goroutine
}

var (
	backendsMu sync.Mutex
	backends   = make(map[string]*backendEntry)
)

var serverMap = make(map[string]*grpc.ClientConn, 10)

func ConnectServer() error {
	targets := map[string]string{
		config.AGENT_NAME: config.AGENT_ADDRESS,
	}

	//多个服务进行注册
	for name, addr := range targets {
		addOrUpdateBackend(name, addr)
	}

	//作为client 连接其他服务
	//conn, err := grpc.Dial(
	//	config.AGENT_ADDRESS,
	//	grpc.WithTransportCredentials(insecure.NewCredentials()),
	//)
	//if err != nil {
	//	fmt.Printf("连接服务端失败: %s", err)
	//	return err
	//}
	//serverMap[config.AGENT_NAME] = conn

	//... 多个服务注册map
	return nil
}

// addOrUpdateBackend 新增或更新后端，并启动后台管理协程（首次添加时）
func addOrUpdateBackend(name, addr string) {
	backendsMu.Lock()
	defer backendsMu.Unlock()

	e := &backendEntry{
		name: name,
		addr: addr,
	}
	backends[name] = e

	ctx, cancel := context.WithCancel(context.Background())
	e.cancel = cancel

	// 启动管理 goroutine 负责连接/重连/健康检查
	go manageBackend(ctx, e)
}

// manageBackend 负责建立连接、检测连接状态以及在断连时重试
func manageBackend(ctx context.Context, e *backendEntry) {
	const (
		baseDelay     = 500 * time.Millisecond
		maxDelay      = 30 * time.Second
		checkInterval = 5 * time.Second // 健康检查间隔
		dialTimeout   = 5 * time.Second // Dial 超时
	)

	var backoff time.Duration = baseDelay

	for {
		// 先尝试建立连接
		dialCtx, dialCancel := context.WithTimeout(ctx, dialTimeout)
		conn, err := grpc.DialContext(
			dialCtx,
			e.addr,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithBlock(), // 等待连接建立完成
			//配置保活，避免连接闲置断开
			grpc.WithKeepaliveParams(keepalive.ClientParameters{
				Time:    30 * time.Second, // 每30秒发一次保活请求
				Timeout: 5 * time.Second,  // 保活请求超时时间
			}),
		)
		dialCancel()
		fmt.Printf("连接服务: %s", e.name)
		if err != nil {
			// 建连失败，标记为不健康，等待并重试（指数退避）
			markBackendUnhealthy(e)
			select {
			case <-time.After(backoff):
				// 增加退避
				backoff *= 2
				if backoff > maxDelay {
					backoff = maxDelay
				}
				// 继续下一轮重试
			case <-ctx.Done():
				return
			}
			continue
		}

		// 建连成功，重置退避
		backoff = baseDelay
		setBackendConn(e, conn)

		// 进入健康轮询循环，根据 conn 状态判断是否需要重连
		// 如果连接变为非 READY，则清理并重新 dial
		// 也会周期性地（checkInterval）检查状态，以便尽快发现问题
		for {
			// 快速判断 ctx
			select {
			case <-ctx.Done():
				// 外部要求停止：关闭连接并返回
				cleanupBackend(e)
				return
			default:
			}

			// 检查当前连接状态
			state := conn.GetState()
			if state == connectivity.Ready {
				markBackendHealthy(e)
			} else {
				// 非 Ready 状态，认为不稳定，触发重连流程
				markBackendUnhealthy(e)
				// 先 close 当前 conn（若未被替换）
				e.mu.Lock()
				if e.conn == conn {
					_ = e.conn.Close()
					e.conn = nil
				}
				e.mu.Unlock()
				break // 跳出内部循环，进入上层的重连逻辑
			}

			// 等待 state 变化或超时轮询
			// 使用 WaitForStateChange 以更快速响应状态变化
			waitCtx, waitCancel := context.WithTimeout(ctx, checkInterval)
			changed := conn.WaitForStateChange(waitCtx, state)
			waitCancel()

			if !changed {
				// 超时且状态未变化，继续下一次轮询（再次检查状态）
				continue
			}
			// 若 changed 为 true -> 状态变了，立刻在下一次循环中检查新状态
		}
	}
}

// setBackendConn 安全设置 conn 字段
func setBackendConn(e *backendEntry, conn *grpc.ClientConn) {
	e.mu.Lock()
	// 关闭旧 conn（若存在且跟新 conn 不同）
	if e.conn != nil && e.conn != conn {
		_ = e.conn.Close()
	}
	e.conn = conn
	e.mu.Unlock()
}

// cleanupBackend 关闭并清理 entry 上的连接（在 stop 时调用）
func cleanupBackend(e *backendEntry) {
	e.mu.Lock()
	if e.conn != nil {
		_ = e.conn.Close()
		e.conn = nil
	}
	e.alive = false
	e.mu.Unlock()
}

// markBackendHealthy / Unhealthy 安全设置健康状态
func markBackendHealthy(e *backendEntry) {
	e.mu.Lock()
	e.alive = true
	e.mu.Unlock()
}
func markBackendUnhealthy(e *backendEntry) {
	e.mu.Lock()
	e.alive = false
	e.mu.Unlock()
}

// GetServerCon 返回连接（并发安全）以及是否健康
func GetServerCon(name string) (*grpc.ClientConn, bool) {
	backendsMu.Lock()
	e, ok := backends[name]
	backendsMu.Unlock()
	if !ok || e == nil {
		return nil, false
	}
	e.mu.Lock()
	conn := e.conn
	alive := e.alive
	e.mu.Unlock()
	return conn, alive
}

// CloseAllServerConns 停止所有管理 goroutine 并关闭连接，供程序优雅退出时调用
func CloseAllServerConns() error {
	backendsMu.Lock()
	defer backendsMu.Unlock()

	var lastErr error
	for name, e := range backends {
		// 取消管理 goroutine
		if e.cancel != nil {
			e.cancel()
		}
		// 关闭 conn
		e.mu.Lock()
		if e.conn != nil {
			if err := e.conn.Close(); err != nil {
				lastErr = err
			}
			e.conn = nil
		}
		e.mu.Unlock()

		delete(backends, name)
	}
	return lastErr
}

//func GetServerCon(name string) (*grpc.ClientConn,bool){
//	value, exists := serverMap[name]
//	return value, exists
//}
