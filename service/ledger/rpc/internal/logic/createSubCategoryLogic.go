package logic

import (
	"context"

	"xledger/service/ledger/rpc/internal/svc"
	"xledger/service/ledger/rpc/ledger"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateSubCategoryLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateSubCategoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateSubCategoryLogic {
	return &CreateSubCategoryLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 子分类操作
func (l *CreateSubCategoryLogic) CreateSubCategory(in *ledger.CreateSubCategoryRequest) (*ledger.CreateSubCategoryResponse, error) {
	// todo: add your logic here and delete this line

	return &ledger.CreateSubCategoryResponse{}, nil
}
