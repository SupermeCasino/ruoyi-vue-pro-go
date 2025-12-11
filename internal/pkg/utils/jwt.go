package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("yudao-backend-go-secret") // TODO: Move to config

type Claims struct {
	UserID   int64  `json:"userId"`
	UserType int    `json:"userType"` // 0: Member, 1: Admin
	TenantID int64  `json:"tenantId"`
	Nickname string `json:"nickname"`
	jwt.RegisteredClaims
}

// GenerateToken 生成 JWT Token
func GenerateToken(userID int64, duration time.Duration) (string, error) {
	return GenerateTokenWithInfo(userID, 0, 0, "", duration)
}

// GenerateTokenWithInfo 生成包含完整信息的 JWT Token
func GenerateTokenWithInfo(userID int64, userType int, tenantID int64, nickname string, duration time.Duration) (string, error) {
	claims := Claims{
		UserID:   userID,
		UserType: userType,
		TenantID: tenantID,
		Nickname: nickname,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "yudao-go",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// ParseToken 解析 Token
func ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
