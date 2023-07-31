package main

import (
	"context"
	"tiktok/cmd/user/dal/pack"
	"tiktok/cmd/user/dal/service"
	user "tiktok/kitex_gen/user"
	"tiktok/pkg/errno"
)

// UserServiceImpl implements the last service interface defined in the IDL.
type UserServiceImpl struct{}

// UserRegister implements the UserServiceImpl interface.
func (s *UserServiceImpl) UserRegister(ctx context.Context, req *user.UserRegisterRequest) (resp *user.UserRegisterResponse, err error) {
	resp = new(user.UserRegisterResponse)

	if len(req.Username) == 0 || len(req.Password) == 0 {
		r := pack.BuildBaseResp(errno.ParamErr)
		resp.StatusCode = r.StatusCode
		resp.StatusMsg = r.StatusMsg
		return
	}

	uid, err := service.NewUserRegisterService(ctx).UserRegister(req)
	r := pack.BuildBaseResp(err)
	resp.StatusCode = r.StatusCode
	resp.StatusMsg = r.StatusMsg
	resp.UserId = uid

	return
}

// UserLogin implements the UserServiceImpl interface.
func (s *UserServiceImpl) UserLogin(ctx context.Context, req *user.UserLoginRequest) (resp *user.UserLoginResponse, err error) {
	resp = new(user.UserLoginResponse)

	if len(req.Username) == 0 || len(req.Password) == 0 {
		r := pack.BuildBaseResp(errno.ParamErr)
		resp.StatusCode = r.StatusCode
		resp.StatusMsg = r.StatusMsg
		return
	}

	uid, err := service.NewUserLoginService(ctx).UserLogin(req)
	r := pack.BuildBaseResp(err)
	resp.StatusCode = r.StatusCode
	resp.StatusMsg = r.StatusMsg
	resp.UserId = uid

	return
}
