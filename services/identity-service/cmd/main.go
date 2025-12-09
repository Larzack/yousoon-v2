package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gorilla/websocket"

	mongodb "github.com/yousoon/services/identity/internal/infrastructure/mongodb"
	"github.com/yousoon/services/identity/internal/interface/graphql/resolver"
	"github.com/yousoon/services/shared/config"
	sharedmongo "github.com/yousoon/services/shared/infrastructure/mongodb"
	"github.com/yousoon/services/shared/infrastructure/nats"
)

const (
	defaultPort = "8080"
	serviceName = "identity-service"
)

func main() {
	// Initialize logger
	logger := observability.NewLogger(serviceName, os.Getenv("LOG_LEVEL"))
	logger.Info("Starting identity service...")

	// Get configuration from environment
	port := config.GetEnv("PORT", defaultPort)
	mongoURI := config.GetEnv("MONGODB_URI", "mongodb://localhost:27017")
	mongoDatabase := config.GetEnv("MONGODB_DATABASE", "identity_db")
	natsURL := config.GetEnv("NATS_URL", "nats://localhost:4222")
	enablePlayground := config.GetEnvBool("ENABLE_PLAYGROUND", true)

	// Initialize MongoDB client
	mongoClient, err := sharedmongo.NewClient(context.Background(), sharedmongo.Config{
		URI:            mongoURI,
		Database:       mongoDatabase,
		ConnectTimeout: 10 * time.Second,
		MaxPoolSize:    100,
	})
	if err != nil {
		logger.Error("Failed to connect to MongoDB", "error", err)
		os.Exit(1)
	}
	defer mongoClient.Disconnect(context.Background())

	logger.Info("Connected to MongoDB", "database", mongoDatabase)

	// Initialize NATS client
	natsClient, err := nats.NewClient(nats.Config{
		URL:       natsURL,
		Name:      serviceName,
		Reconnect: true,
	})
	if err != nil {
		logger.Error("Failed to connect to NATS", "error", err)
		os.Exit(1)
	}
	defer natsClient.Close()

	logger.Info("Connected to NATS")

	// Initialize event publisher
	eventPublisher, err := nats.NewEventPublisher(natsClient, "identity")
	if err != nil {
		logger.Error("Failed to create event publisher", "error", err)
		os.Exit(1)
	}

	// Initialize repositories
	userRepo := mongodb.NewUserRepository(mongoClient.Database())

	// Ensure indexes
	if err := userRepo.EnsureIndexes(context.Background()); err != nil {
		logger.Warn("Failed to ensure indexes", "error", err)
	}

	// Initialize GraphQL resolver
	graphqlResolver := resolver.NewResolver(userRepo, eventPublisher)

	// Note: In a real setup, you would generate the schema first using:
	// go run github.com/99designs/gqlgen generate
	// For now, we'll create a placeholder server

	// Create HTTP server
	mux := http.NewServeMux()

	// Health check endpoint
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"healthy"}`))
	})

	// Ready check endpoint
	mux.HandleFunc("/ready", func(w http.ResponseWriter, r *http.Request) {
		// Check MongoDB connection
		if err := mongoClient.Ping(r.Context()); err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte(`{"status":"not ready","reason":"mongodb unavailable"}`))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ready"}`))
	})

	// GraphQL endpoint placeholder
	// After running gqlgen generate, replace this with:
	// srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: graphqlResolver}))
	mux.HandleFunc("/query", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"message":"GraphQL endpoint - run 'go generate ./...' to generate schema"}`))
	})

	// Playground endpoint (development only)
	if enablePlayground {
		mux.Handle("/", playground.Handler("Identity Service", "/query"))
		logger.Info("GraphQL Playground enabled at /")
	}

	// Create HTTP server with timeouts
	server := &http.Server{
		Addr:         ":" + port,
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 60 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// Start server in goroutine
	go func() {
		logger.Info("Server starting", "port", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("Server failed", "error", err)
			os.Exit(1)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Error("Server shutdown failed", "error", err)
		os.Exit(1)
	}

	logger.Info("Server stopped")

	// Suppress unused variable warnings
	_ = graphqlResolver
	_ = handler.New(nil)
	_ = extension.Introspection{}
	_ = lru.New(0)
	_ = transport.Options{}
	_ = websocket.Upgrader{}
}
