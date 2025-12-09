// Package domain provides core DDD building blocks for all Yousoon services.
package domain

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// DomainEvent represents an event that occurred in the domain.
// Events are immutable records of something that happened.
type DomainEvent interface {
	// EventID returns the unique identifier of this event instance
	EventID() string
	// EventName returns the name of the event (e.g., "UserRegistered")
	EventName() string
	// OccurredAt returns when the event occurred
	OccurredAt() time.Time
	// AggregateID returns the ID of the aggregate that produced this event
	AggregateID() string
	// AggregateType returns the type of aggregate (e.g., "User", "Offer")
	AggregateType() string
	// Version returns the event schema version for compatibility
	Version() int
	// Payload returns the event data as JSON bytes
	Payload() ([]byte, error)
}

// BaseEvent provides a base implementation for domain events.
type BaseEvent struct {
	ID            string    `json:"event_id"`
	Name          string    `json:"event_name"`
	Timestamp     time.Time `json:"occurred_at"`
	AggID         string    `json:"aggregate_id"`
	AggType       string    `json:"aggregate_type"`
	SchemaVersion int       `json:"version"`
	Data          any       `json:"data,omitempty"`
}

// NewBaseEvent creates a new base event with generated ID and current timestamp.
func NewBaseEvent(name, aggregateID, aggregateType string, data any) BaseEvent {
	return BaseEvent{
		ID:            uuid.New().String(),
		Name:          name,
		Timestamp:     time.Now().UTC(),
		AggID:         aggregateID,
		AggType:       aggregateType,
		SchemaVersion: 1,
		Data:          data,
	}
}

// NewEventID generates a new unique event ID.
func NewEventID() string {
	return uuid.New().String()
}

// EventID implements DomainEvent.
func (e BaseEvent) EventID() string {
	return e.ID
}

// EventName implements DomainEvent.
func (e BaseEvent) EventName() string {
	return e.Name
}

// OccurredAt implements DomainEvent.
func (e BaseEvent) OccurredAt() time.Time {
	return e.Timestamp
}

// AggregateID implements DomainEvent.
func (e BaseEvent) AggregateID() string {
	return e.AggID
}

// AggregateType implements DomainEvent.
func (e BaseEvent) AggregateType() string {
	return e.AggType
}

// Version implements DomainEvent.
func (e BaseEvent) Version() int {
	return e.SchemaVersion
}

// Payload implements DomainEvent.
func (e BaseEvent) Payload() ([]byte, error) {
	return json.Marshal(e.Data)
}

// EventMetadata contains metadata about an event for messaging systems.
type EventMetadata struct {
	CorrelationID string            `json:"correlation_id,omitempty"`
	CausationID   string            `json:"causation_id,omitempty"`
	UserID        string            `json:"user_id,omitempty"`
	TraceID       string            `json:"trace_id,omitempty"`
	SpanID        string            `json:"span_id,omitempty"`
	Headers       map[string]string `json:"headers,omitempty"`
}

// EventEnvelope wraps an event with metadata for transport.
type EventEnvelope struct {
	Event    DomainEvent   `json:"event"`
	Metadata EventMetadata `json:"metadata"`
}

// NewEventEnvelope creates an envelope for an event.
func NewEventEnvelope(event DomainEvent, metadata EventMetadata) EventEnvelope {
	return EventEnvelope{
		Event:    event,
		Metadata: metadata,
	}
}

// EventHandler handles domain events.
type EventHandler interface {
	// Handle processes the event
	Handle(event DomainEvent) error
	// HandledEvents returns the list of event names this handler processes
	HandledEvents() []string
}

// EventPublisher publishes domain events to a message broker.
type EventPublisher interface {
	// Publish publishes a single event
	Publish(event DomainEvent) error
	// PublishWithMetadata publishes an event with metadata
	PublishWithMetadata(event DomainEvent, metadata EventMetadata) error
	// PublishAll publishes multiple events
	PublishAll(events []DomainEvent) error
}

// EventSubscriber subscribes to domain events from a message broker.
type EventSubscriber interface {
	// Subscribe subscribes to events with the given handler
	Subscribe(eventName string, handler EventHandler) error
	// SubscribeAll subscribes to multiple event types
	SubscribeAll(eventNames []string, handler EventHandler) error
	// Unsubscribe removes a subscription
	Unsubscribe(eventName string) error
}
