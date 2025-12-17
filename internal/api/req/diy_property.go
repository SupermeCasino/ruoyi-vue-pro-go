package req

import "gorm.io/datatypes"

// DiyPagePropertyUpdateReq 装修页面属性更新请求
// Java: DiyPagePropertyUpdateRequestVO
type DiyPagePropertyUpdateReq struct {
	ID       int64          `json:"id" binding:"required"`
	Property datatypes.JSON `json:"property" binding:"required"`
}

// DiyTemplatePropertyUpdateReq 装修模板属性更新请求
// Java: DiyTemplatePropertyUpdateRequestVO
type DiyTemplatePropertyUpdateReq struct {
	ID       int64          `json:"id" binding:"required"`
	Property datatypes.JSON `json:"property" binding:"required"`
}
