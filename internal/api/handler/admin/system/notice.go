package system

import (
	"strconv"

	"github.com/wxlbd/ruoyi-mall-go/internal/api/req"
	"github.com/wxlbd/ruoyi-mall-go/internal/service/system"
	"github.com/wxlbd/ruoyi-mall-go/pkg/errors"
	"github.com/wxlbd/ruoyi-mall-go/pkg/response"

	"github.com/gin-gonic/gin"
)

type NoticeHandler struct {
	noticeSvc *system.NoticeService
}

func NewNoticeHandler(noticeSvc *system.NoticeService) *NoticeHandler {
	return &NoticeHandler{
		noticeSvc: noticeSvc,
	}
}

// GetNoticePage 获取通知公告分页
func (h *NoticeHandler) GetNoticePage(c *gin.Context) {
	var req req.NoticePageReq
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
	var req req.NoticeSaveReq
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
	var req req.NoticeSaveReq
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

// Push 推送通知公告
func (h *NoticeHandler) Push(c *gin.Context) {
	idStr := c.Query("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	if id == 0 {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	// TODO: WebSocket push not implemented yet
	response.WriteSuccess(c, true)
}
