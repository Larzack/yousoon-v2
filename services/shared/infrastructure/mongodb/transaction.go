package mongodb

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

// TransactionManager handles MongoDB transactions.
type TransactionManager struct {
	client *Client
}

// NewTransactionManager creates a new transaction manager.
func NewTransactionManager(client *Client) *TransactionManager {
	return &TransactionManager{
		client: client,
	}
}

// TransactionFunc is a function that executes within a transaction.
type TransactionFunc func(sessCtx mongo.SessionContext) error

// WithTransaction executes a function within a MongoDB transaction.
// It automatically handles commit and rollback.
func (tm *TransactionManager) WithTransaction(ctx context.Context, fn TransactionFunc) error {
	// Create session
	session, err := tm.client.Client().StartSession()
	if err != nil {
		return fmt.Errorf("failed to start session: %w", err)
	}
	defer session.EndSession(ctx)

	// Configure transaction options
	txnOpts := options.Transaction().
		SetWriteConcern(writeconcern.Majority()).
		SetReadPreference(nil) // Use default read preference

	// Execute transaction
	_, err = session.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
		if err := fn(sessCtx); err != nil {
			return nil, err
		}
		return nil, nil
	}, txnOpts)

	if err != nil {
		return fmt.Errorf("transaction failed: %w", err)
	}

	return nil
}

// WithTransactionResult executes a function within a transaction and returns a result.
func (tm *TransactionManager) WithTransactionResult(ctx context.Context, fn func(sessCtx mongo.SessionContext) (interface{}, error)) (interface{}, error) {
	session, err := tm.client.Client().StartSession()
	if err != nil {
		return nil, fmt.Errorf("failed to start session: %w", err)
	}
	defer session.EndSession(ctx)

	txnOpts := options.Transaction().
		SetWriteConcern(writeconcern.Majority())

	result, err := session.WithTransaction(ctx, fn, txnOpts)
	if err != nil {
		return nil, fmt.Errorf("transaction failed: %w", err)
	}

	return result, nil
}

// UnitOfWork represents a unit of work pattern for managing transactions.
type UnitOfWork struct {
	tm         *TransactionManager
	operations []func(mongo.SessionContext) error
}

// NewUnitOfWork creates a new unit of work.
func NewUnitOfWork(tm *TransactionManager) *UnitOfWork {
	return &UnitOfWork{
		tm:         tm,
		operations: make([]func(mongo.SessionContext) error, 0),
	}
}

// Add adds an operation to the unit of work.
func (uow *UnitOfWork) Add(op func(mongo.SessionContext) error) {
	uow.operations = append(uow.operations, op)
}

// Commit executes all operations within a single transaction.
func (uow *UnitOfWork) Commit(ctx context.Context) error {
	return uow.tm.WithTransaction(ctx, func(sessCtx mongo.SessionContext) error {
		for i, op := range uow.operations {
			if err := op(sessCtx); err != nil {
				return fmt.Errorf("operation %d failed: %w", i, err)
			}
		}
		return nil
	})
}

// Clear removes all pending operations.
func (uow *UnitOfWork) Clear() {
	uow.operations = make([]func(mongo.SessionContext) error, 0)
}

// Count returns the number of pending operations.
func (uow *UnitOfWork) Count() int {
	return len(uow.operations)
}

// SessionRepository wraps a repository to work within a session context.
type SessionRepository[T any] struct {
	repo    *Repository[T]
	session mongo.Session
}

// NewSessionRepository creates a repository bound to a specific session.
func NewSessionRepository[T any](repo *Repository[T], session mongo.Session) *SessionRepository[T] {
	return &SessionRepository[T]{
		repo:    repo,
		session: session,
	}
}

// WithinSession executes the provided function within the session context.
func (sr *SessionRepository[T]) WithinSession(ctx context.Context, fn func(mongo.SessionContext, *Repository[T]) error) error {
	return mongo.WithSession(ctx, sr.session, func(sessCtx mongo.SessionContext) error {
		return fn(sessCtx, sr.repo)
	})
}
