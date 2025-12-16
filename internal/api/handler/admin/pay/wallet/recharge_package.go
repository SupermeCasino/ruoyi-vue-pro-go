package wallet

import (
	"strconv"

	"github.com/wxlbd/ruoyi-mall-go/internal/api/req"
	"github.com/wxlbd/ruoyi-mall-go/internal/api/resp"
	"github.com/wxlbd/ruoyi-mall-go/internal/model/pay"
	payData "github.com/wxlbd/ruoyi-mall-go/internal/service/pay/wallet"
	"github.com/wxlbd/ruoyi-mall-go/pkg/errors"
	"github.com/wxlbd/ruoyi-mall-go/pkg/pagination"
	"github.com/wxlbd/ruoyi-mall-go/pkg/response"

	"github.com/gin-gonic/gin"
)

type PayWalletRechargePackageHandler struct {
	svc *payData.PayWalletRechargePackageService
}

func NewPayWalletRechargePackageHandler(svc *payData.PayWalletRechargePackageService) *PayWalletRechargePackageHandler {
	return &PayWalletRechargePackageHandler{svc: svc}
}

// CreateWalletRechargePackage 创建充值套餐
func (h *PayWalletRechargePackageHandler) CreateWalletRechargePackage(c *gin.Context) {
	var r req.PayWalletRechargePackageCreateReq
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(200, errors.ErrParam)
		return
	}
	id, err := h.svc.CreateWalletRechargePackage(c, &r)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, response.Success(id))
}

// UpdateWalletRechargePackage 更新充值套餐
func (h *PayWalletRechargePackageHandler) UpdateWalletRechargePackage(c *gin.Context) {
	var r req.PayWalletRechargePackageUpdateReq
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(200, errors.ErrParam)
		return
	}
	err := h.svc.UpdateWalletRechargePackage(c, &r)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, response.Success(true))
}

// DeleteWalletRechargePackage 删除充值套餐
func (h *PayWalletRechargePackageHandler) DeleteWalletRechargePackage(c *gin.Context) {
	idStr := c.Query("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(200, errors.ErrParam)
		return
	}
	err = h.svc.DeleteWalletRechargePackage(c, id)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, response.Success(true))
}

// GetWalletRechargePackage 获得充值套餐
func (h *PayWalletRechargePackageHandler) GetWalletRechargePackage(c *gin.Context) {
	idStr := c.Query("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(200, errors.ErrParam)
		return
	}
	pkg, err := h.svc.GetWalletRechargePackage(c, id)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, response.Success(convertPackageResp(pkg)))
}

// GetWalletRechargePackagePage 获得充值套餐分页
func (h *PayWalletRechargePackageHandler) GetWalletRechargePackagePage(c *gin.Context) {
	var r req.PayWalletRechargePackagePageReq
	if err := c.ShouldBindQuery(&r); err != nil {
		c.JSON(200, errors.ErrParam)
		return
	}
	res, err := h.svc.GetWalletRechargePackagePage(c, &r)
	if err != nil {
		c.Error(err)
		return
	}

	// Convert list
	newRes := pagination.NewPageResult(make([]*resp.PayWalletRechargePackageResp, 0, len(res.List)), res.Total)
	for _, item := range res.List {
		newRes.List = append(newRes.List, convertPackageResp(item))
	}
	c.JSON(200, response.Success(newRes))
}

func convertPackageResp(item *pay.PayWalletRechargePackage) *resp.PayWalletRechargePackageResp {
	if item == nil {
		return nil
	}
	return &resp.PayWalletRechargePackageResp{
		ID:         item.ID,
		Name:       item.Name,
		PayPrice:   item.PayPrice,
		BonusPrice: item.BonusPrice,
		Status:     item.Status,
		CreateTime: item.CreatedAt,
	}
}
