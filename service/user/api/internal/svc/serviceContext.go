// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package svc

import (
	"xledger/service/user/api/internal/config"
	"xledger/service/user/api/internal/middleware"
	"xledger/service/user/rpc/userservice"

	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config      config.Config
	JwtAuth     rest.Middleware
	UserRpc     userservice.UserService
	RedisClient *redis.Client
}

func NewServiceContext(c config.Config) *ServiceContext {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     c.Redis.Host,
		Password: c.Redis.Pass,
		DB:       c.Redis.DB,
	})

	return &ServiceContext{
		Config:      c,
		JwtAuth:     middleware.NewJwtAuthMiddleware(redisClient, c.Auth.AccessSecret).Handle,
		UserRpc:     userservice.NewUserService(zrpc.MustNewClient(c.UserRpc)),
		RedisClient: redisClient,
	}
}
