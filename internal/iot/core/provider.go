package core

import (
	"github.com/google/wire"
	"github.com/wxlbd/ruoyi-mall-go/pkg/config"
)

// ProviderSet IOT Core 模块依赖注入
var ProviderSet = wire.NewSet(
	// MessageBus
	ProvideLocalMessageBusConfig,
	NewLocalMessageBus,
	wire.Bind(new(MessageBus), new(*LocalMessageBus)),

	// MQ Producer Factory
	NewMQProducerFactory,

	// Device Auth
	NewDeviceAuthUtils,
)

// ProvideLocalMessageBusConfig 提供消息总线默认配置
func ProvideLocalMessageBusConfig() LocalMessageBusConfig {
	cfg := config.C.IoT.Core.MessageBus
	return LocalMessageBusConfig{
		WorkerNum: cfg.WorkerNum,
		QueueSize: cfg.QueueSize,
	}
}
