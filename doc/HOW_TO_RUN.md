# XLedger Zero - 项目运行指南

本指南将帮助您快速启动和运行 XLedger Zero 项目。

## 📋 前置要求

### 必需软件
- **Go** >= 1.21
- **Docker** & **Docker Compose**
- **goctl** (go-zero 工具)

### 安装 goctl
```bash
go install github.com/zeromicro/go-zero/tools/goctl@latest
```

---

## 🚀 快速启动

### 1. 启动基础服务

使用 Docker Compose 启动所需的基础服务：

```bash
# 在项目根目录执行
docker-compose up -d
```

这将启动以下服务：
- **PostgreSQL** (端口 15432)
- **Redis** (端口 16379)  
- **etcd** (端口 2379)

验证服务状态：
```bash
docker-compose ps
```

### 2. 初始化数据库

创建数据库和表：

```bash
# 连接到 PostgreSQL
docker exec -it xledger-postgres psql -U admin -d xledger

# 在 psql 中执行（如果有初始化脚本）
# 或手动创建表
```

**用户表结构**（如果需要手动创建）：
```sql
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username VARCHAR(50) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    gender VARCHAR(10),
    avatar TEXT,
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_deleted_at ON users(deleted_at);
```

### 3. 启动 User RPC 服务

```bash
cd service/user/rpc
go run user.go
```

服务将在 **8201** 端口启动并注册到 etcd。

### 4. 启动 User API 服务

在新终端中：

```bash
cd service/user/api
go run user.go
```

API 服务将在 **8101** 端口启动。

---

## 🏭 生产环境部署

**⚠️ 以上是开发环境的快速启动方式，不适合生产环境！**

生产环境部署请查看详细指南：
- **[生产部署文档](PRODUCTION_DEPLOYMENT.md)**
  - systemd 服务配置
  - Docker 容器化部署
  - Nginx 反向代理
  - SSL/HTTPS 配置
  - 监控和日志管理
  - 安全加固

---

## 📡 测试 API

### 方式1：使用 OpenAPI 文档 + Apifox

1. 导入 OpenAPI 文档到 Apifox：
   ```
   doc/user.openapi.yaml
   ```

2. 配置环境：
   - Base URL: `http://localhost:8101`

3. 测试接口

### 方式2：使用 curl

**注册用户**：
```bash
curl -X POST http://localhost:8101/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "张三",
    "email": "zhangsan@example.com",
    "password": "password123"
  }'
```

**登录**：
```bash
curl -X POST http://localhost:8101/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "zhangsan@example.com",
    "password": "password123"
  }'
```

**获取用户信息**（需要携带 token）：
```bash
curl -X GET http://localhost:8101/api/v1/user/{user_id} \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

---

## 🔧 配置说明

### 数据库配置

**RPC**: `service/user/rpc/etc/user.yaml`
```yaml
DataSource: postgres://admin:123456@localhost:15432/xledger?sslmode=disable
```

### Redis 配置

**RPC**: `service/user/rpc/etc/user.yaml`
```yaml
CacheRedis:
  - Host: localhost:16379
    Pass: redis123
    Type: node
```

### JWT 配置

⚠️ **重要**：API 和 RPC 的 JWT Secret 必须相同

**RPC**: `service/user/rpc/etc/user.yaml`
```yaml
JwtSecret: xledger-secret-key-change-in-production
```

**API**: `service/user/api/etc/user-api.yaml`
```yaml
Auth:
  AccessSecret: xledger-secret-key-change-in-production
```

### 使用环境变量（推荐）

```bash
export JWT_SECRET="your-production-secret-key"
export DB_PASSWORD="your-db-password"
export REDIS_PASSWORD="your-redis-password"
```

---

## 📁 项目结构

```
xledger-zero/
├── service/
│   └── user/
│       ├── api/              # API Gateway 服务
│       │   ├── etc/          # 配置文件
│       │   ├── internal/
│       │   │   ├── handler/  # HTTP 处理器
│       │   │   ├── logic/    # 业务逻辑
│       │   │   ├── response/ # 统一响应
│       │   │   └── types/    # 类型定义
│       │   └── user.go       # 主入口
│       ├── rpc/              # RPC 服务
│       │   ├── etc/          # 配置文件
│       │   ├── internal/
│       │   │   ├── logic/    # RPC 业务逻辑
│       │   │   └── svc/      # 服务上下文
│       │   └── user.go       # 主入口
│       └── model/            # 数据模型层
│           ├── user.go       # GORM 模型
│           └── userModel.go  # 模型接口
├── doc/
│   └── user.openapi.yaml    # OpenAPI 文档
└── docker-compose.yml        # Docker 编排
```

---

## 🛠️ 开发模式

### 热重载（推荐）

使用 `air` 进行热重载开发：

```bash
# 安装 air
go install github.com/cosmtrek/air@latest

# RPC 服务热重载
cd service/user/rpc
air

# API 服务热重载
cd service/user/api
air
```

### 重新生成代码

当修改 `.api` 或 `.proto` 文件后：

```bash
# 重新生成 API 代码
cd service/user/api
goctl api go -api user.api --style=goZero -dir .

# 重新生成 RPC 代码
cd service/user/rpc
goctl rpc protoc user.proto --go_out=. --go-grpc_out=. --zrpc_out=. --style=goZero
```

---

## 🐛 故障排查

### 问题1：无法连接到数据库

**检查**：
```bash
docker ps | grep postgres
docker logs xledger-postgres
```

**解决**：
```bash
docker-compose restart postgres
```

### 问题2：etcd 连接失败

**检查**：
```bash
docker logs xledger-etcd
```

**解决**：
```bash
docker-compose restart etcd
```

### 问题3：RPC 服务未注册

**检查 etcd 中的服务**：
```bash
docker exec xledger-etcd etcdctl get --prefix /user.rpc
```

### 问题4：端口被占用

修改配置文件中的端口：
- API: `service/user/api/etc/user-api.yaml` → `Port: 8888`
- RPC: `service/user/rpc/etc/user.yaml` → `ListenOn: 0.0.0.0:8001`

---

## 📚 API 文档

- **OpenAPI 文档**: `doc/user.openapi.yaml`
- **所有接口前缀**: `/api/v1`

### 可用接口

| 方法 | 路径 | 说明 | 需要认证 |
|------|------|------|----------|
| POST | `/api/v1/auth/register` | 用户注册 | ❌ |
| POST | `/api/v1/auth/login` | 用户登录 | ❌ |
| POST | `/api/v1/auth/logout` | 用户登出 | ✅ |
| GET | `/api/v1/user/:id` | 获取用户信息 | ✅ |
| POST | `/api/v1/user/password` | 修改密码 | ✅ |
| POST | `/api/v1/user/email` | 修改邮箱 | ✅ |
| PUT | `/api/v1/user/profile` | 更新资料 | ✅ |
| DELETE | `/api/v1/user/:id` | 删除用户 | ✅ |

---

## ⚙️ 生产部署建议

1. **使用环境变量管理敏感信息**
2. **配置反向代理**（Nginx/Traefik）
3. **启用 HTTPS**
4. **配置日志收集**
5. **设置监控和告警**
6. **使用专业的密钥管理工具**

---

## 💡 提示

- **热重载**：使用 `air` 提升开发效率
- **OpenAPI**：导入到 Apifox/Postman 快速测试
- **日志**：所有服务都输出到标准输出
- **JWT Secret**：生产环境务必修改默认值

---

需要帮助？查看 `doc/` 目录下的其他文档。
