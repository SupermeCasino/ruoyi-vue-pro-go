package trade

import (
	"context"

	"github.com/wxlbd/ruoyi-mall-go/internal/model/trade"
	"github.com/wxlbd/ruoyi-mall-go/internal/repo"
)

type AfterSaleLogService struct {
	repo *repo.AfterSaleLogRepository
}

func NewAfterSaleLogService(repo *repo.AfterSaleLogRepository) *AfterSaleLogService {
	return &AfterSaleLogService{
		repo: repo,
	}
}

// CreateAfterSaleLog 创建售后日志
func (s *AfterSaleLogService) CreateAfterSaleLog(ctx context.Context, logDTO *trade.AfterSaleLog) error {
	return s.repo.Create(ctx, logDTO)
}

// GetAfterSaleLogList 获得售后日志列表
func (s *AfterSaleLogService) GetAfterSaleLogList(ctx context.Context, afterSaleId int64) ([]*trade.AfterSaleLog, error) {
	return s.repo.GetListByAfterSaleId(ctx, afterSaleId)
}
