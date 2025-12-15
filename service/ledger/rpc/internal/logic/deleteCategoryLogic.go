package logic

import (
	"context"

	"xledger/service/ledger/rpc/internal/svc"
	"xledger/service/ledger/rpc/ledger"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteCategoryLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteCategoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteCategoryLogic {
	return &DeleteCategoryLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteCategoryLogic) DeleteCategory(in *ledger.DeleteCategoryRequest) (*ledger.DeleteCategoryResponse, error) {
	// todo: add your logic here and delete this line

	return &ledger.DeleteCategoryResponse{}, nil
}
