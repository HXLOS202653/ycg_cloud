//go:build integration
// +build integration

package tests

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"cloud-backend/internal/app"
	"cloud-backend/internal/config"
	"cloud-backend/internal/pkg/database"
)

// TestHealthEndpointsIntegration 测试健康检查端点的集成测试
func TestHealthEndpointsIntegration(t *testing.T) {
	// 加载测试配置
	cfg, err := config.LoadConfig()
	require.NoError(t, err)

	// 创建数据库连接
	dbManager, err := database.NewManager(&cfg.Database)
	require.NoError(t, err)
	defer func() {
		if err := dbManager.Close(); err != nil {
			t.Logf("Failed to close database connections: %w", err)
		}
	}()

	// 创建服务器
	server := &app.Server{
		Config: cfg,
		DB:     dbManager,
	}

	// 启动服务器
	go func() {
		if err := server.Start(); err != nil && err != http.ErrServerClosed {
			t.Logf("Server start error: %w", err)
		}
	}()

	// 等待服务器启动
	time.Sleep(2 * time.Second)

	// 测试基础健康检查
	t.Run("Basic Health Check", func(t *testing.T) {
		resp, err := http.Get("http://localhost:8080/health")
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	// 测试存活检查
	t.Run("Liveness Check", func(t *testing.T) {
		resp, err := http.Get("http://localhost:8080/health/live")
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	// 测试就绪检查
	t.Run("Readiness Check", func(t *testing.T) {
		resp, err := http.Get("http://localhost:8080/health/ready")
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	// 测试监控端点
	t.Run("Metrics Endpoint", func(t *testing.T) {
		resp, err := http.Get("http://localhost:8080/metrics")
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	// 关闭服务器
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		t.Logf("Server shutdown error: %w", err)
	}
}

// TestDatabaseConnectionsIntegration 测试数据库连接的集成测试
func TestDatabaseConnectionsIntegration(t *testing.T) {
	// 加载测试配置
	cfg, err := config.LoadConfig()
	require.NoError(t, err)

	// 创建数据库连接
	dbManager, err := database.NewManager(&cfg.Database)
	require.NoError(t, err)
	defer func() {
		if err := dbManager.Close(); err != nil {
			t.Logf("Failed to close database connections: %w", err)
		}
	}()

	// 测试数据库健康检查
	t.Run("Database Health Check", func(t *testing.T) {
		health := dbManager.HealthCheck()

		// 至少应该有一个数据库连接是健康的
		assert.True(t, health["mysql"].(bool) || health["mongodb"].(bool) || health["redis"].(bool),
			"At least one database should be healthy")
	})

	// 测试连接统计
	t.Run("Connection Stats", func(t *testing.T) {
		stats := dbManager.GetConnectionStats()

		// 验证统计信息结构
		assert.Contains(t, stats, "mysql")
		assert.Contains(t, stats, "mongodb")
		assert.Contains(t, stats, "redis")
	})
}
