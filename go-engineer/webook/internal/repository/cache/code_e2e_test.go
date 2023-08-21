package cache

import (
	"context"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestRedisCodeCache_Set_e2e(t *testing.T) {
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		t.Fatal(err)
	}
	// 因为我们这里不再需要使用 mock 的 Redis 客户端
	// 所以我们直接创建出来
	c := NewRedisCodeCache(rdb).(*RedisCodeCache)

	testCases := []struct {
		name string
		// 在 Redis 中准备数据
		before func(t *testing.T)
		// 验证 Redis 中的数据，
		// 也可以验证之后清理掉测试产生的数据
		after func(t *testing.T)

		// 输入
		ctx   context.Context
		biz   string
		phone string
		code  string

		// 预期输出
		wantErr error
	}{
		{
			name: "验证码存储成功",
			before: func(t *testing.T) {
				// 什么也不需要干
			},
			// 在设置成功的情况下，我们预期在 Redis 里面会有这个数据
			after: func(t *testing.T) {
				ctx := context.Background()
				key := c.key("login", "15212345678")
				val, err := rdb.Get(ctx, key).Result()
				// 断言必然取到了数据
				assert.NoError(t, err)
				assert.Equal(t, "123456", val)
				// 这里可以考虑进一步断言过期时间
				ttl, err := rdb.TTL(ctx, key).Result()
				assert.NoError(t, err)
				// 过期时间是十分钟，所以这里肯定会大于 9 分钟
				assert.True(t, ttl > time.Minute*9)
				err = rdb.Del(ctx, key).Err()
				assert.NoError(t, err)
			},
			ctx:   context.Background(),
			biz:   "login",
			phone: "15212345678",
			code:  "123456",
		},
		{
			// 注意，不要让这个用例依赖于上一个用例
			// 所以用了一个新的手机号码
			name: "发送太频繁",
			before: func(t *testing.T) {
				// 先准备一个数据，假装我们已经发送了一个验证码
				ctx := context.Background()
				key := c.key("login", "15212345679")
				// 简单验证这里咩有出错
				err := rdb.Set(ctx, key, "123456", time.Minute*9+time.Second*30).Err()
				assert.NoError(t, err)
			},
			// 发送太频繁的时候，我们预期还是原本的 123456
			after: func(t *testing.T) {
				ctx := context.Background()
				key := c.key("login", "15212345679")
				val, err := rdb.Get(ctx, key).Result()
				// 断言必然取到了数据
				assert.NoError(t, err)
				assert.Equal(t, "123456", val)
				// 这里没必要断言过期时间了，因为前面已经确定了 123456 还在
				err = rdb.Del(ctx, key).Err()
				assert.NoError(t, err)
			},
			ctx:     context.Background(),
			biz:     "login",
			phone:   "15212345679",
			code:    "234567",
			wantErr: ErrCodeSendTooMany,
		},
		{
			// 再次换了一个手机号码
			name: "未知错误",
			before: func(t *testing.T) {
				// 假装有人放了一个验证码，但是没有设置过期时间
				ctx := context.Background()
				key := c.key("login", "15212345670")
				// 传入 0 就是没有过期时间
				err := rdb.Set(ctx, key, "123456", 0).Err()
				assert.NoError(t, err)
				err = rdb.Del(ctx, key).Err()
				assert.NoError(t, err)
			},
			// 遇到未知错误的时候，Redis 中的数据也还是原本的数据
			// 也就是还保持那个没有设置过期时间的错误数据
			after: func(t *testing.T) {
				ctx := context.Background()
				key := c.key("login", "15212345670")
				val, err := rdb.Get(ctx, key).Result()
				// 断言必然取到了数据
				assert.NoError(t, err)
				assert.Equal(t, "123456", val)
				// 这里没必要断言过期时间了，因为前面已经确定了 123456 还在
			},
			ctx:     context.Background(),
			biz:     "login",
			phone:   "15212345670",
			code:    "234567",
			wantErr: ErrUnknownForCode,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.before(t)
			err := c.Set(tc.ctx, tc.biz, tc.phone, tc.code)
			assert.Equal(t, tc.wantErr, err)
			tc.after(t)
		})
	}
}
