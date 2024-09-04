package utils

import (
	"fmt"
	"github.com/opentracing/opentracing-go"
	"io"
	"serverList/config"
	//"github.com/jaegertracing/jaeger-client-go/config"
	jaegerConfig "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
	"sync"
)

//traver单例
//var tracerOne sync.Once

type tracerStruct struct {
	Tracer opentracing.Tracer
	tracerIo io.Closer
}
//连接池结构体
type JaegerPool struct {
	muLock sync.Mutex
	tracerPool chan *tracerStruct
	poolSize int
	poolNum int
}
//用来记录jaeger连接池
var JaegerPoolPtr *JaegerPool
//记录jaeger配置
var cfg jaegerConfig.Configuration

func NewJaegerPool(poolNum int) error{
	// Jaeger 配置
	cfg = jaegerConfig.Configuration{
		ServiceName: config.JAEGER_SERVICE_NAME,
		Sampler: &jaegerConfig.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &jaegerConfig.ReporterConfig{
			LogSpans: true,
			LocalAgentHostPort: config.JAEGER_IP_PORT,
		},
	}

	//先初始化一个JaegerPoolPtr
	JaegerPoolPtr = &JaegerPool{
		tracerPool: make(chan *tracerStruct, poolNum),// 创建 Jaeger 追踪器/初始化 Jaeger Tracer
		poolSize: config.JAEGER_POOL_SIZE,
		poolNum: 0,
	}

	err := AddJaegerPoolTracer(poolNum)
	return err
}

func createJaegerTracer()  (opentracing.Tracer,io.Closer,error){
	tracer, tracerIo, err := cfg.NewTracer(jaegerConfig.Logger(jaegerlog.StdLogger))
	return tracer,tracerIo,err
}

func AddJaegerPoolTracer(num int)  error{
	for i:= 0 ;i < num ;i++ {
		tracer, tracerIo, err := createJaegerTracer()
		if err != nil{
			return err
		}
		opentracing.SetGlobalTracer(tracer)
		tStruct				:= &tracerStruct{}
		tStruct.Tracer 		= tracer
		tStruct.tracerIo 	= tracerIo
		// 记得在程序结束时关闭 tracer
		//defer tracerIo.Close()
		JaegerPoolPtr.tracerPool <- tStruct
		JaegerPoolPtr.poolNum++
	}
	return nil
}

func CreateTracerLog(span string,event string,tag string, msg string) error {
	tStruct,err := GetTracer()
	if err != nil{
		return err
	}
	defer ReleaseTracer(tStruct)
	//这边进行日志记录
	//在需要追踪的代码段中创建 span
	tracerSpan := tStruct.Tracer.StartSpan(span)
	defer tracerSpan.Finish()

	if tag != ""{
		tracerSpan.SetTag(event, tag)
	}

	//日志(Logs)，日志也定义为名值对。用于捕获调试信息，或者相关Span的相关信息
	tracerSpan.LogKV(event, msg)
	return nil
}

func GetTracer() (*tracerStruct,error) {
	JaegerPoolPtr.muLock.Lock()
	defer JaegerPoolPtr.muLock.Unlock()

	select {
		case tStruct := <-JaegerPoolPtr.tracerPool:
			JaegerPoolPtr.poolNum--
			fmt.Printf("获取一个jaeger tracer:%v \n",JaegerPoolPtr.poolNum)
			return tStruct, nil
		default:
			tracer, tracerIo, err := createJaegerTracer()
			if err != nil{
				return nil, err
			}
			fmt.Printf("新建一个jaeger tracer \n")
			tStruct	:= &tracerStruct{
				tracer,
				tracerIo,
			}
			return tStruct, nil
	}
}

func ReleaseTracer(tStruct *tracerStruct)  {
	JaegerPoolPtr.muLock.Lock()
	defer JaegerPoolPtr.muLock.Unlock()

	//如果没有到达连接池的上限则放回去
	//多了则直接关闭
	if JaegerPoolPtr.poolNum < JaegerPoolPtr.poolSize{
		JaegerPoolPtr.tracerPool <- tStruct
		JaegerPoolPtr.poolNum ++
	}else{
		_= tStruct.tracerIo.Close()
	}
	fmt.Printf("放回一个jaeger tracer:%v \n",JaegerPoolPtr.poolNum)
}
