package router

import (
	payAdmin "backend-go/internal/api/handler/admin/pay"
	payWallet "backend-go/internal/api/handler/admin/pay/wallet"

	"github.com/gin-gonic/gin"
)

// RegisterPayRoutes 注册支付模块路由
func RegisterPayRoutes(engine *gin.Engine,
	payAppHandler *payAdmin.PayAppHandler,
	payChannelHandler *payAdmin.PayChannelHandler,
	payOrderHandler *payAdmin.PayOrderHandler,
	payRefundHandler *payAdmin.PayRefundHandler,
	payNotifyHandler *payAdmin.PayNotifyHandler,
	// Wallet
	payWalletHandler *payWallet.PayWalletHandler,
	payWalletRechargeHandler *payWallet.PayWalletRechargeHandler,
	payWalletRechargePackageHandler *payWallet.PayWalletRechargePackageHandler,
	payWalletTransactionHandler *payWallet.PayWalletTransactionHandler,
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
			// 回调接口 - 无需认证 (对齐 Java @PermitAll @TenantIgnore)
			payNotify.POST("/order/:channelId", payNotifyHandler.NotifyOrder)
			payNotify.POST("/refund/:channelId", payNotifyHandler.NotifyRefund)
			payNotify.POST("/transfer/:channelId", payNotifyHandler.NotifyTransfer)

			// 管理接口 - 需要认证
			payNotify.GET("/get-detail", payNotifyHandler.GetNotifyTaskDetail)
			payNotify.GET("/page", payNotifyHandler.GetNotifyTaskPage)
		}

		// Pay Wallet
		payWallet := payGroup.Group("/wallet")
		{
			payWallet.GET("/get", payWalletHandler.GetWallet)
			payWallet.GET("/page", payWalletHandler.GetWalletPage)
		}

		// Pay Wallet Recharge
		payWalletRecharge := payGroup.Group("/wallet-recharge")
		{
			payWalletRecharge.GET("/page", payWalletRechargeHandler.GetWalletRechargePage)
		}

		// Pay Wallet Transaction
		payWalletTransaction := payGroup.Group("/wallet-transaction")
		{
			payWalletTransaction.GET("/page", payWalletTransactionHandler.GetWalletTransactionPage)
		}

		// Pay Wallet Recharge Package
		payWalletPackage := payGroup.Group("/wallet-recharge-package")
		{
			payWalletPackage.POST("/create", payWalletRechargePackageHandler.CreateWalletRechargePackage)
			payWalletPackage.PUT("/update", payWalletRechargePackageHandler.UpdateWalletRechargePackage)
			payWalletPackage.DELETE("/delete", payWalletRechargePackageHandler.DeleteWalletRechargePackage)
			payWalletPackage.GET("/get", payWalletRechargePackageHandler.GetWalletRechargePackage)
			payWalletPackage.GET("/page", payWalletRechargePackageHandler.GetWalletRechargePackagePage)
		}
	}
}
