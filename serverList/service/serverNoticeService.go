package service

import (
	"fmt"
	"gorm.io/gorm"
	daoNotice "serverList/dao"
	"serverList/enum"
	"serverList/model"
	"time"
)

type Notice struct {
	ServerId int  `json:"serverId"`
	Notice string `json:"notice"`
}

var ServerNoticeMap = make(map[int]string)

var AllNotice []model.ServerNotice

func InitServerNotice() (bool, string) {
	var result *gorm.DB
	result, AllNotice = daoNotice.GetAllNotEndServerNoticeData()
	if result.Error != nil {
		return false, result.Error.Error()
	}
	fmt.Println("start init notice")
	if AllNotice != nil{
		for i := 0; i < len(AllNotice); i++ {
			//先注册到map里
			ServerNoticeMap[AllNotice[i].ServerId] = AllNotice[i].Notice
			go func(notice model.ServerNotice) {
				tickerDb := time.NewTicker(2 * time.Second)
				for  {
					select {
					case <-tickerDb.C:
						now := time.Now().Unix()
						if notice.EndTime < uint64(now) || notice.IsEnd == enum.NOTICE_IS_END_YES{
							delete(ServerNoticeMap,notice.ServerId)
							goto endNotice
						}
					default:
					}
				}
			endNotice:
				fmt.Printf("server : %v notice is end",notice.ServerId)
			}(AllNotice[i])
		}
	}
	fmt.Println("end init notice")
	return true, ""
}


