package trade

import (
	"context"
	"strings"

	"github.com/wxlbd/ruoyi-mall-go/internal/api/req"
	"github.com/wxlbd/ruoyi-mall-go/internal/api/resp"
	"github.com/wxlbd/ruoyi-mall-go/internal/model"
	"github.com/wxlbd/ruoyi-mall-go/internal/model/trade"
	"github.com/wxlbd/ruoyi-mall-go/internal/repo/query"
)

type TradeConfigService struct {
	q *query.Query
}

func NewTradeConfigService(q *query.Query) *TradeConfigService {
	return &TradeConfigService{q: q}
}

// GetTradeConfig 获取交易配置 (Admin/App)
func (s *TradeConfigService) GetTradeConfig(ctx context.Context) (*resp.TradeConfigResp, error) {
	qc := s.q.TradeConfig
	config, err := qc.WithContext(ctx).First()
	if err != nil {
		// 如果不存在，返回默认配置或空
		return &resp.TradeConfigResp{}, nil // 或者返回默认值
	}

	return &resp.TradeConfigResp{
		ID:                          config.ID,
		AppID:                       config.AppID,
		AfterSaleDeadlineDays:       config.AfterSaleDeadlineDays,
		PayTimeoutMinutes:           config.PayTimeoutMinutes,
		AutoReceiveDays:             config.AutoReceiveDays,
		AutoCommentDays:             config.AutoCommentDays,
		BrokerageWithdrawMinPrice:   config.BrokerageWithdrawMinPrice,
		BrokerageWithdrawFeePercent: config.BrokerageWithdrawFeePercent,
		BrokerageEnabled:            bool(config.BrokerageEnabled),
		BrokerageFrozenDays:         config.BrokerageFrozenDays,
		BrokerageFirstPercent:       config.BrokerageFirstPercent,
		BrokerageSecondPercent:      config.BrokerageSecondPercent,
		BrokeragePosterUrls: func() []string {
			if config.BrokeragePosterUrls == "" {
				return []string{}
			}
			return strings.Split(config.BrokeragePosterUrls, ",")
		}(),
	}, nil
}

// SaveTradeConfig 保存交易配置 (Admin)
func (s *TradeConfigService) SaveTradeConfig(ctx context.Context, r *req.TradeConfigSaveReq) error {
	qc := s.q.TradeConfig
	existing, err := qc.WithContext(ctx).First()
	if err == nil {
		// Update
		existing.AfterSaleDeadlineDays = *r.AfterSaleDeadlineDays
		existing.PayTimeoutMinutes = *r.PayTimeoutMinutes
		existing.AutoReceiveDays = *r.AutoReceiveDays
		existing.AutoCommentDays = *r.AutoCommentDays
		if r.BrokerageWithdrawMinPrice != nil {
			existing.BrokerageWithdrawMinPrice = *r.BrokerageWithdrawMinPrice
		}
		if r.BrokerageWithdrawFeePercent != nil {
			existing.BrokerageWithdrawFeePercent = *r.BrokerageWithdrawFeePercent
		}
		if r.BrokerageEnabled != nil {
			existing.BrokerageEnabled = model.BitBool(*r.BrokerageEnabled)
		}
		if r.BrokerageFrozenDays != nil {
			existing.BrokerageFrozenDays = *r.BrokerageFrozenDays
		}
		if r.BrokerageFirstPercent != nil {
			existing.BrokerageFirstPercent = *r.BrokerageFirstPercent
		}
		if r.BrokerageSecondPercent != nil {
			existing.BrokerageSecondPercent = *r.BrokerageSecondPercent
		}
		if r.BrokeragePosterUrls != nil {
			existing.BrokeragePosterUrls = strings.Join(r.BrokeragePosterUrls, ",")
		}
		return qc.WithContext(ctx).Save(existing)
	}

	// Create
	newConfig := &trade.TradeConfig{
		AfterSaleDeadlineDays:       *r.AfterSaleDeadlineDays,
		PayTimeoutMinutes:           *r.PayTimeoutMinutes,
		AutoReceiveDays:             *r.AutoReceiveDays,
		AutoCommentDays:             *r.AutoCommentDays,
		BrokerageWithdrawMinPrice:   0,
		BrokerageWithdrawFeePercent: 0,
		BrokerageEnabled:            false,
		BrokerageFrozenDays:         0,
		BrokerageFirstPercent:       0,
		BrokerageSecondPercent:      0,
		BrokeragePosterUrls:         "",
	}
	if r.BrokerageWithdrawMinPrice != nil {
		newConfig.BrokerageWithdrawMinPrice = *r.BrokerageWithdrawMinPrice
	}
	if r.BrokerageWithdrawFeePercent != nil {
		newConfig.BrokerageWithdrawFeePercent = *r.BrokerageWithdrawFeePercent
	}
	if r.BrokerageEnabled != nil {
		newConfig.BrokerageEnabled = model.BitBool(*r.BrokerageEnabled)
	}
	if r.BrokerageFrozenDays != nil {
		newConfig.BrokerageFrozenDays = *r.BrokerageFrozenDays
	}
	if r.BrokerageFirstPercent != nil {
		newConfig.BrokerageFirstPercent = *r.BrokerageFirstPercent
	}
	if r.BrokerageSecondPercent != nil {
		newConfig.BrokerageSecondPercent = *r.BrokerageSecondPercent
	}
	if r.BrokeragePosterUrls != nil {
		newConfig.BrokeragePosterUrls = strings.Join(r.BrokeragePosterUrls, ",")
	}
	return qc.WithContext(ctx).Create(newConfig)
}
