package resp

import "time"

// ProductFavoriteResp (Admin)
type ProductFavoriteResp struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"userId"`
	SpuID     int64     `json:"spuId"`
	SpuName   string    `json:"spuName"` // Extra info often needed for admin list
	PicURL    string    `json:"picUrl"`  // SPU Pic
	Price     int64     `json:"price"`   // SPU Price
	CreatedAt time.Time `json:"createdAt"`
}

// AppFavoriteResp (App)
type AppFavoriteResp struct {
	ID        int64     `json:"id"`
	SpuID     int64     `json:"spuId"`
	SpuName   string    `json:"spuName"`
	PicURL    string    `json:"picUrl"`
	Price     int64     `json:"price"`
	CreatedAt time.Time `json:"createdAt"`
}
