package pay

import (
	"github.com/wxlbd/ruoyi-mall-go/internal/api/contract/admin/pay"
)

type AppPayOrderSubmitReq struct {
	pay.PayOrderSubmitReq
}

type AppPayOrderSubmitResp struct {
	Status         int    `json:"status"`
	DisplayMode    string `json:"displayMode"`
	DisplayContent string `json:"displayContent"`
}
