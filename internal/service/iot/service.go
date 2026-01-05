package iot

import (
	"github.com/google/wire"
	iotcore "github.com/wxlbd/ruoyi-mall-go/internal/iot/core"
)

// ProviderSet 提供所有 IOT 服务的依赖注入
var ProviderSet = wire.NewSet(
	iotcore.ProviderSet,
	NewProductService,
	NewDeviceService,
	NewThingModelService,
	NewDeviceGroupService,
	NewOtaFirmwareService,
	NewOtaTaskService,
	NewAlertConfigService,
	NewAlertRecordService,
	NewDataSinkService,
	NewDataRuleService,
	NewSceneRuleService,
	NewProductCategoryService,
	NewStatisticsService,
	NewDeviceMessageService,
	NewDevicePropertyService,
	NewIotDeviceCommonApiImpl,
)
