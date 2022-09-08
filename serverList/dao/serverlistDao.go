package dao

import (
	"gorm.io/gorm"
	"serverList/model"
)

type ServerList struct {
	model.ServerList
}

func (s ServerList) GetAllServerListData() (*gorm.DB, *[]model.ServerList) {
	var allServer []model.ServerList
	result := model.Db.Find(&allServer)
	return result, &allServer
}
