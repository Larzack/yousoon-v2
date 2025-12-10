// Package config provides configuration for the Discovery service.
package config

import (
	"strings"
	"time"

	sharedConfig "github.com/yousoon/shared/config"
)

// Config holds all configuration for the Discovery service.
type Config struct {
	// Service
	ServiceName string
	Environment string
	Version     string

	// Server
	HTTPPort string
	GRPCPort string

	// MongoDB
	MongoURI         string
	MongoDB          string
	MongoTimeout     time.Duration
	MongoMaxPoolSize uint64

	// Elasticsearch
	ElasticsearchURLs     []string
	ElasticsearchUsername string
	ElasticsearchPassword string

	// Redis
	RedisAddr     string
	RedisPassword string
	RedisDB       int

	// NATS
	NATSURL       string
	NATSClusterID string

	// gRPC Clients
	PartnerServiceAddr string

	// Schema Registry
	SchemaRegistryURL string

	// Observability
	JaegerEndpoint string
	MetricsPort    string
	LogLevel       string
	LogFormat      string

	// Cache
	CacheTTL         time.Duration
	CategoryCacheTTL time.Duration

	// Rate Limiting
	RateLimitRequests int
	RateLimitWindow   time.Duration

	// Search
	DefaultSearchRadius float64
	MaxSearchRadius     float64
	DefaultPageSize     int
	MaxPageSize         int
}

// Load loads configuration from environment variables.
func Load() (*Config, error) {
	esURLs := sharedConfig.GetEnv("ELASTICSEARCH_URLS", "http://localhost:9200")

	cfg := &Config{
		// Service
		ServiceName: sharedConfig.GetEnv("SERVICE_NAME", "discovery-service"),
		Environment: sharedConfig.GetEnv("ENVIRONMENT", "development"),
		Version:     sharedConfig.GetEnv("VERSION", "1.0.0"),

		// Server
		HTTPPort: sharedConfig.GetEnv("HTTP_PORT", "8080"),
		GRPCPort: sharedConfig.GetEnv("GRPC_PORT", "9090"),

		// MongoDB
		MongoURI:         sharedConfig.GetEnv("MONGO_URI", "mongodb://localhost:27017"),
		MongoDB:          sharedConfig.GetEnv("MONGO_DB", "discovery_db"),
		MongoTimeout:     sharedConfig.GetEnvDuration("MONGO_TIMEOUT", 10*time.Second),
		MongoMaxPoolSize: uint64(sharedConfig.GetEnvInt("MONGO_MAX_POOL_SIZE", 100)),

		// Elasticsearch
		ElasticsearchURLs:     strings.Split(esURLs, ","),
		ElasticsearchUsername: sharedConfig.GetEnv("ELASTICSEARCH_USERNAME", ""),
		ElasticsearchPassword: sharedConfig.GetEnv("ELASTICSEARCH_PASSWORD", ""),

		// Redis
		RedisAddr:     sharedConfig.GetEnv("REDIS_ADDR", "localhost:6379"),
		RedisPassword: sharedConfig.GetEnv("REDIS_PASSWORD", ""),
		RedisDB:       sharedConfig.GetEnvInt("REDIS_DB", 0),

		// NATS
		NATSURL:       sharedConfig.GetEnv("NATS_URL", "nats://localhost:4222"),
		NATSClusterID: sharedConfig.GetEnv("NATS_CLUSTER_ID", "yousoon-cluster"),

		// gRPC Clients
		PartnerServiceAddr: sharedConfig.GetEnv("PARTNER_SERVICE_ADDR", "localhost:9091"),

		// Schema Registry
		SchemaRegistryURL: sharedConfig.GetEnv("SCHEMA_REGISTRY_URL", "http://localhost:8081"),

		// Observability
		JaegerEndpoint: sharedConfig.GetEnv("JAEGER_ENDPOINT", "http://localhost:14268/api/traces"),
		MetricsPort:    sharedConfig.GetEnv("METRICS_PORT", "9100"),
		LogLevel:       sharedConfig.GetEnv("LOG_LEVEL", "info"),
		LogFormat:      sharedConfig.GetEnv("LOG_FORMAT", "json"),

		// Cache
		CacheTTL:         sharedConfig.GetEnvDuration("CACHE_TTL", 5*time.Minute),
		CategoryCacheTTL: sharedConfig.GetEnvDuration("CATEGORY_CACHE_TTL", 1*time.Hour),

		// Rate Limiting
		RateLimitRequests: sharedConfig.GetEnvInt("RATE_LIMIT_REQUESTS", 100),
		RateLimitWindow:   sharedConfig.GetEnvDuration("RATE_LIMIT_WINDOW", 1*time.Minute),

		// Search
		DefaultSearchRadius: float64(sharedConfig.GetEnvInt("DEFAULT_SEARCH_RADIUS", 10)),
		MaxSearchRadius:     float64(sharedConfig.GetEnvInt("MAX_SEARCH_RADIUS", 50)),
		DefaultPageSize:     sharedConfig.GetEnvInt("DEFAULT_PAGE_SIZE", 20),
		MaxPageSize:         sharedConfig.GetEnvInt("MAX_PAGE_SIZE", 100),
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
