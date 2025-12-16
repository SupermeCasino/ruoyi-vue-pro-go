package middleware

import (
	"fmt"
	"strings"

	"github.com/wxlbd/ruoyi-mall-go/pkg/cache"
	"github.com/wxlbd/ruoyi-mall-go/pkg/context"
	"github.com/wxlbd/ruoyi-mall-go/pkg/response"
	"github.com/wxlbd/ruoyi-mall-go/pkg/utils"

	"github.com/gin-gonic/gin"
)

const (
	// Redis Key 前缀：访问令牌，与 Java 保持一致
	RedisKeyAccessToken = "oauth2_access_token:%s"
)

// OAuth2AccessToken 访问令牌结构（用于从 Redis 解析）
type OAuth2AccessToken struct {
	AccessToken  string            `json:"accessToken"`
	RefreshToken string            `json:"refreshToken"`
	UserID       int64             `json:"userId"`
	UserType     int               `json:"userType"`
	TenantID     int64             `json:"tenantId"`
	UserInfo     map[string]string `json:"userInfo"`
}

// Auth Middleware for JWT authentication
// 使用 JWT + Redis 白名单双重验证机制
func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := obtainAuthorization(c)
		if token == "" {
			c.AbortWithStatusJSON(401, response.Error(401, "未登录"))
			return
		}

		// 1. 先验证 JWT 格式和签名
		claims, err := utils.ParseToken(token)
		if err != nil {
			c.AbortWithStatusJSON(401, response.Error(401, "Token无效"))
			return
		}

		// 2. 再检查 Redis 白名单（如果 Redis 可用）
		if cache.RDB != nil {
			redisKey := fmt.Sprintf(RedisKeyAccessToken, token)
			exists, err := cache.RDB.Exists(c.Request.Context(), redisKey).Result()
			if err == nil && exists == 0 {
				// Token 不在白名单中（已登出）
				c.AbortWithStatusJSON(401, response.Error(401, "Token已失效，请重新登录"))
				return
			}
		}

		// 3. 从 JWT Claims 构建 LoginUser（JWT 中已包含完整信息）
		loginUser := &context.LoginUser{
			UserID:   claims.UserID,
			UserType: claims.UserType,
			TenantID: claims.TenantID,
			Nickname: claims.Nickname,
		}

		// 4. Set LoginUser to Context
		context.SetLoginUser(c, loginUser)
		c.Next()
	}
}

// obtainAuthorization 从请求头或参数中获取 Authorization Token
// 支持 Header 和 Parameter 两种方式，与 Java SecurityFrameworkUtils.obtainAuthorization 对齐
func obtainAuthorization(c *gin.Context) string {
	// 1. 先从 Authorization Header 获取
	token := c.GetHeader("Authorization")
	if token != "" {
		// Remove Bearer prefix if present
		if len(token) > 7 && strings.ToUpper(token[0:7]) == "BEARER " {
			token = token[7:]
		}
		return token
	}

	// 2. 再从 Query Parameter 获取（备选方案）
	token = c.Query("Authorization")
	if token != "" {
		// Remove Bearer prefix if present
		if len(token) > 7 && strings.ToUpper(token[0:7]) == "BEARER " {
			token = token[7:]
		}
		return token
	}

	// 3. 最后从 Form Parameter 获取
	token = c.PostForm("Authorization")
	if token != "" {
		// Remove Bearer prefix if present
		if len(token) > 7 && strings.ToUpper(token[0:7]) == "BEARER " {
			token = token[7:]
		}
		return token
	}

	return ""
}
