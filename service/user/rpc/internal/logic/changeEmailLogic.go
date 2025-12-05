package logic

import (
	"context"

	"xledger/service/user/rpc/internal/svc"
	"xledger/service/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
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
	// todo: add your logic here and delete this line

	return &user.ChangeEmailResponse{}, nil
}
