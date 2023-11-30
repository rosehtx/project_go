package rmqConsume

import "errors"

type ServerData struct {

}

func (sData *ServerData)  BasicConsumer(queueName string) (bool,error){
	return true,errors.New("wwwwww")
}