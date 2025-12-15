// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package ledger

import (
	"context"

	"xledger/service/ledger/api/internal/svc"
	"xledger/service/ledger/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteSubCategoryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteSubCategoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteSubCategoryLogic {
	return &DeleteSubCategoryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteSubCategoryLogic) DeleteSubCategory(req *types.DeleteSubCategoryRequest) (resp *types.CommonResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
