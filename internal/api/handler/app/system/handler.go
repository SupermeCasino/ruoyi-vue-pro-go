package system

import "github.com/google/wire"

var ProviderSet = wire.NewSet(
	NewAppTenantHandler,
	NewHandlers,
)

type Handlers struct {
	Tenant *AppTenantHandler
}

func NewHandlers(
	tenant *AppTenantHandler,
) *Handlers {
	return &Handlers{
		Tenant: tenant,
	}
}
