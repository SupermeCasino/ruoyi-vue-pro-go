package admin

import (
	"github.com/wxlbd/ruoyi-mall-go/internal/api/handler/admin/infra"
	"github.com/wxlbd/ruoyi-mall-go/internal/api/handler/admin/mall"
	"github.com/wxlbd/ruoyi-mall-go/internal/api/handler/admin/member"
	"github.com/wxlbd/ruoyi-mall-go/internal/api/handler/admin/pay"
	"github.com/wxlbd/ruoyi-mall-go/internal/api/handler/admin/statistics"
	"github.com/wxlbd/ruoyi-mall-go/internal/api/handler/admin/system"
)

type AdminHandlers struct {
	Infra      *infra.Handlers
	Mall       *mall.Handlers
	Member     *member.Handlers
	Pay        *pay.Handlers
	Statistics *statistics.Handlers
	System     *system.Handlers
}
