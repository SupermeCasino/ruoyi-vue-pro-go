package consts

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

// MemberExperienceBizType 会员经验的业务类型
// 对应 Java: MemberExperienceBizTypeEnum
type MemberExperienceBizType struct {
	Type        int    // 类型
	Name        string // 名字（标题）
	Description string // 描述模板，{} 占位符表示经验值
	Add         bool   // 是否为增加经验
}

// MemberExperienceBizType 常量定义
var (
	// MemberExperienceBizTypeAdmin 管理员调整
	MemberExperienceBizTypeAdmin = MemberExperienceBizType{Type: 0, Name: "管理员调整", Description: "管理员调整获得 {} 经验", Add: true}
	// MemberExperienceBizTypeInviteRegister 邀新奖励
	MemberExperienceBizTypeInviteRegister = MemberExperienceBizType{Type: 1, Name: "邀新奖励", Description: "邀请好友获得 {} 经验", Add: true}
	// MemberExperienceBizTypeSignIn 签到奖励
	MemberExperienceBizTypeSignIn = MemberExperienceBizType{Type: 4, Name: "签到奖励", Description: "签到获得 {} 经验", Add: true}
	// MemberExperienceBizTypeLottery 抽奖奖励
	MemberExperienceBizTypeLottery = MemberExperienceBizType{Type: 5, Name: "抽奖奖励", Description: "抽奖获得 {} 经验", Add: true}
	// MemberExperienceBizTypeOrderGive 下单奖励
	MemberExperienceBizTypeOrderGive = MemberExperienceBizType{Type: 11, Name: "下单奖励", Description: "下单获得 {} 经验", Add: true}
	// MemberExperienceBizTypeOrderGiveCancel 下单奖励（整单取消）
	MemberExperienceBizTypeOrderGiveCancel = MemberExperienceBizType{Type: 12, Name: "下单奖励（整单取消）", Description: "取消订单获得 {} 经验", Add: false}
	// MemberExperienceBizTypeOrderGiveCancelItem 下单奖励（单个退款）
	MemberExperienceBizTypeOrderGiveCancelItem = MemberExperienceBizType{Type: 13, Name: "下单奖励（单个退款）", Description: "退款订单获得 {} 经验", Add: false}
)

// GetMemberExperienceBizTypeByType 根据类型获取业务类型
func GetMemberExperienceBizTypeByType(bizType int) *MemberExperienceBizType {
	switch bizType {
	case 0:
		return &MemberExperienceBizTypeAdmin
	case 1:
		return &MemberExperienceBizTypeInviteRegister
	case 4:
		return &MemberExperienceBizTypeSignIn
	case 5:
		return &MemberExperienceBizTypeLottery
	case 11:
		return &MemberExperienceBizTypeOrderGive
	case 12:
		return &MemberExperienceBizTypeOrderGiveCancel
	case 13:
		return &MemberExperienceBizTypeOrderGiveCancelItem
	default:
		return nil
	}
}
