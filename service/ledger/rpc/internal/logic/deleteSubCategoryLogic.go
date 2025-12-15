package logic

import (
	"context"

	"xledger/service/ledger/rpc/internal/svc"
	"xledger/service/ledger/rpc/ledger"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteSubCategoryLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteSubCategoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteSubCategoryLogic {
	return &DeleteSubCategoryLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteSubCategoryLogic) DeleteSubCategory(in *ledger.DeleteSubCategoryRequest) (*ledger.DeleteSubCategoryResponse, error) {
	// todo: add your logic here and delete this line

	return &ledger.DeleteSubCategoryResponse{}, nil
}
