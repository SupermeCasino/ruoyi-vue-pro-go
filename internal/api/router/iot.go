package router

import (
	"github.com/gin-gonic/gin"
	"github.com/wxlbd/ruoyi-mall-go/internal/api/handler/admin/iot"
	"github.com/wxlbd/ruoyi-mall-go/internal/middleware"
)

func RegisterIotRoutes(r *gin.Engine, h *iot.Handlers, casbin *middleware.CasbinMiddleware) {
	// 管理端 API
	adminGroup := r.Group("/admin-api/iot")
	adminGroup.Use(middleware.Auth())
	{
		// 物模型
		thingModel := adminGroup.Group("/thing-model")
		{
			thingModel.POST("/create", casbin.RequirePermission("iot:thing-model:create"), h.ThingModel.Create)
			thingModel.PUT("/update", casbin.RequirePermission("iot:thing-model:update"), h.ThingModel.Update)
			thingModel.DELETE("/delete", casbin.RequirePermission("iot:thing-model:delete"), h.ThingModel.Delete)
			thingModel.GET("/get", casbin.RequirePermission("iot:thing-model:query"), h.ThingModel.Get)
			thingModel.GET("/get-tsl", casbin.RequirePermission("iot:thing-model:query"), h.ThingModel.GetTSL)
			thingModel.GET("/list", casbin.RequirePermission("iot:thing-model:query"), h.ThingModel.List)
			thingModel.GET("/page", casbin.RequirePermission("iot:thing-model:query"), h.ThingModel.Page)
		}

		// 产品管理
		product := adminGroup.Group("/product")
		{
			product.POST("/create", casbin.RequirePermission("iot:product:create"), h.Product.Create)
			product.PUT("/update", casbin.RequirePermission("iot:product:update"), h.Product.Update)
			product.PUT("/update-status", casbin.RequirePermission("iot:product:update"), h.Product.UpdateStatus)
			product.DELETE("/delete", casbin.RequirePermission("iot:product:delete"), h.Product.Delete)
			product.GET("/get", casbin.RequirePermission("iot:product:query"), h.Product.Get)
			product.GET("/get-by-key", casbin.RequirePermission("iot:product:query"), h.Product.GetByKey)
			product.GET("/simple-list", casbin.RequirePermission("iot:product:query"), h.Product.SimpleList)
			product.GET("/page", casbin.RequirePermission("iot:product:query"), h.Product.Page)
		}

		// 设备管理
		device := adminGroup.Group("/device")
		{
			device.POST("/create", casbin.RequirePermission("iot:device:create"), h.Device.Create)
			device.PUT("/update", casbin.RequirePermission("iot:device:update"), h.Device.Update)
			device.PUT("/update-group", casbin.RequirePermission("iot:device:update"), h.Device.UpdateGroup)
			device.DELETE("/delete", casbin.RequirePermission("iot:device:delete"), h.Device.Delete)
			device.DELETE("/delete-list", casbin.RequirePermission("iot:device:delete"), h.Device.DeleteList)
			device.GET("/get", casbin.RequirePermission("iot:device:query"), h.Device.Get)
			device.GET("/count", casbin.RequirePermission("iot:device:query"), h.Device.GetCount)
			device.GET("/get-auth-info", casbin.RequirePermission("iot:device:query"), h.Device.GetAuthInfo)
			device.GET("/page", casbin.RequirePermission("iot:device:query"), h.Device.Page)
			device.GET("/list-by-product-key-and-names", casbin.RequirePermission("iot:device:query"), h.Device.GetListByProductKeyAndNames)
		}

		// 设备分组管理
		deviceGroup := adminGroup.Group("/device-group")
		{
			deviceGroup.POST("/create", casbin.RequirePermission("iot:device-group:create"), h.DeviceGroup.Create)
			deviceGroup.PUT("/update", casbin.RequirePermission("iot:device-group:update"), h.DeviceGroup.Update)
			deviceGroup.DELETE("/delete", casbin.RequirePermission("iot:device-group:delete"), h.DeviceGroup.Delete)
			deviceGroup.GET("/get", casbin.RequirePermission("iot:device-group:query"), h.DeviceGroup.Get)
			deviceGroup.GET("/page", casbin.RequirePermission("iot:device-group:query"), h.DeviceGroup.Page)
			deviceGroup.GET("/simple-list", h.DeviceGroup.SimpleList)
		}

		// OTA 固件管理
		otaFirmware := adminGroup.Group("/ota-firmware")
		{
			otaFirmware.POST("/create", casbin.RequirePermission("iot:ota-firmware:create"), h.OtaFirmware.Create)
			otaFirmware.PUT("/update", casbin.RequirePermission("iot:ota-firmware:update"), h.OtaFirmware.Update)
			otaFirmware.DELETE("/delete", casbin.RequirePermission("iot:ota-firmware:delete"), h.OtaFirmware.Delete)
			otaFirmware.GET("/get", casbin.RequirePermission("iot:ota-firmware:query"), h.OtaFirmware.Get)
			otaFirmware.GET("/page", casbin.RequirePermission("iot:ota-firmware:query"), h.OtaFirmware.Page)
		}

		// OTA 任务管理
		otaTask := adminGroup.Group("/ota-task")
		{
			otaTask.POST("/create", casbin.RequirePermission("iot:ota-task:create"), h.OtaTask.Create)
			otaTask.GET("/page", casbin.RequirePermission("iot:ota-task:query"), h.OtaTask.Page)
		}

		// OTA 记录管理
		otaTaskRecord := adminGroup.Group("/ota-task-record")
		{
			otaTaskRecord.GET("/page", casbin.RequirePermission("iot:ota-task:query"), h.OtaTask.RecordPage)
		}

		// 告警配置管理
		alertConfig := adminGroup.Group("/alert-config")
		{
			alertConfig.POST("/create", casbin.RequirePermission("iot:alert-config:create"), h.AlertConfig.Create)
			alertConfig.PUT("/update", casbin.RequirePermission("iot:alert-config:update"), h.AlertConfig.Update)
			alertConfig.DELETE("/delete", casbin.RequirePermission("iot:alert-config:delete"), h.AlertConfig.Delete)
			alertConfig.GET("/get", casbin.RequirePermission("iot:alert-config:query"), h.AlertConfig.Get)
			alertConfig.GET("/page", casbin.RequirePermission("iot:alert-config:query"), h.AlertConfig.Page)
			alertConfig.GET("/simple-list", casbin.RequirePermission("iot:alert-config:query"), h.AlertConfig.SimpleList)
		}

		// 告警记录管理
		alertRecord := adminGroup.Group("/alert-record")
		{
			alertRecord.GET("/get", casbin.RequirePermission("iot:alert-record:query"), h.AlertRecord.Get)
			alertRecord.GET("/page", casbin.RequirePermission("iot:alert-record:query"), h.AlertRecord.Page)
			alertRecord.PUT("/process", casbin.RequirePermission("iot:alert-record:update"), h.AlertRecord.Process)
		}

		// 数据目的管理
		dataSink := adminGroup.Group("/data-sink")
		{
			dataSink.POST("/create", casbin.RequirePermission("iot:data-sink:create"), h.DataSink.Create)
			dataSink.PUT("/update", casbin.RequirePermission("iot:data-sink:update"), h.DataSink.Update)
			dataSink.DELETE("/delete", casbin.RequirePermission("iot:data-sink:delete"), h.DataSink.Delete)
			dataSink.GET("/get", casbin.RequirePermission("iot:data-sink:query"), h.DataSink.Get)
			dataSink.GET("/page", casbin.RequirePermission("iot:data-sink:query"), h.DataSink.Page)
			dataSink.GET("/simple-list", casbin.RequirePermission("iot:data-sink:query"), h.DataSink.SimpleList)
		}

		// 数据规则管理
		dataRule := adminGroup.Group("/data-rule")
		{
			dataRule.POST("/create", casbin.RequirePermission("iot:data-rule:create"), h.DataRule.Create)
			dataRule.PUT("/update", casbin.RequirePermission("iot:data-rule:update"), h.DataRule.Update)
			dataRule.DELETE("/delete", casbin.RequirePermission("iot:data-rule:delete"), h.DataRule.Delete)
			dataRule.GET("/get", casbin.RequirePermission("iot:data-rule:query"), h.DataRule.Get)
			dataRule.GET("/page", casbin.RequirePermission("iot:data-rule:query"), h.DataRule.Page)
		}

		// 场景联动管理
		sceneRule := adminGroup.Group("/scene-rule")
		{
			sceneRule.POST("/create", casbin.RequirePermission("iot:scene-rule:create"), h.SceneRule.Create)
			sceneRule.PUT("/update", casbin.RequirePermission("iot:scene-rule:update"), h.SceneRule.Update)
			sceneRule.PUT("/update-status", casbin.RequirePermission("iot:scene-rule:update"), h.SceneRule.UpdateStatus)
			sceneRule.DELETE("/delete", casbin.RequirePermission("iot:scene-rule:delete"), h.SceneRule.Delete)
			sceneRule.GET("/get", casbin.RequirePermission("iot:scene-rule:query"), h.SceneRule.Get)
			sceneRule.GET("/page", casbin.RequirePermission("iot:scene-rule:query"), h.SceneRule.Page)
			sceneRule.GET("/simple-list", casbin.RequirePermission("iot:scene-rule:query"), h.SceneRule.SimpleList)
		}

		// 产品分类管理
		productCategory := adminGroup.Group("/product-category")
		{
			productCategory.POST("/create", casbin.RequirePermission("iot:product-category:create"), h.ProductCategory.Create)
			productCategory.PUT("/update", casbin.RequirePermission("iot:product-category:update"), h.ProductCategory.Update)
			productCategory.DELETE("/delete", casbin.RequirePermission("iot:product-category:delete"), h.ProductCategory.Delete)
			productCategory.GET("/get", casbin.RequirePermission("iot:product-category:query"), h.ProductCategory.Get)
			productCategory.GET("/page", casbin.RequirePermission("iot:product-category:query"), h.ProductCategory.Page)
			productCategory.GET("/simple-list", casbin.RequirePermission("iot:product-category:query"), h.ProductCategory.SimpleList)
		}

		// 数据统计
		statistics := adminGroup.Group("/statistics")
		{
			statistics.GET("/get-summary", casbin.RequirePermission("iot:statistics:query"), h.Statistics.GetSummary)
			statistics.GET("/get-device-message-summary-by-date", casbin.RequirePermission("iot:statistics:query"), h.Statistics.GetDeviceMessageSummaryByDate)
		}

		// 设备消息
		deviceMessage := adminGroup.Group("/device/message")
		{
			deviceMessage.GET("/page", casbin.RequirePermission("iot:device:message-query"), h.DeviceMessage.GetPage)
			deviceMessage.GET("/pair-page", casbin.RequirePermission("iot:device:message-query"), h.DeviceMessage.GetPairPage)
			deviceMessage.POST("/send", casbin.RequirePermission("iot:device:message-end"), h.DeviceMessage.Send)
		}

		// 设备属性
		deviceProperty := adminGroup.Group("/device/property")
		{
			deviceProperty.GET("/get-latest", casbin.RequirePermission("iot:device:property-query"), h.DeviceProperty.GetLatest)
			deviceProperty.GET("/history-list", casbin.RequirePermission("iot:device:property-query"), h.DeviceProperty.GetHistoryList)
		}
	}
}
