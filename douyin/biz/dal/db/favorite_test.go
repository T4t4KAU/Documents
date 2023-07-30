package db

import (
	"douyin/biz/mw/redis"
	"fmt"
	"testing"
)

func TestAddNewFavorite(t *testing.T) {
	Init()
	redis.Init()

	_, err := AddNewFavorite(&Favorites{
		UserId:  1000,
		VideoId: 115,
	})
	if err != nil {
		t.Errorf(err.Error())
	}

	fmt.Println("add new favorite successfully")
}
