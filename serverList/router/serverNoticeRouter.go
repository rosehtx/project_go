package router

import (
	"github.com/gin-gonic/gin"
	"serverList/controller"
)

func initServerNoticeRouter(e *gin.Engine) {
	route := controller.ServerNoticeReturnData{}
	e.GET("/notice/getNotice", route.GetNotice)
}
