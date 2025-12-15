package logic

import (
	"context"

	"xledger/service/ledger/rpc/internal/svc"
	"xledger/service/ledger/rpc/ledger"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCategoryListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetCategoryListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCategoryListLogic {
	return &GetCategoryListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetCategoryListLogic) GetCategoryList(in *ledger.CategoryListRequest) (*ledger.CategoryListResponse, error) {
	// 查询所有状态为 1 的分类
	categories, err := l.svcCtx.CategoriesModel.FindAll(l.ctx)
	if err != nil {
		l.Logger.Errorf("Failed to query categories: %v", err)
		return &ledger.CategoryListResponse{
			Code:    500,
			Message: "Failed to query categories",
			Data:    []*ledger.LedgerCategory{},
		}, nil
	}

	// 转换为 proto 格式
	result := make([]*ledger.LedgerCategory, 0, len(categories))
	for _, cat := range categories {
		result = append(result, &ledger.LedgerCategory{
			Id:        cat.Id,
			UserId:    cat.UserId,
			Name:      cat.Name,
			Icon:      cat.Icon.String,
			Color:     cat.Color.String,
			Type:      cat.Type,
			SortOrder: cat.SortOrder,
			IsSystem:  cat.IsSystem,
			Status:    cat.Status,
			CreatedAt: cat.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: cat.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return &ledger.CategoryListResponse{
		Code:    200,
		Message: "Success",
		Data:    result,
	}, nil
}
