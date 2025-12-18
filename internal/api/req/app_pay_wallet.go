package req

type AppPayWalletRechargeCreateReq struct {
	PayPrice  *int   `json:"payPrice"`
	PackageID *int64 `json:"packageId"`
}
