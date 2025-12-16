package pay

import (
	"github.com/wxlbd/ruoyi-mall-go/internal/api/req"
	"github.com/wxlbd/ruoyi-mall-go/internal/api/resp"
	"github.com/wxlbd/ruoyi-mall-go/internal/model/pay"
	paySvc "github.com/wxlbd/ruoyi-mall-go/internal/service/pay"
	"github.com/wxlbd/ruoyi-mall-go/pkg/errors"
	"github.com/wxlbd/ruoyi-mall-go/pkg/pagination"
	"github.com/wxlbd/ruoyi-mall-go/pkg/response"
	"github.com/wxlbd/ruoyi-mall-go/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
)

type PayRefundHandler struct {
	svc    *paySvc.PayRefundService
	appSvc *paySvc.PayAppService
}

func NewPayRefundHandler(svc *paySvc.PayRefundService, appSvc *paySvc.PayAppService) *PayRefundHandler {
	return &PayRefundHandler{
		svc:    svc,
		appSvc: appSvc,
	}
}

// GetRefund 获得退款订单
func (h *PayRefundHandler) GetRefund(c *gin.Context) {
	id := utils.ParseInt64(c.Query("id"))
	refund, err := h.svc.GetRefund(c, id)
	if err != nil {
		c.Error(err)
		return
	}
	if refund == nil {
		c.JSON(200, response.Success(&resp.PayRefundDetailsResp{}))
		return
	}

	app, _ := h.appSvc.GetApp(c, refund.AppID)

	r := convertRefundDetailsResp(refund, app)
	c.JSON(200, response.Success(r))
}

// GetRefundPage 获得退款订单分页
func (h *PayRefundHandler) GetRefundPage(c *gin.Context) {
	var r req.PayRefundPageReq
	if err := c.ShouldBindQuery(&r); err != nil {
		c.JSON(200, errors.ErrParam)
		return
	}
	pageResult, err := h.svc.GetRefundPage(c, &r)
	if err != nil {
		c.Error(err)
		return
	}

	// Enrich App Info
	var appIds []int64
	for _, item := range pageResult.List {
		appIds = append(appIds, item.AppID)
	}
	appMap, _ := h.appSvc.GetAppMap(c, appIds)

	list := make([]*resp.PayRefundResp, 0, len(pageResult.List))
	for _, item := range pageResult.List {
		list = append(list, convertRefundResp(item, appMap[item.AppID]))
	}

	c.JSON(200, response.Success(pagination.PageResult[*resp.PayRefundResp]{
		List:  list,
		Total: pageResult.Total,
	}))
}

// Helpers

func convertRefundResp(refund *pay.PayRefund, app *pay.PayApp) *resp.PayRefundResp {
	r := &resp.PayRefundResp{}
	copier.Copy(r, refund)
	if app != nil {
		r.AppName = app.Name
	}
	return r
}

func convertRefundDetailsResp(refund *pay.PayRefund, app *pay.PayApp) *resp.PayRefundDetailsResp {
	r := &resp.PayRefundDetailsResp{}
	copier.Copy(&r.PayRefundResp, refund)
	if app != nil {
		r.AppName = app.Name
		r.App = &resp.PayAppResp{}
		copier.Copy(r.App, app)
	}
	return r
}
