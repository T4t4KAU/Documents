package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"gitee.com/geekbang/basic-go/webook/internal/web"
	"gitee.com/geekbang/basic-go/webook/ioc"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestUserHandler_SendSMSLoginCode(t *testing.T) {
	const sendSMSCodeUrl = "/users/login_sms/code/send"
	// 使用依赖注入的 server
	server := InitWebServer()
	rdb := ioc.InitRedis()
	testCases := []struct {
		// 名字
		name string
		// 要提前准备数据
		before func(t *testing.T)
		// 验证并且删除数据
		after func(t *testing.T)
		// 目前前端只有一个手机号码作为输入
		phone string

		// 预期响应
		wantCode   int
		wantResult web.Result
	}{
		{
			name: "发送成功",
			before: func(t *testing.T) {
				// 啥也不做
			},
			// 在设置成功的情况下，我们预期在 Redis 里面会有这个数据
			after: func(t *testing.T) {
				ctx := context.Background()
				key := "phone_code:login:15212345678"
				val, err := rdb.Get(ctx, key).Result()
				// 断言必然取到了数据
				assert.NoError(t, err)
				assert.True(t, len(val) == 6)
				// 这里可以考虑进一步断言过期时间
				ttl, err := rdb.TTL(ctx, key).Result()
				assert.NoError(t, err)
				// 过期时间是十分钟，所以这里肯定会大于 9 分钟
				assert.True(t, ttl > time.Minute*9)

				// 删除数据
				err = rdb.Del(ctx, key).Err()
				assert.NoError(t, err)
			},
			phone:    "15212345678",
			wantCode: 200,
			wantResult: web.Result{
				Msg: "发送成功",
			},
		},
		{
			name: "空的手机号码",
			before: func(t *testing.T) {
				// 啥也不做
			},
			after: func(t *testing.T) {

			},
			wantCode: 200,
			wantResult: web.Result{
				Code: 4,
				Msg:  "请输入手机号码",
			},
		},
		{
			name: "发送太频繁",
			before: func(t *testing.T) {
				// 先准备一个数据，假装我们已经发送了一个验证码
				ctx := context.Background()
				// 这个 key 和 cache 强耦合，不过这是测试，没办法的事情
				key := "phone_code:login:15212345679"
				// 简单验证这里咩有出错
				err := rdb.Set(ctx, key, "123456", time.Minute*9+time.Second*30).Err()
				assert.NoError(t, err)
			},
			// 发送太频繁的时候，我们预期还是原本的 123456
			after: func(t *testing.T) {
				ctx := context.Background()
				key := "phone_code:login:15212345679"
				val, err := rdb.Get(ctx, key).Result()
				// 断言必然取到了数据
				assert.NoError(t, err)
				assert.Equal(t, "123456", val)
				// 这里没必要断言过期时间了，因为前面已经确定了 123456 还在

				// 删除数据
				err = rdb.Del(ctx, key).Err()
				assert.NoError(t, err)
			},
			phone:    "15212345679",
			wantCode: 200,
			wantResult: web.Result{
				Code: 4,
				Msg:  "短信发送太频繁，请稍后再试",
			},
		},
		{
			name: "未知错误",
			before: func(t *testing.T) {
				// 假装有人放了一个验证码，但是没有设置过期时间
				ctx := context.Background()
				key := "phone_code:login:15212345670"
				// 传入 0 就是没有过期时间
				err := rdb.Set(ctx, key, "123456", 0).Err()
				assert.NoError(t, err)
			},
			// 遇到未知错误的时候，Redis 中的数据也还是原本的数据
			// 也就是还保持那个没有设置过期时间的错误数据
			after: func(t *testing.T) {
				ctx := context.Background()
				key := "phone_code:login:15212345670"
				val, err := rdb.Get(ctx, key).Result()
				// 断言必然取到了数据
				assert.NoError(t, err)
				assert.Equal(t, "123456", val)
				// 这里没必要断言过期时间了，因为前面已经确定了 123456 还在

				// 删除数据
				err = rdb.Del(ctx, key).Err()
				assert.NoError(t, err)
			},
			phone:    "15212345670",
			wantCode: 200,
			wantResult: web.Result{
				Code: 5,
				Msg:  "系统错误",
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.before(t)

			body := fmt.Sprintf(`{"phone": "%s"}`, tc.phone)
			req, err := http.NewRequest(http.MethodPost, sendSMSCodeUrl,
				bytes.NewBuffer([]byte(body)))
			req.Header.Set("Content-Type", "application/json")
			assert.NoError(t, err)

			recorder := httptest.NewRecorder()
			server.ServeHTTP(recorder, req)

			code := recorder.Code
			// 反序列化为结果
			var result web.Result
			err = json.Unmarshal(recorder.Body.Bytes(), &result)
			assert.NoError(t, err)
			assert.Equal(t, tc.wantCode, code)
			assert.Equal(t, tc.wantResult, result)
			tc.after(t)
		})
	}
}
