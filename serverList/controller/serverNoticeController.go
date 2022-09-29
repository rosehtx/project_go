package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"serverList/config"
	"serverList/enum"
	"serverList/service"
	"strconv"
)

type ServerNoticeReturnData struct {
	CommonReturnData
	OtherData *service.Notice `json:"data"`
}

func (noticeReturnData ServerNoticeReturnData) GetNotice (c *gin.Context){
	noticeReturnData   = initNoticeReturnData()
	serverId, _ := strconv.Atoi(c.DefaultQuery("serverId", "0"))
	if serverId == 0 {
		noticeReturnData.Status = enum.STATUS_FAIL
		noticeReturnData.Msg 	  = config.ParamError
		c.JSON(http.StatusOK, noticeReturnData)
		return
	}
	noticeReturnData.OtherData.ServerId = serverId
	notice, checkNotice := service.ServerNoticeMap[serverId]
	if checkNotice != false{
		noticeReturnData.OtherData.Notice   = notice
	}
	c.JSON(http.StatusOK, noticeReturnData)
}

//直接结束公告
func (noticeReturnData ServerNoticeReturnData) EndNotice (c *gin.Context){
	noticeReturnData   = initNoticeReturnData()
	serverId, _ := strconv.Atoi(c.DefaultQuery("serverId", "0"))
	if serverId == 0 {
		noticeReturnData.Status = enum.STATUS_FAIL
		noticeReturnData.Msg 	  = config.ParamError
		c.JSON(http.StatusOK, noticeReturnData)
		return
	}
	service.EndServerNotice(serverId)
	c.JSON(http.StatusOK, noticeReturnData)
}

func initNoticeReturnData() ServerNoticeReturnData {
	return ServerNoticeReturnData{
		CommonReturnData:CommonReturnData{
			enum.STATUS_SUCC,
			config.Success,
		},
		OtherData:&service.Notice{},
	}
}
