package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Article struct {
	ID        string         `json:"id" gorm:"type:varchar(36);primaryKey"`
	Title     string         `json:"title" gorm:"not null;size:255"`
	Content   string         `json:"content" gorm:"type:longtext"`
	Slug      string         `json:"slug" gorm:"unique;not null;size:255"`
	Status    string         `json:"status" gorm:"default:'draft';size:20"`
	ExpiresAt *time.Time     `json:"expires_at"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// BeforeCreate 在创建前自动生成UUID
func (a *Article) BeforeCreate(tx *gorm.DB) error {
	if a.ID == "" {
		a.ID = uuid.New().String()
	}
	return nil
}

type User struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Username  string         `json:"username" gorm:"unique;not null;size:100"`
	Email     string         `json:"email" gorm:"unique;not null;size:255"`
	Password  string         `json:"-" gorm:"not null"`
	Role      string         `json:"role" gorm:"default:'user';size:20"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type APIKey struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" gorm:"not null;size:100"`
	Key         string         `json:"key" gorm:"unique;not null;size:255"`
	IsActive    bool           `json:"is_active" gorm:"default:true"`
	LastUsedAt  *time.Time     `json:"last_used_at"`
	ExpiresAt   *time.Time     `json:"expires_at"`
	Permissions string         `json:"permissions" gorm:"type:text"` // JSON格式存储权限
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

type Certificate struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Domain    string         `json:"domain" gorm:"unique;not null;size:255"`
	CertPath  string         `json:"cert_path" gorm:"not null"`
	KeyPath   string         `json:"key_path" gorm:"not null"`
	ExpiresAt time.Time      `json:"expires_at"`
	AutoRenew bool           `json:"auto_renew" gorm:"default:true"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}
