package service

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"serverList/config"
	"strconv"
	"sync"
)

// RabbitMQConnectionPool 包装 RabbitMQ 连接池
type RabbitMQConnectionPool struct {
	pool      	chan *amqp.Connection
	channelPool map[*amqp.Connection]*RabbitMQChannelPool//用来记录对应的channel连接池
	mu        sync.Mutex
	poolSize  int
	available int
}

// RabbitMQChannelPool 包装 RabbitMQ 的channel池
type RabbitMQChannelPool struct {
	pool      chan *amqp.Channel
	mu        sync.Mutex
	poolSize  int
	available int
}
//用来记录连接池
var RabbitMQConnectionPoolPtr *RabbitMQConnectionPool

// NewRabbitMQConnectionPool 建立连接池
func NewRabbitMQConnectionPool(poolSize int) (*RabbitMQConnectionPool, error) {
	pool := make(chan *amqp.Connection, poolSize)
	channelPool := make(map[*amqp.Connection]*RabbitMQChannelPool)//channel连接池的map
	for i := 0; i < poolSize; i++ {
		fmt.Println("创建第"+strconv.Itoa(i+1)+"个rmq连接")
		conn, err := amqp.Dial("amqp://"+config.RMQ_USER+":"+config.RMQ_PASS+"@"+config.RMQ_IP+":"+strconv.Itoa(config.RMQ_PORT)+"/"+config.RMQ_VHOST)
		if err != nil {
			return nil, fmt.Errorf("Failed to create RabbitMQ connection: %v", err)
		}
		pool <- conn
		//channel的池子处理
		rabbitMQChannelPool,_ := NewRabbitMQChannelPool(conn,i)
		if rabbitMQChannelPool == nil {
			return nil, fmt.Errorf("创建channenl异常")
		}
		channelPool[conn] = rabbitMQChannelPool
	}

	//获取连接池使用
	RabbitMQConnectionPoolPtr = &RabbitMQConnectionPool{
		pool:      		pool,
		channelPool:    channelPool,
		poolSize:  		poolSize,
		available: 		poolSize,
	}
	return RabbitMQConnectionPoolPtr, nil
}

// NewRabbitMQChannelPool 建立channel连接池
func NewRabbitMQChannelPool(conn *amqp.Connection,num int) (*RabbitMQChannelPool, error) {
	pool := make(chan *amqp.Channel, config.RMQ_CHANNEL_NUM)
	fmt.Println("创建第"+strconv.Itoa(num + 1)+"连接的channel通道\n")
	//channel的池子处理
	for i := 0; i < config.RMQ_CHANNEL_NUM; i++ {
		channel, channelErr := conn.Channel()
		if channelErr != nil{
			fmt.Println("创建rabbitmq通道error"+channelErr.Error())
			return nil,nil
		}
		pool <- channel
	}
	return &RabbitMQChannelPool{
		pool:pool,
		poolSize:config.RMQ_CHANNEL_NUM,
		available:config.RMQ_CHANNEL_NUM,
	}, nil
}

// GetConnection 从连接池获取一个 RabbitMQ 连接,
func (p *RabbitMQConnectionPool) GetConnection() (*amqp.Connection,error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	select {
		case conn := <-p.pool:
			p.available--
			fmt.Printf("conn 获取一个rmq连接剩余:%v \n",p.available)
			return conn, nil
		default:
			/**
				可以有俩种方式处理
				1.直接返回没有可用连接，比较便捷，防止创建过多连接且未释放
				2.重新创建新的连接
			 */
			return nil, fmt.Errorf("没有可用的连接")
	}
}

// ReleaseConnection 将连接放回连接池
func (p *RabbitMQConnectionPool) ReleaseConnection(conn *amqp.Connection) {
	p.mu.Lock()
	defer p.mu.Unlock()

	//数量小于池子总数 既定的数量
	if p.available < p.poolSize{
		p.pool <- conn
		p.available++
		fmt.Printf("conn 放回一个rmq连接剩余:%v,总数:%v \n",p.available,p.poolSize)
	} else {
		_ = conn.Close()
	}
}

// GetChannel 从连接池获取一个 channel 连接
func (p *RabbitMQConnectionPool) GetChannel() (*amqp.Connection,*amqp.Channel, error) {
	conn,err  := p.GetConnection()
	if err != nil {
		return nil, nil,err
	}

	p.channelPool[conn].mu.Lock()
	defer p.channelPool[conn].mu.Unlock()
	//获取完channel直接吧连接放回池子 提升并发能力
	p.ReleaseConnection(conn)

	select {
	case channel := <-p.channelPool[conn].pool:
		p.channelPool[conn].available--
		fmt.Printf("获取一个rmq channel:%v \n",p.channelPool[conn].available)
		return conn , channel, nil
	default:
		//return nil, nil ,fmt.Errorf("没有可用的channel")
		channel, channelErr := conn.Channel()
		if channelErr != nil{
			fmt.Println("创建rabbitmq通道error"+channelErr.Error())
			return nil, nil ,fmt.Errorf("没有可用的channel")
		}
		fmt.Printf("创建一个新的channel:%v \n",p.channelPool[conn].available)
		return conn , channel, nil
	}
}

// ReleaseChannel 将连接放回连接池
func (p *RabbitMQConnectionPool) ReleaseChannel(conn *amqp.Connection , channel *amqp.Channel) {
	p.channelPool[conn].mu.Lock()
	defer p.channelPool[conn].mu.Unlock()

	if p.channelPool[conn].available < p.channelPool[conn].poolSize {
		p.channelPool[conn].pool <- channel
		p.channelPool[conn].available++
		fmt.Printf("放回一个rmq channel剩余:%v 总共:%v \n",p.channelPool[conn].available,p.channelPool[conn].poolSize)
	} else {
		_ = channel.Close()
	}
}

func RmqBasicPublish(exchange string,routeKey string,msg []byte)  (bool,string){
	conn,channel,_ := RabbitMQConnectionPoolPtr.GetChannel()
	if channel == nil{
		return false,"channel error"
	}
	defer RabbitMQConnectionPoolPtr.ReleaseChannel(conn,channel)

	// 启用 Confirm 模式
	if err := channel.Confirm(false); err != nil {
		log.Fatalf("Failed to enable publisher confirms: %v", err)
		return false,"Failed to enable publisher confirms: " + err.Error()
	}

	// 启用 Confirm 模式后，可以监听 Confirm 的通道
	confirms := channel.NotifyPublish(make(chan amqp.Confirmation, 1))

	body := amqp.Publishing{
		ContentType:"text/plain",
		Body:msg,
	}
	err  := channel.Publish(exchange,routeKey,false,false,body)
	if err != nil{
		log.Fatalf("Failed to publish message: %v", err)
		return false,err.Error()
	}

	// 等待确认
	if confirmed := <-confirms; confirmed.Ack {
		fmt.Println("Rabbitmq Message confirmed")
		return true,""
	} else {
		log.Println("Failed to confirm message")
		return false,"Failed to confirm message"
	}
}



