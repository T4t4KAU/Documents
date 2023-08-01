package auth

import (
	"context"
	"douyin/cmd/api/handlers"
	"douyin/cmd/api/rpc"
	"douyin/kitex_gen/user"
	"douyin/pkg/constants"
	"douyin/pkg/errno"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/hertz-contrib/jwt"
	"time"
)

var MW *jwt.HertzJWTMiddleware

func Init() {
	// 创建jwt middleware
	MW, _ = jwt.New(&jwt.HertzJWTMiddleware{
		Key:        []byte(constants.SecretKey),
		Timeout:    time.Hour * 24,
		MaxRefresh: time.Hour * 24,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(int64); ok {
				return jwt.MapClaims{
					constants.IdentityKey: v,
				}
			}
			return jwt.MapClaims{}
		},
		HTTPStatusMessageFunc: func(e error, ctx context.Context, c *app.RequestContext) string {
			switch e.(type) {
			case errno.ErrNo:
				return e.(errno.ErrNo).ErrMsg
			default:
				return e.Error()
			}
		},
		LoginResponse: func(ctx context.Context, c *app.RequestContext, code int, token string, expire time.Time) {
			c.JSON(consts.StatusOK, map[string]interface{}{
				"status_code": errno.SuccessCode,
				"status_msg":  errno.SuccessMsg,
				"token":       token,
			})
		},
		Unauthorized: func(ctx context.Context, c *app.RequestContext, code int, message string) {
			c.JSON(code, map[string]interface{}{
				"status_code": errno.AuthorizationFailedErrCode,
				"status_msg":  message,
			})
		},
		Authenticator: func(ctx context.Context, c *app.RequestContext) (interface{}, error) {
			var loginVar handlers.UserParam

			loginVar.UserName = c.Query("username")
			loginVar.PassWord = c.Query("password")

			if len(loginVar.UserName) == 0 || len(loginVar.PassWord) == 0 {
				return "", jwt.ErrMissingLoginValues
			}

			return rpc.UserLogin(context.Background(),
				&user.UserLoginRequest{Username: loginVar.UserName, Password: loginVar.PassWord})
		},
		IdentityKey: constants.IdentityKey,

		TokenLookup:   "header: Authorization, query: token, cookie: jwt",
		TokenHeadName: "Bearer",
		TimeFunc:      time.Now,
	})
}
