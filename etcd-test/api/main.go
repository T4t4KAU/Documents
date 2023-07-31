package main

import (
	"etcd-test/kitex_gen/api/addservice"
	"etcd-test/kitex_gen/api/echoservice"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/kitex/client"
	etcd "github.com/kitex-contrib/registry-etcd"
)

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
