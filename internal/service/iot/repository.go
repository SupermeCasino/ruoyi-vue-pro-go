package iot

import (
	"context"

	"time"

	"github.com/wxlbd/ruoyi-mall-go/internal/api/contract/admin/iot"
	"github.com/wxlbd/ruoyi-mall-go/internal/model"
	"github.com/wxlbd/ruoyi-mall-go/pkg/pagination"
)

type ProductRepository interface {
	Create(ctx context.Context, product *model.IotProductDO) error
	Update(ctx context.Context, product *model.IotProductDO) error
	Delete(ctx context.Context, id int64) error
	GetByID(ctx context.Context, id int64) (*model.IotProductDO, error)
	GetByKey(ctx context.Context, key string) (*model.IotProductDO, error)
	GetPage(ctx context.Context, req *iot.IotProductPageReqVO) (*pagination.PageResult[*model.IotProductDO], error)
	ListAll(ctx context.Context) ([]*model.IotProductDO, error)
	Count(ctx context.Context, startTime *time.Time) (int64, error)
}

type DeviceRepository interface {
	Create(ctx context.Context, device *model.IotDeviceDO) error
	Update(ctx context.Context, device *model.IotDeviceDO) error
	Delete(ctx context.Context, id int64) error
	DeleteList(ctx context.Context, ids []int64) error
	GetByID(ctx context.Context, id int64) (*model.IotDeviceDO, error)
	GetPage(ctx context.Context, req *iot.IotDevicePageReqVO) (*pagination.PageResult[*model.IotDeviceDO], error)
	CountByProductID(ctx context.Context, productID int64) (int64, error)
	CountByGatewayID(ctx context.Context, gatewayID int64) (int64, error)
	ListByCondition(ctx context.Context, deviceType *int8, productID *int64) ([]*model.IotDeviceDO, error)
	ListByProductKeyAndNames(ctx context.Context, productKey string, names []string) ([]*model.IotDeviceDO, error)
	GetByProductKeyAndName(ctx context.Context, productKey string, name string) (*model.IotDeviceDO, error)
	Count(ctx context.Context, startTime *time.Time) (int64, error)
	GetStateCountMap(ctx context.Context) (map[int8]int64, error)
}

type ThingModelRepository interface {
	Create(ctx context.Context, tm *model.IotThingModelDO) error
	Update(ctx context.Context, tm *model.IotThingModelDO) error
	Delete(ctx context.Context, id int64) error
	GetByID(ctx context.Context, id int64) (*model.IotThingModelDO, error)
	ListByProductID(ctx context.Context, productID int64) ([]*model.IotThingModelDO, error)
	ListByProductIDAndType(ctx context.Context, productID int64, tmType int8) ([]*model.IotThingModelDO, error)
	GetPage(ctx context.Context, req *iot.IotThingModelPageReqVO) (*pagination.PageResult[*model.IotThingModelDO], error)
}

type DeviceGroupRepository interface {
	Create(ctx context.Context, group *model.IotDeviceGroupDO) error
	Update(ctx context.Context, group *model.IotDeviceGroupDO) error
	Delete(ctx context.Context, id int64) error
	GetByID(ctx context.Context, id int64) (*model.IotDeviceGroupDO, error)
	GetPage(ctx context.Context, req *iot.IotDeviceGroupPageReqVO) (*pagination.PageResult[*model.IotDeviceGroupDO], error)
	ListByStatus(ctx context.Context, status int8) ([]*model.IotDeviceGroupDO, error)
}

type OtaFirmwareRepository interface {
	Create(ctx context.Context, firmware *model.IotOtaFirmwareDO) error
	Update(ctx context.Context, firmware *model.IotOtaFirmwareDO) error
	Delete(ctx context.Context, id int64) error
	GetByID(ctx context.Context, id int64) (*model.IotOtaFirmwareDO, error)
	GetPage(ctx context.Context, req *iot.IotOtaFirmwarePageReqVO) (*pagination.PageResult[*model.IotOtaFirmwareDO], error)
}

type OtaTaskRepository interface {
	Create(ctx context.Context, task *model.IotOtaTaskDO) error
	Update(ctx context.Context, task *model.IotOtaTaskDO) error
	Delete(ctx context.Context, id int64) error
	GetByID(ctx context.Context, id int64) (*model.IotOtaTaskDO, error)
	GetPage(ctx context.Context, req *iot.IotOtaTaskPageReqVO) (*pagination.PageResult[*model.IotOtaTaskDO], error)
}

type OtaTaskRecordRepository interface {
	Create(ctx context.Context, records []*model.IotOtaTaskRecordDO) error
	Update(ctx context.Context, record *model.IotOtaTaskRecordDO) error
	GetPage(ctx context.Context, req *iot.IotOtaTaskRecordPageReqVO) (*pagination.PageResult[*model.IotOtaTaskRecordDO], error)
	CreateBatch(ctx context.Context, records []*model.IotOtaTaskRecordDO) error
	GetListByDeviceIdAndStatus(ctx context.Context, deviceID int64, statuses []int) ([]*model.IotOtaTaskRecordDO, error)
	GetListByTaskIdAndStatus(ctx context.Context, taskID int64, statuses []int) ([]*model.IotOtaTaskRecordDO, error)
}

type AlertConfigRepository interface {
	Create(ctx context.Context, config *model.IotAlertConfigDO) error
	Update(ctx context.Context, config *model.IotAlertConfigDO) error
	Delete(ctx context.Context, id int64) error
	GetByID(ctx context.Context, id int64) (*model.IotAlertConfigDO, error)
	GetPage(ctx context.Context, req *iot.IotAlertConfigPageReqVO) (*pagination.PageResult[*model.IotAlertConfigDO], error)
	GetListByStatus(ctx context.Context, status int8) ([]*model.IotAlertConfigDO, error)
	GetListBySceneRuleIdAndStatus(ctx context.Context, sceneRuleID int64, status int8) ([]*model.IotAlertConfigDO, error)
}

type AlertRecordRepository interface {
	Create(ctx context.Context, record *model.IotAlertRecordDO) error
	Update(ctx context.Context, record *model.IotAlertRecordDO) error
	GetByID(ctx context.Context, id int64) (*model.IotAlertRecordDO, error)
	GetPage(ctx context.Context, req *iot.IotAlertRecordPageReqVO) (*pagination.PageResult[*model.IotAlertRecordDO], error)
	GetListBySceneRuleId(ctx context.Context, sceneRuleID int64, deviceID *int64, processStatus *bool) ([]*model.IotAlertRecordDO, error)
}

type DataRuleRepository interface {
	Create(ctx context.Context, rule *model.IotDataRuleDO) error
	Update(ctx context.Context, rule *model.IotDataRuleDO) error
	Delete(ctx context.Context, id int64) error
	GetByID(ctx context.Context, id int64) (*model.IotDataRuleDO, error)
	GetPage(ctx context.Context, req *iot.IotDataRulePageReqVO) (*pagination.PageResult[*model.IotDataRuleDO], error)
}

type DataSinkRepository interface {
	Create(ctx context.Context, sink *model.IotDataSinkDO) error
	Update(ctx context.Context, sink *model.IotDataSinkDO) error
	Delete(ctx context.Context, id int64) error
	GetByID(ctx context.Context, id int64) (*model.IotDataSinkDO, error)
	GetPage(ctx context.Context, req *iot.IotDataSinkPageReqVO) (*pagination.PageResult[*model.IotDataSinkDO], error)
	CountBySinkID(ctx context.Context, id int64) (int64, error)
	GetListByStatus(ctx context.Context, status int8) ([]*model.IotDataSinkDO, error)
}

type SceneRuleRepository interface {
	Create(ctx context.Context, rule *model.IotSceneRuleDO) error
	Update(ctx context.Context, rule *model.IotSceneRuleDO) error
	Delete(ctx context.Context, id int64) error
	GetByID(ctx context.Context, id int64) (*model.IotSceneRuleDO, error)
	GetPage(ctx context.Context, req *iot.IotSceneRulePageReqVO) (*pagination.PageResult[*model.IotSceneRuleDO], error)
	CountBySceneRuleID(ctx context.Context, id int64) (int64, error)
	GetListByStatus(ctx context.Context, status int8) ([]*model.IotSceneRuleDO, error)
}

type ProductCategoryRepository interface {
	Create(ctx context.Context, category *model.IotProductCategoryDO) error
	Update(ctx context.Context, category *model.IotProductCategoryDO) error
	Delete(ctx context.Context, id int64) error
	GetByID(ctx context.Context, id int64) (*model.IotProductCategoryDO, error)
	GetPage(ctx context.Context, req *iot.IotProductCategoryPageReqVO) (*pagination.PageResult[*model.IotProductCategoryDO], error)
	GetListByStatus(ctx context.Context, status int8) ([]*model.IotProductCategoryDO, error)
	Count(ctx context.Context, startTime *time.Time) (int64, error)
	GetProductCategoryDeviceCountMap(ctx context.Context) (map[string]int64, error)
}

type DeviceMessageRepository interface {
	Create(ctx context.Context, message *model.IotDeviceMessageDO) error
	Count(ctx context.Context, startTime *time.Time) (int64, error)
	GetSummaryByDate(ctx context.Context, interval int, startTime, endTime *time.Time) ([]*iot.IotStatisticsDeviceMessageSummaryByDateRespVO, error)
	GetPage(ctx context.Context, req *iot.IotDeviceMessagePageReqVO) (*pagination.PageResult[*model.IotDeviceMessageDO], error)
	GetListByRequestIdsAndReply(ctx context.Context, deviceID int64, requestIDs []string, reply bool) ([]*model.IotDeviceMessageDO, error)
}

type DevicePropertyRepository interface {
	GetHistoryList(ctx context.Context, req *iot.IotDevicePropertyHistoryListReqVO) ([]*model.IotDevicePropertyDO, error)
	GetLatestProperties(ctx context.Context, deviceID int64) ([]*model.IotDevicePropertyDO, error)
	SaveProperty(ctx context.Context, property *model.IotDevicePropertyDO) error
}
