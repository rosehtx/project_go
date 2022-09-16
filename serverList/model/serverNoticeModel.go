package model

type ServerNotice struct {
	Id       	int64  `gorm:"primaryKey"`
	ServerId 	int    `gorm:"column:server_id"`
	Notice   	string `gorm:"column:notice"`
	StartTime   uint64 `gorm:"column:start_time"`
	EndTime     uint64 `gorm:"column:end_time"`
}

func (ServerNotice) TableName() string {
	return "server_notice"
}
