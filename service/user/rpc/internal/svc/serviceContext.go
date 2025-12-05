package svc

import (
	"xledger/service/user/rpc/internal/config"
	"xledger/service/user/model"
)

type ServiceContext struct {
	Config config.Config
	UserModel model.UserModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	db, err := gorm.Open(postgres.Open(c.DataSource), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     c.CacheRedis.Host,
		Password: c.CacheRedis.Pass,
		DB:       c.CacheRedis.DB,
	})

	return &ServiceContext{
		Config:    c,
		UserModel: model.NewUserModel(db, rdb),
	}
}
