package router

import (
	"backend-go/internal/api/handler/admin"
	"backend-go/internal/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterStatisticsRoutes 注册统计模块路由
func RegisterStatisticsRoutes(engine *gin.Engine,
	tradeStatisticsHandler *admin.TradeStatisticsHandler,
	productStatisticsHandler *admin.ProductStatisticsHandler,
	memberStatisticsHandler *admin.MemberStatisticsHandler,
	payStatisticsHandler *admin.PayStatisticsHandler,
	casbinMiddleware *middleware.CasbinMiddleware,
) {
	// 后台统计路由组 - 与 Java 保持一致使用 /admin-api/statistics
	adminGroup := engine.Group("/admin-api/statistics")
	adminGroup.Use(middleware.Auth())

	// 交易统计路由
	tradeGroup := adminGroup.Group("/trade")
	{
		tradeGroup.GET("/summary", tradeStatisticsHandler.GetTradeSummaryComparison)
		tradeGroup.GET("/analyse", tradeStatisticsHandler.GetTradeStatisticsAnalyse)
		tradeGroup.GET("/list", tradeStatisticsHandler.GetTradeStatisticsList)
		tradeGroup.GET("/order-count", tradeStatisticsHandler.GetOrderCount)
		tradeGroup.GET("/order-comparison", tradeStatisticsHandler.GetOrderComparison)
		tradeGroup.GET("/order-count-trend", tradeStatisticsHandler.GetOrderCountTrendComparison)
		tradeGroup.GET("/export-excel", tradeStatisticsHandler.ExportTradeStatisticsExcel)
	}

	// 商品统计路由
	productGroup := adminGroup.Group("/product")
	{
		productGroup.GET("/analyse", productStatisticsHandler.GetProductStatisticsAnalyse)
		productGroup.GET("/list", productStatisticsHandler.GetProductStatisticsList)
		productGroup.GET("/rank-page", productStatisticsHandler.GetProductStatisticsRankPage)
		productGroup.GET("/export-excel", productStatisticsHandler.ExportProductStatisticsExcel)
	}

	// 会员统计路由
	memberGroup := adminGroup.Group("/member")
	{
		memberGroup.GET("/summary", memberStatisticsHandler.GetMemberSummary)
		memberGroup.GET("/analyse", memberStatisticsHandler.GetMemberAnalyse)
		memberGroup.GET("/area-statistics-list", memberStatisticsHandler.GetMemberAreaStatisticsList)
		memberGroup.GET("/sex-statistics-list", memberStatisticsHandler.GetMemberSexStatisticsList)
		memberGroup.GET("/terminal-statistics-list", memberStatisticsHandler.GetMemberTerminalStatisticsList)
		memberGroup.GET("/user-count-comparison", memberStatisticsHandler.GetUserCountComparison)
		memberGroup.GET("/register-count-list", memberStatisticsHandler.GetMemberRegisterCountList)
	}

	// 支付统计路由
	payGroup := adminGroup.Group("/pay")
	{
		payGroup.GET("/summary", payStatisticsHandler.GetWalletRechargePrice)
	}
}
