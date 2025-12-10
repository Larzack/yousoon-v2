// Package metrics provides Prometheus metrics utilities.
package metrics

import (
	"fmt"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Metrics holds common metrics for a service.
type Metrics struct {
	// HTTP/gRPC request metrics
	RequestsTotal    *prometheus.CounterVec
	RequestDuration  *prometheus.HistogramVec
	RequestsInFlight prometheus.Gauge

	// Database metrics
	DBQueriesTotal   *prometheus.CounterVec
	DBQueryDuration  *prometheus.HistogramVec
	DBConnectionPool *prometheus.GaugeVec

	// Cache metrics
	CacheHits   *prometheus.CounterVec
	CacheMisses *prometheus.CounterVec

	// Event metrics
	EventsPublished *prometheus.CounterVec
	EventsReceived  *prometheus.CounterVec
	EventsProcessed *prometheus.CounterVec
	EventsFailed    *prometheus.CounterVec

	// Business metrics
	UsersActive   prometheus.Gauge
	BookingsTotal prometheus.Counter
	OffersActive  prometheus.Gauge
}

// NewMetrics creates a new metrics instance for a service.
func NewMetrics(namespace, subsystem string) *Metrics {
	return &Metrics{
		RequestsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "requests_total",
				Help:      "Total number of requests",
			},
			[]string{"method", "path", "status"},
		),
		RequestDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "request_duration_seconds",
				Help:      "Request duration in seconds",
				Buckets:   []float64{.001, .005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10},
			},
			[]string{"method", "path"},
		),
		RequestsInFlight: promauto.NewGauge(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "requests_in_flight",
				Help:      "Number of requests currently being processed",
			},
		),
		DBQueriesTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "db_queries_total",
				Help:      "Total number of database queries",
			},
			[]string{"operation", "collection", "status"},
		),
		DBQueryDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "db_query_duration_seconds",
				Help:      "Database query duration in seconds",
				Buckets:   []float64{.001, .005, .01, .025, .05, .1, .25, .5, 1},
			},
			[]string{"operation", "collection"},
		),
		DBConnectionPool: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "db_connection_pool",
				Help:      "Database connection pool metrics",
			},
			[]string{"state"}, // active, idle, total
		),
		CacheHits: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "cache_hits_total",
				Help:      "Total number of cache hits",
			},
			[]string{"cache"},
		),
		CacheMisses: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "cache_misses_total",
				Help:      "Total number of cache misses",
			},
			[]string{"cache"},
		),
		EventsPublished: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "events_published_total",
				Help:      "Total number of events published",
			},
			[]string{"event_type"},
		),
		EventsReceived: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "events_received_total",
				Help:      "Total number of events received",
			},
			[]string{"event_type"},
		),
		EventsProcessed: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "events_processed_total",
				Help:      "Total number of events processed successfully",
			},
			[]string{"event_type"},
		),
		EventsFailed: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "events_failed_total",
				Help:      "Total number of events that failed processing",
			},
			[]string{"event_type", "error"},
		),
	}
}

// RecordRequest records a request metric.
func (m *Metrics) RecordRequest(method, path string, statusCode int, duration time.Duration) {
	m.RequestsTotal.WithLabelValues(method, path, fmt.Sprintf("%d", statusCode)).Inc()
	m.RequestDuration.WithLabelValues(method, path).Observe(duration.Seconds())
}

// RecordDBQuery records a database query metric.
func (m *Metrics) RecordDBQuery(operation, collection string, err error, duration time.Duration) {
	status := "success"
	if err != nil {
		status = "error"
	}
	m.DBQueriesTotal.WithLabelValues(operation, collection, status).Inc()
	m.DBQueryDuration.WithLabelValues(operation, collection).Observe(duration.Seconds())
}

// RecordCacheHit records a cache hit.
func (m *Metrics) RecordCacheHit(cache string) {
	m.CacheHits.WithLabelValues(cache).Inc()
}

// RecordCacheMiss records a cache miss.
func (m *Metrics) RecordCacheMiss(cache string) {
	m.CacheMisses.WithLabelValues(cache).Inc()
}

// RecordEventPublished records an event publication.
func (m *Metrics) RecordEventPublished(eventType string) {
	m.EventsPublished.WithLabelValues(eventType).Inc()
}

// RecordEventReceived records an event reception.
func (m *Metrics) RecordEventReceived(eventType string) {
	m.EventsReceived.WithLabelValues(eventType).Inc()
}

// RecordEventProcessed records a successful event processing.
func (m *Metrics) RecordEventProcessed(eventType string) {
	m.EventsProcessed.WithLabelValues(eventType).Inc()
}

// RecordEventFailed records a failed event processing.
func (m *Metrics) RecordEventFailed(eventType, errorType string) {
	m.EventsFailed.WithLabelValues(eventType, errorType).Inc()
}

// Server provides an HTTP server for metrics.
type Server struct {
	server *http.Server
}

// NewServer creates a new metrics server.
func NewServer(port int) *Server {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	return &Server{
		server: &http.Server{
			Addr:    fmt.Sprintf(":%d", port),
			Handler: mux,
		},
	}
}

// Start starts the metrics server.
func (s *Server) Start() error {
	return s.server.ListenAndServe()
}

// Shutdown gracefully shuts down the metrics server.
func (s *Server) Shutdown() error {
	return s.server.Close()
}
