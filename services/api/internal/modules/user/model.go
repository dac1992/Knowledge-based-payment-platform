package user

import (
	"time"

	"gorm.io/gorm"
)

// 角色常量。Phase 2 升级为 RBAC 多角色模型，此处先行约定取值。
const (
	RoleAdmin    = "admin"    // 平台管理员
	RoleLecturer = "lecturer" // 讲师
	RoleStudent  = "student"  // 学员（默认）
)

// 账号状态。
const (
	StatusDisabled int8 = 0
	StatusActive   int8 = 1
)

// User 用户模型，对应 users 表（由本模块独占，其他模块必须通过 Service 访问）。
type User struct {
	ID           uint64         `gorm:"primaryKey;autoIncrement"`
	Username     string         `gorm:"size:64;not null"`
	Phone        string         `gorm:"size:20;not null;default:''"`
	Email        string         `gorm:"size:128;not null;default:''"`
	PasswordHash string         `gorm:"size:255;not null"`
	Nickname     string         `gorm:"size:64;not null;default:''"`
	Avatar       string         `gorm:"size:512;not null;default:''"`
	Role         string         `gorm:"size:32;not null;default:student;index"`
	Status       int8           `gorm:"not null;default:1"`
	LastLoginAt  *time.Time     `gorm:"column:last_login_at"`
	CreatedAt    time.Time      `gorm:"column:created_at"`
	UpdatedAt    time.Time      `gorm:"column:updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

// TableName 固定表名，避免 GORM 复数化带来的歧义。
func (User) TableName() string { return "users" }

// IsActive 账号是否可用。
func (u *User) IsActive() bool { return u.Status == StatusActive }

// DisplayName 用于日志与展示的友好名称。
func (u *User) DisplayName() string {
	if u.Nickname != "" {
		return u.Nickname
	}
	return u.Username
}
