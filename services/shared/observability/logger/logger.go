// Package logger provides a structured logging interface.
package logger

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"runtime"
	"time"
)

// Level represents the log level.
type Level = slog.Level

const (
	LevelDebug = slog.LevelDebug
	LevelInfo  = slog.LevelInfo
	LevelWarn  = slog.LevelWarn
	LevelError = slog.LevelError
)

// Logger is a structured logger interface.
type Logger interface {
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
	With(args ...any) Logger
	WithContext(ctx context.Context) Logger
}

// logger wraps slog.Logger.
type logger struct {
	slog *slog.Logger
}

// Options configures the logger.
type Options struct {
	Level       Level
	AddSource   bool
	JSON        bool
	ServiceName string
	Version     string
	Environment string
	Writer      io.Writer
}

// DefaultOptions returns default logger options.
func DefaultOptions() Options {
	return Options{
		Level:       LevelInfo,
		AddSource:   false,
		JSON:        false,
		ServiceName: "yousoon",
		Version:     "1.0.0",
		Environment: "development",
		Writer:      os.Stdout,
	}
}

// New creates a new logger with the given options.
func New(opts Options) Logger {
	var handler slog.Handler

	handlerOpts := &slog.HandlerOptions{
		Level:     opts.Level,
		AddSource: opts.AddSource,
	}

	if opts.Writer == nil {
		opts.Writer = os.Stdout
	}

	if opts.JSON {
		handler = slog.NewJSONHandler(opts.Writer, handlerOpts)
	} else {
		handler = slog.NewTextHandler(opts.Writer, handlerOpts)
	}

	slogger := slog.New(handler).With(
		slog.String("service", opts.ServiceName),
		slog.String("version", opts.Version),
		slog.String("environment", opts.Environment),
	)

	return &logger{slog: slogger}
}

// NewFromEnv creates a logger from environment configuration.
func NewFromEnv(serviceName, version string) Logger {
	env := os.Getenv("ENVIRONMENT")
	if env == "" {
		env = "development"
	}

	levelStr := os.Getenv("LOG_LEVEL")
	level := parseLevel(levelStr)

	return New(Options{
		Level:       level,
		AddSource:   env != "development",
		JSON:        env != "development",
		ServiceName: serviceName,
		Version:     version,
		Environment: env,
	})
}

func parseLevel(s string) Level {
	switch s {
	case "debug", "DEBUG":
		return LevelDebug
	case "warn", "WARN", "warning", "WARNING":
		return LevelWarn
	case "error", "ERROR":
		return LevelError
	default:
		return LevelInfo
	}
}

// Debug logs a debug message.
func (l *logger) Debug(msg string, args ...any) {
	l.slog.Debug(msg, args...)
}

// Info logs an info message.
func (l *logger) Info(msg string, args ...any) {
	l.slog.Info(msg, args...)
}

// Warn logs a warning message.
func (l *logger) Warn(msg string, args ...any) {
	l.slog.Warn(msg, args...)
}

// Error logs an error message.
func (l *logger) Error(msg string, args ...any) {
	l.slog.Error(msg, args...)
}

// With returns a new logger with the given attributes.
func (l *logger) With(args ...any) Logger {
	return &logger{slog: l.slog.With(args...)}
}

// WithContext returns a new logger with context values.
func (l *logger) WithContext(ctx context.Context) Logger {
	// Extract trace ID from context if present
	if traceID := ctx.Value(traceIDKey{}); traceID != nil {
		return l.With(slog.String("trace_id", fmt.Sprintf("%v", traceID)))
	}
	return l
}

type traceIDKey struct{}

// SetTraceID adds a trace ID to the context.
func SetTraceID(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, traceIDKey{}, traceID)
}

// GetTraceID extracts the trace ID from context.
func GetTraceID(ctx context.Context) string {
	if traceID := ctx.Value(traceIDKey{}); traceID != nil {
		return fmt.Sprintf("%v", traceID)
	}
	return ""
}

// LogRequest logs an HTTP/gRPC request.
func (l *logger) LogRequest(method, path string, statusCode int, duration time.Duration, err error) {
	if err != nil {
		l.Error("request failed",
			slog.String("method", method),
			slog.String("path", path),
			slog.Int("status_code", statusCode),
			slog.Duration("duration", duration),
			slog.String("error", err.Error()),
		)
	} else {
		l.Info("request completed",
			slog.String("method", method),
			slog.String("path", path),
			slog.Int("status_code", statusCode),
			slog.Duration("duration", duration),
		)
	}
}

// LogPanic logs a panic with stack trace.
func LogPanic(l Logger, recovered interface{}) {
	buf := make([]byte, 4096)
	n := runtime.Stack(buf, false)
	l.Error("panic recovered",
		slog.Any("panic", recovered),
		slog.String("stack", string(buf[:n])),
	)
}

// Global logger instance
var global Logger = New(DefaultOptions())

// SetGlobal sets the global logger.
func SetGlobal(l Logger) {
	global = l
}

// Global returns the global logger.
func Global() Logger {
	return global
}

// Debug logs a debug message using the global logger.
func Debug(msg string, args ...any) {
	global.Debug(msg, args...)
}

// Info logs an info message using the global logger.
func Info(msg string, args ...any) {
	global.Info(msg, args...)
}

// Warn logs a warning message using the global logger.
func Warn(msg string, args ...any) {
	global.Warn(msg, args...)
}

// Error logs an error message using the global logger.
func Error(msg string, args ...any) {
	global.Error(msg, args...)
}
