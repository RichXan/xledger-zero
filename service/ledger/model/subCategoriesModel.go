package model

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ SubCategoriesModel = (*customSubCategoriesModel)(nil)

type (
	// SubCategoriesModel is an interface to be customized, add more methods here,
	// and implement the added methods in customSubCategoriesModel.
	SubCategoriesModel interface {
		subCategoriesModel
	}

	customSubCategoriesModel struct {
		*defaultSubCategoriesModel
	}
)

// NewSubCategoriesModel returns a model for the database table.
func NewSubCategoriesModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) SubCategoriesModel {
	return &customSubCategoriesModel{
		defaultSubCategoriesModel: newSubCategoriesModel(conn, c, opts...),
	}
}
