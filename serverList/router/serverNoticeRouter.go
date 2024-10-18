package router

import (
	"github.com/gin-gonic/gin"
	"serverList/controller"
)

func initServerNoticeRouter(e *gin.Engine) {
	route := controller.InitNoticeReturnData()
	e.GET("/notice/getNotice", route.GetNotice)
	e.GET("/notice/endNotice", route.EndNotice)
}
