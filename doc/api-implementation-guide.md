# User API 实现步骤总结

## 🎯 实现概览

User API 层通过调用 User RPC 服务实现用户管理功能，包含 8 个核心接口。

---

## 📝 实现步骤

### 1️⃣ 定义 API（user.api）

```api
// 定义请求/响应类型
type LoginRequest {
    Email    string `json:"email"`
    Password string `json:"password"`
}

type LoginResponse {
    Code    int64     `json:"code"`
    Message string    `json:"message"`
    Data    LoginData `json:"data"`
}

// 定义路由和 handler
@server(group: auth)
service user-api {
    @handler LoginHandler
    post /api/v1/auth/login (LoginRequest) returns (LoginResponse)
}
```

### 2️⃣ 生成代码

```bash
cd service/user/api
goctl api go -api user.api --style=goZero -dir .
```

**生成内容**：
- `internal/types/types.go` - 类型定义
- `internal/handler/` - HTTP 处理器
- `internal/logic/` - 业务逻辑（需要填充实现）

### 3️⃣ 配置 RPC 客户端

**config.go**:
```go
type Config struct {
    rest.RestConf
    Auth struct {
        AccessSecret string
        AccessExpire int64
    }
    UserRpc zrpc.RpcClientConf  // ← 添加 RPC 配置
}
```

**user-api.yaml**:
```yaml
Port: 8101
Auth:
  AccessSecret: your-jwt-secret
  AccessExpire: 7200

UserRpc:  # ← RPC 服务发现配置
  Etcd:
    Hosts:
      - localhost:2379
    Key: user.rpc
```

### 4️⃣ 初始化服务上下文

**serviceContext.go**:
```go
type ServiceContext struct {
    Config  config.Config
    JwtAuth rest.Middleware
    UserRpc userservice.UserService  // ← RPC 客户端
}

func NewServiceContext(c config.Config) *ServiceContext {
    return &ServiceContext{
        Config:  c,
        JwtAuth: middleware.NewJwtAuthMiddleware().Handle,
        UserRpc: userservice.NewUserService(zrpc.MustNewClient(c.UserRpc)),
    }
}
```

### 5️⃣ 实现业务逻辑

**loginLogic.go**:
```go
func (l *LoginLogic) Login(req *types.LoginRequest, w http.ResponseWriter) {
    // 1. 调用 RPC 服务
    rpcResp, err := l.svcCtx.UserRpc.Login(l.ctx, &user.LoginRequest{
        Email:    req.Email,
        Password: req.Password,
    })
    if err != nil {
        response.Error(w, 500, "Internal server error")
        return
    }

    // 2. 检查业务错误
    if rpcResp.Code != 200 {
        response.BusinessError(w, rpcResp.Code, rpcResp.Message)
        return
    }

    // 3. 转换并返回数据
    response.SuccessWithMessage(w, "Login successful", loginData)
}
```

### 6️⃣ 更新 Handler

**loginHandler.go**:
```go
func LoginHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var req types.LoginRequest
        if err := httpx.Parse(r, &req); err != nil {
            httpx.ErrorCtx(r.Context(), w, err)
            return
        }

        l := auth.NewLoginLogic(r.Context(), svcCtx)
        l.Login(&req, w)  // ← 直接调用，logic 负责写响应
    }
}
```

### 7️⃣ 实现 JWT 中间件（可选）

**jwtauthMiddleware.go**:
```go
func (m *JwtAuthMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // 1. 提取 token
        authHeader := r.Header.Get("Authorization")
        
        // 2. 验证 token
        token, err := jwt.Parse(tokenString, ...)
        
        // 3. 提取 userId 到 context
        ctx := context.WithValue(r.Context(), "userId", userId)
        r = r.WithContext(ctx)
        
        next(w, r)
    }
}
```

---

## 🔄 统一响应格式

创建 `response` 包统一处理响应：

```go
// response/response.go
func Success(w http.ResponseWriter, data interface{}) {
    httpx.OkJson(w, Response{
        Success: true,
        Code:    200,
        Message: "Success",
        Data:    data,
    })
}

func Error(w http.ResponseWriter, code int64, message string) {
    httpx.OkJson(w, Response{
        Success: false,
        Code:    code,
        Message: message,
    })
}
```

---

## 📋 已实现的接口

| 路径 | 方法 | 功能 | 认证 |
|------|------|------|------|
| `/api/v1/auth/register` | POST | 用户注册 | ❌ |
| `/api/v1/auth/login` | POST | 用户登录 | ❌ |
| `/api/v1/auth/logout` | POST | 用户登出 | ✅ |
| `/api/v1/user/:id` | GET | 用户信息 | ✅ |
| `/api/v1/user/password` | POST | 修改密码 | ✅ |
| `/api/v1/user/email` | POST | 修改邮箱 | ✅ |
| `/api/v1/user/profile` | PUT | 更新资料 | ✅ |
| `/api/v1/user/:id` | DELETE | 删除用户 | ✅ |

---

## 🚀 启动服务

```bash
# 1. 启动基础服务
docker-compose up -d

# 2. 启动 RPC（端口 8201）
cd service/user/rpc
go run user.go

# 3. 启动 API（端口 8101）
cd service/user/api
go run user.go
```

---

## 💡 核心要点

1. **API 定义** → goctl 生成代码框架
2. **配置 RPC** → 通过 etcd 自动发现服务
3. **实现 Logic** → 调用 RPC，处理响应
4. **统一响应** → 使用 response 包标准化输出
5. **JWT 中间件** → 保护需要认证的接口

参考：[完整文档](HOW_TO_RUN.md)
