package admin

import (
	"backend-go/internal/api/req"
	"backend-go/internal/api/resp"
	"backend-go/internal/pkg/core"
	"backend-go/internal/service"

	"github.com/gin-gonic/gin"
)

// MemberStatisticsHandler 会员统计处理器
type MemberStatisticsHandler struct {
	memberStatisticsService       service.MemberStatisticsService
	tradeOrderStatisticsService   service.TradeOrderStatisticsService
	apiAccessLogStatisticsService service.ApiAccessLogStatisticsService
}

// NewMemberStatisticsHandler 创建会员统计处理器
func NewMemberStatisticsHandler(
	memberStatisticsService service.MemberStatisticsService,
	tradeOrderStatisticsService service.TradeOrderStatisticsService,
	apiAccessLogStatisticsService service.ApiAccessLogStatisticsService,
) *MemberStatisticsHandler {
	return &MemberStatisticsHandler{
		memberStatisticsService:       memberStatisticsService,
		tradeOrderStatisticsService:   tradeOrderStatisticsService,
		apiAccessLogStatisticsService: apiAccessLogStatisticsService,
	}
}

// GetMemberSummary 获得会员统计摘要
// GET /statistics/member/summary
func (h *MemberStatisticsHandler) GetMemberSummary(c *gin.Context) {
	result, err := h.memberStatisticsService.GetMemberSummary(c)
	if err != nil {
		core.WriteError(c, core.ServerErrCode, err.Error())
		return
	}

	core.WriteSuccess(c, result)
}

// GetMemberAnalyse 获得会员分析数据
// GET /statistics/member/analyse
func (h *MemberStatisticsHandler) GetMemberAnalyse(c *gin.Context) {
	var reqVO req.MemberAnalyseReqVO
	if err := c.ShouldBindQuery(&reqVO); err != nil {
		core.WriteError(c, core.ParamErrCode, err.Error())
		return
	}

	if len(reqVO.Times) != 2 {
		core.WriteError(c, core.ParamErrCode, "时间范围参数错误")
		return
	}

	beginTime := reqVO.Times[0]
	endTime := reqVO.Times[1]

	// 1.1 查询分析对照数据
	comparisonData, err := h.memberStatisticsService.GetMemberAnalyseComparisonData(c, beginTime, endTime)
	if err != nil {
		core.WriteError(c, core.ServerErrCode, err.Error())
		return
	}

	// 1.2 查询成交用户数量
	payUserCount, err := h.tradeOrderStatisticsService.GetPayUserCount(c, beginTime, endTime)
	if err != nil {
		core.WriteError(c, core.ServerErrCode, err.Error())
		return
	}

	// 1.3 计算客单价
	atv := int64(0)
	if payUserCount > 0 {
		payPrice, err := h.tradeOrderStatisticsService.GetOrderPayPrice(c, beginTime, endTime)
		if err != nil {
			core.WriteError(c, core.ServerErrCode, err.Error())
			return
		}
		atv = payPrice / payUserCount
	}

	// 1.4 查询访客数量
	visitUserCount, err := h.apiAccessLogStatisticsService.GetIpCount(c, 0, beginTime, endTime)
	if err != nil {
		core.WriteError(c, core.ServerErrCode, err.Error())
		return
	}

	// 1.5 下单用户数量
	orderUserCount, err := h.tradeOrderStatisticsService.GetOrderUserCount(c, beginTime, endTime)
	if err != nil {
		core.WriteError(c, core.ServerErrCode, err.Error())
		return
	}

	// 2. 拼接返回
	result := &resp.MemberAnalyseRespVO{
		VisitUserCount: visitUserCount,
		OrderUserCount: orderUserCount,
		PayUserCount:   payUserCount,
		ATV:            atv,
		ComparisonData: *comparisonData,
	}

	core.WriteSuccess(c, result)
}

// GetMemberAreaStatisticsList 按照省份获得会员统计列表
// GET /statistics/member/area-statistics-list
func (h *MemberStatisticsHandler) GetMemberAreaStatisticsList(c *gin.Context) {
	result, err := h.memberStatisticsService.GetMemberAreaStatisticsList(c)
	if err != nil {
		core.WriteError(c, core.ServerErrCode, err.Error())
		return
	}

	core.WriteSuccess(c, result)
}

// GetMemberSexStatisticsList 按照性别获得会员统计列表
// GET /statistics/member/sex-statistics-list
func (h *MemberStatisticsHandler) GetMemberSexStatisticsList(c *gin.Context) {
	result, err := h.memberStatisticsService.GetMemberSexStatisticsList(c)
	if err != nil {
		core.WriteError(c, core.ServerErrCode, err.Error())
		return
	}

	core.WriteSuccess(c, result)
}

// GetMemberTerminalStatisticsList 按照终端获得会员统计列表
// GET /statistics/member/terminal-statistics-list
func (h *MemberStatisticsHandler) GetMemberTerminalStatisticsList(c *gin.Context) {
	result, err := h.memberStatisticsService.GetMemberTerminalStatisticsList(c)
	if err != nil {
		core.WriteError(c, core.ServerErrCode, err.Error())
		return
	}

	core.WriteSuccess(c, result)
}

// GetUserCountComparison 获得用户数量对比
// GET /statistics/member/user-count-comparison
func (h *MemberStatisticsHandler) GetUserCountComparison(c *gin.Context) {
	result, err := h.memberStatisticsService.GetUserCountComparison(c)
	if err != nil {
		core.WriteError(c, core.ServerErrCode, err.Error())
		return
	}

	core.WriteSuccess(c, result)
}

// GetMemberRegisterCountList 获得会员注册数量列表
// GET /statistics/member/register-count-list
func (h *MemberStatisticsHandler) GetMemberRegisterCountList(c *gin.Context) {
	var reqVO req.MemberAnalyseReqVO
	if err := c.ShouldBindQuery(&reqVO); err != nil {
		core.WriteError(c, core.ParamErrCode, err.Error())
		return
	}

	if len(reqVO.Times) != 2 {
		core.WriteError(c, core.ParamErrCode, "时间范围参数错误")
		return
	}

	result, err := h.memberStatisticsService.GetMemberRegisterCountList(c, reqVO.Times[0], reqVO.Times[1])
	if err != nil {
		core.WriteError(c, core.ServerErrCode, err.Error())
		return
	}

	core.WriteSuccess(c, result)
}
