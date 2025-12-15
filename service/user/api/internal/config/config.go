// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package config

import (
	"xledger/common/config"

	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	rest.RestConf
	Auth struct {
		AccessSecret string
		AccessExpire int64
	}
	Redis   config.Redis
	UserRpc zrpc.RpcClientConf
}
