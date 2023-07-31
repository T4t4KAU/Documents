# Kitex微服务开发实践(ETCD服务注册)

服务注册通常用于分布式系统或微服务架构中，是一种用于管理和发现这些分布式服务的机制。它的目标是让服务能够动态地找到其他服务，并能够与其进行通信，而无需显式地配置其位置信息

本文简单讲述使用etcd进行服务注册，基于kitex和hertz框架简单实现微服务应用

## 接口定义

使用thrift编写如下idl

add.thrift

```thrift
namespace go api

struct AddRequest {
    1: i32 first
    2: i32 second
}

struct AddResponse {
    1: i32 sum
}

service AddService {
    AddResponse Add(AddRequest req)
}
```

echo.thrift

```thrift
namespace go api

struct EchoRequest {
    1: string message
}

struct EchoResponse {
    2: string message
}

service EchoService {
    EchoResponse Echo(1:EchoRequest req)
}
```

使用Kitex生成代码(此处不作赘述):

```powershell
.
├── add
│   ├── handler.go
│   └── main.go
├── build.sh
├── echo
│   ├── handler.go
│   └── main.go
├── go.mod
├── go.sum
├── idl
│   ├── add.thrift
│   └── echo.thrift
├── kitex_gen
│   └── api
│       ├── add.go
│       ├── addservice
│       │   ├── addservice.go
│       │   ├── client.go
│       │   ├── invoker.go
│       │   └── server.go
│       ├── echo.go
│       ├── echoservice
│       │   ├── client.go
│       │   ├── echoservice.go
│       │   ├── invoker.go
│       │   └── server.go
│       ├── k-add.go
│       ├── k-consts.go
│       └── k-echo.go
├── kitex_info.yaml
└── script
    └── bootstrap.sh
```

echo服务将参数message原样返回，add服务将参数中的两个整数求和后返回

## 代码实现

首先实现一个无服务注册的版本，在上述生成代码的基础上完善方法实现:

add/hander.go

```go
// Add implements the AddServiceImpl interface.
func (s *AddServiceImpl) Add(ctx context.Context, req *api.AddRequest) (resp *api.AddResponse, err error) {
	resp = new(api.AddResponse)
	resp.Sum = req.First + req.Second
	return
}
```

add/main.go

```go
package main

func main() {
    // 服务地址
	addr, _ := net.ResolveTCPAddr("tcp", "127.0.0.1:9091")

	svr := api.NewServer(new(AddServiceImpl),
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: "add"}),
		server.WithServiceAddr(addr),
	)

	err := svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
```

echo/handler.go

```go
// Echo implements the EchoServiceImpl interface.
func (s *EchoServiceImpl) Echo(ctx context.Context, req *api.EchoRequest) (resp *api.EchoResponse, err error) {
	resp = new(api.EchoResponse)
	resp.Message = req.Message
	return
}
```

echo/main.go

```go
package main

func main() {
    // 服务地址
	addr, _ := net.ResolveTCPAddr("tcp", "127.0.0.1:9092")

	svr := api.NewServer(new(EchoServiceImpl),
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: "echo"}),
		server.WithServiceAddr(addr),
	)

	err := svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
```

创建一个api包，作为两个服务的API网关，使用hertz来处理用户请求

目录结构：

```powershell
.
├── add
│   ├── handler.go
│   └── main.go
├── api
│   ├── add.go
│   ├── echo.go
│   └── main.go
├── build.sh
├── echo
│   ├── handler.go
│   └── main.go
├── go.mod
├── go.sum
├── idl
│   ├── add.thrift
│   └── echo.thrift
├── kitex_gen
│   └── api
│       ├── add.go
│       ├── addservice
│       │   ├── addservice.go
│       │   ├── client.go
│       │   ├── invoker.go
│       │   └── server.go
│       ├── echo.go
│       ├── echoservice
│       │   ├── client.go
│       │   ├── echoservice.go
│       │   ├── invoker.go
│       │   └── server.go
│       ├── k-add.go
│       ├── k-consts.go
│       └── k-echo.go
├── kitex_info.yaml
└── script
    └── bootstrap.sh
```

实现web handler函数，在handler中使用RPC请求服务

add:

```go
package main

var addClient addservice.Client

func AddHandler(ctx context.Context, c *app.RequestContext) {
	num1, _ := strconv.Atoi(c.Query("first"))
	num2, _ := strconv.Atoi(c.Query("second"))
    
    // RPC请求
	resp, err := addClient.Add(ctx, &api.AddRequest{
		First: int32(num1), Second: int32(num2),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, utils.H{
		"message": resp.Sum,
	})
}
```

echo:

```go
package main

var echoClient echoservice.Client

func EchoHandler(ctx context.Context, c *app.RequestContext) {
    // RPC请求
	resp, err := echoClient.Echo(ctx, &api.EchoRequest{
		Message: c.Query("message"),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.H{
			"message": err.Error(),
		})
	}

	c.JSON(http.StatusOK, utils.H{
		"message": resp.Message,
	})
}
```

在主调函数中进行RPC客户端初始化和路由注册:

```go
package main

func Init() {
	var err error
	addClient, err = addservice.NewClient("addservice",
		client.WithHostPorts("127.0.0.1:9091"))
	if err != nil {
		panic(err)
	}
	echoClient, err = echoservice.NewClient("echoservice",
		client.WithHostPorts("127.0.0.1:9092"))
	if err != nil {
		panic(err)
	}
}

func main() {
	Init()
	r := server.Default()
	r.POST("/add", AddHandler)
	r.GET("/echo", EchoHandler)

	r.Spin()
}
```

启动add和echo服务，再启动api，请求URL即可访问，此处省略测试过程

## 引入服务注册

在上述实现中，要显式的告诉API网关服务的地址，一旦这个地址发生变化，会导致服务不可用，接下来使用服务注册

使用docker启动etc:

```powershell
docker run -d -p 10079:2379 --name etcd \
  -e ETCD_LISTEN_CLIENT_URLS=http://0.0.0.0:2379 \
  -e ETCD_ADVERTISE_CLIENT_URLS=http://0.0.0.0:2379 \
  -e ETCDCTL_API=3 \
  quay.io/coreos/etcd:v3.5.5
```

etcd服务器启动在10079端口

接下来修改服务代码：

add/main.go

```go
package main

func main() {
    // 服务注册
	r, err := etcd.NewEtcdRegistry([]string{"127.0.0.1:10079"})
	if err != nil {
		panic(err)
	}

	addr, _ := net.ResolveTCPAddr("tcp", "127.0.0.1:9091")

	svr := api.NewServer(new(AddServiceImpl),
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: "add"}), // 指定服务名称
		server.WithServiceAddr(addr),
		server.WithRegistry(r),  // 设置服务注册
	)

	err = svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
```

echo/main.go

```go
package main
func main() {
    // 服务注册
	r, err := etcd.NewEtcdRegistry([]string{"127.0.0.1:10079"})
	if err != nil {
		panic(err)
	}

	addr, _ := net.ResolveTCPAddr("tcp", "127.0.0.1:9092")

	svr := api.NewServer(new(EchoServiceImpl),
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: "echo"}), // 指定服务名称
		server.WithServiceAddr(addr),
		server.WithRegistry(r),  // 设置服务注册
	)

	err = svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
```

修改api代码，无须给出两个服务的地址了:

```go
package main

func Init() {
	var err error

	r, err := etcd.NewEtcdResolver([]string{"127.0.0.1:10079"})
	if err != nil {
		panic(err)
	}

	addClient, err = addservice.NewClient("add",
		client.WithResolver(r),
	)
	if err != nil {
		panic(err)
	}
	echoClient, err = echoservice.NewClient("echo",
		client.WithResolver(r),
	)
	if err != nil {
		panic(err)
	}
}

func main() {
	Init()
	r := server.Default()
	r.POST("/add", AddHandler)
	r.GET("/echo", EchoHandler)

	r.Spin()
}
```

启动之后，服务正常使用
