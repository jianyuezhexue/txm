package util

import (
	"database/sql"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitGorm() *gorm.DB {
	dsn := "root:root@tcp(127.0.0.1:3306)/txm"

	// 连接&设置连接池参数
	sqlDB, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Second * 28800) // SHOW VARIABLES LIKE '%timeout%';

	// 生成gorm连接
	gormDB, err := gorm.Open(mysql.New(mysql.Config{Conn: sqlDB}))
	if err != nil {
		panic(err)
	}
	return gormDB
}
