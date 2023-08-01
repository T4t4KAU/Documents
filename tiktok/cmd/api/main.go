package main

import (
	"context"
	"douyin/cmd/api/auth"
	"douyin/cmd/api/handlers"
	"douyin/cmd/api/rpc"
	"douyin/pkg/constants"
	"douyin/pkg/errno"
	tracer "douyin/pkg/trace"
	"fmt"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/middlewares/server/recovery"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

func Init() {
	tracer.InitJaeger(constants.ApiServiceName)
	rpc.InitRPC()
	auth.Init()
}

func main() {
	Init()
	r := server.New(
		server.WithHostPorts("0.0.0.0:8888"),
		server.WithHandleMethodNotAllowed(true),
	)

	r.Use(recovery.Recovery(recovery.WithRecoveryHandler(
		func(ctx context.Context, c *app.RequestContext, err interface{}, stack []byte) {
			hlog.SystemLogger().CtxErrorf(ctx, "[Recovery] err=%v\nstack=%s", err, stack)
			c.JSON(consts.StatusInternalServerError, map[string]interface{}{
				"status_code": errno.ServiceErrCode,
				"status_msg":  fmt.Sprintf("[Recovery] err=%v\nstack=%s", err, stack),
			})
		})))

	// 注册路由
	router := r.Group("/tiktok")
	userRouter := router.Group("/user")
	userRouter.POST("/register", handlers.RegisterHandler)
	userRouter.POST("/login", auth.MW.LoginHandler)

	relationRouter := router.Group("/relation", auth.MW.MiddlewareFunc())
	relationRouter.POST("/action", handlers.RelationActionHandler)

	r.NoRoute(func(ctx context.Context, c *app.RequestContext) {
		c.String(consts.StatusOK, "no route")
	})
	r.NoMethod(func(ctx context.Context, c *app.RequestContext) {
		c.String(consts.StatusOK, "no method")
	})

	r.Spin()
}
