// Package nats provides NATS JetStream infrastructure components.
package nats

import (
	"context"
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
)

// Config holds NATS connection configuration.
type Config struct {
	// URL is the NATS server URL.
	URL string
	// Name is the client name.
	Name string
	// MaxReconnects is the maximum number of reconnection attempts.
	MaxReconnects int
	// ReconnectWait is the wait time between reconnection attempts.
	ReconnectWait time.Duration
	// Timeout is the connection timeout.
	Timeout time.Duration
	// Token is the authentication token.
	Token string
	// User is the username for authentication.
	User string
	// Password is the password for authentication.
	Password string
}

// DefaultConfig returns a default NATS configuration.
func DefaultConfig() Config {
	return Config{
		URL:           "nats://localhost:4222",
		Name:          "yousoon-service",
		MaxReconnects: -1, // Unlimited
		ReconnectWait: 2 * time.Second,
		Timeout:       10 * time.Second,
	}
}

// Client wraps the NATS connection with JetStream support.
type Client struct {
	conn   *nats.Conn
	js     nats.JetStreamContext
	config Config
}

// NewClient creates a new NATS client with JetStream.
func NewClient(ctx context.Context, cfg Config) (*Client, error) {
	opts := []nats.Option{
		nats.Name(cfg.Name),
		nats.MaxReconnects(cfg.MaxReconnects),
		nats.ReconnectWait(cfg.ReconnectWait),
		nats.Timeout(cfg.Timeout),
		nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
			if err != nil {
				fmt.Printf("NATS disconnected: %v\n", err)
			}
		}),
		nats.ReconnectHandler(func(nc *nats.Conn) {
			fmt.Printf("NATS reconnected to %s\n", nc.ConnectedUrl())
		}),
		nats.ClosedHandler(func(nc *nats.Conn) {
			fmt.Println("NATS connection closed")
		}),
	}

	if cfg.Token != "" {
		opts = append(opts, nats.Token(cfg.Token))
	}

	if cfg.User != "" && cfg.Password != "" {
		opts = append(opts, nats.UserInfo(cfg.User, cfg.Password))
	}

	conn, err := nats.Connect(cfg.URL, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to NATS: %w", err)
	}

	js, err := conn.JetStream()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to create JetStream context: %w", err)
	}

	return &Client{
		conn:   conn,
		js:     js,
		config: cfg,
	}, nil
}

// Conn returns the underlying NATS connection.
func (c *Client) Conn() *nats.Conn {
	return c.conn
}

// JetStream returns the JetStream context.
func (c *Client) JetStream() nats.JetStreamContext {
	return c.js
}

// Close closes the NATS connection.
func (c *Client) Close() {
	c.conn.Close()
}

// HealthCheck performs a health check on the NATS connection.
func (c *Client) HealthCheck(ctx context.Context) error {
	if !c.conn.IsConnected() {
		return fmt.Errorf("NATS not connected")
	}
	return nil
}

// IsConnected returns whether the client is connected.
func (c *Client) IsConnected() bool {
	return c.conn.IsConnected()
}

// Drain drains the connection.
func (c *Client) Drain() error {
	return c.conn.Drain()
}

// StreamConfig holds stream configuration.
type StreamConfig struct {
	Name        string
	Description string
	Subjects    []string
	Retention   nats.RetentionPolicy
	MaxAge      time.Duration
	MaxMsgs     int64
	MaxBytes    int64
	Replicas    int
	Storage     nats.StorageType
}

// DefaultStreamConfig returns a default stream configuration.
func DefaultStreamConfig(name string, subjects []string) StreamConfig {
	return StreamConfig{
		Name:      name,
		Subjects:  subjects,
		Retention: nats.LimitsPolicy,
		MaxAge:    7 * 24 * time.Hour, // 7 days
		MaxMsgs:   -1,                 // Unlimited
		MaxBytes:  -1,                 // Unlimited
		Replicas:  1,
		Storage:   nats.FileStorage,
	}
}

// CreateOrUpdateStream creates or updates a JetStream stream.
func (c *Client) CreateOrUpdateStream(ctx context.Context, cfg StreamConfig) (*nats.StreamInfo, error) {
	streamCfg := &nats.StreamConfig{
		Name:        cfg.Name,
		Description: cfg.Description,
		Subjects:    cfg.Subjects,
		Retention:   cfg.Retention,
		MaxAge:      cfg.MaxAge,
		MaxMsgs:     cfg.MaxMsgs,
		MaxBytes:    cfg.MaxBytes,
		Replicas:    cfg.Replicas,
		Storage:     cfg.Storage,
	}

	// Try to get existing stream
	stream, err := c.js.StreamInfo(cfg.Name)
	if err == nil {
		// Stream exists, update it
		stream, err = c.js.UpdateStream(streamCfg)
		if err != nil {
			return nil, fmt.Errorf("failed to update stream: %w", err)
		}
		return stream, nil
	}

	// Create new stream
	stream, err = c.js.AddStream(streamCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create stream: %w", err)
	}

	return stream, nil
}

// DeleteStream deletes a stream.
func (c *Client) DeleteStream(name string) error {
	return c.js.DeleteStream(name)
}

// StreamInfo returns information about a stream.
func (c *Client) StreamInfo(name string) (*nats.StreamInfo, error) {
	return c.js.StreamInfo(name)
}

// ConsumerConfig holds consumer configuration.
type ConsumerConfig struct {
	Durable       string
	Description   string
	DeliverPolicy nats.DeliverPolicy
	AckPolicy     nats.AckPolicy
	AckWait       time.Duration
	MaxDeliver    int
	FilterSubject string
	ReplayPolicy  nats.ReplayPolicy
	MaxAckPending int
}

// DefaultConsumerConfig returns a default consumer configuration.
func DefaultConsumerConfig(durable string) ConsumerConfig {
	return ConsumerConfig{
		Durable:       durable,
		DeliverPolicy: nats.DeliverAllPolicy,
		AckPolicy:     nats.AckExplicitPolicy,
		AckWait:       30 * time.Second,
		MaxDeliver:    5,
		ReplayPolicy:  nats.ReplayInstantPolicy,
		MaxAckPending: 1000,
	}
}

// CreateOrUpdateConsumer creates or updates a consumer.
func (c *Client) CreateOrUpdateConsumer(stream string, cfg ConsumerConfig) (*nats.ConsumerInfo, error) {
	consumerCfg := &nats.ConsumerConfig{
		Durable:       cfg.Durable,
		Description:   cfg.Description,
		DeliverPolicy: cfg.DeliverPolicy,
		AckPolicy:     cfg.AckPolicy,
		AckWait:       cfg.AckWait,
		MaxDeliver:    cfg.MaxDeliver,
		FilterSubject: cfg.FilterSubject,
		ReplayPolicy:  cfg.ReplayPolicy,
		MaxAckPending: cfg.MaxAckPending,
	}

	consumer, err := c.js.AddConsumer(stream, consumerCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create consumer: %w", err)
	}

	return consumer, nil
}

// DeleteConsumer deletes a consumer.
func (c *Client) DeleteConsumer(stream, consumer string) error {
	return c.js.DeleteConsumer(stream, consumer)
}

// Yousoon-specific streams

const (
	// StreamEvents is the main event stream.
	StreamEvents = "YOUSOON_EVENTS"
	// StreamNotifications is the notification stream.
	StreamNotifications = "YOUSOON_NOTIFICATIONS"
)

// Subject patterns for events.
const (
	// SubjectIdentityEvents matches all identity context events.
	SubjectIdentityEvents = "yousoon.events.identity.>"
	// SubjectPartnerEvents matches all partner context events.
	SubjectPartnerEvents = "yousoon.events.partner.>"
	// SubjectDiscoveryEvents matches all discovery context events.
	SubjectDiscoveryEvents = "yousoon.events.discovery.>"
	// SubjectBookingEvents matches all booking context events.
	SubjectBookingEvents = "yousoon.events.booking.>"
	// SubjectEngagementEvents matches all engagement context events.
	SubjectEngagementEvents = "yousoon.events.engagement.>"
	// SubjectNotificationEvents matches all notification events.
	SubjectNotificationEvents = "yousoon.notifications.>"
)

// InitializeStreams creates the required JetStream streams.
func (c *Client) InitializeStreams(ctx context.Context) error {
	// Main events stream
	eventsConfig := StreamConfig{
		Name:        StreamEvents,
		Description: "Yousoon domain events stream",
		Subjects: []string{
			"yousoon.events.>",
		},
		Retention: nats.WorkQueuePolicy,
		MaxAge:    30 * 24 * time.Hour, // 30 days
		Replicas:  1,
		Storage:   nats.FileStorage,
	}

	if _, err := c.CreateOrUpdateStream(ctx, eventsConfig); err != nil {
		return fmt.Errorf("failed to create events stream: %w", err)
	}

	// Notifications stream
	notifConfig := StreamConfig{
		Name:        StreamNotifications,
		Description: "Yousoon notifications stream",
		Subjects: []string{
			"yousoon.notifications.>",
		},
		Retention: nats.WorkQueuePolicy,
		MaxAge:    7 * 24 * time.Hour, // 7 days
		Replicas:  1,
		Storage:   nats.FileStorage,
	}

	if _, err := c.CreateOrUpdateStream(ctx, notifConfig); err != nil {
		return fmt.Errorf("failed to create notifications stream: %w", err)
	}

	return nil
}
