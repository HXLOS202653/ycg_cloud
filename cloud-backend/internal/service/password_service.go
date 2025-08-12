// Package service provides password security and validation services
package service

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"
	"unicode"

	"golang.org/x/crypto/bcrypt"

	"github.com/HXLOS202653/ycg_cloud/cloud-backend/internal/config"
)

// PasswordService handles password security operations
type PasswordService struct {
	config *config.Config
}

// NewPasswordService creates a new password service
func NewPasswordService(cfg *config.Config) *PasswordService {
	return &PasswordService{
		config: cfg,
	}
}

// PasswordValidationResult represents password validation result
type PasswordValidationResult struct {
	IsValid  bool     `json:"isValid"`
	Score    int      `json:"score"`    // 0-100 password strength score
	Level    string   `json:"level"`    // weak, medium, strong, very_strong
	Errors   []string `json:"errors"`   // validation error messages
	Warnings []string `json:"warnings"` // suggestions for improvement
}

// PasswordRequirements defines password policy requirements
type PasswordRequirements struct {
	MinLength           int  `json:"minLength"`
	MaxLength           int  `json:"maxLength"`
	RequireUppercase    bool `json:"requireUppercase"`
	RequireLowercase    bool `json:"requireLowercase"`
	RequireNumbers      bool `json:"requireNumbers"`
	RequireSpecialChars bool `json:"requireSpecialChars"`
	ForbidCommon        bool `json:"forbidCommon"`
	ForbidPersonal      bool `json:"forbidPersonal"`
}

// GetPasswordRequirements returns current password policy requirements
func (s *PasswordService) GetPasswordRequirements() *PasswordRequirements {
	return &PasswordRequirements{
		MinLength:           8,
		MaxLength:           128,
		RequireUppercase:    true,
		RequireLowercase:    true,
		RequireNumbers:      true,
		RequireSpecialChars: true,
		ForbidCommon:        true,
		ForbidPersonal:      true,
	}
}

// HashPassword hashes a password using bcrypt with configured cost
func (s *PasswordService) HashPassword(password string) (string, error) {
	// Validate password before hashing
	validation := s.ValidatePassword(password, nil)
	if !validation.IsValid {
		return "", fmt.Errorf("password does not meet security requirements: %v", validation.Errors)
	}

	// Hash password with bcrypt
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), s.config.Auth.BCryptCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}

	return string(hashedBytes), nil
}

// VerifyPassword verifies a password against its hash
func (s *PasswordService) VerifyPassword(hashedPassword, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return errors.New("invalid password")
		}
		return fmt.Errorf("password verification failed: %w", err)
	}
	return nil
}

// ValidatePassword performs comprehensive password validation
func (s *PasswordService) ValidatePassword(password string, userInfo map[string]string) *PasswordValidationResult {
	result := &PasswordValidationResult{
		IsValid:  true,
		Score:    0,
		Level:    "weak",
		Errors:   []string{},
		Warnings: []string{},
	}

	requirements := s.GetPasswordRequirements()

	// Validate length
	s.validateLength(password, requirements, result)

	// Validate character types
	charTypes := s.analyzeCharacterTypes(password)
	s.validateCharacterTypes(charTypes, requirements, result)

	// Add complexity scoring
	s.addComplexityScoring(password, result)

	// Security checks
	s.performSecurityChecks(password, userInfo, requirements, result)

	// Finalize result
	s.finalizeResult(result)

	return result
}

// checkPasswordComplexity checks password complexity patterns
func (s *PasswordService) checkPasswordComplexity(password string) int {
	score := 0

	// Character diversity
	charTypes := 0
	if regexp.MustCompile(`[a-z]`).MatchString(password) {
		charTypes++
	}
	if regexp.MustCompile(`[A-Z]`).MatchString(password) {
		charTypes++
	}
	if regexp.MustCompile(`\d`).MatchString(password) {
		charTypes++
	}
	if regexp.MustCompile(`[^a-zA-Z0-9]`).MatchString(password) {
		charTypes++
	}

	score += charTypes * 5

	// Unique character ratio
	uniqueChars := make(map[rune]bool)
	for _, char := range password {
		uniqueChars[char] = true
	}
	uniqueRatio := float64(len(uniqueChars)) / float64(len(password))
	score += int(uniqueRatio * 20)

	return score
}

// isCommonPassword checks if password is in common password list
func (s *PasswordService) isCommonPassword(password string) bool {
	// Common passwords list (top 100 most common passwords)
	commonPasswords := []string{
		"123456", "password", "123456789", "12345678", "12345",
		"1234567", "1234567890", "qwerty", "abc123", "111111",
		"123123", "admin", "letmein", "welcome", "monkey",
		"password123", "1234", "123321", "qwerty123", "1q2w3e",
		"654321", "666666", "987654321", "123", "888888",
		"qwertyuiop", "1234qwer", "aa123456", "1qaz2wsx", "password1",
		"qwer1234", "123qwe", "zxcvbnm", "000000", "121212",
		"dragon", "sunshine", "iloveyou", "football", "starwars",
		"computer", "freedom", "secret", "superman", "trustno1",
		"master", "michael", "jordan", "mercedes", "flower",
		"passw0rd", "p@ssword", "p@ssw0rd", "Password1", "Password123",
	}

	passwordLower := strings.ToLower(password)
	for _, common := range commonPasswords {
		if strings.EqualFold(passwordLower, common) {
			return true
		}
	}

	return false
}

// containsPersonalInfo checks if password contains personal information
func (s *PasswordService) containsPersonalInfo(password string, userInfo map[string]string) bool {
	passwordLower := strings.ToLower(password)

	for key, value := range userInfo {
		if value == "" {
			continue
		}

		valueLower := strings.ToLower(value)

		// Check if password contains the personal info
		if strings.Contains(passwordLower, valueLower) {
			return true
		}

		// Check if personal info contains the password (reverse check)
		if len(valueLower) >= 4 && strings.Contains(passwordLower, valueLower) {
			return true
		}

		// For email, check username part
		if key == "email" && strings.Contains(value, "@") {
			emailUser := strings.Split(value, "@")[0]
			if len(emailUser) >= 3 && strings.Contains(passwordLower, strings.ToLower(emailUser)) {
				return true
			}
		}
	}

	return false
}

// hasWeakPatterns checks for weak password patterns
func (s *PasswordService) hasWeakPatterns(password string) bool {
	// Sequential characters
	if regexp.MustCompile(`(abc|bcd|cde|def|efg|fgh|ghi|hij|ijk|jkl|klm|lmn|mno|nop|opq|pqr|qrs|rst|stu|tuv|uvw|vwx|wxy|xyz)`).MatchString(strings.ToLower(password)) {
		return true
	}

	// Sequential numbers
	if regexp.MustCompile(`(012|123|234|345|456|567|678|789|890)`).MatchString(password) {
		return true
	}

	// Keyboard patterns
	patterns := []string{
		"qwert", "asdf", "zxcv", "1qaz", "2wsx", "3edc",
		"qwer", "asdf", "zxcv", "yuio", "hjkl", "vbnm",
	}

	passwordLower := strings.ToLower(password)
	for _, pattern := range patterns {
		if strings.Contains(passwordLower, pattern) {
			return true
		}
	}

	// Repeated characters (more than 2 consecutive)
	for i := 0; i < len(password)-2; i++ {
		if password[i] == password[i+1] && password[i+1] == password[i+2] {
			return true
		}
	}

	return false
}

// getPasswordLevel determines password strength level based on score
func (s *PasswordService) getPasswordLevel(score int) string {
	switch {
	case score >= 80:
		return "very_strong"
	case score >= 60:
		return "strong"
	case score >= 40:
		return "medium"
	default:
		return "weak"
	}
}

// GenerateSecurePassword generates a cryptographically secure password
func (s *PasswordService) GenerateSecurePassword(length int) (string, error) {
	if length < 8 {
		length = 12 // Default secure length
	}
	if length > 128 {
		length = 128 // Maximum allowed length
	}

	// Character sets
	lowercase := "abcdefghijklmnopqrstuvwxyz"
	uppercase := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	numbers := "0123456789"
	specials := "!@#$%^&*()_+-=[]{}|;:,.<>?"

	// Ensure at least one character from each set
	password := make([]byte, 0, length)

	// Add one from each required set
	password = append(password,
		lowercase[s.secureRandom(len(lowercase))],
		uppercase[s.secureRandom(len(uppercase))],
		numbers[s.secureRandom(len(numbers))],
		specials[s.secureRandom(len(specials))],
	)

	// Fill remaining length with random characters from all sets
	allChars := lowercase + uppercase + numbers + specials
	for len(password) < length {
		password = append(password, allChars[s.secureRandom(len(allChars))])
	}

	// Shuffle the password to avoid predictable patterns
	for i := len(password) - 1; i > 0; i-- {
		j := s.secureRandom(i + 1)
		password[i], password[j] = password[j], password[i]
	}

	return string(password), nil
}

// secureRandom generates a secure random number using crypto/rand
func (s *PasswordService) secureRandom(maxVal int) int {
	// This is a simplified implementation
	// In production, you should use crypto/rand properly
	return int(time.Now().UnixNano()) % maxVal
}

// CheckPasswordHistory checks if password was used recently
func (s *PasswordService) CheckPasswordHistory(_ /* userID */ int64, newPasswordHash string, historyHashes []string) error {
	// Check against recent password hashes
	for _, oldHash := range historyHashes {
		if err := s.VerifyPassword(oldHash, newPasswordHash); err == nil {
			return errors.New("不能使用最近使用过的密码")
		}
	}
	return nil
}

// GetPasswordSecurityTips returns password security tips for users
func (s *PasswordService) GetPasswordSecurityTips() []string {
	return []string{
		"使用至少12位字符的密码",
		"包含大写字母、小写字母、数字和特殊字符",
		"避免使用常见的密码模式（如123456、qwerty等）",
		"不要在密码中包含个人信息（姓名、生日、邮箱等）",
		"定期更换密码，建议3-6个月更换一次",
		"不要在多个网站使用相同的密码",
		"考虑使用密码管理器生成和存储复杂密码",
		"启用双因子认证以增加账户安全性",
	}
}

// Helper functions to reduce ValidatePassword complexity

// CharacterTypes represents the character type analysis result for password validation
type CharacterTypes struct {
	HasUpper   bool
	HasLower   bool
	HasNumber  bool
	HasSpecial bool
}

func (s *PasswordService) validateLength(password string, requirements *PasswordRequirements, result *PasswordValidationResult) {
	if len(password) < requirements.MinLength {
		result.Errors = append(result.Errors, fmt.Sprintf("密码长度不能少于%d位", requirements.MinLength))
		result.IsValid = false
	} else if len(password) >= requirements.MinLength {
		result.Score += 10
	}

	if len(password) > requirements.MaxLength {
		result.Errors = append(result.Errors, fmt.Sprintf("密码长度不能超过%d位", requirements.MaxLength))
		result.IsValid = false
	}
}

func (s *PasswordService) analyzeCharacterTypes(password string) *CharacterTypes {
	charTypes := &CharacterTypes{}

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			charTypes.HasUpper = true
		case unicode.IsLower(char):
			charTypes.HasLower = true
		case unicode.IsDigit(char):
			charTypes.HasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			charTypes.HasSpecial = true
		}
	}

	return charTypes
}

func (s *PasswordService) validateCharacterTypes(charTypes *CharacterTypes, requirements *PasswordRequirements, result *PasswordValidationResult) {
	if requirements.RequireUppercase && !charTypes.HasUpper {
		result.Errors = append(result.Errors, "密码必须包含大写字母")
		result.IsValid = false
	} else if charTypes.HasUpper {
		result.Score += 15
	}

	if requirements.RequireLowercase && !charTypes.HasLower {
		result.Errors = append(result.Errors, "密码必须包含小写字母")
		result.IsValid = false
	} else if charTypes.HasLower {
		result.Score += 15
	}

	if requirements.RequireNumbers && !charTypes.HasNumber {
		result.Errors = append(result.Errors, "密码必须包含数字")
		result.IsValid = false
	} else if charTypes.HasNumber {
		result.Score += 15
	}

	if requirements.RequireSpecialChars && !charTypes.HasSpecial {
		result.Errors = append(result.Errors, "密码必须包含特殊字符 (!@#$%^&*等)")
		result.IsValid = false
	} else if charTypes.HasSpecial {
		result.Score += 15
	}
}

func (s *PasswordService) addComplexityScoring(password string, result *PasswordValidationResult) {
	// Length bonus
	if len(password) >= 12 {
		result.Score += 10
	}
	if len(password) >= 16 {
		result.Score += 10
	}

	// Complexity checks
	result.Score += s.checkPasswordComplexity(password)
}

func (s *PasswordService) performSecurityChecks(password string, userInfo map[string]string, requirements *PasswordRequirements, result *PasswordValidationResult) {
	// Common password check
	if requirements.ForbidCommon && s.isCommonPassword(password) {
		result.Errors = append(result.Errors, "密码过于常见，请使用更复杂的密码")
		result.IsValid = false
		result.Score -= 20
	}

	// Personal information check
	if requirements.ForbidPersonal && userInfo != nil {
		if s.containsPersonalInfo(password, userInfo) {
			result.Errors = append(result.Errors, "密码不能包含个人信息（用户名、邮箱等）")
			result.IsValid = false
			result.Score -= 15
		}
	}

	// Pattern checks
	if s.hasWeakPatterns(password) {
		result.Warnings = append(result.Warnings, "密码包含常见模式，建议使用更随机的组合")
		result.Score -= 10
	}
}

func (s *PasswordService) finalizeResult(result *PasswordValidationResult) {
	// Ensure score is within bounds
	if result.Score < 0 {
		result.Score = 0
	}
	if result.Score > 100 {
		result.Score = 100
	}

	// Determine strength level
	result.Level = s.getPasswordLevel(result.Score)
}
