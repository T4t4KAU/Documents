package utils

import (
	"context"
	"douyin/biz/mw/minio"
	"fmt"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"strings"
)

func NewFileName(uid, time int64) string {
	return fmt.Sprintf("%d.%d", uid, time)
}

func URLConvert(ctx context.Context, c *app.RequestContext, path string) (url string) {
	if len(path) == 0 {
		return ""
	}
	arr := strings.Split(path, "/")
	u, err := minio.GetObjectURL(ctx, arr[0], arr[1])
	if err != nil {
		hlog.CtxInfof(ctx, err.Error())
		return ""
	}
	u.Scheme = string(c.URI().Scheme())
	u.Host = string(c.URI().Host())
	u.Path = "/src" + u.Path
	return u.String()
}
