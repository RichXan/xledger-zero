package model

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var _ UserModel = (*customUserModel)(nil)

type (
	// UserModel is an interface to be customized, add more methods here,
	// and implement the added methods in customUserModel.
	UserModel interface {
		FindByID(ctx context.Context, id uuid.UUID) (*User, error)
		FindByEmail(ctx context.Context, email string) (*User, error)
		Create(ctx context.Context, user *User) error
		Update(ctx context.Context, user *User) error
		UpdatePassword(ctx context.Context, id uuid.UUID, password string) error
		UpdateEmail(ctx context.Context, id uuid.UUID, email string) error
		Delete(ctx context.Context, id uuid.UUID) error
	}

	customUserModel struct {
		db    *gorm.DB
		redis *redis.Client
	}
)

// NewUserModel returns a model for the database table.
func NewUserModel(db *gorm.DB, redisClient *redis.Client) UserModel {
	return &customUserModel{
		db:    db,
		redis: redisClient,
	}
}

// FindByID 根据用户ID查找用户
func (m *customUserModel) FindByID(ctx context.Context, id uuid.UUID) (*User, error) {
	// 尝试从缓存获取
	cacheKey := fmt.Sprintf("user:id:%s", id.String())
	if m.redis != nil {
		val, err := m.redis.Get(ctx, cacheKey).Result()
		if err == nil {
			var user User
			if err := json.Unmarshal([]byte(val), &user); err == nil {
				return &user, nil
			}
		}
	}

	// 从数据库查询
	var user User
	err := m.db.WithContext(ctx).Where("id = ?", id).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}

	// 写入缓存
	if m.redis != nil {
		data, _ := json.Marshal(user)
		m.redis.Set(ctx, cacheKey, data, 10*time.Minute)
	}

	return &user, nil
}

// FindByEmail 根据邮箱查找用户
func (m *customUserModel) FindByEmail(ctx context.Context, email string) (*User, error) {
	// 尝试从缓存获取
	cacheKey := fmt.Sprintf("user:email:%s", email)
	if m.redis != nil {
		val, err := m.redis.Get(ctx, cacheKey).Result()
		if err == nil {
			var user User
			if err := json.Unmarshal([]byte(val), &user); err == nil {
				return &user, nil
			}
		}
	}

	// 从数据库查询
	var user User
	err := m.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}

	// 写入缓存
	if m.redis != nil {
		data, _ := json.Marshal(user)
		m.redis.Set(ctx, cacheKey, data, 10*time.Minute)
		// 同时缓存 ID 索引
		idCacheKey := fmt.Sprintf("user:id:%s", user.ID.String())
		m.redis.Set(ctx, idCacheKey, data, 10*time.Minute)
	}

	return &user, nil
}

// Create 创建新用户
func (m *customUserModel) Create(ctx context.Context, user *User) error {
	err := m.db.WithContext(ctx).Create(user).Error
	if err != nil {
		return err
	}

	// 写入缓存
	if m.redis != nil {
		data, _ := json.Marshal(user)
		idCacheKey := fmt.Sprintf("user:id:%s", user.ID.String())
		emailCacheKey := fmt.Sprintf("user:email:%s", user.Email)
		m.redis.Set(ctx, idCacheKey, data, 10*time.Minute)
		m.redis.Set(ctx, emailCacheKey, data, 10*time.Minute)
	}

	return nil
}

// Update 更新用户信息
func (m *customUserModel) Update(ctx context.Context, user *User) error {
	// 先获取旧的邮箱用于清除缓存
	var oldUser User
	if err := m.db.WithContext(ctx).Where("id = ?", user.ID).First(&oldUser).Error; err != nil {
		return err
	}

	err := m.db.WithContext(ctx).Save(user).Error
	if err != nil {
		return err
	}

	// 清除缓存
	if m.redis != nil {
		idCacheKey := fmt.Sprintf("user:id:%s", user.ID.String())
		oldEmailCacheKey := fmt.Sprintf("user:email:%s", oldUser.Email)
		newEmailCacheKey := fmt.Sprintf("user:email:%s", user.Email)
		m.redis.Del(ctx, idCacheKey, oldEmailCacheKey, newEmailCacheKey)
	}

	return nil
}

// UpdatePassword 更新用户密码
func (m *customUserModel) UpdatePassword(ctx context.Context, id uuid.UUID, password string) error {
	err := m.db.WithContext(ctx).Model(&User{}).Where("id = ?", id).Update("password", password).Error
	if err != nil {
		return err
	}

	// 清除缓存
	if m.redis != nil {
		cacheKey := fmt.Sprintf("user:id:%s", id.String())
		m.redis.Del(ctx, cacheKey)
	}

	return nil
}

// UpdateEmail 更新用户邮箱
func (m *customUserModel) UpdateEmail(ctx context.Context, id uuid.UUID, email string) error {
	// 先获取旧的邮箱用于清除缓存
	var oldUser User
	if err := m.db.WithContext(ctx).Where("id = ?", id).First(&oldUser).Error; err != nil {
		return err
	}

	err := m.db.WithContext(ctx).Model(&User{}).Where("id = ?", id).Update("email", email).Error
	if err != nil {
		return err
	}

	// 清除缓存
	if m.redis != nil {
		idCacheKey := fmt.Sprintf("user:id:%s", id.String())
		oldEmailCacheKey := fmt.Sprintf("user:email:%s", oldUser.Email)
		newEmailCacheKey := fmt.Sprintf("user:email:%s", email)
		m.redis.Del(ctx, idCacheKey, oldEmailCacheKey, newEmailCacheKey)
	}

	return nil
}

// Delete 软删除用户
func (m *customUserModel) Delete(ctx context.Context, id uuid.UUID) error {
	// 先获取用户信息用于清除缓存
	var user User
	if err := m.db.WithContext(ctx).Where("id = ?", id).First(&user).Error; err != nil {
		return err
	}

	err := m.db.WithContext(ctx).Delete(&User{}, id).Error
	if err != nil {
		return err
	}

	// 清除缓存
	if m.redis != nil {
		idCacheKey := fmt.Sprintf("user:id:%s", id.String())
		emailCacheKey := fmt.Sprintf("user:email:%s", user.Email)
		m.redis.Del(ctx, idCacheKey, emailCacheKey)
	}

	return nil
}
