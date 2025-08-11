// Package database provides MongoDB 6.0 database connection.
package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"github.com/HXLOS202653/ycg_cloud/cloud-backend/internal/config"
)

// MongoDBManager manages MongoDB database connections.
type MongoDBManager struct {
	client   *mongo.Client
	database *mongo.Database
	config   config.MongoDBConfig
}

// NewMongoDBConnection creates a new MongoDB 6.0 connection.
func NewMongoDBConnection(cfg *config.MongoDBConfig) (*mongo.Client, error) {
	// Build connection options
	clientOptions := buildMongoDBOptions(cfg)

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), cfg.ConnectTimeout)
	defer cancel()

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// Test the connection
	ctx, cancel = context.WithTimeout(context.Background(), cfg.ServerSelTimeout)
	defer cancel()

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	log.Printf("MongoDB 6.0 connected successfully to %s/%s", cfg.URI, cfg.Database)
	return client, nil
}

// buildMongoDBOptions builds MongoDB connection options.
func buildMongoDBOptions(cfg *config.MongoDBConfig) *options.ClientOptions {
	clientOptions := options.Client().ApplyURI(cfg.URI)

	// Set authentication if provided
	if cfg.Username != "" && cfg.Password != "" {
		credential := options.Credential{
			Username: cfg.Username,
			Password: cfg.Password,
		}
		clientOptions.SetAuth(credential)
	}

	// Set timeouts
	if cfg.ConnectTimeout > 0 {
		clientOptions.SetConnectTimeout(cfg.ConnectTimeout)
	}
	if cfg.SocketTimeout > 0 {
		clientOptions.SetSocketTimeout(cfg.SocketTimeout)
	}
	if cfg.ServerSelTimeout > 0 {
		clientOptions.SetServerSelectionTimeout(cfg.ServerSelTimeout)
	}

	// Set connection pool options
	if cfg.MaxPoolSize > 0 {
		clientOptions.SetMaxPoolSize(cfg.MaxPoolSize)
	}
	if cfg.MinPoolSize > 0 {
		clientOptions.SetMinPoolSize(cfg.MinPoolSize)
	}
	if cfg.MaxIdleTimeMS > 0 {
		idleTime := time.Duration(cfg.MaxIdleTimeMS) * time.Millisecond
		clientOptions.SetMaxConnIdleTime(idleTime)
	}

	// Enable retryable writes for MongoDB 6.0
	clientOptions.SetRetryWrites(true)
	clientOptions.SetRetryReads(true)

	log.Printf("MongoDB connection pool configured: MaxPool=%d, MinPool=%d, MaxIdleTime=%dms",
		cfg.MaxPoolSize, cfg.MinPoolSize, cfg.MaxIdleTimeMS)

	return clientOptions
}

// NewMongoDBManager creates a new MongoDB manager.
func NewMongoDBManager(cfg *config.MongoDBConfig) (*MongoDBManager, error) {
	client, err := NewMongoDBConnection(cfg)
	if err != nil {
		return nil, err
	}

	database := client.Database(cfg.Database)

	return &MongoDBManager{
		client:   client,
		database: database,
		config:   *cfg,
	}, nil
}

// GetClient returns the MongoDB client.
func (m *MongoDBManager) GetClient() *mongo.Client {
	return m.client
}

// GetDatabase returns the MongoDB database.
func (m *MongoDBManager) GetDatabase() *mongo.Database {
	return m.database
}

// GetCollection returns a MongoDB collection.
func (m *MongoDBManager) GetCollection(name string) *mongo.Collection {
	return m.database.Collection(name)
}

// Close closes the MongoDB connection.
func (m *MongoDBManager) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return m.client.Disconnect(ctx)
}

// Health checks the MongoDB health.
func (m *MongoDBManager) Health() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return m.client.Ping(ctx, readpref.Primary())
}

// MongoDB 6.0 specific collection definitions for cloud storage system.
const (
	// Collections for file metadata and content
	CollectionFileMetadata  = "file_metadata"
	CollectionChatMessages  = "chat_messages"
	CollectionActivityLogs  = "activity_logs"
	CollectionSearchIndexes = "search_indexes"
	CollectionSyncRecords   = "sync_records"
	CollectionAnalytics     = "analytics"
	CollectionNotifications = "notifications"
	CollectionUserSessions  = "user_sessions"
)

// GetFileMetadataCollection returns file metadata collection.
func (m *MongoDBManager) GetFileMetadataCollection() *mongo.Collection {
	return m.GetCollection(CollectionFileMetadata)
}

// GetChatMessagesCollection returns chat messages collection.
func (m *MongoDBManager) GetChatMessagesCollection() *mongo.Collection {
	return m.GetCollection(CollectionChatMessages)
}

// GetActivityLogsCollection returns activity logs collection.
func (m *MongoDBManager) GetActivityLogsCollection() *mongo.Collection {
	return m.GetCollection(CollectionActivityLogs)
}

// GetSearchIndexesCollection returns search indexes collection.
func (m *MongoDBManager) GetSearchIndexesCollection() *mongo.Collection {
	return m.GetCollection(CollectionSearchIndexes)
}

// GetSyncRecordsCollection returns sync records collection.
func (m *MongoDBManager) GetSyncRecordsCollection() *mongo.Collection {
	return m.GetCollection(CollectionSyncRecords)
}

// GetAnalyticsCollection returns analytics collection.
func (m *MongoDBManager) GetAnalyticsCollection() *mongo.Collection {
	return m.GetCollection(CollectionAnalytics)
}
