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

type DeleteUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteUserLogic {
	return &DeleteUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// Account Management
func (l *DeleteUserLogic) DeleteUser(in *user.DeleteUserRequest) (*user.DeleteUserResponse, error) {
	// 1. 验证输入
	if in.Id == "" {
		return &user.DeleteUserResponse{
			Code:    400,
			Message: "User ID is required",
		}, nil
	}

	// 2. 解析 UUID
	userID, err := uuid.Parse(in.Id)
	if err != nil {
		return &user.DeleteUserResponse{
			Code:    400,
			Message: "Invalid user ID format",
		}, nil
	}

	// 3. 检查用户是否存在
	existingUser, err := l.svcCtx.UserModel.FindByID(l.ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &user.DeleteUserResponse{
				Code:    404,
				Message: "User not found",
			}, nil
		}
		l.Logger.Errorf("Failed to find user: %v", err)
		return &user.DeleteUserResponse{
			Code:    500,
			Message: "Internal server error",
		}, nil
	}

	// 4. 软删除用户（GORM 自dynamic处理 deleted_at）
	err = l.svcCtx.UserModel.Delete(l.ctx, userID)
	if err != nil {
		l.Logger.Errorf("Failed to delete user: %v", err)
		return &user.DeleteUserResponse{
			Code:    500,
			Message: "Failed to delete user",
		}, nil
	}

	l.Logger.Infof("User deleted: %s (%s)", existingUser.Email, userID.String())

	return &user.DeleteUserResponse{
		Code:    200,
		Message: "User deleted successfully",
	}, nil
}
