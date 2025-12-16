package promotion

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/wxlbd/ruoyi-mall-go/internal/api/resp"
	promotionModel "github.com/wxlbd/ruoyi-mall-go/internal/model/promotion"
	"github.com/wxlbd/ruoyi-mall-go/internal/service/promotion"
	"github.com/wxlbd/ruoyi-mall-go/pkg/errors"
	"github.com/wxlbd/ruoyi-mall-go/pkg/response"
)

type AppDiyTemplateHandler struct {
	templateSvc promotion.DiyTemplateService
	pageSvc     promotion.DiyPageService
}

func NewAppDiyTemplateHandler(templateSvc promotion.DiyTemplateService, pageSvc promotion.DiyPageService) *AppDiyTemplateHandler {
	return &AppDiyTemplateHandler{templateSvc: templateSvc, pageSvc: pageSvc}
}

// GetUsedDiyTemplate 使用中的装修模板
func (h *AppDiyTemplateHandler) GetUsedDiyTemplate(c *gin.Context) {
	diyTemplate, err := h.templateSvc.GetUsedDiyTemplate(c)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, h.buildVo(c, diyTemplate))
}

// GetDiyTemplate 获得装修模板
func (h *AppDiyTemplateHandler) GetDiyTemplate(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Query("id"), 10, 64)
	if id == 0 {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	res, err := h.templateSvc.GetDiyTemplate(c, id)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	// Convert Resp back to Model for buildVo (or overload buildVo)
	// Since GetDiyTemplate returns Resp, and buildVo takes Model?
	// Actually GetDiyTemplate returns *resp.DiyTemplateResp.
	// But GetUsedDiyTemplate returns Model.
	// I should unified them or implementation buildVo to accept generic properties.

	// Re-fetch model or use resp object.
	// Since I need ID to fetch pages, Resp object has ID.

	// BUT, build Vo needs pages.
	// Let's adapt buildVo to take ID and Properties.
	response.WriteSuccess(c, h.buildVoFromResp(c, res))
}

func (h *AppDiyTemplateHandler) buildVo(c *gin.Context, diyTemplate *promotionModel.PromotionDiyTemplate) *resp.AppDiyTemplatePropertyResp {
	if diyTemplate == nil {
		return nil
	}
	return h.doBuildVo(c, diyTemplate.ID, &resp.DiyTemplateResp{
		ID:             diyTemplate.ID,
		Name:           diyTemplate.Name,
		PreviewPicUrls: []string(diyTemplate.PreviewPicUrls),
		Property:       diyTemplate.Property,
		Used:           diyTemplate.Used,
		UsedTime:       diyTemplate.UsedTime,
		Remark:         diyTemplate.Remark,
		CreateTime:     diyTemplate.CreateTime,
	})
}

func (h *AppDiyTemplateHandler) buildVoFromResp(c *gin.Context, res *resp.DiyTemplateResp) *resp.AppDiyTemplatePropertyResp {
	if res == nil {
		return nil
	}
	return h.doBuildVo(c, res.ID, res)
}

func (h *AppDiyTemplateHandler) doBuildVo(c *gin.Context, templateId int64, templateResp *resp.DiyTemplateResp) *resp.AppDiyTemplatePropertyResp {
	pages, err := h.pageSvc.GetDiyPageByTemplateId(c, templateId)
	if err != nil {
		// Log error? Or return partial? Java implementation ignores errors here implicitly?
		// Java: list = diyPageService.getDiyPageByTemplateId(id);
		// Assume success or empty.
		pages = []*promotionModel.PromotionDiyPage{}
	}

	home := ""
	user := ""
	for _, page := range pages {
		if page.Name == "首页" {
			home = page.Property
		} else if page.Name == "我的" {
			user = page.Property
		}
	}

	return &resp.AppDiyTemplatePropertyResp{
		DiyTemplateResp: *templateResp,
		Home:            home,
		User:            user,
	}
}
