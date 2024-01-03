package rmqConsume

import (
	"encoding/json"
	"fmt"
)

type ServerData struct {

}

type RmqServerData struct {
	ServerId int `json:"ServerId"`
	ServerName string `json:"ServerName"`
}

func (sData *ServerData)  BasicConsumer(msg []byte) (bool,error){
	data 	:= &RmqServerData{}

	err  	:= json.Unmarshal(msg, &data)
	if err != nil {
		fmt.Println(err)
		return  false,err
	}

	fmt.Println("ServerData.BasicConsumer")
	fmt.Println(data)

	return  true,nil
}