package logic

import (
	"context"
	"errors"

	"xledger/service/user/rpc/internal/svc"
	"xledger/service/user/rpc/user"

	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type UpdateProfileLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateProfileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateProfileLogic {
	return &UpdateProfileLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateProfileLogic) UpdateProfile(in *user.UpdateProfileRequest) (*user.UpdateProfileResponse, error) {
	// 1. 验证输入
	if in.Id == "" {
		return &user.UpdateProfileResponse{
			Code:    400,
			Message: "User ID is required",
		}, nil
	}

	// 2. 解析 UUID
	userID, err := uuid.Parse(in.Id)
	if err != nil {
		return &user.UpdateProfileResponse{
			Code:    400,
			Message: "Invalid user ID format",
		}, nil
	}

	// 3. 查找用户
	existingUser, err := l.svcCtx.UserModel.FindByID(l.ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &user.UpdateProfileResponse{
				Code:    404,
				Message: "User not found",
			}, nil
		}
		l.Logger.Errorf("Failed to find user: %v", err)
		return &user.UpdateProfileResponse{
			Code:    500,
			Message: "Internal server error",
		}, nil
	}

	// 4. 更新字段（只更新非空字段）
	needUpdate := false

	if in.Username != "" && in.Username != existingUser.Username {
		existingUser.Username = in.Username
		needUpdate = true
	}

	if in.Gender != "" {
		existingUser.Gender = &in.Gender
		needUpdate = true
	}

	if in.Avatar != "" {
		existingUser.Avatar = &in.Avatar
		needUpdate = true
	}

	if !needUpdate {
		return &user.UpdateProfileResponse{
			Code:    400,
			Message: "No fields to update",
		}, nil
	}

	// 5. 保存更新
	err = l.svcCtx.UserModel.Update(l.ctx, existingUser)
	if err != nil {
		l.Logger.Errorf("Failed to update user: %v", err)
		return &user.UpdateProfileResponse{
			Code:    500,
			Message: "Failed to update profile",
		}, nil
	}

	// 6. 返回更新后的用户信息
	var gender, avatar string
	if existingUser.Gender != nil {
		gender = *existingUser.Gender
	}
	if existingUser.Avatar != nil {
		avatar = *existingUser.Avatar
	}

	return &user.UpdateProfileResponse{
		Code:    200,
		Message: "Profile updated successfully",
		Data: &user.User{
			Id:        existingUser.ID.String(),
			Username:  existingUser.Username,
			Email:     existingUser.Email,
			Gender:    gender,
			Avatar:    avatar,
			Status:    existingUser.Status,
			CreatedAt: existingUser.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: existingUser.UpdatedAt.Format("2006-01-02 15:04:05"),
		},
	}, nil
}
