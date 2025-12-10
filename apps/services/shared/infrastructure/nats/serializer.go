package nats

import (
	"encoding/json"
	"fmt"
	"reflect"
	"sync"
	"time"

	"github.com/yousoon/services/shared/domain"
)

// Serializer handles event serialization and deserialization.
type Serializer struct {
	registry map[string]reflect.Type
	mu       sync.RWMutex
}

// NewSerializer creates a new event serializer.
func NewSerializer() *Serializer {
	return &Serializer{
		registry: make(map[string]reflect.Type),
	}
}

// Register registers an event type for deserialization.
func (s *Serializer) Register(eventType string, example domain.DomainEvent) {
	s.mu.Lock()
	defer s.mu.Unlock()

	t := reflect.TypeOf(example)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	s.registry[eventType] = t
}

// Serialize serializes a domain event to JSON bytes.
func (s *Serializer) Serialize(event domain.DomainEvent) ([]byte, error) {
	envelope := SerializedEvent{
		EventID:     event.EventID(),
		EventType:   event.EventName(),
		AggregateID: event.AggregateID(),
		OccurredAt:  event.OccurredAt(),
		Version:     1,
	}

	payload, err := json.Marshal(event)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize payload: %w", err)
	}
	envelope.Payload = payload

	return json.Marshal(envelope)
}

// Deserialize deserializes JSON bytes to a domain event.
func (s *Serializer) Deserialize(data []byte) (domain.DomainEvent, error) {
	var envelope SerializedEvent
	if err := json.Unmarshal(data, &envelope); err != nil {
		return nil, fmt.Errorf("failed to deserialize envelope: %w", err)
	}

	s.mu.RLock()
	eventType, exists := s.registry[envelope.EventType]
	s.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("unknown event type: %s", envelope.EventType)
	}

	// Create new instance of the event type
	eventPtr := reflect.New(eventType).Interface()

	if err := json.Unmarshal(envelope.Payload, eventPtr); err != nil {
		return nil, fmt.Errorf("failed to deserialize payload: %w", err)
	}

	event, ok := eventPtr.(domain.DomainEvent)
	if !ok {
		return nil, fmt.Errorf("event does not implement DomainEvent interface")
	}

	return event, nil
}

// SerializedEvent is the wire format for events.
type SerializedEvent struct {
	EventID     string          `json:"event_id"`
	EventType   string          `json:"event_type"`
	AggregateID string          `json:"aggregate_id"`
	OccurredAt  time.Time       `json:"occurred_at"`
	Version     int             `json:"version"`
	Payload     json.RawMessage `json:"payload"`
}

// EventRegistry provides a global registry for event types.
type EventRegistry struct {
	serializer *Serializer
}

// NewEventRegistry creates a new event registry.
func NewEventRegistry() *EventRegistry {
	return &EventRegistry{
		serializer: NewSerializer(),
	}
}

// Register registers an event type.
func (r *EventRegistry) Register(event domain.DomainEvent) {
	r.serializer.Register(event.EventName(), event)
}

// Serializer returns the serializer.
func (r *EventRegistry) Serializer() *Serializer {
	return r.serializer
}

// Identity Context Events Registration
func (r *EventRegistry) RegisterIdentityEvents() {
	// These will be defined in the identity service
	// r.Register(&UserRegistered{})
	// r.Register(&UserIdentityVerified{})
	// r.Register(&UserSubscribed{})
}

// Partner Context Events Registration
func (r *EventRegistry) RegisterPartnerEvents() {
	// These will be defined in the partner service
	// r.Register(&PartnerRegistered{})
	// r.Register(&EstablishmentAdded{})
}

// Discovery Context Events Registration
func (r *EventRegistry) RegisterDiscoveryEvents() {
	// These will be defined in the discovery service
	// r.Register(&OfferCreated{})
	// r.Register(&OfferPublished{})
}

// Booking Context Events Registration
func (r *EventRegistry) RegisterBookingEvents() {
	// These will be defined in the booking service
	// r.Register(&OutingBooked{})
	// r.Register(&OutingCheckedIn{})
}

// Engagement Context Events Registration
func (r *EventRegistry) RegisterEngagementEvents() {
	// These will be defined in the engagement service
	// r.Register(&ReviewSubmitted{})
	// r.Register(&OfferFavorited{})
}

// RegisterAllEvents registers all known event types.
func (r *EventRegistry) RegisterAllEvents() {
	r.RegisterIdentityEvents()
	r.RegisterPartnerEvents()
	r.RegisterDiscoveryEvents()
	r.RegisterBookingEvents()
	r.RegisterEngagementEvents()
}

// Global registry instance
var globalRegistry = NewEventRegistry()

// GlobalRegistry returns the global event registry.
func GlobalRegistry() *EventRegistry {
	return globalRegistry
}

// RegisterEvent registers an event with the global registry.
func RegisterEvent(event domain.DomainEvent) {
	globalRegistry.Register(event)
}
