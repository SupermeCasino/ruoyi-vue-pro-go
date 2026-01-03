package pay

import (
	"strconv"

	"github.com/wxlbd/ruoyi-mall-go/internal/api/contract/admin/pay"
	paySvc "github.com/wxlbd/ruoyi-mall-go/internal/service/pay"
	"github.com/wxlbd/ruoyi-mall-go/pkg/errors"
	"github.com/wxlbd/ruoyi-mall-go/pkg/response"

	"github.com/gin-gonic/gin"
)

type PayAppHandler struct {
	svc *paySvc.PayAppService
}

func NewPayAppHandler(svc *paySvc.PayAppService) *PayAppHandler {
	return &PayAppHandler{svc: svc}
}

// CreateApp 创建支付应用
func (h *PayAppHandler) CreateApp(c *gin.Context) {
	var r pay.PayAppCreateReq
	if err := c.ShouldBindJSON(&r); err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	id, err := h.svc.CreateApp(c, &r)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, id)
}

// UpdateApp 更新支付应用
func (h *PayAppHandler) UpdateApp(c *gin.Context) {
	var r pay.PayAppUpdateReq
	if err := c.ShouldBindJSON(&r); err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	err := h.svc.UpdateApp(c, &r)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, true)
}

// UpdateAppStatus 更新支付应用状态
func (h *PayAppHandler) UpdateAppStatus(c *gin.Context) {
	var r pay.PayAppUpdateStatusReq
	if err := c.ShouldBindJSON(&r); err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	err := h.svc.UpdateAppStatus(c, &r)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, true)
}

// DeleteApp 删除支付应用
func (h *PayAppHandler) DeleteApp(c *gin.Context) {
	idStr := c.Query("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	err = h.svc.DeleteApp(c, id)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, true)
}

// GetApp 获得支付应用
func (h *PayAppHandler) GetApp(c *gin.Context) {
	idStr := c.Query("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	app, err := h.svc.GetApp(c, id)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}

	res := &pay.PayAppResp{
		ID:                app.ID,
		AppKey:            app.AppKey,
		Name:              app.Name,
		Status:            app.Status,
		Remark:            app.Remark,
		OrderNotifyURL:    app.OrderNotifyURL,
		RefundNotifyURL:   app.RefundNotifyURL,
		TransferNotifyURL: app.TransferNotifyURL,
		CreateTime:        app.CreateTime,
	}
	response.WriteSuccess(c, res)
}

// GetAppList 获得支付应用列表
func (h *PayAppHandler) GetAppList(c *gin.Context) {
	list, err := h.svc.GetAppList(c)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}

	// Convert
	resList := make([]*pay.PayAppResp, 0, len(list))
	for _, app := range list {
		resList = append(resList, &pay.PayAppResp{
			ID:                app.ID,
			AppKey:            app.AppKey,
			Name:              app.Name,
			Status:            app.Status,
			Remark:            app.Remark,
			OrderNotifyURL:    app.OrderNotifyURL,
			RefundNotifyURL:   app.RefundNotifyURL,
			TransferNotifyURL: app.TransferNotifyURL,
			CreateTime:        app.CreateTime,
		})
	}
	response.WriteSuccess(c, resList)
}

// GetAppPage 获得支付应用分页
func (h *PayAppHandler) GetAppPage(c *gin.Context) {
	var r pay.PayAppPageReq
	if err := c.ShouldBindQuery(&r); err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	res, err := h.svc.GetAppPage(c, &r)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, res)
}
