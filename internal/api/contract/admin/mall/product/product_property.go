package product

import "time"

// ProductPropertyCreateReq 创建属性项 Request
type ProductPropertyCreateReq struct {
	Name   string `json:"name" binding:"required"`
	Remark string `json:"remark"`
}

// ProductPropertyUpdateReq 更新属性项 Request
type ProductPropertyUpdateReq struct {
	ID     int64  `json:"id" binding:"required"`
	Name   string `json:"name" binding:"required"`
	Remark string `json:"remark"`
}

// ProductPropertyPageReq 分页查询属性项 Request
type ProductPropertyPageReq struct {
	PageNo   int    `form:"pageNo" binding:"required,min=1"`
	PageSize int    `form:"pageSize" binding:"required,min=1,max=100"`
	Name     string `form:"name"`
}

// ProductPropertyListReq 列表查询属性项 Request
type ProductPropertyListReq struct {
	Name string `form:"name"`
}

// ProductPropertyValueCreateReq 创建属性值 Request
type ProductPropertyValueCreateReq struct {
	PropertyID int64  `json:"propertyId" binding:"required"`
	Name       string `json:"name" binding:"required"`
	Remark     string `json:"remark"`
}

// ProductPropertyValueUpdateReq 更新属性值 Request
type ProductPropertyValueUpdateReq struct {
	ID         int64  `json:"id" binding:"required"`
	PropertyID int64  `json:"propertyId" binding:"required"`
	Name       string `json:"name" binding:"required"`
	Remark     string `json:"remark"`
}

// ProductPropertyValuePageReq 分页查询属性值 Request
type ProductPropertyValuePageReq struct {
	PageNo     int    `form:"pageNo" binding:"required,min=1"`
	PageSize   int    `form:"pageSize" binding:"required,min=1,max=100"`
	PropertyID int64  `form:"propertyId"`
	Name       string `form:"name"`
}

// ProductPropertyResp 属性项响应
type ProductPropertyResp struct {
	ID         int64     `json:"id"`
	Name       string    `json:"name"`
	Remark     string    `json:"remark"`
	CreateTime time.Time `json:"createTime"`
}

// ProductPropertyValueResp 属性值响应
type ProductPropertyValueResp struct {
	ID         int64     `json:"id"`
	PropertyID int64     `json:"propertyId"`
	Name       string    `json:"name"`
	Remark     string    `json:"remark"`
	CreateTime time.Time `json:"createTime"`
}
