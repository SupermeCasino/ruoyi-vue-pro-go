package model

import (
	"github.com/wxlbd/ruoyi-mall-go/pkg/errors"
)

var (
	// ========== 产品相关 1-050-001-000 ============
	ErrProductNotExists                = errors.NewBizError(1050001000, "产品不存在")
	ErrProductKeyExists                = errors.NewBizError(1050001001, "产品标识已经存在")
	ErrProductStatusNotDelete          = errors.NewBizError(1050001002, "产品状是发布状态，不允许删除")
	ErrProductStatusNotAllowThingModel = errors.NewBizError(1050001003, "产品状是发布状态，不允许操作物模型")
	ErrProductDeleteFailHasDevice      = errors.NewBizError(1050001004, "产品下存在设备，不允许删除")

	// ========== 产品分类相关 1-050-002-000 ============
	ErrProductCategoryNotExists = errors.NewBizError(1050002000, "产品分类不存在")

	// ========== 数据流转 1-050-004-000 ============
	ErrDataSinkNotExists  = errors.NewBizError(1050004000, "数据目的不存在")
	ErrDataSinkUsedByRule = errors.NewBizError(1050004001, "数据目的已被数据规则使用，无法删除")
	ErrDataRuleNotExists  = errors.NewBizError(1050004100, "数据规则不存在")

	// ========== 物模型 1-050-005-000 ============
	ErrThingModelNotExists = errors.NewBizError(1050005000, "产品物模型不存在")

	// ========== 告警配置 1-050-006-000 ============
	ErrAlertConfigNotExists = errors.NewBizError(1050006000, "告警配置不存在")
	ErrAlertRecordNotExists = errors.NewBizError(1050006001, "告警记录不存在")

	// ========== 场景规则 1-050-007-000 ============
	ErrSceneRuleNotExists = errors.NewBizError(1050007000, "场景规则不存在")

	// ========== 设备 1-050-003-000 ============
	ErrDeviceNotExists          = errors.NewBizError(1050003000, "设备不存在")
	ErrDeviceNameExists         = errors.NewBizError(1050003001, "设备名称在同一产品下必须唯一")
	ErrDeviceHasChildren        = errors.NewBizError(1050003002, "有子设备，不允许删除")
	ErrDeviceKeyExists          = errors.NewBizError(1050003003, "设备标识已经存在")
	ErrDeviceGatewayNotExists   = errors.NewBizError(1050003004, "网关设备不存在")
	ErrDeviceNotGateway         = errors.NewBizError(1050003005, "设备不是网关设备")
	ErrDeviceImportListIsEmpty  = errors.NewBizError(1050003006, "导入设备数据不能为空！")
	ErrDeviceSerialNumberExists = errors.NewBizError(1050003008, "设备序列号已存在，序列号必须全局唯一")
	ErrDeviceSecretInvalid      = errors.NewBizError(1050003009, "设备密钥不正确")

	// ========== OTA 固件 1-050-013-000 ============
	ErrOtaFirmwareNotExists = errors.NewBizError(1050013000, "固件不存在")

	// ========== OTA 任务 1-050-014-000 ============
	ErrOtaTaskNotExists             = errors.NewBizError(1050014000, "任务不存在")
	ErrOtaTaskStatusNotAllowCancel  = errors.NewBizError(1050014001, "任务状态不支持取消")
	ErrOtaTaskRecordNotExists       = errors.NewBizError(1050014100, "升级记录不存在")
	ErrOtaTaskRecordUpdateFailNoRec = errors.NewBizError(1050014101, "无进行中的升级记录")
)
