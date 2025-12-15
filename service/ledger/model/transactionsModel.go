package model

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ TransactionsModel = (*customTransactionsModel)(nil)

type (
	// TransactionsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customTransactionsModel.
	TransactionsModel interface {
		transactionsModel
	}

	customTransactionsModel struct {
		*defaultTransactionsModel
	}
)

// NewTransactionsModel returns a model for the database table.
func NewTransactionsModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) TransactionsModel {
	return &customTransactionsModel{
		defaultTransactionsModel: newTransactionsModel(conn, c, opts...),
	}
}
