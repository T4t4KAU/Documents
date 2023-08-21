package web

import (
	"bytes"
	"errors"
	"gitee.com/geekbang/basic-go/webook/internal/service"
	svcmocks "gitee.com/geekbang/basic-go/webook/internal/service/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUserHandler_SignUp(t *testing.T) {
	const signupUrl = "/users/signup"
	testCases := []struct {
		// 名字
		name string

		// 准备 mock
		// 因为 UserHandler 用到了 UserService 和 CodeService
		// 所以我们需要准备这两个的 mock 实例。
		// 因此你能看到它返回了 UserService 和 CodeService
		mock func(ctrl *gomock.Controller) (service.UserService, service.CodeService)

		// 输入，因为 request 的构造过程可能很复杂
		// 所以我们在这里定义一个 Builder
		reqBuilder func(t *testing.T) *http.Request

		// 预期响应
		wantCode int
		wantBody string
	}{
		{
			name: "注册成功",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				usersvc := svcmocks.NewMockUserService(ctrl)
				usersvc.EXPECT().Signup(gomock.Any(), gomock.Any()).
					// 注册成功，也就是 UserService 返回了 nil
					Return(nil)

				// 在 signup 这个接口里面，并没有用到的 codesvc，
				// 所以什么不需要准备模拟调用
				codesvc := svcmocks.NewMockCodeService(ctrl)
				return usersvc, codesvc
			},
			reqBuilder: func(t *testing.T) *http.Request {
				body := bytes.NewBuffer([]byte(`{"email":"123@qq.com","password":"hello@world123","confirmPassword":"hello@world123"}`))
				req, err := http.NewRequest(http.MethodPost, signupUrl, body)
				req.Header.Set("Content-Type", "application/json")
				if err != nil {
					t.Fatal(err)
				}
				return req
			},
			wantCode: 200,
			wantBody: "hello, 注册成功",
		},
		{
			name: "非 JSON 输入",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				// 因为根本没有跑到 singup 那里，所以直接返回 nil 都可以
				return nil, nil
			},
			reqBuilder: func(t *testing.T) *http.Request {
				// 准备一个错误的JSON 串
				body := bytes.NewBuffer([]byte(`{"email":"123@qq.com",`))
				req, err := http.NewRequest(http.MethodPost, signupUrl, body)
				req.Header.Set("Content-Type", "application/json")
				if err != nil {
					t.Fatal(err)
				}
				return req
			},
			wantCode: 400,
		},

		{
			name: "邮箱格式不对",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				// 因为根本没有跑到 signup 那里，所以直接返回 nil 都可以
				return nil, nil
			},
			reqBuilder: func(t *testing.T) *http.Request {
				// 准备一个不合法的邮箱
				body := bytes.NewBuffer([]byte(`{"email":"123@","password":"hello@world123","confirmPassword":"hello@world123"}`))
				req, err := http.NewRequest(http.MethodPost, signupUrl, body)
				req.Header.Set("Content-Type", "application/json")
				if err != nil {
					t.Fatal(err)
				}
				return req
			},
			wantCode: 200,
			wantBody: "邮箱不正确",
		},
		{
			name: "两次密码输入不同",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				// 因为根本没有跑到 signup 那里，所以直接返回 nil 都可以
				return nil, nil
			},
			reqBuilder: func(t *testing.T) *http.Request {
				// 准备一个不合法的邮箱
				body := bytes.NewBuffer([]byte(`{"email":"123@qq.com","password":"hello","confirmPassword":"hello@world123"}`))
				req, err := http.NewRequest(http.MethodPost, signupUrl, body)
				req.Header.Set("Content-Type", "application/json")
				if err != nil {
					t.Fatal(err)
				}
				return req
			},
			wantCode: 200,
			wantBody: "两次输入的密码不相同",
		},
		{
			name: "密码格式不对",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				// 因为根本没有跑到 signup 那里，所以直接返回 nil 都可以
				return nil, nil
			},
			reqBuilder: func(t *testing.T) *http.Request {
				// 准备一个不合法的邮箱
				body := bytes.NewBuffer([]byte(`{"email":"123@qq.com","password":"hello","confirmPassword":"hello"}`))
				req, err := http.NewRequest(http.MethodPost, signupUrl, body)
				req.Header.Set("Content-Type", "application/json")
				if err != nil {
					t.Fatal(err)
				}
				return req
			},
			wantCode: 200,
			wantBody: "密码必须包含数字、特殊字符，并且长度不能小于 8 位",
		},
		{
			name: "邮箱冲突",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				usersvc := svcmocks.NewMockUserService(ctrl)
				usersvc.EXPECT().Signup(gomock.Any(), gomock.Any()).
					// 模拟返回邮箱冲突的异常
					Return(service.ErrUserDuplicateEmail)

				// 在 signup 这个接口里面，并没有用到的 codesvc，
				// 所以什么不需要准备模拟调用
				codesvc := svcmocks.NewMockCodeService(ctrl)
				return usersvc, codesvc
			},
			reqBuilder: func(t *testing.T) *http.Request {
				body := bytes.NewBuffer([]byte(`{"email":"123@qq.com","password":"hello@world123","confirmPassword":"hello@world123"}`))
				req, err := http.NewRequest(http.MethodPost, signupUrl, body)
				req.Header.Set("Content-Type", "application/json")
				if err != nil {
					t.Fatal(err)
				}
				return req
			},
			wantCode: 200,
			wantBody: "重复邮箱，请换一个邮箱",
		},
		{
			name: "系统异常",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				usersvc := svcmocks.NewMockUserService(ctrl)
				usersvc.EXPECT().Signup(gomock.Any(), gomock.Any()).
					// 注册失败，系统本身的异常
					Return(errors.New("模拟系统异常"))

				// 在 signup 这个接口里面，并没有用到的 codesvc，
				// 所以什么不需要准备模拟调用
				codesvc := svcmocks.NewMockCodeService(ctrl)
				return usersvc, codesvc
			},
			reqBuilder: func(t *testing.T) *http.Request {
				body := bytes.NewBuffer([]byte(`{"email":"123@qq.com","password":"hello@world123","confirmPassword":"hello@world123"}`))
				req, err := http.NewRequest(http.MethodPost, signupUrl, body)
				req.Header.Set("Content-Type", "application/json")
				if err != nil {
					t.Fatal(err)
				}
				return req
			},
			wantCode: 200,
			wantBody: "服务器异常，注册失败",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			usersvc, codesvc := tc.mock(ctrl)
			// 利用 mock 来构造 UserHandler
			hdl := NewUserHandler(usersvc, codesvc)

			// 注册路由
			server := gin.Default()
			hdl.RegisterRoutes(server)
			// 准备请求
			req := tc.reqBuilder(t)
			// 准备记录响应
			recorder := httptest.NewRecorder()
			// 执行
			server.ServeHTTP(recorder, req)
			// 断言
			assert.Equal(t, tc.wantCode, recorder.Code)
			assert.Equal(t, tc.wantBody, recorder.Body.String())
		})
	}
}

func TestMock(t *testing.T) {
	// 先创建一个控制 mock 的控制器
	ctrl := gomock.NewController(t)
	// 每个测试结束都要调用 Finish，
	// 然后 mock 就会验证你的测试流程是否符合预期
	defer ctrl.Finish()
	usersvc := svcmocks.NewMockUserService(ctrl)
	// 开始设计一个个模拟调用
	// 预期第一个是 Signup 的调用
	// 模拟的条件是 gomock.Any, gomock.Any。
	// 然后返回
	usersvc.EXPECT().Signup(gomock.Any(), gomock.Any()).
		Return(errors.New("模拟的错误"))

	// 后面就可以使用这个 usersvc 了
	//hdl := NewUserHandler(usersvc, nil)
}

//func TestHTTPDemo(t *testing.T) {
//
//	// 注册路由
//	server := gin.Default()
//	// 要用 mock 来模拟，这里还没处理
//	var usersvc service.UserService
//	var codesvc service.CodeService
//	hdl := NewUserHandler(usersvc, codesvc)
//	hdl.RegisterRoutes(server)
//
//	// 构造请求
//	req, err := http.NewRequest(http.MethodPost,
//		"/users/signup", nil)
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	// 准备接收响应
//	recorder := httptest.NewRecorder()
//
//	// 本地直接假装收到了 HTTP 请求
//	server.ServeHTTP(recorder, req)
//
//	// 断言结果
//	assert.Equal(t, 200, recorder.Code)
//}

// TestEmailPattern 用来验证我们的邮箱正则表达式对不对
func TestEmailPattern(t *testing.T) {
	testCases := []struct {
		name  string
		email string
		match bool
	}{
		{
			name:  "不带@",
			email: "123456",
			match: false,
		},
		{
			name:  "带@ 但是没后缀",
			email: "123456@",
			match: false,
		},
		{
			name:  "合法邮箱",
			email: "123456@qq.com",
			match: true,
		},
	}

	h := NewUserHandler(nil, nil)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			match, err := h.emailRegexExp.MatchString(tc.email)
			require.NoError(t, err)
			assert.Equal(t, tc.match, match)
		})
	}
}

func TestComplete(t *testing.T) {
	// 演示完整的测试用例
	testCases := []struct {
		// 用例名字，简明扼要说清测试的场景
		name string

		// 这边需要有预期输入，根据你的方法参数、接收器来设计

		// 这边需要有预期输出，根据你的方法返回值、接收器来设计

		// mock 数据，在单元测试里面很常见
		mock func(ctrl *gomock.Controller)
		// 测试用例准备环境、数据等
		before func(t *testing.T)
		// 数据清理等
		after func(t *testing.T)
	}{
		{},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

		})
	}
}

func TestPasswordPattern(t *testing.T) {
	testCases := []struct {
		name     string
		password string
		match    bool
	}{
		{
			name:     "合法密码",
			password: "Hello#world123",
			match:    true,
		},
		{
			name:     "没有数字",
			password: "Hello#world",
			match:    false,
		},
		{
			name:     "没有特殊字符",
			password: "Helloworld123",
			match:    false,
		},
		{
			name:     "长度不足",
			password: "he!123",
			match:    false,
		},
	}

	h := NewUserHandler(nil, nil)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			match, err := h.passwordRegexExp.MatchString(tc.password)
			require.NoError(t, err)
			assert.Equal(t, tc.match, match)
		})
	}
}
