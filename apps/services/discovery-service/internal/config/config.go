// Package config provides configuration for the Discovery service.
package config

import (
	"time"

	"github.com/yousoon/shared/config"
)

// Config holds all configuration for the Discovery service.
type Config struct {
	// Service
	ServiceName string `env:"SERVICE_NAME" envDefault:"discovery-service"`
	Environment string `env:"ENVIRONMENT" envDefault:"development"`
	Version     string `env:"VERSION" envDefault:"1.0.0"`

	// Server
	HTTPPort string `env:"HTTP_PORT" envDefault:"8080"`
	GRPCPort string `env:"GRPC_PORT" envDefault:"9090"`

	// MongoDB
	MongoURI         string        `env:"MONGO_URI" envDefault:"mongodb://localhost:27017"`
	MongoDB          string        `env:"MONGO_DB" envDefault:"discovery_db"`
	MongoTimeout     time.Duration `env:"MONGO_TIMEOUT" envDefault:"10s"`
	MongoMaxPoolSize uint64        `env:"MONGO_MAX_POOL_SIZE" envDefault:"100"`

	// Elasticsearch
	ElasticsearchURLs     []string `env:"ELASTICSEARCH_URLS" envDefault:"http://localhost:9200"`
	ElasticsearchUsername string   `env:"ELASTICSEARCH_USERNAME"`
	ElasticsearchPassword string   `env:"ELASTICSEARCH_PASSWORD"`

	// Redis
	RedisAddr     string `env:"REDIS_ADDR" envDefault:"localhost:6379"`
	RedisPassword string `env:"REDIS_PASSWORD"`
	RedisDB       int    `env:"REDIS_DB" envDefault:"0"`

	// NATS
	NATSURL       string `env:"NATS_URL" envDefault:"nats://localhost:4222"`
	NATSClusterID string `env:"NATS_CLUSTER_ID" envDefault:"yousoon-cluster"`

	// gRPC Clients
	PartnerServiceAddr string `env:"PARTNER_SERVICE_ADDR" envDefault:"localhost:9091"`

	// Schema Registry
	SchemaRegistryURL string `env:"SCHEMA_REGISTRY_URL" envDefault:"http://localhost:8081"`

	// Observability
	JaegerEndpoint string `env:"JAEGER_ENDPOINT" envDefault:"http://localhost:14268/api/traces"`
	MetricsPort    string `env:"METRICS_PORT" envDefault:"9100"`
	LogLevel       string `env:"LOG_LEVEL" envDefault:"info"`
	LogFormat      string `env:"LOG_FORMAT" envDefault:"json"`

	// Cache
	CacheTTL         time.Duration `env:"CACHE_TTL" envDefault:"5m"`
	CategoryCacheTTL time.Duration `env:"CATEGORY_CACHE_TTL" envDefault:"1h"`

	// Rate Limiting
	RateLimitRequests int           `env:"RATE_LIMIT_REQUESTS" envDefault:"100"`
	RateLimitWindow   time.Duration `env:"RATE_LIMIT_WINDOW" envDefault:"1m"`

	// Search
	DefaultSearchRadius float64 `env:"DEFAULT_SEARCH_RADIUS" envDefault:"10.0"` // km
	MaxSearchRadius     float64 `env:"MAX_SEARCH_RADIUS" envDefault:"50.0"`     // km
	DefaultPageSize     int     `env:"DEFAULT_PAGE_SIZE" envDefault:"20"`
	MaxPageSize         int     `env:"MAX_PAGE_SIZE" envDefault:"100"`
}

// Load loads configuration from environment variables.
func Load() (*Config, error) {
	cfg := &Config{}
	if err := config.Load(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

// IsDevelopment returns true if running in development mode.
func (c *Config) IsDevelopment() bool {
	return c.Environment == "development"
}

// IsProduction returns true if running in production mode.
func (c *Config) IsProduction() bool {
	return c.Environment == "production"
}
