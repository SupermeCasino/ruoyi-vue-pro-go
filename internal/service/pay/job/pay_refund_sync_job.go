package job

import (
	"context"

	"github.com/wxlbd/ruoyi-mall-go/internal/service/pay"
)

// PayRefundSyncJob 支付退款同步 Job
type PayRefundSyncJob struct {
	payRefundService *pay.PayRefundService
}

func NewPayRefundSyncJob(payRefundService *pay.PayRefundService) *PayRefundSyncJob {
	return &PayRefundSyncJob{
		payRefundService: payRefundService,
	}
}

func (j *PayRefundSyncJob) Execute(ctx context.Context, param string) error {
	_, err := j.payRefundService.SyncRefund(ctx)
	return err
}

func (j *PayRefundSyncJob) GetHandlerName() string {
	return "payRefundSyncJob"
}
