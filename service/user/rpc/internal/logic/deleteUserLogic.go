package logic

import (
	"context"

	"xledger/service/user/rpc/internal/svc"
	"xledger/service/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
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
	// todo: add your logic here and delete this line

	return &user.DeleteUserResponse{}, nil
}
