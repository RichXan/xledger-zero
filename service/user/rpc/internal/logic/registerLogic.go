package logic

import (
	"context"
	"errors"

	"xledger/service/user/model"
	"xledger/service/user/rpc/internal/svc"
	"xledger/service/user/rpc/user"

	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type RegisterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// Basic Auth
func (l *RegisterLogic) Register(in *user.RegisterRequest) (*user.RegisterResponse, error) {
	// 1. 验证输入
	if in.Email == "" || in.Password == "" || in.Username == "" {
		return &user.RegisterResponse{
			Code:    400,
			Message: "Missing required fields: username, email, and password are required",
		}, nil
	}

	// 验证邮箱格式
	if !ValidateEmail(in.Email) {
		return &user.RegisterResponse{
			Code:    400,
			Message: "Invalid email format",
		}, nil
	}

	// 验证密码强度
	if err := ValidatePassword(in.Password); err != nil {
		return &user.RegisterResponse{
			Code:    400,
			Message: err.Error(),
		}, nil
	}

	// 验证用户名
	if err := ValidateUsername(in.Username); err != nil {
		return &user.RegisterResponse{
			Code:    400,
			Message: err.Error(),
		}, nil
	}

	// 2. 检查邮箱是否已存在
	existingUser, err := l.svcCtx.UserModel.FindByEmail(l.ctx, in.Email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		l.Logger.Errorf("Failed to check existing user: %v", err)
		return &user.RegisterResponse{
			Code:    500,
			Message: "Internal server error",
		}, nil
	}
	if existingUser != nil {
		return &user.RegisterResponse{
			Code:    400,
			Message: "Email already exists",
		}, nil
	}

	// 3. 密码加密
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	if err != nil {
		l.Logger.Errorf("Failed to hash password: %v", err)
		return &user.RegisterResponse{
			Code:    500,
			Message: "Internal server error",
		}, nil
	}

	// 4. 创建用户
	newUser := &model.User{
		ID:       uuid.New(),
		Username: in.Username,
		Email:    in.Email,
		Password: string(hashedPassword),
		Status:   "active",
	}

	err = l.svcCtx.UserModel.Create(l.ctx, newUser)
	if err != nil {
		l.Logger.Errorf("Failed to create user: %v", err)
		return &user.RegisterResponse{
			Code:    500,
			Message: "Failed to create user",
		}, nil
	}

	// 5. 返回结果
	return &user.RegisterResponse{
		Code:    200,
		Message: "User registered successfully",
		Data: &user.User{
			Id:        newUser.ID.String(),
			Username:  newUser.Username,
			Email:     newUser.Email,
			Status:    newUser.Status,
			CreatedAt: newUser.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: newUser.UpdatedAt.Format("2006-01-02 15:04:05"),
		},
	}, nil
}
