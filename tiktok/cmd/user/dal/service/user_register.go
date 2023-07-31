package service

import (
	"context"
	"tiktok/cmd/user/dal/db"
	"tiktok/kitex_gen/user"
	"tiktok/pkg/constants"
	"tiktok/pkg/errno"
	"tiktok/utils"
)

type UserRegisterService struct {
	ctx context.Context
}

func NewUserRegisterService(ctx context.Context) *UserRegisterService {
	return &UserRegisterService{ctx: ctx}
}

func (s *UserRegisterService) UserRegister(req *user.UserRegisterRequest) (int64, error) {
	u, err := db.QueryUserByName(s.ctx, req.Username)
	if err != nil {
		return int64(0), err
	}
	if *u != (db.User{}) {
		return int64(0), errno.UserAlreadyExistErr
	}

	hashedPassword, _ := utils.EncryptPassword(req.Password)
	uid, err := db.CreateUser(s.ctx, &db.User{
		UserName:        req.Username,
		Password:        hashedPassword,
		Avatar:          constants.TestAva,
		BackgroundImage: constants.TestBackground,
	})

	return uid, err
}
