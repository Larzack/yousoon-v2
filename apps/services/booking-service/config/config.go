package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	// Server
	ServerPort  string
	ServerHost  string
	Environment string
	ServiceName string

	// GraphQL
	GraphQLPath       string
	PlaygroundEnabled bool

	// MongoDB
	MongoURI      string
	MongoDatabase string

	// Redis
	RedisAddr     string
	RedisPassword string
	RedisDB       int

	// NATS
	NatsURL     string
	NatsCluster string

	// gRPC clients
	IdentityServiceAddr  string
	DiscoveryServiceAddr string

	// Booking settings
	BookingExpirationMinutes int
	MaxBookingsPerUser       int

	// Observability
	JaegerEndpoint string
	MetricsPort    string

	// JWT
	JWTSecret string
}

func Load() *Config {
	return &Config{
		// Server
		ServerPort:  getEnv("SERVER_PORT", "8083"),
		ServerHost:  getEnv("SERVER_HOST", "0.0.0.0"),
		Environment: getEnv("ENVIRONMENT", "development"),
		ServiceName: getEnv("SERVICE_NAME", "booking-service"),

		// GraphQL
		GraphQLPath:       getEnv("GRAPHQL_PATH", "/graphql"),
		PlaygroundEnabled: getEnvBool("PLAYGROUND_ENABLED", true),

		// MongoDB
		MongoURI:      getEnv("MONGO_URI", "mongodb://localhost:27017"),
		MongoDatabase: getEnv("MONGO_DATABASE", "booking_db"),

		// Redis
		RedisAddr:     getEnv("REDIS_ADDR", "localhost:6379"),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),
		RedisDB:       getEnvInt("REDIS_DB", 0),

		// NATS
		NatsURL:     getEnv("NATS_URL", "nats://localhost:4222"),
		NatsCluster: getEnv("NATS_CLUSTER", "yousoon-cluster"),

		// gRPC clients
		IdentityServiceAddr:  getEnv("IDENTITY_SERVICE_ADDR", "identity-service:50051"),
		DiscoveryServiceAddr: getEnv("DISCOVERY_SERVICE_ADDR", "discovery-service:50052"),

		// Booking settings
		BookingExpirationMinutes: getEnvInt("BOOKING_EXPIRATION_MINUTES", 30),
		MaxBookingsPerUser:       getEnvInt("MAX_BOOKINGS_PER_USER", 5),

		// Observability
		JaegerEndpoint: getEnv("JAEGER_ENDPOINT", "http://localhost:14268/api/traces"),
		MetricsPort:    getEnv("METRICS_PORT", "9093"),

		// JWT
		JWTSecret: getEnv("JWT_SECRET", "your-secret-key-change-in-prod"),
	}
}

func (c *Config) GetServerAddr() string {
	return c.ServerHost + ":" + c.ServerPort
}

func (c *Config) GetBookingExpiration() time.Duration {
	return time.Duration(c.BookingExpirationMinutes) * time.Minute
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolVal, err := strconv.ParseBool(value); err == nil {
			return boolVal
		}
	}
	return defaultValue
}
