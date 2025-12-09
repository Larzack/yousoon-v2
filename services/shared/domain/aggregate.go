package domain

import (
	"sync"
	"time"
)

// AggregateRoot is the base for all aggregate roots in the domain.
// An aggregate root is the entry point to an aggregate and ensures consistency.
type AggregateRoot struct {
	BaseEntity
	events []DomainEvent
	mu     sync.RWMutex
}

// NewAggregateRoot creates a new aggregate root with timestamps.
func NewAggregateRoot() AggregateRoot {
	return AggregateRoot{
		BaseEntity: NewBaseEntity(),
		events:     make([]DomainEvent, 0),
	}
}

// AddDomainEvent adds a domain event to be dispatched after persistence.
func (a *AggregateRoot) AddDomainEvent(event DomainEvent) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.events = append(a.events, event)
}

// GetDomainEvents returns all pending domain events.
func (a *AggregateRoot) GetDomainEvents() []DomainEvent {
	a.mu.RLock()
	defer a.mu.RUnlock()
	// Return a copy to prevent external modification
	events := make([]DomainEvent, len(a.events))
	copy(events, a.events)
	return events
}

// ClearDomainEvents removes all pending domain events.
// Call this after events have been successfully published.
func (a *AggregateRoot) ClearDomainEvents() {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.events = make([]DomainEvent, 0)
}

// HasPendingEvents returns true if there are unpublished events.
func (a *AggregateRoot) HasPendingEvents() bool {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return len(a.events) > 0
}

// PendingEventCount returns the number of pending events.
func (a *AggregateRoot) PendingEventCount() int {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return len(a.events)
}

// =============================================================================
// Versioned Aggregate Root (with optimistic locking)
// =============================================================================

// VersionedAggregateRoot adds version tracking for optimistic concurrency control.
type VersionedAggregateRoot struct {
	AggregateRoot
	Version int64 `json:"version" bson:"version"`
}

// NewVersionedAggregateRoot creates a new versioned aggregate root.
func NewVersionedAggregateRoot() VersionedAggregateRoot {
	return VersionedAggregateRoot{
		AggregateRoot: NewAggregateRoot(),
		Version:       1,
	}
}

// IncrementVersion increments the version for optimistic locking.
// Call this before saving to detect concurrent modifications.
func (a *VersionedAggregateRoot) IncrementVersion() {
	a.Version++
	a.MarkUpdated()
}

// GetVersion returns the current version.
func (a *VersionedAggregateRoot) GetVersion() int64 {
	return a.Version
}

// ExpectedVersion returns the version that should be in the database.
// Used for optimistic locking checks.
func (a *VersionedAggregateRoot) ExpectedVersion() int64 {
	return a.Version - 1
}

// =============================================================================
// Aggregate Repository Interface
// =============================================================================

// Repository is the base interface for all aggregate repositories.
type Repository[T any, ID any] interface {
	// FindByID retrieves an aggregate by its ID
	FindByID(id ID) (*T, error)
	// Save persists an aggregate (insert or update)
	Save(aggregate *T) error
	// Delete removes an aggregate
	Delete(id ID) error
}

// =============================================================================
// Aggregate Helpers
// =============================================================================

// AggregateMetadata contains metadata about an aggregate.
type AggregateMetadata struct {
	ID        string    `json:"id"`
	Type      string    `json:"type"`
	Version   int64     `json:"version"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// GetMetadata extracts metadata from a versioned aggregate.
func GetMetadata(id, aggregateType string, agg *VersionedAggregateRoot) AggregateMetadata {
	return AggregateMetadata{
		ID:        id,
		Type:      aggregateType,
		Version:   agg.Version,
		CreatedAt: agg.CreatedAt,
		UpdatedAt: agg.UpdatedAt,
	}
}
