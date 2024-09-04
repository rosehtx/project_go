package router

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
	"serverList/utils"
)

func InitJaegerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("jaeger中间件链路追踪开始")
		tStruct,_		:= utils.GetTracer()
		tracer 			:= tStruct.Tracer
		//放回tracer到池子里
		defer utils.ReleaseTracer(tStruct)
		// 从请求上下文中提取 trace 信息，如果不存在则创建一个新的 span
		var span opentracing.Span
		wireContext, _ := tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(c.Request.Header))
		if wireContext != nil {
			span = tracer.StartSpan(c.Request.URL.Path, ext.RPCServerOption(wireContext))
		} else {
			span = tracer.StartSpan(c.Request.URL.Path)
		}
		defer span.Finish()

		// 将 span 注入到请求上下文中
		c.Set("tracer", tracer)
		c.Set("parentSpanContext", span.Context())

		// 记录请求 URL 和参数
		span.SetTag("http.method", c.Request.Method)
		span.SetTag("http.url", c.Request.URL.Path)
		span.LogFields(
			log.String("query_params", c.Request.URL.RawQuery),
			log.String("form_params", c.Request.Form.Encode()),
		)

		// 继续处理请求
		c.Next()

		// 记录响应状态码
		statusCode := c.Writer.Status()
		span.SetTag("http.status_code", statusCode)
		if statusCode >= 400 {
			span.SetTag("error", true)
		}
	}
}
