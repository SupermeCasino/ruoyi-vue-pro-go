package handler

import (
	"github.com/wxlbd/ruoyi-mall-go/internal/api/req"
	"github.com/wxlbd/ruoyi-mall-go/internal/pkg/core"
	"github.com/wxlbd/ruoyi-mall-go/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type NoticeHandler struct {
	noticeSvc *service.NoticeService
}

func NewNoticeHandler(noticeSvc *service.NoticeService) *NoticeHandler {
	return &NoticeHandler{
		noticeSvc: noticeSvc,
	}
}

// GetNoticePage 获取通知公告分页
func (h *NoticeHandler) GetNoticePage(c *gin.Context) {
	var req req.NoticePageReq
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(400, core.Error(400, err.Error()))
		return
	}
	res, err := h.noticeSvc.GetNoticePage(c, &req)
	if err != nil {
		c.JSON(500, core.Error(500, err.Error()))
		return
	}
	c.JSON(200, core.Success(res))
}

// GetNotice 获得通知公告
func (h *NoticeHandler) GetNotice(c *gin.Context) {
	idStr := c.Query("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	if id == 0 {
		c.JSON(400, core.Error(400, "id is required"))
		return
	}
	res, err := h.noticeSvc.GetNotice(c, id)
	if err != nil {
		c.JSON(500, core.Error(500, err.Error()))
		return
	}
	c.JSON(200, core.Success(res))
}

// CreateNotice 创建通知公告
func (h *NoticeHandler) CreateNotice(c *gin.Context) {
	var req req.NoticeSaveReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, core.Error(400, err.Error()))
		return
	}
	id, err := h.noticeSvc.CreateNotice(c, &req)
	if err != nil {
		c.JSON(500, core.Error(500, err.Error()))
		return
	}
	c.JSON(200, core.Success(id))
}

// UpdateNotice 更新通知公告
func (h *NoticeHandler) UpdateNotice(c *gin.Context) {
	var req req.NoticeSaveReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, core.Error(400, err.Error()))
		return
	}
	if err := h.noticeSvc.UpdateNotice(c, &req); err != nil {
		c.JSON(500, core.Error(500, err.Error()))
		return
	}
	c.JSON(200, core.Success(true))
}

// DeleteNotice 删除通知公告
func (h *NoticeHandler) DeleteNotice(c *gin.Context) {
	idStr := c.Query("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	if id == 0 {
		c.JSON(400, core.Error(400, "id is required"))
		return
	}
	if err := h.noticeSvc.DeleteNotice(c, id); err != nil {
		c.JSON(500, core.Error(500, err.Error()))
		return
	}
	c.JSON(200, core.Success(true))
}

// Push 推送通知公告
func (h *NoticeHandler) Push(c *gin.Context) {
	idStr := c.Query("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	if id == 0 {
		c.JSON(400, core.Error(400, "id is required"))
		return
	}
	// TODO: WebSocket push not implemented yet
	c.JSON(200, core.Success(true))
}
