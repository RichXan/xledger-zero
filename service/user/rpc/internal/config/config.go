package config

import (
	"github.com/zeromicro/go-zero/zrpc"
	"github.com/zeromicro/go-zero/zrpc/redisconf"
)

type Config struct {
	zrpc.RpcServerConf
	DataSource string
	CacheRedis redisconf.RedisConf
}
