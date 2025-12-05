# XLedger-Zero 微服务记账系统

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://golang.org)
[![Go-Zero](https://img.shields.io/badge/Go--Zero-1.7+-7C3AED?style=flat)](https://go-zero.dev)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-15+-336791?style=flat&logo=postgresql)](https://www.postgresql.org)
[![Redis](https://img.shields.io/badge/Redis-7+-DC382D?style=flat&logo=redis)](https://redis.io)
[![ory/fosite](https://img.shields.io/badge/ory-fosite-5528FF?style=flat)](https://github.com/ory/fosite)

基于 Go-Zero 框架构建的企业级微服务记账系统，专注于个人财务管理和资产规划。

## 🌟 项目特点

- 📊 **完整记账功能**: 支持收入/支出记录，多维度分类管理
- 📈 **财务报表**: 日报/周报/月报，可视化收支趋势分析
- 💰 **资产规划**: 资产总览、预算管理、财务目标追踪
- 🔐 **企业级认证**: 基于 ory/fosite 的 OAuth2.0/OIDC 认证体系
- 🚀 **微服务架构**: Go-Zero 标准微服务设计，服务独立部署扩展
- 🤖 **AI 增强**: 可选 AI 功能（智能总结、消费洞察）
- 🛡️ **安全可靠**: 多层安全防护，数据隔离，完整审计日志
- 📦 **容器化部署**: Docker Compose 一键启动

## 📐 系统架构

### 技术栈

| 技术 | 版本 | 用途 |
|------|------|------|
| **Go** | 1.21+ | 编程语言 |
| **Go-Zero** | 1.7+ | 微服务框架 |
| **PostgreSQL** | 15+ | 主数据库 |
| **Redis** | 7+ | 缓存/会话存储 |
| **ory/fosite** | latest | OAuth2.0/OIDC 认证 |
| **GORM** | v2+ | ORM 框架 |
| **Docker** | latest | 容器化部署 |

### 微服务架构

```
┌─────────────────────────────────────────────────────────┐
│                     API Gateway                         │
│                  (Go-Zero Gateway)                      │
└────────────┬────────────┬────────────┬──────────────────┘
             │            │            │
     ┌───────▼──────┐ ┌──▼─────────┐ ┌▼────────────────┐
     │ User Service │ │   Ledger   │ │ Report Service  │
     │  (ory/fosite)│ │  Service   │ │                 │
     └──────┬───────┘ └──┬─────────┘ └─┬───────────────┘
            │            │              │
     ┌──────▼────────────▼──────────────▼────────┐
     │           PostgreSQL + Redis              │
     └───────────────────────────────────────────┘
```

## 📁 项目结构

```
xledger-zero/
├── service/                    # 微服务目录
│   ├── user/                   # 用户服务
│   │   ├── api/               # HTTP API 层
│   │   │   ├── etc/           # 配置文件
│   │   │   ├── internal/      # 内部实现
│   │   │   │   ├── config/    # 配置结构
│   │   │   │   ├── handler/   # HTTP 处理器
│   │   │   │   ├── logic/     # 业务逻辑
│   │   │   │   ├── svc/       # 服务上下文
│   │   │   │   └── types/     # 类型定义
│   │   │   └── user.api       # API 定义
│   │   ├── rpc/               # gRPC 服务
│   │   │   ├── etc/
│   │   │   ├── internal/
│   │   │   ├── pb/            # Protobuf 文件
│   │   │   └── user.proto
│   │   └── model/             # 数据模型
│   ├── ledger/                # 账本服务
│   │   ├── api/
│   │   ├── rpc/
│   │   └── model/
│   ├── report/                # 报表服务
│   │   ├── api/
│   │   └── model/
│   └── category/              # 分类服务
│       ├── api/
│       └── model/
├── model/                     # 数据库模型（SQL 迁移）
├── pkg/                       # 共享工具库
│   ├── oauth/                # OAuth2 工具
│   ├── utils/                # 通用工具
│   └── middleware/           # 中间件
├── deploy/                    # 部署配置
│   ├── docker/
│   │   ├── Dockerfile
│   │   └── docker-compose.yml
│   └── k8s/                  # Kubernetes 配置
├── scripts/                   # 脚本工具
└── doc/                       # 文档
```

## 🚀 快速开始

### 环境要求

- Go 1.21+
- Docker & Docker Compose
- PostgreSQL 15+
- Redis 7+

### 本地开发

```bash
# 1. 克隆项目
git clone https://github.com/your-username/xledger-zero.git
cd xledger-zero

# 2. 启动基础服务
docker-compose up -d postgres redis

# 3. 初始化数据库
make migrate-up

# 4. 启动微服务
make run-user    # 用户服务: http://localhost:8001
make run-ledger  # 账本服务: http://localhost:8002
make run-report  # 报表服务: http://localhost:8003

# 或一键启动所有服务
make run-all
```

### Docker 部署

```bash
# 构建并启动所有服务
docker-compose up -d

# 查看服务状态
docker-compose ps

# 查看日志
docker-compose logs -f user-service
```

## 📊 核心功能

### 1. 用户认证 (ory/fosite)

基于 ory/fosite 的企业级 OAuth2.0/OIDC 认证：

```bash
# 注册
POST /api/v1/auth/register
{
  "username": "testuser",
  "email": "test@example.com",
  "password": "SecurePass123!"
}

# 登录获取 Access Token
POST /api/v1/auth/login
{
  "email": "test@example.com",
  "password": "SecurePass123!"
}

# 响应
{
  "access_token": "eyJhbGci...",
  "refresh_token": "eyJhbGci...",
  "expires_in": 3600,
  "token_type": "Bearer"
}
```

### 2. 账本管理

记录和管理日常收支：

```bash
# 创建交易记录
POST /api/v1/transactions
Authorization: Bearer {access_token}
{
  "type": "expense",           # expense | income
  "amount": 128.50,
  "category_id": "uuid",
  "description": "午餐",
  "date": "2025-12-05"
}

# 查询交易记录
GET /api/v1/transactions?start_date=2025-12-01&end_date=2025-12-31

# 响应
{
  "code": 200,
  "data": {
    "transactions": [...],
    "total": 50,
    "page": 1
  }
}
```

### 3. 财务报表

查看日/周/月财务报表：

```bash
# 月度报表
GET /api/v1/reports/monthly?year=2025&month=12
Authorization: Bearer {access_token}

# 响应
{
  "code": 200,
  "data": {
    "period": "2025-12",
    "total_income": 15000.00,
    "total_expense": 8524.50,
    "balance": 6475.50,
    "expense_by_category": [
      {"category": "餐饮", "amount": 2450.00, "percentage": 28.7},
      {"category": "交通", "amount": 1200.00, "percentage": 14.1}
    ],
    "daily_trend": [...]
  }
}

# 周报表
GET /api/v1/reports/weekly?year=2025&week=49

# 日报表
GET /api/v1/reports/daily?date=2025-12-05
```

### 4. 资产规划

```bash
# 查看资产总览
GET /api/v1/assets/overview
Authorization: Bearer {access_token}

# 响应
{
  "code": 200,
  "data": {
    "total_assets": 125000.00,
    "net_worth": 98000.00,
    "monthly_income_avg": 15000.00,
    "monthly_expense_avg": 8500.00,
    "savings_rate": 43.3
  }
}

# 设置预算
POST /api/v1/budgets
{
  "category_id": "uuid",
  "amount": 3000.00,
  "period": "monthly"
}
```

### 5. 分类管理

```bash
# 创建分类
POST /api/v1/categories
{
  "name": "餐饮",
  "type": "expense",
  "icon": "🍔"
}

# 创建子分类
POST /api/v1/categories/{id}/subcategories
{
  "name": "早餐"
}
```

## 🔐 认证与授权

### OAuth2.0 流程 (ory/fosite)

支持标准 OAuth2.0 授权流程：

- **Authorization Code Flow**: Web 应用
- **Client Credentials Flow**: 服务间调用
- **Refresh Token Flow**: Token 刷新

### JWT Token 验证

所有受保护的 API 需要携带有效的 JWT Token：

```bash
curl -H "Authorization: Bearer {access_token}" \
     http://localhost:8002/api/v1/transactions
```

## 🤖 AI 增强功能 (可选)

系统支持可选的 AI 功能增强：

- **智能总结**: 自动生成月度/年度财务总结
- **消费洞察**: AI 分析消费习惯，提供优化建议
- **预测分析**: 基于历史数据预测未来支出
- **异常检测**: 识别异常消费模式

## 🛠️ 开发指南

### 添加新服务

```bash
# 使用 goctl 创建新服务
goctl api new service-name
goctl rpc new service-name
```

### 数据库迁移

```bash
# 创建新迁移
make migrate-create name=add_users_table

# 执行迁移
make migrate-up

# 回滚迁移
make migrate-down
```

### 代码生成

```bash
# 从 .api 文件生成代码
goctl api go -api service.api -dir .

# 从 .proto 文件生成代码
goctl rpc protoc service.proto --go_out=. --go-grpc_out=. --zrpc_out=.
```

## 🧪 测试

```bash
# 运行单元测试
make test

# 运行集成测试
make test-integration

# 代码覆盖率
make coverage
```

## 📦 部署

### 环境变量

```bash
# .env 文件
DB_HOST=localhost
DB_PORT=5432
DB_USER=xledger
DB_PASSWORD=your_password
DB_NAME=xledger

REDIS_HOST=localhost
REDIS_PORT=6379

JWT_SECRET=your_jwt_secret
OAUTH_CLIENT_ID=your_client_id
OAUTH_CLIENT_SECRET=your_client_secret
```

### 生产部署

```bash
# Docker 生产环境
docker-compose -f docker-compose.prod.yml up -d

# Kubernetes
kubectl apply -f deploy/k8s/
```

## 📈 性能指标

- **API 响应时间**: P99 < 100ms
- **并发能力**: 1000+ QPS (单实例)
- **数据库连接池**: 20-50 连接
- **缓存命中率**: > 80%

## 🔒 安全最佳实践

1. **密码加密**: bcrypt 哈希存储
2. **Token 安全**: JWT 短期有效期 + Refresh Token
3. **SQL 注入防护**: 使用参数化查询 (GORM)
4. **XSS 防护**: 输入验证和输出编码
5. **HTTPS**: 生产环境强制 HTTPS
6. **审计日志**: 记录所有敏感操作

## 📝 API 文档

完整的 API 文档请访问：

- **Swagger UI**: http://localhost:8001/swagger
- **API 文档**: [doc/api.md](doc/api.md)

## 🤝 贡献指南

欢迎贡献代码！请阅读 [CONTRIBUTING.md](CONTRIBUTING.md) 了解详情。

## 📄 许可证

本项目采用 MIT 许可证 - 详见 [LICENSE](LICENSE) 文件

## 📧 联系方式

- 作者: xan
- Email: xan@example.com
- 项目主页: https://github.com/your-username/xledger-zero

---

**最后更新**: 2025-12-05  
**版本**: v1.0.0  
**状态**: 开发中 🚧