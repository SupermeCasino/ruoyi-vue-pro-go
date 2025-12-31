package job

import (
	"context"
	"time"

	"github.com/wxlbd/ruoyi-mall-go/internal/service/pay"
)

// PayOrderSyncJob 支付订单同步 Job
type PayOrderSyncJob struct {
	payOrderService *pay.PayOrderService
}

func NewPayOrderSyncJob(payOrderService *pay.PayOrderService) *PayOrderSyncJob {
	return &PayOrderSyncJob{
		payOrderService: payOrderService,
	}
}

func (j *PayOrderSyncJob) Execute(ctx context.Context, param string) error {
	// 对齐 Java: Duration.ofMinutes(10)
	minCreateTime := time.Now().Add(-10 * time.Minute)
	_, err := j.payOrderService.SyncOrder(ctx, minCreateTime)
	return err
}

func (j *PayOrderSyncJob) GetHandlerName() string {
	return "payOrderSyncJob"
}
