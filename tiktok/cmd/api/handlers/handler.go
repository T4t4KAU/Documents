package handlers

import (
	"context"
	"douyin/cmd/api/rpc"
	"douyin/kitex_gen/relation"
	"douyin/kitex_gen/user"
	"douyin/pkg/constants"
	"douyin/pkg/errno"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/hertz-contrib/jwt"
	"strconv"
)

// RegisterHandler 用户注册
func RegisterHandler(ctx context.Context, c *app.RequestContext) {
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

func RelationActionHandler(ctx context.Context, c *app.RequestContext) {
	var relationVar RelationParam

	claims := jwt.ExtractClaims(ctx, c)
	userId := int64(claims[constants.IdentityKey].(float64))

	uid, err := strconv.Atoi(c.Query("to_user_id"))
	if err != nil {
		SendRelationActionResponse(c, relation.RelationActionResponse{
			StatusCode: errno.ParamErrCode,
			StatusMsg:  errno.ParamErrMsg,
		})
	}

	relationVar.ToUserId = int64(uid)

	action, err := strconv.Atoi(c.Query("action_type"))
	if err != nil {
		SendRelationActionResponse(c, relation.RelationActionResponse{
			StatusCode: errno.ParamErrCode,
			StatusMsg:  errno.ParamErrMsg,
		})
	}

	relationVar.ActionType = int32(action)

	_, err = rpc.RelationAction(ctx, &relation.RelationActionRequest{
		CurrentUserId: userId,
		ToUserId:      relationVar.ToUserId,
		ActionType:    int32(action),
	})
	if err != nil {
		SendRelationActionResponse(c, relation.RelationActionResponse{
			StatusCode: errno.ServiceErrCode,
			StatusMsg:  err.Error(),
		})
		return
	}

	SendRelationActionResponse(c, relation.RelationActionResponse{
		StatusCode: errno.SuccessCode,
		StatusMsg:  errno.SuccessMsg,
	})
}

func TestHandler(ctx context.Context, c *app.RequestContext) {

}
