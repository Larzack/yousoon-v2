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

	// MongoDB
	MongoURI      string
	MongoDatabase string

	// NATS
	NatsURL     string
	NatsCluster string

	// OneSignal (Push)
	OneSignalAppID  string
	OneSignalAPIKey string

	// AWS (Email/SMS)
	AWSRegion       string
	AWSSESFromEmail string
	AWSSNSTopicARN  string

	// Observability
	JaegerEndpoint string
	MetricsPort    string
}

func Load() *Config {
	return &Config{
		ServerPort:        getEnv("SERVER_PORT", "8085"),
		ServerHost:        getEnv("SERVER_HOST", "0.0.0.0"),
		Environment:       getEnv("ENVIRONMENT", "development"),
		ServiceName:       getEnv("SERVICE_NAME", "notification-service"),
		GraphQLPath:       getEnv("GRAPHQL_PATH", "/graphql"),
		PlaygroundEnabled: getEnvBool("PLAYGROUND_ENABLED", true),

		MongoURI:      getEnv("MONGO_URI", "mongodb://localhost:27017"),
		MongoDatabase: getEnv("MONGO_DATABASE", "notification_db"),

		NatsURL:     getEnv("NATS_URL", "nats://localhost:4222"),
		NatsCluster: getEnv("NATS_CLUSTER", "yousoon-cluster"),

		OneSignalAppID:  getEnv("ONESIGNAL_APP_ID", ""),
		OneSignalAPIKey: getEnv("ONESIGNAL_API_KEY", ""),

		AWSRegion:       getEnv("AWS_REGION", "eu-west-1"),
		AWSSESFromEmail: getEnv("AWS_SES_FROM_EMAIL", "noreply@yousoon.com"),
		AWSSNSTopicARN:  getEnv("AWS_SNS_TOPIC_ARN", ""),

		JaegerEndpoint: getEnv("JAEGER_ENDPOINT", "http://localhost:14268/api/traces"),
		MetricsPort:    getEnv("METRICS_PORT", "9095"),
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

func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolVal, err := strconv.ParseBool(value); err == nil {
			return boolVal
		}
	}
	return defaultValue
}
