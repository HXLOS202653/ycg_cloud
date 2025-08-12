// Package database provides transaction management for database operations.
package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"gorm.io/gorm"
)

// TransactionManager provides transaction management functionality.
type TransactionManager struct {
	db *gorm.DB
}

// TransactionOptions contains options for transaction execution.
type TransactionOptions struct {
	Isolation sql.IsolationLevel
	ReadOnly  bool
	Timeout   time.Duration
}

// DefaultTransactionOptions returns default transaction options.
func DefaultTransactionOptions() *TransactionOptions {
	return &TransactionOptions{
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  false,
		Timeout:   30 * time.Second,
	}
}

// NewTransactionManager creates a new transaction manager.
func NewTransactionManager(db *gorm.DB) *TransactionManager {
	return &TransactionManager{
		db: db,
	}
}

// WithTransaction executes a function within a database transaction.
// If the function returns an error, the transaction is rolled back.
// Otherwise, the transaction is committed.
func (tm *TransactionManager) WithTransaction(ctx context.Context, fn func(*gorm.DB) error) error {
	return tm.WithTransactionOptions(ctx, DefaultTransactionOptions(), fn)
}

// WithTransactionOptions executes a function within a database transaction with custom options.
func (tm *TransactionManager) WithTransactionOptions(ctx context.Context, opts *TransactionOptions, fn func(*gorm.DB) error) error {
	// Create context with timeout if specified
	if opts.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, opts.Timeout)
		defer cancel()
	}

	// Begin transaction with context
	tx := tm.db.WithContext(ctx).Begin(&sql.TxOptions{
		Isolation: opts.Isolation,
		ReadOnly:  opts.ReadOnly,
	})

	if tx.Error != nil {
		return fmt.Errorf("failed to begin transaction: %w", tx.Error)
	}

	// Defer rollback in case of panic or error
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r) // Re-panic after rollback
		}
	}()

	// Execute the function
	if err := fn(tx); err != nil {
		// Rollback on error
		if rbErr := tx.Rollback().Error; rbErr != nil {
			log.Printf("Failed to rollback transaction: %v", rbErr)
			return fmt.Errorf("transaction error: %w, rollback error: %w", err, rbErr)
		}
		return err
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// WithReadOnlyTransaction executes a read-only transaction.
func (tm *TransactionManager) WithReadOnlyTransaction(ctx context.Context, fn func(*gorm.DB) error) error {
	opts := &TransactionOptions{
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  true,
		Timeout:   30 * time.Second,
	}
	return tm.WithTransactionOptions(ctx, opts, fn)
}

// WithRepeatableReadTransaction executes a transaction with repeatable read isolation.
func (tm *TransactionManager) WithRepeatableReadTransaction(ctx context.Context, fn func(*gorm.DB) error) error {
	opts := &TransactionOptions{
		Isolation: sql.LevelRepeatableRead,
		ReadOnly:  false,
		Timeout:   30 * time.Second,
	}
	return tm.WithTransactionOptions(ctx, opts, fn)
}

// WithSerializableTransaction executes a transaction with serializable isolation.
func (tm *TransactionManager) WithSerializableTransaction(ctx context.Context, fn func(*gorm.DB) error) error {
	opts := &TransactionOptions{
		Isolation: sql.LevelSerializable,
		ReadOnly:  false,
		Timeout:   30 * time.Second,
	}
	return tm.WithTransactionOptions(ctx, opts, fn)
}

// BatchOperation represents a batch operation for bulk processing.
type BatchOperation struct {
	tm        *TransactionManager
	batchSize int
	timeout   time.Duration
}

// NewBatchOperation creates a new batch operation manager.
func (tm *TransactionManager) NewBatchOperation(batchSize int, timeout time.Duration) *BatchOperation {
	if batchSize <= 0 {
		batchSize = 1000 // Default batch size
	}
	if timeout <= 0 {
		timeout = 5 * time.Minute // Default timeout
	}

	return &BatchOperation{
		tm:        tm,
		batchSize: batchSize,
		timeout:   timeout,
	}
}

// ExecuteBatch executes operations in batches with transaction management.
func (bo *BatchOperation) ExecuteBatch(ctx context.Context, operations []func(*gorm.DB) error) error {
	if len(operations) == 0 {
		return nil
	}

	totalBatches := (len(operations) + bo.batchSize - 1) / bo.batchSize
	log.Printf("Executing %d operations in %d batches (batch size: %d)", len(operations), totalBatches, bo.batchSize)

	for i := 0; i < len(operations); i += bo.batchSize {
		end := i + bo.batchSize
		if end > len(operations) {
			end = len(operations)
		}

		batchOps := operations[i:end]
		batchNum := (i / bo.batchSize) + 1

		log.Printf("Processing batch %d/%d (%d operations)", batchNum, totalBatches, len(batchOps))

		// Execute batch within transaction
		err := bo.tm.WithTransactionOptions(ctx, &TransactionOptions{
			Isolation: sql.LevelReadCommitted,
			ReadOnly:  false,
			Timeout:   bo.timeout,
		}, func(tx *gorm.DB) error {
			for opIndex, op := range batchOps {
				if err := op(tx); err != nil {
					return fmt.Errorf("operation %d in batch %d failed: %w", opIndex+1, batchNum, err)
				}
			}
			return nil
		})

		if err != nil {
			return fmt.Errorf("batch %d failed: %w", batchNum, err)
		}

		// Check context cancellation between batches
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			// Continue with next batch
		}
	}

	log.Printf("Successfully completed all %d batches", totalBatches)
	return nil
}

// Savepoint represents a database savepoint within a transaction.
type Savepoint struct {
	name string
	tx   *gorm.DB
}

// SavepointManager manages savepoints within transactions.
type SavepointManager struct {
	tx             *gorm.DB
	savepointCount int
}

// NewSavepointManager creates a new savepoint manager for a transaction.
func NewSavepointManager(tx *gorm.DB) *SavepointManager {
	return &SavepointManager{
		tx:             tx,
		savepointCount: 0,
	}
}

// CreateSavepoint creates a new savepoint.
func (sm *SavepointManager) CreateSavepoint() (*Savepoint, error) {
	sm.savepointCount++
	name := fmt.Sprintf("sp_%d", sm.savepointCount)

	if err := sm.tx.Exec(fmt.Sprintf("SAVEPOINT %s", name)).Error; err != nil {
		return nil, fmt.Errorf("failed to create savepoint %s: %w", name, err)
	}

	return &Savepoint{
		name: name,
		tx:   sm.tx,
	}, nil
}

// RollbackToSavepoint rolls back to a specific savepoint.
func (sm *SavepointManager) RollbackToSavepoint(sp *Savepoint) error {
	if err := sm.tx.Exec(fmt.Sprintf("ROLLBACK TO SAVEPOINT %s", sp.name)).Error; err != nil {
		return fmt.Errorf("failed to rollback to savepoint %s: %w", sp.name, err)
	}
	return nil
}

// ReleaseSavepoint releases a savepoint.
func (sm *SavepointManager) ReleaseSavepoint(sp *Savepoint) error {
	if err := sm.tx.Exec(fmt.Sprintf("RELEASE SAVEPOINT %s", sp.name)).Error; err != nil {
		return fmt.Errorf("failed to release savepoint %s: %w", sp.name, err)
	}
	return nil
}

// TransactionStats provides transaction statistics.
type TransactionStats struct {
	ActiveTransactions int64
	CommittedTxns      int64
	RolledBackTxns     int64
	FailedTxns         int64
	AverageLatency     time.Duration
}

// TransactionMonitor monitors transaction performance and statistics.
type TransactionMonitor struct {
	stats TransactionStats
}

// NewTransactionMonitor creates a new transaction monitor.
func NewTransactionMonitor() *TransactionMonitor {
	return &TransactionMonitor{}
}

// GetStats returns current transaction statistics.
func (tm *TransactionMonitor) GetStats() TransactionStats {
	return tm.stats
}

// RecordCommit records a successful transaction commit.
func (tm *TransactionMonitor) RecordCommit(duration time.Duration) {
	tm.stats.CommittedTxns++
	// Update average latency (simplified)
	tm.stats.AverageLatency = duration
}

// RecordRollback records a transaction rollback.
func (tm *TransactionMonitor) RecordRollback() {
	tm.stats.RolledBackTxns++
}

// RecordFailure records a transaction failure.
func (tm *TransactionMonitor) RecordFailure() {
	tm.stats.FailedTxns++
}

// WithTransactionMonitoring wraps transaction execution with monitoring.
func (tm *TransactionManager) WithTransactionMonitoring(ctx context.Context, monitor *TransactionMonitor, fn func(*gorm.DB) error) error {
	start := time.Now()
	monitor.stats.ActiveTransactions++

	defer func() {
		monitor.stats.ActiveTransactions--

		if r := recover(); r != nil {
			monitor.RecordFailure()
			panic(r)
		}
	}()

	err := tm.WithTransaction(ctx, fn)
	duration := time.Since(start)

	if err != nil {
		monitor.RecordRollback()
		return err
	}

	monitor.RecordCommit(duration)
	return nil
}
