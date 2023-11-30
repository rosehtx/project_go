package rmqConsume

import (
	"errors"
	"fmt"
	"reflect"
	"testing"
)

var RabbitmqConsumerName = []string{
	"ServerData",
}
type TestData struct {
	QueueName string
}

func (sData *TestData)  BasicConsumer(queueName string) (bool,error){
	return true,errors.New("wwwwww")
}

//var RabbitmqConsumerNameSlice = []RmqConsumer{
//	&ServerData{},
//}

func TestConsume(t *testing.T)  {
	for _,consumer  := range RabbitmqConsumerName{
		op  		:= reflect.ValueOf(consumer)
		consumerVal := op.MethodByName("BasicConsumer")
		aa := consumerVal.Call([]reflect.Value{
			reflect.ValueOf("2222"),
		})
		fmt.Println(op,aa)
	}

	//for _,consumer := range RabbitmqConsumerNameSlice{
	//	//a,b := consumer.BasicConsumer("aaa")
	//	//fmt.Println(a,b)
	//	op  		:= reflect.ValueOf(consumer)
	//	consumerVal := op.MethodByName("BasicConsumer")
	//	aa := consumerVal.Call([]reflect.Value{
	//		reflect.ValueOf("2222"),
	//	})
	//	fmt.Println(op,aa)
	//}
}
