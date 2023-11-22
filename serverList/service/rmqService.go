package service

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"serverList/config"
	"strconv"
)

//var channel *amqp.Channel
var conn *amqp.Connection

func init() {
	fmt.Println("init rmq start")
	var err error
	conn, err = amqp.Dial("amqp://"+config.RMQ_USER+":"+config.RMQ_PASS+"@"+config.RMQ_IP+":"+strconv.Itoa(config.RMQ_PORT)+"/"+config.RMQ_VHOST)
	if err != nil{
		fmt.Println(err.Error())
		return
	}
	fmt.Println("init rmq end")
}

func getChannel()  *amqp.Channel{
	// 创建一个通道
	//var channelErr error
	channel, channelErr := conn.Channel()
	if channelErr != nil{
		fmt.Println(channelErr.Error())
		return nil
	}
	return channel
}

func RmqBasicPublish(exchange string,routeKey string,msg []byte)  (bool,string){
	channel := getChannel()
	if channel == nil{
		return false,"channel error"
	}
	defer channel.Close()

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



