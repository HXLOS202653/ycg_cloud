package model

import (
	"time"

	"gorm.io/gorm"
)

// UserStatus 用户状态枚举
type UserStatus string

const (
	UserStatusActive    UserStatus = "active"    // 活跃
	UserStatusInactive  UserStatus = "inactive"  // 非活跃
	UserStatusSuspended UserStatus = "suspended" // 暂停
	UserStatusDeleted   UserStatus = "deleted"   // 删除
)

// UserType 用户类型枚举
type UserType string

const (
	UserTypeNormal UserType = "normal" // 普通用户
	UserTypeAdmin  UserType = "admin"  // 管理员
)

// User 用户模型
type User struct {
	ID           uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	Username     string     `gorm:"type:varchar(50);uniqueIndex;not null" json:"username"`
	Email        string     `gorm:"type:varchar(100);uniqueIndex;not null" json:"email"`
	PasswordHash string     `gorm:"type:varchar(255);not null" json:"-"`
	Nickname     string     `gorm:"type:varchar(100)" json:"nickname"`
	Avatar       string     `gorm:"type:varchar(500)" json:"avatar"`
	Phone        string     `gorm:"type:varchar(20);index" json:"phone"`
	UserType     UserType   `gorm:"type:varchar(20);default:'normal';index" json:"user_type"`
	Status       UserStatus `gorm:"type:varchar(20);default:'active';index" json:"status"`

	// 存储配额相关
	StorageQuota int64 `gorm:"default:5368709120;comment:存储配额(字节)" json:"storage_quota"` // 默认5GB
	UsedStorage  int64 `gorm:"default:0;comment:已使用存储(字节)" json:"used_storage"`

	// 权限模板关联
	PermissionTemplateID *uint               `gorm:"index;comment:权限模板ID" json:"permission_template_id"`
	PermissionTemplate   *PermissionTemplate `gorm:"foreignKey:PermissionTemplateID" json:"permission_template"`

	// 安全相关
	LastLoginAt    *time.Time `gorm:"comment:最后登录时间" json:"last_login_at"`
	LastLoginIP    string     `gorm:"type:varchar(45);comment:最后登录IP" json:"last_login_ip"`
	LoginFailCount int        `gorm:"default:0;comment:登录失败次数" json:"login_fail_count"`
	LockedUntil    *time.Time `gorm:"comment:锁定到期时间" json:"locked_until"`

	// MFA相关
	MFAEnabled bool   `gorm:"default:false;comment:是否启用MFA" json:"mfa_enabled"`
	MFASecret  string `gorm:"type:varchar(255);comment:MFA密钥" json:"-"`

	// 时间戳
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// 关联关系
	Files               []File               `gorm:"foreignKey:OwnerID" json:"files,omitempty"`
	TeamMembers         []TeamMember         `gorm:"foreignKey:UserID" json:"team_members,omitempty"`
	UserPermissions     []userPermission     `gorm:"foreignKey:UserID" json:"user_permissions,omitempty"`
	OperationLogs       []OperationLog       `gorm:"foreignKey:UserID" json:"operation_logs,omitempty"`
	Messages            []Message            `gorm:"foreignKey:SenderID" json:"messages,omitempty"`
	ConversationMembers []ConversationMember `gorm:"foreignKey:UserID" json:"conversation_members,omitempty"`
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}

// BeforeCreate GORM钩子：创建前
func (u *User) BeforeCreate(tx *gorm.DB) error {
	// 设置默认值
	if u.UserType == "" {
		u.UserType = UserTypeNormal
	}
	if u.Status == "" {
		u.Status = UserStatusActive
	}
	if u.StorageQuota == 0 {
		u.StorageQuota = 5368709120 // 5GB
	}
	return nil
}

// IsActive 检查用户是否活跃
func (u *User) IsActive() bool {
	return u.Status == UserStatusActive
}

// IsAdmin 检查用户是否为管理员
func (u *User) IsAdmin() bool {
	return u.UserType == UserTypeAdmin
}

// IsStorageExceeded 检查存储是否超限
func (u *User) IsStorageExceeded() bool {
	return u.UsedStorage >= u.StorageQuota
}

// GetAvailableStorage 获取可用存储空间
func (u *User) GetAvailableStorage() int64 {
	return u.StorageQuota - u.UsedStorage
}

// IsLocked 检查用户是否被锁定
func (u *User) IsLocked() bool {
	return u.LockedUntil != nil && u.LockedUntil.After(time.Now())
}
