package gateway

import (
	"time"

	"github.com/google/wire"
	"github.com/wxlbd/ruoyi-mall-go/internal/iot/core"
	"github.com/wxlbd/ruoyi-mall-go/internal/iot/gateway/codec"
	"github.com/wxlbd/ruoyi-mall-go/pkg/config"
)

// ProviderSet IOT Gateway 模块依赖注入
var ProviderSet = wire.NewSet(
	// Codec Registry
	codec.DefaultRegistry,

	// Connection Manager Provider
	ProvideConnectionManager,

	// MQTT Client Provider
	ProvideMQTTClientConfig,

	// Message Subscribers
	NewDeviceMessageSubscriber,
)

// ProvideConnectionManager 提供连接管理器
func ProvideConnectionManager(messageBus core.MessageBus) *ConnectionManager {
	cfg := config.C.IoT.Gateway
	// 默认 ServerID
	serverID := cfg.ServerID
	if serverID == "" {
		serverID = "iot-gateway-default"
	}
	manager := NewConnectionManager(messageBus, serverID)
	if cfg.MaxConnections > 0 {
		manager.SetMaxConnections(cfg.MaxConnections)
	}
	return manager
}

// ProvideMQTTClientConfig 提供 MQTT 客户端配置（默认值）
func ProvideMQTTClientConfig() *MQTTClientConfig {
	cfg := config.C.IoT.Gateway.MQTT

	keepAlive, _ := time.ParseDuration(cfg.KeepAlive)
	if keepAlive == 0 {
		keepAlive = 60 * time.Second
	}

	connectTimeout, _ := time.ParseDuration(cfg.ConnectTimeout)
	if connectTimeout == 0 {
		connectTimeout = 10 * time.Second
	}

	return &MQTTClientConfig{
		Broker:           cfg.Broker,
		ClientID:         cfg.ClientID,
		Username:         cfg.Username,
		Password:         cfg.Password,
		KeepAlive:        keepAlive,
		ConnectTimeout:   connectTimeout,
		AutoReconnect:    cfg.AutoReconnect,
		CleanSession:     cfg.CleanSession,
		SubscribeTopics:  cfg.SubscribeTopics,
		DefaultCodecType: cfg.DefaultCodecType,
		TopicPrefix:      cfg.TopicPrefix,
	}
}
