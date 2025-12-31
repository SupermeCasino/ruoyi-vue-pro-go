package app

import (
	"github.com/wxlbd/ruoyi-mall-go/internal/api/handler/app/mall"
	"github.com/wxlbd/ruoyi-mall-go/internal/api/handler/app/member"
	"github.com/wxlbd/ruoyi-mall-go/internal/api/handler/app/pay"
)

type AppHandlers struct {
	Mall   *mall.Handlers
	Member *member.Handlers
	Pay    *pay.Handlers
}
