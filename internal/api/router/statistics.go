package router

import (
	"github.com/wxlbd/ruoyi-mall-go/internal/api/handler/admin/statistics"
	"github.com/wxlbd/ruoyi-mall-go/internal/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterStatisticsRoutes 注册统计模块路由
func RegisterStatisticsRoutes(engine *gin.Engine,
	handlers *statistics.Handlers,
	casbinMiddleware *middleware.CasbinMiddleware,
) {
	// 后台统计路由组 - 与 Java 保持一致使用 /admin-api/statistics
	adminGroup := engine.Group("/admin-api/statistics")
	adminGroup.Use(middleware.Auth())

	// 交易统计路由
	tradeGroup := adminGroup.Group("/trade")
	{
		tradeGroup.GET("/summary", handlers.Trade.GetTradeSummaryComparison)
		tradeGroup.GET("/analyse", handlers.Trade.GetTradeStatisticsAnalyse)
		tradeGroup.GET("/list", handlers.Trade.GetTradeStatisticsList)
		tradeGroup.GET("/order-count", handlers.Trade.GetOrderCount)
		tradeGroup.GET("/order-comparison", handlers.Trade.GetOrderComparison)
		tradeGroup.GET("/order-count-trend", handlers.Trade.GetOrderCountTrendComparison)
		tradeGroup.GET("/export-excel", handlers.Trade.ExportTradeStatisticsExcel)
	}

	// 商品统计路由
	productGroup := adminGroup.Group("/product")
	{
		productGroup.GET("/analyse", handlers.Product.GetProductStatisticsAnalyse)
		productGroup.GET("/list", handlers.Product.GetProductStatisticsList)
		productGroup.GET("/rank-page", handlers.Product.GetProductStatisticsRankPage)
		productGroup.GET("/export-excel", handlers.Product.ExportProductStatisticsExcel)
	}

	// 会员统计路由
	memberGroup := adminGroup.Group("/member")
	{
		memberGroup.GET("/summary", handlers.Member.GetMemberSummary)
		memberGroup.GET("/analyse", handlers.Member.GetMemberAnalyse)
		memberGroup.GET("/area-statistics-list", handlers.Member.GetMemberAreaStatisticsList)
		memberGroup.GET("/sex-statistics-list", handlers.Member.GetMemberSexStatisticsList)
		memberGroup.GET("/terminal-statistics-list", handlers.Member.GetMemberTerminalStatisticsList)
		memberGroup.GET("/user-count-comparison", handlers.Member.GetUserCountComparison)
		memberGroup.GET("/register-count-list", handlers.Member.GetMemberRegisterCountList)
	}

	// 支付统计路由
	payGroup := adminGroup.Group("/pay")
	{
		payGroup.GET("/summary", handlers.Pay.GetWalletRechargePrice)
	}
}
