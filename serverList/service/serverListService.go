package service

import (
	"fmt"
	daoServerList "serverList/dao"
	"sync"
)

//server配置
type Server struct {
	ServerId int  `json:"serverId"`
	Type     int  `json:"type"`
	Ip     string `json:"ip"`
	Port   int    `json:"port"`
	Status int    `json:"status"`
}
//server配置集合
type ServerList struct {
	mu 		 sync.Mutex
	Server   []Server
}

var ServerListPtr *ServerList

//更新serverlist
func (s *ServerList) AddOrUpdateServerList(serverId int, serverType int, ip string, port int,status int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i,ss := range s.Server{
		//存在更新
		if ss.ServerId == serverId && ss.Type == serverType {
			err := daoServerList.UpdateServerListData(serverId,ip,port,serverType,status)
			if err != nil {
				return
			}
			s.Server[i].Ip 		= ip
			s.Server[i].Port 	= port
			s.Server[i].Status 	= status
			return
		}
	}
	//没有则新增
	err := daoServerList.InsertServerListData(serverId,ip,port,serverType,status)
	if err != nil {
		return
	}
	s.Server = append(s.Server,Server{
		ServerId:serverId,
		Type:serverType,
		Ip:ip,
		Port:port,
		//Status:status,
	})
}

//初始化serverlist
func InitServerList() error {
	result, ss := daoServerList.GetAllServerListData()
	if result != nil {
		return result
	}

	//server slice
	var sSlice []Server
	for _,modelServer := range ss{
		sSlice = append(sSlice,Server{
			ServerId: modelServer.ServerId,
			Type:modelServer.Type,
			Ip:modelServer.Ip,
			Port:modelServer.Port,
			//Status:modelServer.Status,
		})
	}
	fmt.Println("获取serverList结果:",sSlice)
	ServerListPtr = &ServerList{
		Server:sSlice,
	}
	return nil
}


