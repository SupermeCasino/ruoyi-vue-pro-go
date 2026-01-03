package promotion

import (
	"strconv"

	promotion2 "github.com/wxlbd/ruoyi-mall-go/internal/api/contract/admin/mall/promotion"
	"github.com/wxlbd/ruoyi-mall-go/internal/service/mall/promotion"
	"github.com/wxlbd/ruoyi-mall-go/pkg/errors"
	"github.com/wxlbd/ruoyi-mall-go/pkg/response"

	"github.com/gin-gonic/gin"
)

type DiscountActivityHandler struct {
	svc promotion.DiscountActivityService
}

func NewDiscountActivityHandler(svc promotion.DiscountActivityService) *DiscountActivityHandler {
	return &DiscountActivityHandler{svc: svc}
}

func (h *DiscountActivityHandler) CreateDiscountActivity(c *gin.Context) {
	var r promotion2.DiscountActivityCreateReq
	if err := c.ShouldBindJSON(&r); err != nil {
		response.WriteBizError(c, errors.NewBizError(400, "参数错误"))
		return
	}
	id, err := h.svc.CreateDiscountActivity(c, r)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, id)
}

func (h *DiscountActivityHandler) UpdateDiscountActivity(c *gin.Context) {
	var r promotion2.DiscountActivityUpdateReq
	if err := c.ShouldBindJSON(&r); err != nil {
		response.WriteBizError(c, errors.NewBizError(400, "参数错误"))
		return
	}
	if err := h.svc.UpdateDiscountActivity(c, r); err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, true)
}

func (h *DiscountActivityHandler) CloseDiscountActivity(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Query("id"), 10, 64)
	if id == 0 {
		response.WriteBizError(c, errors.NewBizError(400, "参数错误"))
		return
	}
	if err := h.svc.CloseDiscountActivity(c, id); err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, true)
}

func (h *DiscountActivityHandler) DeleteDiscountActivity(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Query("id"), 10, 64)
	if id == 0 {
		response.WriteBizError(c, errors.NewBizError(400, "参数错误"))
		return
	}
	if err := h.svc.DeleteDiscountActivity(c, id); err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, true)
}

func (h *DiscountActivityHandler) GetDiscountActivity(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Query("id"), 10, 64)
	if id == 0 {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	res, err := h.svc.GetDiscountActivity(c, id)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, res)
}

func (h *DiscountActivityHandler) GetDiscountActivityPage(c *gin.Context) {
	var r promotion2.DiscountActivityPageReq
	if err := c.ShouldBindQuery(&r); err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	res, err := h.svc.GetDiscountActivityPage(c, r)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, res)
}
