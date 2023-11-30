package config

import (
	"serverList/service/rmqConsume"
)

//与rmq的消费者名称一一匹配 用来注册,作为消费后的转发处理
var RabbitmqConsumerNameSlice = []rmqConsume.RmqConsumer{

}


