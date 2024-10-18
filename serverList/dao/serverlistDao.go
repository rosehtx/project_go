package dao

import (
	"serverList/model"
)

//获取所有的serverList数据
func GetAllServerListData() (error, []model.ServerList) {
	var allServer []model.ServerList
	PoolChan,err := model.GetMysqlConn()
	if err != nil{
		return err,allServer
	}
	defer model.ReleaseMysqlConn(PoolChan)
	result := PoolChan.Db.Find(&allServer)
	return result.Error, allServer
}

//新增serverList数据
func InsertServerListData(serverId int,ip string,port int,serverType int,status int) error {
	PoolChan,err := model.GetMysqlConn()
	if err != nil{
		return err
	}
	defer model.ReleaseMysqlConn(PoolChan)
	result := PoolChan.Db.Create(&model.ServerList{
		ServerId:serverId,
		Ip:ip,
		Port:port,
		Type:serverType,
		Status:status,
	})
	return result.Error
}

//UpdateServerListData 更新serverList数据
func UpdateServerListData(serverId int,ip string,port int,serverType int,status int) error {
	PoolChan,err := model.GetMysqlConn()
	if err != nil{
		return err
	}
	defer model.ReleaseMysqlConn(PoolChan)
	result := PoolChan.Db.Model(&model.ServerList{}).
		Where("server_id = ? and type = ?", serverId,serverType).
		Updates(model.ServerList{Ip: ip,Port: port,Status: status})
	return result.Error
}