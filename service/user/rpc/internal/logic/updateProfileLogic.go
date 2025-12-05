package logic

import (
	"context"

	"xledger/service/user/rpc/internal/svc"
	"xledger/service/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
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
	// todo: add your logic here and delete this line

	return &user.UpdateProfileResponse{}, nil
}
