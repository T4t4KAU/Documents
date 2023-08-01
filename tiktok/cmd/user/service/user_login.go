package service

import (
	"context"
	"douyin/cmd/user/dal/db"
	"douyin/kitex_gen/user"
	"douyin/pkg/errno"
	"douyin/utils"
)

type UserLoginService struct {
	ctx context.Context
}

func NewUserLoginService(ctx context.Context) *UserLoginService {
	return &UserLoginService{ctx: ctx}
}

func (s *UserLoginService) UserLogin(req *user.UserLoginRequest) (int64, error) {
	u, err := db.QueryUserByName(s.ctx, req.Username)
	if err != nil {
		return int64(0), err
	}
	if *u == (db.User{}) {
		return int64(0), errno.UserIsNotExistErr
	}

	if !utils.VerifyPassword(req.Password, u.Password) {
		return int64(0), errno.AuthorizationFailedErr
	}
	return u.ID, nil
}
