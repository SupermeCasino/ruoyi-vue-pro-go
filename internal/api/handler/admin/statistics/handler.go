package statistics

import "github.com/google/wire"

var ProviderSet = wire.NewSet(
	NewMemberStatisticsHandler,
	NewPayStatisticsHandler,
	NewProductStatisticsHandler,
	NewTradeStatisticsHandler,
	NewHandlers,
)

type Handlers struct {
	Member  *MemberStatisticsHandler
	Pay     *PayStatisticsHandler
	Product *ProductStatisticsHandler
	Trade   *TradeStatisticsHandler
}

func NewHandlers(
	member *MemberStatisticsHandler,
	pay *PayStatisticsHandler,
	product *ProductStatisticsHandler,
	trade *TradeStatisticsHandler,
) *Handlers {
	return &Handlers{
		Member:  member,
		Pay:     pay,
		Product: product,
		Trade:   trade,
	}
}
