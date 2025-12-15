package model

import (
	"context"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ CategoriesModel = (*customCategoriesModel)(nil)

type (
	// CategoriesModel is an interface to be customized, add more methods here,
	// and implement the added methods in customCategoriesModel.
	CategoriesModel interface {
		categoriesModel
		FindAll(ctx context.Context) ([]*Categories, error)
	}

	customCategoriesModel struct {
		*defaultCategoriesModel
	}
)

// NewCategoriesModel returns a model for the database table.
func NewCategoriesModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) CategoriesModel {
	return &customCategoriesModel{
		defaultCategoriesModel: newCategoriesModel(conn, c, opts...),
	}
}

// FindAll 查询所有分类（状态为 1 的）
func (m *customCategoriesModel) FindAll(ctx context.Context) ([]*Categories, error) {
	var resp []*Categories
	query := `SELECT id, user_id, name, icon, color, type, sort_order, is_system, status, created_at, updated_at 
	          FROM categories 
	          WHERE status = 1 
	          ORDER BY type, sort_order`

	err := m.QueryRowsNoCacheCtx(ctx, &resp, query)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
