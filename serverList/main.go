package main

import (
	"fmt"
	"serverList/config"
	"serverList/model"
	"serverList/router"
	"serverList/service"
)

func main() {
	//建立mysql链接
	_,sqlErr := model.InitSqlPool()
	if sqlErr != nil{
		fmt.Println(sqlErr.Error())
		return
	}

	//初始化serverlist数据
	resInitServer 	:= service.InitServerList()
	if resInitServer != nil {
		fmt.Println(resInitServer.Error())
		return
	}

	//初始化公告数据
	//resInitNotice, resInitNoticeMsg := service.InitServerNotice()
	//if resInitNotice == false {
	//	fmt.Println(resInitNoticeMsg)
	//	return
	//}

	//初始化rmq
	_ , rmqError := service.NewRabbitMQConnectionPool(config.RMQ_CON_NUM)
	if rmqError != nil{
		fmt.Println("start server error" + rmqError.Error())
		return
	}
	//初始化rmq消费者
	for queueName, _ := range config.RabbitmqBasicConsumer {
		for i := 0; i < config.RMQ_CONSUME_NUM; i++ {
			go service.BasicConsumer(queueName)
		}
	}

	//初始化路由
	r 	:= router.InitRouter()
	//r.POST("/gin/test",test)
	err := r.Run(":8090")
	if err != nil{
		fmt.Println("start server error" + err.Error())
		return
	}

	// 利用通道读取的阻塞来执行上面协程
	//done := make(chan bool)
	//<-done
	select {}
}
