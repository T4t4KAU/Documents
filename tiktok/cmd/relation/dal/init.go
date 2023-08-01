package dal

import (
	"douyin/cmd/relation/dal/db"
	"douyin/cmd/relation/dal/redis"
)

func Init() {
	db.Init()
	redis.Init()
}
