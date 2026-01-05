package gateway

import (
	"github.com/google/wire"
	"github.com/wxlbd/ruoyi-mall-go/internal/iot/core"
	"github.com/wxlbd/ruoyi-mall-go/internal/iot/gateway/codec"
)

// ProviderSet IOT Gateway 模块依赖注入
var ProviderSet = wire.NewSet(
	// Codec Registry
	codec.DefaultRegistry,

	// Connection Manager Provider
	ProvideConnectionManager,

	// MQTT Client Provider (需要配置，暂时提供工厂函数) TODO: 集成依赖注入
	ProvideMQTTClientConfig,

	// Message Subscribers
	NewDeviceMessageSubscriber,
)

// ProvideConnectionManager 提供连接管理器
func ProvideConnectionManager(messageBus core.MessageBus) *ConnectionManager {
	// 默认 ServerID，实际应从配置读取
	return NewConnectionManager(messageBus, "iot-gateway-01")
}

// ProvideMQTTClientConfig 提供 MQTT 客户端配置（默认值）
func ProvideMQTTClientConfig() *MQTTClientConfig {
	return &MQTTClientConfig{
		Broker:           "tcp://localhost:1883",
		ClientID:         "iot-gateway-go",
		Username:         "",
		Password:         "",
		KeepAlive:        60,
		ConnectTimeout:   10,
		AutoReconnect:    true,
		CleanSession:     true,
		SubscribeTopics:  []string{"/sys/#"},
		DefaultCodecType: "Alink",
	}
}
