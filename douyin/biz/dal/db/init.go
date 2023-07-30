package db

import (
	"douyin/pkg/constants"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var dbConn *gorm.DB

func Init() {
	var err error
	dbConn, err = gorm.Open(mysql.Open(constants.MySQLDSN),
		&gorm.Config{
			PrepareStmt:            true,
			SkipDefaultTransaction: true,
		})
	if err != nil {
		panic(err)
	}
	if !dbConn.Migrator().HasTable(&User{}) {
		err = dbConn.Migrator().CreateTable(&User{})
		if err != nil {
			panic(err)
		}
	}

	if !dbConn.Migrator().HasTable(&Video{}) {
		err = dbConn.Migrator().CreateTable(&Video{})
		if err != nil {
			panic(err)
		}
	}

	if !dbConn.Migrator().HasTable(&Comment{}) {
		err = dbConn.Migrator().CreateTable(&Comment{})
		if err != nil {
			panic(err)
		}
	}

	if !dbConn.Migrator().HasTable(&Messages{}) {
		err = dbConn.Migrator().CreateTable(&Messages{})
		if err != nil {
			panic(err)
		}
	}

	if !dbConn.Migrator().HasTable(&Favorites{}) {
		err = dbConn.Migrator().CreateTable(&Favorites{})
		if err != nil {
			panic(err)
		}
	}

	if !dbConn.Migrator().HasTable(&Follows{}) {
		err = dbConn.Migrator().CreateTable(&Follows{})
		if err != nil {
			panic(err)
		}
	}
}
