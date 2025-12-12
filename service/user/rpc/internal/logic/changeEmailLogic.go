package logic

import (
	"context"
	"errors"

	"xledger/service/user/rpc/internal/svc"
	"xledger/service/user/rpc/user"

	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type ChangeEmailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewChangeEmailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChangeEmailLogic {
	return &ChangeEmailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ChangeEmailLogic) ChangeEmail(in *user.ChangeEmailRequest) (*user.ChangeEmailResponse, error) {
	// 1. 验证输入
	if in.Id == "" || in.NewEmail == "" || in.Password == "" {
		return &user.ChangeEmailResponse{
			Code:    400,
			Message: "User ID, new email, and password are required",
		}, nil
	}

	// 2. 解析 UUID
	userID, err := uuid.Parse(in.Id)
	if err != nil {
		return &user.ChangeEmailResponse{
			Code:    400,
			Message: "Invalid user ID format",
		}, nil
	}

	// 验证新邮箱格式
	if !ValidateEmail(in.NewEmail) {
		return &user.ChangeEmailResponse{
			Code:    400,
			Message: "Invalid email format",
		}, nil
	}

	// 3. 查找用户
	existingUser, err := l.svcCtx.UserModel.FindByID(l.ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &user.ChangeEmailResponse{
				Code:    404,
				Message: "User not found",
			}, nil
		}
		l.Logger.Errorf("Failed to find user: %v", err)
		return &user.ChangeEmailResponse{
			Code:    500,
			Message: "Internal server error",
		}, nil
	}

	// 4. 验证密码
	err = bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(in.Password))
	if err != nil {
		return &user.ChangeEmailResponse{
			Code:    401,
			Message: "Password is incorrect",
		}, nil
	}

	// 5. 检查新邮箱是否已被使用
	existingEmailUser, err := l.svcCtx.UserModel.FindByEmail(l.ctx, in.NewEmail)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		l.Logger.Errorf("Failed to check email: %v", err)
		return &user.ChangeEmailResponse{
			Code:    500,
			Message: "Internal server error",
		}, nil
	}
	if existingEmailUser != nil && existingEmailUser.ID != userID {
		return &user.ChangeEmailResponse{
			Code:    400,
			Message: "Email already exists",
		}, nil
	}

	// TODO: 验证 code（验证码）如果需要的话

	// 6. 更新邮箱
	err = l.svcCtx.UserModel.UpdateEmail(l.ctx, userID, in.NewEmail)
	if err != nil {
		l.Logger.Errorf("Failed to update email: %v", err)
		return &user.ChangeEmailResponse{
			Code:    500,
			Message: "Failed to update email",
		}, nil
	}

	return &user.ChangeEmailResponse{
		Code:    200,
		Message: "Email changed successfully",
	}, nil
}
