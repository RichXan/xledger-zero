package logic

import (
	"context"

	"xledger/service/user/rpc/internal/svc"
	"xledger/service/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
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
	// todo: add your logic here and delete this line

	return &user.ChangePasswordResponse{}, nil
}
