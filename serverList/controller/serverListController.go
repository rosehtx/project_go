package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"serverList/config"
	"serverList/enum"
	"serverList/service"
	"strconv"
)

type ServerListReturnData struct {
	CommonReturnData
	OtherData []*service.ServerList `json:"data"`
}

func (returnData ServerListReturnData) AddAndUpdateServer(c *gin.Context) {
	returnData   = initReturnData()
	serverId, _ := strconv.Atoi(c.DefaultQuery("serverId", "0"))
	serverType, _ := strconv.Atoi(c.DefaultQuery("type", "0"))
	ip := c.Query("ip")
	port, _ := strconv.Atoi(c.DefaultQuery("port", "0"))

	if serverId == 0 || serverType == 0 || ip == "" || port == 0 {
		returnData.Status = enum.STATUS_FAIL
		returnData.Msg 	  = config.ParamError
		c.JSON(http.StatusOK, returnData)
		return
	}

	service.AddAndUpdateServerList(serverId, serverType, ip, port,true)

	returnData.OtherData = service.OtherData
	c.JSON(http.StatusOK, returnData)

}

func (returnData ServerListReturnData) GetList(c *gin.Context) {
	returnData   = initReturnData()
	returnData.OtherData = service.OtherData
	c.JSON(http.StatusOK, returnData)
}

func initReturnData() ServerListReturnData {
	return ServerListReturnData{
		CommonReturnData:CommonReturnData{
			enum.STATUS_SUCC,
			config.Success,
		},
		OtherData:[]*service.ServerList{},
	}
}
