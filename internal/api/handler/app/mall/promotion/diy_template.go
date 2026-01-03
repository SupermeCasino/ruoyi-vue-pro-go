package promotion

import (
	"strconv"

	"github.com/gin-gonic/gin"
	appPromotionContract "github.com/wxlbd/ruoyi-mall-go/internal/api/contract/app/mall/promotion"
	promotionModel "github.com/wxlbd/ruoyi-mall-go/internal/model/promotion"
	"github.com/wxlbd/ruoyi-mall-go/internal/service/mall/promotion"
	"github.com/wxlbd/ruoyi-mall-go/pkg/errors"
	"github.com/wxlbd/ruoyi-mall-go/pkg/response"
	"gorm.io/datatypes"
)

type AppDiyTemplateHandler struct {
	templateSvc promotion.DiyTemplateService
	pageSvc     promotion.DiyPageService
}

func NewAppDiyTemplateHandler(templateSvc promotion.DiyTemplateService, pageSvc promotion.DiyPageService) *AppDiyTemplateHandler {
	return &AppDiyTemplateHandler{templateSvc: templateSvc, pageSvc: pageSvc}
}

// GetUsedDiyTemplate 使用中的装修模板 (对齐 Java: AppDiyTemplateController.getUsedDiyTemplate)
func (h *AppDiyTemplateHandler) GetUsedDiyTemplate(c *gin.Context) {
	diyTemplate, err := h.templateSvc.GetUsedDiyTemplate(c)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, h.buildVo(c, diyTemplate))
}

// GetDiyTemplate 获得装修模板 (对齐 Java: AppDiyTemplateController.getDiyTemplate)
func (h *AppDiyTemplateHandler) GetDiyTemplate(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Query("id"), 10, 64)
	if id == 0 {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	diyTemplate, err := h.templateSvc.GetDiyTemplateModel(c, id)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, h.buildVo(c, diyTemplate))
}

// buildVo 构建响应 (对齐 Java: AppDiyTemplateController.buildVo)
func (h *AppDiyTemplateHandler) buildVo(c *gin.Context, diyTemplate *promotionModel.PromotionDiyTemplate) *appPromotionContract.AppDiyTemplatePropertyResp {
	if diyTemplate == nil {
		return nil
	}

	// 查询模板下的页面
	pages, err := h.pageSvc.GetDiyPageByTemplateId(c, diyTemplate.ID)
	if err != nil {
		pages = []*promotionModel.PromotionDiyPage{}
	}

	// 查找首页和我的页面 (对齐 Java: DiyPageEnum.INDEX/MY)
	var home, user datatypes.JSON
	for _, page := range pages {
		switch page.Name {
		case "首页": // DiyPageEnum.INDEX
			home = page.Property
		case "我的": // DiyPageEnum.MY
			user = page.Property
		}
	}

	// 拼接返回 (对齐 Java: DiyTemplateConvert.convertPropertyVo2)
	return &appPromotionContract.AppDiyTemplatePropertyResp{
		ID:       diyTemplate.ID,
		Name:     diyTemplate.Name,
		Property: diyTemplate.Property,
		Home:     home,
		User:     user,
	}
}
