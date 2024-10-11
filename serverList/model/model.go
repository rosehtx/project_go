package model

import (
	"context"
	"database/sql"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"serverList/config"
	"strconv"
	"sync"
	"time"
)

type PoolChan struct{
	Db *gorm.DB //mysql连接
	Cancel context.CancelFunc //context取消方法
}
//mysql连接池
type SqlPool struct {
	pool      chan *PoolChan
	poolSize  int
	available int//可用的连接数
	mu        sync.Mutex
}
var SqlPoolPtr *SqlPool

//初始化mysql连接池
func InitSqlPool() error{
	pool  		:= make(chan *PoolChan,config.MysqlConNum)
	SqlPoolPtr 	 = &SqlPool{
		pool:      pool,
		poolSize:  config.MysqlConNum,
		available: config.MysqlConNum,
	}
	for i := 0; i < config.MysqlConNum; i++ {
		fmt.Println("----init mysql start "+strconv.Itoa(i)+"----")
		Db,cancel,createError := CreateSqlCon()
		if createError != nil{
			return createError
		}
		poolChanStruct := &PoolChan{
			Db,
			cancel,
		}
		SqlPoolPtr.pool <- poolChanStruct
	}
	return nil
}

func CreateSqlCon()  (*gorm.DB,context.CancelFunc,error){
	dsn 	 := config.MysqlUser + ":" + config.MysqlPass + "@tcp(" + config.MysqlIp + ":" + strconv.Itoa(config.MysqlPort) + ")/serverlist?charset=utf8mb4&parseTime=True&loc=Local"
	Db, err  := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil{
		return nil,nil,err
	}
	sqlDB, _ 	:= Db.DB()
	// 创建一个可以手动取消的上下文
	ctx, cancel := context.WithCancel(context.Background())
	go pingDb(sqlDB,ctx)
	return Db,cancel,nil
}

//获取mysql连接
func GetMysqlConn() (*PoolChan,error) {
	SqlPoolPtr.mu.Lock()
	defer SqlPoolPtr.mu.Unlock()

	if SqlPoolPtr.available == 0 {
		//return nil,fmt.Errorf("mysql连接池缺失连接")
		Db,cancel,createError := CreateSqlCon()
		if createError != nil{
			return nil,createError
		}
		var returnPoolChanStruct *PoolChan
		returnPoolChanStruct.Db 	= Db
		returnPoolChanStruct.Cancel = cancel
		return returnPoolChanStruct,nil
	}
	SqlPoolPtr.available--
	return <-SqlPoolPtr.pool,nil
}

//放回mysql连接
func ReleaseMysqlConn(poolChan *PoolChan) {
	SqlPoolPtr.mu.Lock()
	defer SqlPoolPtr.mu.Unlock()

	//蜀国数量超出了设定的上限则直接删除即可
	if SqlPoolPtr.available >= config.MysqlConNum {
		fmt.Println("mysql连接数超过数量")
		poolChan.Cancel()
		return
	}
	SqlPoolPtr.available++
	SqlPoolPtr.pool <- poolChan
}

//ping mysql
func pingDb(sqlDB *sql.DB,ctx context.Context) {
	tickerDb := time.NewTicker(20 * time.Second)
	for {
		select {
		case <-tickerDb.C:
			//fmt.Println("mysql ping")
			_ = sqlDB.Ping()
		case <-ctx.Done():
			fmt.Printf("一个pingDb退出...\n")
			return
		default:
		}
	}
}
