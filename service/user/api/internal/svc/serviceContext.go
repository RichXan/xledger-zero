// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package svc

import (
	"xledger/service/user/api/internal/config"
	"xledger/service/user/api/internal/middleware"
	"xledger/service/user/rpc/userservice"

	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config  config.Config
	JwtAuth rest.Middleware
	UserRpc userservice.UserService
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:  c,
		JwtAuth: middleware.NewJwtAuthMiddleware().Handle,
		UserRpc: userservice.NewUserService(zrpc.MustNewClient(c.UserRpc)),
	}
}
