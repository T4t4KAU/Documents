package service

import (
	"context"
	"gitee.com/geekbang/basic-go/webook/internal/domain"
	"gitee.com/geekbang/basic-go/webook/internal/repository"
	repomocks "gitee.com/geekbang/basic-go/webook/internal/repository/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"
	"testing"
	"time"
)

func TestUserService_Login(t *testing.T) {
	//固定使用一个时间
	ctime := time.Now()
	testCases := []struct {
		name string

		mock func(ctrl *gomock.Controller) repository.UserRepository

		// 输入
		ctx      context.Context
		email    string
		password string

		// 预期中的输出
		wantErr  error
		wantUser domain.User
	}{
		{
			name: "登录成功",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				// 这边演示一下不用 gomock.Any
				repo.EXPECT().
					FindByEmail(context.Background(), "123@qq.com").
					Return(domain.User{
						Id:    123,
						Email: "123@qq.com",
						// 这里你要用 bcrypt 生成一个合法的密码
						Password: "$2a$10$s51GBcU20dkNUVTpUAQqpe6febjXkRYvhEwa5OkN5rU6rw2KTbNUi",
						Phone:    "15261890000",
						Ctime:    ctime,
					}, nil)
				return repo
			},
			ctx:   context.Background(),
			email: "123@qq.com",
			// 这是原始的密码。然后你用这个密码调用 bcrypt 生成一个加密后的密码
			password: "hello#world123",
			// 这边这个返回的是，实际上就是在 mock 中返回的
			wantUser: domain.User{
				Id:       123,
				Email:    "123@qq.com",
				Password: "$2a$10$s51GBcU20dkNUVTpUAQqpe6febjXkRYvhEwa5OkN5rU6rw2KTbNUi",
				Phone:    "15261890000",
				Ctime:    ctime,
			},
		},
		{
			name: "用户未找到",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().
					FindByEmail(context.Background(), "123@qq.com").
					// 在这里，模拟返回 ErrUserNotFound 错误
					Return(domain.User{}, repository.ErrUserNotFound)
				return repo
			},
			ctx:      context.Background(),
			email:    "123@qq.com",
			password: "hello#world123",
			// 返回密码错误
			wantErr: ErrInvalidUserOrPassword,
		},
		{
			name: "密码错误",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				// 这边演示一下不用 gomock.Any
				repo.EXPECT().
					FindByEmail(context.Background(), "123@qq.com").
					Return(domain.User{
						Id:    123,
						Email: "123@qq.com",
						// 这里你要用 bcrypt 生成一个合法的密码
						Password: "$2a$10$s51GBcU20dkNUVTpUAQqpe6febjXkRYvhEwa5OkN5rU6rw2KTbNUi",
						Phone:    "15261890000",
						Ctime:    ctime,
					}, nil)
				return repo
			},
			ctx:   context.Background(),
			email: "123@qq.com",
			// 用的是 hello#world123 加密后的密码
			// 这里我们用一个错误的密码
			password: "hello#world",
			// 返回密码错误
			wantErr: ErrInvalidUserOrPassword,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			repo := tc.mock(ctrl)
			svc := NewUserService(repo)
			user, err := svc.Login(tc.ctx, tc.email, tc.password)
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantUser, user)
		})
	}
}

func TestPasswordEncrypt(t *testing.T) {
	pwd := []byte("hello#world123")
	// 加密
	encrypted, err := bcrypt.GenerateFromPassword(pwd, bcrypt.DefaultCost)
	// 比较
	println(string(encrypted))
	err = bcrypt.CompareHashAndPassword(encrypted, pwd)
	require.NoError(t, err)
}
