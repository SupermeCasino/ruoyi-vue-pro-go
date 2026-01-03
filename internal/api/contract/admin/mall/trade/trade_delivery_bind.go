package trade

// DeliveryPickUpBindReq 自提门店绑定核销员工 Request
type DeliveryPickUpBindReq struct {
	ID            int64   `json:"id" binding:"required"`
	VerifyUserIds []int64 `json:"verifyUserIds" binding:"required"`
}
