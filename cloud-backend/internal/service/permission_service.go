// Package service provides permission and authorization services
package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"

	"github.com/HXLOS202653/ycg_cloud/cloud-backend/internal/config"
	"github.com/HXLOS202653/ycg_cloud/cloud-backend/internal/model"
	"github.com/HXLOS202653/ycg_cloud/cloud-backend/internal/pkg/database"
)

// PermissionService handles permission and authorization operations
type PermissionService struct {
	config       *config.Config
	db           *gorm.DB
	redisManager *database.RedisManager
}

// NewPermissionService creates a new permission service
func NewPermissionService(cfg *config.Config, db *gorm.DB, redis *database.RedisManager) *PermissionService {
	return &PermissionService{
		config:       cfg,
		db:           db,
		redisManager: redis,
	}
}

// PermissionCheckResult represents the result of permission checking
type PermissionCheckResult struct {
	Allowed     bool   `json:"allowed"`
	Reason      string `json:"reason"`
	MatchedRule string `json:"matchedRule,omitempty"`
}

// CheckPermission checks if a user has permission to perform an action on a resource
func (s *PermissionService) CheckPermission(userID int64, resource, action string, resourceID ...string) *PermissionCheckResult {
	// Try cache first
	cacheKey := s.getPermissionCacheKey(userID, resource, action, resourceID...)
	if result := s.getPermissionFromCache(cacheKey); result != nil {
		return result
	}

	// Perform permission check
	result := s.performPermissionCheck(userID, resource, action, resourceID...)

	// Cache result for 5 minutes
	s.cachePermissionResult(cacheKey, result, 5*time.Minute)

	return result
}

// performPermissionCheck performs the actual permission checking logic
func (s *PermissionService) performPermissionCheck(userID int64, resource, action string, resourceID ...string) *PermissionCheckResult {
	// Get user info
	var user model.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return &PermissionCheckResult{
			Allowed: false,
			Reason:  "User not found",
		}
	}

	// Check if user is super admin
	if user.Role == model.RoleSuperAdmin {
		return &PermissionCheckResult{
			Allowed:     true,
			Reason:      "Super admin access",
			MatchedRule: "super_admin",
		}
	}

	// Check user-specific permissions first (highest priority)
	if result := s.checkUserPermissions(userID, resource, action, resourceID...); result.Allowed || result.Reason == "explicitly_denied" {
		return result
	}

	// Check role-based permissions
	if result := s.checkRolePermissions(user.Role, resource, action); result.Allowed {
		return result
	}

	// Default deny
	return &PermissionCheckResult{
		Allowed: false,
		Reason:  "No permission found",
	}
}

// checkUserPermissions checks user-specific permissions
func (s *PermissionService) checkUserPermissions(userID int64, resource, action string, resourceID ...string) *PermissionCheckResult {
	var userPerms []model.UserPermission

	query := s.db.Joins("Permission").Where("user_permissions.user_id = ?", userID)

	// Build permission query
	permQuery := s.db.Where("resource = ? AND (action = ? OR action = ?)", resource, action, model.ActionAll)

	var permIDs []int64
	permQuery.Model(&model.Permission{}).Pluck("id", &permIDs)

	if len(permIDs) > 0 {
		query = query.Where("user_permissions.permission_id IN ?", permIDs)
	} else {
		// No matching permissions found
		return &PermissionCheckResult{Allowed: false, Reason: "no_user_permissions"}
	}

	// Filter by resource ID if specified
	if len(resourceID) > 0 && resourceID[0] != "" {
		query = query.Where("(user_permissions.resource_id = ? OR user_permissions.resource_id IS NULL)", resourceID[0])
	}

	if err := query.Find(&userPerms).Error; err != nil {
		return &PermissionCheckResult{Allowed: false, Reason: "query_error"}
	}

	// Evaluate user permissions (explicit deny takes precedence)
	hasAllow := false
	hasDeny := false

	for i := range userPerms {
		if !userPerms[i].IsValid() {
			continue // Skip expired permissions
		}

		switch userPerms[i].Effect {
		case model.EffectDeny:
			hasDeny = true
		case model.EffectAllow:
			hasAllow = true
		}
	}

	if hasDeny {
		return &PermissionCheckResult{
			Allowed:     false,
			Reason:      "explicitly_denied",
			MatchedRule: "user_permission_deny",
		}
	}

	if hasAllow {
		return &PermissionCheckResult{
			Allowed:     true,
			Reason:      "user_permission_allow",
			MatchedRule: "user_permission",
		}
	}

	return &PermissionCheckResult{Allowed: false, Reason: "no_user_permissions"}
}

// checkRolePermissions checks role-based permissions
func (s *PermissionService) checkRolePermissions(roleName, resource, action string) *PermissionCheckResult {
	var role model.Role
	if err := s.db.Preload("Permissions").Where("name = ?", roleName).First(&role).Error; err != nil {
		return &PermissionCheckResult{Allowed: false, Reason: "role_not_found"}
	}

	// Check role permissions
	for i := range role.Permissions {
		perm := &role.Permissions[i]
		if perm.Resource == resource && (perm.Action == action || perm.Action == model.ActionAll) {
			if perm.Effect == model.EffectAllow {
				return &PermissionCheckResult{
					Allowed:     true,
					Reason:      "role_permission",
					MatchedRule: fmt.Sprintf("role:%s", roleName),
				}
			}
		}
	}

	return &PermissionCheckResult{Allowed: false, Reason: "no_role_permissions"}
}

// HasPermission is a simplified wrapper for CheckPermission
func (s *PermissionService) HasPermission(userID int64, resource, action string, resourceID ...string) bool {
	result := s.CheckPermission(userID, resource, action, resourceID...)
	return result.Allowed
}

// GrantUserPermission grants a specific permission to a user
func (s *PermissionService) GrantUserPermission(userID, permissionID int64, resourceID *string, expiresAt *time.Time) error {
	userPerm := &model.UserPermission{
		UserID:       userID,
		PermissionID: permissionID,
		ResourceID:   resourceID,
		Effect:       model.EffectAllow,
		ExpiresAt:    expiresAt,
	}

	if err := s.db.Create(userPerm).Error; err != nil {
		return fmt.Errorf("failed to grant user permission: %w", err)
	}

	// Clear cache for this user
	s.clearUserPermissionCache(userID)

	return nil
}

// RevokeUserPermission revokes a specific permission from a user
func (s *PermissionService) RevokeUserPermission(userID, permissionID int64, resourceID *string) error {
	query := s.db.Where("user_id = ? AND permission_id = ?", userID, permissionID)

	if resourceID != nil {
		query = query.Where("resource_id = ?", *resourceID)
	} else {
		query = query.Where("resource_id IS NULL")
	}

	if err := query.Delete(&model.UserPermission{}).Error; err != nil {
		return fmt.Errorf("failed to revoke user permission: %w", err)
	}

	// Clear cache for this user
	s.clearUserPermissionCache(userID)

	return nil
}

// AssignRoleToUser assigns a role to a user
func (s *PermissionService) AssignRoleToUser(userID int64, roleName string) error {
	// Verify role exists
	var role model.Role
	if err := s.db.Where("name = ?", roleName).First(&role).Error; err != nil {
		return fmt.Errorf("role not found: %w", err)
	}

	// Update user role
	if err := s.db.Model(&model.User{}).Where("id = ?", userID).Update("role", roleName).Error; err != nil {
		return fmt.Errorf("failed to assign role: %w", err)
	}

	// Clear cache for this user
	s.clearUserPermissionCache(userID)

	return nil
}

// GetUserPermissions returns all permissions for a user (both role-based and user-specific)
func (s *PermissionService) GetUserPermissions(userID int64) ([]model.Permission, error) {
	var user model.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// Get role permissions
	var role model.Role
	rolePerms := []model.Permission{}
	if err := s.db.Preload("Permissions").Where("name = ?", user.Role).First(&role).Error; err == nil {
		rolePerms = role.Permissions
	}

	// Get user-specific permissions
	var userPerms []model.UserPermission
	s.db.Preload("Permission").Where("user_id = ?", userID).Find(&userPerms)

	// Combine permissions (remove duplicates)
	permMap := make(map[int64]model.Permission)

	// Add role permissions
	for i := range rolePerms {
		permMap[rolePerms[i].ID] = rolePerms[i]
	}

	// Add/override with user-specific permissions
	for i := range userPerms {
		if userPerms[i].IsValid() {
			permMap[userPerms[i].Permission.ID] = userPerms[i].Permission
		}
	}

	// Convert map to slice
	permissions := make([]model.Permission, 0, len(permMap))
	for _, perm := range permMap {
		permissions = append(permissions, perm)
	}

	return permissions, nil
}

// InitializeSystemPermissions creates default roles and permissions
func (s *PermissionService) InitializeSystemPermissions() error {
	// Create system permissions
	systemPerms := model.GetSystemPermissions()
	for _, perm := range systemPerms {
		var existing model.Permission
		if err := s.db.Where("name = ?", perm.Name).First(&existing).Error; errors.Is(err, gorm.ErrRecordNotFound) {
			if err := s.db.Create(&perm).Error; err != nil {
				return fmt.Errorf("failed to create permission %s: %w", perm.Name, err)
			}
		}
	}

	// Create system roles
	systemRoles := model.GetSystemRoles()
	for _, role := range systemRoles {
		var existing model.Role
		if err := s.db.Where("name = ?", role.Name).First(&existing).Error; errors.Is(err, gorm.ErrRecordNotFound) {
			if err := s.db.Create(&role).Error; err != nil {
				return fmt.Errorf("failed to create role %s: %w", role.Name, err)
			}
		}
	}

	// Assign permissions to roles
	if err := s.assignDefaultRolePermissions(); err != nil {
		return fmt.Errorf("failed to assign default role permissions: %w", err)
	}

	return nil
}

// assignDefaultRolePermissions assigns default permissions to system roles
func (s *PermissionService) assignDefaultRolePermissions() error {
	// Super Admin - all permissions
	if err := s.assignPermissionsToRole(model.RoleSuperAdmin, []string{"system:*"}); err != nil {
		return err
	}

	// Admin - most permissions except system admin
	adminPerms := []string{
		"user:*", "file:*", "folder:*", "team:*", "share:*", "notification:*", "log:read",
	}
	if err := s.assignPermissionsToRole(model.RoleAdmin, adminPerms); err != nil {
		return err
	}

	// User - basic permissions
	userPerms := []string{
		"file:create", "file:read", "file:update", "file:delete", "file:share",
		"folder:create", "folder:read", "folder:update", "folder:delete", "folder:share",
		"share:create", "share:read", "notification:read",
		"user:read", "user:update", // Can read and update own profile
	}
	if err := s.assignPermissionsToRole(model.RoleUser, userPerms); err != nil {
		return err
	}

	// Guest - read-only permissions
	guestPerms := []string{
		"file:read", "folder:read", "share:read",
	}
	return s.assignPermissionsToRole(model.RoleGuest, guestPerms)
}

// assignPermissionsToRole assigns permissions to a role
func (s *PermissionService) assignPermissionsToRole(roleName string, permissionNames []string) error {
	var role model.Role
	if err := s.db.Where("name = ?", roleName).First(&role).Error; err != nil {
		return fmt.Errorf("role %s not found: %w", roleName, err)
	}

	for _, permName := range permissionNames {
		var permission model.Permission
		if err := s.db.Where("name = ?", permName).First(&permission).Error; err != nil {
			return fmt.Errorf("permission %s not found: %w", permName, err)
		}

		// Check if association already exists
		var count int64
		s.db.Model(&model.RolePermission{}).Where("role_id = ? AND permission_id = ?", role.ID, permission.ID).Count(&count)

		if count == 0 {
			rolePermission := &model.RolePermission{
				RoleID:       role.ID,
				PermissionID: permission.ID,
			}
			if err := s.db.Create(rolePermission).Error; err != nil {
				return fmt.Errorf("failed to assign permission %s to role %s: %w", permName, roleName, err)
			}
		}
	}

	return nil
}

// Cache-related methods

// getPermissionCacheKey generates a cache key for permission checking
func (s *PermissionService) getPermissionCacheKey(userID int64, resource, action string, resourceID ...string) string {
	key := fmt.Sprintf("perm:%d:%s:%s", userID, resource, action)
	if len(resourceID) > 0 && resourceID[0] != "" {
		key += ":" + resourceID[0]
	}
	return key
}

// getPermissionFromCache retrieves permission result from cache
func (s *PermissionService) getPermissionFromCache(key string) *PermissionCheckResult {
	ctx := context.Background()
	var result PermissionCheckResult
	if err := s.redisManager.GetStruct(ctx, key, &result); err == nil {
		return &result
	}
	return nil
}

// cachePermissionResult stores permission result in cache
func (s *PermissionService) cachePermissionResult(key string, result *PermissionCheckResult, ttl time.Duration) {
	ctx := context.Background()
	if err := s.redisManager.SetStruct(ctx, key, result, ttl); err != nil {
		// Log cache error but don't fail the permission check
	}
}

// clearUserPermissionCache clears all cached permissions for a user
func (s *PermissionService) clearUserPermissionCache(userID int64) {
	ctx := context.Background()
	_ = userID // Placeholder - would be used for key pattern matching

	// This is a simplified implementation
	// In production, you might want to use Redis SCAN for better performance
	keys := []string{} // You would need to implement key scanning here
	for _, key := range keys {
		if err := s.redisManager.Del(ctx, key); err != nil {
			// Log cache deletion error but continue
		}
	}
}

// ValidateResourceAccess validates if user can access a specific resource instance
func (s *PermissionService) ValidateResourceAccess(userID int64, resourceType, action, resourceID string) bool {
	// Check basic permission first
	if !s.HasPermission(userID, resourceType, action) {
		return false
	}

	// For files and folders, check ownership or team access
	switch resourceType {
	case model.ResourceFile, model.ResourceFolder:
		return s.validateFileAccess(userID, resourceID, action)
	case model.ResourceTeam:
		return s.validateTeamAccess(userID, resourceID, action)
	default:
		return true // For other resources, basic permission is enough
	}
}

// validateFileAccess checks if user can access a specific file/folder
func (s *PermissionService) validateFileAccess(userID int64, resourceID, action string) bool {
	// This would check file ownership, sharing permissions, team access, etc.
	// Implementation depends on your file model structure
	_ = userID
	_ = resourceID
	_ = action

	// Simplified implementation - check if user owns the file or has team access
	// You would implement this based on your actual file model
	return true // Placeholder
}

// validateTeamAccess checks if user can access a specific team
func (s *PermissionService) validateTeamAccess(userID int64, resourceID, action string) bool {
	// This would check team membership, roles within team, etc.
	// Implementation depends on your team model structure
	_ = userID
	_ = resourceID
	_ = action

	// Simplified implementation
	return true // Placeholder
}
