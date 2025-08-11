// Package database provides Redis 7.0+ database connection.
package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"

	"github.com/HXLOS202653/ycg_cloud/cloud-backend/internal/config"
)

// RedisManager manages Redis database connections.
type RedisManager struct {
	client *redis.Client
	config config.RedisConfig
}

// NewRedisConnection creates a new Redis 7.0+ connection.
func NewRedisConnection(cfg *config.RedisConfig) (*redis.Client, error) {
	// Build Redis options
	options := buildRedisOptions(cfg)

	// Create Redis client
	client := redis.NewClient(options)

	// Test the connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	log.Printf("Redis 7.0+ connected successfully to %s:%d", cfg.Host, cfg.Port)
	return client, nil
}

// buildRedisOptions builds Redis connection options.
func buildRedisOptions(cfg *config.RedisConfig) *redis.Options {
	options := &redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	}

	// Set username if provided (Redis 6.0+ ACL support)
	if cfg.Username != "" {
		options.Username = cfg.Username
	}

	// Set connection pool options
	if cfg.PoolSize > 0 {
		options.PoolSize = cfg.PoolSize
	} else {
		options.PoolSize = 20 // Default
	}

	if cfg.MinIdleConns > 0 {
		options.MinIdleConns = cfg.MinIdleConns
	} else {
		options.MinIdleConns = 5 // Default
	}

	if cfg.MaxConnAge > 0 {
		options.MaxConnAge = cfg.MaxConnAge
	} else {
		options.MaxConnAge = time.Hour // Default
	}

	if cfg.PoolTimeout > 0 {
		options.PoolTimeout = cfg.PoolTimeout
	} else {
		options.PoolTimeout = 4 * time.Second // Default
	}

	if cfg.IdleTimeout > 0 {
		options.IdleTimeout = cfg.IdleTimeout
	} else {
		options.IdleTimeout = 5 * time.Minute // Default
	}

	if cfg.IdleCheckFreq > 0 {
		options.IdleCheckFrequency = cfg.IdleCheckFreq
	} else {
		options.IdleCheckFrequency = time.Minute // Default
	}

	// Set timeout options
	if cfg.DialTimeout > 0 {
		options.DialTimeout = cfg.DialTimeout
	} else {
		options.DialTimeout = 5 * time.Second // Default
	}

	if cfg.ReadTimeout > 0 {
		options.ReadTimeout = cfg.ReadTimeout
	} else {
		options.ReadTimeout = 3 * time.Second // Default
	}

	if cfg.WriteTimeout > 0 {
		options.WriteTimeout = cfg.WriteTimeout
	} else {
		options.WriteTimeout = 3 * time.Second // Default
	}

	log.Printf("Redis connection pool configured: PoolSize=%d, MinIdle=%d, MaxConnAge=%v",
		options.PoolSize, options.MinIdleConns, options.MaxConnAge)

	return options
}

// NewRedisManager creates a new Redis manager.
func NewRedisManager(cfg *config.RedisConfig) (*RedisManager, error) {
	client, err := NewRedisConnection(cfg)
	if err != nil {
		return nil, err
	}

	return &RedisManager{
		client: client,
		config: *cfg,
	}, nil
}

// GetClient returns the Redis client.
func (r *RedisManager) GetClient() *redis.Client {
	return r.client
}

// Close closes the Redis connection.
func (r *RedisManager) Close() error {
	return r.client.Close()
}

// Health checks the Redis health.
func (r *RedisManager) Health() error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := r.client.Ping(ctx).Result()
	return err
}

// GetStats returns Redis connection statistics.
func (r *RedisManager) GetStats() *redis.PoolStats {
	return r.client.PoolStats()
}

// Redis 7.0+ specific key patterns for cloud storage system.
const (
	// Cache key patterns
	KeyUserSession   = "session:user:%s"
	KeyUserProfile   = "profile:user:%d"
	KeyFileMetadata  = "file:metadata:%d"
	KeyFolderTree    = "folder:tree:%d"
	KeyShareLink     = "share:link:%s"
	KeyUploadSession = "upload:session:%s"
	KeyDownloadToken = "download:token:%s"

	// Lock key patterns
	LockFileOperation = "lock:file:op:%d"
	LockUserOperation = "lock:user:op:%d"
	LockUpload        = "lock:upload:%s"

	// Rate limiting key patterns
	RateLimitAPI      = "ratelimit:api:%s:%s"
	RateLimitUpload   = "ratelimit:upload:%d"
	RateLimitDownload = "ratelimit:download:%d"

	// Queue key patterns
	QueueFileProcess  = "queue:file:process"
	QueueNotification = "queue:notification"
	QueueThumbnail    = "queue:thumbnail"
	QueueVirusScan    = "queue:virus:scan"

	// PubSub channel patterns
	ChannelFileUpdate   = "channel:file:update"
	ChannelUserOnline   = "channel:user:online"
	ChannelNotification = "channel:notification:%d"
)

// SetUserSession sets user session with expiration.
func (r *RedisManager) SetUserSession(ctx context.Context, userID string, sessionData interface{}, expiration time.Duration) error {
	key := fmt.Sprintf(KeyUserSession, userID)
	return r.client.Set(ctx, key, sessionData, expiration).Err()
}

// GetUserSession gets user session.
func (r *RedisManager) GetUserSession(ctx context.Context, userID string) (string, error) {
	key := fmt.Sprintf(KeyUserSession, userID)
	return r.client.Get(ctx, key).Result()
}

// DeleteUserSession deletes user session.
func (r *RedisManager) DeleteUserSession(ctx context.Context, userID string) error {
	key := fmt.Sprintf(KeyUserSession, userID)
	return r.client.Del(ctx, key).Err()
}

// SetFileMetadataCache sets file metadata cache.
func (r *RedisManager) SetFileMetadataCache(ctx context.Context, fileID int, metadata interface{}, expiration time.Duration) error {
	key := fmt.Sprintf(KeyFileMetadata, fileID)
	return r.client.Set(ctx, key, metadata, expiration).Err()
}

// GetFileMetadataCache gets file metadata from cache.
func (r *RedisManager) GetFileMetadataCache(ctx context.Context, fileID int) (string, error) {
	key := fmt.Sprintf(KeyFileMetadata, fileID)
	return r.client.Get(ctx, key).Result()
}

// AcquireLock acquires a distributed lock.
func (r *RedisManager) AcquireLock(ctx context.Context, lockKey string, expiration time.Duration) (bool, error) {
	result, err := r.client.SetNX(ctx, lockKey, "locked", expiration).Result()
	return result, err
}

// ReleaseLock releases a distributed lock.
func (r *RedisManager) ReleaseLock(ctx context.Context, lockKey string) error {
	return r.client.Del(ctx, lockKey).Err()
}
