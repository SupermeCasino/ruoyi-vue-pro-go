package router

import (
	payAdmin "github.com/wxlbd/ruoyi-mall-go/internal/api/handler/admin/pay"
	payWallet "github.com/wxlbd/ruoyi-mall-go/internal/api/handler/admin/pay/wallet"
	"github.com/wxlbd/ruoyi-mall-go/internal/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterPayRoutes 注册支付模块路由
func RegisterPayRoutes(engine *gin.Engine,
	payAppHandler *payAdmin.PayAppHandler,
	payChannelHandler *payAdmin.PayChannelHandler,
	payOrderHandler *payAdmin.PayOrderHandler,
	payRefundHandler *payAdmin.PayRefundHandler,
	payNotifyHandler *payAdmin.PayNotifyHandler,
	payTransferHandler *payAdmin.PayTransferHandler,
	// Wallet
	payWalletHandler *payWallet.PayWalletHandler,
	payWalletRechargeHandler *payWallet.PayWalletRechargeHandler,
	payWalletRechargePackageHandler *payWallet.PayWalletRechargePackageHandler,
	payWalletTransactionHandler *payWallet.PayWalletTransactionHandler,
	casbinMiddleware *middleware.CasbinMiddleware,
) {
	api := engine.Group("/admin-api")
	payGroup := api.Group("/pay", middleware.Auth())
	{
		// Pay App
		payApp := payGroup.Group("/app")
		{
			payApp.POST("/create", casbinMiddleware.RequirePermission("pay:app:create"), payAppHandler.CreateApp)
			payApp.PUT("/update", casbinMiddleware.RequirePermission("pay:app:update"), payAppHandler.UpdateApp)
			payApp.PUT("/update-status", casbinMiddleware.RequirePermission("pay:app:update"), payAppHandler.UpdateAppStatus)
			payApp.DELETE("/delete", casbinMiddleware.RequirePermission("pay:app:delete"), payAppHandler.DeleteApp)
			payApp.GET("/get", casbinMiddleware.RequirePermission("pay:app:query"), payAppHandler.GetApp)
			payApp.GET("/page", casbinMiddleware.RequirePermission("pay:app:query"), payAppHandler.GetAppPage)
			payApp.GET("/list", casbinMiddleware.RequirePermission("pay:app:query"), payAppHandler.GetAppList)
		}

		// Pay Channel
		payChannel := payGroup.Group("/channel")
		{
			payChannel.POST("/create", casbinMiddleware.RequirePermission("pay:channel:create"), payChannelHandler.CreateChannel)
			payChannel.PUT("/update", casbinMiddleware.RequirePermission("pay:channel:update"), payChannelHandler.UpdateChannel)
			payChannel.DELETE("/delete", casbinMiddleware.RequirePermission("pay:channel:delete"), payChannelHandler.DeleteChannel)
			payChannel.GET("/get", casbinMiddleware.RequirePermission("pay:channel:query"), payChannelHandler.GetChannel)
			payChannel.GET("/get-enable-code-list", casbinMiddleware.RequirePermission("pay:channel:query"), payChannelHandler.GetEnableChannelCodeList)
		}

		// Pay Order
		payOrder := payGroup.Group("/order")
		{
			payOrder.GET("/get", casbinMiddleware.RequirePermission("pay:order:query"), payOrderHandler.GetOrder)
			payOrder.GET("/get-detail", casbinMiddleware.RequirePermission("pay:order:query"), payOrderHandler.GetOrderDetail)
			payOrder.GET("/page", casbinMiddleware.RequirePermission("pay:order:query"), payOrderHandler.GetOrderPage)
			payOrder.GET("/export-excel", casbinMiddleware.RequirePermission("pay:order:export"), payOrderHandler.ExportOrderExcel)
			payOrder.POST("/submit", payOrderHandler.SubmitPayOrder)
		}

		// Pay Refund
		payRefund := payGroup.Group("/refund")
		{
			payRefund.GET("/get", casbinMiddleware.RequirePermission("pay:refund:query"), payRefundHandler.GetRefund)
			payRefund.GET("/page", casbinMiddleware.RequirePermission("pay:refund:query"), payRefundHandler.GetRefundPage)
			payRefund.GET("/export-excel", casbinMiddleware.RequirePermission("pay:refund:export"), payRefundHandler.ExportRefundExcel)
		}

		// Pay Notify 管理接口 - 需要认证
		payNotify := payGroup.Group("/notify")
		{
			payNotify.GET("/get-detail", casbinMiddleware.RequirePermission("pay:notify:query"), payNotifyHandler.GetNotifyTaskDetail)
			payNotify.GET("/page", casbinMiddleware.RequirePermission("pay:notify:query"), payNotifyHandler.GetNotifyTaskPage)
		}

		// Pay Transfer
		payTransfer := payGroup.Group("/transfer")
		{
			payTransfer.GET("/get", casbinMiddleware.RequirePermission("pay:transfer:query"), payTransferHandler.GetTransfer)
			payTransfer.GET("/page", casbinMiddleware.RequirePermission("pay:transfer:query"), payTransferHandler.GetTransferPage)
		}

		// Pay Wallet
		payWallet := payGroup.Group("/wallet")
		{
			payWallet.GET("/get", casbinMiddleware.RequirePermission("pay:wallet:query"), payWalletHandler.GetWallet)
			payWallet.GET("/page", casbinMiddleware.RequirePermission("pay:wallet:query"), payWalletHandler.GetWalletPage)
			payWallet.PUT("/update-balance", casbinMiddleware.RequirePermission("pay:wallet:update-balance"), payWalletHandler.UpdateWalletBalance)
		}

		// Pay Wallet Recharge
		payWalletRecharge := payGroup.Group("/wallet-recharge")
		{
			payWalletRecharge.GET("/page", casbinMiddleware.RequirePermission("pay:wallet-recharge:query"), payWalletRechargeHandler.GetWalletRechargePage)
			payWalletRecharge.PUT("/update-paid", casbinMiddleware.RequirePermission("pay:wallet-recharge:update"), payWalletRechargeHandler.UpdateWalletRechargePaid)
			payWalletRecharge.POST("/refund", casbinMiddleware.RequirePermission("pay:wallet-recharge:refund"), payWalletRechargeHandler.RefundWalletRecharge)
			payWalletRecharge.PUT("/update-refunded", casbinMiddleware.RequirePermission("pay:wallet-recharge:update"), payWalletRechargeHandler.UpdateWalletRechargeRefunded)
		}

		// Pay Wallet Transaction
		payWalletTransaction := payGroup.Group("/wallet-transaction")
		{
			payWalletTransaction.GET("/page", payWalletTransactionHandler.GetWalletTransactionPage)
		}

		// Pay Wallet Recharge Package
		payWalletPackage := payGroup.Group("/wallet-recharge-package")
		{
			payWalletPackage.POST("/create", casbinMiddleware.RequirePermission("pay:wallet-recharge-package:create"), payWalletRechargePackageHandler.CreateWalletRechargePackage)
			payWalletPackage.PUT("/update", casbinMiddleware.RequirePermission("pay:wallet-recharge-package:update"), payWalletRechargePackageHandler.UpdateWalletRechargePackage)
			payWalletPackage.DELETE("/delete", casbinMiddleware.RequirePermission("pay:wallet-recharge-package:delete"), payWalletRechargePackageHandler.DeleteWalletRechargePackage)
			payWalletPackage.GET("/get", casbinMiddleware.RequirePermission("pay:wallet-recharge-package:query"), payWalletRechargePackageHandler.GetWalletRechargePackage)
			payWalletPackage.GET("/page", casbinMiddleware.RequirePermission("pay:wallet-recharge-package:query"), payWalletRechargePackageHandler.GetWalletRechargePackagePage)
		}
	}

	// Pay Notify 回调路由 (无需认证)
	apiPublic := engine.Group("/admin-api/pay/notify")
	{
		apiPublic.POST("/order/:channelId", payNotifyHandler.NotifyOrder)
		apiPublic.POST("/refund/:channelId", payNotifyHandler.NotifyRefund)
		apiPublic.POST("/transfer/:channelId", payNotifyHandler.NotifyTransfer)
	}
}
