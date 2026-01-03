package product

// AppProductSpuDetailResp 用户 APP - 商品 SPU 详情 Response VO - 对齐Java版本AppProductSpuDetailRespVO
type AppProductSpuDetailResp struct {
	// ========== 基本信息 =========
	ID            int64    `json:"id"`            // 商品SPU编号
	Name          string   `json:"name"`          // 商品名称
	Introduction  string   `json:"introduction"`  // 商品简介
	Description   string   `json:"description"`   // 商品详情
	CategoryID    int64    `json:"categoryId"`    // 商品分类编号
	PicURL        string   `json:"picUrl"`        // 商品封面图
	SliderPicURLs []string `json:"sliderPicUrls"` // 商品轮播图

	// ========== SKU 相关字段 =========
	SpecType    bool `json:"specType"`    // 规格类型
	Price       int  `json:"price"`       // 商品价格，单位：分
	MarketPrice int  `json:"marketPrice"` // 市场价，单位：分
	Stock       int  `json:"stock"`       // 库存

	// ========== 统计相关字段 =========
	SalesCount int `json:"salesCount"` // 商品销量

	// ========== SKU 数组 =========
	Skus []AppProductSpuDetailSkuResp `json:"skus"` // SKU数组
}

// AppProductSpuDetailSkuResp 用户 APP - 商品 SPU 详情的 SKU 信息 - 对齐Java版本AppProductSpuDetailRespVO.Sku
type AppProductSpuDetailSkuResp struct {
	ID          int64                           `json:"id"`          // 商品SKU编号
	Properties  []AppProductPropertyValueDetail `json:"properties"`  // 商品属性数组
	Price       int                             `json:"price"`       // 销售价格，单位：分
	MarketPrice int                             `json:"marketPrice"` // 市场价，单位：分
	VipPrice    int                             `json:"vipPrice"`    // VIP价格，单位：分
	PicURL      string                          `json:"picUrl"`      // 图片地址
	Stock       int                             `json:"stock"`       // 库存
	Weight      float64                         `json:"weight"`      // 商品重量，单位：kg
	Volume      float64                         `json:"volume"`      // 商品体积，单位：m³
}

// AppProductPropertyValueDetail 商品属性值详情 - 对齐Java版本AppProductPropertyValueDetailRespVO
type AppProductPropertyValueDetail struct {
	PropertyID   int64  `json:"propertyId"`   // 属性编号
	PropertyName string `json:"propertyName"` // 属性名称
	ValueID      int64  `json:"valueId"`      // 属性值编号
	ValueName    string `json:"valueName"`    // 属性值名称
}
