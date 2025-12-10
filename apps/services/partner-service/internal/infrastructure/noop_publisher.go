// Package infrastructure provides infrastructure implementations.
package infrastructure

import (
	sharedomain "github.com/yousoon/shared/domain"
)

// NoOpEventPublisher is a no-operation event publisher for development.
type NoOpEventPublisher struct{}

// NewNoOpEventPublisher creates a new NoOpEventPublisher.
func NewNoOpEventPublisher() *NoOpEventPublisher {
	return &NoOpEventPublisher{}
}

// Publish does nothing (no-op).
func (p *NoOpEventPublisher) Publish(event sharedomain.DomainEvent) error {
	return nil
}

// PublishWithMetadata does nothing (no-op).
func (p *NoOpEventPublisher) PublishWithMetadata(event sharedomain.DomainEvent, metadata sharedomain.EventMetadata) error {
	return nil
}

// PublishAll does nothing (no-op).
func (p *NoOpEventPublisher) PublishAll(events []sharedomain.DomainEvent) error {
	return nil
}
