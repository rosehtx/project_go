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
	pool      chan *amqp.Connection
	mu        sync.Mutex
	poolSize  int
	available int
}
var RabbitMQConnectionPoolPtr *RabbitMQConnectionPool

func GetRabbitMQConnectionPoolPtr()  *RabbitMQConnectionPool{
	return RabbitMQConnectionPoolPtr
}

//建立连接池
func NewRabbitMQConnectionPool(poolSize int) (*RabbitMQConnectionPool, error) {
	pool := make(chan *amqp.Connection, poolSize)
	for i := 0; i < poolSize; i++ {
		fmt.Println("创建第"+strconv.Itoa(i+1)+"个rmq连接")
		conn, err := amqp.Dial("amqp://"+config.RMQ_USER+":"+config.RMQ_PASS+"@"+config.RMQ_IP+":"+strconv.Itoa(config.RMQ_PORT)+"/"+config.RMQ_VHOST)
		if err != nil {
			return nil, fmt.Errorf("Failed to create RabbitMQ connection: %v", err)
		}
		pool <- conn
	}

	//获取连接池使用
	RabbitMQConnectionPoolPtr = &RabbitMQConnectionPool{
		pool:      pool,
		poolSize:  poolSize,
		available: poolSize,
	}
	return RabbitMQConnectionPoolPtr, nil
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
			return nil, fmt.Errorf("没有可用的连接")
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

func getChannel()  (*amqp.Channel,*amqp.Connection){
	//获取一个连接
	poolPtr 	:= GetRabbitMQConnectionPoolPtr()
	conn ,err 	:= poolPtr.GetConnection()
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
	channel,conn := getChannel()
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



