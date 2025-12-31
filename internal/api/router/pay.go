package router

import (
	pay2 "github.com/wxlbd/ruoyi-mall-go/internal/api/handler/admin/pay"
	"github.com/wxlbd/ruoyi-mall-go/internal/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterPayRoutes 注册支付模块路由
func RegisterPayRoutes(engine *gin.Engine,
	handlers *pay2.Handlers,
	casbinMiddleware *middleware.CasbinMiddleware,
) {
	api := engine.Group("/admin-api")
	payGroup := api.Group("/pay", middleware.Auth())
	{
		// Pay App
		payApp := payGroup.Group("/app")
		{
			payApp.POST("/create", casbinMiddleware.RequirePermission("pay:app:create"), handlers.App.CreateApp)
			payApp.PUT("/update", casbinMiddleware.RequirePermission("pay:app:update"), handlers.App.UpdateApp)
			payApp.PUT("/update-status", casbinMiddleware.RequirePermission("pay:app:update"), handlers.App.UpdateAppStatus)
			payApp.DELETE("/delete", casbinMiddleware.RequirePermission("pay:app:delete"), handlers.App.DeleteApp)
			payApp.GET("/get", casbinMiddleware.RequirePermission("pay:app:query"), handlers.App.GetApp)
			payApp.GET("/page", casbinMiddleware.RequirePermission("pay:app:query"), handlers.App.GetAppPage)
			payApp.GET("/list", casbinMiddleware.RequirePermission("pay:app:query"), handlers.App.GetAppList)
		}

		// Pay Channel
		payChannel := payGroup.Group("/channel")
		{
			payChannel.POST("/create", casbinMiddleware.RequirePermission("pay:channel:create"), handlers.Channel.CreateChannel)
			payChannel.PUT("/update", casbinMiddleware.RequirePermission("pay:channel:update"), handlers.Channel.UpdateChannel)
			payChannel.DELETE("/delete", casbinMiddleware.RequirePermission("pay:channel:delete"), handlers.Channel.DeleteChannel)
			payChannel.GET("/get", casbinMiddleware.RequirePermission("pay:channel:query"), handlers.Channel.GetChannel)
			payChannel.GET("/get-enable-code-list", casbinMiddleware.RequirePermission("pay:channel:query"), handlers.Channel.GetEnableChannelCodeList)
		}

		// Pay Order
		payOrder := payGroup.Group("/order")
		{
			payOrder.GET("/get", casbinMiddleware.RequirePermission("pay:order:query"), handlers.Order.GetOrder)
			payOrder.GET("/get-detail", casbinMiddleware.RequirePermission("pay:order:query"), handlers.Order.GetOrderDetail)
			payOrder.GET("/page", casbinMiddleware.RequirePermission("pay:order:query"), handlers.Order.GetOrderPage)
			payOrder.GET("/export-excel", casbinMiddleware.RequirePermission("pay:order:export"), handlers.Order.ExportOrderExcel)
			payOrder.POST("/submit", handlers.Order.SubmitPayOrder)
		}

		// Pay Refund
		payRefund := payGroup.Group("/refund")
		{
			payRefund.GET("/get", casbinMiddleware.RequirePermission("pay:refund:query"), handlers.Refund.GetRefund)
			payRefund.GET("/page", casbinMiddleware.RequirePermission("pay:refund:query"), handlers.Refund.GetRefundPage)
			payRefund.GET("/export-excel", casbinMiddleware.RequirePermission("pay:refund:export"), handlers.Refund.ExportRefundExcel)
		}

		// Pay Notify 管理接口 - 需要认证
		payNotify := payGroup.Group("/notify")
		{
			payNotify.GET("/get-detail", casbinMiddleware.RequirePermission("pay:notify:query"), handlers.Notify.GetNotifyTaskDetail)
			payNotify.GET("/page", casbinMiddleware.RequirePermission("pay:notify:query"), handlers.Notify.GetNotifyTaskPage)
		}

		// Pay Transfer
		payTransfer := payGroup.Group("/transfer")
		{
			payTransfer.GET("/get", casbinMiddleware.RequirePermission("pay:transfer:query"), handlers.Transfer.GetTransfer)
			payTransfer.GET("/page", casbinMiddleware.RequirePermission("pay:transfer:query"), handlers.Transfer.GetTransferPage)
		}

		// Pay Wallet
		payWallet := payGroup.Group("/wallet")
		{
			payWallet.GET("/get", casbinMiddleware.RequirePermission("pay:wallet:query"), handlers.Wallet.Wallet.GetWallet)
			payWallet.GET("/page", casbinMiddleware.RequirePermission("pay:wallet:query"), handlers.Wallet.Wallet.GetWalletPage)
			payWallet.PUT("/update-balance", casbinMiddleware.RequirePermission("pay:wallet:update-balance"), handlers.Wallet.Wallet.UpdateWalletBalance)
		}

		// Pay Wallet Recharge
		payWalletRecharge := payGroup.Group("/wallet-recharge")
		{
			payWalletRecharge.GET("/page", casbinMiddleware.RequirePermission("pay:wallet-recharge:query"), handlers.Wallet.Recharge.GetWalletRechargePage)
			payWalletRecharge.PUT("/update-paid", casbinMiddleware.RequirePermission("pay:wallet-recharge:update"), handlers.Wallet.Recharge.UpdateWalletRechargePaid)
			payWalletRecharge.POST("/refund", casbinMiddleware.RequirePermission("pay:wallet-recharge:refund"), handlers.Wallet.Recharge.RefundWalletRecharge)
			payWalletRecharge.PUT("/update-refunded", casbinMiddleware.RequirePermission("pay:wallet-recharge:update"), handlers.Wallet.Recharge.UpdateWalletRechargeRefunded)
		}

		// Pay Wallet Transaction
		payWalletTransaction := payGroup.Group("/wallet-transaction")
		{
			payWalletTransaction.GET("/page", handlers.Wallet.Transaction.GetWalletTransactionPage)
		}

		// Pay Wallet Recharge Package
		payWalletPackage := payGroup.Group("/wallet-recharge-package")
		{
			payWalletPackage.POST("/create", casbinMiddleware.RequirePermission("pay:wallet-recharge-package:create"), handlers.Wallet.RechargePackage.CreateWalletRechargePackage)
			payWalletPackage.PUT("/update", casbinMiddleware.RequirePermission("pay:wallet-recharge-package:update"), handlers.Wallet.RechargePackage.UpdateWalletRechargePackage)
			payWalletPackage.DELETE("/delete", casbinMiddleware.RequirePermission("pay:wallet-recharge-package:delete"), handlers.Wallet.RechargePackage.DeleteWalletRechargePackage)
			payWalletPackage.GET("/get", casbinMiddleware.RequirePermission("pay:wallet-recharge-package:query"), handlers.Wallet.RechargePackage.GetWalletRechargePackage)
			payWalletPackage.GET("/page", casbinMiddleware.RequirePermission("pay:wallet-recharge-package:query"), handlers.Wallet.RechargePackage.GetWalletRechargePackagePage)
		}
	}

	// Pay Notify 回调路由 (无需认证)
	apiPublic := engine.Group("/admin-api/pay/notify")
	{
		apiPublic.POST("/order/:channelId", handlers.Notify.NotifyOrder)
		apiPublic.POST("/refund/:channelId", handlers.Notify.NotifyRefund)
		apiPublic.POST("/transfer/:channelId", handlers.Notify.NotifyTransfer)
	}
}
