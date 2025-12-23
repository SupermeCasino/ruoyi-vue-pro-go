package member

import (
	"time"

	"github.com/wxlbd/ruoyi-mall-go/internal/model"
)

// MemberPointRecord 用户积分记录
// Table: member_point_record
type MemberPointRecord struct {
	ID          int64         `gorm:"primaryKey;autoIncrement;comment:自增主键" json:"id"`
	UserID      int64         `gorm:"column:user_id;not null;comment:用户编号" json:"userId"`
	BizID       string        `gorm:"column:biz_id;size:64;comment:业务编码" json:"bizId"`
	BizType     int           `gorm:"column:biz_type;not null;comment:业务类型" json:"bizType"` // MemberPointBizTypeEnum
	Title       string        `gorm:"column:title;size:64;not null;comment:积分标题" json:"title"`
	Description string        `gorm:"column:description;size:255;default:'';comment:积分描述" json:"description"`
	Point       int           `gorm:"column:point;not null;comment:变动积分" json:"point"`              // 1、正数表示获得积分 2、负数表示消耗积分
	TotalPoint  int           `gorm:"column:total_point;not null;comment:变动后的积分" json:"totalPoint"` // 变动后的积分
	Creator     string        `gorm:"column:creator;size:64;default:'';comment:创建者"`
	Updater     string        `gorm:"column:updater;size:64;default:'';comment:更新者"`
	CreateTime   time.Time     `gorm:"column:create_time;autoCreateTime;comment:创建时间"`
	UpdateTime   time.Time     `gorm:"column:update_time;autoUpdateTime;comment:更新时间"`
	Deleted     model.BitBool `gorm:"column:deleted;type:tinyint(1);not null;default:0;comment:是否删除"`
	TenantID    int64         `gorm:"column:tenant_id;default:0;comment:租户编号" json:"tenantId"`
}

func (MemberPointRecord) TableName() string {
	return "member_point_record"
}

// MemberPointBizType 会员积分的业务类型
// 对应 Java: MemberPointBizTypeEnum
type MemberPointBizType struct {
	Type        int    // 类型
	Name        string // 名字（标题）
	Description string // 描述模板，{} 占位符表示积分值
	Add         bool   // 是否为增加积分
}

// MemberPointBizType 常量定义
var (
	// MemberPointBizTypeSign 签到
	MemberPointBizTypeSign = MemberPointBizType{Type: 1, Name: "签到", Description: "签到获得 {} 积分", Add: true}
	// MemberPointBizTypeAdmin 管理员修改
	MemberPointBizTypeAdmin = MemberPointBizType{Type: 2, Name: "管理员修改", Description: "管理员修改 {} 积分", Add: true}
	// MemberPointBizTypeOrderUse 订单积分抵扣
	MemberPointBizTypeOrderUse = MemberPointBizType{Type: 11, Name: "订单积分抵扣", Description: "下单使用 {} 积分", Add: false}
	// MemberPointBizTypeOrderUseCancel 订单积分抵扣（整单取消）
	MemberPointBizTypeOrderUseCancel = MemberPointBizType{Type: 12, Name: "订单积分抵扣（整单取消）", Description: "订单取消，退还 {} 积分", Add: true}
	// MemberPointBizTypeOrderUseCancelItem 订单积分抵扣（单个退款）
	MemberPointBizTypeOrderUseCancelItem = MemberPointBizType{Type: 13, Name: "订单积分抵扣（单个退款）", Description: "订单退款，退还 {} 积分", Add: true}
	// MemberPointBizTypeOrderGive 订单积分奖励
	MemberPointBizTypeOrderGive = MemberPointBizType{Type: 21, Name: "订单积分奖励", Description: "下单获得 {} 积分", Add: true}
	// MemberPointBizTypeOrderGiveCancel 订单积分奖励（整单取消）
	MemberPointBizTypeOrderGiveCancel = MemberPointBizType{Type: 22, Name: "订单积分奖励（整单取消）", Description: "订单取消，扣除赠送的 {} 积分", Add: false}
	// MemberPointBizTypeOrderGiveCancelItem 订单积分奖励（单个退款）
	MemberPointBizTypeOrderGiveCancelItem = MemberPointBizType{Type: 23, Name: "订单积分奖励（单个退款）", Description: "订单退款，扣除赠送的 {} 积分", Add: false}
)

// GetMemberPointBizTypeByType 根据类型获取业务类型
func GetMemberPointBizTypeByType(bizType int) *MemberPointBizType {
	switch bizType {
	case 1:
		return &MemberPointBizTypeSign
	case 2:
		return &MemberPointBizTypeAdmin
	case 11:
		return &MemberPointBizTypeOrderUse
	case 12:
		return &MemberPointBizTypeOrderUseCancel
	case 13:
		return &MemberPointBizTypeOrderUseCancelItem
	case 21:
		return &MemberPointBizTypeOrderGive
	case 22:
		return &MemberPointBizTypeOrderGiveCancel
	case 23:
		return &MemberPointBizTypeOrderGiveCancelItem
	default:
		return nil
	}
}
