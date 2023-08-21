package middleware

import (
	"gitee.com/geekbang/basic-go/webook/internal/web"
	"github.com/ecodeclub/ekit/set"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"net/http"
	"strings"
	"time"
)

type JWTLoginMiddlewareBuilder struct {
	publicPaths set.Set[string]
}

func NewJWTLoginMiddlewareBuilder() *JWTLoginMiddlewareBuilder {
	s := set.NewMapSet[string](3)
	s.Add("/users/signup")
	s.Add("/users/login_sms/code/send")
	s.Add("/users/login_sms")
	s.Add("/users/login")
	return &JWTLoginMiddlewareBuilder{
		publicPaths: s,
	}
}
func (j *JWTLoginMiddlewareBuilder) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 不需要校验
		if j.publicPaths.Exist(ctx.Request.URL.Path) {
			return
		}

		// Authorization 头部
		// 得到的格式 Bearer token
		authCode := ctx.GetHeader("Authorization")
		if authCode == "" {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		// SplitN 的意思是切割字符串，但是最多 N 段
		// 如果要是 N 为 0 或者负数，则是另外的含义，可以看它的文档
		authSegments := strings.SplitN(authCode, " ", 2)
		if len(authSegments) != 2 {
			// 格式不对
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		tokenStr := authSegments[1]
		uc := web.UserClaims{}
		token, err := jwt.ParseWithClaims(tokenStr, &uc, func(token *jwt.Token) (interface{}, error) {
			return web.JWTKey, nil
		})
		if err != nil || !token.Valid {
			// 不正确的 token
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		expireTime, err := uc.GetExpirationTime()
		if err != nil {
			// 拿不到过期时间
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if expireTime.Before(time.Now()) {
			// 已经过期
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if ctx.GetHeader("User-Agent") != uc.UserAgent {
			// 换了一个 User-Agent，可能是攻击者
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		// 每 10 秒刷新一次
		if expireTime.Sub(time.Now()) < time.Second*50 {
			uc.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Minute))
			newToken, err := token.SignedString(web.JWTKey)
			if err != nil {
				// 因为刷新这个事情，并不是一定要做的，所以这里可以考虑打印日志
				// 暂时这样打印
				log.Println(err)
			} else {
				ctx.Header("x-jwt-token", newToken)
			}

		}

		// 说明 token 是合法的
		// 我们把这个 token 里面的数据放到 ctx 里面，后面用的时候就不用再次 Parse 了
		ctx.Set("user", uc)
	}
}
