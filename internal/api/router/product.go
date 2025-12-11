package router

import (
	productHandler "backend-go/internal/api/handler/admin/product"
	"backend-go/internal/middleware"

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
) {
	productGroup := engine.Group("/admin-api/product")
	productGroup.Use(middleware.Auth())
	{
		// Category Routes
		categoryGroup := productGroup.Group("/category")
		{
			categoryGroup.POST("/create", productCategoryHandler.CreateCategory)
			categoryGroup.PUT("/update", productCategoryHandler.UpdateCategory)
			categoryGroup.DELETE("/delete", productCategoryHandler.DeleteCategory)
			categoryGroup.GET("/get", productCategoryHandler.GetCategory)
			categoryGroup.GET("/list", productCategoryHandler.GetCategoryList)
		}

		// Property Routes
		propertyGroup := productGroup.Group("/property")
		{
			propertyGroup.POST("/create", productPropertyHandler.CreateProperty)
			propertyGroup.PUT("/update", productPropertyHandler.UpdateProperty)
			propertyGroup.DELETE("/delete", productPropertyHandler.DeleteProperty)
			propertyGroup.GET("/get", productPropertyHandler.GetProperty)
			propertyGroup.GET("/page", productPropertyHandler.GetPropertyPage)
			propertyGroup.GET("/simple-list", productPropertyHandler.GetPropertySimpleList)

			// Property Value Routes
			valueGroup := propertyGroup.Group("/value")
			{
				valueGroup.POST("/create", productPropertyHandler.CreatePropertyValue)
				valueGroup.PUT("/update", productPropertyHandler.UpdatePropertyValue)
				valueGroup.DELETE("/delete", productPropertyHandler.DeletePropertyValue)
				valueGroup.GET("/get", productPropertyHandler.GetPropertyValue)
				valueGroup.GET("/page", productPropertyHandler.GetPropertyValuePage)
				valueGroup.GET("/simple-list", productPropertyHandler.GetPropertyValueSimpleList)
			}
		}

		// Brand Routes
		brandGroup := productGroup.Group("/brand")
		{
			brandGroup.POST("/create", productBrandHandler.CreateBrand)
			brandGroup.PUT("/update", productBrandHandler.UpdateBrand)
			brandGroup.DELETE("/delete", productBrandHandler.DeleteBrand)
			brandGroup.GET("/get", productBrandHandler.GetBrand)
			brandGroup.GET("/page", productBrandHandler.GetBrandPage)
			brandGroup.GET("/list", productBrandHandler.GetBrandList)
			brandGroup.GET("/list-all-simple", productBrandHandler.GetBrandList)
		}

		// SPU Routes
		spuGroup := productGroup.Group("/spu")
		{
			spuGroup.POST("/create", productSpuHandler.CreateSpu)
			spuGroup.PUT("/update", productSpuHandler.UpdateSpu)
			spuGroup.PUT("/update-status", productSpuHandler.UpdateSpuStatus)
			spuGroup.DELETE("/delete", productSpuHandler.DeleteSpu)
			spuGroup.GET("/get-detail", productSpuHandler.GetSpuDetail)
			spuGroup.GET("/page", productSpuHandler.GetSpuPage)
			spuGroup.GET("/get-count", productSpuHandler.GetTabsCount)
			spuGroup.GET("/list-all-simple", productSpuHandler.GetSpuSimpleList)
			spuGroup.GET("/list", productSpuHandler.GetSpuList)
			spuGroup.GET("/export", productSpuHandler.ExportSpuList)
		}

		// Comment Routes
		commentGroup := productGroup.Group("/comment")
		{
			commentGroup.GET("/page", productCommentHandler.GetCommentPage)
			commentGroup.PUT("/update-visible", productCommentHandler.UpdateCommentVisible)
			commentGroup.PUT("/reply", productCommentHandler.ReplyComment)
			commentGroup.POST("/create", productCommentHandler.CreateComment)
		}

		// Favorite Routes (Admin)
		favoriteGroup := productGroup.Group("/favorite")
		{
			favoriteGroup.GET("/page", productFavoriteHandler.GetFavoritePage)
		}

		// Browse History Routes (Admin)
		browseHistoryGroup := productGroup.Group("/browse-history")
		{
			browseHistoryGroup.GET("/page", productBrowseHistoryHandler.GetBrowseHistoryPage)
		}
	}
}
