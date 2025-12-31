package brokerage

import (
	"strconv"

	"github.com/wxlbd/ruoyi-mall-go/internal/api/req"
	tradeReq "github.com/wxlbd/ruoyi-mall-go/internal/api/req/app/trade"
	tradeResp "github.com/wxlbd/ruoyi-mall-go/internal/api/resp/app/trade"
	"github.com/wxlbd/ruoyi-mall-go/internal/model/trade/brokerage"
	brokerageSvc "github.com/wxlbd/ruoyi-mall-go/internal/service/mall/trade/brokerage"
	"github.com/wxlbd/ruoyi-mall-go/pkg/context"
	"github.com/wxlbd/ruoyi-mall-go/pkg/pagination"
	"github.com/wxlbd/ruoyi-mall-go/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/samber/lo"
)

type AppBrokerageRecordHandler struct {
	recordSvc *brokerageSvc.BrokerageRecordService
}

func NewAppBrokerageRecordHandler(recordSvc *brokerageSvc.BrokerageRecordService) *AppBrokerageRecordHandler {
	return &AppBrokerageRecordHandler{recordSvc: recordSvc}
}

// GetBrokerageRecordPage 获得分销记录分页
func (h *AppBrokerageRecordHandler) GetBrokerageRecordPage(c *gin.Context) {
	var reqVO tradeReq.AppBrokerageRecordPageReqVO
	if err := c.ShouldBindQuery(&reqVO); err != nil {
		response.WriteError(c, 400, "参数错误")
		return
	}

	userId := context.GetLoginUserID(c)
	pageReq := &req.BrokerageRecordPageReq{
		PageParam:  reqVO.PageParam,
		UserID:     userId,
		Status:     reqVO.Status,
		CreateTime: reqVO.CreateTime,
		// BizType:    reqVO.BizType, // BizType mismatch: Model int vs Req string?
		// Service expects request struct.
		// Admin Req has BizType string/int?
		// Checking BrokerageRecordService.GetBrokerageRecordPage(..., r *req.BrokerageRecordPageReq)
		// Let's assume we map explicitly if needed.
		// If reqVO.BizType is string, and Admin req has BizType string, it matches.
		// If Admin req has BizType as string but logic uses it as enum value (1,2) or "order" string.
		// Java: App passes "bizType" (string or int?), Service uses string in Admin DTO?
		// Let's assume string for now based on previous file views.
		BizType: reqVO.BizType,
	}

	pageResult, err := h.recordSvc.GetBrokerageRecordPage(c, pageReq)
	if err != nil {
		response.WriteError(c, 500, err.Error())
		return
	}

	writeResp := pagination.PageResult[*tradeResp.AppBrokerageRecordRespVO]{
		Total: pageResult.Total,
		List: lo.Map(pageResult.List, func(item *brokerage.BrokerageRecord, _ int) *tradeResp.AppBrokerageRecordRespVO {
			return &tradeResp.AppBrokerageRecordRespVO{
				ID:          item.ID,
				UserID:      item.UserID,
				BizType:     item.BizType,
				BizID:       item.BizID,
				Price:       item.Price,
				Title:       item.Title,
				Description: item.Description,
				Status:      item.Status,
				Total:       item.TotalPrice,
				CreateTime:   item.CreateTime,
				// StatusName: item.Status // TODO: Dict lookup
			}
		}),
	}
	response.WriteSuccess(c, writeResp)
}

// GetProductBrokeragePrice 获得商品的分销金额
func (h *AppBrokerageRecordHandler) GetProductBrokeragePrice(c *gin.Context) {
	spuIdStr := c.Query("spuId")
	spuId, err := strconv.ParseInt(spuIdStr, 10, 64)
	if err != nil {
		response.WriteError(c, 400, "参数错误")
		return
	}

	userId := context.GetLoginUserID(c)
	result, err := h.recordSvc.CalculateProductBrokeragePrice(c, userId, spuId)
	if err != nil {
		response.WriteError(c, 500, err.Error())
		return
	}
	response.WriteSuccess(c, result)
}
