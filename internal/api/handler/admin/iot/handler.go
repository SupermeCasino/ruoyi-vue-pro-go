package iot

import (
	"github.com/google/wire"
	iotsvc "github.com/wxlbd/ruoyi-mall-go/internal/service/iot"
)

var ProviderSet = wire.NewSet(
	NewProductHandler,
	NewDeviceHandler,
	NewThingModelHandler,
	NewDeviceGroupHandler,
	NewOtaFirmwareHandler,
	NewOtaTaskHandler,
	NewAlertConfigHandler,
	NewAlertRecordHandler,
	NewDataSinkHandler,
	NewDataRuleHandler,
	NewSceneRuleHandler,
	NewProductCategoryHandler,
	NewStatisticsHandler,
	NewDeviceMessageHandler,
	NewDevicePropertyHandler,
	NewHandlers,
)

type Handlers struct {
	Product         *ProductHandler
	Device          *DeviceHandler
	ThingModel      *ThingModelHandler
	DeviceGroup     *DeviceGroupHandler
	OtaFirmware     *OtaFirmwareHandler
	OtaTask         *OtaTaskHandler
	AlertConfig     *AlertConfigHandler
	AlertRecord     *AlertRecordHandler
	DataSink        *DataSinkHandler
	DataRule        *DataRuleHandler
	SceneRule       *SceneRuleHandler
	ProductCategory *ProductCategoryHandler
	Statistics      *StatisticsHandler
	DeviceMessage   *DeviceMessageHandler
	DeviceProperty  *DevicePropertyHandler
}

func NewHandlers(
	product *ProductHandler,
	device *DeviceHandler,
	thingModel *ThingModelHandler,
	deviceGroup *DeviceGroupHandler,
	otaFirmware *OtaFirmwareHandler,
	otaTask *OtaTaskHandler,
	alertConfig *AlertConfigHandler,
	alertRecord *AlertRecordHandler,
	dataSink *DataSinkHandler,
	dataRule *DataRuleHandler,
	sceneRule *SceneRuleHandler,
	productCategory *ProductCategoryHandler,
	statistics *StatisticsHandler,
	deviceMessage *DeviceMessageHandler,
	deviceProperty *DevicePropertyHandler,
) *Handlers {
	return &Handlers{
		Product:         product,
		Device:          device,
		ThingModel:      thingModel,
		DeviceGroup:     deviceGroup,
		OtaFirmware:     otaFirmware,
		OtaTask:         otaTask,
		AlertConfig:     alertConfig,
		AlertRecord:     alertRecord,
		DataSink:        dataSink,
		DataRule:        dataRule,
		SceneRule:       sceneRule,
		ProductCategory: productCategory,
		Statistics:      statistics,
		DeviceMessage:   deviceMessage,
		DeviceProperty:  deviceProperty,
	}
}

// ProductHandler 产品处理器
type ProductHandler struct {
	svc                *iotsvc.ProductService
	productCategorySvc *iotsvc.ProductCategoryService
}

func NewProductHandler(svc *iotsvc.ProductService, productCategorySvc *iotsvc.ProductCategoryService) *ProductHandler {
	return &ProductHandler{
		svc:                svc,
		productCategorySvc: productCategorySvc,
	}
}

// DeviceHandler 设备处理器
type DeviceHandler struct {
	svc *iotsvc.DeviceService
}

func NewDeviceHandler(svc *iotsvc.DeviceService) *DeviceHandler {
	return &DeviceHandler{svc: svc}
}

// ThingModelHandler 物模型处理器
type ThingModelHandler struct {
	svc *iotsvc.ThingModelService
}

func NewThingModelHandler(svc *iotsvc.ThingModelService) *ThingModelHandler {
	return &ThingModelHandler{svc: svc}
}

// DeviceGroupHandler 设备分组处理器
type DeviceGroupHandler struct {
	svc *iotsvc.DeviceGroupService
}

func NewDeviceGroupHandler(svc *iotsvc.DeviceGroupService) *DeviceGroupHandler {
	return &DeviceGroupHandler{svc: svc}
}

// OtaFirmwareHandler OTA固件处理器
type OtaFirmwareHandler struct {
	svc        *iotsvc.OtaFirmwareService
	productSvc *iotsvc.ProductService
}

func NewOtaFirmwareHandler(svc *iotsvc.OtaFirmwareService, productSvc *iotsvc.ProductService) *OtaFirmwareHandler {
	return &OtaFirmwareHandler{
		svc:        svc,
		productSvc: productSvc,
	}
}

// OtaTaskHandler OTA任务处理器
type OtaTaskHandler struct {
	svc *iotsvc.OtaTaskService
}

func NewOtaTaskHandler(svc *iotsvc.OtaTaskService) *OtaTaskHandler {
	return &OtaTaskHandler{svc: svc}
}

// AlertConfigHandler 告警配置处理器
type AlertConfigHandler struct {
	svc *iotsvc.AlertConfigService
}

func NewAlertConfigHandler(svc *iotsvc.AlertConfigService) *AlertConfigHandler {
	return &AlertConfigHandler{svc: svc}
}

// AlertRecordHandler 告警记录处理器
type AlertRecordHandler struct {
	svc *iotsvc.AlertRecordService
}

func NewAlertRecordHandler(svc *iotsvc.AlertRecordService) *AlertRecordHandler {
	return &AlertRecordHandler{svc: svc}
}

// DataSinkHandler 数据目的处理器
type DataSinkHandler struct {
	svc *iotsvc.DataSinkService
}

func NewDataSinkHandler(svc *iotsvc.DataSinkService) *DataSinkHandler {
	return &DataSinkHandler{svc: svc}
}

// DataRuleHandler 数据规则处理器
type DataRuleHandler struct {
	svc *iotsvc.DataRuleService
}

func NewDataRuleHandler(svc *iotsvc.DataRuleService) *DataRuleHandler {
	return &DataRuleHandler{svc: svc}
}

// SceneRuleHandler 场景联动处理器
type SceneRuleHandler struct {
	svc *iotsvc.SceneRuleService
}

func NewSceneRuleHandler(svc *iotsvc.SceneRuleService) *SceneRuleHandler {
	return &SceneRuleHandler{svc: svc}
}

// ProductCategoryHandler 产品分类处理器
type ProductCategoryHandler struct {
	svc *iotsvc.ProductCategoryService
}

func NewProductCategoryHandler(svc *iotsvc.ProductCategoryService) *ProductCategoryHandler {
	return &ProductCategoryHandler{svc: svc}
}
