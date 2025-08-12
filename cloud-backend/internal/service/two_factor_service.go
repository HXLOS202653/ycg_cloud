// Package service provides business logic services
package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/HXLOS202653/ycg_cloud/cloud-backend/internal/model"
	"github.com/HXLOS202653/ycg_cloud/cloud-backend/internal/pkg/database"
	"gorm.io/gorm"
)

// TwoFactorService handles two-factor authentication operations
type TwoFactorService struct {
	db           *gorm.DB
	redisManager *database.RedisManager
	totpService  *TOTPService
}

// NewTwoFactorService creates a new two-factor authentication service
func NewTwoFactorService(db *gorm.DB, redisManager *database.RedisManager, totpService *TOTPService) *TwoFactorService {
	return &TwoFactorService{
		db:           db,
		redisManager: redisManager,
		totpService:  totpService,
	}
}

// BackupCode represents a backup code for 2FA recovery
type BackupCode struct {
	Code      string     `json:"code"`
	Used      bool       `json:"used"`
	UsedAt    *time.Time `json:"used_at,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
}

// TwoFactorSetup represents the setup data for 2FA
type TwoFactorSetup struct {
	Secret      string       `json:"secret"`
	QRCodeURL   string       `json:"qr_code_url"`
	BackupCodes []BackupCode `json:"backup_codes"`
}

// Setup2FA generates setup data for enabling 2FA
func (s *TwoFactorService) Setup2FA(userID int64, userEmail string) (*TwoFactorSetup, error) {
	// Check if user already has 2FA enabled
	var user model.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	if user.TwoFactorEnabled {
		return nil, fmt.Errorf("2FA is already enabled for this user")
	}

	// Generate TOTP setup data
	setupData, err := s.totpService.GenerateSecret(userEmail)
	if err != nil {
		return nil, fmt.Errorf("failed to generate TOTP secret: %w", err)
	}

	// Create backup codes
	backupCodes := make([]BackupCode, len(setupData.BackupCodes))
	now := time.Now()

	for i, code := range setupData.BackupCodes {
		backupCodes[i] = BackupCode{
			Code:      code,
			Used:      false,
			CreatedAt: now,
		}
	}

	// Store temporary setup data in Redis (expires in 10 minutes)
	setupKey := fmt.Sprintf("2fa_setup:%d", userID)
	setupInfo := TwoFactorSetup{
		Secret:      setupData.Secret,
		QRCodeURL:   setupData.QRCode,
		BackupCodes: backupCodes,
	}

	ctx := context.Background()
	if err := s.redisManager.SetStruct(ctx, setupKey, setupInfo, 10*time.Minute); err != nil {
		return nil, fmt.Errorf("failed to store setup data: %w", err)
	}

	return &setupInfo, nil
}

// Enable2FA enables 2FA for a user after verifying the TOTP code
func (s *TwoFactorService) Enable2FA(userID int64, totpCode string) error {
	// Get setup data from Redis
	setupKey := fmt.Sprintf("2fa_setup:%d", userID)
	ctx := context.Background()

	var setupData TwoFactorSetup
	if err := s.redisManager.GetStruct(ctx, setupKey, &setupData); err != nil {
		return fmt.Errorf("2FA setup not found or expired")
	}

	// Validate TOTP code
	if !s.totpService.ValidateCode(setupData.Secret, totpCode) {
		return fmt.Errorf("invalid TOTP code")
	}

	// Enable 2FA for user
	if err := s.db.Model(&model.User{}).Where("id = ?", userID).Updates(map[string]interface{}{
		"two_factor_enabled": true,
		"two_factor_secret":  setupData.Secret,
	}).Error; err != nil {
		return fmt.Errorf("failed to enable 2FA: %w", err)
	}

	// Store backup codes
	if err := s.storeBackupCodes(userID, setupData.BackupCodes); err != nil {
		return fmt.Errorf("failed to store backup codes: %w", err)
	}

	// Clean up setup data
	s.redisManager.Del(ctx, setupKey)

	return nil
}

// Disable2FA disables 2FA for a user
func (s *TwoFactorService) Disable2FA(userID int64, _ /* password */, totpCode string) error {
	// Get user
	var user model.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	if !user.TwoFactorEnabled {
		return fmt.Errorf("2FA is not enabled for this user")
	}

	// TODO: Validate password using PasswordService
	// For now, assuming password validation is handled elsewhere

	// Validate TOTP code
	if user.TwoFactorSecret == nil {
		return fmt.Errorf("2FA not properly configured")
	}
	if !s.totpService.ValidateCode(*user.TwoFactorSecret, totpCode) {
		return fmt.Errorf("invalid TOTP code")
	}

	// Disable 2FA
	if err := s.db.Model(&user).Updates(map[string]interface{}{
		"two_factor_enabled": false,
		"two_factor_secret":  nil,
	}).Error; err != nil {
		return fmt.Errorf("failed to disable 2FA: %w", err)
	}

	// Remove backup codes
	ctx := context.Background()
	backupKey := fmt.Sprintf("backup_codes:%d", userID)
	s.redisManager.Del(ctx, backupKey)

	return nil
}

// Verify2FA verifies a 2FA code (TOTP or backup code)
func (s *TwoFactorService) Verify2FA(userID int64, code string) error {
	// Get user
	var user model.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	if !user.TwoFactorEnabled {
		return fmt.Errorf("2FA is not enabled for this user")
	}

	// Clean the code (remove spaces, dashes, etc.)
	cleanCode := strings.ReplaceAll(strings.ReplaceAll(strings.ToUpper(code), " ", ""), "-", "")

	// Try TOTP code first
	if len(cleanCode) == 6 && user.TwoFactorSecret != nil && s.totpService.ValidateCode(*user.TwoFactorSecret, cleanCode) {
		return nil
	}

	// Try backup code
	if len(cleanCode) == 8 && s.totpService.ValidateBackupCode(cleanCode) {
		return s.useBackupCode(userID, cleanCode)
	}

	return fmt.Errorf("invalid 2FA code")
}

// GenerateBackupCodes generates new backup codes for a user
func (s *TwoFactorService) GenerateBackupCodes(userID int64, _ /* password */, totpCode string) ([]string, error) {
	// Get user
	var user model.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	if !user.TwoFactorEnabled {
		return nil, fmt.Errorf("2FA is not enabled for this user")
	}

	// TODO: Validate password using PasswordService
	// For now, assuming password validation is handled elsewhere

	// Validate TOTP code
	if user.TwoFactorSecret == nil {
		return nil, fmt.Errorf("2FA not properly configured")
	}
	if !s.totpService.ValidateCode(*user.TwoFactorSecret, totpCode) {
		return nil, fmt.Errorf("invalid TOTP code")
	}

	// Generate new backup codes
	codes, err := s.totpService.generateBackupCodes()
	if err != nil {
		return nil, fmt.Errorf("failed to generate backup codes: %w", err)
	}

	// Create backup code objects
	backupCodes := make([]BackupCode, len(codes))
	now := time.Now()

	for i, code := range codes {
		backupCodes[i] = BackupCode{
			Code:      code,
			Used:      false,
			CreatedAt: now,
		}
	}

	// Store backup codes
	if err := s.storeBackupCodes(userID, backupCodes); err != nil {
		return nil, fmt.Errorf("failed to store backup codes: %w", err)
	}

	return codes, nil
}

// GetBackupCodes returns unused backup codes for a user
func (s *TwoFactorService) GetBackupCodes(userID int64) ([]string, error) {
	backupCodes, err := s.getBackupCodes(userID)
	if err != nil {
		return nil, err
	}

	var unusedCodes []string
	for _, bc := range backupCodes {
		if !bc.Used {
			unusedCodes = append(unusedCodes, bc.Code)
		}
	}

	return unusedCodes, nil
}

// storeBackupCodes stores backup codes in Redis
func (s *TwoFactorService) storeBackupCodes(userID int64, codes []BackupCode) error {
	ctx := context.Background()
	backupKey := fmt.Sprintf("backup_codes:%d", userID)

	// Store for 1 year (backup codes don't expire)
	return s.redisManager.SetStruct(ctx, backupKey, codes, 365*24*time.Hour)
}

// getBackupCodes retrieves backup codes from Redis
func (s *TwoFactorService) getBackupCodes(userID int64) ([]BackupCode, error) {
	ctx := context.Background()
	backupKey := fmt.Sprintf("backup_codes:%d", userID)

	var codes []BackupCode
	if err := s.redisManager.GetStruct(ctx, backupKey, &codes); err != nil {
		return nil, fmt.Errorf("backup codes not found: %w", err)
	}

	return codes, nil
}

// useBackupCode marks a backup code as used
func (s *TwoFactorService) useBackupCode(userID int64, code string) error {
	codes, err := s.getBackupCodes(userID)
	if err != nil {
		return err
	}

	// Find and mark the code as used
	found := false
	for i := range codes {
		if codes[i].Code != code || codes[i].Used {
			continue
		}
		now := time.Now()
		codes[i].Used = true
		codes[i].UsedAt = &now
		found = true
		break
	}

	if !found {
		return fmt.Errorf("backup code not found or already used")
	}

	// Update stored codes
	return s.storeBackupCodes(userID, codes)
}

// IsEnabled checks if 2FA is enabled for a user
func (s *TwoFactorService) IsEnabled(userID int64) (bool, error) {
	var user model.User
	if err := s.db.Select("two_factor_enabled").First(&user, userID).Error; err != nil {
		return false, fmt.Errorf("user not found: %w", err)
	}

	return user.TwoFactorEnabled, nil
}
