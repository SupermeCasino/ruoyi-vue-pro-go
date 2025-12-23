package promotion

import (
	"github.com/wxlbd/ruoyi-mall-go/internal/model"

	"time"
)

// PromotionRewardActivity 满减送活动
// Table: promotion_reward_activity
type PromotionRewardActivity struct {
	ID                 int64     `gorm:"primaryKey;autoIncrement;comment:活动编号" json:"id"`
	Name               string    `gorm:"column:name;type:varchar(64);not null;comment:活动名称" json:"name"`
	Status             int       `gorm:"column:status;type:int;not null;default:0;comment:状态" json:"status"` // 0: 开启, 1: 关闭
	StartTime          time.Time `gorm:"column:start_time;not null;comment:开始时间" json:"startTime"`
	EndTime            time.Time `gorm:"column:end_time;not null;comment:结束时间" json:"endTime"`
	ProductScope       int       `gorm:"column:product_scope;type:int;not null;default:1;comment:商品范围" json:"productScope"`   // 1: 全部商品, 2: 指定商品, 3: 指定分类
	ProductScopeValues string    `gorm:"column:product_scope_values;type:json;comment:商品范围值" json:"productScopeValues"`       // Array of IDs
	ConditionType      int       `gorm:"column:condition_type;type:int;not null;default:1;comment:条件类型" json:"conditionType"` // 10: 满N元, 20: 满N件
	Rules              string    `gorm:"column:rules;type:json;not null;comment:优惠规则" json:"rules"`                           // List<Rule>
	Sort               int       `gorm:"column:sort;type:int;not null;default:0;comment:排序" json:"sort"`
	Remark             string    `gorm:"column:remark;type:varchar(255);comment:备注" json:"remark"`
	model.TenantBaseDO
}

func (PromotionRewardActivity) TableName() string {
	return "promotion_reward_activity"
}
