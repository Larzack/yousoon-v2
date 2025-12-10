// Package domain contains repository interfaces for the Discovery service.
package domain

import "context"

// =============================================================================
// Offer Repository
// =============================================================================

// OfferRepository defines the interface for offer persistence.
type OfferRepository interface {
	// Save persists an offer (create or update).
	Save(ctx context.Context, offer *Offer) error

	// FindByID retrieves an offer by ID.
	FindByID(ctx context.Context, id OfferID) (*Offer, error)

	// FindByPartnerID retrieves all offers for a partner.
	FindByPartnerID(ctx context.Context, partnerID PartnerID) ([]*Offer, error)

	// FindByEstablishmentID retrieves all offers for an establishment.
	FindByEstablishmentID(ctx context.Context, establishmentID EstablishmentID) ([]*Offer, error)

	// FindByCategory retrieves offers in a category.
	FindByCategory(ctx context.Context, categoryID CategoryID, offset, limit int) ([]*Offer, error)

	// List retrieves offers with filters.
	List(ctx context.Context, filter OfferFilter) (*OfferListResult, error)

	// Delete soft-deletes an offer.
	Delete(ctx context.Context, id OfferID) error

	// Count returns the total count of offers matching the filter.
	Count(ctx context.Context, filter OfferFilter) (int64, error)

	// ExistsActiveForPartner checks if a partner has any active offers.
	ExistsActiveForPartner(ctx context.Context, partnerID PartnerID) (bool, error)
}

// OfferFilter defines filters for listing offers.
type OfferFilter struct {
	// Location-based search
	Location  *GeoLocation
	Latitude  *float64
	Longitude *float64
	RadiusKm  float64

	// Filtering
	PartnerID       *PartnerID
	EstablishmentID *EstablishmentID
	CategoryID      *CategoryID
	Status          *OfferStatus
	Tags            []string
	DiscountType    *string
	MinRating       *float64

	// Search
	SearchQuery string

	// Availability
	OnlyActive       bool
	OnlyAvailableNow bool
	ActiveOnly       bool

	// Moderation
	ModerationStatus *ModerationStatus

	// Pagination
	Offset int
	Limit  int

	// Sorting
	SortBy    string // "distance", "created_at", "discount", "popularity"
	SortOrder string // "asc", "desc"
}

// OfferListResult represents the result of listing offers.
type OfferListResult struct {
	Offers     []*Offer
	TotalCount int64
	Offset     int
	Limit      int
}

// =============================================================================
// Offer Read Repository (for optimized queries)
// =============================================================================

// OfferReadRepository defines read-only operations for offers.
type OfferReadRepository interface {
	// GetOfferSummaries returns offer summaries for lists.
	GetOfferSummaries(ctx context.Context, filter OfferFilter) ([]OfferSummary, int64, error)

	// GetOffersNearLocation returns offers near a location.
	GetOffersNearLocation(ctx context.Context, location GeoLocation, radiusKm float64, limit int) ([]OfferSummary, error)

	// GetOffersByIDs returns offers by IDs.
	GetOffersByIDs(ctx context.Context, ids []OfferID) ([]*Offer, error)

	// SearchOffers performs a full-text search on offers.
	SearchOffers(ctx context.Context, query string, location *GeoLocation, limit int) ([]OfferSummary, error)

	// GetRecommendedOffers returns personalized offers for a user.
	GetRecommendedOffers(ctx context.Context, userID UserID, userCategories []CategoryID, location GeoLocation, limit int) ([]OfferSummary, error)

	// GetTrendingOffers returns trending offers.
	GetTrendingOffers(ctx context.Context, location *GeoLocation, limit int) ([]OfferSummary, error)

	// GetNewOffers returns recently published offers.
	GetNewOffers(ctx context.Context, location *GeoLocation, limit int) ([]OfferSummary, error)

	// GetExpiringOffers returns offers expiring soon.
	GetExpiringOffers(ctx context.Context, withinDays int, limit int) ([]OfferSummary, error)

	// GetOfferCountByCategory returns offer counts per category.
	GetOfferCountByCategory(ctx context.Context) (map[CategoryID]int, error)

	// GetPartnerStats returns statistics for a partner's offers.
	GetPartnerOfferStats(ctx context.Context, partnerID PartnerID) (*PartnerOfferStats, error)
}

// PartnerOfferStats represents statistics for a partner's offers.
type PartnerOfferStats struct {
	TotalOffers   int
	ActiveOffers  int
	TotalViews    int
	TotalBookings int
	TotalCheckins int
	AverageRating float64
}

// =============================================================================
// Category Repository
// =============================================================================

// CategoryRepository defines the interface for category persistence.
type CategoryRepository interface {
	// Save persists a category (create or update).
	Save(ctx context.Context, category *Category) error

	// FindByID retrieves a category by ID.
	FindByID(ctx context.Context, id CategoryID) (*Category, error)

	// FindBySlug retrieves a category by slug.
	FindBySlug(ctx context.Context, slug string) (*Category, error)

	// FindAll retrieves all categories.
	FindAll(ctx context.Context) ([]*Category, error)

	// FindActive retrieves all active categories.
	FindActive(ctx context.Context) ([]*Category, error)

	// FindByParentID retrieves categories by parent ID.
	FindByParentID(ctx context.Context, parentID *CategoryID) ([]*Category, error)

	// FindRootCategories retrieves root categories (no parent).
	FindRootCategories(ctx context.Context) ([]*Category, error)

	// Delete deletes a category.
	Delete(ctx context.Context, id CategoryID) error

	// ExistsBySlug checks if a category with the slug exists.
	ExistsBySlug(ctx context.Context, slug string) (bool, error)

	// GetCategoryTree returns the full category tree.
	GetCategoryTree(ctx context.Context) ([]*CategoryTree, error)

	// GetCategorySummaries returns category summaries with offer counts.
	GetCategorySummaries(ctx context.Context) ([]CategorySummary, error)
}

// =============================================================================
// Unit of Work (for transactions)
// =============================================================================

// UnitOfWork provides transactional support for repositories.
type UnitOfWork interface {
	// Begin starts a new transaction.
	Begin(ctx context.Context) (context.Context, error)

	// Commit commits the transaction.
	Commit(ctx context.Context) error

	// Rollback rolls back the transaction.
	Rollback(ctx context.Context) error
}
