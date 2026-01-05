package core

import (
	"sync"
)

// MessageBus IoT 消息总线接口
// 用于在 IoT 系统中发布和订阅消息
type MessageBus interface {
	// Post 发布消息到指定主题
	Post(topic string, message interface{})

	// Register 注册消息订阅者
	Register(subscriber MessageSubscriber)
}

// MessageSubscriber 消息订阅者接口
type MessageSubscriber interface {
	// Topic 返回订阅的主题
	Topic() string

	// Group 返回订阅者分组（用于负载均衡）
	Group() string

	// OnMessage 处理消息
	OnMessage(message any)
}

// LocalMessageBus 基于 Go channel 的本地高性能消息总线实现
type LocalMessageBus struct {
	mu          sync.RWMutex
	subscribers map[string][]MessageSubscriber
}

// NewLocalMessageBus 创建本地消息总线
func NewLocalMessageBus() *LocalMessageBus {
	return &LocalMessageBus{
		subscribers: make(map[string][]MessageSubscriber),
	}
}

// Post 发布消息到指定主题
func (b *LocalMessageBus) Post(topic string, message any) {
	b.mu.RLock()
	subs := b.subscribers[topic]
	b.mu.RUnlock()

	// 异步分发消息给所有订阅者
	for _, sub := range subs {
		go sub.OnMessage(message)
	}
}

// Register 注册消息订阅者
func (b *LocalMessageBus) Register(subscriber MessageSubscriber) {
	b.mu.Lock()
	defer b.mu.Unlock()

	topic := subscriber.Topic()
	b.subscribers[topic] = append(b.subscribers[topic], subscriber)
}

var _ MessageBus = (*LocalMessageBus)(nil)
