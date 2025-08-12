// Package database provides tests for transaction management functionality.
package database

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestTransactionManager_WithTransaction(t *testing.T) {
	db := createTestDB(t)
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	tm := NewTransactionManager(db)
	ctx := context.Background()

	// Test successful transaction
	executed := false
	err := tm.WithTransaction(ctx, func(tx *gorm.DB) error {
		executed = true
		// Verify we're in a transaction
		assert.NotEqual(t, db, tx)
		return nil
	})

	assert.NoError(t, err)
	assert.True(t, executed)
}

func TestTransactionManager_WithTransaction_Rollback(t *testing.T) {
	db := createTestDB(t)
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	tm := NewTransactionManager(db)
	ctx := context.Background()

	// Test transaction rollback on error
	testError := errors.New("test error")
	err := tm.WithTransaction(ctx, func(_ *gorm.DB) error {
		return testError
	})

	assert.Error(t, err)
	assert.Equal(t, testError, err)
}

func TestTransactionManager_WithTransactionOptions(t *testing.T) {
	db := createTestDB(t)
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	tm := NewTransactionManager(db)
	ctx := context.Background()

	opts := &TransactionOptions{
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  false,
		Timeout:   5 * time.Second,
	}

	executed := false
	err := tm.WithTransactionOptions(ctx, opts, func(_ *gorm.DB) error {
		executed = true
		return nil
	})

	assert.NoError(t, err)
	assert.True(t, executed)
}

func TestTransactionManager_WithReadOnlyTransaction(t *testing.T) {
	db := createTestDB(t)
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	tm := NewTransactionManager(db)
	ctx := context.Background()

	executed := false
	err := tm.WithReadOnlyTransaction(ctx, func(_ *gorm.DB) error {
		executed = true
		// In a real test, you might verify read-only behavior
		return nil
	})

	assert.NoError(t, err)
	assert.True(t, executed)
}

func TestTransactionManager_WithRepeatableReadTransaction(t *testing.T) {
	db := createTestDB(t)
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	tm := NewTransactionManager(db)
	ctx := context.Background()

	executed := false
	err := tm.WithRepeatableReadTransaction(ctx, func(tx *gorm.DB) error {
		executed = true
		return nil
	})

	assert.NoError(t, err)
	assert.True(t, executed)
}

func TestTransactionManager_WithSerializableTransaction(t *testing.T) {
	db := createTestDB(t)
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	tm := NewTransactionManager(db)
	ctx := context.Background()

	executed := false
	err := tm.WithSerializableTransaction(ctx, func(tx *gorm.DB) error {
		executed = true
		return nil
	})

	assert.NoError(t, err)
	assert.True(t, executed)
}

func TestTransactionManager_Timeout(t *testing.T) {
	db := createTestDB(t)
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	tm := NewTransactionManager(db)
	ctx := context.Background()

	opts := &TransactionOptions{
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  false,
		Timeout:   100 * time.Millisecond,
	}

	err := tm.WithTransactionOptions(ctx, opts, func(_ *gorm.DB) error {
		// Sleep longer than timeout
		time.Sleep(200 * time.Millisecond)
		return nil
	})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "context deadline exceeded")
}

func TestBatchOperation_ExecuteBatch(t *testing.T) {
	db := createTestDB(t)
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	tm := NewTransactionManager(db)
	batchOp := tm.NewBatchOperation(2, 5*time.Second)

	ctx := context.Background()
	executed := 0

	operations := []func(*gorm.DB) error{
		func(tx *gorm.DB) error {
			executed++
			return nil
		},
		func(tx *gorm.DB) error {
			executed++
			return nil
		},
		func(tx *gorm.DB) error {
			executed++
			return nil
		},
		func(tx *gorm.DB) error {
			executed++
			return nil
		},
	}

	err := batchOp.ExecuteBatch(ctx, operations)
	assert.NoError(t, err)
	assert.Equal(t, 4, executed)
}

func TestBatchOperation_ExecuteBatch_WithError(t *testing.T) {
	db := createTestDB(t)
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	tm := NewTransactionManager(db)
	batchOp := tm.NewBatchOperation(2, 5*time.Second)

	ctx := context.Background()
	executed := 0

	operations := []func(*gorm.DB) error{
		func(tx *gorm.DB) error {
			executed++
			return nil
		},
		func(tx *gorm.DB) error {
			executed++
			return errors.New("batch error")
		},
		func(tx *gorm.DB) error {
			executed++
			return nil
		},
	}

	err := batchOp.ExecuteBatch(ctx, operations)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "batch error")
	// First batch should succeed, second should fail
	assert.Equal(t, 2, executed)
}

func TestBatchOperation_EmptyOperations(t *testing.T) {
	db := createTestDB(t)
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	tm := NewTransactionManager(db)
	batchOp := tm.NewBatchOperation(10, 5*time.Second)

	ctx := context.Background()
	operations := []func(*gorm.DB) error{}

	err := batchOp.ExecuteBatch(ctx, operations)
	assert.NoError(t, err)
}

func TestSavepointManager(t *testing.T) {
	db := createTestDB(t)
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	tm := NewTransactionManager(db)
	ctx := context.Background()

	err := tm.WithTransaction(ctx, func(tx *gorm.DB) error {
		spManager := NewSavepointManager(tx)

		// Create savepoint
		sp1, err := spManager.CreateSavepoint()
		assert.NoError(t, err)
		assert.NotNil(t, sp1)
		assert.Equal(t, "sp_1", sp1.name)

		// Create another savepoint
		sp2, err := spManager.CreateSavepoint()
		assert.NoError(t, err)
		assert.NotNil(t, sp2)
		assert.Equal(t, "sp_2", sp2.name)

		// In a real test with actual database, you would test:
		// - Rollback to savepoint
		// - Release savepoint
		// For this mock test, we just verify the savepoints were created

		return nil
	})

	assert.NoError(t, err)
}

func TestTransactionMonitor(t *testing.T) {
	monitor := NewTransactionMonitor()

	// Initial stats
	stats := monitor.GetStats()
	assert.Equal(t, int64(0), stats.CommittedTxns)
	assert.Equal(t, int64(0), stats.RolledBackTxns)
	assert.Equal(t, int64(0), stats.FailedTxns)

	// Record some transactions
	monitor.RecordCommit(100 * time.Millisecond)
	monitor.RecordCommit(200 * time.Millisecond)
	monitor.RecordRollback()
	monitor.RecordFailure()

	stats = monitor.GetStats()
	assert.Equal(t, int64(2), stats.CommittedTxns)
	assert.Equal(t, int64(1), stats.RolledBackTxns)
	assert.Equal(t, int64(1), stats.FailedTxns)
	assert.Equal(t, 200*time.Millisecond, stats.AverageLatency) // Last recorded latency
}

func TestTransactionManager_WithTransactionMonitoring(t *testing.T) {
	db := createTestDB(t)
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	tm := NewTransactionManager(db)
	monitor := NewTransactionMonitor()
	ctx := context.Background()

	// Successful transaction
	err := tm.WithTransactionMonitoring(ctx, monitor, func(tx *gorm.DB) error {
		time.Sleep(10 * time.Millisecond)
		return nil
	})

	assert.NoError(t, err)
	stats := monitor.GetStats()
	assert.Equal(t, int64(1), stats.CommittedTxns)
	assert.True(t, stats.AverageLatency > 0)

	// Failed transaction
	err = tm.WithTransactionMonitoring(ctx, monitor, func(tx *gorm.DB) error {
		return errors.New("test error")
	})

	assert.Error(t, err)
	stats = monitor.GetStats()
	assert.Equal(t, int64(1), stats.CommittedTxns)
	assert.Equal(t, int64(1), stats.RolledBackTxns)
}

func TestDefaultTransactionOptions(t *testing.T) {
	opts := DefaultTransactionOptions()

	assert.Equal(t, sql.LevelReadCommitted, opts.Isolation)
	assert.False(t, opts.ReadOnly)
	assert.Equal(t, 30*time.Second, opts.Timeout)
}

// Benchmark tests
func BenchmarkTransactionManager_WithTransaction(b *testing.B) {
	db := createTestDB(b)
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	tm := NewTransactionManager(db)
	ctx := context.Background()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = tm.WithTransaction(ctx, func(tx *gorm.DB) error {
				// Simulate lightweight operation
				return nil
			})
		}
	})
}

func BenchmarkBatchOperation_ExecuteBatch(b *testing.B) {
	db := createTestDB(b)
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	tm := NewTransactionManager(db)
	batchOp := tm.NewBatchOperation(100, time.Minute)
	ctx := context.Background()

	// Create operations for benchmark
	operations := make([]func(*gorm.DB) error, 1000)
	for i := range operations {
		operations[i] = func(tx *gorm.DB) error {
			return nil
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = batchOp.ExecuteBatch(ctx, operations)
	}
}
