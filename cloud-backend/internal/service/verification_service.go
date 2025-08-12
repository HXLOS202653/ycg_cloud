// Package service provides verification code management
package service

import (
	"context"
	"fmt"
	"time"

	"github.com/HXLOS202653/ycg_cloud/cloud-backend/internal/config"
	"github.com/HXLOS202653/ycg_cloud/cloud-backend/internal/pkg/database"
)

// VerificationService handles verification code operations
type VerificationService struct {
	config       *config.Config
	redisManager *database.RedisManager
	emailService *EmailService
}

// NewVerificationService creates a new verification service
func NewVerificationService(cfg *config.Config, redis *database.RedisManager, emailService *EmailService) *VerificationService {
	return &VerificationService{
		config:       cfg,
		redisManager: redis,
		emailService: emailService,
	}
}

// VerificationCode represents a verification code record
type VerificationCode struct {
	Code        string    `json:"code"`
	Email       string    `json:"email"`
	Purpose     string    `json:"purpose"`
	CreatedAt   time.Time `json:"created_at"`
	ExpiresAt   time.Time `json:"expires_at"`
	Attempts    int       `json:"attempts"`
	MaxAttempts int       `json:"max_attempts"`
}

// SendVerificationCode sends a verification code to email
func (s *VerificationService) SendVerificationCode(email, purpose string) error {
	// Validate email format
	if !s.emailService.ValidateEmail(email) {
		return fmt.Errorf("invalid email format")
	}

	// Check rate limiting - prevent spam
	if err := s.checkRateLimit(email, purpose); err != nil {
		return err
	}

	// Generate and send verification code
	code, err := s.emailService.SendVerificationCode(email, purpose)
	if err != nil {
		return fmt.Errorf("failed to send verification code: %w", err)
	}

	// Store verification code in Redis
	if err := s.storeVerificationCode(email, code, purpose); err != nil {
		return fmt.Errorf("failed to store verification code: %w", err)
	}

	return nil
}

// VerifyCode verifies a verification code
func (s *VerificationService) VerifyCode(email, code, purpose string) error {
	// Get stored verification code
	storedCode, err := s.getVerificationCode(email, purpose)
	if err != nil {
		return fmt.Errorf("verification code not found or expired")
	}

	// Check if maximum attempts exceeded
	if storedCode.Attempts >= storedCode.MaxAttempts {
		// Delete the code to prevent further attempts
		if err := s.deleteVerificationCode(email, purpose); err != nil {
			// Log deletion error but continue
			fmt.Printf("Failed to delete verification code after max attempts: %v\n", err)
		}
		return fmt.Errorf("maximum verification attempts exceeded")
	}

	// Increment attempt count
	storedCode.Attempts++
	if err := s.updateVerificationCode(email, purpose, storedCode); err != nil {
		// Log update error but continue verification
		fmt.Printf("Failed to update verification code attempt count: %v\n", err)
	}

	// Verify the code
	if storedCode.Code != code {
		return fmt.Errorf("invalid verification code")
	}

	// Code is valid, delete it to prevent reuse
	if err := s.deleteVerificationCode(email, purpose); err != nil {
		// Log deletion error but don't fail verification
	}

	return nil
}

// checkRateLimit checks if user is sending verification codes too frequently
func (s *VerificationService) checkRateLimit(email, purpose string) error {
	ctx := context.Background()
	key := s.getRateLimitKey(email, purpose)

	// Check current count
	count, err := s.redisManager.Get(ctx, key)
	if err != nil && err.Error() != "redis: nil" {
		return fmt.Errorf("failed to check rate limit: %w", err)
	}

	// Parse count
	var currentCount int
	if count != "" {
		if _, err := fmt.Sscanf(count, "%d", &currentCount); err != nil {
			// If parsing fails, assume 0 count
			currentCount = 0
		}
	}

	// Define rate limits based on purpose
	var maxAttempts int
	var window time.Duration

	switch purpose {
	case "registration":
		maxAttempts = 3 // 3 verification codes per hour for registration
		window = time.Hour
	case "password_reset":
		maxAttempts = 3 // 3 verification codes per hour for password reset
		window = time.Hour
	default:
		maxAttempts = 5 // 5 verification codes per hour for other purposes
		window = time.Hour
	}

	if currentCount >= maxAttempts {
		return fmt.Errorf("too many verification code requests, please try again later")
	}

	// Increment count
	pipe := s.redisManager.Pipeline()
	pipe.Incr(ctx, key)
	pipe.Expire(ctx, key, window)
	_, err = pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to update rate limit: %w", err)
	}

	return nil
}

// storeVerificationCode stores verification code in Redis
func (s *VerificationService) storeVerificationCode(email, code, purpose string) error {
	ctx := context.Background()
	key := s.getVerificationKey(email, purpose)

	// Set expiry based on purpose
	var expiry time.Duration
	var maxAttempts int

	switch purpose {
	case "registration":
		expiry = 5 * time.Minute // Registration codes expire in 5 minutes
		maxAttempts = 3
	case "password_reset":
		expiry = 10 * time.Minute // Password reset codes expire in 10 minutes
		maxAttempts = 3
	default:
		expiry = 5 * time.Minute // Default 5 minutes
		maxAttempts = 3
	}

	verificationCode := VerificationCode{
		Code:        code,
		Email:       email,
		Purpose:     purpose,
		CreatedAt:   time.Now(),
		ExpiresAt:   time.Now().Add(expiry),
		Attempts:    0,
		MaxAttempts: maxAttempts,
	}

	// Store in Redis with expiry
	return s.redisManager.SetEX(ctx, key, verificationCode, expiry)
}

// getVerificationCode retrieves verification code from Redis
func (s *VerificationService) getVerificationCode(email, purpose string) (*VerificationCode, error) {
	ctx := context.Background()
	key := s.getVerificationKey(email, purpose)

	var verificationCode VerificationCode
	err := s.redisManager.GetStruct(ctx, key, &verificationCode)
	if err != nil {
		return nil, err
	}

	// Check if expired
	if time.Now().After(verificationCode.ExpiresAt) {
		if err := s.deleteVerificationCode(email, purpose); err != nil {
			// Log deletion error but continue
		}
		return nil, fmt.Errorf("verification code expired")
	}

	return &verificationCode, nil
}

// updateVerificationCode updates verification code in Redis
func (s *VerificationService) updateVerificationCode(email, purpose string, code *VerificationCode) error {
	ctx := context.Background()
	key := s.getVerificationKey(email, purpose)

	// Calculate remaining TTL
	ttl := time.Until(code.ExpiresAt)
	if ttl <= 0 {
		return fmt.Errorf("verification code expired")
	}

	return s.redisManager.SetEX(ctx, key, code, ttl)
}

// deleteVerificationCode deletes verification code from Redis
func (s *VerificationService) deleteVerificationCode(email, purpose string) error {
	ctx := context.Background()
	key := s.getVerificationKey(email, purpose)
	return s.redisManager.Del(ctx, key)
}

// getVerificationKey generates Redis key for verification code
func (s *VerificationService) getVerificationKey(email, purpose string) string {
	return fmt.Sprintf("verification:%s:%s", purpose, email)
}

// getRateLimitKey generates Redis key for rate limiting
func (s *VerificationService) getRateLimitKey(email, purpose string) string {
	return fmt.Sprintf("rate_limit:verification:%s:%s", purpose, email)
}

// ResendVerificationCode resends verification code if allowed
func (s *VerificationService) ResendVerificationCode(email, purpose string) error {
	// Check if there's an existing code
	existingCode, err := s.getVerificationCode(email, purpose)
	if err == nil {
		// Code exists, check if enough time has passed since last send
		timeSinceCreated := time.Since(existingCode.CreatedAt)
		if timeSinceCreated < time.Minute {
			return fmt.Errorf("please wait %d seconds before requesting a new code",
				60-int(timeSinceCreated.Seconds()))
		}
	}

	// Send new verification code (this will override the existing one)
	return s.SendVerificationCode(email, purpose)
}

// CleanupExpiredCodes removes expired verification codes (can be called by a cleanup job)
func (s *VerificationService) CleanupExpiredCodes() error {
	// This would typically be implemented with Redis key scanning and TTL checking
	// For now, Redis TTL automatically handles cleanup
	return nil
}

// GetVerificationCodeInfo returns non-sensitive info about verification code
func (s *VerificationService) GetVerificationCodeInfo(email, purpose string) (map[string]interface{}, error) {
	code, err := s.getVerificationCode(email, purpose)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"email":              email,
		"purpose":            purpose,
		"created_at":         code.CreatedAt,
		"expires_at":         code.ExpiresAt,
		"attempts":           code.Attempts,
		"max_attempts":       code.MaxAttempts,
		"remaining_attempts": code.MaxAttempts - code.Attempts,
		"time_remaining":     int(time.Until(code.ExpiresAt).Seconds()),
	}, nil
}
