package middleware

import (
	"encoding/gob"
	"github.com/ecodeclub/ekit/set"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type LoginMiddlewareBuilder struct {
	publicPaths set.Set[string]
}

func NewLoginMiddlewareBuilder() *LoginMiddlewareBuilder {
	s := set.NewMapSet[string](3)
	s.Add("/users/signup")
	s.Add("/users/login_sms/code/send")
	s.Add("/users/login_sms")
	s.Add("/users/login")
	return &LoginMiddlewareBuilder{
		publicPaths: s,
	}
}

func (l *LoginMiddlewareBuilder) CheckLogin() gin.HandlerFunc {
	gob.Register(time.Time{})
	return func(ctx *gin.Context) {
		// 不需要校验
		if l.publicPaths.Exist(ctx.Request.URL.Path) {
			return
		}
		sess := sessions.Default(ctx)
		// 验证一下就可以
		if sess.Get("userId") == nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		//ctx.Next()
		const timeKey = "update_time"
		val := sess.Get(timeKey)
		updateTime, ok := val.(time.Time)
		// 处于演示效果，整个 session 的过期时间是 1 分钟，所以我这里十秒钟刷新一次。
		// val == nil 是说明刚登录成功
		// 我们不在登录里面初始化这个 update_time，是因为它属于"刷新"机制，而不属于登录机制
		if val == nil || (ok && time.Now().Sub(updateTime) > time.Second*10) {
			sess.Options(sessions.Options{
				MaxAge: 60,
			})
			sess.Set(timeKey, time.Now())
			if err := sess.Save(); err != nil {
				panic(err)
			}
		}
	}
}
