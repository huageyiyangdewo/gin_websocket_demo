package model

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	//_ "github.com/jinzhu/gorm/dialects/mysql"
	 _ "github.com/go-sql-driver/mysql"
	logging "github.com/sirupsen/logrus"

	"time"
)

var DB *gorm.DB


func ConnectMysql(conn string)  {
	fmt.Println(conn)
	db, err := gorm.Open("mysql", conn)
	if err != nil {
		logging.Fatalf("gorm open mysql failed, err:%s \n", err)
		return
	}

	logging.Println("connect mysql success")


	// 是否开启打印 sql 语句
	db.LogMode(true)

	if gin.Mode() == "release" {
		db.LogMode(false)
	}

	db.SingularTable(true)  // 表默认不加复数 s
	db.DB().SetMaxIdleConns(20) // 设置连接池，空闲
	db.DB().SetMaxOpenConns(200)  // 设置最大的连接数
	db.DB().SetConnMaxLifetime(time.Second*30)
	DB = db

	migration()
	logging.Println("mysql migration success")

}

func migration()  {
	// 自动迁移模式
	DB.Set("gorm:table_options", "charset=utf8mb4").AutoMigrate(&User{})
}