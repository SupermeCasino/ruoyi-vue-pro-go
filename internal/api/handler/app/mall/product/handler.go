package product

import "github.com/google/wire"

var ProviderSet = wire.NewSet(
	NewAppProductBrowseHistoryHandler,
	NewAppCategoryHandler,
	NewAppProductCommentHandler,
	NewAppProductFavoriteHandler,
	NewAppProductSpuHandler,
	NewHandlers,
)

type Handlers struct {
	BrowseHistory *AppProductBrowseHistoryHandler
	Category      *AppCategoryHandler
	Comment       *AppProductCommentHandler
	Favorite      *AppProductFavoriteHandler
	Spu           *AppProductSpuHandler
}

func NewHandlers(
	browseHistory *AppProductBrowseHistoryHandler,
	category *AppCategoryHandler,
	comment *AppProductCommentHandler,
	favorite *AppProductFavoriteHandler,
	spu *AppProductSpuHandler,
) *Handlers {
	return &Handlers{
		BrowseHistory: browseHistory,
		Category:      category,
		Comment:       comment,
		Favorite:      favorite,
		Spu:           spu,
	}
}
