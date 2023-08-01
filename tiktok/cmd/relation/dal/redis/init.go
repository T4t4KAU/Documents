package redis

import (
	"douyin/pkg/constants"
	"github.com/go-redis/redis/v8"
	"time"
)

var (
	rdbFollows *redis.Client
	expireTime = 1 * time.Hour
)

func Init() {
	rdbFollows = redis.NewClient(&redis.Options{
		Addr:     constants.RedisAddr,
		Password: constants.RedisPassword,
		DB:       0,
	})
}
