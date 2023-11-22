package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"serverList/config"
	"serverList/enum"
	"serverList/service"
)

type RmqReturnData struct {
	CommonReturnData
}

func (returnData RmqReturnData) RmqPublishMessage(c *gin.Context) {
	msg := c.DefaultQuery("msg", "")
	if msg == "" {
		returnData.Status = enum.STATUS_FAIL
		returnData.Msg 	  = config.ParamError
		c.JSON(http.StatusOK, returnData)
		return
	}

	byteMsg := []byte(msg)
	service.RmqBasicPublish("odoo_billIn","odoo_billIn",byteMsg)
	returnData.Status = enum.STATUS_SUCC
	c.JSON(http.StatusOK, returnData)
}
