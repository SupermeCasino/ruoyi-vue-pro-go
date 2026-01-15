package core

import (
	"context"
	"log"
	"sync"
)

// MessageBus IoT 消息总线接口
// 用于在 IoT 系统中发布和订阅消息
type MessageBus interface {
	// Start 启动消息总线
	Start()
	// Stop 停止消息总线
	Stop()
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

// messageTask 内部消息任务
type messageTask struct {
	topic   string
	payload any
}

// LocalMessageBus 基于 Worker Pool 的本地高性能消息总线实现
type LocalMessageBus struct {
	mu          sync.RWMutex
	subscribers map[string][]MessageSubscriber
	taskChan    chan messageTask
	workerNum   int
	ctx         context.Context
	cancel      context.CancelFunc
	wg          sync.WaitGroup
}

// LocalMessageBusConfig 配置
type LocalMessageBusConfig struct {
	WorkerNum int // 工作协程数量
	QueueSize int // 缓冲队列大小
}

// NewLocalMessageBus 创建本地消息总线
func NewLocalMessageBus(config LocalMessageBusConfig) *LocalMessageBus {
	if config.WorkerNum <= 0 {
		config.WorkerNum = 10 // 默认值
	}
	if config.QueueSize <= 0 {
		config.QueueSize = 1000 // 默认值
	}

	ctx, cancel := context.WithCancel(context.Background())
	return &LocalMessageBus{
		subscribers: make(map[string][]MessageSubscriber),
		taskChan:    make(chan messageTask, config.QueueSize),
		workerNum:   config.WorkerNum,
		ctx:         ctx,
		cancel:      cancel,
	}
}

// Start 启动 Worker Pool
func (b *LocalMessageBus) Start() {
	for i := 0; i < b.workerNum; i++ {
		b.wg.Add(1)
		go b.worker()
	}
	log.Printf("[MessageBus] Started with %d workers", b.workerNum)
}

// Stop 停止消息总线
func (b *LocalMessageBus) Stop() {
	b.cancel()
	b.wg.Wait()
	close(b.taskChan)
	log.Printf("[MessageBus] Stopped")
}

// Post 发布消息到指定主题
func (b *LocalMessageBus) Post(topic string, message any) {
	select {
	case b.taskChan <- messageTask{topic: topic, payload: message}:
		// 消息入队成功
	default:
		// 队列已满，记录错误或采取丢弃策略
		log.Printf("[MessageBus] Queue full, dropping message for topic: %s", topic)
	}
}

// Register 注册消息订阅者
func (b *LocalMessageBus) Register(subscriber MessageSubscriber) {
	b.mu.Lock()
	defer b.mu.Unlock()

	topic := subscriber.Topic()
	b.subscribers[topic] = append(b.subscribers[topic], subscriber)
}

// worker 工作协程
func (b *LocalMessageBus) worker() {
	defer b.wg.Done()

	for {
		select {
		case <-b.ctx.Done():
			// 处理剩余任务（可选）
			return
		case task := <-b.taskChan:
			b.process(task)
		}
	}
}

// process 处理单个消息任务
func (b *LocalMessageBus) process(task messageTask) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[MessageBus] Panic recovered in worker: %v", r)
		}
	}()

	b.mu.RLock()
	subs := b.subscribers[task.topic]
	b.mu.RUnlock()

	for _, sub := range subs {
		// 仍然在同一个 Worker 中串行执行，确保单个 Worker 不会被阻塞太久
		// 如果订阅者处理逻辑很重，应该自己开协程
		sub.OnMessage(task.payload)
	}
}

var _ MessageBus = (*LocalMessageBus)(nil)
