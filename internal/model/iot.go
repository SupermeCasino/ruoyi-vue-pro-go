package model

import (
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/datatypes"
)

// IotProductDO IoT 产品 DO
type IotProductDO struct {
	TenantBaseDO
	ID           int64  `gorm:"column:id;primaryKey;autoIncrement;comment:产品 ID" json:"id"`
	Name         string `gorm:"column:name;size:64;not null;comment:产品名称" json:"name"`
	ProductKey   string `gorm:"column:product_key;size:64;not null;comment:产品标识" json:"productKey"`
	CategoryID   int64  `gorm:"column:category_id;comment:产品分类编号" json:"categoryId"`
	Icon         string `gorm:"column:icon;size:255;comment:产品图标" json:"icon"`
	PicURL       string `gorm:"column:pic_url;size:255;comment:产品图片" json:"picUrl"`
	Description  string `gorm:"column:description;size:255;comment:产品描述" json:"description"`
	Status       int8   `gorm:"column:status;not null;default:0;comment:产品状态" json:"status"`
	DeviceType   int8   `gorm:"column:device_type;not null;default:0;comment:设备类型" json:"deviceType"`
	NetType      int8   `gorm:"column:net_type;not null;default:0;comment:联网方式" json:"netType"`
	LocationType int8   `gorm:"column:location_type;not null;default:0;comment:定位方式" json:"locationType"`
	CodecType    string `gorm:"column:codec_type;size:64;comment:数据格式" json:"codecType"`
}

// TableName 表名
func (IotProductDO) TableName() string {
	return "iot_product"
}

// IotDeviceDO IoT 设备 DO
type IotDeviceDO struct {
	TenantBaseDO
	ID           int64            `gorm:"column:id;primaryKey;autoIncrement;comment:设备 ID" json:"id"`
	DeviceName   string           `gorm:"column:device_name;size:64;not null;comment:设备名称" json:"deviceName"`
	Nickname     string           `gorm:"column:nickname;size:64;comment:设备备注名称" json:"nickname"`
	SerialNumber string           `gorm:"column:serial_number;size:64;comment:设备序列号" json:"serialNumber"`
	PicURL       string           `gorm:"column:pic_url;size:255;comment:设备图片" json:"picUrl"`
	GroupIDs     datatypes.JSON   `gorm:"column:group_ids;comment:设备分组编号集合" json:"groupIds"`
	ProductID    int64            `gorm:"column:product_id;not null;comment:产品编号" json:"productId"`
	ProductKey   string           `gorm:"column:product_key;size:64;not null;comment:产品标识" json:"productKey"`
	DeviceType   int8             `gorm:"column:device_type;not null;comment:设备类型" json:"deviceType"`
	GatewayID    int64            `gorm:"column:gateway_id;default:0;comment:网关设备编号" json:"gatewayId"`
	State        int8             `gorm:"column:state;not null;default:0;comment:设备状态" json:"state"`
	OnlineTime   *time.Time       `gorm:"column:online_time;comment:最后上线时间" json:"onlineTime"`
	OfflineTime  *time.Time       `gorm:"column:offline_time;comment:最后离线时间" json:"offlineTime"`
	ActiveTime   *time.Time       `gorm:"column:active_time;comment:设备激活时间" json:"activeTime"`
	IP           string           `gorm:"column:ip;size:64;comment:设备的 IP 地址" json:"ip"`
	FirmwareID   int64            `gorm:"column:firmware_id;comment:固件编号" json:"firmwareId"`
	DeviceSecret string           `gorm:"column:device_secret;size:64;comment:设备密钥" json:"deviceSecret"`
	AuthType     string           `gorm:"column:auth_type;size:64;comment:认证类型" json:"authType"`
	LocationType int8             `gorm:"column:location_type;comment:定位方式" json:"locationType"`
	Latitude     *decimal.Decimal `gorm:"column:latitude;type:decimal(10,8);comment:纬度" json:"latitude"`
	Longitude    *decimal.Decimal `gorm:"column:longitude;type:decimal(11,8);comment:经度" json:"longitude"`
	AreaID       int32            `gorm:"column:area_id;comment:地区编码" json:"areaId"`
	Address      string           `gorm:"column:address;size:255;comment:详细地址" json:"address"`
	Config       datatypes.JSON   `gorm:"column:config;type:text;comment:设备配置" json:"config"`
}

// TableName 表名
func (IotDeviceDO) TableName() string {
	return "iot_device"
}

// IotThingModelDO IoT 产品物模型功能 DO
type IotThingModelDO struct {
	BaseDO
	ID          int64          `gorm:"column:id;primaryKey;autoIncrement;comment:物模型功能编号" json:"id"`
	Identifier  string         `gorm:"column:identifier;size:64;not null;comment:功能标识" json:"identifier"`
	Name        string         `gorm:"column:name;size:64;not null;comment:功能名称" json:"name"`
	Description string         `gorm:"column:description;size:255;comment:功能描述" json:"description"`
	ProductID   int64          `gorm:"column:product_id;not null;comment:产品编号" json:"productId"`
	ProductKey  string         `gorm:"column:product_key;size:64;not null;comment:产品标识" json:"productKey"`
	Type        int8           `gorm:"column:type;not null;comment:功能类型" json:"type"`
	Property    datatypes.JSON `gorm:"column:property;type:text;comment:属性配置(JSON)" json:"property"`
	Event       datatypes.JSON `gorm:"column:event;type:text;comment:事件配置(JSON)" json:"event"`
	Service     datatypes.JSON `gorm:"column:service;type:text;comment:服务配置(JSON)" json:"service"`
}

// TableName 表名
func (IotThingModelDO) TableName() string {
	return "iot_thing_model"
}

// IotDeviceGroupDO IoT 设备分组 DO
type IotDeviceGroupDO struct {
	TenantBaseDO
	ID          int64  `gorm:"column:id;primaryKey;autoIncrement;comment:分组 ID" json:"id"`
	Name        string `gorm:"column:name;size:64;not null;comment:分组名字" json:"name"`
	Status      int8   `gorm:"column:status;not null;default:0;comment:分组状态" json:"status"`
	Description string `gorm:"column:description;size:255;comment:分组描述" json:"description"`
}

// TableName 表名
func (IotDeviceGroupDO) TableName() string {
	return "iot_device_group"
}

// IotOtaFirmwareDO IoT OTA 固件 DO
type IotOtaFirmwareDO struct {
	TenantBaseDO
	ID                  int64  `gorm:"column:id;primaryKey;autoIncrement;comment:固件编号" json:"id"`
	Name                string `gorm:"column:name;size:64;not null;comment:固件名称" json:"name"`
	Description         string `gorm:"column:description;size:255;comment:固件描述" json:"description"`
	Version             string `gorm:"column:version;size:64;not null;comment:版本号" json:"version"`
	ProductID           int64  `gorm:"column:product_id;not null;comment:产品编号" json:"productId"`
	FileURL             string `gorm:"column:file_url;size:255;not null;comment:固件文件 URL" json:"fileUrl"`
	FileSize            int64  `gorm:"column:file_size;not null;comment:固件文件大小" json:"fileSize"`
	FileDigestAlgorithm string `gorm:"column:file_digest_algorithm;size:32;not null;comment:固件文件签名算法" json:"fileDigestAlgorithm"`
	FileDigestValue     string `gorm:"column:file_digest_value;size:64;not null;comment:固件文件签名结果" json:"fileDigestValue"`
}

// TableName 表名
func (IotOtaFirmwareDO) TableName() string {
	return "iot_ota_firmware"
}

// IotOtaTaskDO IoT OTA 升级任务 DO
type IotOtaTaskDO struct {
	TenantBaseDO
	ID                 int64  `gorm:"column:id;primaryKey;autoIncrement;comment:任务编号" json:"id"`
	Name               string `gorm:"column:name;size:64;not null;comment:任务名称" json:"name"`
	Description        string `gorm:"column:description;size:255;comment:任务描述" json:"description"`
	FirmwareID         int64  `gorm:"column:firmware_id;not null;comment:固件编号" json:"firmwareId"`
	Status             int8   `gorm:"column:status;not null;default:0;comment:任务状态" json:"status"`
	DeviceScope        int8   `gorm:"column:device_scope;not null;comment:设备升级范围" json:"deviceScope"`
	DeviceTotalCount   int32  `gorm:"column:device_total_count;not null;default:0;comment:设备总数数量" json:"deviceTotalCount"`
	DeviceSuccessCount int32  `gorm:"column:device_success_count;not null;default:0;comment:设备成功数量" json:"deviceSuccessCount"`
}

// TableName 表名
func (IotOtaTaskDO) TableName() string {
	return "iot_ota_task"
}

// IotOtaTaskRecordDO IoT OTA 升级任务记录 DO
type IotOtaTaskRecordDO struct {
	TenantBaseDO
	ID             int64  `gorm:"column:id;primaryKey;autoIncrement;comment:升级记录编号" json:"id"`
	FirmwareID     int64  `gorm:"column:firmware_id;not null;comment:固件编号" json:"firmwareId"`
	TaskID         int64  `gorm:"column:task_id;not null;comment:任务编号" json:"taskId"`
	DeviceID       int64  `gorm:"column:device_id;not null;comment:设备编号" json:"deviceId"`
	FromFirmwareID int64  `gorm:"column:from_firmware_id;comment:来源的固件编号" json:"fromFirmwareId"`
	Status         int8   `gorm:"column:status;not null;default:0;comment:升级状态" json:"status"`
	Progress       int32  `gorm:"column:progress;not null;default:0;comment:升级进度" json:"progress"`
	Description    string `gorm:"column:description;size:255;comment:升级进度描述" json:"description"`
}

// TableName 表名
func (IotOtaTaskRecordDO) TableName() string {
	return "iot_ota_task_record"
}

// IotAlertConfigDO IoT 告警配置 DO
type IotAlertConfigDO struct {
	TenantBaseDO
	ID             int64          `gorm:"column:id;primaryKey;autoIncrement;comment:配置编号" json:"id"`
	Name           string         `gorm:"column:name;size:64;not null;comment:配置名称" json:"name"`
	Description    string         `gorm:"column:description;size:255;comment:配置描述" json:"description"`
	Level          int8           `gorm:"column:level;not null;comment:告警级别" json:"level"`
	Status         int8           `gorm:"column:status;not null;default:0;comment:配置状态" json:"status"`
	SceneRuleIDs   datatypes.JSON `gorm:"column:scene_rule_ids;size:255;comment:关联的场景联动规则编号数组" json:"sceneRuleIds"`
	ReceiveUserIDs datatypes.JSON `gorm:"column:receive_user_ids;size:255;comment:接收的用户编号数组" json:"receiveUserIds"`
	ReceiveTypes   datatypes.JSON `gorm:"column:receive_types;size:255;comment:接收的类型数组" json:"receiveTypes"`
}

// TableName 表名
func (IotAlertConfigDO) TableName() string {
	return "iot_alert_config"
}

// IotAlertRecordDO IoT 告警记录 DO
type IotAlertRecordDO struct {
	TenantBaseDO
	ID            int64  `gorm:"column:id;primaryKey;autoIncrement;comment:记录编号" json:"id"`
	ConfigID      int64  `gorm:"column:config_id;not null;comment:告警配置编号" json:"configId"`
	ConfigName    string `gorm:"column:config_name;size:64;not null;comment:告警配置名称" json:"configName"`
	ConfigLevel   int8   `gorm:"column:config_level;not null;comment:告警配置级别" json:"configLevel"`
	SceneRuleID   int64  `gorm:"column:scene_rule_id;comment:场景规则编号" json:"sceneRuleId"`
	ProductID     int64  `gorm:"column:product_id;not null;comment:产品编号" json:"productId"`
	DeviceID      int64  `gorm:"column:device_id;not null;comment:设备编号" json:"deviceId"`
	DeviceMessage string `gorm:"column:device_message;type:text;comment:触发的设备消息" json:"deviceMessage"`
	ProcessStatus bool   `gorm:"column:process_status;not null;default:false;comment:是否处理" json:"processStatus"`
	ProcessRemark string `gorm:"column:process_remark;size:255;comment:处理结果" json:"processRemark"`
}

// TableName 表名
func (IotAlertRecordDO) TableName() string {
	return "iot_alert_record"
}

// IotDataRuleDO IoT 数据流转规则 DO
type IotDataRuleDO struct {
	BaseDO
	ID            int64          `gorm:"column:id;primaryKey;autoIncrement;comment:规则编号" json:"id"`
	Name          string         `gorm:"column:name;size:64;not null;comment:规则名称" json:"name"`
	Description   string         `gorm:"column:description;size:255;comment:规则描述" json:"description"`
	Status        int8           `gorm:"column:status;not null;default:0;comment:状态" json:"status"`
	SourceConfigs datatypes.JSON `gorm:"column:source_configs;type:text;comment:数据源配置(JSON)" json:"sourceConfigs"`
	SinkIDs       datatypes.JSON `gorm:"column:sink_ids;size:255;comment:数据目的编号数组(JSON)" json:"sinkIds"`
}

// TableName 表名
func (IotDataRuleDO) TableName() string {
	return "iot_data_rule"
}

// IotDataSinkDO IoT 数据流转目的 DO
type IotDataSinkDO struct {
	BaseDO
	ID          int64          `gorm:"column:id;primaryKey;autoIncrement;comment:目的编号" json:"id"`
	Name        string         `gorm:"column:name;size:64;not null;comment:目的名称" json:"name"`
	Description string         `gorm:"column:description;size:255;comment:目的描述" json:"description"`
	Status      int8           `gorm:"column:status;not null;default:0;comment:状态" json:"status"`
	Type        int8           `gorm:"column:type;not null;comment:目的类型" json:"type"`
	Config      datatypes.JSON `gorm:"column:config;type:text;comment:目的配置(JSON)" json:"config"`
}

// TableName 表名
func (IotDataSinkDO) TableName() string {
	return "iot_data_sink"
}

// IotSceneRuleDO IoT 场景联动规则 DO
type IotSceneRuleDO struct {
	TenantBaseDO
	ID          int64          `gorm:"column:id;primaryKey;autoIncrement;comment:规则编号" json:"id"`
	Name        string         `gorm:"column:name;size:64;not null;comment:规则名称" json:"name"`
	Description string         `gorm:"column:description;size:255;comment:规则描述" json:"description"`
	Status      int8           `gorm:"column:status;not null;default:0;comment:状态" json:"status"`
	Triggers    datatypes.JSON `gorm:"column:triggers;type:text;comment:触发器配置(JSON)" json:"triggers"`
	Actions     datatypes.JSON `gorm:"column:actions;type:text;comment:动作配置(JSON)" json:"actions"`
}

// TableName 表名
func (IotSceneRuleDO) TableName() string {
	return "iot_scene_rule"
}

// IotProductCategoryDO IoT 产品分类 DO
type IotProductCategoryDO struct {
	BaseDO
	ID          int64  `gorm:"column:id;primaryKey;autoIncrement;comment:分类 ID" json:"id"`
	Name        string `gorm:"column:name;size:64;not null;comment:分类名字" json:"name"`
	Sort        int32  `gorm:"column:sort;comment:分类排序" json:"sort"`
	Status      int8   `gorm:"column:status;not null;default:0;comment:分类状态" json:"status"`
	Description string `gorm:"column:description;size:255;comment:分类描述" json:"description"`
}

// TableName 表名
func (IotProductCategoryDO) TableName() string {
	return "iot_product_category"
}

// IotDeviceMessageDO IoT 设备消息 DO (暂用 MySQL 存储)
type IotDeviceMessageDO struct {
	TenantBaseDO
	ID         string `gorm:"column:id;primaryKey;size:64;comment:消息编号" json:"id"`
	ReportTime int64  `gorm:"column:report_time;not null;comment:上报时间戳" json:"reportTime"`
	TS         int64  `gorm:"column:ts;not null;comment:存储时间戳" json:"ts"`
	DeviceID   int64  `gorm:"column:device_id;not null;comment:设备编号" json:"deviceId"`
	ServerID   string `gorm:"column:server_id;size:64;comment:服务编号" json:"serverId"`
	Upstream   bool   `gorm:"column:upstream;not null;comment:是否上行消息" json:"upstream"`
	Reply      bool   `gorm:"column:reply;not null;comment:是否回复消息" json:"reply"`
	Identifier string `gorm:"column:identifier;size:64;comment:标识符" json:"identifier"`
	RequestID  string `gorm:"column:request_id;size:64;comment:请求编号" json:"requestId"`
	Method     string `gorm:"column:method;size:64;comment:请求方法" json:"method"`
	Params     string `gorm:"column:params;type:text;comment:请求参数" json:"params"`
	Data       string `gorm:"column:data;type:text;comment:响应结果" json:"data"`
	Code       int    `gorm:"column:code;comment:响应错误码" json:"code"`
	Msg        string `gorm:"column:msg;size:255;comment:响应提示" json:"msg"`
}

// TableName 表名
func (IotDeviceMessageDO) TableName() string {
	return "iot_device_message"
}

// IotDevicePropertyDO IoT 设备属性历史记录 DO
type IotDevicePropertyDO struct {
	TenantBaseDO
	ID         int64      `gorm:"column:id;primaryKey;autoIncrement;comment:记录编号" json:"id"`
	DeviceID   int64      `gorm:"column:device_id;not null;index;comment:设备编号" json:"deviceId"`
	Identifier string     `gorm:"column:identifier;size:64;not null;index;comment:属性标识符" json:"identifier"`
	Value      string     `gorm:"column:value;type:text;not null;comment:属性值" json:"value"`
	UpdateTime *time.Time `gorm:"column:update_time;not null;comment:更新时间" json:"updateTime"`
}

// TableName 表名
func (IotDevicePropertyDO) TableName() string {
	return "iot_device_property"
}
