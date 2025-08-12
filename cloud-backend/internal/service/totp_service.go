// Package service provides business logic services
package service

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"net/url"
	"strings"
	"time"

	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
)

// TOTPService handles TOTP (Time-based One-Time Password) operations
type TOTPService struct {
	issuer string
}

// NewTOTPService creates a new TOTP service
func NewTOTPService(issuer string) *TOTPService {
	return &TOTPService{
		issuer: issuer,
	}
}

// TOTPSetupData contains information needed for TOTP setup
type TOTPSetupData struct {
	Secret      string   `json:"secret"`
	QRCode      string   `json:"qr_code"`
	BackupCodes []string `json:"backup_codes"`
}

// GenerateSecret generates a new TOTP secret for a user
func (s *TOTPService) GenerateSecret(userEmail string) (*TOTPSetupData, error) {
	// Generate TOTP key
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      s.issuer,
		AccountName: userEmail,
		SecretSize:  32,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to generate TOTP key: %w", err)
	}

	// Generate QR code URL
	qrCodeURL := key.URL()

	// Generate backup codes
	backupCodes, err := s.generateBackupCodes()
	if err != nil {
		return nil, fmt.Errorf("failed to generate backup codes: %w", err)
	}

	return &TOTPSetupData{
		Secret:      key.Secret(),
		QRCode:      qrCodeURL,
		BackupCodes: backupCodes,
	}, nil
}

// ValidateCode validates a TOTP code against a secret
func (s *TOTPService) ValidateCode(secret, code string) bool {
	// Validate current time window
	valid := totp.Validate(code, secret)
	if valid {
		return true
	}

	// Also check previous and next time windows for clock skew tolerance
	now := time.Now()

	// Check previous 30-second window
	if validPrev, _ := totp.ValidateCustom(code, secret, now.Add(-30*time.Second), totp.ValidateOpts{
		Period:    30,
		Skew:      0,
		Digits:    otp.DigitsSix,
		Algorithm: otp.AlgorithmSHA1,
	}); validPrev {
		return true
	}

	// Check next 30-second window
	if validNext, _ := totp.ValidateCustom(code, secret, now.Add(30*time.Second), totp.ValidateOpts{
		Period:    30,
		Skew:      0,
		Digits:    otp.DigitsSix,
		Algorithm: otp.AlgorithmSHA1,
	}); validNext {
		return true
	}

	return false
}

// GenerateCode generates a TOTP code for testing purposes (mainly for development)
func (s *TOTPService) GenerateCode(secret string) (string, error) {
	code, err := totp.GenerateCode(secret, time.Now())
	if err != nil {
		return "", fmt.Errorf("failed to generate code: %w", err)
	}
	return code, nil
}

// generateBackupCodes generates backup codes for 2FA recovery
func (s *TOTPService) generateBackupCodes() ([]string, error) {
	codes := make([]string, 10) // Generate 10 backup codes

	for i := 0; i < 10; i++ {
		// Generate 8-character alphanumeric code
		code, err := s.generateBackupCode()
		if err != nil {
			return nil, err
		}
		codes[i] = code
	}

	return codes, nil
}

// generateBackupCode generates a single backup code
func (s *TOTPService) generateBackupCode() (string, error) {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const codeLength = 8

	code := make([]byte, codeLength)
	for i := range code {
		randomIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		code[i] = charset[randomIndex.Int64()]
	}

	// Format as XXXX-XXXX
	return fmt.Sprintf("%s-%s", string(code[:4]), string(code[4:])), nil
}

// ValidateBackupCode validates a backup code format
func (s *TOTPService) ValidateBackupCode(code string) bool {
	// Remove any spaces or dashes for flexibility
	cleanCode := strings.ReplaceAll(strings.ReplaceAll(code, " ", ""), "-", "")

	// Should be exactly 8 alphanumeric characters
	if len(cleanCode) != 8 {
		return false
	}

	// Check if all characters are alphanumeric
	for _, char := range cleanCode {
		if (char < 'A' || char > 'Z') && (char < '0' || char > '9') {
			return false
		}
	}

	return true
}

// ParseOTPURL parses an OTP URL and extracts the secret
func (s *TOTPService) ParseOTPURL(otpURL string) (string, error) {
	u, err := url.Parse(otpURL)
	if err != nil {
		return "", fmt.Errorf("invalid OTP URL: %w", err)
	}

	if u.Scheme != "otpauth" {
		return "", fmt.Errorf("invalid OTP scheme: %s", u.Scheme)
	}

	secret := u.Query().Get("secret")
	if secret == "" {
		return "", fmt.Errorf("no secret found in OTP URL")
	}

	return secret, nil
}
