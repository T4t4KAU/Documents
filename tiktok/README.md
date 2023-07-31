# Kitex实践：用户管理服务

本文讲述如何使用kitex开发一个用户管理微服务，负责用户的登录与注册

安装kitex: `go install github.com/cloudwego/kitex/tool/cmd/kitex@latest`

安装thrift-go:``go install github.com/cloudwego/thriftgo@latest`

## IDL代码生成

IDL是接口描述语言，本文使用thrift来编写

```protobuf
namespace go user

struct UserRegisterRequest {
    1: string username   // 用户名
    2: string password   // 密码
}

struct UserRegisterResponse {
    1: i32 status_code   // 状态码
    2: string status_msg // 状态信息
    3: i64 user_id       // 用户id
    4: string token      // 鉴权token
}

struct UserLoginRequest {
    1: string username
    2: string password

}

struct UserLoginResponse {
    1: i32 status_code
    2: string status_msg
    3: i64 user_id
    4: string token
}

service UserService {
    UserRegisterResponse UserRegister(1: UserRegisterRequest req)
    UserLoginResponse UserLogin(1: UserLoginRequest req)
}
```

代码生成：`kitex -module tiktok -service a.b.c idl/user.thrift`

随后执行`go mod tidy`安装依赖

目录结构

```powershell
tree .
.
├── build.sh
├── go.mod
├── go.sum
├── handler.go
├── idl
│   ├── common.thrift
│   └── user.thrift
├── kitex_gen
│   ├── common
│   │   ├── common.go
│   │   ├── k-common.go
│   │   └── k-consts.go
│   └── user
│       ├── k-consts.go
│       ├── k-user.go
│       ├── user.go
│       └── userservice
│           ├── client.go
│           ├── invoker.go
│           ├── server.go
│           └── userservice.go
├── kitex_info.yaml
├── main.go
└── script
    └── bootstrap.sh
```

以上都是kitex自动生成的代码，随后创建一个cmd/user目录: `mkdir -p cmd/user`

此处cmd目录用于存放代码的实现，user目录存放user服务相关的代码，将上述生成的main.go和handler.go移动到user目录

## 实现数据库操作

在user目录下创建了目录结构，dal/db中存放数据相关的代码

```powershell
user
├── dal
│   ├── db
│   │   ├── init.go           # 初始化数据库
│   │   ├── user.go           # 数据库操作
│   │   └── user_test.go      # 单元测试
│   ├── init.go
│   ├── pack
│   │   └── resp.go           # 响应格式
│   └── service
│       ├── user_login.go     # 用户登录服务
│       └── user_register.go  # 用户注册服务
├── handler.go
└── main.go
```

定义user结构体：

```go
type User struct {
	ID              int64  `json:"id"`               // 用户ID
	UserName        string `json:"user_name"`        // 用户名
	Password        string `json:"password"`         // 密码
	Avatar          string `json:"avatar"`           // 头像路径
	BackgroundImage string `json:"background_image"` // 背景路径
	Signature       string `json:"signature"`        // 签名
}
```

定义方法：

```go
// 返回数据库表名
func (User) TableName() string {
	return constants.UserTableName
}

// CreateUser 创建用户
func CreateUser(ctx context.Context, user *User) (int64, error) {
	err := dbConn.WithContext(ctx).Create(user).Error
	if err != nil {
		return 0, err
	}
	return user.ID, err
}

// QueryUserByName 通过用户名查询用户
func QueryUserByName(ctx context.Context, uname string) (*User, error) {
	var user User
	err := dbConn.WithContext(ctx).Where(
		"user_name = ?", uname).Find(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// QueryUserById 通过用户ID查询用户
func QueryUserById(ctx context.Context, userId int64) (*User, error) {
	var user User
	err := dbConn.WithContext(ctx).Where(
		"id = ?", userId).Find(&user).Error
	if err != nil {
		return nil, err
	}

	if user == (User{}) {
		err = errno.UserIsNotExistErr
		return nil, err
	}
	return &user, nil
}
```

上述User结构体实现TableName方法，决定了gorm自动建立表格时的数据库表名，隐含实现了gorm中Tabler接口 (区别于java使用implements显示实现接口)

数据库初始化代码 (以下constants包中的变量为定义好的常量，不作详细解释)：

```go
package db

var dbConn *gorm.DB

func Init() {
	var err error

	// 打开数据库连接
	dbConn, err = gorm.Open(mysql.Open(constants.MySQLDSN), &gorm.Config{
		PrepareStmt:            true,
		SkipDefaultTransaction: true,
	})
	if err != nil {
		panic(err)
	}
	
    // 分布式追踪
	err = dbConn.Use(gormopentracing.New())
	if err != nil {
		panic(err)
	}

	// 创建数据库表
	if !dbConn.Migrator().HasTable(&User{}) {
		err = dbConn.Migrator().CreateTable(&User{})
		if err != nil {
			panic(err)
		}
	}
}

```

下述为单元测试：

```go
package db

import (
	"context"
	"fmt"
	"testing"
)

func TestCreateUser(t *testing.T) {
	Init()
	user := &User{
		ID:       1000,
		UserName: "test",
		Password: "123456",
	}

	uid, err := CreateUser(context.Background(), user)
	if err != nil {
		t.Errorf(err.Error())
	}

	fmt.Printf("%v\n", uid)
}

func TestQueryUserByName(t *testing.T) {
	Init()
	user, err := QueryUserByName(context.Background(), "test")
	if err != nil {
		t.Errorf(err.Error())
		return
	}

	fmt.Printf("%v\n", user)
}

func TestQueryUserById(t *testing.T) {
	Init()
	user, err := QueryUserById(context.Background(), int64(1000))
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	fmt.Println(user)
}
```

## 完善业务逻辑

基于上述生成的代码，简单实现业务逻辑，具体实现在service目录下

用户注册：

```go
package service

type UserRegisterService struct {
	ctx context.Context
}

func NewUserRegisterService(ctx context.Context) *UserRegisterService {
	return &UserRegisterService{ctx: ctx}
}

func (s *UserRegisterService) UserRegister(req *user.UserRegisterRequest) (int64, error) {
    // 使用用户名查询用户是否存在
	u, err := db.QueryUserByName(s.ctx, req.Username)
	if err != nil {
		return int64(0), err
	}
    
    // 用户信息不为空 表明用户存在
	if *u != (db.User{}) {
		return int64(0), errno.UserAlreadyExistErr
	}
	
    // 对密码进行哈希
	hashedPassword, _ := utils.EncryptPassword(req.Password)
	uid, err := db.CreateUser(s.ctx, &db.User{
		UserName:        req.Username,
		Password:        hashedPassword,
		Avatar:          constants.TestAva,
		BackgroundImage: constants.TestBackground,
	})

	return uid, err
}
```

用户注册的流程是先读入请求参数中的username，先查询用户是否存在，如果已经存在则返回错误信息，如果不存在则读入密码并加密 (如何加密暂不做赘述)，将用户信息存入数据库，将用户ID返回

用户登录：

```go
package service

import (
	"context"
	"tiktok/cmd/user/dal/db"
	"tiktok/kitex_gen/user"
	"tiktok/pkg/errno"
	"tiktok/utils"
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
	
    // 校验密码
	if !utils.VerifyPassword(req.Password, u.Password) {
		return int64(0), errno.AuthorizationFailedErr
	}
	return u.ID, nil
}
```

## 完善请求处理

下面编写handler.go文件，完善请求处理逻辑，这部分为自动生成的代码，直接在生成的函数中编写即可：

```go
package main

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
	
    // 创建用户注册服务
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
	
    // 创建用户登录服务
	uid, err := service.NewUserLoginService(ctx).UserLogin(req)
	r := pack.BuildBaseResp(err)
	resp.StatusCode = r.StatusCode
	resp.StatusMsg = r.StatusMsg
	resp.UserId = uid

	return
}
```

在上述代码中，只须要调用service中的代码并返回相关请求即可

最后完善main函数，进行一些简单的配置：

```go
package main

// 服务地址
const serviceAddr = "127.0.0.1:8889"

func main() {
	dal.Init()
	
	addr, err := net.ResolveTCPAddr("tcp", serviceAddr)

	svr := user.NewServer(new(UserServiceImpl),
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: constants.UserServiceName}), // server name
		server.WithServiceAddr(addr),                                       // 设置服务地址
		server.WithLimit(&limit.Option{MaxConnections: 1000, MaxQPS: 100}), // 设置限制连接数
	)
	err = svr.Run()
	if err != nil {
		klog.Fatal(err)
	}
}
```

## 远程过程调用

写一段单元测试，来进行RPC调用

测试代码：

```go
func TestUserRegisterService(t *testing.T) {
	c, err := userservice.NewClient("user", client.WithHostPorts("127.0.0.1:8889"))
	if err != nil {
		t.Errorf("New client error: %#v", err)
		return
	}
	
    // 调用RPC函数
	resp, err := c.UserRegister(context.Background(), &user.UserRegisterRequest{
		Username: "test",
		Password: "123456",
	})
	if err != nil {
		t.Errorf("user register error: %#v\n", err)
		return
	}
	fmt.Printf("%#v", resp)
}

```

运行结果，成功打印出了用户结构体

```go
=== RUN   TestUserRegisterService
&user.UserRegisterResponse{StatusCode:0, StatusMsg:"Success", UserId:1026, Token:""}
--- PASS: TestUserRegisterService (0.02s)
PASS
```

## 完整实现

在目前的完整实现中，加入了API网关，etcd服务注册，jaeger链路追踪，jwt鉴权，并设置了反向代理，这些内容比较繁杂，暂不做过多解释了，目前的项目目录结构：

```powershell
tree .
.
├── README.md
├── build.sh
├── cmd
│   ├── api # API网关
│   │   ├── handlers
│   │   │   ├── handler.go
│   │   │   ├── param.go
│   │   │   └── resp.go
│   │   └── rpc
│   │       ├── init.go
│   │       ├── user.go
│   │       └── user_test.go
│   ├── main.go # 主调函数
│   └── user    # 用户服务
│       ├── dal
│       │   ├── db
│       │   │   ├── init.go
│       │   │   ├── user.go
│       │   │   └── user_test.go
│       │   ├── init.go
│       │   ├── pack
│       │   │   └── resp.go
│       │   └── service
│       │       ├── user_login.go
│       │       └── user_register.go
│       ├── handler.go
│       └── main.go
├── docker-compose.yml  # 配置docker
├── go.mod
├── go.sum
├── idl
│   ├── common.thrift
│   └── user.thrift
├── kitex_gen
│   ├── common
│   │   ├── common.go
│   │   ├── k-common.go
│   │   └── k-consts.go
│   └── user
│       ├── k-consts.go
│       ├── k-user.go
│       ├── user.go
│       └── userservice
│           ├── client.go
│           ├── invoker.go
│           ├── server.go
│           └── userservice.go
├── kitex_info.yaml
├── pkg
│   ├── configs
│   │   └── sql
│   ├── constants
│   │   └── constant.go
│   ├── errno
│   │   └── errno.go
│   ├── mw
│   │   ├── client.go
│   │   ├── common.go
│   │   └── server.go
│   └── trace
│       └── trace.go
├── script
│   └── bootstrap.sh
├── test  # 测试代码
│   ├── common.go
│   └── user_api_test.go
└── utils
    └── utils.go
```

先运行上述的用户服务，再运行API网关，如下是控制台信息

```powershell
2023/07/31 14:26:10 debug logging disabled
2023/07/31 14:26:10 debug logging disabled
2023/07/31 14:26:10.262669 engine.go:617: [Debug] HERTZ: Method=POST   absolutePath=/douyin/user/register     --> handlerName=tiktok/cmd/api/handlers.Register (num=2 handlers)
2023/07/31 14:26:10.262847 engine.go:617: [Debug] HERTZ: Method=POST   absolutePath=/douyin/user/login        --> handlerName=github.com/hertz-contrib/jwt.(*HertzJWTMiddleware).LoginHandler-fm (num=2 handlers)
2023/07/31 14:26:10.263017 engine.go:389: [Info] HERTZ: Using network library=netpoll
2023/07/31 14:26:10.263141 transport.go:115: [Info] HERTZ: HTTP server listening on address=[::]:8888
```

使用这段测试代码来向网关发送HTTP请求:

```go
// 测试用户注册接口
func TestUserRegister(t *testing.T) {
	url := serverAddr + "/douyin/user/register?username=test666&password=123456"
	method := "POST"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("User-Agent", "Apifox/1.0.0 (https://apifox.com)")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))
}
```

运行成功: `{"status_code":0,"status_msg":"Success","user_id":1027,"token":""}`

```go
// 测试用户登录接口
func TestUserLogin(t *testing.T) {
	url := serverAddr + "/douyin/user/login/?username=test&password=123456"
	method := "POST"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("User-Agent", "Apifox/1.0.0 (https://apifox.com)")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))
}
```

运行成功：`{"status_code":0,"status_msg":"Success","token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2OTA4NzE1MDMsIm9yaWdfaWF0IjoxNjkwNzg1MTAzLCJ1c2VyX2lkIjoxMDI2fQ.JBU8j6IzcVZTx6xrnIrPJSVsGfvGDAxXlGguKBrRKUc"}`