package commands

import (
	"context"
	"time"

	"github.com/yousoon/services/identity/internal/domain"
	"github.com/yousoon/shared/infrastructure/nats"
)

// LoginCommand represents a command to authenticate a user.
type LoginCommand struct {
	Email     string
	Password  string
	Platform  string
	IP        string
	UserAgent string
}

// LoginResult represents the result of login.
type LoginResult struct {
	UserID       string
	AccessToken  string
	RefreshToken string
	ExpiresIn    int64
}

// LoginHandler handles user login.
type LoginHandler struct {
	userRepo       domain.UserRepository
	eventPublisher *nats.EventPublisher
}

// NewLoginHandler creates a new LoginHandler.
func NewLoginHandler(
	userRepo domain.UserRepository,
	eventPublisher *nats.EventPublisher,
) *LoginHandler {
	return &LoginHandler{
		userRepo:       userRepo,
		eventPublisher: eventPublisher,
	}
}

// Handle executes the login command.
func (h *LoginHandler) Handle(ctx context.Context, cmd LoginCommand) (*LoginResult, error) {
	// Validate email
	email, err := domain.NewEmail(cmd.Email)
	if err != nil {
		return nil, domain.ErrInvalidCredentials
	}

	// Find user
	user, err := h.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, domain.ErrInvalidCredentials
	}

	// Check password
	password := domain.NewPasswordFromHash(user.PasswordHash)
	if !password.Matches(cmd.Password) {
		return nil, domain.ErrInvalidCredentials
	}

	// Check user status
	if user.Status != domain.UserStatusActive {
		switch user.Status {
		case domain.UserStatusSuspended:
			return nil, domain.ErrUserSuspended
		case domain.UserStatusDeleted:
			return nil, domain.ErrUserDeleted
		default:
			return nil, domain.ErrUserNotActive
		}
	}

	// Update last login
	user.RecordLogin()
	if err := h.userRepo.Update(ctx, user); err != nil {
		// Log but don't fail login
	}

	// Parse platform
	platform := domain.Platform(cmd.Platform)

	// Emit login event
	event := domain.NewUserLoggedInEvent(user.ID, platform, cmd.IP, cmd.UserAgent)
	if err := h.eventPublisher.Publish(ctx, event); err != nil {
		// Log but don't fail login
	}

	// TODO: Generate JWT tokens
	// For now, return user ID and placeholder tokens

	return &LoginResult{
		UserID:    user.ID.String(),
		ExpiresIn: int64(6 * time.Hour / time.Second), // 6 hours
	}, nil
}
