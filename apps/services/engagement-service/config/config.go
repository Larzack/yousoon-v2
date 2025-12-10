package config

import (
	"os"
	"strconv"
)

type Config struct {
	ServerPort        string
	ServerHost        string
	Environment       string
	ServiceName       string
	GraphQLPath       string
	PlaygroundEnabled bool
	MongoURI          string
	MongoDatabase     string
	RedisAddr         string
	RedisPassword     string
	RedisDB           int
	NatsURL           string
	JaegerEndpoint    string
	MetricsPort       string
	JWTSecret         string
}

func Load() *Config {
	return &Config{
		ServerPort:        getEnv("SERVER_PORT", "8084"),
		ServerHost:        getEnv("SERVER_HOST", "0.0.0.0"),
		Environment:       getEnv("ENVIRONMENT", "development"),
		ServiceName:       getEnv("SERVICE_NAME", "engagement-service"),
		GraphQLPath:       getEnv("GRAPHQL_PATH", "/graphql"),
		PlaygroundEnabled: getEnvBool("PLAYGROUND_ENABLED", true),
		MongoURI:          getEnv("MONGO_URI", "mongodb://localhost:27017"),
		MongoDatabase:     getEnv("MONGO_DATABASE", "engagement_db"),
		RedisAddr:         getEnv("REDIS_ADDR", "localhost:6379"),
		RedisPassword:     getEnv("REDIS_PASSWORD", ""),
		RedisDB:           getEnvInt("REDIS_DB", 0),
		NatsURL:           getEnv("NATS_URL", "nats://localhost:4222"),
		JaegerEndpoint:    getEnv("JAEGER_ENDPOINT", "http://localhost:14268/api/traces"),
		MetricsPort:       getEnv("METRICS_PORT", "9094"),
		JWTSecret:         getEnv("JWT_SECRET", "your-secret-key-change-in-prod"),
	}
}

func (c *Config) GetServerAddr() string {
	return c.ServerHost + ":" + c.ServerPort
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
