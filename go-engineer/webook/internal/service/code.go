package service

import (
	"context"
	"fmt"
	"gitee.com/geekbang/basic-go/webook/internal/repository"
	"gitee.com/geekbang/basic-go/webook/internal/service/sms"
	"math/rand"
)

var ErrCodeSendTooMany = repository.ErrCodeSendTooMany

const codeTplId = "1877556"

type CodeService interface {
	Send(ctx context.Context, biz string, phone string) error
	Verify(ctx context.Context, biz string, phone string, inputCode string) (bool, error)
}

// SMSCodeService 短信验证码的实现
type SMSCodeService struct {
	sms  sms.Service
	repo repository.CodeRepository
}

func NewSMSCodeService(svc sms.Service, repo repository.CodeRepository) CodeService {
	return &SMSCodeService{
		sms:  svc,
		repo: repo,
	}
}

// Send 生成一个随机验证码，并发送
func (c *SMSCodeService) Send(ctx context.Context, biz string, phone string) error {
	code := c.generate()
	err := c.repo.Store(ctx, biz, phone, code)
	if err != nil {
		return err
	}
	err = c.sms.Send(ctx, codeTplId, []string{code}, phone)
	return err
}

// Verify 验证验证码
func (c *SMSCodeService) Verify(ctx context.Context,
	biz string,
	phone string,
	inputCode string) (bool, error) {
	ok, err := c.repo.Verify(ctx, biz, phone, inputCode)
	// 这里我们在 service 层面上对 Handler 屏蔽了最为特殊的错误
	if err == repository.ErrCodeVerifyTooManyTimes {
		// 在接入了告警之后，这边要告警
		// 因为这意味着有人在搞你
		return false, nil
	}
	return ok, err
}

func (c *SMSCodeService) generate() string {
	// 用随机数生成一个
	num := rand.Intn(999999)
	return fmt.Sprintf("%6d", num)
}
