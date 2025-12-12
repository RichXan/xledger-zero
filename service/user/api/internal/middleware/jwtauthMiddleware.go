// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v4"
	"github.com/redis/go-redis/v9"
)

type JwtAuthMiddleware struct {
	redisClient *redis.Client
	jwtSecret   string
}

func NewJwtAuthMiddleware() *JwtAuthMiddleware {
	// TODO: 从配置注入 Redis 和 JWT Secret
	return &JwtAuthMiddleware{}
}

func (m *JwtAuthMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 提取 Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing authorization header", http.StatusUnauthorized)
			return
		}

		// 移除 "Bearer " 前缀
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			http.Error(w, "Invalid authorization format", http.StatusUnauthorized)
			return
		}

		// 检查 token 是否在黑名单中
		if m.redisClient != nil {
			blacklistKey := fmt.Sprintf("blacklist:token:%s", tokenString)
			exists, err := m.redisClient.Exists(r.Context(), blacklistKey).Result()
			if err == nil && exists > 0 {
				http.Error(w, "Token has been revoked", http.StatusUnauthorized)
				return
			}
		}

		// 解析 JWT token 提取用户 ID
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// 验证签名方法
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			// 使用配置的密钥
			secret := m.jwtSecret
			if secret == "" {
				secret = "xledger-secret-key-change-in-production"
			}
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// 提取用户 ID 并设置到 context
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			if userId, ok := claims["user_id"].(string); ok {
				ctx := context.WithValue(r.Context(), "userId", userId)
				r = r.WithContext(ctx)
			}
		}

		// 放行到下一个 handler
		next(w, r)
	}
}
