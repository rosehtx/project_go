package dao

import (
	"gorm.io/gorm"
	"serverList/model"
)

// GetAllServerListData 获取所有的serverList数据
func GetAllServerListData() (*gorm.DB, []model.ServerList) {
	var allServer []model.ServerList
	result := model.Db.Find(&allServer)
	return result, allServer
}

//InsertServerListData 新增serverList数据
func InsertServerListData(serverId int,ip string,port int,serverType int,status int) (*gorm.DB) {
	result := model.Db.Create(&model.ServerList{
		ServerId:serverId,
		Ip:ip,
		Port:port,
		Type:serverType,
		Status:status,
	})
	return result
}

//UpdateServerListData 更新serverList数据
func UpdateServerListData(serverId int,ip string,port int,serverType int,status int) (*gorm.DB) {
	result := model.Db.Model(&model.ServerList{}).
		Where("server_id = ? and type = ?", serverId,serverType).
		Updates(model.ServerList{Ip: ip,Port: port,Status: status})
	return result
}