package sms

import "context"

// Service 发送短信的抽象
// 目前你可以理解为，这是一个为了适配不同的短信供应商的抽象
type Service interface {
	Send(ctx context.Context, tplId string,
		args []string, numbers ...string) error
}
