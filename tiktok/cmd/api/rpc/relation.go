package rpc

import (
	"context"
	"douyin/kitex_gen/relation"
	"douyin/kitex_gen/relation/relationservice"
	"douyin/pkg/constants"
	"douyin/pkg/errno"
	"douyin/pkg/mw"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/retry"
	etcd "github.com/kitex-contrib/registry-etcd"
	trace "github.com/kitex-contrib/tracer-opentracing"
	"time"
)

var relationClient relationservice.Client

func initRelationRPC() {
	r, err := etcd.NewEtcdResolver([]string{constants.EtcdAddress})
	if err != nil {
		panic(err)
	}

	c, err := relationservice.NewClient(
		constants.RelationServiceName,
		client.WithMiddleware(mw.CommonMiddleware),
		client.WithInstanceMW(mw.ClientMiddleware),
		client.WithMuxConnection(1),                       // mux
		client.WithRPCTimeout(3*time.Second),              // rpc timeout
		client.WithConnectTimeout(50*time.Millisecond),    // conn timeout
		client.WithFailureRetry(retry.NewFailurePolicy()), // retry
		client.WithSuite(trace.NewDefaultClientSuite()),   // tracer
		client.WithResolver(r),                            // resolver
	)
	if err != nil {
		panic(err)
	}

	relationClient = c
}

func RelationAction(ctx context.Context, req *relation.RelationActionRequest) (*relation.RelationActionResponse, error) {
	resp, err := relationClient.RelationAction(ctx, req)
	if err != nil {
		return resp, err
	}

	if resp.StatusCode != 0 {
		return resp, errno.NewErrNo(resp.StatusCode, resp.StatusMsg)
	}
	return resp, nil
}
