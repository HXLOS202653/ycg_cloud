// Package database provides usage examples for connection pool and transaction management.
package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"gorm.io/gorm"

	"github.com/HXLOS202653/ycg_cloud/cloud-backend/internal/config"
)

// Example demonstrates various database operations using the enhanced features.
func Example() {
	// Initialize database configuration
	cfg := config.GetDefaultDatabaseConfig()

	// Create MySQL manager with enhanced features
	mysqlManager, err := NewMySQLManager(&cfg.MySQL)
	if err != nil {
		log.Fatalf("Failed to create MySQL manager: %v", err)
	}
	defer func() { _ = mysqlManager.Close() }()

	// Example 1: Basic transaction usage
	ExampleBasicTransaction(mysqlManager)

	// Example 2: Transaction with custom options
	ExampleTransactionWithOptions(mysqlManager)

	// Example 3: Batch operations
	ExampleBatchOperations(mysqlManager)

	// Example 4: Connection pool monitoring
	ExamplePoolMonitoring(mysqlManager)

	// Example 5: Retry operations
	ExampleRetryOperations(mysqlManager)

	// Example 6: Savepoints
	ExampleSavepoints(mysqlManager)
}

// ExampleBasicTransaction demonstrates basic transaction usage.
func ExampleBasicTransaction(manager *MySQLManager) {
	log.Println("=== Example: Basic Transaction ===")

	ctx := context.Background()
	transactionMgr := manager.GetTransactionManager()

	err := transactionMgr.WithTransaction(ctx, func(_ *gorm.DB) error {
		// Simulate some database operations
		log.Println("Performing database operations within transaction...")

		// Example operations (replace with actual model operations)
		// if err := tx.Create(&user).Error; err != nil {
		//     return err
		// }
		//
		// if err := tx.Create(&profile).Error; err != nil {
		//     return err
		// }

		// Simulate processing time
		time.Sleep(100 * time.Millisecond)
		log.Println("Transaction operations completed successfully")
		return nil
	})

	if err != nil {
		log.Printf("Transaction failed: %v", err)
	} else {
		log.Println("Transaction committed successfully")
	}
}

// ExampleTransactionWithOptions demonstrates transaction with custom options.
func ExampleTransactionWithOptions(manager *MySQLManager) {
	log.Println("=== Example: Transaction with Custom Options ===")

	ctx := context.Background()
	transactionMgr := manager.GetTransactionManager()

	// Custom transaction options
	opts := &TransactionOptions{
		Isolation: sql.LevelRepeatableRead,
		ReadOnly:  false,
		Timeout:   10 * time.Second,
	}

	err := transactionMgr.WithTransactionOptions(ctx, opts, func(_ *gorm.DB) error {
		log.Println("Performing operations with repeatable read isolation...")

		// Simulate operations that require consistent reads
		time.Sleep(200 * time.Millisecond)
		return nil
	})

	if err != nil {
		log.Printf("Custom transaction failed: %v", err)
	} else {
		log.Println("Custom transaction completed successfully")
	}
}

// ExampleBatchOperations demonstrates batch processing with transactions.
func ExampleBatchOperations(manager *MySQLManager) {
	log.Println("=== Example: Batch Operations ===")

	ctx := context.Background()
	transactionMgr := manager.GetTransactionManager()

	// Create batch operation manager
	batchOp := transactionMgr.NewBatchOperation(100, 30*time.Second)

	// Create sample operations
	operations := make([]func(*gorm.DB) error, 500)
	for i := 0; i < len(operations); i++ {
		operationID := i
		operations[i] = func(_ *gorm.DB) error {
			// Simulate a database operation
			log.Printf("Processing operation %d", operationID)
			time.Sleep(10 * time.Millisecond)
			return nil
		}
	}

	// Execute batch operations
	err := batchOp.ExecuteBatch(ctx, operations)
	if err != nil {
		log.Printf("Batch operations failed: %v", err)
	} else {
		log.Println("All batch operations completed successfully")
	}
}

// ExamplePoolMonitoring demonstrates connection pool monitoring.
func ExamplePoolMonitoring(manager *MySQLManager) {
	log.Println("=== Example: Connection Pool Monitoring ===")

	pool := manager.GetConnectionPool()
	if pool == nil {
		log.Println("Connection pool not available")
		return
	}

	// Get pool metrics
	metrics := pool.GetMetrics()
	log.Printf("Pool Status: %s", metrics.PoolStatus)
	log.Printf("Open Connections: %d", metrics.OpenConnections)
	log.Printf("In Use Connections: %d", metrics.InUseConnections)
	log.Printf("Idle Connections: %d", metrics.IdleConnections)
	log.Printf("Total Queries: %d", metrics.TotalQueries)
	log.Printf("Slow Queries: %d", metrics.SlowQueries)
	log.Printf("Failed Queries: %d", metrics.FailedQueries)
	log.Printf("Average Query Time: %v", metrics.AverageQueryTime)
	log.Printf("Health Check Status: %s", metrics.HealthCheckStatus)
	log.Printf("Last Health Check: %v", metrics.LastHealthCheck)

	// Check if pool is healthy
	if manager.IsPoolHealthy() {
		log.Println("Connection pool is healthy")
	} else {
		log.Println("Connection pool is unhealthy")
	}
}

// ExampleRetryOperations demonstrates retry logic for database operations.
func ExampleRetryOperations(manager *MySQLManager) {
	log.Println("=== Example: Retry Operations ===")

	// Example of operation that might fail and need retry
	err := manager.ExecuteWithRetry(func(_ *gorm.DB) error {
		log.Println("Attempting database operation...")

		// Simulate an operation that might fail occasionally
		// In real scenarios, this could be network timeouts, deadlocks, etc.
		if time.Now().UnixNano()%3 == 0 {
			return fmt.Errorf("simulated temporary failure")
		}

		log.Println("Operation succeeded")
		return nil
	})

	if err != nil {
		log.Printf("Operation failed after retries: %v", err)
	} else {
		log.Println("Operation completed successfully (possibly after retries)")
	}
}

// ExampleSavepoints demonstrates savepoint usage within transactions.
func ExampleSavepoints(manager *MySQLManager) {
	log.Println("=== Example: Savepoints ===")

	ctx := context.Background()
	transactionMgr := manager.GetTransactionManager()

	err := transactionMgr.WithTransaction(ctx, func(tx *gorm.DB) error {
		log.Println("Starting transaction with savepoints...")

		// Create savepoint manager
		spManager := NewSavepointManager(tx)

		// Perform some operations
		log.Println("Performing first set of operations...")

		// Create a savepoint
		sp1, err := spManager.CreateSavepoint()
		if err != nil {
			return fmt.Errorf("failed to create savepoint: %w", err)
		}

		// Perform risky operations
		log.Println("Performing risky operations...")

		// Simulate a failure
		shouldRollback := true
		if shouldRollback {
			log.Println("Risky operation failed, rolling back to savepoint...")
			if err := spManager.RollbackToSavepoint(sp1); err != nil {
				return fmt.Errorf("failed to rollback to savepoint: %w", err)
			}
		}

		// Continue with safe operations
		log.Println("Performing safe operations...")

		// Release savepoint
		if err := spManager.ReleaseSavepoint(sp1); err != nil {
			log.Printf("Warning: failed to release savepoint: %v", err)
		}

		return nil
	})

	if err != nil {
		log.Printf("Savepoint transaction failed: %v", err)
	} else {
		log.Println("Savepoint transaction completed successfully")
	}
}

// ExampleTransactionMonitoring demonstrates transaction monitoring.
func ExampleTransactionMonitoring(manager *MySQLManager) {
	log.Println("=== Example: Transaction Monitoring ===")

	monitor := manager.GetTransactionMonitor()

	// Perform some transactions to generate metrics
	ctx := context.Background()
	for i := 0; i < 10; i++ {
		err := manager.ExecuteWithTransaction(ctx, func(_ *gorm.DB) error {
			// Simulate work
			time.Sleep(time.Duration(i*10) * time.Millisecond)

			// Simulate occasional failures
			if i%4 == 0 {
				return fmt.Errorf("simulated error")
			}

			return nil
		})

		if err != nil {
			log.Printf("Transaction %d failed: %v", i, err)
		}
	}

	// Get transaction statistics
	stats := monitor.GetStats()
	log.Printf("Transaction Statistics:")
	log.Printf("  Active Transactions: %d", stats.ActiveTransactions)
	log.Printf("  Committed Transactions: %d", stats.CommittedTxns)
	log.Printf("  Rolled Back Transactions: %d", stats.RolledBackTxns)
	log.Printf("  Failed Transactions: %d", stats.FailedTxns)
	log.Printf("  Average Latency: %v", stats.AverageLatency)
}

// ExampleTimeoutOperations demonstrates timeout handling.
func ExampleTimeoutOperations(manager *MySQLManager) {
	log.Println("=== Example: Timeout Operations ===")

	// Execute operation with timeout
	err := manager.ExecuteWithTimeout(2*time.Second, func(_ *gorm.DB) error {
		log.Println("Starting operation with 2-second timeout...")

		// Simulate long-running operation
		time.Sleep(3 * time.Second)

		return nil
	})

	if err != nil {
		log.Printf("Operation timed out: %v", err)
	} else {
		log.Println("Operation completed within timeout")
	}
}
