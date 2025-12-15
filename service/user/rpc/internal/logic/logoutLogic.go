package logic

import (
	"context"
	"fmt"
	"time"

	"xledger/service/user/rpc/internal/svc"
	"xledger/service/user/rpc/user"

	"github.com/golang-jwt/jwt/v4"
	"github.com/zeromicro/go-zero/core/logx"
)

type LogoutLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLogoutLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LogoutLogic {
	return &LogoutLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LogoutLogic) Logout(in *user.LogoutRequest) (*user.LogoutResponse, error) {
	// 验证输入
	if in.AccessToken == "" {
		return &user.LogoutResponse{
			Code:    400,
			Message: "Access token is required",
		}, nil
	}

	// 解析 JWT Token 获取过期时间
	token, err := jwt.Parse(in.AccessToken, func(token *jwt.Token) (interface{}, error) {
		// 从配置读取密钥
		secret := l.svcCtx.Config.Auth.AccessSecret
		if secret == "" {
			secret = "xledger-secret-key-change-in-production"
		}
		return []byte(secret), nil
	})

	// 即使解析失败，也尝试将 token 加入黑名单（防止重放攻击）
	var ttl time.Duration
	if err == nil && token.Valid {
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			if exp, ok := claims["exp"].(float64); ok {
				expirationTime := time.Unix(int64(exp), 0)
				ttl = time.Until(expirationTime)
				if ttl < 0 {
					// Token 已过期
					return &user.LogoutResponse{
						Code:    200,
						Message: "Logout successful",
					}, nil
				}
			}
		}
	}

	// 如果无法获取过期时间，使用默认值
	if ttl <= 0 {
		ttl = 2 * time.Hour // 使用默认的 access token 过期时间
	}

	// 将 token 加入 Redis 黑名单
	if l.svcCtx.RedisClient != nil {
		key := fmt.Sprintf("blacklist:token:%s", in.AccessToken)
		err = l.svcCtx.RedisClient.Set(l.ctx, key, "1", ttl).Err()
		if err != nil {
			l.Logger.Errorf("Failed to add token to blacklist: %v", err)
			return &user.LogoutResponse{
				Code:    500,
				Message: "Failed to logout",
			}, nil
		}
	}

	return &user.LogoutResponse{
		Code:    200,
		Message: "Logout successful",
	}, nil
}
