package resolver

import (
	"github.com/yousoon/services/identity/internal/application/commands"
	"github.com/yousoon/services/identity/internal/application/queries"
	"github.com/yousoon/services/identity/internal/domain"
	"github.com/yousoon/services/shared/infrastructure/nats"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

// Resolver is the root resolver for the Identity service.
type Resolver struct {
	// Repositories
	userRepo domain.UserRepository

	// Event Publisher
	eventPublisher *nats.EventPublisher

	// Command Handlers
	registerUserHandler               *commands.RegisterUserHandler
	loginHandler                      *commands.LoginHandler
	updateProfileHandler              *commands.UpdateProfileHandler
	submitIdentityVerificationHandler *commands.SubmitIdentityVerificationHandler
	verifyIdentityHandler             *commands.VerifyIdentityHandler
	rejectIdentityHandler             *commands.RejectIdentityHandler

	// Query Handlers
	getUserByIDHandler    *queries.GetUserByIDHandler
	getUserByEmailHandler *queries.GetUserByEmailHandler
	getCurrentUserHandler *queries.GetCurrentUserHandler
}

// NewResolver creates a new Resolver.
func NewResolver(
	userRepo domain.UserRepository,
	eventPublisher *nats.EventPublisher,
) *Resolver {
	return &Resolver{
		userRepo:       userRepo,
		eventPublisher: eventPublisher,

		// Initialize command handlers
		registerUserHandler:               commands.NewRegisterUserHandler(userRepo, eventPublisher),
		loginHandler:                      commands.NewLoginHandler(userRepo, eventPublisher),
		updateProfileHandler:              commands.NewUpdateProfileHandler(userRepo, eventPublisher),
		submitIdentityVerificationHandler: commands.NewSubmitIdentityVerificationHandler(userRepo, eventPublisher),
		verifyIdentityHandler:             commands.NewVerifyIdentityHandler(userRepo, eventPublisher),
		rejectIdentityHandler:             commands.NewRejectIdentityHandler(userRepo, eventPublisher),

		// Initialize query handlers
		getUserByIDHandler:    queries.NewGetUserByIDHandler(userRepo),
		getUserByEmailHandler: queries.NewGetUserByEmailHandler(userRepo),
		getCurrentUserHandler: queries.NewGetCurrentUserHandler(userRepo),
	}
}
