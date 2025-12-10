// Package config provides configuration for the Partner service.
package config

import (
	"os"
	"strconv"
	"time"
)

// Config holds all configuration for the Partner service.
type Config struct {
	// Server configuration
	Server ServerConfig

	// GraphQL configuration
	GraphQL GraphQLConfig

	// MongoDB configuration
	MongoDB MongoDBConfig

	// NATS configuration
	NATS NATSConfig

	// Redis configuration
	Redis RedisConfig

	// Observability configuration
	Observability ObservabilityConfig

	// Service discovery
	ServiceDiscovery ServiceDiscoveryConfig
}

// ServerConfig holds server-related configuration.
type ServerConfig struct {
	Host         string
	Port         int
	GRPCPort     int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

// GraphQLConfig holds GraphQL-related configuration.
type GraphQLConfig struct {
	PlaygroundEnabled       bool
	IntrospectionEnabled    bool
	QueryComplexityLimit    int
	PersistedQueriesEnabled bool
}

// MongoDBConfig holds MongoDB connection configuration.
type MongoDBConfig struct {
	URI             string
	Database        string
	MaxPoolSize     uint64
	MinPoolSize     uint64
	MaxConnIdleTime time.Duration
	ConnectTimeout  time.Duration
}

// NATSConfig holds NATS connection configuration.
type NATSConfig struct {
	URL            string
	ClusterID      string
	ClientID       string
	ConnectTimeout time.Duration
	MaxReconnects  int
	ReconnectWait  time.Duration
}

// RedisConfig holds Redis connection configuration.
type RedisConfig struct {
	Host         string
	Port         int
	Password     string
	DB           int
	PoolSize     int
	MinIdleConns int
}

// ObservabilityConfig holds observability configuration.
type ObservabilityConfig struct {
	// Logging
	LogLevel  string
	LogFormat string // json, text

	// Tracing (OpenTelemetry)
	TracingEnabled bool
	JaegerEndpoint string
	SamplingRate   float64

	// Metrics (Prometheus)
	MetricsEnabled bool
	MetricsPath    string
	MetricsPort    int
}

// ServiceDiscoveryConfig holds service discovery configuration.
type ServiceDiscoveryConfig struct {
	RegistryURL       string
	ServiceName       string
	ServiceURL        string
	HeartbeatInterval time.Duration
}

// Load loads configuration from environment variables.
func Load() *Config {
	return &Config{
		Server: ServerConfig{
			Host:         getEnv("SERVER_HOST", "0.0.0.0"),
			Port:         getEnvAsInt("SERVER_PORT", 8080),
			GRPCPort:     getEnvAsInt("GRPC_PORT", 9090),
			ReadTimeout:  getEnvAsDuration("SERVER_READ_TIMEOUT", 30*time.Second),
			WriteTimeout: getEnvAsDuration("SERVER_WRITE_TIMEOUT", 30*time.Second),
			IdleTimeout:  getEnvAsDuration("SERVER_IDLE_TIMEOUT", 120*time.Second),
		},
		GraphQL: GraphQLConfig{
			PlaygroundEnabled:       getEnvAsBool("GRAPHQL_PLAYGROUND_ENABLED", true),
			IntrospectionEnabled:    getEnvAsBool("GRAPHQL_INTROSPECTION_ENABLED", true),
			QueryComplexityLimit:    getEnvAsInt("GRAPHQL_QUERY_COMPLEXITY_LIMIT", 200),
			PersistedQueriesEnabled: getEnvAsBool("GRAPHQL_PERSISTED_QUERIES_ENABLED", false),
		},
		MongoDB: MongoDBConfig{
			URI:             getEnv("MONGODB_URI", "mongodb://localhost:27017"),
			Database:        getEnv("MONGODB_DATABASE", "partner_db"),
			MaxPoolSize:     uint64(getEnvAsInt("MONGODB_MAX_POOL_SIZE", 100)),
			MinPoolSize:     uint64(getEnvAsInt("MONGODB_MIN_POOL_SIZE", 10)),
			MaxConnIdleTime: getEnvAsDuration("MONGODB_MAX_CONN_IDLE_TIME", 10*time.Minute),
			ConnectTimeout:  getEnvAsDuration("MONGODB_CONNECT_TIMEOUT", 10*time.Second),
		},
		NATS: NATSConfig{
			URL:            getEnv("NATS_URL", "nats://localhost:4222"),
			ClusterID:      getEnv("NATS_CLUSTER_ID", "yousoon-cluster"),
			ClientID:       getEnv("NATS_CLIENT_ID", "partner-service"),
			ConnectTimeout: getEnvAsDuration("NATS_CONNECT_TIMEOUT", 10*time.Second),
			MaxReconnects:  getEnvAsInt("NATS_MAX_RECONNECTS", 10),
			ReconnectWait:  getEnvAsDuration("NATS_RECONNECT_WAIT", 2*time.Second),
		},
		Redis: RedisConfig{
			Host:         getEnv("REDIS_HOST", "localhost"),
			Port:         getEnvAsInt("REDIS_PORT", 6379),
			Password:     getEnv("REDIS_PASSWORD", ""),
			DB:           getEnvAsInt("REDIS_DB", 0),
			PoolSize:     getEnvAsInt("REDIS_POOL_SIZE", 10),
			MinIdleConns: getEnvAsInt("REDIS_MIN_IDLE_CONNS", 5),
		},
		Observability: ObservabilityConfig{
			LogLevel:       getEnv("LOG_LEVEL", "info"),
			LogFormat:      getEnv("LOG_FORMAT", "json"),
			TracingEnabled: getEnvAsBool("TRACING_ENABLED", true),
			JaegerEndpoint: getEnv("JAEGER_ENDPOINT", "http://localhost:14268/api/traces"),
			SamplingRate:   getEnvAsFloat("TRACING_SAMPLING_RATE", 0.1),
			MetricsEnabled: getEnvAsBool("METRICS_ENABLED", true),
			MetricsPath:    getEnv("METRICS_PATH", "/metrics"),
			MetricsPort:    getEnvAsInt("METRICS_PORT", 9091),
		},
		ServiceDiscovery: ServiceDiscoveryConfig{
			RegistryURL:       getEnv("REGISTRY_URL", "http://localhost:4000"),
			ServiceName:       getEnv("SERVICE_NAME", "partner-service"),
			ServiceURL:        getEnv("SERVICE_URL", "http://localhost:8080"),
			HeartbeatInterval: getEnvAsDuration("HEARTBEAT_INTERVAL", 30*time.Second),
		},
	}
}

// Helper functions

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	if value, exists := os.LookupEnv(key); exists {
		if boolVal, err := strconv.ParseBool(value); err == nil {
			return boolVal
		}
	}
	return defaultValue
}

func getEnvAsFloat(key string, defaultValue float64) float64 {
	if value, exists := os.LookupEnv(key); exists {
		if floatVal, err := strconv.ParseFloat(value, 64); err == nil {
			return floatVal
		}
	}
	return defaultValue
}

func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	if value, exists := os.LookupEnv(key); exists {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}
