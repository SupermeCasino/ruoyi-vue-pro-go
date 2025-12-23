package product

import (
	"strconv"

	"github.com/wxlbd/ruoyi-mall-go/internal/api/req"
	"github.com/wxlbd/ruoyi-mall-go/internal/api/resp"
	memberSvc "github.com/wxlbd/ruoyi-mall-go/internal/service/member"
	"github.com/wxlbd/ruoyi-mall-go/internal/service/product"
	"github.com/wxlbd/ruoyi-mall-go/pkg/context"
	"github.com/wxlbd/ruoyi-mall-go/pkg/errors"
	"github.com/wxlbd/ruoyi-mall-go/pkg/pagination"
	"github.com/wxlbd/ruoyi-mall-go/pkg/response"
	"github.com/wxlbd/ruoyi-mall-go/pkg/utils"

	"github.com/gin-gonic/gin"
)

type AppProductSpuHandler struct {
	spuSvc         *product.ProductSpuService
	historySvc     *product.ProductBrowseHistoryService
	memberUserSvc  *memberSvc.MemberUserService
	memberLevelSvc *memberSvc.MemberLevelService
}

func NewAppProductSpuHandler(spuSvc *product.ProductSpuService, historySvc *product.ProductBrowseHistoryService, memberUserSvc *memberSvc.MemberUserService, memberLevelSvc *memberSvc.MemberLevelService) *AppProductSpuHandler {
	return &AppProductSpuHandler{
		spuSvc:         spuSvc,
		historySvc:     historySvc,
		memberUserSvc:  memberUserSvc,
		memberLevelSvc: memberLevelSvc,
	}
}

// GetSpuDetail 获得 SPU 详情 (Trigger History)
// @Summary 获得 SPU 详情
// @Tags 用户 APP - 商品 SPU
// @Produce json
// @Param id query int true "SPU 编号"
// @Success 200 {object} resp.ProductSpuResp
// @Router /app-api/product/spu/get-detail [get]
func (h *AppProductSpuHandler) GetSpuDetail(c *gin.Context) {
	idStr := c.Query("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(200, errors.ErrParam)
		return
	}

	res, err := h.spuSvc.GetSpuDetail(c, id)
	if err != nil {
		c.Error(err)
		return
	}
	if res == nil {
		response.WriteBizError(c, errors.NewBizError(1006000002, "商品不存在"))
		return
	}
	if res.Status != 1 {
		response.WriteBizError(c, errors.NewBizError(1006000003, "商品已下架"))
		return
	}

	userID := context.GetLoginUserID(c)
	if userID > 0 {
		_ = h.historySvc.CreateBrowseHistory(c, userID, id)
		_ = h.spuSvc.UpdateBrowseCount(c, id, 1)
	}

	res.SalesCount += res.VirtualSalesCount

	discountPercent := 100
	if userID > 0 {
		user, _ := h.memberUserSvc.GetUser(c, userID)
		if user != nil && user.LevelID > 0 {
			level, _ := h.memberLevelSvc.GetLevel(c, user.LevelID)
			if level != nil {
				discountPercent = level.DiscountPercent
			}
		}
	}

	for _, sku := range res.Skus {
		if discountPercent < 100 {
			sku.VipPrice = int(int64(sku.Price) * int64(discountPercent) / 100)
		} else {
			sku.VipPrice = sku.Price
		}
	}

	response.WriteSuccess(c, res)
}

// GetSpuList 获得商品 SPU 列表
func (h *AppProductSpuHandler) GetSpuList(c *gin.Context) {
	idsStr := c.Query("ids")
	ids := utils.SplitToInt64(idsStr)
	if len(ids) == 0 {
		response.WriteSuccess(c, []resp.AppProductSpuResp{})
		return
	}

	list, err := h.spuSvc.GetSpuList(c, ids)
	if err != nil {
		response.WriteError(c, 500, err.Error())
		return
	}

	// 转换并计算 VIP 价格 (如果登录)
	resList := make([]resp.AppProductSpuResp, len(list))
	userID := context.GetLoginUserID(c)
	discountPercent := 100
	if userID > 0 {
		user, _ := h.memberUserSvc.GetUser(c, userID)
		if user != nil && user.LevelID > 0 {
			level, _ := h.memberLevelSvc.GetLevel(c, user.LevelID)
			if level != nil {
				discountPercent = level.DiscountPercent
			}
		}
	}

	for i, spu := range list {
		price := spu.Price
		vipPrice := spu.Price
		if discountPercent < 100 {
			vipPrice = int(int64(price) * int64(discountPercent) / 100)
		}

		resList[i] = resp.AppProductSpuResp{
			ID:          spu.ID,
			Name:        spu.Name,
			PicURL:      spu.PicURL,
			Price:       spu.Price,
			MarketPrice: spu.MarketPrice,
			SalesCount:  spu.SalesCount + spu.VirtualSalesCount,
			VIPPrice:    vipPrice,
			Stock:       spu.Stock,
			Status:      spu.Status,
			CreateTime:   spu.CreateTime,
		}
	}
	response.WriteSuccess(c, resList)
}

// GetSpuPage 获得商品 SPU 分页
func (h *AppProductSpuHandler) GetSpuPage(c *gin.Context) {
	var r req.AppProductSpuPageReq
	if err := c.ShouldBindQuery(&r); err != nil {
		response.WriteError(c, 400, err.Error())
		return
	}

	// 调用 Service
	pageResult, err := h.spuSvc.GetSpuPageForApp(c, &r)
	if err != nil {
		response.WriteError(c, 500, err.Error())
		return
	}

	// 转换
	userID := context.GetLoginUserID(c)
	discountPercent := 100
	if userID > 0 {
		user, _ := h.memberUserSvc.GetUser(c, userID)
		if user != nil && user.LevelID > 0 {
			level, _ := h.memberLevelSvc.GetLevel(c, user.LevelID)
			if level != nil {
				discountPercent = level.DiscountPercent
			}
		}
	}

	list := make([]resp.AppProductSpuResp, len(pageResult.List))
	for i, spu := range pageResult.List {
		price := spu.Price
		vipPrice := spu.Price
		if discountPercent < 100 {
			vipPrice = int(int64(price) * int64(discountPercent) / 100)
		}

		list[i] = resp.AppProductSpuResp{
			ID:          spu.ID,
			Name:        spu.Name,
			PicURL:      spu.PicURL,
			Price:       spu.Price,
			MarketPrice: spu.MarketPrice,
			SalesCount:  spu.SalesCount + spu.VirtualSalesCount,
			VIPPrice:    vipPrice,
			Stock:       spu.Stock,
			Status:      spu.Status,
			CreateTime:   spu.CreateTime,
		}
	}

	response.WriteSuccess(c, pagination.PageResult[resp.AppProductSpuResp]{
		List:  list,
		Total: pageResult.Total,
	})
}
