package trade

import (
	"context"

	"go.uber.org/zap"
)

// PriceCalculator 价格计算器接口
// 使用策略模式实现不同类型的价格计算逻辑
type PriceCalculator interface {
	// Calculate 执行价格计算
	Calculate(ctx context.Context, req *TradePriceCalculateReqBO, resp *TradePriceCalculateRespBO) error

	// GetOrder 获取计算器执行顺序（优先级）
	GetOrder() int

	// GetName 获取计算器名称
	GetName() string

	// IsApplicable 判断是否适用于当前订单类型
	IsApplicable(orderType int) bool
}

const (
	// PromotionTypeCoupon 优惠券
	PromotionTypeCoupon = 40
	// PromotionTypeCombination 拼团活动
	PromotionTypeCombination = 50
	// PromotionTypeBargain 砍价活动
	PromotionTypeBargain = 60
	// PromotionTypeSeckill 秒杀活动
	PromotionTypeSeckill = 70
	// PromotionTypePoint 积分抵扣
	PromotionTypePoint = 80
)

// PriceCalculatorFactory 价格计算器工厂
type PriceCalculatorFactory struct {
	calculators []PriceCalculator
	logger      *zap.Logger
}

// NewPriceCalculatorFactory 创建价格计算器工厂
func NewPriceCalculatorFactory(logger *zap.Logger) *PriceCalculatorFactory {
	return &PriceCalculatorFactory{
		calculators: make([]PriceCalculator, 0),
		logger:      logger,
	}
}

// RegisterCalculator 注册价格计算器
func (f *PriceCalculatorFactory) RegisterCalculator(calculator PriceCalculator) {
	f.calculators = append(f.calculators, calculator)
	f.logger.Info("注册价格计算器",
		zap.String("name", calculator.GetName()),
		zap.Int("order", calculator.GetOrder()),
	)
}

// GetCalculators 获取所有已注册的计算器，按优先级排序
func (f *PriceCalculatorFactory) GetCalculators() []PriceCalculator {
	// 按优先级排序（数字越小优先级越高）
	calculators := make([]PriceCalculator, len(f.calculators))
	copy(calculators, f.calculators)

	// 使用冒泡排序按优先级排序
	for i := 0; i < len(calculators)-1; i++ {
		for j := 0; j < len(calculators)-1-i; j++ {
			if calculators[j].GetOrder() > calculators[j+1].GetOrder() {
				calculators[j], calculators[j+1] = calculators[j+1], calculators[j]
			}
		}
	}

	return calculators
}

// GetApplicableCalculators 获取适用于指定订单类型的计算器
func (f *PriceCalculatorFactory) GetApplicableCalculators(orderType int) []PriceCalculator {
	allCalculators := f.GetCalculators()
	applicable := make([]PriceCalculator, 0)

	for _, calculator := range allCalculators {
		if calculator.IsApplicable(orderType) {
			applicable = append(applicable, calculator)
		}
	}

	return applicable
}

// BasePriceCalculator 基础价格计算器，提供通用功能
type BasePriceCalculator struct {
	name   string
	order  int
	Helper *PriceCalculatorHelper
	logger *zap.Logger
}

// NewBasePriceCalculator 创建基础价格计算器
func NewBasePriceCalculator(name string, order int, helper *PriceCalculatorHelper, logger *zap.Logger) *BasePriceCalculator {
	return &BasePriceCalculator{
		name:   name,
		order:  order,
		Helper: helper,
		logger: logger,
	}
}

// GetName 获取计算器名称
func (b *BasePriceCalculator) GetName() string {
	return b.name
}

// GetOrder 获取计算器执行顺序
func (b *BasePriceCalculator) GetOrder() int {
	return b.order
}

// LogCalculation 记录计算过程日志
func (b *BasePriceCalculator) LogCalculation(ctx context.Context, req *TradePriceCalculateReqBO, message string, fields ...zap.Field) {
	logFields := []zap.Field{
		zap.String("calculator", b.name),
		zap.Int64("userId", req.UserID),
		zap.Int("itemCount", len(req.Items)),
		zap.String("message", message),
	}
	logFields = append(logFields, fields...)

	b.logger.Info("价格计算器执行", logFields...)
}

// LogError 记录错误日志
func (b *BasePriceCalculator) LogError(ctx context.Context, req *TradePriceCalculateReqBO, err error, message string) {
	b.logger.Error("价格计算器执行失败",
		zap.String("calculator", b.name),
		zap.Int64("userId", req.UserID),
		zap.Int("itemCount", len(req.Items)),
		zap.String("message", message),
		zap.Error(err),
	)
}
