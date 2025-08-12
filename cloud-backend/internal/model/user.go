// Package model defines data models for the application
package model

import (
	"time"

	"gorm.io/gorm"
)

// User represents the user model in the database - 符合字段映射规范v2.0
type User struct {
	// 基础标识字段
	ID            int64   `json:"id" gorm:"primaryKey;autoIncrement;column:id"`
	Username      string  `json:"username" gorm:"type:varchar(50);uniqueIndex;not null;column:username"`
	Email         string  `json:"email" gorm:"type:varchar(100);uniqueIndex;not null;column:email"`
	EmailVerified bool    `json:"email_verified" gorm:"default:false;column:email_verified"`
	Phone         *string `json:"phone" gorm:"type:varchar(20);column:phone"`
	PhoneVerified bool    `json:"phone_verified" gorm:"default:false;column:phone_verified"`
	PasswordHash  string  `json:"-" gorm:"type:varchar(255);not null;column:password_hash"` // 敏感字段不返回

	// 用户信息字段
	RealName  *string `json:"real_name" gorm:"type:varchar(100);column:real_name"`
	AvatarURL *string `json:"avatar_url" gorm:"type:varchar(500);column:avatar_url"`
	Role      string  `json:"role" gorm:"type:enum('user','vip','admin','super_admin');default:'user';not null;column:role"`
	Status    string  `json:"status" gorm:"type:enum('pending','active','disabled','banned');default:'pending';not null;column:status"`

	// 存储配置字段
	StorageQuota           int64   `json:"storage_quota" gorm:"default:10737418240;column:storage_quota"`
	StorageUsed            int64   `json:"storage_used" gorm:"default:0;column:storage_used"`
	UploadBandwidthLimit   int     `json:"upload_bandwidth_limit" gorm:"default:10485760;column:upload_bandwidth_limit"`
	DownloadBandwidthLimit int     `json:"download_bandwidth_limit" gorm:"default:10485760;column:download_bandwidth_limit"`
	MaxFileSize            int64   `json:"max_file_size" gorm:"default:10737418240;column:max_file_size"`
	AllowedFileTypes       *string `json:"allowed_file_types" gorm:"type:json;column:allowed_file_types"`
	ForbiddenFileTypes     *string `json:"forbidden_file_types" gorm:"type:json;default:'[\"exe\", \"bat\", \"com\", \"scr\", \"pif\"]';column:forbidden_file_types"`

	// 安全相关字段
	LoginAttempts     uint       `json:"-" gorm:"default:0;column:login_attempts"` // 敏感字段不返回
	LockedUntil       *time.Time `json:"-" gorm:"column:locked_until"`             // 敏感字段不返回
	LastLoginAt       *time.Time `json:"last_login_at" gorm:"column:last_login_at"`
	LastLoginIP       *string    `json:"last_login_ip" gorm:"type:varchar(45);column:last_login_ip"`
	PasswordChangedAt time.Time  `json:"password_changed_at" gorm:"default:CURRENT_TIMESTAMP;column:password_changed_at"`
	TwoFactorEnabled  bool       `json:"two_factor_enabled" gorm:"default:false;column:two_factor_enabled"`
	TwoFactorSecret   *string    `json:"-" gorm:"type:varchar(32);column:two_factor_secret"` // 敏感字段不返回

	// 用户偏好设置字段
	Language           string `json:"language" gorm:"type:varchar(10);default:'zh-CN';column:language"`
	Timezone           string `json:"timezone" gorm:"type:varchar(50);default:'Asia/Shanghai';column:timezone"`
	Theme              string `json:"theme" gorm:"type:enum('light','dark','auto');default:'auto';column:theme"`
	EmailNotifications bool   `json:"email_notifications" gorm:"default:true;column:email_notifications"`
	SMSNotifications   bool   `json:"sms_notifications" gorm:"default:false;column:sms_notifications"`

	// 审计字段
	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime;column:created_at"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"autoUpdateTime;column:updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index;column:deleted_at"` // 软删除字段
}

// TableName returns the table name for the User model
func (User) TableName() string {
	return "users"
}

// BeforeCreate hooks runs before creating a user
func (u *User) BeforeCreate(_ *gorm.DB) error {
	// Set default values if not provided
	if u.StorageQuota == 0 {
		u.StorageQuota = 10737418240 // 10GB
	}
	if u.Status == "" {
		u.Status = "active"
	}
	if u.Role == "" {
		u.Role = "user"
	}
	if u.Language == "" {
		u.Language = "zh-CN"
	}
	if u.Timezone == "" {
		u.Timezone = "Asia/Shanghai"
	}
	if u.Theme == "" {
		u.Theme = "light"
	}
	return nil
}

// IsActive checks if the user account is active
func (u *User) IsActive() bool {
	return u.Status == "active"
}

// IsLocked checks if the user account is temporarily locked
func (u *User) IsLocked() bool {
	return u.LockedUntil != nil && u.LockedUntil.After(time.Now())
}

// CanLogin checks if the user can log in
func (u *User) CanLogin() bool {
	return u.IsActive() && !u.IsLocked()
}

// GetStorageUsagePercentage calculates storage usage percentage
func (u *User) GetStorageUsagePercentage() float64 {
	if u.StorageQuota == 0 {
		return 0
	}
	return float64(u.StorageUsed) / float64(u.StorageQuota) * 100
}

// HasStorageSpace checks if user has enough storage space
func (u *User) HasStorageSpace(requiredBytes int64) bool {
	return u.StorageUsed+requiredBytes <= u.StorageQuota
}

// IncrementLoginAttempts increments failed login attempts
func (u *User) IncrementLoginAttempts() {
	u.LoginAttempts++
	// Lock account for 15 minutes after 5 failed attempts
	if u.LoginAttempts >= 5 {
		lockUntil := time.Now().Add(15 * time.Minute)
		u.LockedUntil = &lockUntil
	}
}

// ResetLoginAttempts resets failed login attempts
func (u *User) ResetLoginAttempts() {
	u.LoginAttempts = 0
	u.LockedUntil = nil
}

// UpdateLastLogin updates last login information
func (u *User) UpdateLastLogin(ip string) {
	now := time.Now()
	u.LastLoginAt = &now
	u.LastLoginIP = &ip
	u.ResetLoginAttempts()
}

// UserSession represents user session information
type UserSession struct {
	ID           int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID       int64     `json:"userId" gorm:"not null;index;column:user_id"`
	SessionToken string    `json:"-" gorm:"type:varchar(255);uniqueIndex;not null;column:session_token"` // Hidden from JSON
	RefreshToken string    `json:"-" gorm:"type:varchar(255);uniqueIndex;not null;column:refresh_token"` // Hidden from JSON
	DeviceInfo   string    `json:"deviceInfo" gorm:"type:text;column:device_info"`
	IPAddress    string    `json:"ipAddress" gorm:"type:varchar(45);column:ip_address"`
	UserAgent    string    `json:"userAgent" gorm:"type:text;column:user_agent"`
	IsActive     bool      `json:"isActive" gorm:"default:true;column:is_active"`
	ExpiresAt    time.Time `json:"expiresAt" gorm:"not null;column:expires_at"`
	CreatedAt    time.Time `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt    time.Time `json:"updatedAt" gorm:"autoUpdateTime"`

	// Relationship
	User User `json:"user,omitempty" gorm:"foreignKey:UserID;references:ID"`
}

// TableName returns the table name for the UserSession model
func (UserSession) TableName() string {
	return "user_sessions"
}

// IsExpired checks if the session has expired
func (s *UserSession) IsExpired() bool {
	return time.Now().After(s.ExpiresAt)
}

// IsValid checks if the session is valid and active
func (s *UserSession) IsValid() bool {
	return s.IsActive && !s.IsExpired()
}
