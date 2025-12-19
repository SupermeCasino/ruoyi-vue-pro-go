package trade

import (
	"time"

	"github.com/wxlbd/ruoyi-mall-go/internal/model"
)

// TradeDeliveryExpress 物流公司 DO
type TradeDeliveryExpress struct {
	ID        int64         `gorm:"primaryKey;autoIncrement;comment:编号" json:"id"`
	Code      string        `gorm:"size:64;not null;comment:物流编码" json:"code"`
	Name      string        `gorm:"size:64;not null;comment:物流名称" json:"name"`
	Logo      string        `gorm:"size:256;default:'';comment:物流Logo" json:"logo"`
	Sort      int           `gorm:"default:0;not null;comment:排序" json:"sort"`
	Status    int           `gorm:"default:0;not null;comment:状态" json:"status"`
	Creator   string        `gorm:"size:64;default:'';comment:创建者" json:"creator"`
	Updater   string        `gorm:"size:64;default:'';comment:更新者" json:"updater"`
	CreatedAt time.Time     `gorm:"column:create_time;autoCreateTime;comment:创建时间" json:"createTime"`
	UpdatedAt time.Time     `gorm:"column:update_time;autoUpdateTime;comment:更新时间" json:"updateTime"`
	Deleted   model.BitBool `gorm:"column:deleted;softDelete:flag;default:0;comment:是否删除" json:"deleted"`
	TenantID  int64         `gorm:"column:tenant_id;default:0;comment:租户编号" json:"tenantId"`
}

func (TradeDeliveryExpress) TableName() string {
	return "trade_delivery_express"
}

// TradeDeliveryPickUpStore 自提门店 DO
type TradeDeliveryPickUpStore struct {
	ID            int64         `gorm:"primaryKey;autoIncrement;comment:编号" json:"id"`
	Name          string        `gorm:"size:64;not null;comment:门店名称" json:"name"`
	Introduction  string        `gorm:"size:256;default:'';comment:门店简介" json:"introduction"`
	Phone         string        `gorm:"size:11;not null;comment:联系电话" json:"phone"`
	AreaID        int           `gorm:"column:area_id;not null;comment:区域编号" json:"areaId"`
	DetailAddress string        `gorm:"size:256;not null;comment:详细地址" json:"detailAddress"`
	Logo          string        `gorm:"size:256;not null;comment:门店Logo" json:"logo"`
	OpeningTime   *time.Time    `gorm:"column:opening_time;comment:营业开始时间" json:"openingTime"`
	ClosingTime   *time.Time    `gorm:"column:closing_time;comment:营业结束时间" json:"closingTime"`
	Latitude      float64       `gorm:"type:decimal(10,6);comment:纬度" json:"latitude"`
	Longitude     float64       `gorm:"type:decimal(10,6);comment:经度" json:"longitude"`
	VerifyUserIds model.IntListFromCSV `gorm:"column:verify_user_ids;type:varchar(500);comment:核销员工用户编号数组" json:"verifyUserIds"`
	Status        int           `gorm:"default:0;not null;comment:状态" json:"status"`
	Sort          int           `gorm:"default:0;not null;comment:排序" json:"sort"`
	Creator       string        `gorm:"size:64;default:'';comment:创建者" json:"creator"`
	Updater       string        `gorm:"size:64;default:'';comment:更新者" json:"updater"`
	CreatedAt     time.Time     `gorm:"column:create_time;autoCreateTime;comment:创建时间" json:"createTime"`
	UpdatedAt     time.Time     `gorm:"column:update_time;autoUpdateTime;comment:更新时间" json:"updateTime"`
	Deleted       model.BitBool `gorm:"column:deleted;softDelete:flag;default:0;comment:是否删除" json:"deleted"`
	TenantID      int64         `gorm:"column:tenant_id;default:0;comment:租户编号" json:"tenantId"`
}

func (TradeDeliveryPickUpStore) TableName() string {
	return "trade_delivery_pick_up_store"
}

// TradeDeliveryExpressTemplate 快递运费模板 DO
type TradeDeliveryExpressTemplate struct {
	ID         int64         `gorm:"primaryKey;autoIncrement;comment:编号" json:"id"`
	Name       string        `gorm:"size:64;not null;comment:模板名称" json:"name"`
	ChargeMode int           `gorm:"default:0;not null;comment:配送计费方式" json:"chargeMode"` // 1-按件 2-按重量 3-按体积
	Sort       int           `gorm:"default:0;not null;comment:排序" json:"sort"`
	Creator    string        `gorm:"size:64;default:'';comment:创建者" json:"creator"`
	Updater    string        `gorm:"size:64;default:'';comment:更新者" json:"updater"`
	CreatedAt  time.Time     `gorm:"column:create_time;autoCreateTime;comment:创建时间" json:"createTime"`
	UpdatedAt  time.Time     `gorm:"column:update_time;autoUpdateTime;comment:更新时间" json:"updateTime"`
	Deleted    model.BitBool `gorm:"column:deleted;softDelete:flag;default:0;comment:是否删除" json:"deleted"`
	TenantID   int64         `gorm:"column:tenant_id;default:0;comment:租户编号" json:"tenantId"`
}

func (TradeDeliveryExpressTemplate) TableName() string {
	return "trade_delivery_express_template"
}

// TradeDeliveryExpressTemplateCharge 快递运费模板计费规则 DO
type TradeDeliveryExpressTemplateCharge struct {
	ID         int64         `gorm:"primaryKey;autoIncrement;comment:编号" json:"id"`
	TemplateID int64         `gorm:"column:template_id;not null;comment:模板编号" json:"templateId"`
	AreaIDs    string        `gorm:"column:area_ids;type:text;not null;comment:区域编号列表" json:"areaIds"` // 逗号分隔
	ChargeMode int           `gorm:"column:charge_mode;not null;comment:配送计费方式" json:"chargeMode"`     // 1-按件 2-按重量 3-按体积
	StartCount float64       `gorm:"column:start_count;not null;comment:首件/首重/首体积" json:"startCount"`
	StartPrice int           `gorm:"column:start_price;not null;comment:首费(分)" json:"startPrice"`
	ExtraCount float64       `gorm:"column:extra_count;not null;comment:续件/续重/续体积" json:"extraCount"`
	ExtraPrice int           `gorm:"column:extra_price;not null;comment:续费(分)" json:"extraPrice"`
	Creator    string        `gorm:"size:64;default:'';comment:创建者" json:"creator"`
	Updater    string        `gorm:"size:64;default:'';comment:更新者" json:"updater"`
	CreatedAt  time.Time     `gorm:"column:create_time;autoCreateTime;comment:创建时间" json:"createTime"`
	UpdatedAt  time.Time     `gorm:"column:update_time;autoUpdateTime;comment:更新时间" json:"updateTime"`
	Deleted    model.BitBool `gorm:"column:deleted;softDelete:flag;default:0;comment:是否删除" json:"deleted"`
	TenantID   int64         `gorm:"column:tenant_id;default:0;comment:租户编号" json:"tenantId"`
}

func (TradeDeliveryExpressTemplateCharge) TableName() string {
	return "trade_delivery_express_template_charge"
}

// TradeDeliveryExpressTemplateFree 快递运费模板包邮规则 DO
type TradeDeliveryExpressTemplateFree struct {
	ID         int64         `gorm:"primaryKey;autoIncrement;comment:编号" json:"id"`
	TemplateID int64         `gorm:"column:template_id;not null;comment:模板编号" json:"templateId"`
	AreaIDs    string        `gorm:"column:area_ids;type:text;not null;comment:区域编号列表" json:"areaIds"` // 逗号分隔
	FreePrice  int           `gorm:"column:free_price;not null;comment:包邮金额(分)" json:"freePrice"`
	FreeCount  int           `gorm:"column:free_count;not null;comment:包邮件数/重量/体积" json:"freeCount"`
	Creator    string        `gorm:"size:64;default:'';comment:创建者" json:"creator"`
	Updater    string        `gorm:"size:64;default:'';comment:更新者" json:"updater"`
	CreatedAt  time.Time     `gorm:"column:create_time;autoCreateTime;comment:创建时间" json:"createTime"`
	UpdatedAt  time.Time     `gorm:"column:update_time;autoUpdateTime;comment:更新时间" json:"updateTime"`
	Deleted    model.BitBool `gorm:"column:deleted;softDelete:flag;default:0;comment:是否删除" json:"deleted"`
	TenantID   int64         `gorm:"column:tenant_id;default:0;comment:租户编号" json:"tenantId"`
}

func (TradeDeliveryExpressTemplateFree) TableName() string {
	return "trade_delivery_express_template_free"
}
