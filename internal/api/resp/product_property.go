package resp

import "time"

// ProductPropertyResp 属性项 Response
type ProductPropertyResp struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Remark    string    `json:"remark"`
	CreateTime time.Time `json:"createTime"`
}

// ProductPropertyValueResp 属性值 Response
type ProductPropertyValueResp struct {
	ID         int64     `json:"id"`
	PropertyID int64     `json:"propertyId"`
	Name       string    `json:"name"`
	Remark     string    `json:"remark"`
	CreateTime  time.Time `json:"createTime"`
}
