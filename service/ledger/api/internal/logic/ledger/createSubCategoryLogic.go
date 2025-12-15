// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package ledger

import (
	"context"

	"xledger/service/ledger/api/internal/svc"
	"xledger/service/ledger/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateSubCategoryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateSubCategoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateSubCategoryLogic {
	return &CreateSubCategoryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateSubCategoryLogic) CreateSubCategory(req *types.CreateSubCategoryRequest) (resp *types.SubCategoryResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
