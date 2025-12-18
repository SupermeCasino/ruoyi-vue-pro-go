package trade

import (
	"time"

	"github.com/wxlbd/ruoyi-mall-go/internal/model"
)

// AfterSaleLog 售后日志
type AfterSaleLog struct {
	ID           int64         `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	UserID       int64         `gorm:"column:user_id" json:"userId"`
	UserType     int           `gorm:"column:user_type" json:"userType"`
	AfterSaleID  int64         `gorm:"column:after_sale_id" json:"afterSaleId"`
	BeforeStatus int           `gorm:"column:before_status" json:"beforeStatus"`
	AfterStatus  int           `gorm:"column:after_status" json:"afterStatus"`
	OperateType  int           `gorm:"column:operate_type" json:"operateType"`
	Content      string        `gorm:"column:content" json:"content"`
	Creator      string        `gorm:"column:creator" json:"creator"`
	Updater      string        `gorm:"column:updater" json:"updater"`
	CreatedAt    time.Time     `gorm:"column:create_time" json:"createTime"`
	UpdatedAt    time.Time     `gorm:"column:update_time" json:"updateTime"`
	Deleted      model.BitBool `gorm:"column:deleted" json:"deleted"`
}

func (AfterSaleLog) TableName() string {
	return "trade_after_sale_log"
}
