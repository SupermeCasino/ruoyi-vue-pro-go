package handler

import (
	"backend-go/internal/api/resp"
	"backend-go/internal/pkg/area"
	"backend-go/internal/pkg/core"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/lionsoul2014/ip2region/binding/golang/xdb"
	"go.uber.org/zap"
)

// AreaHandler 地区处理器
type AreaHandler struct {
	searcher *xdb.Searcher
}

var (
	ip2regionOnce     sync.Once
	ip2regionSearcher *xdb.Searcher
)

// NewAreaHandler 创建地区处理器
func NewAreaHandler() *AreaHandler {
	ip2regionOnce.Do(func() {
		dbPath := "configs/ip2region.xdb"

		// 从 xdb 文件加载 header 获取正确的 Version
		header, err := xdb.LoadHeaderFromFile(dbPath)
		if err != nil {
			zap.L().Warn("Failed to load ip2region header", zap.Error(err))
			return
		}

		version, err := xdb.VersionFromHeader(header)
		if err != nil {
			zap.L().Warn("Failed to get version from header", zap.Error(err))
			return
		}

		// 加载完整内容到内存
		cBuff, err := xdb.LoadContentFromFile(dbPath)
		if err != nil {
			zap.L().Warn("Failed to load ip2region.xdb", zap.Error(err))
			return
		}

		// 使用正确的 Version 创建 Searcher
		ip2regionSearcher, err = xdb.NewWithBuffer(version, cBuff)
		if err != nil {
			zap.L().Warn("Failed to create ip2region searcher", zap.Error(err))
			return
		}
		zap.L().Info("ip2region searcher initialized", zap.String("version", version.Name))
	})

	return &AreaHandler{searcher: ip2regionSearcher}
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

	// 如果没有 ip2region 数据库，返回未知
	if h.searcher == nil {
		c.JSON(200, core.Success("未知"))
		return
	}

	// 使用 ip2region 查询
	// 返回格式: 区域ID (如 320100 表示南京市)
	regionStr, err := h.searcher.SearchByStr(ip)
	if err != nil {
		zap.L().Debug("ip2region search failed", zap.String("ip", ip), zap.Error(err))
		c.JSON(200, core.Success("未知"))
		return
	}

	// 将区域ID转换为区域名称
	areaID, err := strconv.Atoi(regionStr)
	if err != nil {
		c.JSON(200, core.Success("未知"))
		return
	}

	// 使用 area.Format 获取格式化的地区名
	formatted := area.Format(areaID)
	if formatted == "" {
		c.JSON(200, core.Success("未知"))
		return
	}

	c.JSON(200, core.Success(formatted))
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
