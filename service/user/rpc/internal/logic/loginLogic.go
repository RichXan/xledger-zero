package logic

import (
	"context"
	"errors"
	"fmt"
	"time"

	"xledger/service/user/rpc/internal/svc"
	"xledger/service/user/rpc/user"

	"github.com/golang-jwt/jwt/v4"
	"github.com/zeromicro/go-zero/core/logx"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type LoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LoginLogic) Login(in *user.LoginRequest) (*user.LoginResponse, error) {
	// 1. 验证输入
	if in.Email == "" || in.Password == "" {
		return &user.LoginResponse{
			Code:    400,
			Message: "Email and password are required",
		}, nil
	}

	// 验证邮箱格式
	if !ValidateEmail(in.Email) {
		return &user.LoginResponse{
			Code:    400,
			Message: "Invalid email format",
		}, nil
	}

	// 2. 查找用户
	existingUser, err := l.svcCtx.UserModel.FindByEmail(l.ctx, in.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &user.LoginResponse{
				Code:    401,
				Message: "Invalid email or password",
			}, nil
		}
		l.Logger.Errorf("Failed to find user: %v", err)
		return &user.LoginResponse{
			Code:    500,
			Message: "Internal server error",
		}, nil
	}

	// 3. 验证密码
	err = bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(in.Password))
	if err != nil {
		return &user.LoginResponse{
			Code:    401,
			Message: "Invalid email or password",
		}, nil
	}

	// 4. 检查用户状态
	if existingUser.Status != "active" {
		return &user.LoginResponse{
			Code:    403,
			Message: fmt.Sprintf("User account is %s", existingUser.Status),
		}, nil
	}

	// 5. 生成 JWT Token
	now := time.Now()
	// 从配置读取过期时间
	accessTokenExpires := l.svcCtx.Config.JwtAccessExpire
	refreshTokenExpires := l.svcCtx.Config.JwtRefreshExpire
	if accessTokenExpires == 0 {
		accessTokenExpires = 7200 // 默认 2 hours
	}
	if refreshTokenExpires == 0 {
		refreshTokenExpires = 604800 // 默认 7 days
	}

	accessToken, err := l.generateToken(existingUser.ID.String(), existingUser.Email, now.Unix(), accessTokenExpires)
	if err != nil {
		l.Logger.Errorf("Failed to generate access token: %v", err)
		return &user.LoginResponse{
			Code:    500,
			Message: "Failed to generate token",
		}, nil
	}

	refreshToken, err := l.generateToken(existingUser.ID.String(), existingUser.Email, now.Unix(), refreshTokenExpires)
	if err != nil {
		l.Logger.Errorf("Failed to generate refresh token: %v", err)
		return &user.LoginResponse{
			Code:    500,
			Message: "Failed to generate refresh token",
		}, nil
	}

	// 6. 返回结果
	return &user.LoginResponse{
		Code:    200,
		Message: "Login successful",
		Data: &user.LoginData{
			User: &user.User{
				Id:        existingUser.ID.String(),
				Username:  existingUser.Username,
				Email:     existingUser.Email,
				Status:    existingUser.Status,
				CreatedAt: existingUser.CreatedAt.Format("2006-01-02 15:04:05"),
				UpdatedAt: existingUser.UpdatedAt.Format("2006-01-02 15:04:05"),
			},
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
			ExpiresIn:    accessTokenExpires,
			TokenType:    "Bearer",
		},
	}, nil
}

func (l *LoginLogic) generateToken(userID, email string, iat, exp int64) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"iat":     iat,
		"exp":     iat + exp,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// 从配置读取密钥
	secret := l.svcCtx.Config.JwtSecret
	if secret == "" {
		secret = "xledger-secret-key-change-in-production"
	}
	return token.SignedString([]byte(secret))
}
