package dao

import (
	"gorm.io/gorm"
	"serverList/model"
)

func GetAllServerNoticeData() (*gorm.DB, *[]model.ServerNotice) {
	var allNotice []model.ServerNotice
	result := model.Db.Find(&allNotice)
	return result, &allNotice
}

