package handler

import (
	"strconv"

	"github.com/wxlbd/ruoyi-mall-go/internal/api/req"
	"github.com/wxlbd/ruoyi-mall-go/internal/service"
	"github.com/wxlbd/ruoyi-mall-go/pkg/errors"
	"github.com/wxlbd/ruoyi-mall-go/pkg/response"

	"github.com/gin-gonic/gin"
)

type SensitiveWordHandler struct {
	svc *service.SensitiveWordService
}

func NewSensitiveWordHandler(svc *service.SensitiveWordService) *SensitiveWordHandler {
	return &SensitiveWordHandler{svc: svc}
}

// CreateSensitiveWord 创建敏感词
func (h *SensitiveWordHandler) CreateSensitiveWord(c *gin.Context) {
	var r req.SensitiveWordCreateReq
	if err := c.ShouldBindJSON(&r); err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	id, err := h.svc.CreateSensitiveWord(c, &r)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, id)
}

// UpdateSensitiveWord 更新敏感词
func (h *SensitiveWordHandler) UpdateSensitiveWord(c *gin.Context) {
	var r req.SensitiveWordUpdateReq
	if err := c.ShouldBindJSON(&r); err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	if err := h.svc.UpdateSensitiveWord(c, &r); err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, true)
}

// DeleteSensitiveWord 删除敏感词
func (h *SensitiveWordHandler) DeleteSensitiveWord(c *gin.Context) {
	idStr := c.Query("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	if id == 0 {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	if err := h.svc.DeleteSensitiveWord(c, id); err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, true)
}

// GetSensitiveWord 获得敏感词
func (h *SensitiveWordHandler) GetSensitiveWord(c *gin.Context) {
	idStr := c.Query("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	if id == 0 {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	word, err := h.svc.GetSensitiveWord(c, id)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, word)
}

// GetSensitiveWordPage 获得敏感词分页
func (h *SensitiveWordHandler) GetSensitiveWordPage(c *gin.Context) {
	var r req.SensitiveWordPageReq
	if err := c.ShouldBindQuery(&r); err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	page, err := h.svc.GetSensitiveWordPage(c, &r)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, page)
}

// ValidateSensitiveWord 验证敏感词
func (h *SensitiveWordHandler) ValidateSensitiveWord(c *gin.Context) {
	text := c.Query("text")
	tag := c.Query("tag") // Single tag for simple test, or array
	var tags []string
	if tag != "" {
		tags = append(tags, tag)
	}

	words := h.svc.ValidateSensitiveWord(c, text, tags)
	response.WriteSuccess(c, words)
}

// ExportSensitiveWord 导出敏感词
func (h *SensitiveWordHandler) ExportSensitiveWord(c *gin.Context) {
	// TODO: Implement Export (Can limit only Page for now or reuse GetPage logic with large size)
	response.WriteSuccess(c, nil)
}
