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

type UserInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserInfoLogic {
	return &UserInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// User Info
func (l *UserInfoLogic) UserInfo(in *user.UserInfoRequest) (*user.UserInfoResponse, error) {
	// 1. 验证输入
	if in.Id == "" {
		return &user.UserInfoResponse{
			Code:    400,
			Message: "User ID is required",
		}, nil
	}

	// 2. 解析 UUID
	userID, err := uuid.Parse(in.Id)
	if err != nil {
		return &user.UserInfoResponse{
			Code:    400,
			Message: "Invalid user ID format",
		}, nil
	}

	// 3. 查询用户
	existingUser, err := l.svcCtx.UserModel.FindByID(l.ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &user.UserInfoResponse{
				Code:    404,
				Message: "User not found",
			}, nil
		}
		l.Logger.Errorf("Failed to find user: %v", err)
		return &user.UserInfoResponse{
			Code:    500,
			Message: "Internal server error",
		}, nil
	}

	// 4. 返回用户信息
	var gender, avatar string
	if existingUser.Gender != nil {
		gender = *existingUser.Gender
	}
	if existingUser.Avatar != nil {
		avatar = *existingUser.Avatar
	}

	return &user.UserInfoResponse{
		Code:    200,
		Message: "Success",
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
