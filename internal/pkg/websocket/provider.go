package websocket

import "github.com/google/wire"

// ProviderSet is the Wire provider set for websocket package
var ProviderSet = wire.NewSet(NewManager)
