package dao

import (
	"serverList/enum"
	"serverList/model"
	"time"
)

func GetAllNotEndServerNoticeData() (error,[]*model.ServerNotice) {
	var allNotice []*model.ServerNotice
	PoolChan,err := model.GetMysqlConn()
	if err != nil{
		return err,nil
	}
	defer model.ReleaseMysqlConn(PoolChan)
	result := PoolChan.Db.Where("is_end = ? and end_time > ?", enum.NOTICE_IS_END_NO,time.Now().Unix()).
		Find(&allNotice)
	return result.Error,allNotice
}

