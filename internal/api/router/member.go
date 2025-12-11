package router

import (
	memberAdmin "backend-go/internal/api/handler/admin/member"
	"backend-go/internal/middleware"

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
		configGroup.PUT("/save", memberConfigHandler.SaveConfig)
		configGroup.GET("/get", memberConfigHandler.GetConfig)
	}

	// Member Group 用户分组
	groupGroup := api.Group("/member/group")
	{
		groupGroup.POST("/create", memberGroupHandler.CreateGroup)
		groupGroup.PUT("/update", memberGroupHandler.UpdateGroup)
		groupGroup.DELETE("/delete", memberGroupHandler.DeleteGroup)
		groupGroup.GET("/get", memberGroupHandler.GetGroup)
		groupGroup.GET("/list-all-simple", memberGroupHandler.GetSimpleGroupList)
		groupGroup.GET("/page", memberGroupHandler.GetGroupPage)
	}

	// Member Level 会员等级
	levelGroup := api.Group("/member/level")
	{
		levelGroup.POST("/create", memberLevelHandler.CreateLevel)
		levelGroup.PUT("/update", memberLevelHandler.UpdateLevel)
		levelGroup.DELETE("/delete", memberLevelHandler.DeleteLevel)
		levelGroup.GET("/get", memberLevelHandler.GetLevel)
		levelGroup.GET("/list-all-simple", memberLevelHandler.GetLevelListSimple)
		levelGroup.GET("/list", memberLevelHandler.GetLevelListSimple)
	}

	// Member Tag 会员标签
	tagGroup := api.Group("/member/tag")
	{
		tagGroup.POST("/create", memberTagHandler.CreateTag)
		tagGroup.PUT("/update", memberTagHandler.UpdateTag)
		tagGroup.DELETE("/delete", memberTagHandler.DeleteTag)
		tagGroup.GET("/get", memberTagHandler.GetTag)
		tagGroup.GET("/list-all-simple", memberTagHandler.GetSimpleTagList)
		tagGroup.GET("/page", memberTagHandler.GetTagPage)
	}

	// Member User 会员用户
	userGroup := api.Group("/member/user")
	{
		userGroup.PUT("/update", memberUserHandler.UpdateUser)
		userGroup.PUT("/update-level", memberUserHandler.UpdateUserLevel)
		userGroup.PUT("/update-point", memberUserHandler.UpdateUserPoint)
		userGroup.GET("/get", memberUserHandler.GetUser)
		userGroup.GET("/page", memberUserHandler.GetUserPage)
	}
}
