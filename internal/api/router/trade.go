package router

import (
	"github.com/wxlbd/ruoyi-mall-go/internal/api/handler/admin/mall/trade"
	"github.com/wxlbd/ruoyi-mall-go/internal/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterTradeRoutes 注册交易订单模块路由
func RegisterTradeRoutes(engine *gin.Engine,
	handlers *trade.Handlers,
	casbinMiddleware *middleware.CasbinMiddleware,
) {
	// Trade Order
	tradeGroup := engine.Group("/admin-api/trade/order")
	tradeGroup.Use(middleware.Auth())
	{
		tradeGroup.GET("/page", handlers.Order.GetOrderPage)
		tradeGroup.GET("/get-detail", handlers.Order.GetOrderDetail)
		tradeGroup.GET("/summary", handlers.Order.GetOrderSummary)
		tradeGroup.GET("/get-express-track-list", handlers.Order.GetOrderExpressTrackList)
		tradeGroup.GET("/get-by-pick-up-verify-code", handlers.Order.GetByPickUpVerifyCode)
		tradeGroup.PUT("/delivery", handlers.Order.DeliveryOrder)
		tradeGroup.PUT("/update-remark", handlers.Order.UpdateOrderRemark)
		tradeGroup.PUT("/update-price", handlers.Order.UpdateOrderPrice)
		tradeGroup.PUT("/update-address", handlers.Order.UpdateOrderAddress)
		tradeGroup.PUT("/pick-up-by-id", handlers.Order.PickUpOrderById)
		tradeGroup.PUT("/pick-up-by-verify-code", handlers.Order.PickUpOrderByVerifyCode)
	}

	// Trade AfterSale
	afterSaleGroup := engine.Group("/admin-api/trade/after-sale")
	afterSaleGroup.Use(middleware.Auth())
	{
		afterSaleGroup.GET("/page", handlers.AfterSale.GetAfterSalePage)
		afterSaleGroup.GET("/get-detail", handlers.AfterSale.GetAfterSaleDetail)
		afterSaleGroup.PUT("/agree", handlers.AfterSale.AgreeAfterSale)
		afterSaleGroup.PUT("/disagree", handlers.AfterSale.DisagreeAfterSale)
		afterSaleGroup.PUT("/receive", handlers.AfterSale.ReceiveAfterSale)
		afterSaleGroup.PUT("/refuse", handlers.AfterSale.RefuseAfterSale)
		afterSaleGroup.PUT("/refund", handlers.AfterSale.RefundAfterSale)
	}

	// Delivery Routes
	deliveryGroup := engine.Group("/admin-api/trade/delivery")
	deliveryGroup.Use(middleware.Auth())
	{
		// Express
		expressGroup := deliveryGroup.Group("/express")
		{
			expressGroup.POST("/create", handlers.DeliveryExpress.CreateDeliveryExpress)
			expressGroup.PUT("/update", handlers.DeliveryExpress.UpdateDeliveryExpress)
			expressGroup.DELETE("/delete", handlers.DeliveryExpress.DeleteDeliveryExpress)
			expressGroup.GET("/get", handlers.DeliveryExpress.GetDeliveryExpress)
			expressGroup.GET("/page", handlers.DeliveryExpress.GetDeliveryExpressPage)
			expressGroup.GET("/list-all-simple", handlers.DeliveryExpress.GetSimpleDeliveryExpressList)
			expressGroup.GET("/export-excel", handlers.DeliveryExpress.ExportDeliveryExpress)
		}

		// Pick Up Store
		pickUpStoreGroup := deliveryGroup.Group("/pick-up-store")
		{
			pickUpStoreGroup.POST("/create", handlers.DeliveryPickUpStore.CreateDeliveryPickUpStore)
			pickUpStoreGroup.PUT("/update", handlers.DeliveryPickUpStore.UpdateDeliveryPickUpStore)
			pickUpStoreGroup.DELETE("/delete", handlers.DeliveryPickUpStore.DeleteDeliveryPickUpStore)
			pickUpStoreGroup.GET("/get", handlers.DeliveryPickUpStore.GetDeliveryPickUpStore)
			pickUpStoreGroup.GET("/page", handlers.DeliveryPickUpStore.GetDeliveryPickUpStorePage)
			pickUpStoreGroup.GET("/simple-list", handlers.DeliveryPickUpStore.GetSimpleDeliveryPickUpStoreList)
			pickUpStoreGroup.POST("/bind", handlers.DeliveryPickUpStore.BindDeliveryPickUpStore)
		}

		// Express Template (运费模板) - 对齐 Java 路径
		expressTemplateGroup := deliveryGroup.Group("/express-template")
		{
			expressTemplateGroup.POST("/create", handlers.DeliveryExpressTemplate.CreateDeliveryExpressTemplate)
			expressTemplateGroup.PUT("/update", handlers.DeliveryExpressTemplate.UpdateDeliveryExpressTemplate)
			expressTemplateGroup.DELETE("/delete", handlers.DeliveryExpressTemplate.DeleteDeliveryExpressTemplate)
			expressTemplateGroup.GET("/get", handlers.DeliveryExpressTemplate.GetDeliveryExpressTemplate)
			expressTemplateGroup.GET("/page", handlers.DeliveryExpressTemplate.GetDeliveryExpressTemplatePage)
			expressTemplateGroup.GET("/list-all-simple", handlers.DeliveryExpressTemplate.GetSimpleDeliveryExpressTemplateList)
		}
	}

	// Trade Config (Admin)
	tradeConfigGroup := engine.Group("/admin-api/trade/config")
	tradeConfigGroup.Use(middleware.Auth())
	{
		tradeConfigGroup.GET("/get", handlers.Config.GetTradeConfig)
		tradeConfigGroup.PUT("/save", handlers.Config.SaveTradeConfig)
	}

	// Brokerage User
	brokerageUserGroup := engine.Group("/admin-api/trade/brokerage-user")
	brokerageUserGroup.Use(middleware.Auth())
	{
		brokerageUserGroup.POST("/create", handlers.Brokerage.BrokerageUser.CreateBrokerageUser)
		brokerageUserGroup.PUT("/update-bind-user", handlers.Brokerage.BrokerageUser.UpdateBindUser)
		brokerageUserGroup.PUT("/clear-bind-user", handlers.Brokerage.BrokerageUser.ClearBindUser)
		brokerageUserGroup.PUT("/update-brokerage-enable", handlers.Brokerage.BrokerageUser.UpdateBrokerageEnabled)
		brokerageUserGroup.GET("/get", handlers.Brokerage.BrokerageUser.GetBrokerageUser)
		brokerageUserGroup.GET("/page", handlers.Brokerage.BrokerageUser.GetBrokerageUserPage)
	}

	// Brokerage Record
	brokerageRecordGroup := engine.Group("/admin-api/trade/brokerage-record")
	brokerageRecordGroup.Use(middleware.Auth())
	{
		brokerageRecordGroup.GET("/get", handlers.Brokerage.BrokerageRecord.GetBrokerageRecord)
		brokerageRecordGroup.GET("/page", handlers.Brokerage.BrokerageRecord.GetBrokerageRecordPage)
	}

	// Brokerage Withdraw
	brokerageWithdrawGroup := engine.Group("/admin-api/trade/brokerage-withdraw")
	brokerageWithdrawGroup.Use(middleware.Auth())
	{
		brokerageWithdrawGroup.PUT("/approve", handlers.Brokerage.BrokerageWithdraw.ApproveBrokerageWithdraw)
		brokerageWithdrawGroup.PUT("/reject", handlers.Brokerage.BrokerageWithdraw.RejectBrokerageWithdraw)
		brokerageWithdrawGroup.POST("/update-transferred", handlers.Brokerage.BrokerageWithdraw.UpdateBrokerageWithdrawTransferred)
		brokerageWithdrawGroup.GET("/get", handlers.Brokerage.BrokerageWithdraw.GetBrokerageWithdraw)
		brokerageWithdrawGroup.GET("/page", handlers.Brokerage.BrokerageWithdraw.GetBrokerageWithdrawPage)
	}

	// Trade AfterSale Callback (No Auth)
	afterSaleCallbackGroup := engine.Group("/admin-api/trade/after-sale")
	{
		afterSaleCallbackGroup.POST("/update-refunded", handlers.AfterSale.UpdateAfterSaleRefunded)
	}
}
