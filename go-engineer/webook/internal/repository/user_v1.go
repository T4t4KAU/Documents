package repository

import (
	"gitee.com/geekbang/basic-go/webook/internal/repository/cache"
	"gitee.com/geekbang/basic-go/webook/internal/repository/dao"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// 用来演示非依赖注入写法

type DBConfig struct {
	DSN string
}

type CacheConfig struct {
	Type string
	// Redis 配置
	Addr string
	// 本地缓存配置
}

// NewUserRepositoryV1 非依赖注入的写法
func NewUserRepositoryV1(dbCfg DBConfig, c CacheConfig) *CachedUserRepository {
	db, err := gorm.Open(mysql.Open(dbCfg.DSN))
	if err != nil {
		panic(err)
	}
	ud := dao.NewGORMUserDAO(db)
	uc := cache.NewRedisUserCache(redis.NewClient(&redis.Options{
		Addr: c.Addr,
	}))
	return &CachedUserRepository{
		dao:   ud,
		cache: uc,
	}
}
