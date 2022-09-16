package service

import (
	daoNotice "serverList/dao"
)

func InitServerNotice() (bool, string) {
	result, ss := daoNotice.GetAllServerNoticeData()
	if result.Error != nil {
		return false, result.Error.Error()
	}

	for i := 0; i < len(*ss); i++ {
		AddAndUpdateServerList((*ss)[i].ServerId, (*ss)[i].Type, (*ss)[i].Ip, (*ss)[i].Port,(*ss)[i].Status,false)
	}
	return true, ""
}


