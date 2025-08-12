// Package service provides business logic for authentication
package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/HXLOS202653/ycg_cloud/cloud-backend/internal/config"
)

// AuthService handles authentication operations
type AuthService struct {
	config          *config.Config
	passwordService *PasswordService
}

// NewAuthService creates a new authentication service
func NewAuthService(cfg *config.Config) *AuthService {
	passwordService := NewPasswordService(cfg)
	return &AuthService{
		config:          cfg,
		passwordService: passwordService,
	}
}

// Claims represents JWT claims structure
type Claims struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// TokenPair represents access and refresh tokens
type TokenPair struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
	TokenType    string    `json:"token_type"`
}

// GenerateTokens generates both access and refresh tokens
func (s *AuthService) GenerateTokens(userID int64, username, email, role string) (*TokenPair, error) {
	// Generate access token
	accessToken, err := s.generateAccessToken(userID, username, email, role)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	// Generate refresh token (longer expiry, simpler claims)
	refreshToken, err := s.generateRefreshToken(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    time.Now().Add(s.config.Auth.TokenExpiry),
		TokenType:    "Bearer",
	}, nil
}

// generateAccessToken creates a new access token with user claims
func (s *AuthService) generateAccessToken(userID int64, username, email, role string) (string, error) {
	expirationTime := time.Now().Add(s.config.Auth.TokenExpiry)

	claims := &Claims{
		UserID:   userID,
		Username: username,
		Email:    email,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "ycgcloud.com",
			Subject:   fmt.Sprintf("%d", userID),
			ID:        fmt.Sprintf("access_%d_%d", userID, time.Now().Unix()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.config.Auth.JWTSecret))
}

// generateRefreshToken creates a new refresh token
func (s *AuthService) generateRefreshToken(userID int64) (string, error) {
	expirationTime := time.Now().Add(s.config.Auth.RefreshExpiry)

	claims := &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(expirationTime),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		NotBefore: jwt.NewNumericDate(time.Now()),
		Issuer:    "ycgcloud.com",
		Subject:   fmt.Sprintf("%d", userID),
		ID:        fmt.Sprintf("refresh_%d_%d", userID, time.Now().Unix()),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.config.Auth.JWTSecret))
}

// ValidateToken validates and parses a JWT token
func (s *AuthService) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Verify signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.config.Auth.JWTSecret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

// ValidateRefreshToken validates a refresh token and returns user ID
func (s *AuthService) ValidateRefreshToken(tokenString string) (int64, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Verify signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.config.Auth.JWTSecret), nil
	})

	if err != nil {
		return 0, fmt.Errorf("failed to parse refresh token: %w", err)
	}

	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok || !token.Valid {
		return 0, errors.New("invalid refresh token")
	}

	// Extract user ID from subject
	var userID int64
	if _, err := fmt.Sscanf(claims.Subject, "%d", &userID); err != nil {
		return 0, fmt.Errorf("invalid user ID in token: %w", err)
	}

	return userID, nil
}

// HashPassword hashes a password using bcrypt with enhanced security validation
func (s *AuthService) HashPassword(password string) (string, error) {
	return s.passwordService.HashPassword(password)
}

// VerifyPassword verifies a password against its hash
func (s *AuthService) VerifyPassword(hashedPassword, password string) error {
	return s.passwordService.VerifyPassword(hashedPassword, password)
}

// ValidatePassword validates password strength and security requirements
func (s *AuthService) ValidatePassword(password string, userInfo map[string]string) *PasswordValidationResult {
	return s.passwordService.ValidatePassword(password, userInfo)
}

// GetPasswordRequirements returns current password policy
func (s *AuthService) GetPasswordRequirements() *PasswordRequirements {
	return s.passwordService.GetPasswordRequirements()
}

// RefreshAccessToken generates a new access token from a valid refresh token
func (s *AuthService) RefreshAccessToken(refreshToken string, userInfo map[string]interface{}) (*TokenPair, error) {
	// Validate refresh token
	userID, err := s.ValidateRefreshToken(refreshToken)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token: %w", err)
	}

	// Extract user information (would typically come from database)
	username, _ := userInfo["username"].(string)
	email, _ := userInfo["email"].(string)
	role, _ := userInfo["role"].(string)

	// Generate new token pair
	return s.GenerateTokens(userID, username, email, role)
}
