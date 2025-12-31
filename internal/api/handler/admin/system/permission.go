package system

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/wxlbd/ruoyi-mall-go/internal/service/system"

	"github.com/wxlbd/ruoyi-mall-go/internal/api/req"
	"github.com/wxlbd/ruoyi-mall-go/pkg/errors"
	"github.com/wxlbd/ruoyi-mall-go/pkg/response"
)

type PermissionHandler struct {
	svc       *system.PermissionService
	tenantSvc *system.TenantService
}

func NewPermissionHandler(svc *system.PermissionService, tenantSvc *system.TenantService) *PermissionHandler {
	return &PermissionHandler{
		svc:       svc,
		tenantSvc: tenantSvc,
	}
}

func (h *PermissionHandler) GetRoleMenuList(c *gin.Context) {
	roleIdStr := c.Query("roleId")
	roleId, _ := strconv.ParseInt(roleIdStr, 10, 64)
	list, err := h.svc.GetRoleMenuListByRoleId(c.Request.Context(), []int64{roleId})
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, list)
}

func (h *PermissionHandler) AssignRoleMenu(c *gin.Context) {
	var r req.PermissionAssignRoleMenuReq
	if err := c.ShouldBindJSON(&r); err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	// Filter menus by tenant
	err := h.tenantSvc.HandleTenantMenu(c, func(allowedMenuIds []int64) {
		if allowedMenuIds == nil {
			return
		}
		// Filter r.MenuIDs
		allowedSet := make(map[int64]bool)
		for _, id := range allowedMenuIds {
			allowedSet[id] = true
		}
		filtered := make([]int64, 0, len(r.MenuIDs))
		for _, id := range r.MenuIDs {
			if allowedSet[id] {
				filtered = append(filtered, id)
			}
		}
		r.MenuIDs = filtered
	})
	if err != nil {
		response.WriteBizError(c, err)
		return
	}

	if err := h.svc.AssignRoleMenu(c.Request.Context(), r.RoleID, r.MenuIDs); err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, true)
}

func (h *PermissionHandler) AssignRoleDataScope(c *gin.Context) {
	var r req.PermissionAssignRoleDataScopeReq
	if err := c.ShouldBindJSON(&r); err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	if err := h.svc.AssignRoleDataScope(c.Request.Context(), r.RoleID, r.DataScope, r.DataScopeDeptIDs); err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, true)
}

func (h *PermissionHandler) GetUserRoleList(c *gin.Context) {
	userIdStr := c.Query("userId")
	userId, _ := strconv.ParseInt(userIdStr, 10, 64)
	list, err := h.svc.GetUserRoleIdListByUserId(c.Request.Context(), userId)
	if err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, list)
}

func (h *PermissionHandler) AssignUserRole(c *gin.Context) {
	var r req.PermissionAssignUserRoleReq
	if err := c.ShouldBindJSON(&r); err != nil {
		response.WriteBizError(c, errors.ErrParam)
		return
	}
	if err := h.svc.AssignUserRole(c.Request.Context(), r.UserID, r.RoleIDs); err != nil {
		response.WriteBizError(c, err)
		return
	}
	response.WriteSuccess(c, true)
}
