package product

import "github.com/wxlbd/ruoyi-mall-go/pkg/errors"

// 商品模块错误码常量
// 参考 Java: yudao-module-mall/yudao-module-product/src/main/java/cn/iocoder/yudao/module/product/enums/ErrorCodeConstants.java
// product 系统，使用 1-008-000-000 段

var (
	// ========== 商品分类相关 1-008-001-000 ============
	ErrCategoryNotExists           = errors.NewBizError(1008001000, "商品分类不存在")
	ErrCategoryParentNotExists     = errors.NewBizError(1008001001, "父分类不存在")
	ErrCategoryParentNotFirstLevel = errors.NewBizError(1008001002, "父分类不能是二级分类")
	ErrCategoryExistsChildren      = errors.NewBizError(1008001003, "存在子分类，无法删除")
	ErrCategoryDisabled            = errors.NewBizError(1008001004, "商品分类已禁用，无法使用")
	ErrCategoryHaveBindSpu         = errors.NewBizError(1008001005, "类别下存在商品，无法删除")

	// ========== 商品品牌相关编号 1-008-002-000 ==========
	ErrBrandNotExists  = errors.NewBizError(1008002000, "品牌不存在")
	ErrBrandDisabled   = errors.NewBizError(1008002001, "品牌已禁用")
	ErrBrandNameExists = errors.NewBizError(1008002002, "品牌名称已存在")

	// ========== 商品属性项 1-008-003-000 ==========
	ErrPropertyNotExists             = errors.NewBizError(1008003000, "属性项不存在")
	ErrPropertyExists                = errors.NewBizError(1008003001, "属性项的名称已存在")
	ErrPropertyDeleteFailValueExists = errors.NewBizError(1008003002, "属性项下存在属性值，无法删除")

	// ========== 商品属性值 1-008-004-000 ==========
	ErrPropertyValueNotExists = errors.NewBizError(1008004000, "属性值不存在")
	ErrPropertyValueExists    = errors.NewBizError(1008004001, "属性值的名称已存在")

	// ========== 商品 SPU 1-008-005-000 ==========
	ErrSpuNotExists                       = errors.NewBizError(1008005000, "商品 SPU 不存在")
	ErrSpuSaveFailCategoryLevelError      = errors.NewBizError(1008005001, "商品分类不正确，原因：必须使用第二级的商品分类及以下")
	ErrSpuSaveFailCouponTemplateNotExists = errors.NewBizError(1008005002, "商品 SPU 保存失败，原因：优惠劵不存在")
	ErrSpuNotEnable                       = errors.NewBizError(1008005003, "商品 SPU 不处于上架状态")
	ErrSpuNotRecycle                      = errors.NewBizError(1008005004, "商品 SPU 不处于回收站状态")

	// ========== 商品 SKU 1-008-006-000 ==========
	ErrSkuNotExists               = errors.NewBizError(1008006000, "商品 SKU 不存在")
	ErrSkuPropertiesDuplicated    = errors.NewBizError(1008006001, "商品 SKU 的属性组合存在重复")
	ErrSpuAttrNumbersMustBeEquals = errors.NewBizError(1008006002, "一个 SPU 下的每个 SKU，其属性项必须一致")
	ErrSpuSkuNotDuplicate         = errors.NewBizError(1008006003, "一个 SPU 下的每个 SKU，必须不重复")
	ErrSkuStockNotEnough          = errors.NewBizError(1008006004, "商品 SKU 库存不足")

	// ========== 商品 评价 1-008-007-000 ==========
	ErrCommentNotExists   = errors.NewBizError(1008007000, "商品评价不存在")
	ErrCommentOrderExists = errors.NewBizError(1008007001, "订单的商品评价已存在")

	// ========== 商品 收藏 1-008-008-000 ==========
	ErrFavoriteExists    = errors.NewBizError(1008008000, "该商品已经被收藏")
	ErrFavoriteNotExists = errors.NewBizError(1008008001, "商品收藏不存在")
)
