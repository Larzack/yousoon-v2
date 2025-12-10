package commands

import (
	"context"

	"github.com/yousoon/services/identity/internal/domain"
	"github.com/yousoon/shared/infrastructure/nats"
)

// RegisterUserCommand represents a command to register a new user.
type RegisterUserCommand struct {
	Email     string
	Password  string
	FirstName string
	LastName  string
	Platform  string
}

// RegisterUserResult represents the result of user registration.
type RegisterUserResult struct {
	UserID       string
	AccessToken  string
	RefreshToken string
}

// RegisterUserHandler handles user registration.
type RegisterUserHandler struct {
	userRepo       domain.UserRepository
	eventPublisher *nats.EventPublisher
}

// NewRegisterUserHandler creates a new RegisterUserHandler.
func NewRegisterUserHandler(
	userRepo domain.UserRepository,
	eventPublisher *nats.EventPublisher,
) *RegisterUserHandler {
	return &RegisterUserHandler{
		userRepo:       userRepo,
		eventPublisher: eventPublisher,
	}
}

// Handle executes the register user command.
func (h *RegisterUserHandler) Handle(ctx context.Context, cmd RegisterUserCommand) (*RegisterUserResult, error) {
	// Validate email
	email, err := domain.NewEmail(cmd.Email)
	if err != nil {
		return nil, err
	}

	// Check if user already exists
	exists, err := h.userRepo.ExistsByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, domain.ErrEmailAlreadyExists
	}

	// Create password hash
	password, err := domain.NewPassword(cmd.Password)
	if err != nil {
		return nil, err
	}

	// Create profile
	profile := domain.NewProfile(cmd.FirstName, cmd.LastName)

	// Create user
	user, err := domain.NewUser(email, password.Hash(), profile)
	if err != nil {
		return nil, err
	}

	// Save user
	if err := h.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	// Publish domain events
	for _, event := range user.GetDomainEvents() {
		if err := h.eventPublisher.Publish(ctx, event); err != nil {
			// Log error but don't fail registration
			// Events can be retried via outbox pattern
		}
	}
	user.ClearDomainEvents()

	// TODO: Generate tokens via auth service
	// For now, return user ID
	return &RegisterUserResult{
		UserID: user.ID.String(),
	}, nil
}
