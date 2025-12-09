// Package main is the entry point for the Discovery service.
package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/websocket"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"

	"github.com/yousoon/discovery-service/internal/config"
	esrepo "github.com/yousoon/discovery-service/internal/infrastructure/elasticsearch"
	mongorepo "github.com/yousoon/discovery-service/internal/infrastructure/mongodb"
	"github.com/yousoon/discovery-service/internal/interface/graphql/generated"
	"github.com/yousoon/discovery-service/internal/interface/graphql/resolver"
	"github.com/yousoon/shared/infrastructure/nats"
	"github.com/yousoon/shared/observability/logger"
	"github.com/yousoon/shared/observability/tracing"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Failed to load configuration: %v\n", err)
		os.Exit(1)
	}

	// Initialize logger
	log, err := logger.New(cfg.LogLevel, cfg.LogFormat, cfg.ServiceName)
	if err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer log.Sync()

	log.Info("Starting Discovery Service",
		zap.String("version", cfg.Version),
		zap.String("environment", cfg.Environment),
	)

	// Initialize tracing
	tp, err := tracing.InitTracer(cfg.ServiceName, cfg.JaegerEndpoint)
	if err != nil {
		log.Error("Failed to initialize tracer", zap.Error(err))
	}
	defer func() {
		if tp != nil {
			_ = tp.Shutdown(context.Background())
		}
	}()

	// Create context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Connect to MongoDB
	mongoClient, err := connectMongoDB(ctx, cfg, log)
	if err != nil {
		log.Fatal("Failed to connect to MongoDB", zap.Error(err))
	}
	defer func() {
		if err := mongoClient.Disconnect(ctx); err != nil {
			log.Error("Failed to disconnect from MongoDB", zap.Error(err))
		}
	}()

	mongoDB := mongoClient.Database(cfg.MongoDB)

	// Connect to Elasticsearch
	esClient, err := connectElasticsearch(cfg, log)
	if err != nil {
		log.Error("Failed to connect to Elasticsearch", zap.Error(err))
		// Continue without Elasticsearch - fallback to MongoDB search
	}

	// Connect to NATS
	natsClient, err := nats.NewClient(cfg.NATSURL, cfg.NATSClusterID, cfg.ServiceName)
	if err != nil {
		log.Fatal("Failed to connect to NATS", zap.Error(err))
	}
	defer natsClient.Close()

	eventPublisher := nats.NewPublisher(natsClient)

	// Initialize repositories
	offerRepo := mongorepo.NewOfferRepository(mongoDB)
	categoryRepo := mongorepo.NewCategoryRepository(mongoDB)

	// Ensure MongoDB indexes
	if err := offerRepo.EnsureIndexes(ctx); err != nil {
		log.Error("Failed to create offer indexes", zap.Error(err))
	}
	if err := categoryRepo.EnsureIndexes(ctx); err != nil {
		log.Error("Failed to create category indexes", zap.Error(err))
	}

	// Initialize Elasticsearch repository
	var searchRepo *esrepo.OfferSearchRepository
	if esClient != nil {
		searchRepo = esrepo.NewOfferSearchRepository(esClient)
		if err := searchRepo.EnsureIndex(ctx); err != nil {
			log.Error("Failed to create Elasticsearch index", zap.Error(err))
		}
	}

	// Create resolver
	res := resolver.NewResolver(offerRepo, categoryRepo, searchRepo, eventPublisher)

	// Setup GraphQL server
	srv := handler.New(generated.NewExecutableSchema(generated.Config{
		Resolvers: res,
	}))

	// Add transports
	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})
	srv.AddTransport(transport.MultipartForm{})
	srv.AddTransport(&transport.Websocket{
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true // Allow all origins in development
			},
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
		KeepAlivePingInterval: 10 * time.Second,
	})

	// Add extensions
	if cfg.IsDevelopment() {
		srv.Use(extension.Introspection{})
	}
	srv.Use(extension.FixedComplexityLimit(100))

	// Setup router
	router := chi.NewRouter()

	// Middleware
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(30 * time.Second))

	// Routes
	router.Handle("/graphql", srv)
	if cfg.IsDevelopment() {
		router.Handle("/", playground.Handler("Discovery Service", "/graphql"))
	}

	// Health check
	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"healthy"}`))
	})

	// Readiness check
	router.Get("/ready", func(w http.ResponseWriter, r *http.Request) {
		// Check MongoDB connection
		if err := mongoClient.Ping(r.Context(), nil); err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte(`{"status":"not ready","reason":"mongodb unavailable"}`))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ready"}`))
	})

	// Start metrics server
	go func() {
		metricsRouter := chi.NewRouter()
		metricsRouter.Handle("/metrics", promhttp.Handler())
		metricsAddr := fmt.Sprintf(":%s", cfg.MetricsPort)
		log.Info("Starting metrics server", zap.String("addr", metricsAddr))
		if err := http.ListenAndServe(metricsAddr, metricsRouter); err != nil {
			log.Error("Metrics server error", zap.Error(err))
		}
	}()

	// Start HTTP server
	httpAddr := fmt.Sprintf(":%s", cfg.HTTPPort)
	httpServer := &http.Server{
		Addr:         httpAddr,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Info("Starting HTTP server", zap.String("addr", httpAddr))
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("HTTP server error", zap.Error(err))
		}
	}()

	// Register with Schema Registry
	if cfg.SchemaRegistryURL != "" {
		go registerWithSchemaRegistry(cfg, log)
	}

	// Wait for shutdown signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down Discovery Service...")

	// Graceful shutdown
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		log.Error("HTTP server shutdown error", zap.Error(err))
	}

	log.Info("Discovery Service stopped")
}

// connectMongoDB establishes a connection to MongoDB.
func connectMongoDB(ctx context.Context, cfg *config.Config, log *zap.Logger) (*mongo.Client, error) {
	log.Info("Connecting to MongoDB", zap.String("uri", maskURI(cfg.MongoURI)))

	clientOptions := options.Client().
		ApplyURI(cfg.MongoURI).
		SetMaxPoolSize(cfg.MongoMaxPoolSize).
		SetServerSelectionTimeout(cfg.MongoTimeout)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// Ping to verify connection
	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	log.Info("Connected to MongoDB successfully")
	return client, nil
}

// connectElasticsearch establishes a connection to Elasticsearch.
func connectElasticsearch(cfg *config.Config, log *zap.Logger) (*elasticsearch.Client, error) {
	log.Info("Connecting to Elasticsearch", zap.Strings("urls", cfg.ElasticsearchURLs))

	esCfg := elasticsearch.Config{
		Addresses: cfg.ElasticsearchURLs,
	}

	if cfg.ElasticsearchUsername != "" {
		esCfg.Username = cfg.ElasticsearchUsername
		esCfg.Password = cfg.ElasticsearchPassword
	}

	client, err := elasticsearch.NewClient(esCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create Elasticsearch client: %w", err)
	}

	// Test connection
	res, err := client.Info()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Elasticsearch: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("Elasticsearch error: %s", res.String())
	}

	log.Info("Connected to Elasticsearch successfully")
	return client, nil
}

// registerWithSchemaRegistry registers the service schema with the registry.
func registerWithSchemaRegistry(cfg *config.Config, log *zap.Logger) {
	// This would typically read the schema file and POST it to the registry
	log.Info("Registering with Schema Registry", zap.String("url", cfg.SchemaRegistryURL))

	// Implementation would be similar to partner-service
	// For now, just log
	log.Info("Schema registration complete")
}

// maskURI masks sensitive parts of a URI for logging.
func maskURI(uri string) string {
	// Simple masking - in production use a proper URL parser
	if len(uri) > 20 {
		return uri[:20] + "..."
	}
	return uri
}
