package jwt

import (
	"context"
	"douyin/biz/dal/db"
	"douyin/biz/model/basic/user"
	"douyin/pkg/errno"
	"douyin/pkg/utils"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/hertz-contrib/jwt"
	"time"
)

var (
	MW       *jwt.HertzJWTMiddleware
	identity = "user_id"
)

func Init() {
	MW, _ = jwt.New(&jwt.HertzJWTMiddleware{
		Key:         []byte("tiktok secret key"),
		TokenLookup: "query:token,form:token",
		Timeout:     24 * time.Hour,
		IdentityKey: identity,

		IdentityHandler: func(ctx context.Context, c *app.RequestContext) any {
			claims := jwt.ExtractClaims(ctx, c)
			return claims[identity]
		},

		// Verify password at login
		Authenticator: func(ctx context.Context, c *app.RequestContext) (any, error) {
			var loginRequest user.DouyinUserLoginRequest
			if err := c.BindAndValidate(&loginRequest); err != nil {
				return nil, err
			}
			u, err := db.QueryUserByName(loginRequest.Username)

			if ok := utils.VerifyPassword(loginRequest.Password, u.Password); !ok {
				err = errno.PasswordIsNotVerified
				return nil, err
			}
			if err != nil {
				return nil, err
			}

			c.Set("user_id", u.ID)
			return u.ID, nil
		},
		// Set the payload in the token
		PayloadFunc: func(data any) jwt.MapClaims {
			if v, ok := data.(int64); ok {
				return jwt.MapClaims{
					identity: v,
				}
			}
			return jwt.MapClaims{}
		},
		// build login response if verify password successfully
		LoginResponse: func(ctx context.Context, c *app.RequestContext, code int, token string, expire time.Time) {
			hlog.CtxInfof(ctx, "Login success ，token is issued clientIP: "+c.ClientIP())
			c.Set("token", token)
		},
		// Verify token and get the id of logged-in user
		Authorizator: func(data any, ctx context.Context, c *app.RequestContext) bool {
			c.Set("current_user_id", 1014)
			if v, ok := data.(int64); ok {
				c.Set("current_user_id", v)
				hlog.CtxInfof(ctx, "Token is verified clientIP: "+c.ClientIP())
				return true
			}
			return true
		},
		// Validation failed, build the message
		Unauthorized: func(ctx context.Context, c *app.RequestContext, code int, message string) {
			c.JSON(consts.StatusOK, user.DouyinUserLoginResponse{
				StatusCode: errno.AuthorizationFailedErrCode,
				StatusMsg:  message,
			})
		},
		HTTPStatusMessageFunc: func(e error, ctx context.Context, c *app.RequestContext) string {
			resp := utils.BuildBaseResp(e)
			return resp.StatusMsg
		},
	})
}
