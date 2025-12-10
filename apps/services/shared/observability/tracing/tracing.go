// Package tracing provides OpenTelemetry tracing utilities.
package tracing

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.opentelemetry.io/otel/trace"
)

// Config holds tracing configuration.
type Config struct {
	ServiceName    string
	ServiceVersion string
	Environment    string
	OTLPEndpoint   string
	Enabled        bool
	SampleRate     float64
}

// DefaultConfig returns a default tracing configuration.
func DefaultConfig(serviceName, version string) Config {
	return Config{
		ServiceName:    serviceName,
		ServiceVersion: version,
		Environment:    "development",
		OTLPEndpoint:   "localhost:4317",
		Enabled:        true,
		SampleRate:     1.0,
	}
}

// Tracer wraps OpenTelemetry tracer.
type Tracer struct {
	tracer   trace.Tracer
	provider *sdktrace.TracerProvider
	config   Config
}

// Init initializes the tracing system.
func Init(ctx context.Context, cfg Config) (*Tracer, error) {
	if !cfg.Enabled {
		return &Tracer{
			tracer: otel.Tracer(cfg.ServiceName),
			config: cfg,
		}, nil
	}

	// Create OTLP exporter
	client := otlptracegrpc.NewClient(
		otlptracegrpc.WithEndpoint(cfg.OTLPEndpoint),
		otlptracegrpc.WithInsecure(),
	)

	exporter, err := otlptrace.New(ctx, client)
	if err != nil {
		return nil, fmt.Errorf("failed to create trace exporter: %w", err)
	}

	// Create resource
	res, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(cfg.ServiceName),
			semconv.ServiceVersion(cfg.ServiceVersion),
			attribute.String("environment", cfg.Environment),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	// Create trace provider
	provider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sdktrace.TraceIDRatioBased(cfg.SampleRate)),
	)

	// Set global provider
	otel.SetTracerProvider(provider)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	return &Tracer{
		tracer:   provider.Tracer(cfg.ServiceName),
		provider: provider,
		config:   cfg,
	}, nil
}

// Shutdown shuts down the tracer provider.
func (t *Tracer) Shutdown(ctx context.Context) error {
	if t.provider != nil {
		return t.provider.Shutdown(ctx)
	}
	return nil
}

// Start starts a new span.
func (t *Tracer) Start(ctx context.Context, name string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	return t.tracer.Start(ctx, name, opts...)
}

// StartWithAttributes starts a new span with attributes.
func (t *Tracer) StartWithAttributes(ctx context.Context, name string, attrs ...attribute.KeyValue) (context.Context, trace.Span) {
	return t.tracer.Start(ctx, name, trace.WithAttributes(attrs...))
}

// SpanFromContext returns the span from context.
func SpanFromContext(ctx context.Context) trace.Span {
	return trace.SpanFromContext(ctx)
}

// AddEvent adds an event to the current span.
func AddEvent(ctx context.Context, name string, attrs ...attribute.KeyValue) {
	span := trace.SpanFromContext(ctx)
	span.AddEvent(name, trace.WithAttributes(attrs...))
}

// SetStatus sets the status of the current span.
func SetStatus(ctx context.Context, code codes.Code, description string) {
	span := trace.SpanFromContext(ctx)
	span.SetStatus(code, description)
}

// RecordError records an error on the current span.
func RecordError(ctx context.Context, err error) {
	span := trace.SpanFromContext(ctx)
	span.RecordError(err)
	span.SetStatus(codes.Error, err.Error())
}

// SetAttributes sets attributes on the current span.
func SetAttributes(ctx context.Context, attrs ...attribute.KeyValue) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attrs...)
}

// TraceID returns the trace ID from context.
func TraceID(ctx context.Context) string {
	span := trace.SpanFromContext(ctx)
	return span.SpanContext().TraceID().String()
}

// SpanID returns the span ID from context.
func SpanID(ctx context.Context) string {
	span := trace.SpanFromContext(ctx)
	return span.SpanContext().SpanID().String()
}

// Common attribute helpers

// UserIDAttr creates a user ID attribute.
func UserIDAttr(userID string) attribute.KeyValue {
	return attribute.String("user.id", userID)
}

// PartnerIDAttr creates a partner ID attribute.
func PartnerIDAttr(partnerID string) attribute.KeyValue {
	return attribute.String("partner.id", partnerID)
}

// OfferIDAttr creates an offer ID attribute.
func OfferIDAttr(offerID string) attribute.KeyValue {
	return attribute.String("offer.id", offerID)
}

// DBOperationAttr creates a database operation attribute.
func DBOperationAttr(operation string) attribute.KeyValue {
	return attribute.String("db.operation", operation)
}

// DBCollectionAttr creates a database collection attribute.
func DBCollectionAttr(collection string) attribute.KeyValue {
	return attribute.String("db.collection", collection)
}

// EventTypeAttr creates an event type attribute.
func EventTypeAttr(eventType string) attribute.KeyValue {
	return attribute.String("event.type", eventType)
}
