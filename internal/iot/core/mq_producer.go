package core

import (
	"fmt"
)

// MQProducer MQ 消息生产者接口
// 用于将消息发送到外部消息队列（Kafka/RocketMQ/NSQ 等）
type MQProducer interface {
	// Publish 发布消息到指定主题
	Publish(topic string, message []byte) error

	// Close 关闭生产者
	Close() error
}

// MQType MQ 类型枚举
type MQType string

const (
	MQTypeKafka    MQType = "kafka"
	MQTypeRocketMQ MQType = "rocketmq"
	MQTypeNSQ      MQType = "nsq"
	MQTypeLocal    MQType = "local" // 本地模式，不实际发送
)

// MQConfig MQ 配置
type MQConfig struct {
	Type    MQType   `yaml:"type"`
	Brokers []string `yaml:"brokers"`
	Topic   string   `yaml:"topic"`
}

// MQProducerFactory MQ 生产者工厂
// 根据配置动态创建对应类型的 MQ 生产者
type MQProducerFactory struct{}

// NewMQProducerFactory 创建 MQ 生产者工厂
func NewMQProducerFactory() *MQProducerFactory {
	return &MQProducerFactory{}
}

// Create 根据配置创建 MQ 生产者
func (f *MQProducerFactory) Create(config *MQConfig) (MQProducer, error) {
	switch config.Type {
	case MQTypeLocal:
		return NewLocalMQProducer(), nil
	case MQTypeKafka:
		// Kafka 生产者待实现
		return nil, fmt.Errorf("kafka producer not implemented yet")
	case MQTypeRocketMQ:
		// RocketMQ 生产者待实现
		return nil, fmt.Errorf("rocketmq producer not implemented yet")
	case MQTypeNSQ:
		// NSQ 生产者待实现
		return nil, fmt.Errorf("nsq producer not implemented yet")
	default:
		return nil, fmt.Errorf("unsupported mq type: %s", config.Type)
	}
}

// LocalMQProducer 本地 MQ 生产者（用于开发/测试，不实际发送）
type LocalMQProducer struct{}

// NewLocalMQProducer 创建本地 MQ 生产者
func NewLocalMQProducer() *LocalMQProducer {
	return &LocalMQProducer{}
}

// Publish 发布消息（本地模式仅打印日志）
func (p *LocalMQProducer) Publish(topic string, message []byte) error {
	// 本地模式：仅记录日志，不实际发送
	// log.Printf("[LocalMQProducer] topic=%s, message=%s", topic, string(message))
	return nil
}

// Close 关闭生产者
func (p *LocalMQProducer) Close() error {
	return nil
}

var _ MQProducer = (*LocalMQProducer)(nil)
