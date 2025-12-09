package nats

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/nats-io/nats.go"
)

// MessageHandler is a function that handles incoming messages.
type MessageHandler func(ctx context.Context, msg *nats.Msg) error

// EventHandler is a function that handles incoming events.
type EventHandler func(ctx context.Context, envelope EventEnvelope) error

// Subscriber handles message subscriptions from NATS JetStream.
type Subscriber struct {
	client        *Client
	subscriptions []*nats.Subscription
	mu            sync.Mutex
	wg            sync.WaitGroup
	done          chan struct{}
}

// NewSubscriber creates a new subscriber.
func NewSubscriber(client *Client) *Subscriber {
	return &Subscriber{
		client:        client,
		subscriptions: make([]*nats.Subscription, 0),
		done:          make(chan struct{}),
	}
}

// SubscribeConfig holds subscription configuration.
type SubscribeConfig struct {
	Stream        string
	Consumer      string
	Subject       string
	Handler       MessageHandler
	MaxInFlight   int
	AckWait       time.Duration
	MaxDeliver    int
	DeliverPolicy nats.DeliverPolicy
}

// DefaultSubscribeConfig returns a default subscription configuration.
func DefaultSubscribeConfig(stream, consumer, subject string, handler MessageHandler) SubscribeConfig {
	return SubscribeConfig{
		Stream:        stream,
		Consumer:      consumer,
		Subject:       subject,
		Handler:       handler,
		MaxInFlight:   100,
		AckWait:       30 * time.Second,
		MaxDeliver:    5,
		DeliverPolicy: nats.DeliverAllPolicy,
	}
}

// Subscribe creates a push-based subscription.
func (s *Subscriber) Subscribe(ctx context.Context, cfg SubscribeConfig) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Create or ensure consumer exists
	consumerCfg := ConsumerConfig{
		Durable:       cfg.Consumer,
		DeliverPolicy: cfg.DeliverPolicy,
		AckWait:       cfg.AckWait,
		MaxDeliver:    cfg.MaxDeliver,
		FilterSubject: cfg.Subject,
		MaxAckPending: cfg.MaxInFlight,
	}

	if _, err := s.client.CreateOrUpdateConsumer(cfg.Stream, consumerCfg); err != nil {
		return fmt.Errorf("failed to create consumer: %w", err)
	}

	// Create subscription
	sub, err := s.client.JetStream().PullSubscribe(
		cfg.Subject,
		cfg.Consumer,
		nats.Bind(cfg.Stream, cfg.Consumer),
	)
	if err != nil {
		return fmt.Errorf("failed to create subscription: %w", err)
	}

	s.subscriptions = append(s.subscriptions, sub)

	// Start message processing goroutine
	s.wg.Add(1)
	go s.processMessages(ctx, sub, cfg.Handler, cfg.MaxInFlight)

	return nil
}

// processMessages processes messages from a pull subscription.
func (s *Subscriber) processMessages(ctx context.Context, sub *nats.Subscription, handler MessageHandler, batchSize int) {
	defer s.wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case <-s.done:
			return
		default:
			msgs, err := sub.Fetch(batchSize, nats.MaxWait(5*time.Second))
			if err != nil {
				if err == nats.ErrTimeout {
					continue
				}
				log.Printf("Error fetching messages: %v", err)
				continue
			}

			for _, msg := range msgs {
				if err := s.handleMessage(ctx, msg, handler); err != nil {
					log.Printf("Error handling message: %v", err)
					// NAK the message for redelivery
					msg.Nak()
				} else {
					// ACK successful processing
					msg.Ack()
				}
			}
		}
	}
}

// handleMessage handles a single message.
func (s *Subscriber) handleMessage(ctx context.Context, msg *nats.Msg, handler MessageHandler) error {
	// Create a timeout context for the handler
	handlerCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	return handler(handlerCtx, msg)
}

// SubscribeEvents subscribes to events with automatic deserialization.
func (s *Subscriber) SubscribeEvents(ctx context.Context, stream, consumer, subject string, handler EventHandler) error {
	messageHandler := func(ctx context.Context, msg *nats.Msg) error {
		var envelope EventEnvelope
		if err := json.Unmarshal(msg.Data, &envelope); err != nil {
			return fmt.Errorf("failed to unmarshal event: %w", err)
		}
		return handler(ctx, envelope)
	}

	cfg := DefaultSubscribeConfig(stream, consumer, subject, messageHandler)
	return s.Subscribe(ctx, cfg)
}

// Close stops all subscriptions and waits for processing to complete.
func (s *Subscriber) Close() error {
	close(s.done)

	s.mu.Lock()
	defer s.mu.Unlock()

	var errs []error
	for _, sub := range s.subscriptions {
		if err := sub.Unsubscribe(); err != nil {
			errs = append(errs, err)
		}
	}

	// Wait for all goroutines to finish
	s.wg.Wait()

	if len(errs) > 0 {
		return fmt.Errorf("errors closing subscriptions: %v", errs)
	}
	return nil
}

// EventRouter routes events to appropriate handlers based on event type.
type EventRouter struct {
	subscriber *Subscriber
	handlers   map[string]EventHandler
	mu         sync.RWMutex
}

// NewEventRouter creates a new event router.
func NewEventRouter(subscriber *Subscriber) *EventRouter {
	return &EventRouter{
		subscriber: subscriber,
		handlers:   make(map[string]EventHandler),
	}
}

// RegisterHandler registers a handler for a specific event type.
func (r *EventRouter) RegisterHandler(eventType string, handler EventHandler) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.handlers[eventType] = handler
}

// Start starts the event router.
func (r *EventRouter) Start(ctx context.Context, stream, consumer, subject string) error {
	return r.subscriber.SubscribeEvents(ctx, stream, consumer, subject, r.handleEvent)
}

// handleEvent routes an event to the appropriate handler.
func (r *EventRouter) handleEvent(ctx context.Context, envelope EventEnvelope) error {
	r.mu.RLock()
	handler, exists := r.handlers[envelope.EventType]
	r.mu.RUnlock()

	if !exists {
		// Log unhandled event type but don't fail
		log.Printf("No handler for event type: %s", envelope.EventType)
		return nil
	}

	return handler(ctx, envelope)
}

// WorkerPool manages a pool of workers for processing messages.
type WorkerPool struct {
	subscriber *Subscriber
	numWorkers int
	wg         sync.WaitGroup
}

// NewWorkerPool creates a new worker pool.
func NewWorkerPool(subscriber *Subscriber, numWorkers int) *WorkerPool {
	return &WorkerPool{
		subscriber: subscriber,
		numWorkers: numWorkers,
	}
}

// Start starts the worker pool.
func (wp *WorkerPool) Start(ctx context.Context, stream, consumer, subject string, handler EventHandler) error {
	for i := 0; i < wp.numWorkers; i++ {
		wp.wg.Add(1)
		go func(workerID int) {
			defer wp.wg.Done()
			workerConsumer := fmt.Sprintf("%s-worker-%d", consumer, workerID)

			if err := wp.subscriber.SubscribeEvents(ctx, stream, workerConsumer, subject, handler); err != nil {
				log.Printf("Worker %d failed to subscribe: %v", workerID, err)
			}
		}(i)
	}
	return nil
}

// Stop stops the worker pool.
func (wp *WorkerPool) Stop() {
	wp.subscriber.Close()
	wp.wg.Wait()
}
