// Package main is the entry point for the Partner service.
package main

import (
	"context"
	"fmt"
	"log/slog"
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
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/yousoon/services/partner/internal/application/commands"
	"github.com/yousoon/services/partner/internal/application/queries"
	"github.com/yousoon/services/partner/internal/config"
	"github.com/yousoon/services/partner/internal/infrastructure"
	"github.com/yousoon/services/partner/internal/infrastructure/mongodb"
	"github.com/yousoon/services/partner/internal/interface/graphql/resolver"
	// "github.com/yousoon/services/partner/internal/interface/graphql/generated"
)

const (
	serviceName    = "partner-service"
	serviceVersion = "1.0.0"
)

func main() {
	// Initialize structured logger
	logger := initLogger()
	slog.SetDefault(logger)

	logger.Info("Starting service",
		slog.String("service", serviceName),
		slog.String("version", serviceVersion),
	)

	// Load configuration
	cfg := config.Load()

	// Create context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initialize MongoDB
	mongoClient, err := initMongoDB(ctx, cfg.MongoDB)
	if err != nil {
		logger.Error("Failed to connect to MongoDB", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer func() {
		if err := mongoClient.Disconnect(ctx); err != nil {
			logger.Error("Failed to disconnect from MongoDB", slog.String("error", err.Error()))
		}
	}()

	// Initialize repositories
	db := mongoClient.Database(cfg.MongoDB.Database)
	partnerRepo := mongodb.NewPartnerRepository(db)
	partnerReadRepo := mongodb.NewPartnerReadRepository(db)

	// Initialize event publisher (NoOp for now, will be replaced with NATS)
	eventPublisher := infrastructure.NewNoOpEventPublisher()

	// Initialize command handlers
	registerPartnerHandler := commands.NewRegisterPartnerHandler(partnerRepo, eventPublisher)
	updatePartnerHandler := commands.NewUpdatePartnerHandler(partnerRepo, eventPublisher)
	verifyPartnerHandler := commands.NewVerifyPartnerHandler(partnerRepo, eventPublisher)
	suspendPartnerHandler := commands.NewSuspendPartnerHandler(partnerRepo, eventPublisher)
	addEstablishmentHandler := commands.NewAddEstablishmentHandler(partnerRepo, eventPublisher)
	inviteTeamMemberHandler := commands.NewInviteTeamMemberHandler(partnerRepo, eventPublisher)
	acceptTeamInvitationHandler := commands.NewAcceptTeamInvitationHandler(partnerRepo, eventPublisher)

	// Initialize query handlers
	getPartnerHandler := queries.NewGetPartnerHandler(partnerRepo)
	getPartnerByOwnerHandler := queries.NewGetPartnerByOwnerHandler(partnerRepo)
	listPartnersHandler := queries.NewListPartnersHandler(partnerRepo)
	getEstablishmentHandler := queries.NewGetEstablishmentHandler(partnerRepo)
	listEstablishmentsHandler := queries.NewListEstablishmentsHandler(partnerRepo)
	getTeamMembersHandler := queries.NewGetTeamMembersHandler(partnerRepo)
	getPartnersForTeamMemberHandler := queries.NewGetPartnersForTeamMemberHandler(partnerRepo)

	// Create resolver with all dependencies
	resolverConfig := &resolver.Resolver{
		PartnerRepo:                     partnerRepo,
		PartnerReadRepo:                 partnerReadRepo,
		RegisterPartnerHandler:          registerPartnerHandler,
		UpdatePartnerHandler:            updatePartnerHandler,
		VerifyPartnerHandler:            verifyPartnerHandler,
		SuspendPartnerHandler:           suspendPartnerHandler,
		AddEstablishmentHandler:         addEstablishmentHandler,
		InviteTeamMemberHandler:         inviteTeamMemberHandler,
		AcceptTeamInvitationHandler:     acceptTeamInvitationHandler,
		GetPartnerHandler:               getPartnerHandler,
		GetPartnerByOwnerHandler:        getPartnerByOwnerHandler,
		ListPartnersHandler:             listPartnersHandler,
		GetEstablishmentHandler:         getEstablishmentHandler,
		ListEstablishmentsHandler:       listEstablishmentsHandler,
		GetTeamMembersHandler:           getTeamMembersHandler,
		GetPartnersForTeamMemberHandler: getPartnersForTeamMemberHandler,
	}

	// Create GraphQL server
	// NOTE: Replace with generated.NewExecutableSchema when gqlgen is run
	// srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{
	//     Resolvers: resolverConfig,
	// }))

	// Placeholder handler until schema is generated
	srv := createGraphQLHandler(resolverConfig, cfg.GraphQL)

	// Setup HTTP routes
	mux := http.NewServeMux()

	// GraphQL endpoint
	mux.Handle("/graphql", corsMiddleware(srv))
	mux.Handle("/query", corsMiddleware(srv)) // Alias for compatibility

	// Playground (development only)
	if cfg.GraphQL.PlaygroundEnabled {
		mux.Handle("/", playground.Handler("Partner Service GraphQL", "/graphql"))
		logger.Info("GraphQL Playground enabled at /")
	}

	// Health checks
	mux.HandleFunc("/health", healthHandler)
	mux.HandleFunc("/ready", readyHandler(mongoClient))

	// Metrics endpoint
	if cfg.Observability.MetricsEnabled {
		go startMetricsServer(cfg.Observability.MetricsPort, cfg.Observability.MetricsPath, logger)
	}

	// Create HTTP server
	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port),
		Handler:      loggingMiddleware(mux, logger),
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	// Start server in goroutine
	go func() {
		logger.Info("GraphQL server starting",
			slog.String("address", server.Addr),
		)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("Server failed", slog.String("error", err.Error()))
			cancel()
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	// Graceful shutdown with timeout
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Error("Server forced to shutdown", slog.String("error", err.Error()))
	}

	logger.Info("Server stopped")
}

// initLogger initializes the structured logger.
func initLogger() *slog.Logger {
	var handler slog.Handler

	logFormat := os.Getenv("LOG_FORMAT")
	if logFormat == "text" {
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: getLogLevel(),
		})
	} else {
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: getLogLevel(),
		})
	}

	return slog.New(handler)
}

func getLogLevel() slog.Level {
	switch os.Getenv("LOG_LEVEL") {
	case "debug":
		return slog.LevelDebug
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

// initMongoDB initializes the MongoDB connection.
func initMongoDB(ctx context.Context, cfg config.MongoDBConfig) (*mongo.Client, error) {
	clientOpts := options.Client().
		ApplyURI(cfg.URI).
		SetMaxPoolSize(cfg.MaxPoolSize).
		SetMinPoolSize(cfg.MinPoolSize).
		SetMaxConnIdleTime(cfg.MaxConnIdleTime)

	connectCtx, cancel := context.WithTimeout(ctx, cfg.ConnectTimeout)
	defer cancel()

	client, err := mongo.Connect(connectCtx, clientOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// Ping to verify connection
	if err := client.Ping(connectCtx, nil); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	slog.Info("Connected to MongoDB",
		slog.String("database", cfg.Database),
	)

	return client, nil
}

// createGraphQLHandler creates the GraphQL handler with all transports and extensions.
func createGraphQLHandler(resolverConfig *resolver.Resolver, cfg config.GraphQLConfig) *handler.Server {
	// NOTE: This is a placeholder. Replace with actual generated schema:
	// srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{
	//     Resolvers: resolverConfig,
	// }))

	// Temporary: create empty handler for compilation
	// In production, uncomment above and remove this placeholder
	srv := handler.New(nil)

	// Add transports
	srv.AddTransport(transport.Websocket{
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true // TODO: Restrict in production
			},
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
		KeepAlivePingInterval: 10 * time.Second,
	})
	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})
	srv.AddTransport(transport.MultipartForm{})

	// Add extensions
	if cfg.IntrospectionEnabled {
		srv.Use(extension.Introspection{})
	}

	// Query caching
	srv.SetQueryCache(lru.New(1000))

	// Automatic persisted queries
	if cfg.PersistedQueriesEnabled {
		srv.Use(extension.AutomaticPersistedQuery{
			Cache: lru.New(100),
		})
	}

	// Complexity limit
	srv.Use(extension.FixedComplexityLimit(cfg.QueryComplexityLimit))

	return srv
}

// healthHandler returns service health status.
func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"status":"healthy","service":"%s","version":"%s"}`, serviceName, serviceVersion)
}

// readyHandler returns service readiness status.
func readyHandler(mongoClient *mongo.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		// Check MongoDB connection
		if err := mongoClient.Ping(ctx, nil); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusServiceUnavailable)
			fmt.Fprintf(w, `{"status":"not_ready","error":"mongodb: %s"}`, err.Error())
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"status":"ready"}`)
	}
}

// startMetricsServer starts a separate server for Prometheus metrics.
func startMetricsServer(port int, path string, logger *slog.Logger) {
	mux := http.NewServeMux()
	mux.Handle(path, promhttp.Handler())

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}

	logger.Info("Metrics server starting",
		slog.Int("port", port),
		slog.String("path", path),
	)

	if err := server.ListenAndServe(); err != nil {
		logger.Error("Metrics server failed", slog.String("error", err.Error()))
	}
}

// Middleware

// loggingMiddleware logs HTTP requests.
func loggingMiddleware(next http.Handler, logger *slog.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Skip logging for health checks
		if r.URL.Path == "/health" || r.URL.Path == "/ready" {
			next.ServeHTTP(w, r)
			return
		}

		next.ServeHTTP(w, r)

		logger.Info("HTTP request",
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.String("remote_addr", r.RemoteAddr),
			slog.Duration("duration", time.Since(start)),
		)
	})
}

// corsMiddleware adds CORS headers.
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, X-CSRF-Token")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Max-Age", "300")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
