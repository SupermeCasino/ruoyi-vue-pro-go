package system

// AreaNodeResp 地区节点响应
// 对应 Java: cn.iocoder.yudao.module.controller.admin.ip.vo.AreaNodeRespVO
type AreaNodeResp struct {
	ID       int             `json:"id"`
	Name     string          `json:"name"`
	Children []*AreaNodeResp `json:"children,omitempty"`
}
