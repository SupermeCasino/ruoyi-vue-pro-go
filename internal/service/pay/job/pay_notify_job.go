package job

import (
	"context"

	"github.com/wxlbd/ruoyi-mall-go/internal/service/pay"
)

// PayNotifyJob 支付通知 Job
// 通过不断扫描待通知的 PayNotifyTaskDO 记录，回调业务线的回调接口
type PayNotifyJob struct {
	payNotifyService *pay.PayNotifyService
}

func NewPayNotifyJob(payNotifyService *pay.PayNotifyService) *PayNotifyJob {
	return &PayNotifyJob{
		payNotifyService: payNotifyService,
	}
}

func (j *PayNotifyJob) Execute(ctx context.Context, param string) error {
	_, err := j.payNotifyService.ExecuteNotify(ctx)
	return err
}

func (j *PayNotifyJob) GetHandlerName() string {
	return "payNotifyJob"
}
