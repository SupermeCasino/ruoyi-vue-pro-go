package core

import (
	"github.com/google/wire"
)

// ProviderSet IOT Core 模块依赖注入
var ProviderSet = wire.NewSet(
	// MessageBus
	NewLocalMessageBus,
	wire.Bind(new(MessageBus), new(*LocalMessageBus)),

	// MQ Producer Factory
	NewMQProducerFactory,

	// Device Auth
	NewDeviceAuthUtils,
)
