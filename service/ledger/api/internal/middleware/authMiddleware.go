// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v4"
	"github.com/redis/go-redis/v9"
)

type AuthMiddleware struct {
	redisClient *redis.Client
	jwtSecret   string
}

func NewAuthMiddleware(redisClient *redis.Client, jwtSecret string) *AuthMiddleware {
	return &AuthMiddleware{
		redisClient: redisClient,
		jwtSecret:   jwtSecret,
	}
}

func (m *AuthMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 从 Authorization header 获取 token
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

		// 验证 token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(m.jwtSecret), nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// 检查 token 是否在黑名单中
		ctx := context.Background()
		exists, err := m.redisClient.Exists(ctx, "blacklist:"+tokenString).Result()
		if err == nil && exists > 0 {
			http.Error(w, "Token has been revoked", http.StatusUnauthorized)
			return
		}

		// 从 token 中提取 user_id 并设置到 context
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			if userID, ok := claims["user_id"]; ok {
				ctx := context.WithValue(r.Context(), "user_id", userID)
				r = r.WithContext(ctx)
			}
		}

		next(w, r)
	}
}
