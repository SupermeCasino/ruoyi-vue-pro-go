package trade

import (
	"context"

	"github.com/wxlbd/ruoyi-mall-go/internal/api/req"
	"github.com/wxlbd/ruoyi-mall-go/internal/api/resp"
	"github.com/wxlbd/ruoyi-mall-go/internal/consts"
	"github.com/wxlbd/ruoyi-mall-go/internal/model"
	"github.com/wxlbd/ruoyi-mall-go/internal/model/trade"
	"github.com/wxlbd/ruoyi-mall-go/internal/repo/query"
	"gorm.io/datatypes"
)

type TradeConfigService struct {
	q *query.Query
}

func NewTradeConfigService(q *query.Query) *TradeConfigService {
	return &TradeConfigService{q: q}
}

// GetTradeConfig 获取交易配置 (Admin)
func (s *TradeConfigService) GetTradeConfig(ctx context.Context) (*resp.TradeConfigResp, error) {
	qc := s.q.TradeConfig
	config, err := qc.WithContext(ctx).First()
	if err != nil {
		// 如果不存在，返回默认配置（对齐 Java application.yaml / 默认行为）
		return &resp.TradeConfigResp{
			AfterSaleDeadlineDays: consts.DefaultAfterSaleDeadlineDays,
			PayTimeoutMinutes:     consts.DefaultPayTimeoutMinutes,
			AutoReceiveDays:       consts.DefaultAutoReceiveDays,
			AutoCommentDays:       consts.DefaultAutoCommentDays,
		}, nil
	}

	// 补充部分默认值，防止数据库中记录存在但值为 0 的情况
	res := &resp.TradeConfigResp{
		ID:                          config.ID,
		AppID:                       config.AppID,
		AfterSaleDeadlineDays:       config.AfterSaleDeadlineDays,
		PayTimeoutMinutes:           config.PayTimeoutMinutes,
		AutoReceiveDays:             config.AutoReceiveDays,
		AutoCommentDays:             config.AutoCommentDays,
		AfterSaleRefundReasons:      []string(config.AfterSaleRefundReasons),
		AfterSaleReturnReasons:      []string(config.AfterSaleReturnReasons),
		DeliveryExpressFreeEnabled:  bool(config.DeliveryExpressFreeEnabled),
		DeliveryExpressFreePrice:    config.DeliveryExpressFreePrice,
		DeliveryPickUpEnabled:       bool(config.DeliveryPickUpEnabled),
		BrokerageWithdrawMinPrice:   config.BrokerageWithdrawMinPrice,
		BrokerageWithdrawFeePercent: config.BrokerageWithdrawFeePercent,
		BrokerageEnabled:            bool(config.BrokerageEnabled),
		BrokerageFrozenDays:         config.BrokerageFrozenDays,
		BrokerageFirstPercent:       config.BrokerageFirstPercent,
		BrokerageSecondPercent:      config.BrokerageSecondPercent,
		BrokerageEnabledCondition:   config.BrokerageEnabledCondition,
		BrokerageBindMode:           config.BrokerageBindMode,
		BrokeragePosterUrls:         []string(config.BrokeragePosterUrls),
		BrokerageWithdrawTypes:      []int(config.BrokerageWithdrawTypes),
	}
	if res.PayTimeoutMinutes <= consts.DefaultPayTimeoutMinutes {
		res.PayTimeoutMinutes = consts.DefaultPayTimeoutMinutes
	}
	if res.AfterSaleDeadlineDays <= consts.DefaultAfterSaleDeadlineDays {
		res.AfterSaleDeadlineDays = consts.DefaultAfterSaleDeadlineDays
	}
	return res, nil
}

// GetAppTradeConfig 获取交易配置 (App) - 对齐 Java: AppTradeConfigController.getTradeConfig
// 对应 Java: TradeConfigConvert.convert02(TradeConfigDO)
func (s *TradeConfigService) GetAppTradeConfig(ctx context.Context) (*resp.AppTradeConfigResp, error) {
	qc := s.q.TradeConfig
	config, err := qc.WithContext(ctx).First()
	if err != nil {
		// 如果不存在，返回默认配置（对齐 Java 的 ObjUtil.defaultIfNull）
		config = &trade.TradeConfig{}
	}

	// 转换响应结构（对齐 Java: TradeConfigConvert.convert02）
	return &resp.AppTradeConfigResp{
		TencentLbsKey:             "",
		DeliveryPickUpEnabled:     bool(config.DeliveryPickUpEnabled),
		AfterSaleRefundReasons:    []string(config.AfterSaleRefundReasons),
		AfterSaleReturnReasons:    []string(config.AfterSaleReturnReasons),
		BrokeragePosterUrls:       []string(config.BrokeragePosterUrls),
		BrokerageFrozenDays:       config.BrokerageFrozenDays,
		BrokerageWithdrawMinPrice: config.BrokerageWithdrawMinPrice,
		BrokerageWithdrawTypes:    []int(config.BrokerageWithdrawTypes),
	}, nil
}

// SaveTradeConfig 保存交易配置 (Admin) - 对齐 Java: TradeConfigService.saveTradeConfig
func (s *TradeConfigService) SaveTradeConfig(ctx context.Context, r *req.TradeConfigSaveReq) error {
	qc := s.q.TradeConfig
	existing, err := qc.WithContext(ctx).First()
	if err == nil {
		// Update
		existing.AfterSaleDeadlineDays = *r.AfterSaleDeadlineDays
		existing.PayTimeoutMinutes = *r.PayTimeoutMinutes
		existing.AutoReceiveDays = *r.AutoReceiveDays
		existing.AutoCommentDays = *r.AutoCommentDays
		if len(r.AfterSaleRefundReasons) > 0 {
			existing.AfterSaleRefundReasons = datatypes.JSONSlice[string](r.AfterSaleRefundReasons)
		}
		if len(r.AfterSaleReturnReasons) > 0 {
			existing.AfterSaleReturnReasons = datatypes.JSONSlice[string](r.AfterSaleReturnReasons)
		}
		if r.DeliveryExpressFreeEnabled != nil {
			existing.DeliveryExpressFreeEnabled = model.BitBool(*r.DeliveryExpressFreeEnabled)
		}
		if r.DeliveryExpressFreePrice != nil {
			existing.DeliveryExpressFreePrice = *r.DeliveryExpressFreePrice
		}
		if r.DeliveryPickUpEnabled != nil {
			existing.DeliveryPickUpEnabled = model.BitBool(*r.DeliveryPickUpEnabled)
		}
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
		if r.BrokerageEnabledCondition != nil {
			existing.BrokerageEnabledCondition = *r.BrokerageEnabledCondition
		}
		if r.BrokerageBindMode != nil {
			existing.BrokerageBindMode = *r.BrokerageBindMode
		}
		if len(r.BrokeragePosterUrls) > 0 {
			existing.BrokeragePosterUrls = datatypes.JSONSlice[string](r.BrokeragePosterUrls)
		}
		if len(r.BrokerageWithdrawTypes) > 0 {
			existing.BrokerageWithdrawTypes = r.BrokerageWithdrawTypes
		}
		return qc.WithContext(ctx).Save(existing)
	}

	// Create
	newConfig := &trade.TradeConfig{
		AfterSaleDeadlineDays:       *r.AfterSaleDeadlineDays,
		PayTimeoutMinutes:           *r.PayTimeoutMinutes,
		AutoReceiveDays:             *r.AutoReceiveDays,
		AutoCommentDays:             *r.AutoCommentDays,
		AfterSaleRefundReasons:      datatypes.JSONSlice[string](r.AfterSaleRefundReasons),
		AfterSaleReturnReasons:      datatypes.JSONSlice[string](r.AfterSaleReturnReasons),
		BrokerageWithdrawMinPrice:   0,
		BrokerageWithdrawFeePercent: 0,
		BrokerageEnabled:            false,
		BrokerageFrozenDays:         0,
		BrokerageFirstPercent:       0,
		BrokerageSecondPercent:      0,
		BrokerageEnabledCondition:   consts.BrokerageEnabledConditionAll,
		BrokerageBindMode:           consts.BrokerageBindModeAnytime,
	}
	if r.DeliveryExpressFreeEnabled != nil {
		newConfig.DeliveryExpressFreeEnabled = model.BitBool(*r.DeliveryExpressFreeEnabled)
	}
	if r.DeliveryExpressFreePrice != nil {
		newConfig.DeliveryExpressFreePrice = *r.DeliveryExpressFreePrice
	}
	if r.DeliveryPickUpEnabled != nil {
		newConfig.DeliveryPickUpEnabled = model.BitBool(*r.DeliveryPickUpEnabled)
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
	if r.BrokerageEnabledCondition != nil {
		newConfig.BrokerageEnabledCondition = *r.BrokerageEnabledCondition
	}
	if r.BrokerageBindMode != nil {
		newConfig.BrokerageBindMode = *r.BrokerageBindMode
	}
	if len(r.BrokeragePosterUrls) > 0 {
		newConfig.BrokeragePosterUrls = datatypes.JSONSlice[string](r.BrokeragePosterUrls)
	}
	if len(r.BrokerageWithdrawTypes) > 0 {
		newConfig.BrokerageWithdrawTypes = r.BrokerageWithdrawTypes
	}
	return qc.WithContext(ctx).Create(newConfig)
}
