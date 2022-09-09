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

var Db *gorm.DB

func InitModel(ch chan string, wg *sync.WaitGroup) {
	fmt.Println("----init mysql start----")
	//链接mysql
	dsn := config.MysqlUser + ":" + config.MysqlPass + "@tcp(" + config.MysqlIp + ":" + strconv.Itoa(config.MysqlPort) + ")/serverlist?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	Db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})

	defer wg.Done()

	if err != nil {
		ch <- err.Error()
		return
	}
	sqlDB, errDb := Db.DB()
	if errDb != nil {
		ch <- errDb.Error()
		return
	}
	fmt.Println("----ticker ping mysql start----")
	go pingDb(sqlDB)
	fmt.Println("----ticker ping mysql end----")

	fmt.Println("----init mysql end----")
}

func pingDb(sqlDB *sql.DB) {
	tickerDb := time.NewTicker(20 * time.Second)
	for {
		select {
		case <-tickerDb.C:
			fmt.Println("mysql ping")
			sqlDB.Ping()
		default:
		}
	}
}
