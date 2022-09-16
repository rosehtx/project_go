package main

import (
	"fmt"
	"serverList/model"
	"serverList/service"
	"sync"
)

func main() {
	//建立mysql链接
	var mysqlChanWg sync.WaitGroup
	checkMysqlChan := make(chan string, 1)
	mysqlChanWg.Add(1)
	go model.InitModel(checkMysqlChan, &mysqlChanWg)
	mysqlChanWg.Wait()

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

	//起http服务1
	go startServer()

	// 利用通道读取的阻塞来执行上面协程
	done := make(chan bool)
	<-done
}
