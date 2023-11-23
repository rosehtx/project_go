package main

import (
	"fmt"
	"serverList/model"
	"serverList/router"
	"serverList/service"
	"sync"
)

func main() {
	//建立mysql链接
	var mysqlChanWg sync.WaitGroup
	checkMysqlChan := make(chan string, 1)
	mysqlChanWg.Add(1)
	go model.InitModel(&checkMysqlChan, &mysqlChanWg)
	mysqlChanWg.Wait()

	checkMsg := <-checkMysqlChan
	//不用了可以关闭哦
	close(checkMysqlChan)
	fmt.Println("init model result :" + checkMsg)
	if checkMsg != "success"{
		fmt.Println(checkMsg)
		return
	}

	//初始化serverlist数据
	resInitServer, resInitServerMsg := service.InitServerList()
	if resInitServer == false {
		fmt.Println(resInitServerMsg)
		return
	}

	//初始化公告数据
	resInitNotice, resInitNoticeMsg := service.InitServerNotice()
	if resInitNotice == false {
		fmt.Println(resInitNoticeMsg)
		return
	}

	//初始化rmq
	_ , rmqError := service.NewRabbitMQConnectionPool(5)
	if rmqError != nil{
		fmt.Println("start server error" + rmqError.Error())
		return
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
