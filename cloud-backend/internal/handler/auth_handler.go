// Package handler provides HTTP request handlers
package handler

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/HXLOS202653/ycg_cloud/cloud-backend/internal/dto"
	"github.com/HXLOS202653/ycg_cloud/cloud-backend/internal/model"
	"github.com/HXLOS202653/ycg_cloud/cloud-backend/internal/service"
)

// AuthHandler handles authentication-related HTTP requests
type AuthHandler struct {
	authService         *service.AuthService
	emailService        *service.EmailService
	verificationService *service.VerificationService
	sessionService      *service.SessionService
	twoFactorService    *service.TwoFactorService
	db                  *gorm.DB
}

// NewAuthHandler creates a new authentication handler
func NewAuthHandler(authService *service.AuthService, emailService *service.EmailService, verificationService *service.VerificationService, sessionService *service.SessionService, twoFactorService *service.TwoFactorService, db *gorm.DB) *AuthHandler {
	return &AuthHandler{
		authService:         authService,
		emailService:        emailService,
		verificationService: verificationService,
		sessionService:      sessionService,
		twoFactorService:    twoFactorService,
		db:                  db,
	}
}

// Register handles user registration
// @Summary Register a new user
// @Description Register a new user account with email verification
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body dto.RegisterRequest true "Registration data"
// @Success 201 {object} dto.UserResponse
// @Failure 400 {object} map[string]interface{} "Validation error"
// @Failure 409 {object} map[string]interface{} "User already exists"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/v1/auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.sendValidationError(c, "USER_VALIDATION_REQUEST_INVALID", "请求参数验证失败", err.Error())
		return
	}

	// Validate the registration request
	if err := h.validateRegistrationRequest(&req); err != nil {
		h.sendValidationErrorFromValidation(c, err)
		return
	}

	// Check for existing users
	if err := h.checkExistingUser(&req, c); err != nil {
		return // Error already sent in function
	}

	// Create and save the user
	user, err := h.createUserFromRequest(&req)
	if err != nil {
		h.sendSystemError(c, "USER_SYSTEM_CREATE_FAILED", "用户创建失败", err.Error())
		return
	}

	// Send verification email and respond
	h.handleRegistrationSuccess(c, user)
}

// Login handles user login
// @Summary User login
// @Description Authenticate user and return JWT tokens
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "Login credentials"
// @Success 200 {object} dto.LoginResponse
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 401 {object} map[string]interface{} "Authentication failed"
// @Failure 423 {object} map[string]interface{} "Account locked"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/v1/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    "AUTH_VALIDATION_REQUEST_INVALID",
			"message": "请求参数验证失败",
			"error": gin.H{
				"type":    "VALIDATION_ERROR",
				"details": err.Error(),
			},
		})
		return
	}

	// Find user by username or email
	var user model.User
	result := h.db.Where("username = ? OR email = ?", req.Username, req.Username).First(&user)
	if result.Error != nil {
		// User not found - return generic error for security
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"code":    "AUTH_BUSINESS_LOGIN_FAILED",
			"message": "用户名或密码错误",
			"error": gin.H{
				"type":    "AUTHENTICATION_ERROR",
				"details": "Invalid username or password",
			},
		})
		return
	}

	// Check if account is locked
	if user.IsLocked() {
		c.JSON(http.StatusLocked, gin.H{
			"success": false,
			"code":    "AUTH_BUSINESS_ACCOUNT_LOCKED",
			"message": "账户已被临时锁定",
			"error": gin.H{
				"type":         "AUTHENTICATION_ERROR",
				"details":      fmt.Sprintf("Account is locked until %s", user.LockedUntil.Format("2006-01-02 15:04:05")),
				"locked_until": user.LockedUntil,
			},
		})
		return
	}

	// Check if account is active
	if !user.IsActive() {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"code":    "AUTH_BUSINESS_ACCOUNT_INACTIVE",
			"message": "账户未激活",
			"error": gin.H{
				"type":    "AUTHENTICATION_ERROR",
				"details": "Account is not active",
				"status":  user.Status,
			},
		})
		return
	}

	// Verify password
	if err := h.authService.VerifyPassword(user.PasswordHash, req.Password); err != nil {
		// Increment login attempts
		user.IncrementLoginAttempts()
		h.db.Save(&user)

		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"code":    "AUTH_BUSINESS_LOGIN_FAILED",
			"message": "用户名或密码错误",
			"error": gin.H{
				"type":               "AUTHENTICATION_ERROR",
				"details":            "Invalid username or password",
				"remaining_attempts": 5 - user.LoginAttempts,
			},
		})
		return
	}

	// Generate tokens
	tokenPair, err := h.authService.GenerateTokens(user.ID, user.Username, user.Email, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code":    "AUTH_SYSTEM_TOKEN_GENERATION_FAILED",
			"message": "令牌生成失败",
			"error": gin.H{
				"type":    "SYSTEM_ERROR",
				"details": "Failed to generate authentication tokens",
			},
		})
		return
	}

	// Update user login information
	clientIP := c.ClientIP()
	user.UpdateLastLogin(clientIP)
	h.db.Save(&user)

	// Extract device information
	userAgent := c.GetHeader("User-Agent")
	deviceInfo := h.sessionService.ExtractDeviceInfo(userAgent)

	// Create user session record
	session, err := h.sessionService.CreateSession(user.ID, tokenPair.RefreshToken, deviceInfo, clientIP, userAgent)
	if err != nil {
		// Log error but don't fail login
		fmt.Printf("Warning: Failed to create session: %v\n", err)
	}

	// Prepare response
	userResponse := h.userToResponse(&user)
	loginResponse := dto.LoginResponse{
		User:         *userResponse,
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresAt:    tokenPair.ExpiresAt,
		TokenType:    tokenPair.TokenType,
	}

	// Add session information if available
	if session != nil {
		loginResponse.SessionInfo = &dto.SessionResponse{
			ID:         session.ID,
			DeviceInfo: session.DeviceInfo,
			IPAddress:  session.IPAddress,
			CreatedAt:  session.CreatedAt,
			ExpiresAt:  session.ExpiresAt,
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"code":    200,
		"message": "登录成功",
		"data":    loginResponse,
	})
}

// isPasswordStrong validates password strength
func (h *AuthHandler) isPasswordStrong(password string) bool {
	// At least 8 characters
	if len(password) < 8 {
		return false
	}

	// Must contain uppercase, lowercase, and number
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasNumber := regexp.MustCompile(`\d`).MatchString(password)

	return hasUpper && hasLower && hasNumber
}

// isValidEmail validates email format
func (h *AuthHandler) isValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$`)
	return emailRegex.MatchString(email)
}

// isValidUsername validates username format
func (h *AuthHandler) isValidUsername(username string) bool {
	// 3-50 characters, letters, numbers, underscores only
	if len(username) < 3 || len(username) > 50 {
		return false
	}
	usernameRegex := regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
	return usernameRegex.MatchString(username)
}

// userToResponse converts User model to UserResponse DTO
func (h *AuthHandler) userToResponse(user *model.User) *dto.UserResponse {
	var realName, phone, avatar, lastLoginIP string
	if user.RealName != nil {
		realName = *user.RealName
	}
	if user.Phone != nil {
		phone = *user.Phone
	}
	if user.AvatarURL != nil {
		avatar = *user.AvatarURL
	}
	if user.LastLoginIP != nil {
		lastLoginIP = *user.LastLoginIP
	}

	return &dto.UserResponse{
		ID:                     user.ID,
		Username:               user.Username,
		Email:                  user.Email,
		RealName:               realName,
		Phone:                  phone,
		Avatar:                 avatar,
		Status:                 user.Status,
		Role:                   user.Role,
		EmailVerified:          user.EmailVerified,
		TwoFactorEnabled:       user.TwoFactorEnabled,
		StorageQuota:           user.StorageQuota,
		StorageUsed:            user.StorageUsed,
		StorageUsagePercentage: user.GetStorageUsagePercentage(),
		Language:               user.Language,
		Timezone:               user.Timezone,
		Theme:                  user.Theme,
		LastLoginAt:            user.LastLoginAt,
		LastLoginIP:            lastLoginIP,
		CreatedAt:              user.CreatedAt,
		UpdatedAt:              user.UpdatedAt,
	}
}

// GetCurrentUser returns current authenticated user info
func (h *AuthHandler) GetCurrentUser(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"code":    "AUTH_NOT_AUTHENTICATED",
			"message": "用户未认证",
		})
		return
	}

	var user model.User
	if err := h.db.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"code":    "USER_NOT_FOUND",
			"message": "用户不存在",
		})
		return
	}

	userResponse := h.userToResponse(&user)
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"code":    200,
		"message": "获取用户信息成功",
		"data":    userResponse,
	})
}

// SendVerificationCode sends a verification code to email
// @Summary Send verification code
// @Description Send a verification code to user's email
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body dto.SendVerificationCodeRequest true "Verification code request"
// @Success 200 {object} dto.VerificationCodeResponse
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 429 {object} map[string]interface{} "Rate limit exceeded"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/v1/auth/send-verification [post]
func (h *AuthHandler) SendVerificationCode(c *gin.Context) {
	var req dto.SendVerificationCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    "USER_VALIDATION_REQUEST_INVALID",
			"message": "请求参数验证失败",
			"error": gin.H{
				"type":    "VALIDATION_ERROR",
				"details": err.Error(),
			},
		})
		return
	}

	// Send verification code
	if err := h.verificationService.SendVerificationCode(req.Email, req.Purpose); err != nil {
		if err.Error() == "too many verification code requests, please try again later" {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"success": false,
				"code":    "USER_BUSINESS_RATE_LIMIT_EXCEEDED",
				"message": "验证码请求过于频繁",
				"error": gin.H{
					"type":    "BUSINESS_ERROR",
					"details": err.Error(),
				},
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code":    "USER_SYSTEM_EMAIL_SEND_FAILED",
			"message": "验证码发送失败",
			"error": gin.H{
				"type":    "SYSTEM_ERROR",
				"details": err.Error(),
			},
		})
		return
	}

	// Get verification code info for response
	info, _ := h.verificationService.GetVerificationCodeInfo(req.Email, req.Purpose)

	response := dto.VerificationCodeResponse{
		Email:   req.Email,
		Purpose: req.Purpose,
		Message: "验证码已发送到您的邮箱",
	}

	if info != nil {
		if expiresAt, ok := info["expires_at"].(time.Time); ok {
			response.ExpiresAt = expiresAt
		}
		if maxAttempts, ok := info["max_attempts"].(int); ok {
			response.MaxAttempts = maxAttempts
		}
		if timeRemaining, ok := info["time_remaining"].(int); ok {
			response.TimeRemaining = timeRemaining
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"code":    200,
		"message": "验证码发送成功",
		"data":    response,
	})
}

// VerifyEmail verifies user's email with verification code
// @Summary Verify email
// @Description Verify user's email address with verification code
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body dto.VerifyEmailRequest true "Email verification data"
// @Success 200 {object} dto.UserResponse
// @Failure 400 {object} map[string]interface{} "Invalid verification code"
// @Failure 404 {object} map[string]interface{} "User not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/v1/auth/verify-email [post]
func (h *AuthHandler) VerifyEmail(c *gin.Context) {
	var req dto.VerifyEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    "USER_VALIDATION_REQUEST_INVALID",
			"message": "请求参数验证失败",
			"error": gin.H{
				"type":    "VALIDATION_ERROR",
				"details": err.Error(),
			},
		})
		return
	}

	// Verify the code
	if err := h.verificationService.VerifyCode(req.Email, req.VerificationCode, "registration"); err != nil {
		statusCode := http.StatusBadRequest
		code := "USER_BUSINESS_VERIFICATION_FAILED"

		if err.Error() == "verification code not found or expired" {
			code = "USER_BUSINESS_VERIFICATION_EXPIRED"
		} else if err.Error() == "maximum verification attempts exceeded" {
			code = "USER_BUSINESS_VERIFICATION_ATTEMPTS_EXCEEDED"
		}

		c.JSON(statusCode, gin.H{
			"success": false,
			"code":    code,
			"message": "邮箱验证失败",
			"error": gin.H{
				"type":    "BUSINESS_ERROR",
				"details": err.Error(),
			},
		})
		return
	}

	// Find and update user
	var user model.User
	result := h.db.Where("email = ?", req.Email).First(&user)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"code":    "USER_RESOURCE_NOT_FOUND",
			"message": "用户不存在",
			"error": gin.H{
				"type":    "RESOURCE_ERROR",
				"details": "User with this email not found",
			},
		})
		return
	}

	// Update user verification status
	user.EmailVerified = true
	user.Status = "active" // Activate user after email verification

	if err := h.db.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code":    "USER_SYSTEM_UPDATE_FAILED",
			"message": "用户状态更新失败",
			"error": gin.H{
				"type":    "SYSTEM_ERROR",
				"details": "Failed to update user verification status",
			},
		})
		return
	}

	// Return updated user info
	userResponse := h.userToResponse(&user)
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"code":    200,
		"message": "邮箱验证成功",
		"data":    userResponse,
	})
}

// ResendVerificationCode resends verification code
// @Summary Resend verification code
// @Description Resend verification code to user's email
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body dto.ResendVerificationRequest true "Resend request"
// @Success 200 {object} dto.VerificationCodeResponse
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 429 {object} map[string]interface{} "Rate limit exceeded"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/v1/auth/resend-verification [post]
func (h *AuthHandler) ResendVerificationCode(c *gin.Context) {
	var req dto.ResendVerificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    "USER_VALIDATION_REQUEST_INVALID",
			"message": "请求参数验证失败",
			"error": gin.H{
				"type":    "VALIDATION_ERROR",
				"details": err.Error(),
			},
		})
		return
	}

	// Resend verification code
	if err := h.verificationService.ResendVerificationCode(req.Email, "registration"); err != nil {
		statusCode := http.StatusBadRequest
		code := "USER_BUSINESS_RESEND_FAILED"

		if err.Error() == "too many verification code requests, please try again later" {
			statusCode = http.StatusTooManyRequests
			code = "USER_BUSINESS_RATE_LIMIT_EXCEEDED"
		}

		c.JSON(statusCode, gin.H{
			"success": false,
			"code":    code,
			"message": "验证码重发失败",
			"error": gin.H{
				"type":    "BUSINESS_ERROR",
				"details": err.Error(),
			},
		})
		return
	}

	// Get verification code info for response
	info, _ := h.verificationService.GetVerificationCodeInfo(req.Email, "registration")

	response := dto.VerificationCodeResponse{
		Email:   req.Email,
		Purpose: "registration",
		Message: "验证码已重新发送到您的邮箱",
	}

	if info != nil {
		if expiresAt, ok := info["expires_at"].(time.Time); ok {
			response.ExpiresAt = expiresAt
		}
		if maxAttempts, ok := info["max_attempts"].(int); ok {
			response.MaxAttempts = maxAttempts
		}
		if timeRemaining, ok := info["time_remaining"].(int); ok {
			response.TimeRemaining = timeRemaining
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"code":    200,
		"message": "验证码重发成功",
		"data":    response,
	})
}

// RefreshToken handles token refresh requests
// @Summary Refresh access token
// @Description Refresh expired access token using refresh token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body dto.RefreshTokenRequest true "Refresh token request"
// @Success 200 {object} dto.TokenResponse
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 401 {object} map[string]interface{} "Invalid refresh token"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/v1/auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req dto.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    "AUTH_VALIDATION_REQUEST_INVALID",
			"message": "请求参数验证失败",
			"error": gin.H{
				"type":    "VALIDATION_ERROR",
				"details": err.Error(),
			},
		})
		return
	}

	// Validate refresh token
	userID, err := h.authService.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"code":    "AUTH_REFRESH_TOKEN_INVALID",
			"message": "刷新令牌无效或已过期",
			"error": gin.H{
				"type":    "AUTHENTICATION_ERROR",
				"details": err.Error(),
			},
		})
		return
	}

	// Get user information
	var user model.User
	if err := h.db.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"code":    "AUTH_USER_NOT_FOUND",
			"message": "用户不存在",
			"error": gin.H{
				"type":    "AUTHENTICATION_ERROR",
				"details": "User associated with refresh token not found",
			},
		})
		return
	}

	// Check if user is still active
	if !user.IsActive() {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"code":    "AUTH_USER_INACTIVE",
			"message": "用户账户已被停用",
			"error": gin.H{
				"type":    "AUTHENTICATION_ERROR",
				"details": "User account is not active",
			},
		})
		return
	}

	// Generate new token pair
	tokenPair, err := h.authService.GenerateTokens(user.ID, user.Username, user.Email, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code":    "AUTH_TOKEN_GENERATION_FAILED",
			"message": "令牌生成失败",
			"error": gin.H{
				"type":    "SYSTEM_ERROR",
				"details": "Failed to generate new tokens",
			},
		})
		return
	}

	// Update session with new refresh token
	clientIP := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")
	deviceInfo := req.DeviceInfo
	if deviceInfo == "" {
		deviceInfo = h.sessionService.ExtractDeviceInfo(userAgent)
	}

	// Find and update existing session
	var session model.UserSession
	if err := h.db.Where("user_id = ? AND refresh_token = ?", userID, req.RefreshToken).First(&session).Error; err == nil {
		// Update existing session with new refresh token
		session.RefreshToken = tokenPair.RefreshToken
		session.DeviceInfo = deviceInfo
		session.IPAddress = clientIP
		session.ExpiresAt = tokenPair.ExpiresAt
		h.db.Save(&session)
	} else {
		// Create new session if old one not found
		if _, err := h.sessionService.CreateSession(user.ID, tokenPair.RefreshToken, deviceInfo, clientIP, userAgent); err != nil {
			// Log error but don't fail the refresh process
			// Session creation failure shouldn't block token refresh
			fmt.Printf("Failed to create session during token refresh: %v\n", err)
		}
	}

	// Prepare response
	tokenResponse := dto.TokenResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresAt:    tokenPair.ExpiresAt,
		TokenType:    tokenPair.TokenType,
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"code":    200,
		"message": "令牌刷新成功",
		"data":    tokenResponse,
	})
}

// Logout handles user logout
// @Summary User logout
// @Description Logout user and invalidate session
// @Tags Authentication
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{} "Not authenticated"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/v1/auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"code":    "AUTH_NOT_AUTHENTICATED",
			"message": "用户未认证",
		})
		return
	}

	// Get refresh token from request body or header
	var req struct {
		RefreshToken string `json:"refreshToken,omitempty"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		// If JSON binding fails, continue to try header-based refresh token
		// This allows flexible token refresh from both body and header
		fmt.Printf("JSON binding failed for refresh token request, trying header: %v\n", err)
	}

	if req.RefreshToken == "" {
		// Try to get from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
			req.RefreshToken = strings.TrimPrefix(authHeader, "Bearer ")
		}
	}

	// Invalidate specific session if refresh token provided
	if req.RefreshToken != "" {
		var session model.UserSession
		if err := h.db.Where("user_id = ? AND refresh_token = ?", userID, req.RefreshToken).First(&session).Error; err == nil {
			if err := h.sessionService.InvalidateSession(session.SessionToken); err != nil {
				// Log error but continue logout process
				fmt.Printf("Failed to invalidate specific session during logout: %v\n", err)
			}
		}
	} else {
		// Invalidate all sessions for user if no specific token provided
		if userIDInt, ok := userID.(int64); ok {
			if err := h.sessionService.InvalidateAllUserSessions(userIDInt); err != nil {
				// Log error but continue logout process
				fmt.Printf("Failed to invalidate all user sessions during logout: %v\n", err)
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"code":    200,
		"message": "退出登录成功",
		"data": gin.H{
			"timestamp": time.Now(),
		},
	})
}

// LogoutAll handles logout from all devices
// @Summary Logout from all devices
// @Description Logout user from all devices and invalidate all sessions
// @Tags Authentication
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{} "Not authenticated"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/v1/auth/logout-all [post]
func (h *AuthHandler) LogoutAll(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"code":    "AUTH_NOT_AUTHENTICATED",
			"message": "用户未认证",
		})
		return
	}

	// Invalidate all sessions for user
	userIDInt, ok := userID.(int64)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code":    "AUTH_USER_ID_INVALID",
			"message": "用户ID格式错误",
		})
		return
	}

	if err := h.sessionService.InvalidateAllUserSessions(userIDInt); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code":    "AUTH_LOGOUT_FAILED",
			"message": "退出登录失败",
			"error": gin.H{
				"type":    "SYSTEM_ERROR",
				"details": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"code":    200,
		"message": "已从所有设备退出登录",
		"data": gin.H{
			"timestamp": time.Now(),
		},
	})
}

// GetActiveSessions returns user's active sessions
// @Summary Get active sessions
// @Description Get list of user's active sessions
// @Tags Authentication
// @Accept json
// @Produce json
// @Success 200 {object} []dto.SessionResponse
// @Failure 401 {object} map[string]interface{} "Not authenticated"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/v1/auth/sessions [get]
func (h *AuthHandler) GetActiveSessions(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"code":    "AUTH_NOT_AUTHENTICATED",
			"message": "用户未认证",
		})
		return
	}

	// Get active sessions
	userIDInt, ok := userID.(int64)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code":    "AUTH_USER_ID_INVALID",
			"message": "用户ID格式错误",
		})
		return
	}

	sessions, err := h.sessionService.GetUserActiveSessions(userIDInt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code":    "AUTH_SESSIONS_FETCH_FAILED",
			"message": "获取会话列表失败",
			"error": gin.H{
				"type":    "SYSTEM_ERROR",
				"details": err.Error(),
			},
		})
		return
	}

	// Convert to response DTOs
	sessionResponses := make([]dto.SessionResponse, len(sessions))
	for i := range sessions {
		sessionResponses[i] = dto.SessionResponse{
			ID:         sessions[i].ID,
			DeviceInfo: sessions[i].DeviceInfo,
			IPAddress:  sessions[i].IPAddress,
			CreatedAt:  sessions[i].CreatedAt,
			ExpiresAt:  sessions[i].ExpiresAt,
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"code":    200,
		"message": "获取会话列表成功",
		"data":    sessionResponses,
	})
}

// Setup2FA handles 2FA setup request
// @Summary Setup two-factor authentication
// @Description Setup 2FA for a user account
// @Tags Authentication
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} map[string]interface{} "2FA setup data"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/v1/auth/2fa/setup [post]
func (h *AuthHandler) Setup2FA(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"code":    "AUTH_USER_NOT_FOUND",
			"message": "用户信息不存在",
		})
		return
	}

	userIDInt, ok := userID.(int64)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code":    "AUTH_USER_ID_INVALID",
			"message": "用户ID格式错误",
		})
		return
	}

	// Get user email
	var user model.User
	if err := h.db.Select("email").First(&user, userIDInt).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code":    "AUTH_USER_FETCH_FAILED",
			"message": "获取用户信息失败",
			"error": gin.H{
				"type":    "DATABASE_ERROR",
				"details": err.Error(),
			},
		})
		return
	}

	// Setup 2FA
	setupData, err := h.twoFactorService.Setup2FA(userIDInt, user.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    "AUTH_2FA_SETUP_FAILED",
			"message": "2FA设置失败",
			"error": gin.H{
				"type":    "SETUP_ERROR",
				"details": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"code":    "AUTH_2FA_SETUP_SUCCESS",
		"message": "2FA设置成功",
		"data":    setupData,
	})
}

// Enable2FA handles 2FA enable request
// @Summary Enable two-factor authentication
// @Description Enable 2FA after verifying TOTP code
// @Tags Authentication
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body dto.Enable2FARequest true "Enable 2FA request"
// @Success 200 {object} map[string]interface{} "Success message"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/v1/auth/2fa/enable [post]
func (h *AuthHandler) Enable2FA(c *gin.Context) {
	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"code":    "AUTH_USER_NOT_FOUND",
			"message": "用户信息不存在",
		})
		return
	}

	var req dto.Enable2FARequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    "AUTH_VALIDATION_REQUEST_INVALID",
			"message": "请求参数验证失败",
			"error": gin.H{
				"type":    "VALIDATION_ERROR",
				"details": err.Error(),
			},
		})
		return
	}

	userIDInt, ok := userID.(int64)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code":    "AUTH_USER_ID_INVALID",
			"message": "用户ID格式错误",
		})
		return
	}

	// Enable 2FA
	if err := h.twoFactorService.Enable2FA(userIDInt, req.TOTPCode); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    "AUTH_2FA_ENABLE_FAILED",
			"message": "启用2FA失败",
			"error": gin.H{
				"type":    "ENABLE_ERROR",
				"details": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"code":    "AUTH_2FA_ENABLE_SUCCESS",
		"message": "2FA启用成功",
	})
}

// Disable2FA handles 2FA disable request
// @Summary Disable two-factor authentication
// @Description Disable 2FA after verifying password and TOTP code
// @Tags Authentication
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body dto.Disable2FARequest true "Disable 2FA request"
// @Success 200 {object} map[string]interface{} "Success message"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/v1/auth/2fa/disable [post]
func (h *AuthHandler) Disable2FA(c *gin.Context) {
	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"code":    "AUTH_USER_NOT_FOUND",
			"message": "用户信息不存在",
		})
		return
	}

	var req dto.Disable2FARequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    "AUTH_VALIDATION_REQUEST_INVALID",
			"message": "请求参数验证失败",
			"error": gin.H{
				"type":    "VALIDATION_ERROR",
				"details": err.Error(),
			},
		})
		return
	}

	userIDInt, ok := userID.(int64)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code":    "AUTH_USER_ID_INVALID",
			"message": "用户ID格式错误",
		})
		return
	}

	// Disable 2FA
	if err := h.twoFactorService.Disable2FA(userIDInt, req.Password, req.TOTPCode); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    "AUTH_2FA_DISABLE_FAILED",
			"message": "禁用2FA失败",
			"error": gin.H{
				"type":    "DISABLE_ERROR",
				"details": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"code":    "AUTH_2FA_DISABLE_SUCCESS",
		"message": "2FA禁用成功",
	})
}

// Verify2FA handles 2FA verification request
// @Summary Verify two-factor authentication code
// @Description Verify TOTP or backup code for 2FA (temporary verification during login flow)
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body dto.Verify2FARequest true "Verify 2FA request"
// @Success 200 {object} map[string]interface{} "Success message"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 401 {object} map[string]interface{} "Invalid code"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/v1/auth/2fa/verify [post]
func (h *AuthHandler) Verify2FA(c *gin.Context) {
	var req dto.Verify2FARequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    "AUTH_VALIDATION_REQUEST_INVALID",
			"message": "请求参数验证失败",
			"error": gin.H{
				"type":    "VALIDATION_ERROR",
				"details": err.Error(),
			},
		})
		return
	}

	// Find user by username
	var user model.User
	if err := h.db.Where("username = ? OR email = ?", req.Username, req.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"code":    "AUTH_USER_NOT_FOUND",
			"message": "用户不存在",
		})
		return
	}

	// Verify 2FA code
	if err := h.twoFactorService.Verify2FA(user.ID, req.TOTPCode); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"code":    "AUTH_2FA_VERIFY_FAILED",
			"message": "2FA验证失败",
			"error": gin.H{
				"type":    "VERIFICATION_ERROR",
				"details": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"code":    "AUTH_2FA_VERIFY_SUCCESS",
		"message": "2FA验证成功",
	})
}

// RegenerateBackupCodes handles backup codes regeneration request
// @Summary Regenerate backup codes
// @Description Generate new backup codes for 2FA recovery
// @Tags Authentication
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body map[string]string true "Password and TOTP code"
// @Success 200 {object} map[string]interface{} "New backup codes"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/v1/auth/2fa/backup-codes [post]
func (h *AuthHandler) RegenerateBackupCodes(c *gin.Context) {
	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"code":    "AUTH_USER_NOT_FOUND",
			"message": "用户信息不存在",
		})
		return
	}

	var req struct {
		Password string `json:"password" binding:"required"`
		TOTPCode string `json:"totpCode" binding:"required,len=6"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    "AUTH_VALIDATION_REQUEST_INVALID",
			"message": "请求参数验证失败",
			"error": gin.H{
				"type":    "VALIDATION_ERROR",
				"details": err.Error(),
			},
		})
		return
	}

	userIDInt, ok := userID.(int64)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code":    "AUTH_USER_ID_INVALID",
			"message": "用户ID格式错误",
		})
		return
	}

	// Generate new backup codes
	codes, err := h.twoFactorService.GenerateBackupCodes(userIDInt, req.Password, req.TOTPCode)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    "AUTH_BACKUP_CODES_FAILED",
			"message": "生成备份码失败",
			"error": gin.H{
				"type":    "GENERATION_ERROR",
				"details": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"code":    "AUTH_BACKUP_CODES_SUCCESS",
		"message": "备份码生成成功",
		"data": gin.H{
			"backup_codes": codes,
		},
	})
}

// Helper functions for Register method to reduce complexity

// ValidationError represents a validation error with structured information
type ValidationError struct {
	Code    string
	Message string
	Details string
}

func (e *ValidationError) Error() string {
	return e.Message
}

func (h *AuthHandler) validateRegistrationRequest(req *dto.RegisterRequest) error {
	// Validate passwords match
	if req.Password != req.ConfirmPassword {
		return &ValidationError{
			Code:    "USER_VALIDATION_PASSWORD_MISMATCH",
			Message: "密码确认不匹配",
			Details: "Password and confirm password do not match",
		}
	}

	// Validate password strength
	if !h.isPasswordStrong(req.Password) {
		return &ValidationError{
			Code:    "USER_VALIDATION_PASSWORD_TOO_WEAK",
			Message: "密码强度不足",
			Details: "Password must be at least 8 characters long and contain uppercase, lowercase, and numbers",
		}
	}

	// Validate email format
	if !h.isValidEmail(req.Email) {
		return &ValidationError{
			Code:    "USER_VALIDATION_EMAIL_INVALID",
			Message: "邮箱格式不正确",
			Details: "Invalid email format",
		}
	}

	// Validate username format
	if !h.isValidUsername(req.Username) {
		return &ValidationError{
			Code:    "USER_VALIDATION_USERNAME_INVALID",
			Message: "用户名格式不正确",
			Details: "Username must be 3-50 characters long and contain only letters, numbers, and underscores",
		}
	}

	// Check if terms are accepted
	if !req.TermsAccepted {
		return &ValidationError{
			Code:    "USER_VALIDATION_TERMS_NOT_ACCEPTED",
			Message: "必须同意服务条款",
			Details: "Terms and conditions must be accepted",
		}
	}

	return nil
}

func (h *AuthHandler) checkExistingUser(req *dto.RegisterRequest, c *gin.Context) error {
	var existingUser model.User
	result := h.db.Where("email = ? OR username = ?", req.Email, req.Username).First(&existingUser)
	if result.Error == nil {
		// User exists, determine which field conflicts
		if existingUser.Email == req.Email {
			c.JSON(http.StatusConflict, gin.H{
				"success": false,
				"code":    "USER_BUSINESS_EMAIL_EXISTS",
				"message": "邮箱已被注册",
				"error": gin.H{
					"type":    "BUSINESS_ERROR",
					"details": "A user with this email address already exists",
					"field":   "email",
				},
			})
		} else {
			c.JSON(http.StatusConflict, gin.H{
				"success": false,
				"code":    "USER_BUSINESS_USERNAME_EXISTS",
				"message": "用户名已被占用",
				"error": gin.H{
					"type":    "BUSINESS_ERROR",
					"details": "A user with this username already exists",
					"field":   "username",
				},
			})
		}
		return fmt.Errorf("user already exists")
	}
	return nil
}

func (h *AuthHandler) createUserFromRequest(req *dto.RegisterRequest) (*model.User, error) {
	// Hash password
	hashedPassword, err := h.authService.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user
	var realName, phone *string
	if req.RealName != "" {
		realName = &req.RealName
	}
	if req.Phone != "" {
		phone = &req.Phone
	}

	user := model.User{
		Username:      req.Username,
		Email:         req.Email,
		PasswordHash:  hashedPassword,
		RealName:      realName,
		Phone:         phone,
		Language:      req.Language,
		Timezone:      req.Timezone,
		Status:        "pending", // 待邮箱验证
		Role:          "user",
		EmailVerified: false,
	}

	// Set default values if not provided
	if user.Language == "" {
		user.Language = "zh-CN"
	}
	if user.Timezone == "" {
		user.Timezone = "Asia/Shanghai"
	}

	// Save user to database
	if err := h.db.Create(&user).Error; err != nil {
		return nil, fmt.Errorf("failed to create user in database: %w", err)
	}

	return &user, nil
}

func (h *AuthHandler) handleRegistrationSuccess(c *gin.Context, user *model.User) {
	// Send verification email
	if err := h.verificationService.SendVerificationCode(user.Email, "registration"); err != nil {
		// Log error but don't fail registration
		fmt.Printf("Failed to send verification email: %v\n", err)

		// Still return success but with a warning
		c.JSON(http.StatusCreated, gin.H{
			"success": true,
			"code":    201,
			"message": "用户注册成功，但验证邮件发送失败",
			"data": gin.H{
				"user":    h.userToResponse(user),
				"warning": "邮件发送失败，请稍后手动请求验证码",
			},
		})
		return
	}

	// Convert to response DTO
	userResponse := h.userToResponse(user)

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"code":    201,
		"message": "用户注册成功，请查收验证邮件",
		"data": gin.H{
			"user": userResponse,
			"verification": gin.H{
				"email_sent": true,
				"message":    "验证码已发送到您的邮箱，有效期5分钟",
			},
		},
	})
}

func (h *AuthHandler) sendValidationError(c *gin.Context, code, message, details string) {
	c.JSON(http.StatusBadRequest, gin.H{
		"success": false,
		"code":    code,
		"message": message,
		"error": gin.H{
			"type":    "VALIDATION_ERROR",
			"details": details,
		},
	})
}

func (h *AuthHandler) sendValidationErrorFromValidation(c *gin.Context, err error) {
	var valErr *ValidationError
	if errors.As(err, &valErr) {
		h.sendValidationError(c, valErr.Code, valErr.Message, valErr.Details)
	} else {
		h.sendValidationError(c, "USER_VALIDATION_UNKNOWN", "验证失败", err.Error())
	}
}

func (h *AuthHandler) sendSystemError(c *gin.Context, code, message, details string) {
	c.JSON(http.StatusInternalServerError, gin.H{
		"success": false,
		"code":    code,
		"message": message,
		"error": gin.H{
			"type":    "SYSTEM_ERROR",
			"details": details,
		},
	})
}
