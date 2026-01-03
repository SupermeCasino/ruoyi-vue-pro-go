package product

import (
	"sort"

	product2 "github.com/wxlbd/ruoyi-mall-go/internal/api/contract/app/mall/product"
	"github.com/wxlbd/ruoyi-mall-go/internal/service/mall/product"
	"github.com/wxlbd/ruoyi-mall-go/pkg/response"
	"github.com/wxlbd/ruoyi-mall-go/pkg/utils"

	"github.com/gin-gonic/gin"
)

type AppCategoryHandler struct {
	svc *product.ProductCategoryService
}

func NewAppCategoryHandler(svc *product.ProductCategoryService) *AppCategoryHandler {
	return &AppCategoryHandler{svc: svc}
}

// GetCategoryList 获得商品分类列表
func (h *AppCategoryHandler) GetCategoryList(c *gin.Context) {
	list, err := h.svc.GetEnableCategoryList(c)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	// sort
	sort.Slice(list, func(i, j int) bool {
		return list[i].Sort < list[j].Sort
	})

	// convert
	resList := make([]product2.AppCategoryResp, len(list))
	for i, v := range list {
		resList[i] = product2.AppCategoryResp{
			ID:       v.ID,
			ParentID: v.ParentID,
			Name:     v.Name,
			PicURL:   v.PicURL,
		}
	}
	response.WriteSuccess(c, resList)
}

// GetCategoryListByIds 获得商品分类列表，指定编号
func (h *AppCategoryHandler) GetCategoryListByIds(c *gin.Context) {
	idsStr := c.Query("ids")
	ids := utils.SplitToInt64(idsStr)
	if len(ids) == 0 {
		response.WriteSuccess(c, []product2.AppCategoryResp{})
		return
	}

	list, err := h.svc.GetEnableCategoryListByIds(c, ids)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	// sort
	sort.Slice(list, func(i, j int) bool {
		return list[i].Sort < list[j].Sort
	})

	// convert
	resList := make([]product2.AppCategoryResp, len(list))
	for i, v := range list {
		resList[i] = product2.AppCategoryResp{
			ID:       v.ID,
			ParentID: v.ParentID,
			Name:     v.Name,
			PicURL:   v.PicURL,
		}
	}
	response.WriteSuccess(c, resList)
}
