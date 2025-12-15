// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package svc

import (
	"xledger/service/ledger/api/internal/config"
	"xledger/service/ledger/api/internal/middleware"
	"xledger/service/ledger/rpc/ledgerservice"

	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config         config.Config
	AuthMiddleware rest.Middleware
	LedgerRpc      ledgerservice.LedgerService
	RedisClient    *redis.Client
}

func NewServiceContext(c config.Config) *ServiceContext {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     c.Redis.Host,
		Password: c.Redis.Pass,
		DB:       c.Redis.DB,
	})

	return &ServiceContext{
		Config:         c,
		AuthMiddleware: middleware.NewAuthMiddleware(redisClient, c.Auth.AccessSecret).Handle,
		LedgerRpc:      ledgerservice.NewLedgerService(zrpc.MustNewClient(c.LedgerRpc)),
		RedisClient:    redisClient,
	}
}
