package main

import (
	"context"
	api "etcd-test/kitex_gen/api"
)

// EchoServiceImpl implements the last service interface defined in the IDL.
type EchoServiceImpl struct{}

// Echo implements the EchoServiceImpl interface.
func (s *EchoServiceImpl) Echo(ctx context.Context, req *api.EchoRequest) (resp *api.EchoResponse, err error) {
	resp = new(api.EchoResponse)
	resp.Message = req.Message
	return
}
