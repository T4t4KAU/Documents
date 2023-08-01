package main

import (
	"douyin/cmd/relation/dal"
	"douyin/kitex_gen/relation/relationservice"
	"douyin/pkg/constants"
	"douyin/pkg/mw"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/pkg/limit"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	etcd "github.com/kitex-contrib/registry-etcd"
	"net"
)

const serviceAddr = "127.0.0.1:9092"

func main() {
	dal.Init()

	r, err := etcd.NewEtcdRegistry([]string{
		constants.EtcdAddress,
	})

	addr, err := net.ResolveTCPAddr("tcp", serviceAddr)

	svr := relationservice.NewServer(new(RelationServiceImpl),
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: constants.RelationServiceName}), // server name
		server.WithMiddleware(mw.CommonMiddleware),                                                         // middleware
		server.WithMiddleware(mw.ServerMiddleware),
		server.WithServiceAddr(addr),                                       // address
		server.WithLimit(&limit.Option{MaxConnections: 1000, MaxQPS: 100}), // limit
		server.WithRegistry(r),                                             // registry
	)
	err = svr.Run()
	if err != nil {
		klog.Fatal(err)
	}
}
