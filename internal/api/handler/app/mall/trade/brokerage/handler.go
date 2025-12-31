package brokerage

import "github.com/google/wire"

var ProviderSet = wire.NewSet(
	NewAppBrokerageRecordHandler,
	NewAppBrokerageUserHandler,
	NewAppBrokerageWithdrawHandler,
	NewHandlers,
)

type Handlers struct {
	BrokerageRecord   *AppBrokerageRecordHandler
	BrokerageUser     *AppBrokerageUserHandler
	BrokerageWithdraw *AppBrokerageWithdrawHandler
}

func NewHandlers(
	brokerageRecord *AppBrokerageRecordHandler,
	brokerageUser *AppBrokerageUserHandler,
	brokerageWithdraw *AppBrokerageWithdrawHandler,
) *Handlers {
	return &Handlers{
		BrokerageRecord:   brokerageRecord,
		BrokerageUser:     brokerageUser,
		BrokerageWithdraw: brokerageWithdraw,
	}
}
