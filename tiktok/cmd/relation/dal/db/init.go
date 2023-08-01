package db

import (
	"douyin/pkg/constants"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormopentracing "gorm.io/plugin/opentracing"
)

var dbConn *gorm.DB

func Init() {
	var err error

	// 打开数据库连接
	dbConn, err = gorm.Open(mysql.Open(constants.MySQLDSN), &gorm.Config{
		PrepareStmt:            true,
		SkipDefaultTransaction: true,
	})
	if err != nil {
		panic(err)
	}

	err = dbConn.Use(gormopentracing.New())
	if err != nil {
		panic(err)
	}

	// 创建数据库表
	if !dbConn.Migrator().HasTable(&Relation{}) {
		err = dbConn.Migrator().CreateTable(&Relation{})
		if err != nil {
			panic(err)
		}
	}
}
