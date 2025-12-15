# Ledger 服务实现 Walkthrough

按照 User 服务的实现模式，完整实现 Ledger 记账服务的开发指南。

## 目标

实现完整的 Ledger 服务，包括：

- **交易记录管理**：创建、更新、删除、查询交易记录
- **分类管理**：一级分类和子分类的 CRUD
- **统计功能**：收支统计、分类统计、趋势分析

## 📋 实现步骤总览

```mermaid
graph LR
    A[1. 数据库表设计] --> B[2. 定义 Proto/API]
    B --> C[3. 生成代码]
    C --> D[4. 生成 Model]
    D --> E[5. 实现 RPC Logic]
    E --> F[6. 实现 API Logic]
    F --> G[7. 测试验证]
```

---

## 第一步：数据库表设计与迁移

### 1.1 创建迁移文件

**文件**: `model/00002_ledger.sql`

```sql
-- 分类表（简化命名）
CREATE TABLE categories (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    name VARCHAR(100) NOT NULL,
    icon VARCHAR(50),
    color VARCHAR(20),
    type INT NOT NULL,           -- 1: 收入, 2: 支出
    sort_order INT DEFAULT 0,
    is_system BOOLEAN DEFAULT FALSE,
    status INT DEFAULT 1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_categories_user ON categories(user_id, status);

-- 子分类表
CREATE TABLE sub_categories (
    id BIGSERIAL PRIMARY KEY,
    category_id BIGINT NOT NULL REFERENCES categories(id),
    user_id BIGINT NOT NULL,
    name VARCHAR(100) NOT NULL,
    icon VARCHAR(50),
    color VARCHAR(20),
    sort_order INT DEFAULT 0,
    status INT DEFAULT 1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_sub_categories_category ON sub_categories(category_id, status);

-- 交易记录表
CREATE TABLE transactions (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    category_id BIGINT NOT NULL REFERENCES categories(id),
    sub_category_id BIGINT REFERENCES sub_categories(id),
    amount DECIMAL(15, 2) NOT NULL,
    type INT NOT NULL,           -- 1: 收入, 2: 支出
    description TEXT,
    note TEXT,
    transaction_date TIMESTAMP NOT NULL,
    tags TEXT[],
    location VARCHAR(200),
    images TEXT[],
    status INT DEFAULT 1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_transactions_user ON transactions(user_id, transaction_date DESC);
CREATE INDEX idx_transactions_category ON transactions(category_id);
CREATE INDEX idx_transactions_date ON transactions(transaction_date);

-- 插入默认分类数据
-- 支出分类
INSERT INTO categories (user_id, name, icon, color, type, sort_order, is_system, status) VALUES
(0, '餐饮', '🍔', '#FF6B6B', 2, 1, true, 1),
(0, '交通', '🚗', '#4ECDC4', 2, 2, true, 1),
(0, '购物', '🛍️', '#95E1D3', 2, 3, true, 1),
(0, '娱乐', '🎮', '#F38181', 2, 4, true, 1),
(0, '住房', '🏠', '#AA96DA', 2, 5, true, 1),
(0, '医疗', '⚕️', '#FCBAD3', 2, 6, true, 1),
(0, '教育', '📚', '#A8D8EA', 2, 7, true, 1),
(0, '其他支出', '💸', '#FFFFD2', 2, 99, true, 1);

-- 收入分类
INSERT INTO categories (user_id, name, icon, color, type, sort_order, is_system, status) VALUES
(0, '工资', '💼', '#51CF66', 1, 1, true, 1),
(0, '奖金', '🎁', '#69DB7C', 1, 2, true, 1),
(0, '投资收益', '📈', '#8CE99A', 1, 3, true, 1),
(0, '兼职', '💻', '#B2F2BB', 1, 4, true, 1),
(0, '其他收入', '💰', '#D3F9D8', 1, 99, true, 1);
```

### 1.2 添加 Makefile 指令

在项目根目录的 `Makefile` 中添加：

```makefile
# 数据库迁移
.PHONY: migrate-ledger
migrate-ledger:
 @echo "Running ledger migration..."
 docker exec -i xledger-postgres psql -U admin -d xledger < model/00002_ledger.sql
 @echo "Ledger migration completed!"

# 回滚迁移（如需要）
.PHONY: migrate-ledger-down
migrate-ledger-down:
 @echo "Rolling back ledger migration..."
 docker exec -i xledger-postgres psql -U admin -d xledger -c "DROP TABLE IF EXISTS transactions CASCADE;"
 docker exec -i xledger-postgres psql -U admin -d xledger -c "DROP TABLE IF EXISTS sub_categories CASCADE;"
 docker exec -i xledger-postgres psql -U admin -d xledger -c "DROP TABLE IF EXISTS categories CASCADE;"
 @echo "Rollback completed!"
```

### 1.3 执行迁移

```bash
make migrate-ledger
```

---

## 第二步：定义 Proto 和 API 文件

> **重要**：先定义好接口规范，再生成代码骨架！

### 2.1 Proto 文件已存在

Proto 文件已定义：`service/ledger/rpc/ledger.proto`

需要确认并调整（将 `LedgerRecord` 改为 `Transaction`）

### 2.2 API 文件已存在

API 文件已定义：`service/ledger/api/ledger.api`

同样需要调整命名以匹配数据库表名

---

## 第三步：生成 RPC 和 API 代码

### 3.1 添加 Makefile 指令

在 `Makefile` 中添加代码生成指令：

```makefile
# 生成 Ledger RPC 代码
.PHONY: gen-ledger-rpc
gen-ledger-rpc:
 @echo "Generating Ledger RPC code..."
 cd service/ledger/rpc && goctl rpc protoc ledger.proto --go_out=. --go-grpc_out=. --zrpc_out=.
 @echo "Ledger RPC code generated!"

# 生成 Ledger API 代码
.PHONY: gen-ledger-api
gen-ledger-api:
 @echo "Generating Ledger API code..."
 cd service/ledger/api && goctl api go -api ledger.api -dir .
 @echo "Ledger API code generated!"

# 一键生成所有 Ledger 代码
.PHONY: gen-ledger
gen-ledger: gen-ledger-rpc gen-ledger-api
 @echo "All Ledger code generated!"
```

### 3.2 执行代码生成

```bash
# 生成 RPC 代码
make gen-ledger-rpc

# 生成 API 代码
make gen-ledger-api

# 或一键生成
make gen-ledger
```

---

## 第四步：生成 GORM Model

> **注意**：使用 goctl 工具从数据库表生成 Model

### 4.1 添加 Makefile 指令

```makefile
# 生成 Ledger Model（从数据库表生成）
.PHONY: gen-ledger-model
gen-ledger-model:
 @echo "Generating Ledger models from database..."
 cd service/ledger/model && goctl model pg datasource \
  "postgres://admin:123456@localhost:15432/xledger?sslmode=disable" \
  table categories \
  table sub_categories \
  table transactions \
  dir .
 @echo "Ledger models generated!"
```

### 4.2 执行 Model 生成

```bash
make gen-ledger-model
```

### 4.3 生成的文件结构

```
service/ledger/model/
├── categoriesmodel.go
├── subcategoriesmodel.go
├── transactionsmodel.go
└── vars.go                  # 共享变量
```

goctl 会自动生成包含基础 CRUD 方法的 Model，如：

- `Insert`
- `FindOne`
- `Update`
- `Delete`

---

## 第五步：配置 RPC 服务

### 5.1 配置文件

**文件**: `service/ledger/rpc/etc/ledger.yaml`

```yaml
Name: ledger.rpc
ListenOn: 0.0.0.0:8301

# etcd 配置
Etcd:
  Hosts:
    - localhost:2379
  Key: ledger.rpc

# PostgreSQL 数据库配置
DataSource: postgres://admin:${DB_PASSWORD:123456}@localhost:15432/xledger?sslmode=disable

# Redis 缓存配置
CacheRedis:
  - Host: ${REDIS_HOST:localhost:16379}
    Pass: ${REDIS_PASSWORD:redis123}
    Type: node
```

### 5.2 Service Context

**文件**: `service/ledger/rpc/internal/svc/servicecontext.go`

```go
package svc

import (
    "xledger/service/ledger/model"
    "xledger/service/ledger/rpc/internal/config"

    "github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
    Config             config.Config
    TransactionModel   model.TransactionsModel
    CategoryModel      model.CategoriesModel
    SubCategoryModel   model.SubCategoriesModel
}

func NewServiceContext(c config.Config) *ServiceContext {
    // 连接数据库
    conn := sqlx.NewSqlConn("postgres", c.DataSource)

    return &ServiceContext{
        Config:           c,
        TransactionModel: model.NewTransactionsModel(conn, c.CacheRedis),
        CategoryModel:    model.NewCategoriesModel(conn, c.CacheRedis),
        SubCategoryModel: model.NewSubCategoriesModel(conn, c.CacheRedis),
    }
}
```

### 5.3 主入口文件

**文件**: `service/ledger/rpc/ledger.go`

```go
package main

import (
    "flag"
    "fmt"

    "xledger/service/ledger/rpc/internal/config"
    "xledger/service/ledger/rpc/internal/server"
    "xledger/service/ledger/rpc/internal/svc"
    "xledger/service/ledger/rpc/ledger"

    "github.com/joho/godotenv"
    "github.com/zeromicro/go-zero/core/conf"
    "github.com/zeromicro/go-zero/core/service"
    "github.com/zeromicro/go-zero/zrpc"
    "google.golang.org/grpc"
    "google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/ledger.yaml", "the config file")

func main() {
    flag.Parse()

    // 加载 .env 文件
    _ = godotenv.Load()

    var c config.Config
    conf.MustLoad(*configFile, &c)
    ctx := svc.NewServiceContext(c)

    s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
        ledger.RegisterLedgerServiceServer(grpcServer, server.NewLedgerServiceServer(ctx))

        if c.Mode == service.DevMode || c.Mode == service.TestMode {
            reflection.Register(grpcServer)
        }
    })
    defer s.Stop()

    fmt.Printf("Starting rpc server at %s...\\n", c.ListenOn)
    s.Start()
}
```

---

## 第六步：实现 RPC Logic

### 示例：创建交易记录

**文件**: `service/ledger/rpc/internal/logic/createledgerrecordlogic.go`

```go
package logic

import (
    "context"
    "time"

    "xledger/service/ledger/model"
    "xledger/service/ledger/rpc/internal/svc"
    "xledger/service/ledger/rpc/ledger"

    "github.com/zeromicro/go-zero/core/logx"
)

type CreateLedgerRecordLogic struct {
    ctx    context.Context
    svcCtx *svc.ServiceContext
    logx.Logger
}

func NewCreateLedgerRecordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateLedgerRecordLogic {
    return &CreateLedgerRecordLogic{
        ctx:    ctx,
        svcCtx: svcCtx,
        Logger: logx.WithContext(ctx),
    }
}

func (l *CreateLedgerRecordLogic) CreateLedgerRecord(in *ledger.CreateLedgerRecordRequest) (*ledger.CreateLedgerRecordResponse, error) {
    // 1. 验证输入
    if in.CategoryId == 0 || in.Amount <= 0 {
        return &ledger.CreateLedgerRecordResponse{
            Code:    400,
            Message: "Invalid input: category_id and amount are required",
        }, nil
    }

    // 2. 获取用户 ID（从 metadata）
    userID := getUserIDFromContext(l.ctx)

    // 3. 解析交易日期
    transactionDate, err := time.Parse("2006-01-02", in.RecordDate)
    if err != nil {
        return &ledger.CreateLedgerRecordResponse{
            Code:    400,
            Message: "Invalid date format, use YYYY-MM-DD",
        }, nil
    }

    // 4. 创建交易记录
    transaction := &model.Transactions{
        UserId:          userID,
        CategoryId:      in.CategoryId,
        Amount:          in.Amount,
        Type:            in.Type,
        Description:     sql.NullString{String: in.Description, Valid: in.Description != ""},
        Note:            sql.NullString{String: in.Note, Valid: in.Note != ""},
        TransactionDate: transactionDate,
        // Tags, Location, Images 需要使用 pq.Array
        Status:          1,
    }

    if in.SubCategoryId > 0 {
        transaction.SubCategoryId = sql.NullInt64{Int64: in.SubCategoryId, Valid: true}
    }

    result, err := l.svcCtx.TransactionModel.Insert(l.ctx, transaction)
    if err != nil {
        l.Logger.Errorf("Failed to create transaction: %v", err)
        return &ledger.CreateLedgerRecordResponse{
            Code:    500,
            Message: "Failed to create transaction",
        }, nil
    }

    id, _ := result.LastInsertId()
    transaction.Id = id

    // 5. 返回结果
    return &ledger.CreateLedgerRecordResponse{
        Code:    200,
        Message: "Transaction created successfully",
        Data:    convertToProtoTransaction(transaction),
    }, nil
}
```

---

## 第七步：配置和实现 API 服务

### 7.1 API 配置

**文件**: `service/ledger/api/etc/ledger-api.yaml`

```yaml
Name: ledger-api
Host: 0.0.0.0
Port: 8102

# JWT 配置
Auth:
  AccessSecret: ${JWT_SECRET:xledger-secret-key-change-in-production}
  AccessExpire: ${JWT_ACCESS_EXPIRE:7200}

# Redis 配置
Redis:
  Host: ${REDIS_HOST:localhost:16379}
  Pass: ${REDIS_PASSWORD:redis123}
  DB: ${REDIS_DB:0}

# Ledger RPC 服务
LedgerRpc:
  Etcd:
    Hosts:
      - localhost:2379
    Key: ledger.rpc
```

### 7.2 Service Context

**文件**: `service/ledger/api/internal/svc/servicecontext.go`

```go
package svc

import (
    commonconfig "xledger/common/config"
    "xledger/service/ledger/api/internal/config"
    "xledger/service/ledger/api/internal/middleware"
    "xledger/service/ledger/rpc/ledgerservice"

    "github.com/redis/go-redis/v9"
    "github.com/zeromicro/go-zero/rest"
    "github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
    Config      config.Config
    JwtAuth     rest.Middleware
    LedgerRpc   ledgerservice.LedgerService
    RedisClient *redis.Client
}

func NewServiceContext(c config.Config) *ServiceContext {
    // 初始化 Redis
    redisClient := redis.NewClient(&redis.Options{
        Addr:     c.Redis.Host,
        Password: c.Redis.Pass,
        DB:       c.Redis.DB,
    })

    return &ServiceContext{
        Config:      c,
        JwtAuth:     middleware.NewJwtAuthMiddleware(redisClient, c.Auth.AccessSecret).Handle,
        LedgerRpc:   ledgerservice.NewLedgerService(zrpc.MustNewClient(c.LedgerRpc)),
        RedisClient: redisClient,
    }
}
```

---

## 第八步：Makefile 完整配置

在项目根目录 `Makefile` 中整合所有命令：

```makefile
# Ledger 服务相关命令
.PHONY: ledger-migrate ledger-gen ledger-run

# 数据库迁移
ledger-migrate: migrate-ledger

# 代码生成
ledger-gen: gen-ledger gen-ledger-model

# 启动服务
ledger-run-rpc:
 cd service/ledger/rpc && go run ledger.go

ledger-run-api:
 cd service/ledger/api && go run ledger.go

# 完整流程（初次部署）
ledger-init: ledger-migrate ledger-gen
 @echo "Ledger service initialized!"
```

使用方式：

```bash
# 1. 初次部署
make ledger-init

# 2. 启动 RPC
make ledger-run-rpc

# 3. 启动 API（新终端）
make ledger-run-api
```

---

## 关键实现要点

### 1. 默认分类初始化

系统分类（`is_system = true`，`user_id = 0`）作为模板，用户首次登录时复制：

```go
func (l *InitUserCategoriesLogic) InitCategories(userID int64) error {
    // 查询系统默认分类
    systemCats, err := l.svcCtx.CategoryModel.FindSystemCategories(l.ctx)
    
    // 为用户创建副本
    for _, cat := range systemCats {
        userCat := cat
        userCat.UserId = userID
        userCat.IsSystem = false
        l.svcCtx.CategoryModel.Insert(l.ctx, &userCat)
    }
}
```

### 2. 数据权限验证

```go
// 在 Logic 中验证数据归属
transaction, err := l.svcCtx.TransactionModel.FindOne(l.ctx, id)
if transaction.UserId != getUserIDFromContext(l.ctx) {
    return nil, errors.New("unauthorized")
}
```

### 3. 统计查询优化

```sql
SELECT 
    category_id,
    SUM(CASE WHEN type = 1 THEN amount ELSE 0 END) as income,
    SUM(CASE WHEN type = 2 THEN amount ELSE 0 END) as expense
FROM transactions
WHERE user_id = ? AND transaction_date BETWEEN ? AND ?
GROUP BY category_id
```

---

## 测试验证

```bash
# 创建交易
curl -X POST http://localhost:8102/api/ledger/records \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "category_id": 1,
    "amount": 50.00,
    "type": 2,
    "description": "午餐",
    "transaction_date": "2025-12-16"
  }'
```

---

## 参考资料

- User 服务: `service/user/`
- Proto: `service/ledger/rpc/ledger.proto`
- API: `service/ledger/api/ledger.api`
- Goctl 文档: <https://go-zero.dev/docs/tutorials>
