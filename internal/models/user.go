package models

import (
	"errors"
	"regexp"
	"strings"
)

// User 定义了用户的数据结构
type User struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	Mail string `json:"mail"`
}

// Validate 验证用户数据是否有效
func (u *User) Validate() error {
	// 验证名称
	if strings.TrimSpace(u.Name) == "" {
		return errors.New("name is required")
	}

	if len(u.Name) > 255 {
		return errors.New("name is too long (maximum 255 characters)")
	}

	// 验证邮箱
	if strings.TrimSpace(u.Mail) == "" {
		return errors.New("email is required")
	}

	// 简单的邮箱格式验证
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(u.Mail) {
		return errors.New("invalid email format")
	}

	if len(u.Mail) > 255 {
		return errors.New("email is too long (maximum 255 characters)")
	}

	return nil
}
