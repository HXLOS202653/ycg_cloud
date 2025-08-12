// Package model defines permission and role related data models
package model

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

// Role represents a user role in the system
type Role struct {
	ID          int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	Name        string    `json:"name" gorm:"type:varchar(50);uniqueIndex;not null"`
	DisplayName string    `json:"displayName" gorm:"type:varchar(100);not null;column:display_name"`
	Description string    `json:"description" gorm:"type:text;column:description"`
	IsSystem    bool      `json:"isSystem" gorm:"default:false;column:is_system"`
	IsActive    bool      `json:"isActive" gorm:"default:true;column:is_active"`
	Priority    int       `json:"priority" gorm:"default:0;column:priority"` // Higher priority = more permissions
	CreatedAt   time.Time `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updatedAt" gorm:"autoUpdateTime"`

	// Relationships
	Permissions []Permission `json:"permissions,omitempty" gorm:"many2many:role_permissions;"`
	Users       []User       `json:"users,omitempty" gorm:"foreignKey:Role;references:Name"`
}

// TableName returns the table name for the Role model
func (Role) TableName() string {
	return "roles"
}

// Permission represents a permission in the system
type Permission struct {
	ID          int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	Name        string    `json:"name" gorm:"type:varchar(100);uniqueIndex;not null"`
	Resource    string    `json:"resource" gorm:"type:varchar(50);not null;index"`
	Action      string    `json:"action" gorm:"type:varchar(50);not null;index"`
	Effect      string    `json:"effect" gorm:"type:varchar(10);not null;default:'allow'"` // allow, deny
	Description string    `json:"description" gorm:"type:text"`
	IsSystem    bool      `json:"isSystem" gorm:"default:false;column:is_system"`
	CreatedAt   time.Time `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updatedAt" gorm:"autoUpdateTime"`

	// Relationships
	Roles []Role `json:"roles,omitempty" gorm:"many2many:role_permissions;"`
}

// TableName returns the table name for the Permission model
func (Permission) TableName() string {
	return "permissions"
}

// RolePermission represents the many-to-many relationship between roles and permissions
type RolePermission struct {
	ID           int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	RoleID       int64     `json:"roleId" gorm:"not null;index;column:role_id"`
	PermissionID int64     `json:"permissionId" gorm:"not null;index;column:permission_id"`
	CreatedAt    time.Time `json:"createdAt" gorm:"autoCreateTime"`

	// Relationships
	Role       Role       `json:"role,omitempty" gorm:"foreignKey:RoleID;references:ID"`
	Permission Permission `json:"permission,omitempty" gorm:"foreignKey:PermissionID;references:ID"`
}

// TableName returns the table name for the RolePermission model
func (RolePermission) TableName() string {
	return "role_permissions"
}

// UserPermission represents user-specific permissions that override role permissions
type UserPermission struct {
	ID           int64      `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID       int64      `json:"userId" gorm:"not null;index;column:user_id"`
	PermissionID int64      `json:"permissionId" gorm:"not null;index;column:permission_id"`
	ResourceID   *string    `json:"resourceId,omitempty" gorm:"type:varchar(255);column:resource_id"` // Specific resource instance
	Effect       string     `json:"effect" gorm:"type:varchar(10);not null;default:'allow'"`          // allow, deny
	ExpiresAt    *time.Time `json:"expiresAt,omitempty" gorm:"column:expires_at"`
	CreatedAt    time.Time  `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt    time.Time  `json:"updatedAt" gorm:"autoUpdateTime"`

	// Relationships
	User       User       `json:"user,omitempty" gorm:"foreignKey:UserID;references:ID"`
	Permission Permission `json:"permission,omitempty" gorm:"foreignKey:PermissionID;references:ID"`
}

// TableName returns the table name for the UserPermission model
func (UserPermission) TableName() string {
	return "user_permissions"
}

// IsExpired checks if the user permission has expired
func (up *UserPermission) IsExpired() bool {
	return up.ExpiresAt != nil && time.Now().After(*up.ExpiresAt)
}

// IsValid checks if the user permission is valid and not expired
func (up *UserPermission) IsValid() bool {
	return !up.IsExpired()
}

// BeforeCreate sets default values before creating a role
func (r *Role) BeforeCreate(_ *gorm.DB) error {
	if r.Name == "" {
		return gorm.ErrInvalidData
	}
	if r.DisplayName == "" {
		r.DisplayName = r.Name
	}
	return nil
}

// BeforeCreate sets default values before creating a permission
func (p *Permission) BeforeCreate(_ *gorm.DB) error {
	if p.Name == "" || p.Resource == "" || p.Action == "" {
		return gorm.ErrInvalidData
	}
	if p.Effect == "" {
		p.Effect = "allow"
	}
	return nil
}

// System-defined roles
const (
	RoleSuperAdmin = "super_admin"
	RoleAdmin      = "admin"
	RoleUser       = "user"
	RoleGuest      = "guest"
)

// System-defined resources
const (
	ResourceUser         = "user"
	ResourceFile         = "file"
	ResourceFolder       = "folder"
	ResourceTeam         = "team"
	ResourceSystem       = "system"
	ResourceNotification = "notification"
	ResourceLog          = "log"
	ResourceShare        = "share"
)

// System-defined actions
const (
	ActionCreate = "create"
	ActionRead   = "read"
	ActionUpdate = "update"
	ActionDelete = "delete"
	ActionShare  = "share"
	ActionAdmin  = "admin"
	ActionAll    = "*"
)

// Permission effects
const (
	EffectAllow = "allow"
	EffectDeny  = "deny"
)

// GetSystemRoles returns predefined system roles
func GetSystemRoles() []Role {
	return []Role{
		{
			Name:        RoleSuperAdmin,
			DisplayName: "超级管理员",
			Description: "拥有系统所有权限",
			IsSystem:    true,
			Priority:    1000,
		},
		{
			Name:        RoleAdmin,
			DisplayName: "管理员",
			Description: "拥有大部分管理权限",
			IsSystem:    true,
			Priority:    500,
		},
		{
			Name:        RoleUser,
			DisplayName: "普通用户",
			Description: "拥有基本用户权限",
			IsSystem:    true,
			Priority:    100,
		},
		{
			Name:        RoleGuest,
			DisplayName: "访客",
			Description: "只读权限",
			IsSystem:    true,
			Priority:    10,
		},
	}
}

// GetSystemPermissions returns predefined system permissions
func GetSystemPermissions() []Permission {
	permissions := []Permission{}

	resources := []string{ResourceUser, ResourceFile, ResourceFolder, ResourceTeam, ResourceSystem, ResourceNotification, ResourceLog, ResourceShare}
	actions := []string{ActionCreate, ActionRead, ActionUpdate, ActionDelete, ActionShare, ActionAdmin}

	for _, resource := range resources {
		for _, action := range actions {
			// Skip invalid combinations
			if resource == ResourceSystem && action == ActionShare {
				continue
			}
			if resource == ResourceLog && (action == ActionCreate || action == ActionUpdate || action == ActionDelete || action == ActionShare) {
				continue
			}

			permission := Permission{
				Name:        fmt.Sprintf("%s:%s", resource, action),
				Resource:    resource,
				Action:      action,
				Effect:      EffectAllow,
				Description: fmt.Sprintf("Permission to %s %s", action, resource),
				IsSystem:    true,
			}
			permissions = append(permissions, permission)
		}

		// Add wildcard permission for each resource
		permission := Permission{
			Name:        fmt.Sprintf("%s:*", resource),
			Resource:    resource,
			Action:      ActionAll,
			Effect:      EffectAllow,
			Description: fmt.Sprintf("Full access to %s", resource),
			IsSystem:    true,
		}
		permissions = append(permissions, permission)
	}

	// Add global admin permission
	permissions = append(permissions, Permission{
		Name:        "system:*",
		Resource:    ResourceSystem,
		Action:      ActionAll,
		Effect:      EffectAllow,
		Description: "Full system access",
		IsSystem:    true,
	})

	return permissions
}
