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
		fmt.Printf("%v",i)
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

// GetConnection 从连接池获取一个 RabbitMQ 连接
func (p *RabbitMQConnectionPool) GetConnection() (*amqp.Connection, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	select {
		case conn := <-p.pool:
			p.available--
			fmt.Printf("获取一个rmq连接剩余:%v \n",p.available)
			return conn, nil
		default:
			//return nil, fmt.Errorf("没有可用的连接")
			//这边可以临时生成，放回池子如果超出数量则直接删除

	}
}

// ReleaseConnection 将连接放回连接池
func (p *RabbitMQConnectionPool) ReleaseConnection(conn *amqp.Connection) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.available < p.poolSize {
		p.pool <- conn
		p.available++
		fmt.Printf("放回一个rmq连接剩余:%v \n",p.available)
	} else {
		_ = conn.Close()
	}
}

// GetChannel 从连接池获取一个 channel 连接
func (p *RabbitMQChannelPool) GetChannel() (*amqp.Channel, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	select {
	case conn := <-p.pool:
		p.available--
		fmt.Printf("获取一个rmq channel:%v \n",p.available)
		return conn, nil
	default:
		return nil, fmt.Errorf("没有可用的channel")
	}
}

// ReleaseChannel 将连接放回连接池
func (p *RabbitMQChannelPool) ReleaseChannel(channel *amqp.Channel) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.available < p.poolSize {
		p.pool <- channel
		p.available++
		fmt.Printf("放回一个rmq channel剩余:%v \n",p.available)
	} else {
		_ = channel.Close()
	}
}

//这边获取了channel之后立马吧 connect放回池子里，便于其他场景获取
func GetChannel()  (*amqp.Channel,*amqp.Connection){
	//获取一个连接
	conn ,err 	:= RabbitMQConnectionPoolPtr.GetConnection()
	if conn == nil{
		fmt.Println("rmq连接池获取连接失败:" + err.Error())
		return nil,nil
	}
	// 创建一个通道
	//var channelErr error
	channel, channelErr := conn.Channel()
	if channelErr != nil{
		fmt.Println(channelErr.Error())
		poolPtr.ReleaseConnection(conn)
		return nil,nil
	}
	return channel,conn
}

func RmqBasicPublish(exchange string,routeKey string,msg []byte)  (bool,string){
	channel,conn := GetChannel()
	if channel == nil{
		return false,"channel error"
	}
	defer channel.Close()
	//放回连接池
	poolPtr 	:= GetRabbitMQConnectionPoolPtr()
	defer poolPtr.ReleaseConnection(conn)

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
		fmt.Println("Rabbitmq Message confirmed \n")
		return true,""
	} else {
		log.Println("Failed to confirm message\n")
		return false,"Failed to confirm message"
	}

}



