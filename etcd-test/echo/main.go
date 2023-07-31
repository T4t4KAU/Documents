package main

import (
	api "etcd-test/kitex_gen/api/echoservice"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	etcd "github.com/kitex-contrib/registry-etcd"
	"log"
	"net"
)

func main() {
	r, err := etcd.NewEtcdRegistry([]string{"127.0.0.1:10079"})
	if err != nil {
		panic(err)
	}

	addr, _ := net.ResolveTCPAddr("tcp", "127.0.0.1:9092")

	svr := api.NewServer(new(EchoServiceImpl),
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: "echo"}),
		server.WithServiceAddr(addr),
		server.WithRegistry(r),
	)

	err = svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
