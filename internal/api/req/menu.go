package req

// MenuListReq 菜单列表请求参数
type MenuListReq struct {
	Name   string `json:"name" form:"name"`
	Status *int32 `json:"status" form:"status"`
}

// MenuCreateReq 创建菜单请求参数
type MenuCreateReq struct {
	ID            int64  `json:"id"`
	ParentID      int64  `json:"parentId"`                            // 父菜单 ID
	Name          string `json:"name" binding:"required,max=50"`      // 菜单名称
	Type          int32  `json:"type" binding:"required,oneof=1 2 3"` // 菜单类型
	Sort          int32  `json:"sort"`                                // 显示顺序
	Path          string `json:"path" binding:"max=200"`              // 路由地址
	Icon          string `json:"icon"`                                // 菜单图标
	Component     string `json:"component" binding:"max=255"`         // 组件路径
	ComponentName string `json:"componentName"`                       // 组件名
	Permission    string `json:"permission" binding:"max=100"`        // 权限标识
	Status        int32  `json:"status" binding:"oneof=0 1"`          // 状态
	Visible       *bool  `json:"visible"`                             // 是否可见
	KeepAlive     *bool  `json:"keepAlive"`                           // 是否缓存
	AlwaysShow    *bool  `json:"alwaysShow"`                          // 是否总是显示
}

// MenuUpdateReq 更新菜单请求参数
type MenuUpdateReq struct {
	ID int64 `json:"id" binding:"required"`
	MenuCreateReq
}
