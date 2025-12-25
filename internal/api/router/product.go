package router

import (
	productHandler "github.com/wxlbd/ruoyi-mall-go/internal/api/handler/admin/product"
	"github.com/wxlbd/ruoyi-mall-go/internal/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterProductRoutes 注册商品管理模块路由
func RegisterProductRoutes(engine *gin.Engine,
	productCategoryHandler *productHandler.ProductCategoryHandler,
	productBrandHandler *productHandler.ProductBrandHandler,
	productPropertyHandler *productHandler.ProductPropertyHandler,
	productSpuHandler *productHandler.ProductSpuHandler,
	productCommentHandler *productHandler.ProductCommentHandler,
	productFavoriteHandler *productHandler.ProductFavoriteHandler,
	productBrowseHistoryHandler *productHandler.ProductBrowseHistoryHandler,
	casbinMiddleware *middleware.CasbinMiddleware,
) {
	productGroup := engine.Group("/admin-api/product")
	productGroup.Use(middleware.Auth())
	{
		// Category Routes
		categoryGroup := productGroup.Group("/category")
		{
			categoryGroup.POST("/create", casbinMiddleware.RequirePermission("product:category:create"), productCategoryHandler.CreateCategory)
			categoryGroup.PUT("/update", casbinMiddleware.RequirePermission("product:category:update"), productCategoryHandler.UpdateCategory)
			categoryGroup.DELETE("/delete", casbinMiddleware.RequirePermission("product:category:delete"), productCategoryHandler.DeleteCategory)
			categoryGroup.GET("/get", casbinMiddleware.RequirePermission("product:category:query"), productCategoryHandler.GetCategory)
			categoryGroup.GET("/list", casbinMiddleware.RequirePermission("product:category:query"), productCategoryHandler.GetCategoryList)
		}

		// Property Routes
		propertyGroup := productGroup.Group("/property")
		{
			propertyGroup.POST("/create", casbinMiddleware.RequirePermission("product:property:create"), productPropertyHandler.CreateProperty)
			propertyGroup.PUT("/update", casbinMiddleware.RequirePermission("product:property:update"), productPropertyHandler.UpdateProperty)
			propertyGroup.DELETE("/delete", casbinMiddleware.RequirePermission("product:property:delete"), productPropertyHandler.DeleteProperty)
			propertyGroup.GET("/get", casbinMiddleware.RequirePermission("product:property:query"), productPropertyHandler.GetProperty)
			propertyGroup.GET("/page", casbinMiddleware.RequirePermission("product:property:query"), productPropertyHandler.GetPropertyPage)
			propertyGroup.GET("/simple-list", productPropertyHandler.GetPropertySimpleList)

			// Property Value Routes
			valueGroup := propertyGroup.Group("/value")
			{
				valueGroup.POST("/create", casbinMiddleware.RequirePermission("product:property:create"), productPropertyHandler.CreatePropertyValue)
				valueGroup.PUT("/update", casbinMiddleware.RequirePermission("product:property:update"), productPropertyHandler.UpdatePropertyValue)
				valueGroup.DELETE("/delete", casbinMiddleware.RequirePermission("product:property:delete"), productPropertyHandler.DeletePropertyValue)
				valueGroup.GET("/get", casbinMiddleware.RequirePermission("product:property:query"), productPropertyHandler.GetPropertyValue)
				valueGroup.GET("/page", casbinMiddleware.RequirePermission("product:property:query"), productPropertyHandler.GetPropertyValuePage)
				valueGroup.GET("/simple-list", productPropertyHandler.GetPropertyValueSimpleList)
			}
		}

		// Brand Routes
		brandGroup := productGroup.Group("/brand")
		{
			brandGroup.POST("/create", casbinMiddleware.RequirePermission("product:brand:create"), productBrandHandler.CreateBrand)
			brandGroup.PUT("/update", casbinMiddleware.RequirePermission("product:brand:update"), productBrandHandler.UpdateBrand)
			brandGroup.DELETE("/delete", casbinMiddleware.RequirePermission("product:brand:delete"), productBrandHandler.DeleteBrand)
			brandGroup.GET("/get", casbinMiddleware.RequirePermission("product:brand:query"), productBrandHandler.GetBrand)
			brandGroup.GET("/page", casbinMiddleware.RequirePermission("product:brand:query"), productBrandHandler.GetBrandPage)
			brandGroup.GET("/list", productBrandHandler.GetBrandList)
			brandGroup.GET("/list-all-simple", productBrandHandler.GetBrandList)
		}

		// SPU Routes
		spuGroup := productGroup.Group("/spu")
		spuGroup.Use(middleware.ProductErrorHandler()) // 使用商品模块错误处理中间件
		{
			spuGroup.POST("/create", casbinMiddleware.RequirePermission("product:spu:create"), productSpuHandler.CreateSpu)
			spuGroup.PUT("/update", casbinMiddleware.RequirePermission("product:spu:update"), productSpuHandler.UpdateSpu)
			spuGroup.PUT("/update-status", casbinMiddleware.RequirePermission("product:spu:update"), productSpuHandler.UpdateSpuStatus)
			spuGroup.DELETE("/delete", casbinMiddleware.RequirePermission("product:spu:delete"), productSpuHandler.DeleteSpu)
			spuGroup.GET("/get-detail", casbinMiddleware.RequirePermission("product:spu:query"), productSpuHandler.GetSpuDetail)
			spuGroup.GET("/page", casbinMiddleware.RequirePermission("product:spu:query"), productSpuHandler.GetSpuPage)
			spuGroup.GET("/get-count", casbinMiddleware.RequirePermission("product:spu:query"), productSpuHandler.GetTabsCount)
			spuGroup.GET("/list-all-simple", productSpuHandler.GetSpuSimpleList)
			spuGroup.GET("/list", productSpuHandler.GetSpuList)
			spuGroup.GET("/export", casbinMiddleware.RequirePermission("product:spu:export"), productSpuHandler.ExportSpuList)
		}

		// Comment Routes
		commentGroup := productGroup.Group("/comment")
		{
			commentGroup.GET("/page", casbinMiddleware.RequirePermission("product:comment:query"), productCommentHandler.GetCommentPage)
			commentGroup.PUT("/update-visible", casbinMiddleware.RequirePermission("product:comment:update"), productCommentHandler.UpdateCommentVisible)
			commentGroup.PUT("/reply", casbinMiddleware.RequirePermission("product:comment:update"), productCommentHandler.ReplyComment)
			commentGroup.POST("/create", casbinMiddleware.RequirePermission("product:comment:create"), productCommentHandler.CreateComment)
		}

		// Favorite Routes (Admin)
		favoriteGroup := productGroup.Group("/favorite")
		{
			favoriteGroup.GET("/page", casbinMiddleware.RequirePermission("product:favorite:query"), productFavoriteHandler.GetFavoritePage)
		}

		// Browse History Routes (Admin)
		browseHistoryGroup := productGroup.Group("/browse-history")
		{
			browseHistoryGroup.GET("/page", casbinMiddleware.RequirePermission("product:browse-history:query"), productBrowseHistoryHandler.GetBrowseHistoryPage)
		}
	}
}
