package handlers

import (
	"douyin/kitex_gen/relation"
	"douyin/kitex_gen/user"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

// SendRegisterResponse 发送注册响应信息
func SendRegisterResponse(c *app.RequestContext, resp user.UserRegisterResponse) {
	c.JSON(consts.StatusOK, resp)
}

// SendLoginResponse 发送登录响应信息
func SendLoginResponse(c *app.RequestContext, resp user.UserLoginResponse) {
	c.JSON(consts.StatusOK, resp)
}

func SendRelationActionResponse(c *app.RequestContext, resp relation.RelationActionResponse) {
	c.JSON(consts.StatusOK, resp)
}
