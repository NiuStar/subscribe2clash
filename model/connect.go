package model

import (
	"fmt"
	"gorm.io/driver/mysql"
	_ "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

func init() {
	var err error
	dialector := mysql.Open("root:12345678@tcp(127.0.0.1:3306)/subscribe2clash?charset=utf8mb4&parseTime=True&loc=Local")
	db, err = gorm.Open(dialector)
	if err != nil {
		fmt.Println("连接数据库失败，请检查参数：", err)
		return
	}
	db1, _ := db.DB()
	db1.SetMaxIdleConns(10)
	db1.SetMaxOpenConns(100)

}
