package router

import (
	"github.com/gin-gonic/gin"
	"reflect"
	"serverList/controller"
)

func initServerListRouter(e *gin.Engine) {
	routeServer := controller.ServerListReturnData{}
	//(1)注册方法一  利用反射
	Routers = append(Routers, Router{
		path:       "/server/addOrUpdateServer", //路由
		httpMethod: "get",                        //http方法 get post
		MethodName: "AddOrUpdateServerList",      //方法名
		Controller: reflect.ValueOf(routeServer), //方法
	})

	//(2)注册方法二  直接注册不过这边有新的controller话不停的添加
	e.GET("/server/getList", routeServer.GetList)
	//e.POST("/server/initServerList", routeServer.InitServerList)
}
