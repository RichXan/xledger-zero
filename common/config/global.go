package config

// Redis 配置（API 服务使用）
type Redis struct {
	Host string
	Pass string
	DB   int
}

// Auth JWT 认证配置（API 和 RPC 共用）
type Auth struct {
	AccessSecret  string // JWT 签名密钥
	AccessExpire  int64  // Access Token 过期时间（秒）
	RefreshExpire int64  // Refresh Token 过期时间（秒）
}
