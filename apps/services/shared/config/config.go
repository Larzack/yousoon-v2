// Package config provides configuration utilities.
package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// GetEnv retrieves an environment variable or returns a default value.
func GetEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// GetEnvInt retrieves an environment variable as an integer.
func GetEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}

// GetEnvInt64 retrieves an environment variable as an int64.
func GetEnvInt64(key string, defaultValue int64) int64 {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.ParseInt(value, 10, 64); err == nil {
			return intVal
		}
	}
	return defaultValue
}

// GetEnvBool retrieves an environment variable as a boolean.
func GetEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolVal, err := strconv.ParseBool(value); err == nil {
			return boolVal
		}
	}
	return defaultValue
}

// GetEnvDuration retrieves an environment variable as a duration.
func GetEnvDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}

// ServiceConfig holds common service configuration.
type ServiceConfig struct {
	// Name is the service name.
	Name string
	// Version is the service version.
	Version string
	// Environment is the deployment environment (dev, staging, prod).
	Environment string
	// LogLevel is the logging level.
	LogLevel string
	// GRPCPort is the gRPC server port.
	GRPCPort int
	// HTTPPort is the HTTP server port (for GraphQL).
	HTTPPort int
	// MetricsPort is the metrics server port.
	MetricsPort int
}

// NewServiceConfig creates a service config from environment variables.
func NewServiceConfig(name, version string) ServiceConfig {
	return ServiceConfig{
		Name:        name,
		Version:     version,
		Environment: GetEnv("ENVIRONMENT", "development"),
		LogLevel:    GetEnv("LOG_LEVEL", "info"),
		GRPCPort:    GetEnvInt("GRPC_PORT", 50051),
		HTTPPort:    GetEnvInt("HTTP_PORT", 8080),
		MetricsPort: GetEnvInt("METRICS_PORT", 9090),
	}
}

// IsDevelopment returns true if running in development mode.
func (c ServiceConfig) IsDevelopment() bool {
	return c.Environment == "development" || c.Environment == "dev"
}

// IsProduction returns true if running in production mode.
func (c ServiceConfig) IsProduction() bool {
	return c.Environment == "production" || c.Environment == "prod"
}

// MongoDBConfig holds MongoDB configuration.
type MongoDBConfig struct {
	URI                    string
	Database               string
	MaxPoolSize            uint64
	MinPoolSize            uint64
	MaxConnIdleTime        time.Duration
	ConnectTimeout         time.Duration
	ServerSelectionTimeout time.Duration
}

// buildMongoDBURI constructs a MongoDB URI from separate environment variables.
func buildMongoDBURI() string {
	// Check if MONGODB_URI is explicitly set
	if uri := os.Getenv("MONGODB_URI"); uri != "" {
		return uri
	}

	// Build URI from separate variables
	host := GetEnv("MONGODB_HOST", "localhost")
	port := GetEnv("MONGODB_PORT", "27017")
	username := GetEnv("MONGODB_USERNAME", "")
	password := GetEnv("MONGODB_PASSWORD", "")

	if username != "" && password != "" {
		return fmt.Sprintf("mongodb://%s:%s@%s:%s", username, password, host, port)
	}
	return fmt.Sprintf("mongodb://%s:%s", host, port)
}

// NewMongoDBConfig creates a MongoDB config from environment variables.
func NewMongoDBConfig(database string) MongoDBConfig {
	return MongoDBConfig{
		URI:                    buildMongoDBURI(),
		Database:               GetEnv("MONGODB_DATABASE", database),
		MaxPoolSize:            uint64(GetEnvInt("MONGODB_MAX_POOL_SIZE", 100)),
		MinPoolSize:            uint64(GetEnvInt("MONGODB_MIN_POOL_SIZE", 10)),
		MaxConnIdleTime:        GetEnvDuration("MONGODB_MAX_CONN_IDLE_TIME", 30*time.Minute),
		ConnectTimeout:         GetEnvDuration("MONGODB_CONNECT_TIMEOUT", 10*time.Second),
		ServerSelectionTimeout: GetEnvDuration("MONGODB_SERVER_SELECTION_TIMEOUT", 5*time.Second),
	}
}

// RedisConfig holds Redis configuration.
type RedisConfig struct {
	Address      string
	Password     string
	DB           int
	PoolSize     int
	MinIdleConns int
	DialTimeout  time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

// NewRedisConfig creates a Redis config from environment variables.
func NewRedisConfig() RedisConfig {
	return RedisConfig{
		Address:      GetEnv("REDIS_ADDRESS", "localhost:6379"),
		Password:     GetEnv("REDIS_PASSWORD", ""),
		DB:           GetEnvInt("REDIS_DB", 0),
		PoolSize:     GetEnvInt("REDIS_POOL_SIZE", 100),
		MinIdleConns: GetEnvInt("REDIS_MIN_IDLE_CONNS", 10),
		DialTimeout:  GetEnvDuration("REDIS_DIAL_TIMEOUT", 5*time.Second),
		ReadTimeout:  GetEnvDuration("REDIS_READ_TIMEOUT", 3*time.Second),
		WriteTimeout: GetEnvDuration("REDIS_WRITE_TIMEOUT", 3*time.Second),
	}
}

// NATSConfig holds NATS configuration.
type NATSConfig struct {
	URL           string
	Name          string
	MaxReconnects int
	ReconnectWait time.Duration
	Timeout       time.Duration
	Token         string
	User          string
	Password      string
}

// NewNATSConfig creates a NATS config from environment variables.
func NewNATSConfig(serviceName string) NATSConfig {
	return NATSConfig{
		URL:           GetEnv("NATS_URL", "nats://localhost:4222"),
		Name:          serviceName,
		MaxReconnects: GetEnvInt("NATS_MAX_RECONNECTS", -1),
		ReconnectWait: GetEnvDuration("NATS_RECONNECT_WAIT", 2*time.Second),
		Timeout:       GetEnvDuration("NATS_TIMEOUT", 10*time.Second),
		Token:         GetEnv("NATS_TOKEN", ""),
		User:          GetEnv("NATS_USER", ""),
		Password:      GetEnv("NATS_PASSWORD", ""),
	}
}

// JWTConfig holds JWT configuration.
type JWTConfig struct {
	Secret          string
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
	Issuer          string
	Audience        string
}

// NewJWTConfig creates a JWT config from environment variables.
func NewJWTConfig() JWTConfig {
	return JWTConfig{
		Secret:          GetEnv("JWT_SECRET", "your-secret-key-change-in-production"),
		AccessTokenTTL:  GetEnvDuration("JWT_ACCESS_TOKEN_TTL", 6*time.Hour),
		RefreshTokenTTL: GetEnvDuration("JWT_REFRESH_TOKEN_TTL", 30*24*time.Hour),
		Issuer:          GetEnv("JWT_ISSUER", "yousoon"),
		Audience:        GetEnv("JWT_AUDIENCE", "yousoon-api"),
	}
}

// S3Config holds AWS S3 configuration.
type S3Config struct {
	Region          string
	Bucket          string
	AccessKeyID     string
	SecretAccessKey string
	Endpoint        string
}

// NewS3Config creates an S3 config from environment variables.
func NewS3Config() S3Config {
	return S3Config{
		Region:          GetEnv("AWS_REGION", "eu-west-1"),
		Bucket:          GetEnv("S3_BUCKET", "yousoon-media"),
		AccessKeyID:     GetEnv("AWS_ACCESS_KEY_ID", ""),
		SecretAccessKey: GetEnv("AWS_SECRET_ACCESS_KEY", ""),
		Endpoint:        GetEnv("S3_ENDPOINT", ""),
	}
}

// FullConfig holds all configuration for a service.
type FullConfig struct {
	Service ServiceConfig
	MongoDB MongoDBConfig
	Redis   RedisConfig
	NATS    NATSConfig
	JWT     JWTConfig
	S3      S3Config
}

// NewFullConfig creates a full configuration from environment variables.
func NewFullConfig(serviceName, version, database string) FullConfig {
	return FullConfig{
		Service: NewServiceConfig(serviceName, version),
		MongoDB: NewMongoDBConfig(database),
		Redis:   NewRedisConfig(),
		NATS:    NewNATSConfig(serviceName),
		JWT:     NewJWTConfig(),
		S3:      NewS3Config(),
	}
}
