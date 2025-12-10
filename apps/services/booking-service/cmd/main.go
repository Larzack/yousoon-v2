package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/yousoon/apps/services/booking-service/config"
	"github.com/yousoon/apps/services/booking-service/internal/application/commands"
	"github.com/yousoon/apps/services/booking-service/internal/application/queries"
	"github.com/yousoon/apps/services/booking-service/internal/domain"
	"github.com/yousoon/apps/services/booking-service/internal/infrastructure/mongodb"
	"github.com/yousoon/apps/services/booking-service/internal/interface/graphql/resolver"
)

func main() {
	// Load configuration
	cfg := config.Load()

	log.Printf("Starting %s in %s mode", cfg.ServiceName, cfg.Environment)

	// Context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Connect to MongoDB
	mongoClient, err := connectMongoDB(ctx, cfg.MongoURI)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer mongoClient.Disconnect(ctx)

	db := mongoClient.Database(cfg.MongoDatabase)

	// Initialize repositories
	outingRepo := mongodb.NewOutingRepository(db)

	// Initialize services (stubs for now - would be gRPC clients)
	offerService := &stubOfferService{}
	userService := &stubUserService{}
	notifyService := &stubNotificationService{}

	// Initialize command handlers
	bookOutingHandler := commands.NewBookOutingHandler(
		outingRepo,
		offerService,
		userService,
		notifyService,
		cfg.BookingExpirationMinutes,
	)
	checkInHandler := commands.NewCheckInOutingHandler(outingRepo, notifyService)
	cancelOutingHandler := commands.NewCancelOutingHandler(outingRepo, offerService, notifyService)

	// Initialize query handlers
	getOutingHandler := queries.NewGetOutingHandler(outingRepo)
	getOutingByQRHandler := queries.NewGetOutingByQRHandler(outingRepo)
	listUserOutingsHandler := queries.NewListUserOutingsHandler(outingRepo)
	listPartnerOutingsHandler := queries.NewListPartnerOutingsHandler(outingRepo)
	listEstablishmentOutingsHandler := queries.NewListEstablishmentOutingsHandler(outingRepo)
	getBookingStatsHandler := queries.NewGetBookingStatsHandler(outingRepo)

	// Initialize resolver
	resolv := resolver.NewResolver(
		bookOutingHandler,
		checkInHandler,
		cancelOutingHandler,
		getOutingHandler,
		getOutingByQRHandler,
		listUserOutingsHandler,
		listPartnerOutingsHandler,
		listEstablishmentOutingsHandler,
		getBookingStatsHandler,
	)

	// Create GraphQL server
	srv := handler.NewDefaultServer(nil) // Would use generated.NewExecutableSchema

	// HTTP server
	mux := http.NewServeMux()

	// GraphQL endpoint
	mux.Handle(cfg.GraphQLPath, srv)

	// Playground (development only)
	if cfg.PlaygroundEnabled {
		mux.Handle("/", playground.Handler("Booking Service", cfg.GraphQLPath))
	}

	// Health check
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"healthy"}`))
	})

	// Ready check
	mux.HandleFunc("/ready", func(w http.ResponseWriter, r *http.Request) {
		// Check MongoDB connection
		if err := mongoClient.Ping(ctx, nil); err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte(`{"status":"not ready"}`))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ready"}`))
	})

	server := &http.Server{
		Addr:         cfg.GetServerAddr(),
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server
	go func() {
		log.Printf("Server listening on %s", cfg.GetServerAddr())
		log.Printf("GraphQL endpoint: %s", cfg.GraphQLPath)
		if cfg.PlaygroundEnabled {
			log.Printf("Playground: http://%s/", cfg.GetServerAddr())
		}
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited properly")

	// Keep resolver reference
	_ = resolv
}

func connectMongoDB(ctx context.Context, uri string) (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI(uri)

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	log.Println("Connected to MongoDB")
	return client, nil
}

// =============================================================================
// STUB SERVICES (would be replaced with gRPC clients)
// =============================================================================

type stubOfferService struct{}

func (s *stubOfferService) GetOfferSnapshot(ctx context.Context, offerID string) (*domain.OfferSnapshot, error) {
	// Would call Discovery service via gRPC
	snapshot := domain.NewOfferSnapshot(
		offerID, "partner-1", "establishment-1",
		"Sample Offer", "Description",
		"percentage", 20,
		"restaurant",
		"Restaurant Name", "123 Main St",
		48.8566, 2.3522,
		"https://example.com/image.jpg",
	)
	return &snapshot, nil
}

func (s *stubOfferService) CanBook(ctx context.Context, offerID string) error {
	return nil
}

func (s *stubOfferService) IncrementBookingCount(ctx context.Context, offerID string) error {
	return nil
}

func (s *stubOfferService) DecrementBookingCount(ctx context.Context, offerID string) error {
	return nil
}

type stubUserService struct{}

func (s *stubUserService) GetUserSnapshot(ctx context.Context, userID string) (*domain.UserSnapshot, error) {
	snapshot := domain.NewUserSnapshot(userID, "John", "Doe", "john@example.com")
	return &snapshot, nil
}

func (s *stubUserService) CanBook(ctx context.Context, userID string) error {
	return nil
}

type stubNotificationService struct{}

func (s *stubNotificationService) SendBookingConfirmation(ctx context.Context, outing *domain.Outing) error {
	return nil
}

func (s *stubNotificationService) SendCheckInConfirmation(ctx context.Context, outing *domain.Outing) error {
	return nil
}

func (s *stubNotificationService) SendCancellationNotification(ctx context.Context, outing *domain.Outing) error {
	return nil
}

func (s *stubNotificationService) SendExpirationReminder(ctx context.Context, outing *domain.Outing) error {
	return nil
}
