// Package middleware provides HTTP middleware functions
package middleware

import (
	"net/http"
	"strings"

	"github.com/HXLOS202653/ycg_cloud/cloud-backend/internal/service"
	"github.com/gin-gonic/gin"
)

// AuthMiddleware creates JWT authentication middleware
func AuthMiddleware(authService *service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"code":    "AUTH_TOKEN_MISSING",
				"message": "Authorization header is required",
				"error": gin.H{
					"type":    "AUTHENTICATION_ERROR",
					"details": "No authorization token provided",
				},
			})
			c.Abort()
			return
		}

		// Check Bearer token format
		tokenParts := strings.SplitN(authHeader, " ", 2)
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"code":    "AUTH_TOKEN_INVALID_FORMAT",
				"message": "Invalid authorization header format",
				"error": gin.H{
					"type":    "AUTHENTICATION_ERROR",
					"details": "Authorization header must be in format: Bearer <token>",
				},
			})
			c.Abort()
			return
		}

		token := tokenParts[1]

		// Validate token
		claims, err := authService.ValidateToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"code":    "AUTH_TOKEN_INVALID",
				"message": "Invalid or expired token",
				"error": gin.H{
					"type":    "AUTHENTICATION_ERROR",
					"details": err.Error(),
				},
			})
			c.Abort()
			return
		}

		// Set user information in context
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("email", claims.Email)
		c.Set("role", claims.Role)
		c.Set("claims", claims)

		c.Next()
	}
}

// OptionalAuthMiddleware creates optional JWT authentication middleware
// This middleware does not abort on missing/invalid tokens but sets user info if valid
func OptionalAuthMiddleware(authService *service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		// Check Bearer token format
		tokenParts := strings.SplitN(authHeader, " ", 2)
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.Next()
			return
		}

		token := tokenParts[1]

		// Validate token
		claims, err := authService.ValidateToken(token)
		if err != nil {
			c.Next()
			return
		}

		// Set user information in context
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("email", claims.Email)
		c.Set("role", claims.Role)
		c.Set("claims", claims)
		c.Set("authenticated", true)

		c.Next()
	}
}

// RequireRole creates middleware to check user roles
func RequireRole(_ *service.AuthService, allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// First ensure user is authenticated
		userRole, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"code":    "AUTH_NOT_AUTHENTICATED",
				"message": "Authentication required",
				"error": gin.H{
					"type":    "AUTHENTICATION_ERROR",
					"details": "User must be authenticated to access this resource",
				},
			})
			c.Abort()
			return
		}

		role, ok := userRole.(string)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"code":    "AUTH_ROLE_ERROR",
				"message": "Invalid role information",
				"error": gin.H{
					"type":    "SYSTEM_ERROR",
					"details": "Unable to determine user role",
				},
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
			"code":    "AUTH_PERMISSION_DENIED",
			"message": "Insufficient permissions",
			"error": gin.H{
				"type":           "PERMISSION_ERROR",
				"details":        "User does not have required permissions for this resource",
				"required_roles": allowedRoles,
				"user_role":      role,
			},
		})
		c.Abort()
	}
}

// GetUserID extracts user ID from context
func GetUserID(c *gin.Context) (int64, bool) {
	userID, exists := c.Get("user_id")
	if !exists {
		return 0, false
	}

	id, ok := userID.(int64)
	return id, ok
}

// GetUserInfo extracts user information from context
func GetUserInfo(c *gin.Context) map[string]interface{} {
	userInfo := make(map[string]interface{})

	if userID, exists := c.Get("user_id"); exists {
		userInfo["user_id"] = userID
	}

	if username, exists := c.Get("username"); exists {
		userInfo["username"] = username
	}

	if email, exists := c.Get("email"); exists {
		userInfo["email"] = email
	}

	if role, exists := c.Get("role"); exists {
		userInfo["role"] = role
	}

	return userInfo
}

// IsAuthenticated checks if user is authenticated
func IsAuthenticated(c *gin.Context) bool {
	_, exists := c.Get("user_id")
	return exists
}
