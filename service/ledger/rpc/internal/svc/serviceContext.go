package svc

import (
	"xledger/service/ledger/model"
	"xledger/service/ledger/rpc/internal/config"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config             config.Config
	CategoriesModel    model.CategoriesModel
	SubCategoriesModel model.SubCategoriesModel
	TransactionsModel  model.TransactionsModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewSqlConn("postgres", c.DataSource)

	return &ServiceContext{
		Config:             c,
		CategoriesModel:    model.NewCategoriesModel(conn, c.CacheRedis),
		SubCategoriesModel: model.NewSubCategoriesModel(conn, c.CacheRedis),
		TransactionsModel:  model.NewTransactionsModel(conn, c.CacheRedis),
	}
}
