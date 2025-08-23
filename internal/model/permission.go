package model

import (
	"time"

	"gorm.io/gorm"
)

// PermissionAction 权限操作枚举
type PermissionAction string

const (
	// 文件权限
	PermissionRead     PermissionAction = "read"     // 读取
	PermissionWrite    PermissionAction = "write"    // 写入
	PermissionDelete   PermissionAction = "delete"   // 删除
	PermissionShare    PermissionAction = "share"    // 分享
	PermissionDownload PermissionAction = "download" // 下载
	PermissionUpload   PermissionAction = "upload"   // 上传
	PermissionPreview  PermissionAction = "preview"  // 预览

	// 系统权限
	PermissionUserManage    PermissionAction = "user_manage"    // 用户管理
	PermissionSystemConfig  PermissionAction = "system_config"  // 系统配置
	PermissionLogView       PermissionAction = "log_view"       // 日志查看
	PermissionTeamManage    PermissionAction = "team_manage"    // 团队管理
	PermissionStorageManage PermissionAction = "storage_manage" // 存储管理
)

// ResourceType 资源类型枚举
type ResourceType string

const (
	ResourceTypeFile   ResourceType = "file"   // 文件
	ResourceTypeFolder ResourceType = "folder" // 文件夹
	ResourceTypeTeam   ResourceType = "team"   // 团队
	ResourceTypeSystem ResourceType = "system" // 系统
)

// PermissionTemplate 权限模板
type PermissionTemplate struct {
	ID          uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string `gorm:"type:varchar(100);not null;uniqueIndex" json:"name" validate:"required"`
	Description string `gorm:"type:text" json:"description"`
	IsDefault   bool   `gorm:"default:false;index" json:"is_default"`
	IsSystem    bool   `gorm:"default:false;index" json:"is_system"`

	// 存储配额设置
	StorageQuota int64 `gorm:"default:5368709120;comment:存储配额(字节)" json:"storage_quota"`

	// 功能权限设置(JSON格式存储)
	Permissions string `gorm:"type:text;comment:权限配置(JSON)" json:"permissions"`

	// 时间戳
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// 关联关系
	Users               []User               `gorm:"foreignKey:PermissionTemplateID;constraint:OnDelete:SET NULL" json:"users,omitempty"`
	TemplatePermissions []templatePermission `gorm:"foreignKey:TemplateID;constraint:OnDelete:CASCADE" json:"template_permissions,omitempty"`
}

// TableName 指定表名
func (PermissionTemplate) TableName() string {
	return "permission_templates"
}

// templatePermission 模板权限详情 (私有)
type templatePermission struct {
	ID           uint               `gorm:"primaryKey;autoIncrement" json:"id"`
	TemplateID   uint               `gorm:"not null;index" json:"template_id"`
	Template     PermissionTemplate `gorm:"foreignKey:TemplateID;constraint:OnDelete:CASCADE" json:"template,omitempty"`
	ResourceType ResourceType       `gorm:"type:varchar(20);not null;index" json:"resource_type"`
	Action       PermissionAction   `gorm:"type:varchar(50);not null;index" json:"action"`
	Allowed      bool               `gorm:"default:false" json:"allowed"`

	// 时间戳
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// TableName 指定表名
func (templatePermission) TableName() string {
	return "template_permissions"
}

// userPermission 用户权限 (私有)
type userPermission struct {
	ID           uint             `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID       uint             `gorm:"not null;index" json:"user_id"`
	User         User             `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"user,omitempty"`
	ResourceType ResourceType     `gorm:"type:varchar(20);not null;index" json:"resource_type"`
	ResourceID   *uint            `gorm:"index;comment:资源ID(可为空表示全局权限)" json:"resource_id"`
	Action       PermissionAction `gorm:"type:varchar(50);not null;index" json:"action"`
	Allowed      bool             `gorm:"default:false" json:"allowed"`

	// 权限来源
	GrantedBy *uint     `gorm:"index;comment:授权人ID" json:"granted_by"`
	Granter   *User     `gorm:"foreignKey:GrantedBy" json:"granter,omitempty"`
	GrantedAt time.Time `gorm:"autoCreateTime" json:"granted_at"`

	// 过期时间
	ExpiresAt *time.Time `gorm:"index;comment:权限过期时间" json:"expires_at"`

	// 时间戳
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

// TableName 指定表名
func (userPermission) TableName() string {
	return "user_permissions"
}

// filePermission 文件权限 (私有)
type filePermission struct {
	ID      uint             `gorm:"primaryKey;autoIncrement" json:"id"`
	FileID  uint             `gorm:"not null;index" json:"file_id"`
	File    File             `gorm:"foreignKey:FileID;constraint:OnDelete:CASCADE" json:"file,omitempty"`
	UserID  *uint            `gorm:"index;comment:用户ID" json:"user_id"`
	User    *User            `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"user,omitempty"`
	TeamID  *uint            `gorm:"index;comment:团队ID" json:"team_id"`
	Team    *Team            `gorm:"foreignKey:TeamID;constraint:OnDelete:CASCADE" json:"team,omitempty"`
	Action  PermissionAction `gorm:"type:varchar(50);not null;index" json:"action"`
	Allowed bool             `gorm:"default:false" json:"allowed"`

	// 权限来源
	GrantedBy *uint     `gorm:"index;comment:授权人ID" json:"granted_by"`
	Granter   *User     `gorm:"foreignKey:GrantedBy" json:"granter,omitempty"`
	GrantedAt time.Time `gorm:"autoCreateTime" json:"granted_at"`

	// 过期时间
	ExpiresAt *time.Time `gorm:"index;comment:权限过期时间" json:"expires_at"`

	// 时间戳
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

// TableName 指定表名
func (filePermission) TableName() string {
	return "file_permissions"
}

// Role 角色定义
type Role struct {
	ID          uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string `gorm:"type:varchar(100);not null;uniqueIndex" json:"name" validate:"required"`
	Description string `gorm:"type:text" json:"description"`
	IsSystem    bool   `gorm:"default:false;index" json:"is_system"`
	Level       int    `gorm:"default:0;index;comment:角色级别" json:"level"`

	// 权限配置(JSON格式)
	Permissions string `gorm:"type:text;comment:角色权限(JSON)" json:"permissions"`

	// 时间戳
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// 关联关系
	UserRoles []userRole `gorm:"foreignKey:RoleID;constraint:OnDelete:CASCADE" json:"user_roles,omitempty"`
	TeamRoles []TeamRole `gorm:"foreignKey:RoleID" json:"team_roles,omitempty"`
}

// TableName 指定表名
func (Role) TableName() string {
	return "roles"
}

// userRole 用户角色关联
type userRole struct {
	ID     uint `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID uint `gorm:"not null;index" json:"user_id"`
	User   User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"user,omitempty"`
	RoleID uint `gorm:"not null;index" json:"role_id"`
	Role   Role `gorm:"foreignKey:RoleID;constraint:OnDelete:CASCADE" json:"role,omitempty"`

	// 权限来源
	GrantedBy *uint     `gorm:"index;comment:授权人ID" json:"granted_by"`
	Granter   *User     `gorm:"foreignKey:GrantedBy;constraint:OnDelete:RESTRICT" json:"granter,omitempty"`
	GrantedAt time.Time `gorm:"autoCreateTime" json:"granted_at"`

	// 过期时间
	ExpiresAt *time.Time `gorm:"index;comment:角色过期时间" json:"expires_at"`

	// 时间戳
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

// TableName 指定表名
func (userRole) TableName() string {
	return "user_roles"
}

// IsExpired 检查权限是否过期
func (up *userPermission) IsExpired() bool {
	return up.ExpiresAt != nil && up.ExpiresAt.Before(time.Now())
}

// IsExpired 检查文件权限是否过期
func (fp *filePermission) IsExpired() bool {
	return fp.ExpiresAt != nil && fp.ExpiresAt.Before(time.Now())
}

// IsExpired 检查角色是否过期
func (ur *userRole) IsExpired() bool {
	return ur.ExpiresAt != nil && ur.ExpiresAt.Before(time.Now())
}
