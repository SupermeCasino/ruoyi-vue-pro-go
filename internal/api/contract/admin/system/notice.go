package system

import (
	"time"

	"github.com/wxlbd/ruoyi-mall-go/pkg/pagination"
)

type NoticeSaveReq struct {
	ID      int64  `json:"id"`
	// 对应Java: @NotBlank(message = "公告标题不能为空") + @Size(max = 50, message = "公告标题不能超过50个字符")
	Title   string `json:"title" binding:"required,max=50"`
	// 对应Java: @NotNull(message = "公告类型不能为空")
	Type    *int32 `json:"type" binding:"required,min=1,max=2"`
	// 对应Java: 公告内容无特殊验证，但应该限制长度
	Content string `json:"content" binding:"required,max=5000"`
	// 对应Java: 状态为CommonStatusEnum，有效值为0(关闭)和1(开启)
	Status  *int32 `json:"status" binding:"required,min=0,max=1"`
}

type NoticePageReq struct {
	pagination.PageParam
	Title  string `form:"title"`
	Status *int32 `form:"status"`
}

type NoticeRespVO struct {
	ID         int64     `json:"id"`
	Title      string    `json:"title"`
	Type       int32     `json:"type"`
	Content    string    `json:"content"`
	Status     int32     `json:"status"`
	CreateTime time.Time `json:"createTime"`
	// ← 以下字段对齐Java的NoticeRespVO (继承自BaseDO的审计字段)
	UpdateTime time.Time `json:"updateTime"`  // 更新时间 (对应Java: BaseDO.updateTime)
	Creator    string    `json:"creator"`     // 创建人 (对应Java: BaseDO.creator)
	Updater    string    `json:"updater"`     // 更新人 (对应Java: BaseDO.updater)
}
