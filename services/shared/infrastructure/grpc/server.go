// Package grpc provides gRPC infrastructure components.
package grpc

import (
	"context"
	"fmt"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
)

// ServerConfig holds gRPC server configuration.
type ServerConfig struct {
	// Port is the port to listen on.
	Port int
	// MaxRecvMsgSize is the maximum message size in bytes.
	MaxRecvMsgSize int
	// MaxSendMsgSize is the maximum send message size in bytes.
	MaxSendMsgSize int
	// EnableReflection enables gRPC reflection for debugging.
	EnableReflection bool
	// EnableHealthCheck enables the health check service.
	EnableHealthCheck bool
	// KeepaliveTime is the keepalive time.
	KeepaliveTime time.Duration
	// KeepaliveTimeout is the keepalive timeout.
	KeepaliveTimeout time.Duration
}

// DefaultServerConfig returns a default server configuration.
func DefaultServerConfig(port int) ServerConfig {
	return ServerConfig{
		Port:              port,
		MaxRecvMsgSize:    4 * 1024 * 1024, // 4MB
		MaxSendMsgSize:    4 * 1024 * 1024, // 4MB
		EnableReflection:  true,
		EnableHealthCheck: true,
		KeepaliveTime:     30 * time.Second,
		KeepaliveTimeout:  10 * time.Second,
	}
}

// Server wraps a gRPC server with additional functionality.
type Server struct {
	server       *grpc.Server
	healthServer *health.Server
	config       ServerConfig
	listener     net.Listener
}

// NewServer creates a new gRPC server.
func NewServer(cfg ServerConfig, opts ...grpc.ServerOption) *Server {
	// Default options
	defaultOpts := []grpc.ServerOption{
		grpc.MaxRecvMsgSize(cfg.MaxRecvMsgSize),
		grpc.MaxSendMsgSize(cfg.MaxSendMsgSize),
		grpc.KeepaliveParams(keepalive.ServerParameters{
			Time:    cfg.KeepaliveTime,
			Timeout: cfg.KeepaliveTimeout,
		}),
		grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
			MinTime:             5 * time.Second,
			PermitWithoutStream: true,
		}),
	}

	// Merge with provided options
	allOpts := append(defaultOpts, opts...)

	server := grpc.NewServer(allOpts...)
	s := &Server{
		server: server,
		config: cfg,
	}

	// Enable reflection
	if cfg.EnableReflection {
		reflection.Register(server)
	}

	// Enable health check
	if cfg.EnableHealthCheck {
		s.healthServer = health.NewServer()
		healthpb.RegisterHealthServer(server, s.healthServer)
	}

	return s
}

// Server returns the underlying gRPC server.
func (s *Server) Server() *grpc.Server {
	return s.server
}

// SetServingStatus sets the health status of a service.
func (s *Server) SetServingStatus(service string, serving bool) {
	if s.healthServer == nil {
		return
	}

	status := healthpb.HealthCheckResponse_NOT_SERVING
	if serving {
		status = healthpb.HealthCheckResponse_SERVING
	}

	s.healthServer.SetServingStatus(service, status)
}

// Start starts the gRPC server.
func (s *Server) Start() error {
	addr := fmt.Sprintf(":%d", s.config.Port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	s.listener = listener

	fmt.Printf("gRPC server listening on %s\n", addr)
	return s.server.Serve(listener)
}

// StartAsync starts the gRPC server in a goroutine.
func (s *Server) StartAsync() error {
	addr := fmt.Sprintf(":%d", s.config.Port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	s.listener = listener

	go func() {
		fmt.Printf("gRPC server listening on %s\n", addr)
		if err := s.server.Serve(listener); err != nil {
			fmt.Printf("gRPC server error: %v\n", err)
		}
	}()

	return nil
}

// Stop gracefully stops the gRPC server.
func (s *Server) Stop() {
	if s.healthServer != nil {
		s.healthServer.Shutdown()
	}
	s.server.GracefulStop()
}

// ForceStop immediately stops the gRPC server.
func (s *Server) ForceStop() {
	s.server.Stop()
}

// Address returns the server address.
func (s *Server) Address() string {
	if s.listener == nil {
		return ""
	}
	return s.listener.Addr().String()
}

// ClientConfig holds gRPC client configuration.
type ClientConfig struct {
	// Address is the server address.
	Address string
	// Timeout is the connection timeout.
	Timeout time.Duration
	// MaxRetries is the maximum number of retries.
	MaxRetries int
	// RetryBackoff is the backoff between retries.
	RetryBackoff time.Duration
	// Insecure disables TLS.
	Insecure bool
}

// DefaultClientConfig returns a default client configuration.
func DefaultClientConfig(address string) ClientConfig {
	return ClientConfig{
		Address:      address,
		Timeout:      10 * time.Second,
		MaxRetries:   3,
		RetryBackoff: 100 * time.Millisecond,
		Insecure:     true, // For development
	}
}

// NewClientConn creates a new gRPC client connection.
func NewClientConn(ctx context.Context, cfg ClientConfig, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	defaultOpts := []grpc.DialOption{
		grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(4*1024*1024),
			grpc.MaxCallSendMsgSize(4*1024*1024),
		),
	}

	if cfg.Insecure {
		defaultOpts = append(defaultOpts, grpc.WithInsecure())
	}

	allOpts := append(defaultOpts, opts...)

	ctx, cancel := context.WithTimeout(ctx, cfg.Timeout)
	defer cancel()

	conn, err := grpc.DialContext(ctx, cfg.Address, allOpts...)
	if err != nil {
		return nil, fmt.Errorf("failed to connect: %w", err)
	}

	return conn, nil
}
