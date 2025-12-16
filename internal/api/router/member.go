package router

import (
	memberAdmin "github.com/wxlbd/ruoyi-mall-go/internal/api/handler/admin/member"
	"github.com/wxlbd/ruoyi-mall-go/internal/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterMemberRoutes 注册会员管理模块路由
func RegisterMemberRoutes(engine *gin.Engine,
	memberSignInConfigHandler *memberAdmin.MemberSignInConfigHandler,
	memberSignInRecordHandler *memberAdmin.MemberSignInRecordHandler,
	memberPointRecordHandler *memberAdmin.MemberPointRecordHandler,
	memberConfigHandler *memberAdmin.MemberConfigHandler,
	memberGroupHandler *memberAdmin.MemberGroupHandler,
	memberLevelHandler *memberAdmin.MemberLevelHandler,
	memberTagHandler *memberAdmin.MemberTagHandler,
	memberUserHandler *memberAdmin.MemberUserHandler,
	casbinMiddleware *middleware.CasbinMiddleware,
) {
	api := engine.Group("/admin-api")
	api.Use(middleware.Auth())

	// Member Point Record
	pointRecordGroup := api.Group("/member/point/record")
	{
		pointRecordGroup.GET("/page", memberPointRecordHandler.GetPointRecordPage)
	}

	// Member Sign-in Config
	signInConfigGroup := api.Group("/member/sign-in/config")
	{
		signInConfigGroup.POST("/create", memberSignInConfigHandler.CreateSignInConfig)
		signInConfigGroup.PUT("/update", memberSignInConfigHandler.UpdateSignInConfig)
		signInConfigGroup.DELETE("/delete", memberSignInConfigHandler.DeleteSignInConfig)
		signInConfigGroup.GET("/get", memberSignInConfigHandler.GetSignInConfig)
		signInConfigGroup.GET("/list", memberSignInConfigHandler.GetSignInConfigList)
	}

	// Member Sign-in Record (Admin)
	signInRecordGroup := api.Group("/member/sign-in/record")
	{
		signInRecordGroup.GET("/page", memberSignInRecordHandler.GetSignInRecordPage)
	}

	// Member Config 会员配置
	configGroup := api.Group("/member/config")
	{
		configGroup.PUT("/save", casbinMiddleware.RequirePermission("member:config:save"), memberConfigHandler.SaveConfig)
		configGroup.GET("/get", casbinMiddleware.RequirePermission("member:config:query"), memberConfigHandler.GetConfig)
	}

	// Member Group 用户分组
	groupGroup := api.Group("/member/group")
	{
		groupGroup.POST("/create", casbinMiddleware.RequirePermission("member:group:create"), memberGroupHandler.CreateGroup)
		groupGroup.PUT("/update", casbinMiddleware.RequirePermission("member:group:update"), memberGroupHandler.UpdateGroup)
		groupGroup.DELETE("/delete", casbinMiddleware.RequirePermission("member:group:delete"), memberGroupHandler.DeleteGroup)
		groupGroup.GET("/get", casbinMiddleware.RequirePermission("member:group:query"), memberGroupHandler.GetGroup)
		groupGroup.GET("/list-all-simple", memberGroupHandler.GetSimpleGroupList)
		groupGroup.GET("/page", casbinMiddleware.RequirePermission("member:group:query"), memberGroupHandler.GetGroupPage)
	}

	// Member Level 会员等级
	levelGroup := api.Group("/member/level")
	{
		levelGroup.POST("/create", casbinMiddleware.RequirePermission("member:level:create"), memberLevelHandler.CreateLevel)
		levelGroup.PUT("/update", casbinMiddleware.RequirePermission("member:level:update"), memberLevelHandler.UpdateLevel)
		levelGroup.DELETE("/delete", casbinMiddleware.RequirePermission("member:level:delete"), memberLevelHandler.DeleteLevel)
		levelGroup.GET("/get", casbinMiddleware.RequirePermission("member:level:query"), memberLevelHandler.GetLevel)
		levelGroup.GET("/list-all-simple", memberLevelHandler.GetLevelListSimple)
		levelGroup.GET("/list", memberLevelHandler.GetLevelListSimple)
	}

	// Member Tag 会员标签
	tagGroup := api.Group("/member/tag")
	{
		tagGroup.POST("/create", casbinMiddleware.RequirePermission("member:tag:create"), memberTagHandler.CreateTag)
		tagGroup.PUT("/update", casbinMiddleware.RequirePermission("member:tag:update"), memberTagHandler.UpdateTag)
		tagGroup.DELETE("/delete", casbinMiddleware.RequirePermission("member:tag:delete"), memberTagHandler.DeleteTag)
		tagGroup.GET("/get", casbinMiddleware.RequirePermission("member:tag:query"), memberTagHandler.GetTag)
		tagGroup.GET("/list-all-simple", memberTagHandler.GetSimpleTagList)
		tagGroup.GET("/page", casbinMiddleware.RequirePermission("member:tag:query"), memberTagHandler.GetTagPage)
	}

	// Member User 会员用户
	userGroup := api.Group("/member/user")
	{
		userGroup.PUT("/update", casbinMiddleware.RequirePermission("member:user:update"), memberUserHandler.UpdateUser)
		userGroup.PUT("/update-level", casbinMiddleware.RequirePermission("member:user:update-level"), memberUserHandler.UpdateUserLevel)
		userGroup.PUT("/update-point", casbinMiddleware.RequirePermission("member:user:update-point"), memberUserHandler.UpdateUserPoint)
		userGroup.GET("/get", casbinMiddleware.RequirePermission("member:user:query"), memberUserHandler.GetUser)
		userGroup.GET("/page", casbinMiddleware.RequirePermission("member:user:query"), memberUserHandler.GetUserPage)
	}
}
