package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User GORM 模型定义
type User struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Username  string         `gorm:"type:varchar(100);not null" json:"username"`
	Email     string         `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
	Password  string         `gorm:"type:varchar(255);not null" json:"-"`
	Gender    *string        `gorm:"type:varchar(10)" json:"gender,omitempty"`
	Avatar    *string        `gorm:"type:varchar(500)" json:"avatar,omitempty"`
	Status    string         `gorm:"type:varchar(20);default:'active';not null" json:"status"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (User) TableName() string {
	return "user"
}
