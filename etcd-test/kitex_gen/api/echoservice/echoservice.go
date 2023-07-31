// Code generated by Kitex v0.6.2. DO NOT EDIT.

package echoservice

import (
	"context"
	api "etcd-test/kitex_gen/api"
	client "github.com/cloudwego/kitex/client"
	kitex "github.com/cloudwego/kitex/pkg/serviceinfo"
)

func serviceInfo() *kitex.ServiceInfo {
	return echoServiceServiceInfo
}

var echoServiceServiceInfo = NewServiceInfo()

func NewServiceInfo() *kitex.ServiceInfo {
	serviceName := "EchoService"
	handlerType := (*api.EchoService)(nil)
	methods := map[string]kitex.MethodInfo{
		"Echo": kitex.NewMethodInfo(echoHandler, newEchoServiceEchoArgs, newEchoServiceEchoResult, false),
	}
	extra := map[string]interface{}{
		"PackageName": "api",
	}
	svcInfo := &kitex.ServiceInfo{
		ServiceName:     serviceName,
		HandlerType:     handlerType,
		Methods:         methods,
		PayloadCodec:    kitex.Thrift,
		KiteXGenVersion: "v0.6.2",
		Extra:           extra,
	}
	return svcInfo
}

func echoHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*api.EchoServiceEchoArgs)
	realResult := result.(*api.EchoServiceEchoResult)
	success, err := handler.(api.EchoService).Echo(ctx, realArg.Req)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newEchoServiceEchoArgs() interface{} {
	return api.NewEchoServiceEchoArgs()
}

func newEchoServiceEchoResult() interface{} {
	return api.NewEchoServiceEchoResult()
}

type kClient struct {
	c client.Client
}

func newServiceClient(c client.Client) *kClient {
	return &kClient{
		c: c,
	}
}

func (p *kClient) Echo(ctx context.Context, req *api.EchoRequest) (r *api.EchoResponse, err error) {
	var _args api.EchoServiceEchoArgs
	_args.Req = req
	var _result api.EchoServiceEchoResult
	if err = p.c.Call(ctx, "Echo", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}
