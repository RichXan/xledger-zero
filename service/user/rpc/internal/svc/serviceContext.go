package svc

import (
	"xledger/service/user/model"
	"xledger/service/user/rpc/internal/config"

	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config      config.Config
	UserModel   model.UserModel
	RedisClient *redis.Client
}

func NewServiceContext(c config.Config) *ServiceContext {
	db, err := gorm.Open(postgres.Open(c.DataSource), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	// 从 CacheConf 中提取 Redis 配置
	redisConf := c.CacheRedis[0].RedisConf
	rdb := redis.NewClient(&redis.Options{
		Addr:     redisConf.Host,
		Password: redisConf.Pass,
		DB:       0, // 默认使用 DB 0
	})

	return &ServiceContext{
		Config:      c,
		UserModel:   model.NewUserModel(db, rdb),
		RedisClient: rdb,
	}
}
