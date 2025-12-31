package router

import (
	"github.com/wxlbd/ruoyi-mall-go/internal/api/handler/admin/mall/product"
	"github.com/wxlbd/ruoyi-mall-go/internal/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterProductRoutes 注册商品管理模块路由
func RegisterProductRoutes(engine *gin.Engine,
	handlers *product.Handlers,
	casbinMiddleware *middleware.CasbinMiddleware,
) {
	productGroup := engine.Group("/admin-api/product")
	productGroup.Use(middleware.Auth())
	{
		// Category Routes
		categoryGroup := productGroup.Group("/category")
		{
			categoryGroup.POST("/create", casbinMiddleware.RequirePermission("product:category:create"), handlers.Category.CreateCategory)
			categoryGroup.PUT("/update", casbinMiddleware.RequirePermission("product:category:update"), handlers.Category.UpdateCategory)
			categoryGroup.DELETE("/delete", casbinMiddleware.RequirePermission("product:category:delete"), handlers.Category.DeleteCategory)
			categoryGroup.GET("/get", casbinMiddleware.RequirePermission("product:category:query"), handlers.Category.GetCategory)
			categoryGroup.GET("/list", casbinMiddleware.RequirePermission("product:category:query"), handlers.Category.GetCategoryList)
		}

		// Property Routes
		propertyGroup := productGroup.Group("/property")
		{
			propertyGroup.POST("/create", casbinMiddleware.RequirePermission("product:property:create"), handlers.Property.CreateProperty)
			propertyGroup.PUT("/update", casbinMiddleware.RequirePermission("product:property:update"), handlers.Property.UpdateProperty)
			propertyGroup.DELETE("/delete", casbinMiddleware.RequirePermission("product:property:delete"), handlers.Property.DeleteProperty)
			propertyGroup.GET("/get", casbinMiddleware.RequirePermission("product:property:query"), handlers.Property.GetProperty)
			propertyGroup.GET("/page", casbinMiddleware.RequirePermission("product:property:query"), handlers.Property.GetPropertyPage)
			propertyGroup.GET("/simple-list", handlers.Property.GetPropertySimpleList)

			// Property Value Routes
			valueGroup := propertyGroup.Group("/value")
			{
				valueGroup.POST("/create", casbinMiddleware.RequirePermission("product:property:create"), handlers.Property.CreatePropertyValue)
				valueGroup.PUT("/update", casbinMiddleware.RequirePermission("product:property:update"), handlers.Property.UpdatePropertyValue)
				valueGroup.DELETE("/delete", casbinMiddleware.RequirePermission("product:property:delete"), handlers.Property.DeletePropertyValue)
				valueGroup.GET("/get", casbinMiddleware.RequirePermission("product:property:query"), handlers.Property.GetPropertyValue)
				valueGroup.GET("/page", casbinMiddleware.RequirePermission("product:property:query"), handlers.Property.GetPropertyValuePage)
				valueGroup.GET("/simple-list", handlers.Property.GetPropertyValueSimpleList)
			}
		}

		// Brand Routes
		brandGroup := productGroup.Group("/brand")
		{
			brandGroup.POST("/create", casbinMiddleware.RequirePermission("product:brand:create"), handlers.Brand.CreateBrand)
			brandGroup.PUT("/update", casbinMiddleware.RequirePermission("product:brand:update"), handlers.Brand.UpdateBrand)
			brandGroup.DELETE("/delete", casbinMiddleware.RequirePermission("product:brand:delete"), handlers.Brand.DeleteBrand)
			brandGroup.GET("/get", casbinMiddleware.RequirePermission("product:brand:query"), handlers.Brand.GetBrand)
			brandGroup.GET("/page", casbinMiddleware.RequirePermission("product:brand:query"), handlers.Brand.GetBrandPage)
			brandGroup.GET("/list", handlers.Brand.GetBrandList)
			brandGroup.GET("/list-all-simple", handlers.Brand.GetBrandList)
		}

		// SPU Routes
		spuGroup := productGroup.Group("/spu")
		spuGroup.Use(middleware.ProductErrorHandler()) // 使用商品模块错误处理中间件
		{
			spuGroup.POST("/create", casbinMiddleware.RequirePermission("product:spu:create"), handlers.Spu.CreateSpu)
			spuGroup.PUT("/update", casbinMiddleware.RequirePermission("product:spu:update"), handlers.Spu.UpdateSpu)
			spuGroup.PUT("/update-status", casbinMiddleware.RequirePermission("product:spu:update"), handlers.Spu.UpdateSpuStatus)
			spuGroup.DELETE("/delete", casbinMiddleware.RequirePermission("product:spu:delete"), handlers.Spu.DeleteSpu)
			spuGroup.GET("/get-detail", casbinMiddleware.RequirePermission("product:spu:query"), handlers.Spu.GetSpuDetail)
			spuGroup.GET("/page", casbinMiddleware.RequirePermission("product:spu:query"), handlers.Spu.GetSpuPage)
			spuGroup.GET("/get-count", casbinMiddleware.RequirePermission("product:spu:query"), handlers.Spu.GetTabsCount)
			spuGroup.GET("/list-all-simple", handlers.Spu.GetSpuSimpleList)
			spuGroup.GET("/list", handlers.Spu.GetSpuList)
			spuGroup.GET("/export", casbinMiddleware.RequirePermission("product:spu:export"), handlers.Spu.ExportSpuList)
		}

		// Comment Routes
		commentGroup := productGroup.Group("/comment")
		{
			commentGroup.GET("/page", casbinMiddleware.RequirePermission("product:comment:query"), handlers.Comment.GetCommentPage)
			commentGroup.PUT("/update-visible", casbinMiddleware.RequirePermission("product:comment:update"), handlers.Comment.UpdateCommentVisible)
			commentGroup.PUT("/reply", casbinMiddleware.RequirePermission("product:comment:update"), handlers.Comment.ReplyComment)
			commentGroup.POST("/create", casbinMiddleware.RequirePermission("product:comment:create"), handlers.Comment.CreateComment)
		}

		// Favorite Routes (Admin)
		favoriteGroup := productGroup.Group("/favorite")
		{
			favoriteGroup.GET("/page", casbinMiddleware.RequirePermission("product:favorite:query"), handlers.Favorite.GetFavoritePage)
		}

		// Browse History Routes (Admin)
		browseHistoryGroup := productGroup.Group("/browse-history")
		{
			browseHistoryGroup.GET("/page", casbinMiddleware.RequirePermission("product:browse-history:query"), handlers.BrowseHistory.GetBrowseHistoryPage)
		}
	}
}
