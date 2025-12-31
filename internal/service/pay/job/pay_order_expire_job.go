package job

import (
	"context"

	"github.com/wxlbd/ruoyi-mall-go/internal/service/pay"
)

// PayOrderExpireJob 支付订单过期清理 Job
type PayOrderExpireJob struct {
	payOrderService *pay.PayOrderService
}

func NewPayOrderExpireJob(payOrderService *pay.PayOrderService) *PayOrderExpireJob {
	return &PayOrderExpireJob{
		payOrderService: payOrderService,
	}
}

func (j *PayOrderExpireJob) Execute(ctx context.Context, param string) error {
	_, err := j.payOrderService.ExpireOrder(ctx)
	return err
}

func (j *PayOrderExpireJob) GetHandlerName() string {
	return "payOrderExpireJob"
}
