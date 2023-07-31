package main

import (
	"context"
	"etcd-test/kitex_gen/api"
	"etcd-test/kitex_gen/api/addservice"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"net/http"
	"strconv"
)

var addClient addservice.Client

func AddHandler(ctx context.Context, c *app.RequestContext) {
	num1, _ := strconv.Atoi(c.Query("first"))
	num2, _ := strconv.Atoi(c.Query("second"))
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
