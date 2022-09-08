package controller

import (
	"fmt"
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

func (returnData *ServerListReturnData) AddAndUpdateServer(c *gin.Context) {
	fmt.Println("serverList get request")
	serverId, _ := strconv.Atoi(c.DefaultQuery("serverId", "0"))
	serverType, _ := strconv.Atoi(c.DefaultQuery("type", "0"))
	ip := c.Query("ip")
	port, _ := strconv.Atoi(c.DefaultQuery("port", "0"))

	if serverId == 0 || serverType == 0 || ip == "" || port == 0 {
		c.JSON(http.StatusOK, CommonReturnData{
			enum.STATUS_FAIL, config.ParamError,
		})
		return
	}

	service.AddAndUpdateServerList(serverId, serverType, ip, port)

	returnData.OtherData = service.OtherData
	c.JSON(http.StatusOK, returnData)

}

func (returnData *ServerListReturnData) GetList(c *gin.Context) {
	returnData.CommonReturnData = CommonReturnData{
		enum.STATUS_SUCC, config.Success,
	}
	returnData.OtherData = service.OtherData
	c.JSON(http.StatusOK, returnData)
}
