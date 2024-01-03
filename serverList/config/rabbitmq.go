package config

import (
	"serverList/service/rmqConsume"
)

type RmqConsumer interface {
	BasicConsumer(msg []byte) (bool,error)//用来处理不同业务场景的逻辑
}

//RabbitmqBasicConsumer 基础消费 消费的队列名称与rmq消费的业务逻辑是一一对应的
var RabbitmqBasicConsumer = map[string]RmqConsumer{
	"odoo_billIn": &rmqConsume.ServerData{},
}


