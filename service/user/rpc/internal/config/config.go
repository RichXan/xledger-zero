package config

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	DataSource       string
	CacheRedis       cache.CacheConf
	JwtSecret        string // JWT 签名密钥
	JwtAccessExpire  int64  // Access Token 过期时间（秒）
	JwtRefreshExpire int64  // Refresh Token 过期时间（秒）
}
