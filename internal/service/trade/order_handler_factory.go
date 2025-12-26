package trade

import (
	"fmt"
	"sync"

	"go.uber.org/zap"
)

// OrderHandlerFactory 订单处理器工厂
type OrderHandlerFactory struct {
	handlers map[string]OrderHandler // 处理器映射
	logger   *zap.Logger
	mu       sync.RWMutex
}

// NewOrderHandlerFactory 创建订单处理器工厂
func NewOrderHandlerFactory(logger *zap.Logger) *OrderHandlerFactory {
	return &OrderHandlerFactory{
		handlers: make(map[string]OrderHandler),
		logger:   logger,
	}
}

// RegisterHandler 注册处理器
func (f *OrderHandlerFactory) RegisterHandler(handler OrderHandler) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	handlerType := handler.GetHandlerType()
	if _, exists := f.handlers[handlerType]; exists {
		return fmt.Errorf("处理器类型 %s 已存在", handlerType)
	}

	f.handlers[handlerType] = handler
	f.logger.Info("注册订单处理器",
		zap.String("handlerType", handlerType),
	)

	return nil
}

// RegisterHandlers 批量注册处理器
func (f *OrderHandlerFactory) RegisterHandlers(handlers []OrderHandler) error {
	for _, handler := range handlers {
		if err := f.RegisterHandler(handler); err != nil {
			return err
		}
	}
	return nil
}

// GetHandler 获取处理器
func (f *OrderHandlerFactory) GetHandler(handlerType string) (OrderHandler, error) {
	f.mu.RLock()
	defer f.mu.RUnlock()

	handler, exists := f.handlers[handlerType]
	if !exists {
		return nil, fmt.Errorf("处理器类型 %s 未找到", handlerType)
	}

	return handler, nil
}

// GetHandlerByOperation 根据操作类型获取处理器
func (f *OrderHandlerFactory) GetHandlerByOperation(operation string) (OrderHandler, error) {
	f.mu.RLock()
	defer f.mu.RUnlock()

	for _, handler := range f.handlers {
		if handler.CanHandle(operation) {
			return handler, nil
		}
	}

	return nil, fmt.Errorf("操作类型 %s 没有对应的处理器", operation)
}

// ListHandlers 列出所有处理器
func (f *OrderHandlerFactory) ListHandlers() []string {
	f.mu.RLock()
	defer f.mu.RUnlock()

	var handlerTypes []string
	for handlerType := range f.handlers {
		handlerTypes = append(handlerTypes, handlerType)
	}

	return handlerTypes
}

// GetHandlers 获取所有处理器
func (f *OrderHandlerFactory) GetHandlers() []OrderHandler {
	f.mu.RLock()
	defer f.mu.RUnlock()

	var handlers []OrderHandler
	for _, handler := range f.handlers {
		handlers = append(handlers, handler)
	}
	return handlers
}

// GetHandlerCount 获取处理器数量
func (f *OrderHandlerFactory) GetHandlerCount() int {
	f.mu.RLock()
	defer f.mu.RUnlock()

	return len(f.handlers)
}
