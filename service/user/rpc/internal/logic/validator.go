package logic

import (
	"errors"
	"regexp"
	"strings"
)

var (
	// emailRegex 邮箱格式验证正则表达式
	emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
)

// ValidateEmail 验证邮箱格式
func ValidateEmail(email string) bool {
	if email == "" {
		return false
	}
	return emailRegex.MatchString(email)
}

// ValidatePassword 验证密码强度
func ValidatePassword(password string) error {
	if password == "" {
		return errors.New("password cannot be empty")
	}
	if len(password) < 6 {
		return errors.New("password must be at least 6 characters")
	}
	if len(password) > 128 {
		return errors.New("password is too long")
	}
	return nil
}

// ValidateUsername 验证用户名
func ValidateUsername(username string) error {
	username = strings.TrimSpace(username)
	if username == "" {
		return errors.New("username cannot be empty")
	}
	if len(username) < 2 {
		return errors.New("username must be at least 2 characters")
	}
	if len(username) > 50 {
		return errors.New("username is too long")
	}
	return nil
}
