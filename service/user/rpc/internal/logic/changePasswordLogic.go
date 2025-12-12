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

type ChangePasswordLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewChangePasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChangePasswordLogic {
	return &ChangePasswordLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// Updates
func (l *ChangePasswordLogic) ChangePassword(in *user.ChangePasswordRequest) (*user.ChangePasswordResponse, error) {
	// 1. 验证输入
	if in.Id == "" || in.OldPassword == "" || in.NewPassword == "" {
		return &user.ChangePasswordResponse{
			Code:    400,
			Message: "User ID, old password, and new password are required",
		}, nil
	}

	// 2. 检查新密码长度
	if err := ValidatePassword(in.NewPassword); err != nil {
		return &user.ChangePasswordResponse{
			Code:    400,
			Message: err.Error(),
		}, nil
	}

	// 检查新旧密码是否相同
	if in.OldPassword == in.NewPassword {
		return &user.ChangePasswordResponse{
			Code:    400,
			Message: "New password must be different from old password",
		}, nil
	}

	// 3. 解析 UUID
	userID, err := uuid.Parse(in.Id)
	if err != nil {
		return &user.ChangePasswordResponse{
			Code:    400,
			Message: "Invalid user ID format",
		}, nil
	}

	// 4. 查找用户
	existingUser, err := l.svcCtx.UserModel.FindByID(l.ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &user.ChangePasswordResponse{
				Code:    404,
				Message: "User not found",
			}, nil
		}
		l.Logger.Errorf("Failed to find user: %v", err)
		return &user.ChangePasswordResponse{
			Code:    500,
			Message: "Internal server error",
		}, nil
	}

	// 5. 验证旧密码
	err = bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(in.OldPassword))
	if err != nil {
		return &user.ChangePasswordResponse{
			Code:    401,
			Message: "Old password is incorrect",
		}, nil
	}

	// 6. 加密新密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(in.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		l.Logger.Errorf("Failed to hash password: %v", err)
		return &user.ChangePasswordResponse{
			Code:    500,
			Message: "Internal server error",
		}, nil
	}

	// 7. 更新密码
	err = l.svcCtx.UserModel.UpdatePassword(l.ctx, userID, string(hashedPassword))
	if err != nil {
		l.Logger.Errorf("Failed to update password: %v", err)
		return &user.ChangePasswordResponse{
			Code:    500,
			Message: "Failed to update password",
		}, nil
	}

	return &user.ChangePasswordResponse{
		Code:    200,
		Message: "Password changed successfully",
	}, nil
}
