package calculators

import (
	"github.com/google/wire"
)

// ProviderSet 价格计算器 Provider Set
var ProviderSet = wire.NewSet(
	NewBargainActivityPriceCalculator,
	NewCombinationActivityPriceCalculator,
	NewCouponPriceCalculator,
	NewDeliveryPriceCalculator,
	NewDiscountActivityPriceCalculator,
	NewPointActivityPriceCalculator,
	NewPointGivePriceCalculator,
	NewPointUsePriceCalculator,
	NewRewardActivityPriceCalculator,
	NewSeckillActivityPriceCalculator,
)
