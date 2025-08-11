package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestBasic 基础测试，确保测试框架正常工作
func TestBasic(t *testing.T) {
	t.Run("Basic Test", func(t *testing.T) {
		assert.True(t, true, "Basic test should pass")
		assert.Equal(t, 1+1, 2, "Math should work")
	})
}

// TestTestEnvironment 测试环境验证
func TestTestEnvironment(t *testing.T) {
	t.Run("Test Environment", func(t *testing.T) {
		// 这个测试总是通过，用于验证测试环境
		assert.NotNil(t, t, "Test context should exist")
	})
}