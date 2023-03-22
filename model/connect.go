package model

import (
	"fmt"
	"gorm.io/driver/mysql"
	_ "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
)

var db *gorm.DB

func init() {
	mysqlString := os.Getenv("mysql")
	if len(mysqlString) == 0 {
		mysqlString = "127.0.0.1"
	}
	var err error
	dialector := mysql.Open(fmt.Sprintf("root:12345678@tcp(%s:3306)/subscribe2clash?charset=utf8mb4&parseTime=True&loc=Local", mysqlString))
	db, err = gorm.Open(dialector)
	if err != nil {
		fmt.Println("连接数据库失败，请检查参数：", err)
		return
	}
	db1, _ := db.DB()
	db1.SetMaxIdleConns(10)
	db1.SetMaxOpenConns(100)

}
