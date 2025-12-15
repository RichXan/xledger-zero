package config

import (
	commonconfig "xledger/common/config"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	DataSource string
	CacheRedis cache.CacheConf
	Auth       commonconfig.Auth // 使用公共 Auth 配置
}
