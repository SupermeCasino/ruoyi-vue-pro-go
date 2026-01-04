package system

import (
	"strconv"

	system2 "github.com/wxlbd/ruoyi-mall-go/internal/api/contract/admin/system"
	infra2 "github.com/wxlbd/ruoyi-mall-go/internal/api/handler/admin/infra"
	"github.com/wxlbd/ruoyi-mall-go/internal/service/system"
	"github.com/wxlbd/ruoyi-mall-go/pkg/errors"
	"github.com/wxlbd/ruoyi-mall-go/pkg/response"

	"github.com/gin-gonic/gin"
)

type NoticeHandler struct {
	noticeSvc *system.NoticeService
	wsHandler *infra2.WebSocketHandler  // WebSocket处理器，用于推送通知
}

func NewNoticeHandler(
	noticeSvc *system.NoticeService,
	wsHandler *infra2.WebSocketHandler,
) *NoticeHandler {
	return &NoticeHandler{
		noticeSvc: noticeSvc,
		wsHandler: wsHandler,
	}
}

// GetNoticePage 获取通知公告分页
func (h *NoticeHandler) GetNoticePage(c *gin.Context) {
	var req system2.NoticePageReq
	if err := c.ShouldBindQuery(&req); err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	res, err := h.noticeSvc.GetNoticePage(c, &req)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, res)
}

// GetNotice 获得通知公告
func (h *NoticeHandler) GetNotice(c *gin.Context) {
	idStr := c.Query("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	if id == 0 {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	res, err := h.noticeSvc.GetNotice(c, id)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, res)
}

// CreateNotice 创建通知公告
func (h *NoticeHandler) CreateNotice(c *gin.Context) {
	var req system2.NoticeSaveReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	id, err := h.noticeSvc.CreateNotice(c, &req)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, id)
}

// UpdateNotice 更新通知公告
func (h *NoticeHandler) UpdateNotice(c *gin.Context) {
	var req system2.NoticeSaveReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	if err := h.noticeSvc.UpdateNotice(c, &req); err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, true)
}

// DeleteNotice 删除通知公告
func (h *NoticeHandler) DeleteNotice(c *gin.Context) {
	idStr := c.Query("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	if id == 0 {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	if err := h.noticeSvc.DeleteNotice(c, id); err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, true)
}

// Push 推送通知公告 - 完整实现，与Java NoticeController.push() 对齐
// 对应Java: cn.iocoder.yudao.module.system.controller.admin.notice.NoticeController.push()
func (h *NoticeHandler) Push(c *gin.Context) {
	idStr := c.Query("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)

	// 1. 参数校验 (对应 Java 的 @RequestParam 校验)
	if id == 0 {
		response.WriteBizError(c, errors.ErrParam)
		return
	}

	// 2. 获取公告 (对应 Java: noticeService.getNotice(id))
	notice, err := h.noticeSvc.GetNotice(c, id)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}

	// 3. 校验公告存在性 (对应 Java: Assert.notNull(notice, "公告不能为空"))
	if notice == nil {
		response.WriteBizError(c, errors.NewBizError(errors.NotFoundCode, "公告不能为空"))
		return
	}

	// 4. 推送给所有ADMIN类型的在线用户
	// 对应 Java: webSocketSenderApi.sendObject(UserTypeEnum.ADMIN.getValue(), "notice-push", notice)
	// UserTypeEnum.ADMIN = 2
	err = h.wsHandler.BroadcastByUserType(2, "notice-push", notice)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}

	// 5. 返回成功 (对应 Java: return success(true))
	response.WriteSuccess(c, true)
}
