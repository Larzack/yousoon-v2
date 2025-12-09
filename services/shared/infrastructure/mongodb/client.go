// Package mongodb provides MongoDB infrastructure components.
package mongodb

import (
	"context"
	"fmt"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Config holds MongoDB connection configuration.
type Config struct {
	// URI is the MongoDB connection string.
	URI string
	// Database is the database name.
	Database string
	// MaxPoolSize is the maximum number of connections in the pool.
	MaxPoolSize uint64
	// MinPoolSize is the minimum number of connections in the pool.
	MinPoolSize uint64
	// MaxConnIdleTime is the maximum time a connection can remain idle.
	MaxConnIdleTime time.Duration
	// ConnectTimeout is the timeout for establishing connections.
	ConnectTimeout time.Duration
	// ServerSelectionTimeout is the timeout for server selection.
	ServerSelectionTimeout time.Duration
}

// DefaultConfig returns a default MongoDB configuration.
func DefaultConfig() Config {
	return Config{
		URI:                    "mongodb://localhost:27017",
		Database:               "yousoon",
		MaxPoolSize:            100,
		MinPoolSize:            10,
		MaxConnIdleTime:        30 * time.Minute,
		ConnectTimeout:         10 * time.Second,
		ServerSelectionTimeout: 5 * time.Second,
	}
}

// Client wraps the MongoDB client with additional functionality.
type Client struct {
	client   *mongo.Client
	database *mongo.Database
	config   Config
	mu       sync.RWMutex
}

// NewClient creates a new MongoDB client with the given configuration.
func NewClient(ctx context.Context, cfg Config) (*Client, error) {
	opts := options.Client().
		ApplyURI(cfg.URI).
		SetMaxPoolSize(cfg.MaxPoolSize).
		SetMinPoolSize(cfg.MinPoolSize).
		SetMaxConnIdleTime(cfg.MaxConnIdleTime).
		SetConnectTimeout(cfg.ConnectTimeout).
		SetServerSelectionTimeout(cfg.ServerSelectionTimeout)

	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// Verify connection
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	return &Client{
		client:   client,
		database: client.Database(cfg.Database),
		config:   cfg,
	}, nil
}

// Database returns the MongoDB database instance.
func (c *Client) Database() *mongo.Database {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.database
}

// Collection returns a collection from the database.
func (c *Client) Collection(name string) *mongo.Collection {
	return c.Database().Collection(name)
}

// Client returns the underlying MongoDB client.
func (c *Client) Client() *mongo.Client {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.client
}

// Ping verifies the connection to MongoDB.
func (c *Client) Ping(ctx context.Context) error {
	return c.client.Ping(ctx, readpref.Primary())
}

// Close disconnects from MongoDB.
func (c *Client) Close(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.client.Disconnect(ctx)
}

// HealthCheck performs a health check on the MongoDB connection.
func (c *Client) HealthCheck(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := c.Ping(ctx); err != nil {
		return fmt.Errorf("MongoDB health check failed: %w", err)
	}
	return nil
}

// WithDatabase returns a new client using a different database.
func (c *Client) WithDatabase(name string) *Client {
	return &Client{
		client:   c.client,
		database: c.client.Database(name),
		config:   c.config,
	}
}

// Stats returns connection pool statistics.
type Stats struct {
	TotalConnections   int
	AvailableConns     int
	InUseConns         int
	WaitQueueLength    int
	WaitQueueTimeoutMS int64
}

// RunCommand runs a MongoDB command.
func (c *Client) RunCommand(ctx context.Context, cmd interface{}) (*mongo.SingleResult, error) {
	result := c.database.RunCommand(ctx, cmd)
	if result.Err() != nil {
		return nil, result.Err()
	}
	return result, nil
}

// ListCollections returns the names of all collections in the database.
func (c *Client) ListCollections(ctx context.Context) ([]string, error) {
	names, err := c.database.ListCollectionNames(ctx, map[string]interface{}{})
	if err != nil {
		return nil, fmt.Errorf("failed to list collections: %w", err)
	}
	return names, nil
}

// CreateCollection creates a new collection with optional options.
func (c *Client) CreateCollection(ctx context.Context, name string, opts ...*options.CreateCollectionOptions) error {
	if err := c.database.CreateCollection(ctx, name, opts...); err != nil {
		return fmt.Errorf("failed to create collection %s: %w", name, err)
	}
	return nil
}

// DropCollection drops a collection.
func (c *Client) DropCollection(ctx context.Context, name string) error {
	if err := c.Collection(name).Drop(ctx); err != nil {
		return fmt.Errorf("failed to drop collection %s: %w", name, err)
	}
	return nil
}
