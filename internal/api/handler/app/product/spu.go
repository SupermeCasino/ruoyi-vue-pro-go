package product

import (
	"backend-go/internal/pkg/core"
	memberSvc "backend-go/internal/service/member"
	"backend-go/internal/service/product"
	"strconv"

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
		c.JSON(200, core.ErrParam)
		return
	}

	res, err := h.spuSvc.GetSpuDetail(c, id)
	if err != nil {
		c.Error(err)
		return
	}
	if res == nil {
		core.WriteBizError(c, core.NewBizError(1006000002, "商品不存在"))
		return
	}
	if res.Status != 1 {
		core.WriteBizError(c, core.NewBizError(1006000003, "商品已下架"))
		return
	}

	userID := core.GetLoginUserID(c)
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

	core.WriteSuccess(c, res)
}
