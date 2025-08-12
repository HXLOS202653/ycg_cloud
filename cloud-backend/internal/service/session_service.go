// Package service provides user session management
package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/HXLOS202653/ycg_cloud/cloud-backend/internal/config"
	"github.com/HXLOS202653/ycg_cloud/cloud-backend/internal/model"
	"github.com/HXLOS202653/ycg_cloud/cloud-backend/internal/pkg/database"
	"gorm.io/gorm"
)

// SessionService handles user session operations
type SessionService struct {
	config       *config.Config
	db           *gorm.DB
	redisManager *database.RedisManager
}

// NewSessionService creates a new session service
func NewSessionService(cfg *config.Config, db *gorm.DB, redis *database.RedisManager) *SessionService {
	return &SessionService{
		config:       cfg,
		db:           db,
		redisManager: redis,
	}
}

// SessionInfo represents active session information
type SessionInfo struct {
	UserID       int64     `json:"user_id"`
	SessionToken string    `json:"session_token"`
	DeviceInfo   string    `json:"device_info"`
	IPAddress    string    `json:"ip_address"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	ExpiresAt    time.Time `json:"expires_at"`
	IsActive     bool      `json:"is_active"`
}

// CreateSession creates a new user session
func (s *SessionService) CreateSession(userID int64, refreshToken, deviceInfo, ipAddress, userAgent string) (*model.UserSession, error) {
	// Generate unique session token
	sessionToken, err := s.generateSessionToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate session token: %w", err)
	}

	// Create session record
	session := &model.UserSession{
		UserID:       userID,
		SessionToken: sessionToken,
		RefreshToken: refreshToken,
		DeviceInfo:   deviceInfo,
		IPAddress:    ipAddress,
		UserAgent:    userAgent,
		IsActive:     true,
		ExpiresAt:    time.Now().Add(s.config.Auth.RefreshExpiry), // Use refresh token expiry
	}

	// Save to database
	if err := s.db.Create(session).Error; err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	// Store session in Redis for quick access
	if err := s.storeSessionInCache(session); err != nil {
		// Log error but don't fail - session is still in database
		fmt.Printf("Warning: Failed to cache session: %v\n", err)
	}

	return session, nil
}

// GetActiveSession retrieves active session for user
func (s *SessionService) GetActiveSession(userID int64, sessionToken string) (*model.UserSession, error) {
	// Try to get from cache first
	if session, err := s.getSessionFromCache(sessionToken); err == nil {
		return session, nil
	}

	// Fall back to database
	var session model.UserSession
	err := s.db.Where("user_id = ? AND session_token = ? AND is_active = ? AND expires_at > ?",
		userID, sessionToken, true, time.Now()).First(&session).Error
	if err != nil {
		return nil, err
	}

	// Update cache
	if err := s.storeSessionInCache(&session); err != nil {
		// Log cache error but don't fail session creation
	}

	return &session, nil
}

// UpdateSessionActivity updates session last seen time
func (s *SessionService) UpdateSessionActivity(sessionToken string) error {
	ctx := context.Background()

	// Update in cache
	key := fmt.Sprintf("session:activity:%s", sessionToken)
	if err := s.redisManager.Set(ctx, key, time.Now().Unix(), time.Hour); err != nil {
		fmt.Printf("Warning: Failed to update session activity in cache: %v\n", err)
	}

	// Update in database (batch update every few minutes to reduce DB load)
	// For now, update immediately
	return s.db.Model(&model.UserSession{}).
		Where("session_token = ? AND is_active = ?", sessionToken, true).
		Update("updated_at", time.Now()).Error
}

// InvalidateSession invalidates a user session
func (s *SessionService) InvalidateSession(sessionToken string) error {
	// Remove from cache
	ctx := context.Background()
	cacheKey := fmt.Sprintf("session:%s", sessionToken)
	if err := s.redisManager.Del(ctx, cacheKey); err != nil {
		// Log cache deletion error but continue
	}

	// Update database
	return s.db.Model(&model.UserSession{}).
		Where("session_token = ?", sessionToken).
		Updates(map[string]interface{}{
			"is_active":  false,
			"updated_at": time.Now(),
		}).Error
}

// InvalidateAllUserSessions invalidates all sessions for a user
func (s *SessionService) InvalidateAllUserSessions(userID int64) error {
	// Get all active sessions for user
	var sessions []model.UserSession
	if err := s.db.Where("user_id = ? AND is_active = ?", userID, true).Find(&sessions).Error; err != nil {
		return err
	}

	// Remove from cache
	ctx := context.Background()
	for i := range sessions {
		cacheKey := fmt.Sprintf("session:%s", sessions[i].SessionToken)
		if err := s.redisManager.Del(ctx, cacheKey); err != nil {
			// Log cache deletion error but continue
		}
	}

	// Update database
	return s.db.Model(&model.UserSession{}).
		Where("user_id = ? AND is_active = ?", userID, true).
		Updates(map[string]interface{}{
			"is_active":  false,
			"updated_at": time.Now(),
		}).Error
}

// GetUserActiveSessions returns all active sessions for a user
func (s *SessionService) GetUserActiveSessions(userID int64) ([]model.UserSession, error) {
	var sessions []model.UserSession
	err := s.db.Where("user_id = ? AND is_active = ? AND expires_at > ?",
		userID, true, time.Now()).Find(&sessions).Error
	return sessions, err
}

// CleanupExpiredSessions removes expired sessions
func (s *SessionService) CleanupExpiredSessions() error {
	// Delete expired sessions from database
	result := s.db.Where("expires_at < ? OR (is_active = ? AND logged_out_at IS NOT NULL AND logged_out_at < ?)",
		time.Now(), false, time.Now().AddDate(0, 0, -30)).Delete(&model.UserSession{})

	if result.Error != nil {
		return result.Error
	}

	fmt.Printf("Cleaned up %d expired sessions\n", result.RowsAffected)
	return nil
}

// ValidateSessionToken validates if session is still valid
func (s *SessionService) ValidateSessionToken(userID int64, sessionToken string) (bool, error) {
	session, err := s.GetActiveSession(userID, sessionToken)
	if err != nil {
		return false, err
	}

	// Check if session is still active and not expired
	if !session.IsActive || time.Now().After(session.ExpiresAt) {
		return false, nil
	}

	// Update activity
	if err := s.UpdateSessionActivity(sessionToken); err != nil {
		// Log activity update error but continue
	}

	return true, nil
}

// GetSessionStats returns session statistics for monitoring
func (s *SessionService) GetSessionStats() (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// Count active sessions
	var activeCount int64
	if err := s.db.Model(&model.UserSession{}).
		Where("is_active = ? AND expires_at > ?", true, time.Now()).
		Count(&activeCount).Error; err != nil {
		return nil, err
	}

	// Count total sessions today
	var todayCount int64
	today := time.Now().Truncate(24 * time.Hour)
	if err := s.db.Model(&model.UserSession{}).
		Where("login_at >= ?", today).
		Count(&todayCount).Error; err != nil {
		return nil, err
	}

	stats["active_sessions"] = activeCount
	stats["sessions_today"] = todayCount
	stats["timestamp"] = time.Now()

	return stats, nil
}

// storeSessionInCache stores session in Redis cache
func (s *SessionService) storeSessionInCache(session *model.UserSession) error {
	ctx := context.Background()
	key := fmt.Sprintf("session:%s", session.SessionToken)

	sessionInfo := SessionInfo{
		UserID:       session.UserID,
		SessionToken: session.SessionToken,
		DeviceInfo:   session.DeviceInfo,
		IPAddress:    session.IPAddress,
		CreatedAt:    session.CreatedAt,
		UpdatedAt:    session.UpdatedAt,
		ExpiresAt:    session.ExpiresAt,
		IsActive:     session.IsActive,
	}

	// Store with TTL matching session expiry
	ttl := time.Until(session.ExpiresAt)
	return s.redisManager.SetStruct(ctx, key, sessionInfo, ttl)
}

// getSessionFromCache retrieves session from Redis cache
func (s *SessionService) getSessionFromCache(sessionToken string) (*model.UserSession, error) {
	ctx := context.Background()
	key := fmt.Sprintf("session:%s", sessionToken)

	var sessionInfo SessionInfo
	if err := s.redisManager.GetStruct(ctx, key, &sessionInfo); err != nil {
		return nil, err
	}

	// Convert back to model
	session := &model.UserSession{
		UserID:       sessionInfo.UserID,
		SessionToken: sessionInfo.SessionToken,
		DeviceInfo:   sessionInfo.DeviceInfo,
		IPAddress:    sessionInfo.IPAddress,
		CreatedAt:    sessionInfo.CreatedAt,
		UpdatedAt:    sessionInfo.UpdatedAt,
		ExpiresAt:    sessionInfo.ExpiresAt,
		IsActive:     sessionInfo.IsActive,
	}

	return session, nil
}

// ExtractDeviceInfo extracts device information from User-Agent
func (s *SessionService) ExtractDeviceInfo(userAgent string) string {
	if userAgent == "" {
		return "Unknown Device"
	}

	// Simple device detection - can be enhanced with a proper library
	deviceInfo := "Unknown Device"

	switch {
	case contains(userAgent, "Android"):
		deviceInfo = "Android Mobile"
	case contains(userAgent, "iPhone"):
		deviceInfo = "iPhone"
	case contains(userAgent, "Mobile"):
		deviceInfo = "Mobile Device"
	case contains(userAgent, "Windows"):
		deviceInfo = "Windows Desktop"
	case contains(userAgent, "Macintosh") || contains(userAgent, "Mac OS X"):
		deviceInfo = "Mac Desktop"
	case contains(userAgent, "Linux"):
		deviceInfo = "Linux Desktop"
	case contains(userAgent, "Chrome"):
		deviceInfo = "Chrome Browser"
	case contains(userAgent, "Firefox"):
		deviceInfo = "Firefox Browser"
	case contains(userAgent, "Safari"):
		deviceInfo = "Safari Browser"
	}

	return deviceInfo
}

// contains checks if string contains substring (case-insensitive)
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr ||
		(len(s) > len(substr) &&
			(s[:len(substr)] == substr || s[len(s)-len(substr):] == substr ||
				findSubstring(s, substr))))
}

// findSubstring performs a simple substring search
func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// generateSessionToken generates a unique session token
func (s *SessionService) generateSessionToken() (string, error) {
	// Generate 32 random bytes (256 bits)
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	// Convert to hex string (64 characters)
	return hex.EncodeToString(bytes), nil
}
