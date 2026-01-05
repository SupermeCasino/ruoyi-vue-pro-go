package iot

import (
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	NewProductRepository,
	NewDeviceRepository,
	NewThingModelRepository,
	NewDeviceGroupRepository,
	NewOtaFirmwareRepository,
	NewOtaTaskRepository,
	NewOtaTaskRecordRepository,
	NewAlertConfigRepository,
	NewAlertRecordRepository,
	NewDataRuleRepository,
	NewDataSinkRepository,
	NewSceneRuleRepository,
	NewProductCategoryRepository,
	NewDeviceMessageRepository,
	NewDevicePropertyRepository,
)
