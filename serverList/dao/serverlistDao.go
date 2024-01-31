package dao

import (
	"serverList/model"
)

//获取所有的serverList数据
func GetAllServerListData() (error, []model.ServerList) {
	var allServer []model.ServerList
	Db,err := model.SqlPoolPtr.GetMysqlConn()
	if err != nil{
		return err,allServer
	}
	defer model.SqlPoolPtr.ReleaseMysqlConn(Db)
	result := Db.Find(&allServer)
	return result.Error, allServer
}

//新增serverList数据
func InsertServerListData(serverId int,ip string,port int,serverType int,status int) error {
	Db,err := model.SqlPoolPtr.GetMysqlConn()
	if err != nil{
		return err
	}
	defer model.SqlPoolPtr.ReleaseMysqlConn(Db)
	result := Db.Create(&model.ServerList{
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
	Db,err := model.SqlPoolPtr.GetMysqlConn()
	if err != nil{
		return err
	}
	defer model.SqlPoolPtr.ReleaseMysqlConn(Db)
	result := Db.Model(&model.ServerList{}).
		Where("server_id = ? and type = ?", serverId,serverType).
		Updates(model.ServerList{Ip: ip,Port: port,Status: status})
	return result.Error
}