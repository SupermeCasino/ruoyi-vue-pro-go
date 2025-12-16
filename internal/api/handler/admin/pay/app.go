package pay

import (
	"github.com/wxlbd/ruoyi-mall-go/internal/api/req"
	"github.com/wxlbd/ruoyi-mall-go/internal/api/resp"
	"github.com/wxlbd/ruoyi-mall-go/internal/pkg/core"
	paySvc "github.com/wxlbd/ruoyi-mall-go/internal/service/pay"
	"strconv"

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
	var r req.PayAppCreateReq
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(200, core.ErrParam)
		return
	}
	id, err := h.svc.CreateApp(c, &r)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, core.Success(id))
}

// UpdateApp 更新支付应用
func (h *PayAppHandler) UpdateApp(c *gin.Context) {
	var r req.PayAppUpdateReq
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(200, core.ErrParam)
		return
	}
	err := h.svc.UpdateApp(c, &r)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, core.Success(true))
}

// UpdateAppStatus 更新支付应用状态
func (h *PayAppHandler) UpdateAppStatus(c *gin.Context) {
	var r req.PayAppUpdateStatusReq
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(200, core.ErrParam)
		return
	}
	err := h.svc.UpdateAppStatus(c, &r)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, core.Success(true))
}

// DeleteApp 删除支付应用
func (h *PayAppHandler) DeleteApp(c *gin.Context) {
	idStr := c.Query("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(200, core.ErrParam)
		return
	}
	err = h.svc.DeleteApp(c, id)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, core.Success(true))
}

// GetApp 获得支付应用
func (h *PayAppHandler) GetApp(c *gin.Context) {
	idStr := c.Query("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(200, core.ErrParam)
		return
	}
	app, err := h.svc.GetApp(c, id)
	if err != nil {
		c.Error(err)
		return
	}

	res := &resp.PayAppResp{
		ID:                app.ID,
		AppKey:            app.AppKey,
		Name:              app.Name,
		Status:            app.Status,
		Remark:            app.Remark,
		OrderNotifyURL:    app.OrderNotifyURL,
		RefundNotifyURL:   app.RefundNotifyURL,
		TransferNotifyURL: app.TransferNotifyURL,
		CreateTime:        app.CreatedAt,
	}
	c.JSON(200, core.Success(res))
}

// GetAppPage 获得支付应用分页
func (h *PayAppHandler) GetAppPage(c *gin.Context) {
	var r req.PayAppPageReq
	if err := c.ShouldBindQuery(&r); err != nil {
		c.JSON(200, core.ErrParam)
		return
	}
	res, err := h.svc.GetAppPage(c, &r)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, core.Success(res))
}

// GetAppList 获得支付应用列表
func (h *PayAppHandler) GetAppList(c *gin.Context) {
	list, err := h.svc.GetAppList(c)
	if err != nil {
		c.Error(err)
		return
	}

	// Convert
	resList := make([]*resp.PayAppResp, 0, len(list))
	for _, app := range list {
		resList = append(resList, &resp.PayAppResp{
			ID:                app.ID,
			AppKey:            app.AppKey,
			Name:              app.Name,
			Status:            app.Status,
			Remark:            app.Remark,
			OrderNotifyURL:    app.OrderNotifyURL,
			RefundNotifyURL:   app.RefundNotifyURL,
			TransferNotifyURL: app.TransferNotifyURL,
			CreateTime:        app.CreatedAt,
		})
	}
	c.JSON(200, core.Success(resList))
}
