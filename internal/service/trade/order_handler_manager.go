package trade

import (
	"context"
	"fmt"

	"go.uber.org/zap"
)

// OrderHandlerManager 订单处理器管理器
type OrderHandlerManager struct {
	factory *OrderHandlerFactory
	router  *OrderHandlerRouter
	logger  *zap.Logger
}

// NewOrderHandlerManager 创建订单处理器管理器
func NewOrderHandlerManager(logger *zap.Logger) *OrderHandlerManager {
	factory := NewOrderHandlerFactory(logger)
	router := NewOrderHandlerRouter(factory, logger)

	return &OrderHandlerManager{
		factory: factory,
		router:  router,
		logger:  logger,
	}
}

// Initialize 初始化处理器管理器
func (m *OrderHandlerManager) Initialize(handlers []OrderHandler) error {
	if err := m.factory.RegisterHandlers(handlers); err != nil {
		return fmt.Errorf("注册订单处理器失败: %w", err)
	}

	m.logger.Info("订单处理器管理器初始化完成",
		zap.Int("handlerCount", m.factory.GetHandlerCount()),
		zap.Strings("handlerTypes", m.factory.ListHandlers()),
	)

	return nil
}

// HandleOrder 处理订单操作
func (m *OrderHandlerManager) HandleOrder(ctx context.Context, req *OrderHandleRequest) (*OrderHandleResponse, error) {
	return m.router.Route(ctx, req)
}

// GetFactory 获取处理器工厂
func (m *OrderHandlerManager) GetFactory() *OrderHandlerFactory {
	return m.factory
}

// GetRouter 获取处理器路由
func (m *OrderHandlerManager) GetRouter() *OrderHandlerRouter {
	return m.router
}

// GetHandlerInfo 获取处理器信息
func (m *OrderHandlerManager) GetHandlerInfo() map[string]interface{} {
	return map[string]interface{}{
		"handlerCount": m.factory.GetHandlerCount(),
		"handlerTypes": m.factory.ListHandlers(),
	}
}
