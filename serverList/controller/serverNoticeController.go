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
	serverId, _ := strconv.Atoi(c.DefaultQuery("serverId", "0"))
	if serverId == 0 {
		noticeReturnData.Status = enum.STATUS_FAIL
		noticeReturnData.Msg 	  = config.ParamError
		c.JSON(http.StatusOK, noticeReturnData)
		return
	}
	//获取公告
	service.NoticeOperationChan <- service.NoticeOperation{
		Action 		: enum.MAP_GET,
		ServerId 	: serverId,
	}
	resNotice := <-service.NoticeOperationResultChan
	noticeReturnData.OtherData.ServerId = serverId
	noticeReturnData.OtherData.Notice   = resNotice.Notice
	c.JSON(http.StatusOK, noticeReturnData)
}

//直接结束公告
func (noticeReturnData ServerNoticeReturnData) EndNotice (c *gin.Context){
	serverId, _ := strconv.Atoi(c.DefaultQuery("serverId", "0"))
	if serverId == 0 {
		noticeReturnData.Status	= enum.STATUS_FAIL
		noticeReturnData.Msg  	= config.ParamError
		c.JSON(http.StatusOK, noticeReturnData)
		return
	}
	service.EndServerNotice(serverId)
	noticeReturnData.OtherData  = nil
	c.JSON(http.StatusOK, noticeReturnData)
}

func InitNoticeReturnData() ServerNoticeReturnData {
	return ServerNoticeReturnData{
		CommonReturnData:CommonReturnData{
			enum.STATUS_SUCC,
			config.Success,
		},
		OtherData:&service.Notice{},
	}
}
