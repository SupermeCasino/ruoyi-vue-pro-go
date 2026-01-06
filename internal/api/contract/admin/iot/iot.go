package iot

import (
	"time"

	"github.com/shopspring/decimal"
	"github.com/wxlbd/ruoyi-mall-go/internal/dto"
	"github.com/wxlbd/ruoyi-mall-go/pkg/pagination"
)

// ================= Iot Product =================

// IotProductSaveReqVO 产品保存请求
type IotProductSaveReqVO struct {
	ID           int64  `json:"id"`
	Name         string `json:"name" binding:"required"`
	ProductKey   string `json:"productKey"`
	CategoryID   int64  `json:"categoryId" binding:"required"`
	Icon         string `json:"icon"`
	PicURL       string `json:"picUrl"`
	Description  string `json:"description"`
	DeviceType   int8   `json:"deviceType" binding:"required"`
	NetType      int8   `json:"netType"`
	LocationType int8   `json:"locationType"`
	CodecType    string `json:"codecType" binding:"required"`
}

// IotProductRespVO 产品响应信息
type IotProductRespVO struct {
	ID           int64     `json:"id"`
	Name         string    `json:"name"`
	ProductKey   string    `json:"productKey"`
	CategoryID   int64     `json:"categoryId"`
	CategoryName string    `json:"categoryName"`
	Icon         string    `json:"icon"`
	PicURL       string    `json:"picUrl"`
	Description  string    `json:"description"`
	Status       int8      `json:"status"`
	DeviceType   int8      `json:"deviceType"`
	NetType      int8      `json:"netType"`
	LocationType int8      `json:"locationType"`
	CodecType    string    `json:"codecType"`
	CreateTime   time.Time `json:"createTime"`
}

// IotProductPageReqVO 产品分页请求
type IotProductPageReqVO struct {
	PageNo     int    `form:"pageNo" binding:"required"`
	PageSize   int    `form:"pageSize" binding:"required"`
	Name       string `form:"name"`
	ProductKey string `form:"productKey"`
}

// ================= Iot Device =================

// IotDeviceSaveReqVO 设备保存请求
type IotDeviceSaveReqVO struct {
	ID           int64            `json:"id"`
	DeviceName   string           `json:"deviceName"`
	Nickname     string           `json:"nickname"`
	SerialNumber string           `json:"serialNumber"`
	PicURL       string           `json:"picUrl"`
	GroupIDs     []int64          `json:"groupIds"`
	ProductID    int64            `json:"productId"`
	GatewayID    int64            `json:"gatewayId"`
	Config       string           `json:"config"`
	LocationType int8             `json:"locationType" binding:"required"`
	Latitude     *decimal.Decimal `json:"latitude"`
	Longitude    *decimal.Decimal `json:"longitude"`
}

// IotDeviceRespVO 设备响应信息
type IotDeviceRespVO struct {
	ID           int64            `json:"id"`
	DeviceName   string           `json:"deviceName"`
	Nickname     string           `json:"nickname"`
	SerialNumber string           `json:"serialNumber"`
	PicURL       string           `json:"picUrl"`
	GroupIDs     []int64          `json:"groupIds"`
	ProductID    int64            `json:"productId"`
	ProductKey   string           `json:"productKey"`
	DeviceType   int8             `json:"deviceType"`
	GatewayID    int64            `json:"gatewayId"`
	State        int8             `json:"state"`
	OnlineTime   *time.Time       `json:"onlineTime"`
	OfflineTime  *time.Time       `json:"offlineTime"`
	ActiveTime   *time.Time       `json:"activeTime"`
	DeviceSecret string           `json:"deviceSecret"`
	AuthType     string           `json:"authType"`
	Config       string           `json:"config"`
	LocationType int8             `json:"locationType"`
	Latitude     *decimal.Decimal `json:"latitude"`
	Longitude    *decimal.Decimal `json:"longitude"`
	CreateTime   time.Time        `json:"createTime"`
}

// IotDevicePageReqVO 设备分页请求
type IotDevicePageReqVO struct {
	PageNo     int    `form:"pageNo" binding:"required"`
	PageSize   int    `form:"pageSize" binding:"required"`
	DeviceName string `form:"deviceName"`
	Nickname   string `form:"nickname"`
	ProductID  int64  `form:"productId"`
	DeviceType int8   `form:"deviceType"`
	Status     int8   `form:"status"`
	GroupID    int64  `form:"groupId"`
}

// IotDeviceUpdateGroupReqVO 设备更新分组请求
type IotDeviceUpdateGroupReqVO struct {
	IDs      []int64 `json:"ids" binding:"required"`
	GroupIDs []int64 `json:"groupIds"`
}

// IotDeviceAuthInfoRespVO 设备认证信息响应
type IotDeviceAuthInfoRespVO struct {
	ProductKey   string `json:"productKey"`
	DeviceName   string `json:"deviceName"`
	DeviceSecret string `json:"deviceSecret"`
	MqttHost     string `json:"mqttHost"`
	MqttPort     int    `json:"mqttPort"`
}

// IotDeviceByProductKeyAndNamesReqVO 根据产品Key和设备名称查询请求
type IotDeviceByProductKeyAndNamesReqVO struct {
	ProductKey  string   `form:"productKey" binding:"required"`
	DeviceNames []string `form:"deviceNames" binding:"required"`
}

// IotDeviceImportExcelVO 设备 Excel 导入 VO
type IotDeviceImportExcelVO struct {
	DeviceName       string `json:"deviceName"`
	ParentDeviceName string `json:"parentDeviceName"`
	ProductKey       string `json:"productKey"`
	GroupNames       string `json:"groupNames"`
	LocationType     int8   `json:"locationType"`
}

// IotDeviceImportRespVO 设备导入响应
type IotDeviceImportRespVO struct {
	CreateDeviceNames  []string          `json:"createDeviceNames"`
	UpdateDeviceNames  []string          `json:"updateDeviceNames"`
	FailureDeviceNames map[string]string `json:"failureDeviceNames"`
}

// IotDeviceAuthReqDTO 设备认证请求 DTO (内部业务逻辑使用)
type IotDeviceAuthReqDTO struct {
	ProductKey   string `json:"productKey" binding:"required"`
	DeviceName   string `json:"deviceName" binding:"required"`
	DeviceSecret string `json:"deviceSecret"`
}

// ================= Iot ThingModel =================

// IotThingModelSaveReqVO 物模型保存请求
type IotThingModelSaveReqVO struct {
	ID          int64                   `json:"id"`
	Identifier  string                  `json:"identifier" binding:"required"`
	Name        string                  `json:"name" binding:"required"`
	Description string                  `json:"description"`
	ProductID   int64                   `json:"productId" binding:"required"`
	ProductKey  string                  `json:"productKey"`
	Type        int8                    `json:"type" binding:"required"`
	Property    *dto.ThingModelProperty `json:"property"`
	Event       *dto.ThingModelEvent    `json:"event"`
	Service     *dto.ThingModelService  `json:"service"`
}

// IotThingModelRespVO 物模型响应信息
type IotThingModelRespVO struct {
	ID          int64                   `json:"id"`
	Identifier  string                  `json:"identifier"`
	Name        string                  `json:"name"`
	Description string                  `json:"description"`
	ProductID   int64                   `json:"productId"`
	ProductKey  string                  `json:"productKey"`
	Type        int8                    `json:"type"`
	Property    *dto.ThingModelProperty `json:"property,omitempty"`
	Event       *dto.ThingModelEvent    `json:"event,omitempty"`
	Service     *dto.ThingModelService  `json:"service,omitempty"`
	CreateTime  time.Time               `json:"createTime"`
}

// IotThingModelTSLRespVO 物模型 TSL 响应
type IotThingModelTSLRespVO struct {
	ProductID  int64                    `json:"productId"`
	ProductKey string                   `json:"productKey"`
	Properties []dto.ThingModelProperty `json:"properties"`
	Services   []dto.ThingModelService  `json:"services"`
	Events     []dto.ThingModelEvent    `json:"events"`
}

// IotThingModelPageReqVO 物模型分页请求
type IotThingModelPageReqVO struct {
	PageNo    int   `form:"pageNo" binding:"required"`
	PageSize  int   `form:"pageSize" binding:"required"`
	ProductID int64 `form:"productId" binding:"required"`
}

// IotThingModelListReqVO 物模型列表请求
type IotThingModelListReqVO struct {
	ProductID int64 `form:"productId" binding:"required"`
}

// ================= Iot Device Group =================

// IotDeviceGroupSaveReqVO 设备分组保存请求
type IotDeviceGroupSaveReqVO struct {
	ID          int64  `json:"id"`
	Name        string `json:"name" binding:"required"`
	Status      int8   `json:"status"`
	Description string `json:"description"`
}

// IotDeviceGroupRespVO 设备分组响应信息
type IotDeviceGroupRespVO struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Status      int8      `json:"status"`
	Description string    `json:"description"`
	CreateTime  time.Time `json:"createTime"`
}

// IotDeviceGroupPageReqVO 设备分组分页请求
type IotDeviceGroupPageReqVO struct {
	PageNo   int    `form:"pageNo" binding:"required"`
	PageSize int    `form:"pageSize" binding:"required"`
	Name     string `form:"name"`
	Status   int8   `form:"status"`
}

// ================= Iot OTA Firmware =================

// IotOtaFirmwareSaveReqVO 固件保存请求 (创建/更新)
type IotOtaFirmwareSaveReqVO struct {
	ID                  int64  `json:"id"`
	Name                string `json:"name" binding:"required"`
	Description         string `json:"description"`
	Version             string `json:"version" binding:"required"`
	ProductID           int64  `json:"productId" binding:"required"`
	FileURL             string `json:"fileUrl" binding:"required"`
	FileSize            int64  `json:"fileSize" binding:"required"`
	FileDigestAlgorithm string `json:"fileDigestAlgorithm" binding:"required"`
	FileDigestValue     string `json:"fileDigestValue" binding:"required"`
}

// IotOtaFirmwareRespVO 固件响应信息
type IotOtaFirmwareRespVO struct {
	ID                  int64     `json:"id"`
	Name                string    `json:"name"`
	Description         string    `json:"description"`
	Version             string    `json:"version"`
	ProductID           int64     `json:"productId"`
	ProductName         string    `json:"productName"`
	FileURL             string    `json:"fileUrl"`
	FileSize            int64     `json:"fileSize"`
	FileDigestAlgorithm string    `json:"fileDigestAlgorithm"`
	FileDigestValue     string    `json:"fileDigestValue"`
	CreateTime          time.Time `json:"createTime"`
}

// IotOtaFirmwarePageReqVO 固件分页请求
type IotOtaFirmwarePageReqVO struct {
	PageNo    int    `form:"pageNo" binding:"required"`
	PageSize  int    `form:"pageSize" binding:"required"`
	Name      string `form:"name"`
	ProductID int64  `form:"productId"`
}

// ================= Iot OTA Task =================

// IotOtaTaskCreateReqVO 固件任务创建请求
type IotOtaTaskCreateReqVO struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description"`
	FirmwareID  int64   `json:"firmwareId" binding:"required"`
	DeviceScope int8    `json:"deviceScope" binding:"required"`
	DeviceIDs   []int64 `json:"deviceIds"`
}

// IotOtaTaskRespVO 固件任务响应信息
type IotOtaTaskRespVO struct {
	ID                 int64     `json:"id"`
	Name               string    `json:"name"`
	Description        string    `json:"description"`
	FirmwareID         int64     `json:"firmwareId"`
	Status             int8      `json:"status"`
	DeviceScope        int8      `json:"deviceScope"`
	DeviceTotalCount   int32     `json:"deviceTotalCount"`
	DeviceSuccessCount int32     `json:"deviceSuccessCount"`
	CreateTime         time.Time `json:"createTime"`
}

// IotOtaTaskPageReqVO 固件任务分页请求
type IotOtaTaskPageReqVO struct {
	PageNo   int    `form:"pageNo" binding:"required"`
	PageSize int    `form:"pageSize" binding:"required"`
	Name     string `form:"name"`
}

// IotOtaTaskRecordRespVO 固件任务记录响应信息
type IotOtaTaskRecordRespVO struct {
	ID             int64     `json:"id"`
	FirmwareID     int64     `json:"firmwareId"`
	TaskID         int64     `json:"taskId"`
	DeviceID       int64     `json:"deviceId"`
	FromFirmwareID int64     `json:"fromFirmwareId"`
	Status         int8      `json:"status"`
	Progress       int32     `json:"progress"`
	Description    string    `json:"description"`
	CreateTime     time.Time `json:"createTime"`
}

// IotOtaTaskRecordPageReqVO 固件任务记录分页请求
type IotOtaTaskRecordPageReqVO struct {
	PageNo   int   `form:"pageNo" binding:"required"`
	PageSize int   `form:"pageSize" binding:"required"`
	TaskID   int64 `form:"taskId" binding:"required"`
	Status   int8  `form:"status"`
}

// ================= Iot Alert Management =================

// IotAlertConfigSaveReqVO 告警配置保存请求
type IotAlertConfigSaveReqVO struct {
	ID             int64   `json:"id"`
	Name           string  `json:"name" binding:"required"`
	Description    string  `json:"description"`
	Level          int8    `json:"level" binding:"required"`
	Status         int8    `json:"status"`
	SceneRuleIDs   []int64 `json:"sceneRuleIds"`
	ReceiveUserIDs []int64 `json:"receiveUserIds"`
	ReceiveTypes   []int   `json:"receiveTypes"`
}

// IotAlertConfigRespVO 告警配置响应信息
type IotAlertConfigRespVO struct {
	ID               int64     `json:"id"`
	Name             string    `json:"name"`
	Description      string    `json:"description"`
	Level            int8      `json:"level"`
	Status           int8      `json:"status"`
	SceneRuleIDs     []int64   `json:"sceneRuleIds"`
	ReceiveUserIDs   []int64   `json:"receiveUserIds"`
	ReceiveUserNames []string  `json:"receiveUserNames"`
	ReceiveTypes     []int     `json:"receiveTypes"`
	CreateTime       time.Time `json:"createTime"`
}

// IotAlertConfigPageReqVO 告警配置分页请求
type IotAlertConfigPageReqVO struct {
	PageNo   int    `form:"pageNo" binding:"required"`
	PageSize int    `form:"pageSize" binding:"required"`
	Name     string `form:"name"`
	Status   int8   `form:"status"`
}

// IotAlertRecordRespVO 告警记录响应信息
type IotAlertRecordRespVO struct {
	ID            int64     `json:"id"`
	ConfigID      int64     `json:"configId"`
	ConfigName    string    `json:"configName"`
	ConfigLevel   int8      `json:"configLevel"`
	SceneRuleID   int64     `json:"sceneRuleId"`
	ProductID     int64     `json:"productId"`
	DeviceID      int64     `json:"deviceId"`
	DeviceMessage string    `json:"deviceMessage"`
	ProcessStatus bool      `json:"processStatus"`
	ProcessRemark string    `json:"processRemark"`
	CreateTime    time.Time `json:"createTime"`
}

// IotAlertRecordPageReqVO 告警记录分页请求
type IotAlertRecordPageReqVO struct {
	PageNo        int   `form:"pageNo" binding:"required"`
	PageSize      int   `form:"pageSize" binding:"required"`
	ConfigID      int64 `form:"configId"`
	ProcessStatus *bool `form:"processStatus"`
}

// IotAlertRecordProcessReqVO 告警记录处理请求
type IotAlertRecordProcessReqVO struct {
	ID            int64  `json:"id" binding:"required"`
	ProcessRemark string `json:"processRemark" binding:"required"`
}

// ================= Iot Data Rule =================

// IotDataRuleSourceConfig 数据源配置
type IotDataRuleSourceConfig struct {
	Method     string `json:"method" binding:"required"`
	ProductID  int64  `json:"productId"`
	DeviceID   int64  `json:"deviceId" binding:"required"`
	Identifier string `json:"identifier"`
}

// IotDataRuleSaveReqVO 数据流转规则保存请求
type IotDataRuleSaveReqVO struct {
	ID            int64                     `json:"id"`
	Name          string                    `json:"name" binding:"required"`
	Description   string                    `json:"description"`
	Status        int8                      `json:"status" binding:"required"`
	SourceConfigs []IotDataRuleSourceConfig `json:"sourceConfigs" binding:"required"`
	SinkIDs       []int64                   `json:"sinkIds" binding:"required"`
}

// IotDataRuleRespVO 数据流转规则响应信息
type IotDataRuleRespVO struct {
	ID            int64                     `json:"id"`
	Name          string                    `json:"name"`
	Description   string                    `json:"description"`
	Status        int8                      `json:"status"`
	SourceConfigs []IotDataRuleSourceConfig `json:"sourceConfigs"`
	SinkIDs       []int64                   `json:"sinkIds"`
	CreateTime    time.Time                 `json:"createTime"`
}

// IotDataRulePageReqVO 数据流转规则分页请求
type IotDataRulePageReqVO struct {
	PageNo   int    `form:"pageNo" binding:"required"`
	PageSize int    `form:"pageSize" binding:"required"`
	Name     string `form:"name"`
	Status   int8   `form:"status"`
}

// ================= Iot Data Sink =================

// IotDataSinkSaveReqVO 数据流转目的保存请求
type IotDataSinkSaveReqVO struct {
	ID          int64       `json:"id"`
	Name        string      `json:"name" binding:"required"`
	Description string      `json:"description"`
	Status      int8        `json:"status" binding:"required"`
	Type        int8        `json:"type" binding:"required"`
	Config      interface{} `json:"config" binding:"required"`
}

// IotDataSinkRespVO 数据流转目的响应信息
type IotDataSinkRespVO struct {
	ID          int64       `json:"id"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Status      int8        `json:"status"`
	Type        int8        `json:"type"`
	Config      interface{} `json:"config"`
	CreateTime  time.Time   `json:"createTime"`
}

// IotDataSinkPageReqVO 数据流转目的分页请求
type IotDataSinkPageReqVO struct {
	PageNo   int    `form:"pageNo" binding:"required"`
	PageSize int    `form:"pageSize" binding:"required"`
	Name     string `form:"name"`
	Status   int8   `form:"status"`
}

// ================= Iot Scene Rule =================

// IotSceneRuleTriggerCondition 触发条件
type IotSceneRuleTriggerCondition struct {
	Type       int8   `json:"type"`
	ProductID  int64  `json:"productId"`
	DeviceID   int64  `json:"deviceId"`
	Identifier string `json:"identifier"`
	Operator   string `json:"operator"`
	Param      string `json:"param"`
}

// IotSceneRuleTrigger 场景联动触发器
type IotSceneRuleTrigger struct {
	Type            int8                             `json:"type"`
	ProductID       int64                            `json:"productId"`
	DeviceID        int64                            `json:"deviceId"`
	Identifier      string                           `json:"identifier"`
	Operator        string                           `json:"operator"`
	Value           string                           `json:"value"`
	CronExpression  string                           `json:"cronExpression"`
	ConditionGroups [][]IotSceneRuleTriggerCondition `json:"conditionGroups"`
}

// IotSceneRuleAction 场景联动动作
type IotSceneRuleAction struct {
	Type          int8   `json:"type"`
	ProductID     int64  `json:"productId"`
	DeviceID      int64  `json:"deviceId"`
	Identifier    string `json:"identifier"`
	Params        string `json:"params"`
	AlertConfigID int64  `json:"alertConfigId"`
}

// IotSceneRuleSaveReqVO 场景联动保存请求
type IotSceneRuleSaveReqVO struct {
	ID          int64                 `json:"id"`
	Name        string                `json:"name" binding:"required"`
	Description string                `json:"description"`
	Status      int8                  `json:"status" binding:"required"`
	Triggers    []IotSceneRuleTrigger `json:"triggers" binding:"required"`
	Actions     []IotSceneRuleAction  `json:"actions" binding:"required"`
}

// IotSceneRuleRespVO 场景联动响应信息
type IotSceneRuleRespVO struct {
	ID          int64                 `json:"id"`
	Name        string                `json:"name"`
	Description string                `json:"description"`
	Status      int8                  `json:"status"`
	Triggers    []IotSceneRuleTrigger `json:"triggers"`
	Actions     []IotSceneRuleAction  `json:"actions"`
	CreateTime  time.Time             `json:"createTime"`
}

// IotSceneRulePageReqVO 场景联动分页请求
type IotSceneRulePageReqVO struct {
	PageNo   int    `form:"pageNo" binding:"required"`
	PageSize int    `form:"pageSize" binding:"required"`
	Name     string `form:"name"`
	Status   int8   `form:"status"`
}

// ================= Iot Product Category =================

// IotProductCategorySaveReqVO 产品分类保存请求
type IotProductCategorySaveReqVO struct {
	ID          int64  `json:"id"`
	Name        string `json:"name" binding:"required"`
	Sort        int32  `json:"sort"`
	Status      int8   `json:"status" binding:"required"`
	Description string `json:"description"`
}

// IotProductCategoryRespVO 产品分类响应信息
type IotProductCategoryRespVO struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Sort        int32     `json:"sort"`
	Status      int8      `json:"status"`
	Description string    `json:"description"`
	CreateTime  time.Time `json:"createTime"`
}

// IotProductCategoryPageReqVO 产品分类分页请求
type IotProductCategoryPageReqVO struct {
	PageNo   int    `form:"pageNo" binding:"required"`
	PageSize int    `form:"pageSize" binding:"required"`
	Name     string `form:"name"`
}

// ================= Iot Statistics =================

// IotStatisticsSummaryRespVO 统计摘要响应
type IotStatisticsSummaryRespVO struct {
	ProductCategoryCount        int64            `json:"productCategoryCount"`
	ProductCount                int64            `json:"productCount"`
	DeviceCount                 int64            `json:"deviceCount"`
	DeviceMessageCount          int64            `json:"deviceMessageCount"`
	ProductCategoryTodayCount   int64            `json:"productCategoryTodayCount"`
	ProductTodayCount           int64            `json:"productTodayCount"`
	DeviceTodayCount            int64            `json:"deviceTodayCount"`
	DeviceMessageTodayCount     int64            `json:"deviceMessageTodayCount"`
	DeviceOnlineCount           int64            `json:"deviceOnlineCount"`
	DeviceOfflineCount          int64            `json:"deviceOfflineCount"`
	DeviceInactiveCount         int64            `json:"deviceInactiveCount"`
	ProductCategoryDeviceCounts map[string]int64 `json:"productCategoryDeviceCounts"`
}

// IotStatisticsDeviceMessageReqVO 设备消息统计请求
type IotStatisticsDeviceMessageReqVO struct {
	Interval int          `form:"interval" binding:"required"`
	Times    []*time.Time `form:"times"`
}

// IotStatisticsDeviceMessageSummaryByDateRespVO 分日期设备消息统计响应
type IotStatisticsDeviceMessageSummaryByDateRespVO struct {
	Time            string `json:"time"`
	UpstreamCount   int64  `json:"upstreamCount"`
	DownstreamCount int64  `json:"downstreamCount"`
}

// ================= Iot Device Message =================

// IotDeviceMessagePageReqVO IoT 设备消息分页查询 Request VO
type IotDeviceMessagePageReqVO struct {
	pagination.PageParam
	DeviceID   int64        `form:"deviceId" binding:"required"`
	Method     string       `form:"method"`
	Upstream   *bool        `form:"upstream"`
	Reply      *bool        `form:"reply"` // 注意：Java Controller 中 getDeviceMessagePairPage 会显式 setReply(false)
	Identifier string       `form:"identifier"`
	Times      []*time.Time `form:"times" time_format:"2006-01-02 15:04:05"`
}

// IotDeviceMessageRespVO IoT 设备消息 Response VO
type IotDeviceMessageRespVO struct {
	ID         string     `json:"id"`
	ReportTime *time.Time `json:"reportTime"`
	TS         *time.Time `json:"ts"`
	DeviceID   int64      `json:"deviceId"`
	ServerID   string     `json:"serverId"`
	Upstream   bool       `json:"upstream"`
	Reply      bool       `json:"reply"`
	Identifier string     `json:"identifier"`
	// codec（编解码）字段
	RequestID string      `json:"requestId"`
	Method    string      `json:"method"`
	Params    interface{} `json:"params"`
	Data      interface{} `json:"data"`
	Code      int         `json:"code"`
	Msg       string      `json:"msg"`
}

// IotDeviceMessageRespPairVO IoT 设备消息对 Response VO
type IotDeviceMessageRespPairVO struct {
	Request *IotDeviceMessageRespVO `json:"request"`
	Reply   *IotDeviceMessageRespVO `json:"reply"`
}

// IotDeviceMessageSendReqVO IoT 设备消息发送 Request VO
type IotDeviceMessageSendReqVO struct {
	Method   string      `json:"method" binding:"required"`
	Params   interface{} `json:"params"`
	DeviceID int64       `json:"deviceId" binding:"required"`
}

// ================= Iot Device Property =================

// IotDevicePropertyRespVO IoT 设备属性 Response VO
type IotDevicePropertyRespVO struct {
	Identifier string      `json:"identifier"`
	Value      interface{} `json:"value"`
	UpdateTime int64       `json:"updateTime"`
}

// IotDevicePropertyDetailRespVO IoT 设备属性详细 Response VO
type IotDevicePropertyDetailRespVO struct {
	IotDevicePropertyRespVO
	Name          string                    `json:"name"`
	DataType      string                    `json:"dataType"`
	DataSpecs     *dto.ThingModelDataSpecs  `json:"dataSpecs"`
	DataSpecsList []dto.ThingModelDataSpecs `json:"dataSpecsList"`
}

// IotDevicePropertyHistoryListReqVO IoT 设备属性历史列表 Request VO
type IotDevicePropertyHistoryListReqVO struct {
	DeviceID   int64        `form:"deviceId" binding:"required"`
	Identifier string       `form:"identifier" binding:"required"`
	Times      []*time.Time `form:"times" time_format:"2006-01-02 15:04:05"`
}
