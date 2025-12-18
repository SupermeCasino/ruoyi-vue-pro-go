package resp

import "time"

type MemberLevelResp struct {
	ID              int64     `json:"id"`
	Name            string    `json:"name"`
	Level           int       `json:"level"`
	Experience      int       `json:"experience"`
	DiscountPercent int       `json:"discountPercent"`
	Icon            string    `json:"icon"`
	BackgroundURL   string    `json:"backgroundUrl"`
	CreatedAt       time.Time `json:"createTime"`
}
