package resp

import "time"

// ProductBrandResp 品牌 Response
type ProductBrandResp struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	PicURL      string    `json:"picUrl"`
	Sort        int       `json:"sort"`
	Description string    `json:"description"`
	Status      int       `json:"status"`
	CreateTime   time.Time `json:"createTime"`
}

// ProductBrandSimpleResp 精简品牌 Response
type ProductBrandSimpleResp struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}
