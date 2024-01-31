package model

import (
	"database/sql"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"serverList/config"
	"strconv"
	"sync"
	"time"
)

//mysql连接池
type SqlPool struct {
	pool      chan *gorm.DB
	poolSize  int
	available int//可用的连接数
	mu        sync.Mutex
}
var SqlPoolPtr *SqlPool

//初始化mysql连接池
func InitSqlPool() (*SqlPool,error){
	pool  := make(chan *gorm.DB,config.MysqlConNum)
	for i := 0; i < config.MysqlConNum; i++ {
		fmt.Println("----init mysql start "+strconv.Itoa(i)+"----")
		dsn := config.MysqlUser + ":" + config.MysqlPass + "@tcp(" + config.MysqlIp + ":" + strconv.Itoa(config.MysqlPort) + ")/serverlist?charset=utf8mb4&parseTime=True&loc=Local"
		Db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil{
			return nil,err
		}
		sqlDB, _ := Db.DB()
		go pingDb(sqlDB)
		pool <- Db
	}

	SqlPoolPtr = &SqlPool{
		pool:      pool,
		poolSize:  config.MysqlConNum,
		available: config.MysqlConNum,
	}
	return SqlPoolPtr,nil
}

//获取mysql连接
func (sql *SqlPool) GetMysqlConn() (*gorm.DB,error) {
	sql.mu.Lock()
	defer sql.mu.Unlock()
	if sql.available == 0 {
		return nil,fmt.Errorf("mysql连接池已满")
	}
	sql.available--
	return <-sql.pool,nil
}

//放回mysql连接
func (sql *SqlPool) ReleaseMysqlConn(db *gorm.DB) {
	sql.mu.Lock()
	defer sql.mu.Unlock()
	sql.available++
	sql.pool <- db
}

//ping mysql
func pingDb(sqlDB *sql.DB) {
	tickerDb := time.NewTicker(20 * time.Second)
	for {
		select {
		case <-tickerDb.C:
			//fmt.Println("mysql ping")
			sqlDB.Ping()
		default:
		}
	}
}
