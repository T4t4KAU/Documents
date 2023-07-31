package handlers

import (
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"tiktok/kitex_gen/user"
)

// SendRegisterResponse 发送注册响应信息
func SendRegisterResponse(c *app.RequestContext, resp user.UserRegisterResponse) {
	c.JSON(consts.StatusOK, resp)
}

// SendLoginResponse 发送登录响应信息
func SendLoginResponse(c *app.RequestContext, resp user.UserLoginResponse) {
	c.JSON(consts.StatusOK, resp)
}
