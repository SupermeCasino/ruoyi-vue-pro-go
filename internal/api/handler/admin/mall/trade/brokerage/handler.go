package brokerage

import "github.com/google/wire"

var ProviderSet = wire.NewSet(
	NewBrokerageRecordHandler,
	NewBrokerageUserHandler,
	NewBrokerageWithdrawHandler,
	NewHandlers,
)

type Handlers struct {
	BrokerageRecord   *BrokerageRecordHandler
	BrokerageUser     *BrokerageUserHandler
	BrokerageWithdraw *BrokerageWithdrawHandler
}

func NewHandlers(
	brokerageRecord *BrokerageRecordHandler,
	brokerageUser *BrokerageUserHandler,
	brokerageWithdraw *BrokerageWithdrawHandler,
) *Handlers {
	return &Handlers{
		BrokerageRecord:   brokerageRecord,
		BrokerageUser:     brokerageUser,
		BrokerageWithdraw: brokerageWithdraw,
	}
}
