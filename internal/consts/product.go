package consts

// ProductSpuStatus 商品 SPU 状态枚举
// 对应 Java: ProductSpuStatusEnum
const (
	// ProductSpuStatusRecycle 回收站
	ProductSpuStatusRecycle = -1
	// ProductSpuStatusDisable 下架
	ProductSpuStatusDisable = 0
	// ProductSpuStatusEnable 上架
	ProductSpuStatusEnable = 1
)

// Product Scope Constants (对齐 Java PromotionProductScopeEnum)
const (
	ProductScopeAll      = 1 // 全部商品 (Java: ALL)
	ProductScopeSpu      = 2 // 指定商品 (Java: SPU)
	ProductScopeCategory = 3 // 指定品类 (Java: CATEGORY)
)

// ProductScopeValues 商品范围值数组
var ProductScopeValues = []int{ProductScopeAll, ProductScopeSpu, ProductScopeCategory}

// IsValidProductScope 验证商品范围是否有效
func IsValidProductScope(scope int) bool {
	for _, v := range ProductScopeValues {
		if v == scope {
			return true
		}
	}
	return false
}

// IsProductScopeAll 判断是否为全部商品范围
func IsProductScopeAll(scope int) bool {
	return scope == ProductScopeAll
}

// IsProductScopeSpu 判断是否为指定商品范围
func IsProductScopeSpu(scope int) bool {
	return scope == ProductScopeSpu
}

// IsProductScopeCategory 判断是否为指定品类范围
func IsProductScopeCategory(scope int) bool {
	return scope == ProductScopeCategory
}

// ProductCommentScore 商品评价评分常量
const (
	// ProductCommentScoreBad 差评 (1-2分)
	ProductCommentScoreBad = 2
	// ProductCommentScoreNormal 中评 (3分)
	ProductCommentScoreNormal = 3
	// ProductCommentScoreGood 好评 (4-5分)
	ProductCommentScoreGood = 4
)
