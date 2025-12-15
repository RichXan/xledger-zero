package logic

import (
	"context"

	"xledger/service/ledger/rpc/internal/svc"
	"xledger/service/ledger/rpc/ledger"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateSubCategoryLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateSubCategoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateSubCategoryLogic {
	return &UpdateSubCategoryLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateSubCategoryLogic) UpdateSubCategory(in *ledger.UpdateSubCategoryRequest) (*ledger.UpdateSubCategoryResponse, error) {
	// todo: add your logic here and delete this line

	return &ledger.UpdateSubCategoryResponse{}, nil
}
