// Code generated by hertz generator.

package user

import (
	"context"
	user "douyin/biz/model/basic/user"
	"douyin/biz/mw/jwt"
	service "douyin/biz/service/user"
	"douyin/pkg/errno"
	"douyin/pkg/utils"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

// User .
// @router /douyin/user/ [GET]
func User(ctx context.Context, c *app.RequestContext) {
	var err error
	var req user.DouyinUserRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		resp := utils.BuildBaseResp(err)
		c.JSON(consts.StatusOK, user.DouyinUserResponse{
			StatusCode: resp.StatusCode,
			StatusMsg:  resp.StatusMsg,
		})
		return
	}

	u, err := service.NewUserService(ctx, c).UserInfo(&req)
	resp := utils.BuildBaseResp(err)
	c.JSON(consts.StatusOK, user.DouyinUserResponse{
		StatusCode: resp.StatusCode,
		StatusMsg:  resp.StatusMsg,
		User:       u,
	})
}

// UserRegister .
// @router /douyin/user/register/ [POST]
func UserRegister(ctx context.Context, c *app.RequestContext) {
	var err error
	var req user.DouyinUserRegisterRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		resp := utils.BuildBaseResp(err)
		c.JSON(consts.StatusOK, user.DouyinUserResponse{
			StatusCode: resp.StatusCode,
			StatusMsg:  resp.StatusMsg,
		})

		return
	}

	_, err = service.NewUserService(ctx, c).UserRegister(&req)
	if err != nil {
		resp := utils.BuildBaseResp(err)
		c.JSON(consts.StatusOK, user.DouyinUserRegisterResponse{
			StatusCode: resp.StatusCode,
			StatusMsg:  resp.StatusMsg,
		})

		return
	}

	jwt.MW.LoginHandler(ctx, c)
	token := c.GetString("token")
	v, _ := c.Get("user_id")
	uid := v.(int64)

	// 返回JSON数据
	c.JSON(consts.StatusOK, user.DouyinUserRegisterResponse{
		StatusCode: errno.SuccessCode,
		StatusMsg:  errno.SuccessMsg,
		Token:      token,
		UserId:     uid,
	})
}

// UserLogin .
// @router /douyin/user/login/ [POST]
func UserLogin(ctx context.Context, c *app.RequestContext) {
	var err error
	var req user.DouyinUserLoginRequest

	err = c.BindAndValidate(&req)
	if err != nil {
		resp := utils.BuildBaseResp(err)
		c.JSON(consts.StatusOK, user.DouyinUserResponse{
			StatusCode: resp.StatusCode,
			StatusMsg:  resp.StatusMsg,
		})
		return
	}

	jwt.MW.LoginHandler(ctx, c)
	token := c.GetString("token")

	// token为空 用户验证失败
	if len(token) == 0 {
		return
	}

	v, _ := c.Get("user_id")
	uid := v.(int64)

	c.JSON(consts.StatusOK, user.DouyinUserLoginResponse{
		StatusCode: errno.SuccessCode,
		StatusMsg:  errno.SuccessMsg,
		Token:      token,
		UserId:     uid,
	})
}