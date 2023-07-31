package handlers

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"tiktok/cmd/api/rpc"
	"tiktok/kitex_gen/user"
	"tiktok/pkg/errno"
)

// Register 用户注册
func Register(ctx context.Context, c *app.RequestContext) {
	var registerVar UserParam

	registerVar.UserName = c.Query("username")
	registerVar.PassWord = c.Query("password")

	if len(registerVar.UserName) == 0 || len(registerVar.PassWord) == 0 {
		SendRegisterResponse(c, user.UserRegisterResponse{
			StatusCode: errno.ParamErrCode,
			StatusMsg:  errno.ParamErrMsg,
		})

		return
	}

	// 使用注册rpc
	resp, _ := rpc.UserRegister(context.Background(), &user.UserRegisterRequest{
		Username: registerVar.UserName,
		Password: registerVar.PassWord,
	})

	SendRegisterResponse(c, *resp)
}
