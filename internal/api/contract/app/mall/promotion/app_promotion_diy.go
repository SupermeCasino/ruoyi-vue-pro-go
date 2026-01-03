package promotion

import "gorm.io/datatypes"

// AppDiyTemplatePropertyResp DIY 模板属性 Response (App)
type AppDiyTemplatePropertyResp struct {
	ID       int64          `json:"id"`
	Name     string         `json:"name"`
	Property datatypes.JSON `json:"property"`
	Home     datatypes.JSON `json:"home"`
	User     datatypes.JSON `json:"user"`
}

// AppDiyPagePropertyResp DIY 页面属性 Response (App)
type AppDiyPagePropertyResp struct {
	ID       int64          `json:"id"`
	Name     string         `json:"name"`
	Property datatypes.JSON `json:"property"`
}
