package req

// MenuListReq 菜单列表请求参数
type MenuListReq struct {
	Name   string `json:"name" form:"name"`
	Status *int32 `json:"status" form:"status"`
}

// MenuCreateReq 创建菜单请求参数
type MenuCreateReq struct {
	ID            int64  `json:"id"`
	ParentID      int64  `json:"parentId" binding:"required"`
	Name          string `json:"name" binding:"required,max=50"`
	Type          int32  `json:"type" binding:"required,oneof=1 2 3"` // 1:目录, 2:菜单, 3:按钮
	Sort          int32  `json:"sort" binding:"required"`
	Path          string `json:"path"`
	Icon          string `json:"icon"`
	Component     string `json:"component"`
	ComponentName string `json:"componentName"` // 组件名，用于 KeepAlive
	Permission    string `json:"permission"`
	Status        int32  `json:"status" binding:"required,oneof=0 1"`
	Visible       bool   `json:"visible"`    // 是否可见
	KeepAlive     bool   `json:"keepAlive"`  // 是否缓存
	AlwaysShow    bool   `json:"alwaysShow"` // 是否总是显示
}

// MenuUpdateReq 更新菜单请求参数
type MenuUpdateReq struct {
	ID int64 `json:"id" binding:"required"`
	MenuCreateReq
}
