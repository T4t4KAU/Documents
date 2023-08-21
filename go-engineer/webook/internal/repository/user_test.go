package repository

import (
	"context"
	"database/sql"
	"gitee.com/geekbang/basic-go/webook/internal/domain"
	"gitee.com/geekbang/basic-go/webook/internal/repository/cache"
	cachemocks "gitee.com/geekbang/basic-go/webook/internal/repository/cache/mocks"
	"gitee.com/geekbang/basic-go/webook/internal/repository/dao"
	daomocks "gitee.com/geekbang/basic-go/webook/internal/repository/dao/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
	"time"
)

func TestCachedUserRepository_FindById(t *testing.T) {
	// 因为存储的是毫秒数，也就是纳秒部分被去掉了
	// 所以我们需要利用 nowMs 来重建一个不含纳秒部分的 time.Time
	nowMs := time.Now().UnixMilli()
	now := time.UnixMilli(nowMs)
	testCases := []struct {
		name string
		// 返回 mock 的 UserDAO 和 UserCache
		mock func(ctrl *gomock.Controller) (dao.UserDAO, cache.UserCache)

		// 输入
		ctx context.Context
		id  int64

		// 预期输出
		wantUser domain.User
		wantErr  error
	}{
		{
			name: "找到了用户，未命中缓存",
			mock: func(ctrl *gomock.Controller) (dao.UserDAO, cache.UserCache) {
				d := daomocks.NewMockUserDAO(ctrl)
				c := cachemocks.NewMockUserCache(ctrl)
				// 注意这边，我们传入的是 int64，
				// 所以要做一个显式的转换，不然默认 12 是 int 类型
				c.EXPECT().Get(gomock.Any(), int64(12)).
					// 模拟缓存未命中
					Return(domain.User{}, cache.ErrKeyNotExist)
				// 模拟回写缓存
				c.EXPECT().Set(gomock.Any(), domain.User{
					Id:       12,
					Email:    "123@qq.com",
					Password: "123456",
					Phone:    "15212345678",
					Ctime:    now,
				}).Return(nil)

				d.EXPECT().FindById(gomock.Any(), int64(12)).
					Return(dao.User{
						Id: 12,
						Email: sql.NullString{
							String: "123@qq.com",
							Valid:  true,
						},
						Password: "123456",
						Phone: sql.NullString{
							String: "15212345678",
							Valid:  true,
						},
						Ctime: nowMs,
						Utime: nowMs,
					}, nil)
				return d, c
			},

			ctx: context.Background(),
			id:  12,

			wantUser: domain.User{
				Id:       12,
				Email:    "123@qq.com",
				Password: "123456",
				Phone:    "15212345678",
				Ctime:    now,
			},
		},
		{
			name: "找到了用户，直接命中缓存",
			mock: func(ctrl *gomock.Controller) (dao.UserDAO, cache.UserCache) {
				d := daomocks.NewMockUserDAO(ctrl)
				c := cachemocks.NewMockUserCache(ctrl)
				// 注意这边，我们传入的是 int64，
				// 所以要做一个显式的转换，不然默认 12 是 int 类型
				c.EXPECT().Get(gomock.Any(), int64(12)).
					// 模拟缓存命中
					Return(domain.User{
						Id:       12,
						Email:    "123@qq.com",
						Password: "123456",
						Phone:    "15212345678",
						Ctime:    now,
					}, nil)
				return d, c
			},

			ctx: context.Background(),
			id:  12,

			wantUser: domain.User{
				Id:       12,
				Email:    "123@qq.com",
				Password: "123456",
				Phone:    "15212345678",
				Ctime:    now,
			},
		},
		{
			name: "没有找到用户",
			mock: func(ctrl *gomock.Controller) (dao.UserDAO, cache.UserCache) {
				d := daomocks.NewMockUserDAO(ctrl)
				c := cachemocks.NewMockUserCache(ctrl)
				// 注意这边，我们传入的是 int64，
				// 所以要做一个显式的转换，不然默认 12 是 int 类型
				c.EXPECT().Get(gomock.Any(), int64(12)).
					// 模拟缓存命中
					Return(domain.User{}, cache.ErrKeyNotExist)
				d.EXPECT().FindById(gomock.Any(), int64(12)).
					Return(dao.User{}, dao.ErrDataNotFound)
				return d, c
			},

			ctx:     context.Background(),
			id:      12,
			wantErr: ErrUserNotFound,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			d, c := tc.mock(ctrl)
			repo := NewCachedUserRepository(d, c)
			u, err := repo.FindById(tc.ctx, tc.id)
			assert.Equal(t, tc.wantUser, u)
			assert.Equal(t, tc.wantErr, err)
		})
	}
}
