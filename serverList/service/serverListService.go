package service

import (
	daoServerList "serverList/dao"
)

type Server struct {
	Type int    `json:"type"`
	Ip   string `json:"ip"`
	Port int    `json:"port"`
}
type ServerList struct {
	ServerId int      `json:"serverId"`
	Server   []Server `json:"server"`
}

//用来存放server和type的map集合
var serverTypeMap = make(map[int]map[int]struct{})

var OtherServerListData []*ServerList

func AddAndUpdateServerList(serverId int, serverType int, ip string, port int,isOperateMysql bool) {
	//判断serverid是否已经定义
	allServerTypeMap, checkServerId := serverTypeMap[serverId]
	_, checkServerType := allServerTypeMap[serverType]
	//server都还没有定义
	if checkServerId == false {
		//先存放集合
		serverTypeMap[serverId] = map[int]struct{}{serverType: {}}
		server := Server{
			serverType,
			ip,
			port,
		}
		serverList := &ServerList{
			serverId,
			[]Server{server},
		}
		OtherServerListData = append(OtherServerListData, serverList)
		//这边对数据库进行新增
		if isOperateMysql == true {
			daoServerList.InsertServerListData(serverId,ip,port,serverType)
		}
	} else {
		//server已有但是 type没有
		if checkServerType == false {
			for i, v := range OtherServerListData {
				if v.ServerId == serverId {
					OtherServerListData[i].Server = append(OtherServerListData[i].Server, Server{serverType,
						ip,
						port})
					//存放type
					serverTypeMap[serverId][serverType] = struct{}{}
					//这边对数据库进行新增
					if isOperateMysql == true {
						daoServerList.InsertServerListData(serverId,ip,port,serverType)
					}
					goto cancelFor
				}
			}
		} else {
			//都有进进行更新
			for i, v := range OtherServerListData {
				if v.ServerId == serverId {
					for severI, severV := range v.Server {
						if severV.Type == serverType {
							OtherServerListData[i].Server[severI].Ip = ip
							OtherServerListData[i].Server[severI].Port = port
							//这边对数据库进行更新
							if isOperateMysql == true {
								daoServerList.UpdateServerListData(serverId,ip,port,serverType)
							}
							goto cancelFor
						}
					}
				}
			}
		}
	}
cancelFor:
}

func InitServerList() (bool, string) {
	result, ss := daoServerList.GetAllServerListData()
	if result.Error != nil {
		return false, result.Error.Error()
	}

	for i := 0; i < len(*ss); i++ {
		AddAndUpdateServerList((*ss)[i].ServerId, (*ss)[i].Type, (*ss)[i].Ip, (*ss)[i].Port,false)
	}
	return true, ""
}


