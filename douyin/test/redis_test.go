package test

import (
	"douyin/pkg/constants"
	"github.com/go-redis/redis/v7"
)

func GetRedisClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     constants.RedisAddr,
		Password: "",
		DB:       0,
	})
}
