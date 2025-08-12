// Package migration provides validation functionality for database migrations.
package migration

import (
	"fmt"
	"strings"
)

// Validator validates migration files and operations.
type Validator struct {
	config *ValidationConfig
}

// NewValidator creates a new migration validator.
func NewValidator(config *ValidationConfig) *Validator {
	return &Validator{
		config: config,
	}
}

// ValidateMigration validates a migration file.
func (v *Validator) ValidateMigration(migration *MySQLMigration) error {
	if migration == nil {
		return fmt.Errorf("migration cannot be nil")
	}

	// Validate version format (based on existing pattern)
	if err := v.validateVersion(migration.Version); err != nil {
		return fmt.Errorf("invalid migration version: %w", err)
	}

	// Validate migration name (using existing pattern)
	if err := v.validateMigrationName(migration.Name); err != nil {
		return fmt.Errorf("invalid migration name: %w", err)
	}

	return nil
}

// ValidateSQL validates SQL content.
func (v *Validator) ValidateSQL(sql string) error {
	if strings.TrimSpace(sql) == "" {
		return fmt.Errorf("SQL content cannot be empty")
	}

	// Basic SQL validation (extend based on config)
	if v.config != nil && v.config.Enabled {
		return v.validateSQLContent(sql)
	}

	return nil
}

// validateVersion validates migration version format (based on existing pattern)
func (v *Validator) validateVersion(version string) error {
	if len(version) != 14 {
		return fmt.Errorf("version must be 14 characters (YYYYMMDDHHMMSS format)")
	}

	for _, char := range version {
		if char < '0' || char > '9' {
			return fmt.Errorf("version must contain only digits")
		}
	}

	return nil
}

// validateMigrationName validates migration name (based on existing pattern from main.go)
func (v *Validator) validateMigrationName(name string) error {
	if name == "" {
		return fmt.Errorf("migration name cannot be empty")
	}

	// Check for special characters (same logic as main.go)
	for _, char := range name {
		if !((char >= 'a' && char <= 'z') ||
			(char >= 'A' && char <= 'Z') ||
			(char >= '0' && char <= '9') ||
			char == '_') {
			return fmt.Errorf("migration name can only contain letters, numbers and underscores")
		}
	}

	return nil
}

// validateSQLContent validates SQL content based on configuration
func (v *Validator) validateSQLContent(sql string) error {
	if v.config.StrictMode {
		// Check for forbidden operations
		upperSQL := strings.ToUpper(sql)
		for _, forbidden := range v.config.ForbiddenOperations {
			if strings.Contains(upperSQL, strings.ToUpper(forbidden)) {
				return fmt.Errorf("forbidden SQL operation detected: %s", forbidden)
			}
		}

		// Check for allowed operations only
		if len(v.config.AllowedOperations) > 0 {
			hasAllowed := false
			for _, allowed := range v.config.AllowedOperations {
				if strings.Contains(upperSQL, strings.ToUpper(allowed)) {
					hasAllowed = true
					break
				}
			}
			if !hasAllowed {
				return fmt.Errorf("SQL must contain at least one allowed operation")
			}
		}
	}

	return nil
}
