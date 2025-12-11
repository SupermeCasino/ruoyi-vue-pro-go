package trade

import (
	"backend-go/internal/model/trade"
	"backend-go/internal/repo" // aggregated repo
	"context"
)

type TradeOrderLogService struct {
	repo *repo.TradeOrderLogRepository
}

func NewTradeOrderLogService(repo *repo.TradeOrderLogRepository) *TradeOrderLogService {
	return &TradeOrderLogService{
		repo: repo,
	}
}

// CreateOrderLog 创建交易订单日志
func (s *TradeOrderLogService) CreateOrderLog(ctx context.Context, order *trade.TradeOrder, userType int, operatorId int64, content string) error {
	log := &trade.TradeOrderLog{
		UserID:      operatorId,
		UserType:    userType,
		OrderID:     order.ID,
		AfterStatus: order.Status, // Assume current status is after status
		Content:     content,
		OperateType: 0, // Default or pass in
	}
	// Enrich log details if needed
	return s.repo.Create(ctx, log)
}

// GetOrderLogListByOrderId 获得交易订单日志列表
func (s *TradeOrderLogService) GetOrderLogListByOrderId(ctx context.Context, orderId int64) ([]*trade.TradeOrderLog, error) {
	return s.repo.GetListByOrderId(ctx, orderId)
}
