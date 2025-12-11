package pay

import (
	"backend-go/internal/api/req"
	"backend-go/internal/api/resp"
	"backend-go/internal/model/pay"
	"backend-go/internal/pkg/core"
	paySvc "backend-go/internal/service/pay"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
)

type PayNotifyHandler struct {
	svc    *paySvc.PayNotifyService
	appSvc *paySvc.PayAppService
}

func NewPayNotifyHandler(svc *paySvc.PayNotifyService, appSvc *paySvc.PayAppService) *PayNotifyHandler {
	return &PayNotifyHandler{
		svc:    svc,
		appSvc: appSvc,
	}
}

// GetNotifyTaskDetail 获得回调通知详情 (Task + Logs)
func (h *PayNotifyHandler) GetNotifyTaskDetail(c *gin.Context) {
	id := core.ParseInt64(c.Query("id"))
	task, err := h.svc.GetNotifyTask(c, id)
	if err != nil {
		c.Error(err)
		return
	}
	if task == nil {
		c.JSON(200, core.Success(&resp.PayNotifyTaskDetailResp{}))
		return
	}

	logs, _ := h.svc.GetNotifyLogList(c, id)
	app, _ := h.appSvc.GetApp(c, task.AppID)

	r := &resp.PayNotifyTaskDetailResp{}
	copier.Copy(&r.PayNotifyTaskResp, task)
	if app != nil {
		r.AppName = app.Name
	}

	logResps := make([]*resp.PayNotifyLogResp, 0, len(logs))
	for _, log := range logs {
		lr := &resp.PayNotifyLogResp{}
		copier.Copy(lr, log)
		logResps = append(logResps, lr)
	}
	r.Logs = logResps

	c.JSON(200, core.Success(r))
}

// GetNotifyTaskPage 获得回调通知分页
func (h *PayNotifyHandler) GetNotifyTaskPage(c *gin.Context) {
	var r req.PayNotifyTaskPageReq
	if err := c.ShouldBindQuery(&r); err != nil {
		c.JSON(200, core.ErrParam)
		return
	}
	pageResult, err := h.svc.GetNotifyTaskPage(c, &r)
	if err != nil {
		c.Error(err)
		return
	}

	var appIds []int64
	for _, item := range pageResult.List {
		appIds = append(appIds, item.AppID)
	}
	appMap, _ := h.appSvc.GetAppMap(c, appIds)

	list := make([]*resp.PayNotifyTaskResp, 0, len(pageResult.List))
	for _, item := range pageResult.List {
		tr := &resp.PayNotifyTaskResp{}
		copier.Copy(tr, item)
		if app, ok := appMap[item.AppID]; ok {
			tr.AppName = app.Name
		}
		list = append(list, tr)
	}

	c.JSON(200, core.Success(core.PageResult[*resp.PayNotifyTaskResp]{
		List:  list,
		Total: pageResult.Total,
	}))
}

// Helpers
func convertNotifyTaskResp(task *pay.PayNotifyTask, app *pay.PayApp) *resp.PayNotifyTaskResp {
	r := &resp.PayNotifyTaskResp{}
	copier.Copy(r, task)
	if app != nil {
		r.AppName = app.Name
	}
	return r
}
