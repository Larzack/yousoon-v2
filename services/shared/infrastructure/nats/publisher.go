package nats

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/nats-io/nats.go"

	"github.com/yousoon/services/shared/domain"
)

// EventPublisher publishes domain events to NATS JetStream.
type EventPublisher struct {
	client *Client
}

// NewEventPublisher creates a new event publisher.
func NewEventPublisher(client *Client) *EventPublisher {
	return &EventPublisher{
		client: client,
	}
}

// Publish publishes a single domain event.
func (p *EventPublisher) Publish(ctx context.Context, event domain.DomainEvent) error {
	subject := p.buildSubject(event)
	
	envelope := EventEnvelope{
		EventID:     event.EventID(),
		EventType:   event.EventName(),
		AggregateID: event.AggregateID(),
		OccurredAt:  event.OccurredAt(),
		Payload:     event,
	}

	data, err := json.Marshal(envelope)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	// Publish with acknowledgment
	ack, err := p.client.JetStream().Publish(subject, data, nats.Context(ctx))
	if err != nil {
		return fmt.Errorf("failed to publish event: %w", err)
	}

	// Log acknowledgment (optional)
	_ = ack

	return nil
}

// PublishAll publishes multiple domain events.
func (p *EventPublisher) PublishAll(ctx context.Context, events []domain.DomainEvent) error {
	for _, event := range events {
		if err := p.Publish(ctx, event); err != nil {
			return err
		}
	}
	return nil
}

// PublishAsync publishes an event asynchronously.
func (p *EventPublisher) PublishAsync(event domain.DomainEvent) (nats.PubAckFuture, error) {
	subject := p.buildSubject(event)
	
	envelope := EventEnvelope{
		EventID:     event.EventID(),
		EventType:   event.EventName(),
		AggregateID: event.AggregateID(),
		OccurredAt:  event.OccurredAt(),
		Payload:     event,
	}

	data, err := json.Marshal(envelope)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal event: %w", err)
	}

	return p.client.JetStream().PublishAsync(subject, data)
}

// buildSubject builds the NATS subject from the event type.
func (p *EventPublisher) buildSubject(event domain.DomainEvent) string {
	return fmt.Sprintf("yousoon.events.%s", event.EventName())
}

// EventEnvelope wraps a domain event for transport.
type EventEnvelope struct {
	EventID     string              `json:"event_id"`
	EventType   string              `json:"event_type"`
	AggregateID string              `json:"aggregate_id"`
	OccurredAt  time.Time           `json:"occurred_at"`
	Payload     domain.DomainEvent  `json:"payload"`
}

// PublishNotification publishes a notification message.
func (p *EventPublisher) PublishNotification(ctx context.Context, notif NotificationMessage) error {
	subject := fmt.Sprintf("yousoon.notifications.%s.%s", notif.Channel, notif.Type)

	data, err := json.Marshal(notif)
	if err != nil {
		return fmt.Errorf("failed to marshal notification: %w", err)
	}

	_, err = p.client.JetStream().Publish(subject, data, nats.Context(ctx))
	if err != nil {
		return fmt.Errorf("failed to publish notification: %w", err)
	}

	return nil
}

// NotificationMessage represents a notification to be sent.
type NotificationMessage struct {
	ID          string                 `json:"id"`
	Type        string                 `json:"type"`        // booking_confirmed, offer_nearby, etc.
	Channel     string                 `json:"channel"`     // push, email, sms
	RecipientID string                 `json:"recipient_id"`
	Title       string                 `json:"title"`
	Body        string                 `json:"body"`
	Data        map[string]interface{} `json:"data"`
	ScheduledAt *time.Time             `json:"scheduled_at,omitempty"`
	CreatedAt   time.Time              `json:"created_at"`
}

// NewNotificationMessage creates a new notification message.
func NewNotificationMessage(notifType, channel, recipientID, title, body string) NotificationMessage {
	return NotificationMessage{
		ID:          domain.NewEventID(),
		Type:        notifType,
		Channel:     channel,
		RecipientID: recipientID,
		Title:       title,
		Body:        body,
		Data:        make(map[string]interface{}),
		CreatedAt:   time.Now(),
	}
}

// AggregatePublisher publishes all uncommitted events from an aggregate.
type AggregatePublisher struct {
	publisher *EventPublisher
}

// NewAggregatePublisher creates a new aggregate publisher.
func NewAggregatePublisher(publisher *EventPublisher) *AggregatePublisher {
	return &AggregatePublisher{
		publisher: publisher,
	}
}

// PublishEvents publishes all uncommitted events from an aggregate and clears them.
func (p *AggregatePublisher) PublishEvents(ctx context.Context, aggregate *domain.AggregateRoot) error {
	events := aggregate.GetDomainEvents()
	if len(events) == 0 {
		return nil
	}

	if err := p.publisher.PublishAll(ctx, events); err != nil {
		return err
	}

	aggregate.ClearDomainEvents()
	return nil
}
