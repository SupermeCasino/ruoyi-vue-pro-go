package product

import "github.com/google/wire"

var ProviderSet = wire.NewSet(
	NewProductBrandHandler,
	NewProductBrowseHistoryHandler,
	NewProductCategoryHandler,
	NewProductCommentHandler,
	NewProductFavoriteHandler,
	NewProductPropertyHandler,
	NewProductSpuHandler,
	NewHandlers,
)

type Handlers struct {
	Brand         *ProductBrandHandler
	BrowseHistory *ProductBrowseHistoryHandler
	Category      *ProductCategoryHandler
	Comment       *ProductCommentHandler
	Favorite      *ProductFavoriteHandler
	Property      *ProductPropertyHandler
	Spu           *ProductSpuHandler
}

func NewHandlers(
	brand *ProductBrandHandler,
	browseHistory *ProductBrowseHistoryHandler,
	category *ProductCategoryHandler,
	comment *ProductCommentHandler,
	favorite *ProductFavoriteHandler,
	property *ProductPropertyHandler,
	spu *ProductSpuHandler,
) *Handlers {
	return &Handlers{
		Brand:         brand,
		BrowseHistory: browseHistory,
		Category:      category,
		Comment:       comment,
		Favorite:      favorite,
		Property:      property,
		Spu:           spu,
	}
}
