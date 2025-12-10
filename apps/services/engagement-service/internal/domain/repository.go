package domain

import (
	"context"
)

// =============================================================================
// FAVORITE REPOSITORY
// =============================================================================

type FavoriteRepository interface {
	Create(ctx context.Context, favorite *Favorite) error
	Delete(ctx context.Context, id string) error
	GetByID(ctx context.Context, id string) (*Favorite, error)
	GetByUserAndOffer(ctx context.Context, userID, offerID string) (*Favorite, error)
	GetByUserID(ctx context.Context, userID string, filter FavoriteFilter) ([]*Favorite, int64, error)
	GetOfferIDsByUserID(ctx context.Context, userID string) ([]string, error)
	CountByOfferID(ctx context.Context, offerID string) (int64, error)
	ExistsByUserAndOffer(ctx context.Context, userID, offerID string) (bool, error)
}

type FavoriteFilter struct {
	Offset    int
	Limit     int
	SortBy    string
	SortOrder string
}

func DefaultFavoriteFilter() FavoriteFilter {
	return FavoriteFilter{
		Offset:    0,
		Limit:     20,
		SortBy:    "created_at",
		SortOrder: "desc",
	}
}

// =============================================================================
// REVIEW REPOSITORY
// =============================================================================

type ReviewRepository interface {
	Create(ctx context.Context, review *Review) error
	Update(ctx context.Context, review *Review) error
	Delete(ctx context.Context, id string) error
	GetByID(ctx context.Context, id string) (*Review, error)
	GetByUserAndOffer(ctx context.Context, userID, offerID string) (*Review, error)
	GetByOfferID(ctx context.Context, offerID string, filter ReviewFilter) ([]*Review, int64, error)
	GetByPartnerID(ctx context.Context, partnerID string, filter ReviewFilter) ([]*Review, int64, error)
	GetByUserID(ctx context.Context, userID string, filter ReviewFilter) ([]*Review, int64, error)
	GetPendingReviews(ctx context.Context, filter ReviewFilter) ([]*Review, int64, error)
	GetReportedReviews(ctx context.Context, filter ReviewFilter) ([]*Review, int64, error)
	GetAverageRating(ctx context.Context, offerID string) (float64, int64, error)
	GetPartnerAverageRating(ctx context.Context, partnerID string) (float64, int64, error)
}

type ReviewFilter struct {
	Status    []ReviewStatus
	Rating    *int
	Offset    int
	Limit     int
	SortBy    string
	SortOrder string
}

func DefaultReviewFilter() ReviewFilter {
	return ReviewFilter{
		Offset:    0,
		Limit:     20,
		SortBy:    "created_at",
		SortOrder: "desc",
	}
}

// =============================================================================
// EXTERNAL SERVICES
// =============================================================================

type OfferService interface {
	GetOfferInfo(ctx context.Context, offerID string) (*OfferInfo, error)
	IncrementFavoriteCount(ctx context.Context, offerID string) error
	DecrementFavoriteCount(ctx context.Context, offerID string) error
	UpdateAverageRating(ctx context.Context, offerID string, rating float64, count int64) error
}

type OfferInfo struct {
	OfferID         string
	PartnerID       string
	EstablishmentID string
	Title           string
	ImageURL        string
	PartnerName     string
}

type UserService interface {
	GetUserInfo(ctx context.Context, userID string) (*UserInfo, error)
	HasBooking(ctx context.Context, userID, offerID string) (bool, *string, error)
}

type UserInfo struct {
	UserID    string
	FirstName string
	Avatar    string
}
