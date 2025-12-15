// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package ledger

import (
	"context"

	"xledger/service/ledger/api/internal/svc"
	"xledger/service/ledger/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateCategoryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateCategoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateCategoryLogic {
	return &UpdateCategoryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateCategoryLogic) UpdateCategory(req *types.UpdateCategoryRequest) (resp *types.CategoryResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
