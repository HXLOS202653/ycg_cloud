// Package middleware provides permission checking middleware
package middleware

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/HXLOS202653/ycg_cloud/cloud-backend/internal/model"
	"github.com/HXLOS202653/ycg_cloud/cloud-backend/internal/service"
	"github.com/gin-gonic/gin"
)

// PermissionMiddleware creates a middleware that checks permissions
func PermissionMiddleware(permissionService *service.PermissionService, resource, action string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user ID from context (set by auth middleware)
		userID, exists := GetUserID(c)
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"code":    "AUTH_NOT_AUTHENTICATED",
				"message": "用户未认证",
				"error": gin.H{
					"type":    "AUTHENTICATION_ERROR",
					"details": "User must be authenticated to access this resource",
				},
			})
			c.Abort()
			return
		}

		// Extract resource ID from URL parameters if needed
		resourceID := extractResourceID(c, resource)

		// Check permission
		result := permissionService.CheckPermission(userID, resource, action, resourceID)
		if !result.Allowed {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"code":    "AUTH_PERMISSION_DENIED",
				"message": "权限不足",
				"error": gin.H{
					"type":     "PERMISSION_ERROR",
					"details":  "Insufficient permissions to perform this action",
					"resource": resource,
					"action":   action,
					"reason":   result.Reason,
					"user_id":  userID,
				},
			})
			c.Abort()
			return
		}

		// Set permission info in context for potential use by handlers
		c.Set("permission_check_result", result)
		c.Set("checked_resource", resource)
		c.Set("checked_action", action)

		c.Next()
	}
}

// RequirePermission creates a middleware that requires specific permission
func RequirePermission(permissionService *service.PermissionService, resource, action string) gin.HandlerFunc {
	return PermissionMiddleware(permissionService, resource, action)
}

// RequireAnyPermission creates a middleware that requires any of the specified permissions
func RequireAnyPermission(permissionService *service.PermissionService, permissions ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := GetUserID(c)
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"code":    "AUTH_NOT_AUTHENTICATED",
				"message": "用户未认证",
			})
			c.Abort()
			return
		}

		// Check if user has any of the required permissions
		for _, perm := range permissions {
			parts := strings.Split(perm, ":")
			if len(parts) != 2 {
				continue
			}

			resource, action := parts[0], parts[1]
			if permissionService.HasPermission(userID, resource, action) {
				c.Set("matched_permission", perm)
				c.Next()
				return
			}
		}

		// No permission found
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"code":    "AUTH_PERMISSION_DENIED",
			"message": "权限不足",
			"error": gin.H{
				"type":                 "PERMISSION_ERROR",
				"details":              "User does not have any of the required permissions",
				"required_permissions": permissions,
				"user_id":              userID,
			},
		})
		c.Abort()
	}
}

// RequireAllPermissions creates a middleware that requires all specified permissions
func RequireAllPermissions(permissionService *service.PermissionService, permissions ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := GetUserID(c)
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"code":    "AUTH_NOT_AUTHENTICATED",
				"message": "用户未认证",
			})
			c.Abort()
			return
		}

		// Check if user has all required permissions
		missingPerms := []string{}
		for _, perm := range permissions {
			parts := strings.Split(perm, ":")
			if len(parts) != 2 {
				missingPerms = append(missingPerms, perm)
				continue
			}

			resource, action := parts[0], parts[1]
			if !permissionService.HasPermission(userID, resource, action) {
				missingPerms = append(missingPerms, perm)
			}
		}

		if len(missingPerms) > 0 {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"code":    "AUTH_PERMISSION_DENIED",
				"message": "权限不足",
				"error": gin.H{
					"type":                "PERMISSION_ERROR",
					"details":             "User does not have all required permissions",
					"missing_permissions": missingPerms,
					"user_id":             userID,
				},
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireSpecificRole creates a middleware that requires specific role
func RequireSpecificRole(_ *service.PermissionService, allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"code":    "AUTH_NOT_AUTHENTICATED",
				"message": "用户未认证",
			})
			c.Abort()
			return
		}

		role, ok := userRole.(string)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"code":    "AUTH_ROLE_ERROR",
				"message": "角色信息错误",
			})
			c.Abort()
			return
		}

		// Check if user role is in allowed roles
		for _, allowedRole := range allowedRoles {
			if role == allowedRole {
				c.Next()
				return
			}
		}

		// Role not allowed
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"code":    "AUTH_ROLE_DENIED",
			"message": "角色权限不足",
			"error": gin.H{
				"type":           "ROLE_ERROR",
				"details":        "User role is not authorized for this resource",
				"required_roles": allowedRoles,
				"user_role":      role,
			},
		})
		c.Abort()
	}
}

// RequireOwnership creates a middleware that requires resource ownership
func RequireOwnership(_ *service.PermissionService, resourceType string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := GetUserID(c)
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"code":    "AUTH_NOT_AUTHENTICATED",
				"message": "用户未认证",
			})
			c.Abort()
			return
		}

		resourceID := extractResourceID(c, resourceType)
		if resourceID == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"code":    "RESOURCE_ID_MISSING",
				"message": "资源ID缺失",
			})
			c.Abort()
			return
		}

		// Check ownership (this would need to be implemented based on your data model)
		if !checkResourceOwnership(userID, resourceType, resourceID) {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"code":    "AUTH_OWNERSHIP_REQUIRED",
				"message": "需要资源所有权",
				"error": gin.H{
					"type":        "OWNERSHIP_ERROR",
					"details":     "User must own the resource to perform this action",
					"resource_id": resourceID,
					"user_id":     userID,
				},
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// AdminOnly creates a middleware that only allows admin users
func AdminOnly(permissionService *service.PermissionService) gin.HandlerFunc {
	return RequireSpecificRole(permissionService, model.RoleAdmin, model.RoleSuperAdmin)
}

// SuperAdminOnly creates a middleware that only allows super admin users
func SuperAdminOnly(permissionService *service.PermissionService) gin.HandlerFunc {
	return RequireSpecificRole(permissionService, model.RoleSuperAdmin)
}

// FilePermission creates a middleware for file operations
func FilePermission(permissionService *service.PermissionService, action string) gin.HandlerFunc {
	return PermissionMiddleware(permissionService, model.ResourceFile, action)
}

// FolderPermission creates a middleware for folder operations
func FolderPermission(permissionService *service.PermissionService, action string) gin.HandlerFunc {
	return PermissionMiddleware(permissionService, model.ResourceFolder, action)
}

// TeamPermission creates a middleware for team operations
func TeamPermission(permissionService *service.PermissionService, action string) gin.HandlerFunc {
	return PermissionMiddleware(permissionService, model.ResourceTeam, action)
}

// SystemPermission creates a middleware for system operations
func SystemPermission(permissionService *service.PermissionService, action string) gin.HandlerFunc {
	return PermissionMiddleware(permissionService, model.ResourceSystem, action)
}

// Helper functions

// extractResourceID extracts resource ID from URL parameters
func extractResourceID(c *gin.Context, resourceType string) string {
	// Try common parameter names
	paramNames := []string{"id", resourceType + "_id", resourceType + "Id"}

	for _, paramName := range paramNames {
		if id := c.Param(paramName); id != "" {
			return id
		}
		if id := c.Query(paramName); id != "" {
			return id
		}
	}

	return ""
}

// checkResourceOwnership checks if user owns the resource
func checkResourceOwnership(userID int64, resourceType, resourceID string) bool {
	// This is a placeholder implementation
	// You would implement this based on your actual data models

	// Convert resourceID to int64 if needed
	id, err := strconv.ParseInt(resourceID, 10, 64)
	if err != nil {
		return false
	}

	// Check ownership based on resource type
	switch resourceType {
	case model.ResourceFile:
		// Check if user owns the file
		return checkFileOwnership(userID, id)
	case model.ResourceFolder:
		// Check if user owns the folder
		return checkFolderOwnership(userID, id)
	case model.ResourceTeam:
		// Check if user is team owner/admin
		return checkTeamOwnership(userID, id)
	default:
		// For other resources, assume ownership
		return true
	}
}

// checkFileOwnership checks if user owns a specific file
func checkFileOwnership(userID, fileID int64) bool {
	// Placeholder implementation
	// You would query your file table here
	_ = userID
	_ = fileID
	return true
}

// checkFolderOwnership checks if user owns a specific folder
func checkFolderOwnership(userID, folderID int64) bool {
	// Placeholder implementation
	// You would query your folder table here
	_ = userID
	_ = folderID
	return true
}

// checkTeamOwnership checks if user owns/administers a specific team
func checkTeamOwnership(userID, teamID int64) bool {
	// Placeholder implementation
	// You would query your team membership table here
	_ = userID
	_ = teamID
	return true
}

// PermissionInfo returns information about permission check
type PermissionInfo struct {
	Resource   string `json:"resource"`
	Action     string `json:"action"`
	ResourceID string `json:"resourceId,omitempty"`
	Allowed    bool   `json:"allowed"`
	Reason     string `json:"reason"`
}

// GetPermissionInfo extracts permission information from context
func GetPermissionInfo(c *gin.Context) *PermissionInfo {
	result, exists := c.Get("permission_check_result")
	if !exists {
		return nil
	}

	checkResult, ok := result.(*service.PermissionCheckResult)
	if !ok {
		return nil
	}

	resource, _ := c.Get("checked_resource")
	action, _ := c.Get("checked_action")

	return &PermissionInfo{
		Resource: resource.(string),
		Action:   action.(string),
		Allowed:  checkResult.Allowed,
		Reason:   checkResult.Reason,
	}
}
