package resp

import "time"

// AppAddressResp 收件地址响应
type AppAddressResp struct {
	ID            int64     `json:"id"`
	Name          string    `json:"name"`
	Mobile        string    `json:"mobile"`
	AreaID        int64     `json:"areaId"`
	DetailAddress string    `json:"detailAddress"`
	DefaultStatus bool      `json:"defaultStatus"`
	CreateTime    time.Time `json:"createTime"`
}
