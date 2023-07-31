package main

import (
	"context"
	"etcd-test/kitex_gen/api"
	"etcd-test/kitex_gen/api/echoservice"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"net/http"
)

var echoClient echoservice.Client

func EchoHandler(ctx context.Context, c *app.RequestContext) {
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
