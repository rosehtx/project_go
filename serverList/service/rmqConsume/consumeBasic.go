package rmqConsume

type RmqConsumer interface {
	BasicConsumer(queueName string) (bool,error)//用来处理不同业务场景的逻辑
}

//用来初始化消费者
func init()  {
	//for _,consumer := range config.RabbitmqConsumerNameSlice{
	//	op  		:= reflect.ValueOf(consumer)
	//	consumerVal := op.MethodByName("BasicConsumer")
	//	consumerVal.Call()
	//}
}