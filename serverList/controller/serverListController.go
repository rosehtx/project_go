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
	OtherData []service.Server `json:"data"`
}

func (returnData *ServerListReturnData) AddOrUpdateServerList(c *gin.Context) {
	serverId, _ := strconv.Atoi(c.DefaultQuery("serverId", "0"))
	serverType, _ := strconv.Atoi(c.DefaultQuery("type", "0"))
	ip := c.Query("ip")
	port, _ := strconv.Atoi(c.DefaultQuery("port", "0"))
	status, _ := strconv.Atoi(c.DefaultQuery("status", "2"))

	if serverId == 0 || serverType == 0 || ip == "" || port == 0 {
		returnData.Status = enum.STATUS_FAIL
		returnData.Msg 	  = config.ParamError
		c.JSON(http.StatusOK, returnData)
		return
	}

	service.ServerListPtr.AddOrUpdateServerList(serverId, serverType, ip, port,status)

	returnData.OtherData = service.ServerListPtr.Server
	c.JSON(http.StatusOK, returnData)
}

func (returnData *ServerListReturnData) GetList(c *gin.Context) {
	returnData.OtherData = service.ServerListPtr.Server
	c.JSON(http.StatusOK, returnData)
}

func (returnData *ServerListReturnData) Test(c *gin.Context) {
	//defer utils.CreateTracerLog("Server","func","thisIsTag","desc")
	c.JSON(http.StatusOK, returnData)
}

func InitReturnData() *ServerListReturnData{
	return &ServerListReturnData{
		CommonReturnData:CommonReturnData{
			enum.STATUS_SUCC,
			config.Success,
		},
		OtherData:[]service.Server{},
	}
}
