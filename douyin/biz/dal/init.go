package dal

import (
	"douyin/biz/dal/db"
	"douyin/biz/mw/redis"
)

func Init() {
	db.Init()
	redis.Init()
}
