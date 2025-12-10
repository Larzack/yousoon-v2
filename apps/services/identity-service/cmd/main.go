package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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
	// Initialize structured logger
	logLevel := slog.LevelInfo
	if os.Getenv("LOG_LEVEL") == "debug" {
		logLevel = slog.LevelDebug
	}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel}))
	slog.SetDefault(logger)

	slog.Info("Starting identity service...")

	// Get configuration from environment
	port := config.GetEnv("PORT", defaultPort)
	mongoURI := config.GetEnv("MONGODB_URI", "mongodb://localhost:27017")
	mongoDatabase := config.GetEnv("MONGODB_DATABASE", "identity_db")
	natsURL := config.GetEnv("NATS_URL", "nats://localhost:4222")

	// Initialize MongoDB client
	mongoClient, err := sharedmongo.NewClient(context.Background(), sharedmongo.Config{
		URI:            mongoURI,
		Database:       mongoDatabase,
		ConnectTimeout: 10 * time.Second,
		MaxPoolSize:    100,
	})
	if err != nil {
		slog.Error("Failed to connect to MongoDB", "error", err)
		os.Exit(1)
	}
	defer mongoClient.Close(context.Background())

	slog.Info("Connected to MongoDB", "database", mongoDatabase)

	// Initialize NATS client
	natsClient, err := nats.NewClient(context.Background(), nats.Config{
		URL:  natsURL,
		Name: serviceName,
	})
	if err != nil {
		slog.Error("Failed to connect to NATS", "error", err)
		os.Exit(1)
	}
	defer natsClient.Close()

	slog.Info("Connected to NATS")

	// Initialize event publisher
	eventPublisher := nats.NewEventPublisher(natsClient)

	// Initialize repositories
	userRepo := mongodb.NewUserRepository(mongoClient.Database())

	// Ensure indexes
	if err := userRepo.EnsureIndexes(context.Background()); err != nil {
		slog.Warn("Failed to ensure indexes", "error", err)
	}

	// Initialize GraphQL resolver
	graphqlResolver := resolver.NewResolver(userRepo, eventPublisher)

	// Create HTTP server
	mux := http.NewServeMux()

	// Health check endpoint
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"healthy"}`))
	})

	// Ready check endpoint
	mux.HandleFunc("/ready", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
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
	// After running gqlgen generate, this will be replaced with the actual handler
	mux.HandleFunc("/query", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"message":"GraphQL endpoint ready. Run 'go generate ./...' to generate schema."}`))
	})

	// GraphQL Playground (development only)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`<!DOCTYPE html>
<html>
<head>
    <title>Identity Service - GraphQL Playground</title>
    <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/graphql-playground-react/build/static/css/index.css" />
    <script src="https://cdn.jsdelivr.net/npm/graphql-playground-react/build/static/js/middleware.js"></script>
</head>
<body>
    <div id="root">
        <style>
            body { margin: 0; }
        </style>
        <script>
            window.addEventListener('load', function (event) {
                GraphQLPlayground.init(document.getElementById('root'), { endpoint: '/query' })
            })
        </script>
    </div>
</body>
</html>`))
	})

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
		slog.Info("Server starting", "port", port, "playground", "/")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Server failed", "error", err)
			os.Exit(1)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("Shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("Server shutdown failed", "error", err)
		os.Exit(1)
	}

	slog.Info("Server stopped")

	// Suppress unused variable warning
	_ = graphqlResolver
}
