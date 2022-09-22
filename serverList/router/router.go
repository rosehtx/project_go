package router

import (
	"github.com/gin-gonic/gin"
	"reflect"
)

//路由结构体
type Router struct {
	path 		string 			//路由
	httpMethod 	string 			//http方法 get post
	MethodName 	string 			//方法名
	Controller 	reflect.Value	//控制器
	//Args   []reflect.Type 	//参数类型
}

var Routers = []Router{}

/**
下面用了俩种方式来注册路由
(1) 利用反射往上面的Router结构体里写入路由数据后用下面的Bind方法统一处理(处理Routers切片里多个Router)
(2) 去路由方法里直接注册 这个方便很多推荐
 */
func InitRouter() *gin.Engine  {
	//初始化路由
	r := gin.Default()
	//这边去注册路由
	Bind(r)
	return r
}

func Bind(e *gin.Engine) {
	initServerListRouter(e)  //注册serverList的router
	initServerNoticeRouter(e)  //注册serverList的router
	//(1)注册方法一  利用反射
	for _, route := range Routers {
		if(route.httpMethod == "get"){
			meth 		:= route.Controller.MethodByName(route.MethodName)
			handlerFunc := match(meth)
			e.GET(route.path,handlerFunc)
		}
		//if(route.httpMethod == "POST"){
		//	e.POST(route.path,match(route.path,route))
		//}
	}

}

func match(Method reflect.Value) gin.HandlerFunc {
	return func(c *gin.Context) {
		arguments 	:= make([]reflect.Value, 1)
		arguments[0] = reflect.ValueOf(c) // *gin.Context
		Method.Call(arguments)
	}
}