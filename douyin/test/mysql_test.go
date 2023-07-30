package test

import (
	"douyin/pkg/constants"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"testing"
)

func GetMySQLDB() (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(constants.MySQLDSN), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}

func TestGetMySQLDB(t *testing.T) {
	_, err := GetMySQLDB()
	if err != nil {
		t.Errorf("MySQL connection is not alive: %s", err)
	}
}
