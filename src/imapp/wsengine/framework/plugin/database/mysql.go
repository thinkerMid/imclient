package database

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"sync"
	"time"
	"ws/framework/env"
)

var (
	masterDB *gorm.DB
	once     sync.Once
)

// MasterDB .
func MasterDB() *gorm.DB {
	once.Do(func() {
		config := env.NacosConfig.MysqlDataBase

		dial := mysql.Open(config.DSN())

		gormConfig := gorm.Config{
			Logger:                 newDBLoggerProxy(),
			SkipDefaultTransaction: true,
			QueryFields:            true,
		}

		db, err := gorm.Open(dial, &gormConfig)
		if err != nil {
			panic("fail to connect database,err:" + err.Error())
		}

		sqlDB, err := db.DB()
		if err != nil {
			panic("fail to create database,err:" + err.Error())
		}

		sqlDB.SetMaxOpenConns(100)
		sqlDB.SetMaxIdleConns(10)
		sqlDB.SetConnMaxIdleTime(20 * time.Second)

		masterDB = db
	})

	return masterDB
}
