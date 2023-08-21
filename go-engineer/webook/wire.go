//go:build wireinject

package main

import (
	"gitee.com/geekbang/basic-go/webook/internal/repository"
	"gitee.com/geekbang/basic-go/webook/internal/repository/cache"
	"gitee.com/geekbang/basic-go/webook/internal/repository/dao"
	"gitee.com/geekbang/basic-go/webook/internal/service"
	"gitee.com/geekbang/basic-go/webook/internal/web"
	"gitee.com/geekbang/basic-go/webook/ioc"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

func InitWebServer() *gin.Engine {
	wire.Build(
		ioc.InitRedis, ioc.InitDB,

		// DAO 部分
		dao.NewGORMUserDAO,

		// Cache 部分
		cache.NewRedisUserCache, cache.NewRedisCodeCache,

		// repository 部分
		repository.NewCachedUserRepository,
		repository.NewCachedCodeRepository,

		// service 部分
		ioc.InitSmsService,
		service.NewSMSCodeService,
		service.NewUserService,

		// handler 部分
		web.NewUserHandler,

		// gin 的中间件
		ioc.GinMiddlewares,

		// Web 服务器
		ioc.InitWebServer,
	)
	// 随便返回一个
	return gin.Default()
}
