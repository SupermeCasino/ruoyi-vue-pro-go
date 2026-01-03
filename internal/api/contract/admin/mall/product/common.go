package product

// ProductSkuPropertyResp SKU 属性响应
type ProductSkuPropertyResp struct {
	PropertyID   int64  `json:"propertyId"`
	PropertyName string `json:"propertyName"`
	ValueID      int64  `json:"valueId"`
	ValueName    string `json:"valueName"`
}
