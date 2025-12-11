package router

import (
	payAdmin "backend-go/internal/api/handler/admin/pay"

	"github.com/gin-gonic/gin"
)

// RegisterPayRoutes 注册支付模块路由
func RegisterPayRoutes(engine *gin.Engine,
	payAppHandler *payAdmin.PayAppHandler,
	payChannelHandler *payAdmin.PayChannelHandler,
	payOrderHandler *payAdmin.PayOrderHandler,
	payRefundHandler *payAdmin.PayRefundHandler,
	payNotifyHandler *payAdmin.PayNotifyHandler,
) {
	api := engine.Group("/admin-api")
	payGroup := api.Group("/pay")
	{
		// Pay App
		payApp := payGroup.Group("/app")
		{
			payApp.POST("/create", payAppHandler.CreateApp)
			payApp.PUT("/update", payAppHandler.UpdateApp)
			payApp.PUT("/update-status", payAppHandler.UpdateAppStatus)
			payApp.DELETE("/delete", payAppHandler.DeleteApp)
			payApp.GET("/get", payAppHandler.GetApp)
			payApp.GET("/page", payAppHandler.GetAppPage)
			payApp.GET("/list", payAppHandler.GetAppList)
		}

		// Pay Channel
		payChannel := payGroup.Group("/channel")
		{
			payChannel.POST("/create", payChannelHandler.CreateChannel)
			payChannel.PUT("/update", payChannelHandler.UpdateChannel)
			payChannel.DELETE("/delete", payChannelHandler.DeleteChannel)
			payChannel.GET("/get", payChannelHandler.GetChannel)
			payChannel.GET("/get-enable-code-list", payChannelHandler.GetEnableChannelCodeList)
		}

		// Pay Order
		payOrder := payGroup.Group("/order")
		{
			payOrder.GET("/get", payOrderHandler.GetOrder)
			payOrder.GET("/get-detail", payOrderHandler.GetOrderDetail)
			payOrder.GET("/page", payOrderHandler.GetOrderPage)
			payOrder.POST("/submit", payOrderHandler.SubmitPayOrder)
		}

		// Pay Refund
		payRefund := payGroup.Group("/refund")
		{
			payRefund.GET("/get", payRefundHandler.GetRefund)
			payRefund.GET("/page", payRefundHandler.GetRefundPage)
		}

		// Pay Notify
		payNotify := payGroup.Group("/notify")
		{
			payNotify.GET("/get-detail", payNotifyHandler.GetNotifyTaskDetail)
			payNotify.GET("/page", payNotifyHandler.GetNotifyTaskPage)
		}
	}
}
