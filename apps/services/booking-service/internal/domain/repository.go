package domain

import (
	"context"
	"time"
)

// =============================================================================
// REPOSITORY INTERFACE
// =============================================================================

// OutingRepository defines the interface for outing persistence
type OutingRepository interface {
	// Create persists a new outing
	Create(ctx context.Context, outing *Outing) error

	// Update updates an existing outing
	Update(ctx context.Context, outing *Outing) error

	// GetByID retrieves an outing by ID
	GetByID(ctx context.Context, id string) (*Outing, error)

	// GetByQRCode retrieves an outing by QR code
	GetByQRCode(ctx context.Context, qrCode string) (*Outing, error)

	// GetByUserID retrieves all outings for a user with pagination
	GetByUserID(ctx context.Context, userID string, filter OutingFilter) ([]*Outing, int64, error)

	// GetByOfferID retrieves all outings for an offer
	GetByOfferID(ctx context.Context, offerID string, filter OutingFilter) ([]*Outing, int64, error)

	// GetByPartnerID retrieves all outings for a partner
	GetByPartnerID(ctx context.Context, partnerID string, filter OutingFilter) ([]*Outing, int64, error)

	// GetByEstablishmentID retrieves all outings for an establishment
	GetByEstablishmentID(ctx context.Context, establishmentID string, filter OutingFilter) ([]*Outing, int64, error)

	// GetActiveByUserAndOffer checks if user has an active outing for the offer
	GetActiveByUserAndOffer(ctx context.Context, userID, offerID string) (*Outing, error)

	// CountByUserAndOffer counts user bookings for an offer (for quota check)
	CountByUserAndOffer(ctx context.Context, userID, offerID string) (int64, error)

	// CountByOfferAndPeriod counts bookings for an offer in a period (for quota check)
	CountByOfferAndPeriod(ctx context.Context, offerID string, start, end time.Time) (int64, error)

	// GetExpiredOutings retrieves outings that have expired but not marked
	GetExpiredOutings(ctx context.Context, before time.Time, limit int) ([]*Outing, error)

	// Delete removes an outing (soft delete via status)
	Delete(ctx context.Context, id string) error
}

// OutingFilter contains filter options for listing outings
type OutingFilter struct {
	Status    []OutingStatus
	StartDate *time.Time
	EndDate   *time.Time
	Offset    int
	Limit     int
	SortBy    string
	SortOrder string
}

func DefaultOutingFilter() OutingFilter {
	return OutingFilter{
		Offset:    0,
		Limit:     20,
		SortBy:    "created_at",
		SortOrder: "desc",
	}
}

// =============================================================================
// DOMAIN SERVICE INTERFACES
// =============================================================================

// OfferService provides offer information for booking
type OfferService interface {
	// GetOfferSnapshot retrieves offer details for creating a snapshot
	GetOfferSnapshot(ctx context.Context, offerID string) (*OfferSnapshot, error)

	// CanBook checks if an offer can be booked
	CanBook(ctx context.Context, offerID string) error

	// IncrementBookingCount increments the offer's booking count
	IncrementBookingCount(ctx context.Context, offerID string) error

	// DecrementBookingCount decrements the offer's booking count (on cancel)
	DecrementBookingCount(ctx context.Context, offerID string) error
}

// UserService provides user information for booking
type UserService interface {
	// GetUserSnapshot retrieves user details for creating a snapshot
	GetUserSnapshot(ctx context.Context, userID string) (*UserSnapshot, error)

	// CanBook checks if a user can book (verified, active subscription, etc.)
	CanBook(ctx context.Context, userID string) error
}

// NotificationService sends booking notifications
type NotificationService interface {
	// SendBookingConfirmation sends booking confirmation to user
	SendBookingConfirmation(ctx context.Context, outing *Outing) error

	// SendCheckInConfirmation sends check-in confirmation
	SendCheckInConfirmation(ctx context.Context, outing *Outing) error

	// SendCancellationNotification sends cancellation notification
	SendCancellationNotification(ctx context.Context, outing *Outing) error

	// SendExpirationReminder sends reminder before expiration
	SendExpirationReminder(ctx context.Context, outing *Outing) error
}
