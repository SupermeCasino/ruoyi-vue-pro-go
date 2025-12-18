package service

import (
	"context"

	"github.com/wxlbd/ruoyi-mall-go/internal/service/pay"
	"go.uber.org/zap"
)

// PayTransferSyncJob 这里的名字要和 Java 保持一致: payTransferSyncJob
type PayTransferSyncJob struct {
	transferService *pay.PayTransferService
	logger          *zap.Logger
}

func NewPayTransferSyncJob(transferService *pay.PayTransferService, logger *zap.Logger) *PayTransferSyncJob {
	return &PayTransferSyncJob{
		transferService: transferService,
		logger:          logger,
	}
}

func (j *PayTransferSyncJob) Execute(ctx context.Context, param string) error {
	count, err := j.transferService.SyncTransfer(ctx)
	if err != nil {
		return err
	}
	j.logger.Info("同步转账单完成", zap.Int("count", count))
	return nil
}

func (j *PayTransferSyncJob) GetHandlerName() string {
	return "payTransferSyncJob"
}
