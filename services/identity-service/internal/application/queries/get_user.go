package queries

import (
	"context"
	"time"

	"github.com/yousoon/services/identity/internal/domain"
)

// UserDTO represents a user for query responses.
type UserDTO struct {
	ID               string     `json:"id"`
	Email            string     `json:"email"`
	Phone            *string    `json:"phone,omitempty"`
	Profile          ProfileDTO `json:"profile"`
	Status           string     `json:"status"`
	Grade            string     `json:"grade"`
	EmailVerified    bool       `json:"emailVerified"`
	PhoneVerified    bool       `json:"phoneVerified"`
	IdentityVerified bool       `json:"identityVerified"`
	HasSubscription  bool       `json:"hasSubscription"`
	CreatedAt        time.Time  `json:"createdAt"`
	LastLoginAt      *time.Time `json:"lastLoginAt,omitempty"`
}

// ProfileDTO represents a user profile.
type ProfileDTO struct {
	FirstName   string     `json:"firstName"`
	LastName    string     `json:"lastName"`
	DisplayName string     `json:"displayName"`
	Avatar      *string    `json:"avatar,omitempty"`
	BirthDate   *time.Time `json:"birthDate,omitempty"`
	Gender      *string    `json:"gender,omitempty"`
}

// GetUserByIDQuery represents a query to get a user by ID.
type GetUserByIDQuery struct {
	UserID string
}

// GetUserByIDHandler handles the query.
type GetUserByIDHandler struct {
	userRepo domain.UserRepository
}

// NewGetUserByIDHandler creates a new handler.
func NewGetUserByIDHandler(userRepo domain.UserRepository) *GetUserByIDHandler {
	return &GetUserByIDHandler{userRepo: userRepo}
}

// Handle executes the query.
func (h *GetUserByIDHandler) Handle(ctx context.Context, q GetUserByIDQuery) (*UserDTO, error) {
	userID, err := domain.ParseUserID(q.UserID)
	if err != nil {
		return nil, err
	}

	user, err := h.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return mapUserToDTO(user), nil
}

// GetUserByEmailQuery represents a query to get a user by email.
type GetUserByEmailQuery struct {
	Email string
}

// GetUserByEmailHandler handles the query.
type GetUserByEmailHandler struct {
	userRepo domain.UserRepository
}

// NewGetUserByEmailHandler creates a new handler.
func NewGetUserByEmailHandler(userRepo domain.UserRepository) *GetUserByEmailHandler {
	return &GetUserByEmailHandler{userRepo: userRepo}
}

// Handle executes the query.
func (h *GetUserByEmailHandler) Handle(ctx context.Context, q GetUserByEmailQuery) (*UserDTO, error) {
	email, err := domain.NewEmail(q.Email)
	if err != nil {
		return nil, err
	}

	user, err := h.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	return mapUserToDTO(user), nil
}

// GetCurrentUserQuery gets the currently authenticated user.
type GetCurrentUserQuery struct {
	UserID string
}

// GetCurrentUserHandler handles the query.
type GetCurrentUserHandler struct {
	userRepo domain.UserRepository
}

// NewGetCurrentUserHandler creates a new handler.
func NewGetCurrentUserHandler(userRepo domain.UserRepository) *GetCurrentUserHandler {
	return &GetCurrentUserHandler{userRepo: userRepo}
}

// Handle executes the query.
func (h *GetCurrentUserHandler) Handle(ctx context.Context, q GetCurrentUserQuery) (*UserDTO, error) {
	userID, err := domain.ParseUserID(q.UserID)
	if err != nil {
		return nil, err
	}

	user, err := h.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return mapUserToDTO(user), nil
}

// mapUserToDTO converts a domain user to a DTO.
func mapUserToDTO(user *domain.User) *UserDTO {
	dto := &UserDTO{
		ID:               user.ID.String(),
		Email:            user.Email.String(),
		Status:           string(user.Status),
		Grade:            user.Grade.String(),
		EmailVerified:    user.EmailVerified,
		PhoneVerified:    user.PhoneVerified,
		IdentityVerified: user.HasVerifiedIdentity(),
		HasSubscription:  user.HasActiveSubscription(),
		CreatedAt:        user.CreatedAt,
		LastLoginAt:      user.LastLoginAt,
		Profile: ProfileDTO{
			FirstName:   user.Profile.FirstName,
			LastName:    user.Profile.LastName,
			DisplayName: user.Profile.DisplayName,
			Avatar:      user.Profile.Avatar,
			BirthDate:   user.Profile.BirthDate,
		},
	}

	if user.Phone != nil {
		phone := user.Phone.String()
		dto.Phone = &phone
	}

	if user.Profile.Gender != nil {
		gender := string(*user.Profile.Gender)
		dto.Profile.Gender = &gender
	}

	return dto
}
