package router

import (
	tradeAdmin "github.com/wxlbd/ruoyi-mall-go/internal/api/handler/admin/trade"
	tradeBrokerage "github.com/wxlbd/ruoyi-mall-go/internal/api/handler/admin/trade/brokerage"
	"github.com/wxlbd/ruoyi-mall-go/internal/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterTradeRoutes 注册交易订单模块路由
func RegisterTradeRoutes(engine *gin.Engine,
	tradeOrderHandler *tradeAdmin.TradeOrderHandler,
	tradeAfterSaleHandler *tradeAdmin.TradeAfterSaleHandler,
	deliveryExpressHandler *tradeAdmin.DeliveryExpressHandler,
	deliveryPickUpStoreHandler *tradeAdmin.DeliveryPickUpStoreHandler,
	deliveryExpressTemplateHandler *tradeAdmin.DeliveryExpressTemplateHandler,
	tradeConfigHandler *tradeAdmin.TradeConfigHandler,
	brokerageUserHandler *tradeBrokerage.BrokerageUserHandler,
	brokerageRecordHandler *tradeBrokerage.BrokerageRecordHandler,
	brokerageWithdrawHandler *tradeBrokerage.BrokerageWithdrawHandler,
	casbinMiddleware *middleware.CasbinMiddleware,
) {
	// Trade Order
	tradeGroup := engine.Group("/admin-api/trade/order")
	tradeGroup.Use(middleware.Auth())
	{
		tradeGroup.GET("/page", tradeOrderHandler.GetOrderPage)
		tradeGroup.GET("/get-detail", tradeOrderHandler.GetOrderDetail)
		tradeGroup.GET("/get-summary", tradeOrderHandler.GetOrderSummary)
		tradeGroup.GET("/get-express-track-list", tradeOrderHandler.GetOrderExpressTrackList)
		tradeGroup.GET("/get-by-pick-up-verify-code", tradeOrderHandler.GetByPickUpVerifyCode)
		tradeGroup.PUT("/delivery", tradeOrderHandler.DeliveryOrder)
		tradeGroup.PUT("/update-remark", tradeOrderHandler.UpdateOrderRemark)
		tradeGroup.PUT("/update-price", tradeOrderHandler.UpdateOrderPrice)
		tradeGroup.PUT("/update-address", tradeOrderHandler.UpdateOrderAddress)
		tradeGroup.PUT("/pick-up-by-id", tradeOrderHandler.PickUpOrderById)
		tradeGroup.PUT("/pick-up-by-verify-code", tradeOrderHandler.PickUpOrderByVerifyCode)
	}

	// Trade AfterSale
	afterSaleGroup := engine.Group("/admin-api/trade/after-sale")
	afterSaleGroup.Use(middleware.Auth())
	{
		afterSaleGroup.GET("/page", tradeAfterSaleHandler.GetAfterSalePage)
		afterSaleGroup.GET("/get-detail", tradeAfterSaleHandler.GetAfterSaleDetail)
		afterSaleGroup.PUT("/agree", tradeAfterSaleHandler.AgreeAfterSale)
		afterSaleGroup.PUT("/disagree", tradeAfterSaleHandler.DisagreeAfterSale)
		afterSaleGroup.PUT("/receive", tradeAfterSaleHandler.ReceiveAfterSale)
		afterSaleGroup.PUT("/refund", tradeAfterSaleHandler.RefundAfterSale)
	}

	// Delivery Routes
	deliveryGroup := engine.Group("/admin-api/trade/delivery")
	deliveryGroup.Use(middleware.Auth())
	{
		// Express
		expressGroup := deliveryGroup.Group("/express")
		{
			expressGroup.POST("/create", deliveryExpressHandler.CreateDeliveryExpress)
			expressGroup.PUT("/update", deliveryExpressHandler.UpdateDeliveryExpress)
			expressGroup.DELETE("/delete", deliveryExpressHandler.DeleteDeliveryExpress)
			expressGroup.GET("/get", deliveryExpressHandler.GetDeliveryExpress)
			expressGroup.GET("/page", deliveryExpressHandler.GetDeliveryExpressPage)
			expressGroup.GET("/list-all-simple", deliveryExpressHandler.GetSimpleDeliveryExpressList)
			expressGroup.GET("/export-excel", deliveryExpressHandler.ExportDeliveryExpress)
		}

		// Pick Up Store
		pickUpStoreGroup := deliveryGroup.Group("/pick-up-store")
		{
			pickUpStoreGroup.POST("/create", deliveryPickUpStoreHandler.CreateDeliveryPickUpStore)
			pickUpStoreGroup.PUT("/update", deliveryPickUpStoreHandler.UpdateDeliveryPickUpStore)
			pickUpStoreGroup.DELETE("/delete", deliveryPickUpStoreHandler.DeleteDeliveryPickUpStore)
			pickUpStoreGroup.GET("/get", deliveryPickUpStoreHandler.GetDeliveryPickUpStore)
			pickUpStoreGroup.GET("/page", deliveryPickUpStoreHandler.GetDeliveryPickUpStorePage)
			pickUpStoreGroup.GET("/simple-list", deliveryPickUpStoreHandler.GetSimpleDeliveryPickUpStoreList)
		}

		// Express Template (运费模板) - 对齐 Java 路径
		expressTemplateGroup := deliveryGroup.Group("/express-template")
		{
			expressTemplateGroup.POST("/create", deliveryExpressTemplateHandler.CreateDeliveryExpressTemplate)
			expressTemplateGroup.PUT("/update", deliveryExpressTemplateHandler.UpdateDeliveryExpressTemplate)
			expressTemplateGroup.DELETE("/delete", deliveryExpressTemplateHandler.DeleteDeliveryExpressTemplate)
			expressTemplateGroup.GET("/get", deliveryExpressTemplateHandler.GetDeliveryExpressTemplate)
			expressTemplateGroup.GET("/page", deliveryExpressTemplateHandler.GetDeliveryExpressTemplatePage)
			expressTemplateGroup.GET("/list-all-simple", deliveryExpressTemplateHandler.GetSimpleDeliveryExpressTemplateList)
		}
	}

	// Trade Config (Admin)
	tradeConfigGroup := engine.Group("/admin-api/trade/config")
	tradeConfigGroup.Use(middleware.Auth())
	{
		tradeConfigGroup.GET("/get", tradeConfigHandler.GetTradeConfig)
		tradeConfigGroup.PUT("/save", tradeConfigHandler.SaveTradeConfig)
	}

	// Brokerage User
	brokerageUserGroup := engine.Group("/admin-api/trade/brokerage-user")
	brokerageUserGroup.Use(middleware.Auth())
	{
		brokerageUserGroup.POST("/create", brokerageUserHandler.CreateBrokerageUser)
		brokerageUserGroup.PUT("/update-bind-user", brokerageUserHandler.UpdateBindUser)
		brokerageUserGroup.PUT("/clear-bind-user", brokerageUserHandler.ClearBindUser)
		brokerageUserGroup.PUT("/update-brokerage-enable", brokerageUserHandler.UpdateBrokerageEnabled)
		brokerageUserGroup.GET("/get", brokerageUserHandler.GetBrokerageUser)
		brokerageUserGroup.GET("/page", brokerageUserHandler.GetBrokerageUserPage)
	}

	// Brokerage Record
	brokerageRecordGroup := engine.Group("/admin-api/trade/brokerage-record")
	brokerageRecordGroup.Use(middleware.Auth())
	{
		brokerageRecordGroup.GET("/get", brokerageRecordHandler.GetBrokerageRecord)
		brokerageRecordGroup.GET("/page", brokerageRecordHandler.GetBrokerageRecordPage)
	}

	// Brokerage Withdraw
	brokerageWithdrawGroup := engine.Group("/admin-api/trade/brokerage-withdraw")
	brokerageWithdrawGroup.Use(middleware.Auth())
	{
		brokerageWithdrawGroup.PUT("/approve", brokerageWithdrawHandler.ApproveBrokerageWithdraw)
		brokerageWithdrawGroup.PUT("/reject", brokerageWithdrawHandler.RejectBrokerageWithdraw)
		brokerageWithdrawGroup.POST("/update-transferred", brokerageWithdrawHandler.UpdateBrokerageWithdrawTransferred)
		brokerageWithdrawGroup.GET("/get", brokerageWithdrawHandler.GetBrokerageWithdraw)
		brokerageWithdrawGroup.GET("/page", brokerageWithdrawHandler.GetBrokerageWithdrawPage)
	}

	// Trade AfterSale Callback (No Auth)
	afterSaleCallbackGroup := engine.Group("/admin-api/trade/after-sale")
	{
		afterSaleCallbackGroup.POST("/update-refunded", tradeAfterSaleHandler.UpdateAfterSaleRefunded)
	}
}
