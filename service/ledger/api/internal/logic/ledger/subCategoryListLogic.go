// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package ledger

import (
	"context"

	"xledger/service/ledger/api/internal/svc"
	"xledger/service/ledger/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SubCategoryListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSubCategoryListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SubCategoryListLogic {
	return &SubCategoryListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SubCategoryListLogic) SubCategoryList(req *types.SubCategoryListRequest) (resp *types.SubCategoryListResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
