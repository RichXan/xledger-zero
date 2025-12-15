package logic

import (
	"context"

	"xledger/service/ledger/rpc/internal/svc"
	"xledger/service/ledger/rpc/ledger"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetSubCategoryListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetSubCategoryListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSubCategoryListLogic {
	return &GetSubCategoryListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetSubCategoryListLogic) GetSubCategoryList(in *ledger.SubCategoryListRequest) (*ledger.SubCategoryListResponse, error) {
	// todo: add your logic here and delete this line

	return &ledger.SubCategoryListResponse{}, nil
}
