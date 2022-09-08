package model

type ServerList struct {
	Id       int64  `gorm:"primaryKey"`
	ServerId int    `gorm:"column:server_id"`
	Ip       string `gorm:"column:ip"`
	Port     int    `gorm:"column:port"`
	Type     int    `gorm:"column:type"`
	UpdateAt int64  `gorm:"column:update_at" gorm:"autoUpdateTime"`
}

func (ServerList) TableName() string {
	return "server_list"
}
