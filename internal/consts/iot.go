package consts

// IotDeviceStateEnum 设备状态
const (
	IotDeviceStateInactive = 0 // 未激活
	IotDeviceStateOnline   = 1 // 在线
	IotDeviceStateOffline  = 2 // 离线
)

// DateIntervalEnum 时间间隔枚举
const (
	DateIntervalHour    = 0 // 小时
	DateIntervalDay     = 1 // 天
	DateIntervalWeek    = 2 // 周
	DateIntervalMonth   = 3 // 月
	DateIntervalQuarter = 4 // 季度
	DateIntervalYear    = 5 // 年
)

// IotDeviceMessageMethodEnum 设备消息方法 (严格对齐 Java: cn.iocoder.yudao.module.iot.core.enums.IotDeviceMessageMethodEnum)
const (
	// ========== 设备状态 ==========
	IotDeviceMessageMethodStateUpdate = "thing.state.update" // 设备状态更新

	// ========== 设备属性 ==========
	IotDeviceMessageMethodPropertyPost = "thing.property.post" // 属性上报
	IotDeviceMessageMethodPropertySet  = "thing.property.set"  // 属性设置

	// ========== 设备事件 ==========
	IotDeviceMessageMethodEventPost = "thing.event.post" // 事件上报

	// ========== 设备服务调用 ==========
	IotDeviceMessageMethodServiceInvoke = "thing.service.invoke" // 服务调用

	// ========== 设备配置 ==========
	IotDeviceMessageMethodConfigPush = "thing.config.push" // 配置推送

	// ========== OTA 固件 ==========
	IotDeviceMessageMethodOtaUpgrade  = "thing.ota.upgrade"  // OTA 固定信息推送
	IotDeviceMessageMethodOtaProgress = "thing.ota.progress" // OTA 升级进度上报
)

// IotOtaTaskStatusEnum OTA 升级任务状态
const (
	IotOtaTaskStatusWait    = 0 // 待发布
	IotOtaTaskStatusRunning = 1 // 发布中
	IotOtaTaskStatusCancel  = 2 // 已取消
	IotOtaTaskStatusFinish  = 3 // 已完成
	IotOtaTaskStatusDone    = 3 // 已结束 (别名)
)

// IotOtaDeviceScopeEnum OTA 升级范围
const (
	IotOtaDeviceScopeAll       = 1 // 全部设备
	IotOtaDeviceScopeSpecified = 2 // 指定设备
)

// IotOtaRecordStatusEnum OTA 升级记录状态
const (
	IotOtaRecordStatusWait        = 0  // 等待升级
	IotOtaRecordStatusPushed      = 1  // 已推送
	IotOtaRecordStatusDownloading = 10 // 下载中
	IotOtaRecordStatusVerifying   = 20 // 校验中
	IotOtaRecordStatusUpgrading   = 30 // 升级中
	IotOtaRecordStatusSuccess     = 40 // 成功
	IotOtaRecordStatusFail        = 50 // 失败
	IotOtaRecordStatusCanceled    = 60 // 已取消
)

// IotProductStatusEnum IoT 产品状态
const (
	IotProductStatusUnpublished = 0 // 开发中
	IotProductStatusPublished   = 1 // 已发布
)

// IotThingModelTypeEnum IoT 产品功能（物模型）类型
const (
	IotThingModelTypeProperty = 1 // 属性
	IotThingModelTypeService  = 2 // 服务
	IotThingModelTypeEvent    = 3 // 事件
)

// IotProductDeviceTypeEnum IoT 产品的设备类型
const (
	IotProductDeviceTypeDirect  = 0 // 直连设备
	IotProductDeviceTypeSub     = 1 // 网关子设备
	IotProductDeviceTypeGateway = 2 // 网关设备
)

// IotAlertLevelEnum IoT 告警级别
const (
	IotAlertLevelInfo  = 1 // INFO
	IotAlertLevelWarn  = 3 // WARN
	IotAlertLevelError = 5 // ERROR
)

// IotAlertReceiveTypeEnum IoT 告警接收方式
const (
	IotAlertReceiveTypeSms    = 1 // 短信
	IotAlertReceiveTypeMail   = 2 // 邮箱
	IotAlertReceiveTypeNotify = 3 // 站内信
)

// IotDataSinkTypeEnum IoT 数据目的类型
const (
	IotDataSinkTypeHttp     = 1  // HTTP
	IotDataSinkTypeDatabase = 20 // Database
	IotDataSinkTypeRedis    = 21 // Redis
	IotDataSinkTypeRocketMQ = 30 // RocketMQ
	IotDataSinkTypeRabbitMQ = 31 // RabbitMQ
	IotDataSinkTypeKafka    = 32 // Kafka
)
