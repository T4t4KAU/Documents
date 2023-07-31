package main

import (
	"context"
	api "etcd-test/kitex_gen/api"
)

// AddServiceImpl implements the last service interface defined in the IDL.
type AddServiceImpl struct{}

// Add implements the AddServiceImpl interface.
func (s *AddServiceImpl) Add(ctx context.Context, req *api.AddRequest) (resp *api.AddResponse, err error) {
	resp = new(api.AddResponse)
	resp.Sum = req.First + req.Second
	return
}
