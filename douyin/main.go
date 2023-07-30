package main

import (
	"context"
	"douyin/biz/dal"
	"douyin/biz/handler"
	"douyin/biz/mw/jwt"
	"douyin/biz/router"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/hertz-contrib/reverseproxy"
)

// customizeRegister registers customize routers.
func customizedRegister(r *server.Hertz) {
	r.GET("/ping", handler.Ping)

	// your code ...
}

// register registers all routers.
func register(r *server.Hertz) {

	router.GeneratedRegister(r)

	customizedRegister(r)
}

// Set up /src/*name route forwarding to access minio from external network
func minioReverseProxy(c context.Context, ctx *app.RequestContext) {
	proxy, _ := reverseproxy.NewSingleHostReverseProxy("http://localhost:18001")
	ctx.URI().SetPath(ctx.Param("name"))
	hlog.CtxInfof(c, string(ctx.Request.URI().Path()))
	proxy.ServeHTTP(c, ctx)
}

func main() {
	dal.Init()
	jwt.Init()

	h := server.Default(
		server.WithStreamBody(true),
		server.WithHostPorts("0.0.0.0:18005"),
	)

	h.GET("/src/*name", minioReverseProxy)
	register(h)

	h.Spin()
}
