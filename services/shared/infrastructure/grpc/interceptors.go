package grpc

import (
	"context"
	"fmt"
	"runtime/debug"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// LoggingInterceptor logs gRPC requests and responses.
func LoggingInterceptor(logger Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		start := time.Now()

		// Extract metadata
		md, _ := metadata.FromIncomingContext(ctx)

		// Call handler
		resp, err := handler(ctx, req)

		// Log request
		duration := time.Since(start)
		if err != nil {
			logger.Error("gRPC request failed",
				"method", info.FullMethod,
				"duration", duration,
				"error", err,
				"metadata", md,
			)
		} else {
			logger.Info("gRPC request completed",
				"method", info.FullMethod,
				"duration", duration,
				"metadata", md,
			)
		}

		return resp, err
	}
}

// RecoveryInterceptor recovers from panics.
func RecoveryInterceptor(logger Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
		defer func() {
			if r := recover(); r != nil {
				stack := string(debug.Stack())
				logger.Error("gRPC panic recovered",
					"method", info.FullMethod,
					"panic", r,
					"stack", stack,
				)
				err = status.Errorf(codes.Internal, "internal server error")
			}
		}()

		return handler(ctx, req)
	}
}

// TimeoutInterceptor adds a timeout to requests.
func TimeoutInterceptor(timeout time.Duration) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		return handler(ctx, req)
	}
}

// ValidationInterceptor validates requests.
func ValidationInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		if validator, ok := req.(Validator); ok {
			if err := validator.Validate(); err != nil {
				return nil, status.Errorf(codes.InvalidArgument, "validation failed: %v", err)
			}
		}

		return handler(ctx, req)
	}
}

// Validator is an interface for validatable requests.
type Validator interface {
	Validate() error
}

// AuthInterceptor handles authentication.
func AuthInterceptor(authenticator Authenticator) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		// Skip auth for health checks
		if info.FullMethod == "/grpc.health.v1.Health/Check" {
			return handler(ctx, req)
		}

		// Extract token from metadata
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Errorf(codes.Unauthenticated, "missing metadata")
		}

		tokens := md.Get("authorization")
		if len(tokens) == 0 {
			return nil, status.Errorf(codes.Unauthenticated, "missing authorization token")
		}

		// Authenticate
		claims, err := authenticator.Authenticate(ctx, tokens[0])
		if err != nil {
			return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %v", err)
		}

		// Add claims to context
		ctx = context.WithValue(ctx, claimsKey{}, claims)

		return handler(ctx, req)
	}
}

// Claims represents JWT claims.
type Claims struct {
	UserID  string
	Email   string
	Role    string
	IsAdmin bool
}

type claimsKey struct{}

// ClaimsFromContext extracts claims from context.
func ClaimsFromContext(ctx context.Context) (*Claims, error) {
	claims, ok := ctx.Value(claimsKey{}).(*Claims)
	if !ok {
		return nil, fmt.Errorf("claims not found in context")
	}
	return claims, nil
}

// Authenticator authenticates requests.
type Authenticator interface {
	Authenticate(ctx context.Context, token string) (*Claims, error)
}

// Logger is a simple logging interface.
type Logger interface {
	Info(msg string, keysAndValues ...interface{})
	Error(msg string, keysAndValues ...interface{})
	Debug(msg string, keysAndValues ...interface{})
	Warn(msg string, keysAndValues ...interface{})
}

// ChainUnaryInterceptors chains multiple unary interceptors.
func ChainUnaryInterceptors(interceptors ...grpc.UnaryServerInterceptor) grpc.ServerOption {
	return grpc.ChainUnaryInterceptor(interceptors...)
}

// ChainStreamInterceptors chains multiple stream interceptors.
func ChainStreamInterceptors(interceptors ...grpc.StreamServerInterceptor) grpc.ServerOption {
	return grpc.ChainStreamInterceptor(interceptors...)
}

// Client interceptors

// ClientLoggingInterceptor logs client requests.
func ClientLoggingInterceptor(logger Logger) grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		start := time.Now()
		err := invoker(ctx, method, req, reply, cc, opts...)
		duration := time.Since(start)

		if err != nil {
			logger.Error("gRPC client call failed",
				"method", method,
				"duration", duration,
				"error", err,
			)
		} else {
			logger.Debug("gRPC client call completed",
				"method", method,
				"duration", duration,
			)
		}

		return err
	}
}

// ClientRetryInterceptor retries failed requests.
func ClientRetryInterceptor(maxRetries int, backoff time.Duration) grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		var err error
		for i := 0; i <= maxRetries; i++ {
			err = invoker(ctx, method, req, reply, cc, opts...)
			if err == nil {
				return nil
			}

			// Only retry on transient errors
			st, ok := status.FromError(err)
			if !ok {
				return err
			}

			switch st.Code() {
			case codes.Unavailable, codes.ResourceExhausted, codes.Aborted:
				// Retry
				select {
				case <-ctx.Done():
					return ctx.Err()
				case <-time.After(backoff * time.Duration(i+1)):
					continue
				}
			default:
				return err
			}
		}

		return err
	}
}

// ClientAuthInterceptor adds authentication to client requests.
func ClientAuthInterceptor(token string) grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		ctx = metadata.AppendToOutgoingContext(ctx, "authorization", token)
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}
