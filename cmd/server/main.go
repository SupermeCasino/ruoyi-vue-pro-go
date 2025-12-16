package main

import (
	"github.com/wxlbd/ruoyi-mall-go/internal/api/handler"
	"github.com/wxlbd/ruoyi-mall-go/internal/api/router"
	"github.com/wxlbd/ruoyi-mall-go/internal/pkg/area"
	"github.com/wxlbd/ruoyi-mall-go/pkg/config"
	"github.com/wxlbd/ruoyi-mall-go/pkg/logger"

	"go.uber.org/zap"
)

func main() {
	// 1. 初始化配置
	if err := config.Load(); err != nil {
		panic(err)
	}

	// 2. 初始化日志
	logger.Init()
	logger.Info("Config and Logger initialized")

	// 3. 初始化地区数据
	if err := area.Init("configs/area.csv"); err != nil {
		logger.Log.Warn("Failed to init area data", zap.Error(err))
	}

	// 4. 初始化应用 (通过 Wire 注入)
	// 注意：InitDB 和 InitRedis 会在 InitApp 中被自动调用
	r, err := InitApp()
	if err != nil {
		logger.Log.Fatal("failed to init app", zap.Error(err))
	}

	// 5. 注册地区路由 (独立于 Wire)
	areaHandler := handler.NewAreaHandler()
	router.RegisterAreaRoutes(r, areaHandler)

	// 6. 启动服务
	addr := config.C.HTTP.Port
	logger.Info("Server starting...", zap.String("addr", addr))
	if err := r.Run(addr); err != nil {
		logger.Log.Fatal("failed to start server", zap.Error(err))
	}
}
