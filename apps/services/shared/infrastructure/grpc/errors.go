package grpc

import (
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/yousoon/shared/domain"
)

// DomainErrorToGRPC converts a domain error to a gRPC status error.
func DomainErrorToGRPC(err error) error {
	if err == nil {
		return nil
	}

	// Check if it's already a gRPC status error
	if _, ok := status.FromError(err); ok {
		return err
	}

	// Handle domain errors by checking wrapped errors
	switch {
	case errors.Is(err, domain.ErrNotFound):
		return status.Errorf(codes.NotFound, "%s", err.Error())
	case errors.Is(err, domain.ErrValidation):
		return status.Errorf(codes.InvalidArgument, "%s", err.Error())
	case errors.Is(err, domain.ErrAlreadyExists), errors.Is(err, domain.ErrConflict):
		return status.Errorf(codes.AlreadyExists, "%s", err.Error())
	case errors.Is(err, domain.ErrUnauthorized):
		return status.Errorf(codes.Unauthenticated, "%s", err.Error())
	case errors.Is(err, domain.ErrForbidden):
		return status.Errorf(codes.PermissionDenied, "%s", err.Error())
	case errors.Is(err, domain.ErrExpired):
		return status.Errorf(codes.DeadlineExceeded, "%s", err.Error())
	case errors.Is(err, domain.ErrQuotaExceeded):
		return status.Errorf(codes.ResourceExhausted, "%s", err.Error())
	case errors.Is(err, domain.ErrInternal):
		return status.Errorf(codes.Internal, "%s", err.Error())
	default:
		return status.Errorf(codes.Unknown, "%s", err.Error())
	}
}

// GRPCToDomainError converts a gRPC status error to a domain error.
func GRPCToDomainError(err error) error {
	if err == nil {
		return nil
	}

	st, ok := status.FromError(err)
	if !ok {
		return domain.ErrInternal
	}

	switch st.Code() {
	case codes.NotFound:
		return domain.ErrNotFound
	case codes.InvalidArgument:
		return domain.ErrValidation
	case codes.AlreadyExists:
		return domain.ErrAlreadyExists
	case codes.Unauthenticated:
		return domain.ErrUnauthorized
	case codes.PermissionDenied:
		return domain.ErrForbidden
	case codes.DeadlineExceeded:
		return domain.ErrExpired
	case codes.ResourceExhausted:
		return domain.ErrQuotaExceeded
	case codes.Internal:
		return domain.ErrInternal
	default:
		return domain.ErrInternal
	}
}

// WrapError wraps an error in a gRPC status error with additional context.
func WrapError(code codes.Code, message string, details ...interface{}) error {
	return status.Errorf(code, message, details...)
}

// NotFoundError creates a NotFound error.
func NotFoundError(resource string) error {
	return status.Errorf(codes.NotFound, "%s not found", resource)
}

// InvalidArgumentError creates an InvalidArgument error.
func InvalidArgumentError(message string) error {
	return status.Errorf(codes.InvalidArgument, message)
}

// InternalError creates an Internal error.
func InternalError(message string) error {
	return status.Errorf(codes.Internal, message)
}

// UnauthenticatedError creates an Unauthenticated error.
func UnauthenticatedError(message string) error {
	return status.Errorf(codes.Unauthenticated, message)
}

// PermissionDeniedError creates a PermissionDenied error.
func PermissionDeniedError(message string) error {
	return status.Errorf(codes.PermissionDenied, message)
}

// AlreadyExistsError creates an AlreadyExists error.
func AlreadyExistsError(resource string) error {
	return status.Errorf(codes.AlreadyExists, "%s already exists", resource)
}

// FailedPreconditionError creates a FailedPrecondition error.
func FailedPreconditionError(message string) error {
	return status.Errorf(codes.FailedPrecondition, message)
}
