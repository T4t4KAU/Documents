// Code generated by Kitex v0.6.2. DO NOT EDIT.

package addservice

import (
	"context"
	api "etcd-test/kitex_gen/api"
	client "github.com/cloudwego/kitex/client"
	kitex "github.com/cloudwego/kitex/pkg/serviceinfo"
)

func serviceInfo() *kitex.ServiceInfo {
	return addServiceServiceInfo
}

var addServiceServiceInfo = NewServiceInfo()

func NewServiceInfo() *kitex.ServiceInfo {
	serviceName := "AddService"
	handlerType := (*api.AddService)(nil)
	methods := map[string]kitex.MethodInfo{
		"Add": kitex.NewMethodInfo(addHandler, newAddServiceAddArgs, newAddServiceAddResult, false),
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

func addHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*api.AddServiceAddArgs)
	realResult := result.(*api.AddServiceAddResult)
	success, err := handler.(api.AddService).Add(ctx, realArg.Req)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newAddServiceAddArgs() interface{} {
	return api.NewAddServiceAddArgs()
}

func newAddServiceAddResult() interface{} {
	return api.NewAddServiceAddResult()
}

type kClient struct {
	c client.Client
}

func newServiceClient(c client.Client) *kClient {
	return &kClient{
		c: c,
	}
}

func (p *kClient) Add(ctx context.Context, req *api.AddRequest) (r *api.AddResponse, err error) {
	var _args api.AddServiceAddArgs
	_args.Req = req
	var _result api.AddServiceAddResult
	if err = p.c.Call(ctx, "Add", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}
