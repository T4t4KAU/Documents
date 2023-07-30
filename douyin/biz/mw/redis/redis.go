package redis

import (
	"douyin/pkg/constants"
	"github.com/go-redis/redis/v7"
	"time"
)

var (
	expireTime  = time.Hour * 1
	rdbFollows  *redis.Client
	rdbFavorite *redis.Client
)

func Init() {
	rdbFollows = redis.NewClient(&redis.Options{
		Addr:     constants.RedisAddr,
		Password: constants.RedisPassword,
		DB:       0,
	})
	rdbFavorite = redis.NewClient(&redis.Options{
		Addr:     constants.RedisAddr,
		Password: constants.RedisPassword,
		DB:       1,
	})
}
