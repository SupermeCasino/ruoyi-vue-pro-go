package handler

import (
	"backend-go/internal/api/resp"
	"backend-go/internal/pkg/area"
	"backend-go/internal/pkg/core"

	"github.com/gin-gonic/gin"
)

// AreaHandler 地区处理器
type AreaHandler struct{}

// NewAreaHandler 创建地区处理器
func NewAreaHandler() *AreaHandler {
	return &AreaHandler{}
}

// GetAreaTree 获得地区树
// GET /admin-api/system/area/tree
func (h *AreaHandler) GetAreaTree(c *gin.Context) {
	tree := area.GetAreaTree()
	if tree == nil {
		c.JSON(200, core.Success([]*resp.AreaNodeResp{}))
		return
	}

	result := convertAreaTree(tree)
	c.JSON(200, core.Success(result))
}

// GetAreaByIP 获得 IP 对应的地区名
// GET /admin-api/system/area/get-by-ip?ip=xxx
func (h *AreaHandler) GetAreaByIP(c *gin.Context) {
	ip := c.Query("ip")
	if ip == "" {
		c.JSON(200, core.ErrParam)
		return
	}

	// TODO: 集成 ip2region 库实现 IP 查询
	// 当前返回未知
	c.JSON(200, core.Success("未知"))
}

// convertAreaTree 转换地区树为响应结构
func convertAreaTree(areas []*area.Area) []*resp.AreaNodeResp {
	if areas == nil {
		return nil
	}

	result := make([]*resp.AreaNodeResp, 0, len(areas))
	for _, a := range areas {
		node := &resp.AreaNodeResp{
			ID:   a.ID,
			Name: a.Name,
		}
		if len(a.Children) > 0 {
			node.Children = convertAreaTree(a.Children)
		}
		result = append(result, node)
	}
	return result
}
