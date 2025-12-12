# RPC 服务业务逻辑实现指南

## 📋 概述

你已经使用 `goctl` 生成了完整的 RPC 代码框架。现在需要实现具体的业务逻辑。

## 🎯 需要修改的文件

### 1. **配置文件** (必须先配置)

#### `service/user/rpc/etc/user.yaml`
添加数据库和 Redis 配置：

```yaml
Name: user.rpc
ListenOn: 0.0.0.0:8001

# 数据库配置
DataSource: postgres://xledger:password@localhost:5432/xledger?sslmode=disable

# Redis 配置
CacheRedis:
  - Host: localhost:6379
    Pass: ""
    Type: node
```

---

### 2. **配置结构** (添加数据库配置)

#### `service/user/rpc/internal/config/config.go`
```go
package config

import (
    "github.com/zeromicro/go-zero/core/stores/cache"
    "github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
    zrpc.RpcServerConf
    DataSource string              // 数据库连接字符串
    CacheRedis cache.CacheConf     // Redis 缓存配置
}
```

---

### 3. **创建数据模型** (使用 GORM)

#### `service/user/rpc/model/userModel.go`
```go
package model

import (
    "context"
    "time"

    "github.com/google/uuid"
    "gorm.io/gorm"
)

// User 用户模型
type User struct {
    ID        uuid.UUID      `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
    Username  string         `gorm:"type:varchar(255);not null"`
    Email     string         `gorm:"type:varchar(255);not null;uniqueIndex"`
    Password  string         `gorm:"type:varchar(255);not null"`
    Gender    *string        `gorm:"type:varchar(50)"`
    Avatar    *string        `gorm:"type:varchar(255)"`
    Status    string         `gorm:"type:varchar(30);default:'active'"`
    CreatedAt time.Time      `gorm:"autoCreateTime"`
    UpdatedAt time.Time      `gorm:"autoUpdateTime"`
    DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (User) TableName() string {
    return "user"
}

// UserModel 用户数据访问接口
type UserModel interface {
    Create(ctx context.Context, user *User) error
    FindByID(ctx context.Context, id uuid.UUID) (*User, error)
    FindByEmail(ctx context.Context, email string) (*User, error)
    Update(ctx context.Context, user *User) error
    Delete(ctx context.Context, id uuid.UUID) error
}

type defaultUserModel struct {
    db *gorm.DB
}

func NewUserModel(db *gorm.DB) UserModel {
    return &defaultUserModel{db: db}
}

func (m *defaultUserModel) Create(ctx context.Context, user *User) error {
    return m.db.WithContext(ctx).Create(user).Error
}

func (m *defaultUserModel) FindByID(ctx context.Context, id uuid.UUID) (*User, error) {
    var user User
    err := m.db.WithContext(ctx).Where("id = ?", id).First(&user).Error
    if err != nil {
        return nil, err
    }
    return &user, nil
}

func (m *defaultUserModel) FindByEmail(ctx context.Context, email string) (*User, error) {
    var user User
    err := m.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
    if err != nil {
        return nil, err
    }
    return &user, nil
}

func (m *defaultUserModel) Update(ctx context.Context, user *User) error {
    return m.db.WithContext(ctx).Save(user).Error
}

func (m *defaultUserModel) Delete(ctx context.Context, id uuid.UUID) error {
    return m.db.WithContext(ctx).Delete(&User{}, "id = ?", id).Error
}
```

---

### 4. **服务上下文** (依赖注入)

#### `service/user/rpc/internal/svc/serviceContext.go`
```go
package svc

import (
    "xledger/service/user/rpc/internal/config"
    "xledger/service/user/rpc/model"

    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

type ServiceContext struct {
    Config    config.Config
    UserModel model.UserModel
}

func NewServiceContext(c config.Config) *ServiceContext {
    // 连接数据库
    db, err := gorm.Open(postgres.Open(c.DataSource), &gorm.Config{})
    if err != nil {
        panic("failed to connect database: " + err.Error())
    }

    return &ServiceContext{
        Config:    c,
        UserModel: model.NewUserModel(db),
    }
}
```

---

### 5. **业务逻辑实现**

现在可以在 `internal/logic/` 目录下实现具体业务逻辑：

#### `internal/logic/registerLogic.go`
```go
package logic

import (
    "context"
    "errors"

    "xledger/service/user/rpc/internal/svc"
    "xledger/service/user/rpc/model"
    "xledger/service/user/rpc/user"

    "github.com/google/uuid"
    "github.com/zeromicro/go-zero/core/logx"
    "golang.org/x/crypto/bcrypt"
)

type RegisterLogic struct {
    ctx    context.Context
    svcCtx *svc.ServiceContext
    logx.Logger
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
    return &RegisterLogic{
        ctx:    ctx,
        svcCtx: svcCtx,
        Logger: logx.WithContext(ctx),
    }
}

func (l *RegisterLogic) Register(in *user.RegisterRequest) (*user.RegisterResponse, error) {
    // 1. 验证输入
    if in.Email == "" || in.Password == "" || in.Username == "" {
        return &user.RegisterResponse{
            Code:    400,
            Message: "Missing required fields",
        }, nil
    }

    // 2. 检查邮箱是否已存在
    existingUser, _ := l.svcCtx.UserModel.FindByEmail(l.ctx, in.Email)
    if existingUser != nil {
        return &user.RegisterResponse{
            Code:    400,
            Message: "Email already exists",
        }, nil
    }

    // 3. 密码加密
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
    if err != nil {
        l.Logger.Errorf("Failed to hash password: %v", err)
        return &user.RegisterResponse{
            Code:    500,
            Message: "Internal server error",
        }, nil
    }

    // 4. 创建用户
    newUser := &model.User{
        ID:       uuid.New(),
        Username: in.Username,
        Email:    in.Email,
        Password: string(hashedPassword),
        Status:   "active",
    }

    err = l.svcCtx.UserModel.Create(l.ctx, newUser)
    if err != nil {
        l.Logger.Errorf("Failed to create user: %v", err)
        return &user.RegisterResponse{
            Code:    500,
            Message: "Failed to create user",
        }, nil
    }

    // 5. 返回结果
    return &user.RegisterResponse{
        Code:    200,
        Message: "User registered successfully",
        Data: &user.User{
            Id:        newUser.ID.String(),
            Username:  newUser.Username,
            Email:     newUser.Email,
            Status:    newUser.Status,
            CreatedAt: newUser.CreatedAt.Format("2006-01-02 15:04:05"),
            UpdatedAt: newUser.UpdatedAt.Format("2006-01-02 15:04:05"),
        },
    }, nil
}
```

#### `internal/logic/loginLogic.go`
```go
package logic

import (
    "context"
    "time"

    "xledger/service/user/rpc/internal/svc"
    "xledger/service/user/rpc/user"

    "github.com/golang-jwt/jwt/v4"
    "github.com/zeromicro/go-zero/core/logx"
    "golang.org/x/crypto/bcrypt"
)

type LoginLogic struct {
    ctx    context.Context
    svcCtx *svc.ServiceContext
    logx.Logger
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
    return &LoginLogic{
        ctx:    ctx,
        svcCtx: svcCtx,
        Logger: logx.WithContext(ctx),
    }
}

func (l *LoginLogic) Login(in *user.LoginRequest) (*user.LoginResponse, error) {
    // 1. 查找用户
    existingUser, err := l.svcCtx.UserModel.FindByEmail(l.ctx, in.Email)
    if err != nil {
        return &user.LoginResponse{
            Code:    401,
            Message: "Invalid email or password",
        }, nil
    }

    // 2. 验证密码
    err = bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(in.Password))
    if err != nil {
        return &user.LoginResponse{
            Code:    401,
            Message: "Invalid email or password",
        }, nil
    }

    // 3. 生成 JWT Token
    now := time.Now()
    expiresIn := int64(7200) // 2 hours
    
    accessToken, err := l.generateToken(existingUser.ID.String(), now.Unix(), expiresIn)
    if err != nil {
        l.Logger.Errorf("Failed to generate token: %v", err)
        return &user.LoginResponse{
            Code:    500,
            Message: "Failed to generate token",
        }, nil
    }

    refreshToken, err := l.generateToken(existingUser.ID.String(), now.Unix(), 86400*7) // 7 days
    if err != nil {
        l.Logger.Errorf("Failed to generate refresh token: %v", err)
        return &user.LoginResponse{
            Code:    500,
            Message: "Failed to generate refresh token",
        }, nil
    }

    // 4. 返回结果
    return &user.LoginResponse{
        Code:    200,
        Message: "Login successful",
        Data: &user.LoginData{
            User: &user.User{
                Id:        existingUser.ID.String(),
                Username:  existingUser.Username,
                Email:     existingUser.Email,
                Status:    existingUser.Status,
                CreatedAt: existingUser.CreatedAt.Format("2006-01-02 15:04:05"),
                UpdatedAt: existingUser.UpdatedAt.Format("2006-01-02 15:04:05"),
            },
            AccessToken:  accessToken,
            RefreshToken: refreshToken,
            ExpiresIn:    expiresIn,
            TokenType:    "Bearer",
        },
    }, nil
}

func (l *LoginLogic) generateToken(userID string, iat, exp int64) (string, error) {
    claims := jwt.MapClaims{
        "user_id": userID,
        "iat":     iat,
        "exp":     iat + exp,
    }
    
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte("your-secret-key")) // TODO: 从配置读取
}
```

---

## 📦 需要安装的依赖

```bash
# GORM PostgreSQL 驱动
go get -u gorm.io/gorm
go get -u gorm.io/driver/postgres

# UUID 生成
go get github.com/google/uuid

# bcrypt 密码加密
go get golang.org/x/crypto/bcrypt

# JWT
go get github.com/golang-jwt/jwt/v4
```

---

## 🚀 实现步骤总结

1. ✅ **修改配置文件** → `etc/user.yaml` (添加数据库配置)
2. ✅ **更新配置结构** → `config/config.go` (添加 DataSource 字段)
3. ✅ **创建数据模型** → `model/userModel.go` (GORM 模型和接口)
4. ✅ **更新服务上下文** → `svc/serviceContext.go` (注入 UserModel)
5. ✅ **实现业务逻辑** → `logic/*Logic.go` (在每个 Logic 文件中实现)

---

## 🧪 测试

启动服务后，可以使用 grpcurl 测试：

```bash
# 注册用户
grpcurl -plaintext -d '{
  "username": "testuser",
  "email": "test@example.com",
  "password": "password123"
}' localhost:8001 user.UserService/Register

# 登录
grpcurl -plaintext -d '{
  "email": "test@example.com",
  "password": "password123"
}' localhost:8001 user.UserService/Login
```
