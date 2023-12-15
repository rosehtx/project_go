package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"serverList/config"
	"serverList/enum"
	"serverList/service"
	"sync"
	"time"
)

type RmqReturnData struct {
	CommonReturnData
}

func (returnData RmqReturnData) RmqPublishMessage(c *gin.Context) {
	msg := c.DefaultQuery("msg", "")
	if msg == "" {
		returnData.Status = enum.STATUS_FAIL
		returnData.Msg 	  = config.ParamError
		c.JSON(http.StatusOK, returnData)
		return
	}

	byteMsg := []byte(msg)
	service.RmqBasicPublish("odoo_billIn","odoo_billIn",byteMsg)
	returnData.Status = enum.STATUS_SUCC
	returnData.Msg = "success"
	c.JSON(http.StatusOK, returnData)
}

func (returnData RmqReturnData) TestRmq(c *gin.Context) {
	// 模拟多个任务需要进行 RabbitMQ 操作 默认启动5个连接的话会有5个获取不到连接
	taskCount := 20
	var wg sync.WaitGroup
	defer func() {
		fmt.Println("nice ok l")
	}()


	//生成一个测试service.RabbitMQConnectionPoolPtr.GetChannel并发的方法
	for i := 0; i < taskCount; i++ {
		wg.Add(1)
		go func(c int) {
			defer wg.Done()
			// 从连接池获取连接
			conn, channel,err := service.RabbitMQConnectionPoolPtr.GetChannel()
			if err != nil {
				fmt.Printf("Failed to get RabbitMQ connection from pool: %v \n", err)
				return
			}
			defer service.RabbitMQConnectionPoolPtr.ReleaseChannel(conn,channel)
			time.Sleep(time.Duration(2)*time.Second)
		}(i)
	}

	// 等待所有任务完成
	wg.Wait()
	returnData.Status = enum.STATUS_SUCC
	returnData.Msg = "success"
	c.JSON(http.StatusOK, returnData)
}

