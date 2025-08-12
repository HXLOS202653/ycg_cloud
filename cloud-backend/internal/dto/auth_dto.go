// Package dto defines data transfer objects for API requests and responses
package dto

import (
	"time"
)

// RegisterRequest represents user registration request
type RegisterRequest struct {
	Username        string `json:"username" binding:"required,min=3,max=50" example:"john_doe"`
	Email           string `json:"email" binding:"required,email" example:"john@example.com"`
	Password        string `json:"password" binding:"required,min=8,max=128" example:"SecurePassword123!"`
	ConfirmPassword string `json:"confirmPassword" binding:"required" example:"SecurePassword123!"`
	RealName        string `json:"realName,omitempty" binding:"max=100" example:"张三"`
	Phone           string `json:"phone,omitempty" binding:"max=20" example:"+86-13800138000"`
	Language        string `json:"language,omitempty" binding:"max=10" example:"zh-CN"`
	Timezone        string `json:"timezone,omitempty" binding:"max=50" example:"Asia/Shanghai"`
	TermsAccepted   bool   `json:"termsAccepted" binding:"required" example:"true"`
}

// LoginRequest represents user login request
type LoginRequest struct {
	Username   string `json:"username" binding:"required" example:"john_doe"`
	Password   string `json:"password" binding:"required" example:"SecurePassword123!"`
	RememberMe bool   `json:"rememberMe,omitempty" example:"false"`
	DeviceInfo string `json:"deviceInfo,omitempty" example:"Chrome on Windows"`
}

// LoginResponse represents successful login response
type LoginResponse struct {
	User         UserResponse     `json:"user"`
	AccessToken  string           `json:"accessToken"`
	RefreshToken string           `json:"refreshToken"`
	ExpiresAt    time.Time        `json:"expiresAt"`
	TokenType    string           `json:"tokenType"`
	SessionInfo  *SessionResponse `json:"sessionInfo,omitempty"`
}

// RefreshTokenRequest represents token refresh request
type RefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken" binding:"required" example:"refresh_token_here"`
	DeviceInfo   string `json:"deviceInfo,omitempty" example:"Chrome on Windows"`
}

// TokenResponse represents token response
type TokenResponse struct {
	AccessToken  string    `json:"accessToken"`
	RefreshToken string    `json:"refreshToken"`
	ExpiresAt    time.Time `json:"expiresAt"`
	TokenType    string    `json:"tokenType"`
}

// UserResponse represents user information in responses
type UserResponse struct {
	ID            int64  `json:"id"`
	Username      string `json:"username"`
	Email         string `json:"email"`
	RealName      string `json:"realName"`
	Phone         string `json:"phone"`
	Avatar        string `json:"avatar"`
	Status        string `json:"status"`
	Role          string `json:"role"`
	EmailVerified bool   `json:"emailVerified"`
	// EmailVerifiedAt removed - 已从数据库设计中移除
	TwoFactorEnabled       bool       `json:"twoFactorEnabled"`
	StorageQuota           int64      `json:"storageQuota"`
	StorageUsed            int64      `json:"storageUsed"`
	StorageUsagePercentage float64    `json:"storageUsagePercentage"`
	Language               string     `json:"language"`
	Timezone               string     `json:"timezone"`
	Theme                  string     `json:"theme"`
	LastLoginAt            *time.Time `json:"lastLoginAt"`
	LastLoginIP            string     `json:"lastLoginIP"`
	CreatedAt              time.Time  `json:"createdAt"`
	UpdatedAt              time.Time  `json:"updatedAt"`
}

// ChangePasswordRequest represents password change request
type ChangePasswordRequest struct {
	CurrentPassword string `json:"currentPassword" binding:"required" example:"OldPassword123!"`
	NewPassword     string `json:"newPassword" binding:"required,min=8,max=128" example:"NewPassword123!"`
	ConfirmPassword string `json:"confirmPassword" binding:"required" example:"NewPassword123!"`
}

// ForgotPasswordRequest represents forgot password request
type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email" example:"john@example.com"`
}

// ResetPasswordRequest represents password reset request
type ResetPasswordRequest struct {
	Email           string `json:"email" binding:"required,email" example:"john@example.com"`
	ResetCode       string `json:"resetCode" binding:"required,len=6" example:"123456"`
	NewPassword     string `json:"newPassword" binding:"required,min=8,max=128" example:"NewPassword123!"`
	ConfirmPassword string `json:"confirmPassword" binding:"required" example:"NewPassword123!"`
}

// VerifyEmailRequest represents email verification request
type VerifyEmailRequest struct {
	Email            string `json:"email" binding:"required,email" example:"john@example.com"`
	VerificationCode string `json:"verificationCode" binding:"required,len=6" example:"123456"`
}

// ResendVerificationRequest represents resend verification code request
type ResendVerificationRequest struct {
	Email string `json:"email" binding:"required,email" example:"john@example.com"`
}

// SendVerificationCodeRequest represents send verification code request
type SendVerificationCodeRequest struct {
	Email   string `json:"email" binding:"required,email" example:"john@example.com"`
	Purpose string `json:"purpose" binding:"required,oneof=registration password_reset" example:"registration"`
}

// VerificationCodeResponse represents verification code operation response
type VerificationCodeResponse struct {
	Email         string    `json:"email"`
	Purpose       string    `json:"purpose"`
	ExpiresAt     time.Time `json:"expiresAt"`
	MaxAttempts   int       `json:"maxAttempts"`
	TimeRemaining int       `json:"timeRemaining"` // seconds
	Message       string    `json:"message"`
}

// SessionResponse represents user session information
type SessionResponse struct {
	ID         int64     `json:"id"`
	DeviceInfo string    `json:"deviceInfo"`
	IPAddress  string    `json:"ipAddress"`
	UserAgent  string    `json:"userAgent"`
	IsActive   bool      `json:"isActive"`
	IsCurrent  bool      `json:"isCurrent"`
	ExpiresAt  time.Time `json:"expiresAt"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

// LogoutRequest represents logout request
type LogoutRequest struct {
	AllDevices bool `json:"allDevices,omitempty" example:"false"`
}

// ValidateTokenResponse represents token validation response
type ValidateTokenResponse struct {
	Valid     bool         `json:"valid"`
	User      UserResponse `json:"user,omitempty"`
	ExpiresAt *time.Time   `json:"expiresAt,omitempty"`
}

// UpdateProfileRequest represents user profile update request
type UpdateProfileRequest struct {
	RealName string `json:"realName,omitempty" binding:"max=100" example:"张三"`
	Phone    string `json:"phone,omitempty" binding:"max=20" example:"+86-13800138000"`
	Language string `json:"language,omitempty" binding:"max=10" example:"zh-CN"`
	Timezone string `json:"timezone,omitempty" binding:"max=50" example:"Asia/Shanghai"`
	Theme    string `json:"theme,omitempty" binding:"max=20" example:"dark"`
}

// Enable2FARequest represents enable two-factor authentication request
type Enable2FARequest struct {
	Password string `json:"password" binding:"required" example:"Password123!"`
	TOTPCode string `json:"totpCode" binding:"required,len=6" example:"123456"`
}

// Disable2FARequest represents disable two-factor authentication request
type Disable2FARequest struct {
	Password string `json:"password" binding:"required" example:"Password123!"`
	TOTPCode string `json:"totpCode" binding:"required,len=6" example:"123456"`
}

// Verify2FARequest represents two-factor authentication verification request
type Verify2FARequest struct {
	Username  string `json:"username" binding:"required" example:"john_doe"`
	TOTPCode  string `json:"totpCode" binding:"required,len=6" example:"123456"`
	TempToken string `json:"tempToken" binding:"required" example:"temp_token_here"`
}
