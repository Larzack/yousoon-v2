// Package queries contains query handlers for offers.
package queries

import (
	"context"

	"github.com/yousoon/discovery-service/internal/domain"
)

// =============================================================================
// Get Offer Query
// =============================================================================

// GetOfferQuery retrieves a single offer by ID.
type GetOfferQuery struct {
	OfferID string
}

// GetOfferHandler handles the get offer query.
type GetOfferHandler struct {
	offerRepo domain.OfferRepository
}

// NewGetOfferHandler creates a new GetOfferHandler.
func NewGetOfferHandler(offerRepo domain.OfferRepository) *GetOfferHandler {
	return &GetOfferHandler{
		offerRepo: offerRepo,
	}
}

// Handle executes the get offer query.
func (h *GetOfferHandler) Handle(ctx context.Context, query GetOfferQuery) (*domain.Offer, error) {
	offer, err := h.offerRepo.FindByID(ctx, domain.OfferID(query.OfferID))
	if err != nil {
		return nil, err
	}
	if offer == nil {
		return nil, domain.ErrOfferNotFound
	}
	return offer, nil
}

// =============================================================================
// List Offers Query
// =============================================================================

// ListOffersQuery retrieves a list of offers with filters.
type ListOffersQuery struct {
	// Location-based search
	Longitude *float64
	Latitude  *float64
	RadiusKm  float64

	// Filtering
	PartnerID       *string
	EstablishmentID *string
	CategoryID      *string
	Status          *string
	Tags            []string

	// Search
	SearchQuery string

	// Availability
	OnlyActive       bool
	OnlyAvailableNow bool

	// Moderation (admin)
	ModerationStatus *string

	// Pagination
	Offset int
	Limit  int

	// Sorting
	SortBy    string // "distance", "created_at", "discount", "popularity"
	SortOrder string // "asc", "desc"
}

// ListOffersResult represents the result of listing offers.
type ListOffersResult struct {
	Offers     []*domain.Offer
	TotalCount int64
	Offset     int
	Limit      int
	HasMore    bool
}

// ListOffersHandler handles the list offers query.
type ListOffersHandler struct {
	offerRepo domain.OfferRepository
}

// NewListOffersHandler creates a new ListOffersHandler.
func NewListOffersHandler(offerRepo domain.OfferRepository) *ListOffersHandler {
	return &ListOffersHandler{
		offerRepo: offerRepo,
	}
}

// Handle executes the list offers query.
func (h *ListOffersHandler) Handle(ctx context.Context, query ListOffersQuery) (*ListOffersResult, error) {
	// Build filter
	filter := domain.OfferFilter{
		RadiusKm:         query.RadiusKm,
		Tags:             query.Tags,
		SearchQuery:      query.SearchQuery,
		OnlyActive:       query.OnlyActive,
		OnlyAvailableNow: query.OnlyAvailableNow,
		Offset:           query.Offset,
		Limit:            query.Limit,
		SortBy:           query.SortBy,
		SortOrder:        query.SortOrder,
	}

	// Set location if provided
	if query.Longitude != nil && query.Latitude != nil {
		location, err := domain.NewGeoLocation(*query.Longitude, *query.Latitude)
		if err != nil {
			return nil, err
		}
		filter.Location = &location
	}

	// Set optional filters
	if query.PartnerID != nil {
		partnerID := domain.PartnerID(*query.PartnerID)
		filter.PartnerID = &partnerID
	}
	if query.EstablishmentID != nil {
		establishmentID := domain.EstablishmentID(*query.EstablishmentID)
		filter.EstablishmentID = &establishmentID
	}
	if query.CategoryID != nil {
		categoryID := domain.CategoryID(*query.CategoryID)
		filter.CategoryID = &categoryID
	}
	if query.Status != nil {
		status := domain.OfferStatus(*query.Status)
		filter.Status = &status
	}
	if query.ModerationStatus != nil {
		modStatus := domain.ModerationStatus(*query.ModerationStatus)
		filter.ModerationStatus = &modStatus
	}

	// Set defaults
	if filter.Limit == 0 {
		filter.Limit = 20
	}
	if filter.RadiusKm == 0 {
		filter.RadiusKm = 10 // Default 10km radius
	}

	// Execute query
	result, err := h.offerRepo.List(ctx, filter)
	if err != nil {
		return nil, err
	}

	return &ListOffersResult{
		Offers:     result.Offers,
		TotalCount: result.TotalCount,
		Offset:     result.Offset,
		Limit:      result.Limit,
		HasMore:    int64(result.Offset+len(result.Offers)) < result.TotalCount,
	}, nil
}

// =============================================================================
// Get Partner Offers Query
// =============================================================================

// GetPartnerOffersQuery retrieves all offers for a partner.
type GetPartnerOffersQuery struct {
	PartnerID string
	Status    *string
	Offset    int
	Limit     int
}

// GetPartnerOffersHandler handles the get partner offers query.
type GetPartnerOffersHandler struct {
	offerRepo domain.OfferRepository
}

// NewGetPartnerOffersHandler creates a new GetPartnerOffersHandler.
func NewGetPartnerOffersHandler(offerRepo domain.OfferRepository) *GetPartnerOffersHandler {
	return &GetPartnerOffersHandler{
		offerRepo: offerRepo,
	}
}

// Handle executes the get partner offers query.
func (h *GetPartnerOffersHandler) Handle(ctx context.Context, query GetPartnerOffersQuery) (*ListOffersResult, error) {
	partnerID := domain.PartnerID(query.PartnerID)
	filter := domain.OfferFilter{
		PartnerID: &partnerID,
		Offset:    query.Offset,
		Limit:     query.Limit,
		SortBy:    "created_at",
		SortOrder: "desc",
	}

	if query.Status != nil {
		status := domain.OfferStatus(*query.Status)
		filter.Status = &status
	}

	if filter.Limit == 0 {
		filter.Limit = 20
	}

	result, err := h.offerRepo.List(ctx, filter)
	if err != nil {
		return nil, err
	}

	return &ListOffersResult{
		Offers:     result.Offers,
		TotalCount: result.TotalCount,
		Offset:     result.Offset,
		Limit:      result.Limit,
		HasMore:    int64(result.Offset+len(result.Offers)) < result.TotalCount,
	}, nil
}

// =============================================================================
// Search Offers Query
// =============================================================================

// SearchOffersQuery searches for offers.
type SearchOffersQuery struct {
	Query     string
	Longitude *float64
	Latitude  *float64
	Limit     int
}

// SearchOffersHandler handles the search offers query.
type SearchOffersHandler struct {
	readRepo domain.OfferReadRepository
}

// NewSearchOffersHandler creates a new SearchOffersHandler.
func NewSearchOffersHandler(readRepo domain.OfferReadRepository) *SearchOffersHandler {
	return &SearchOffersHandler{
		readRepo: readRepo,
	}
}

// Handle executes the search offers query.
func (h *SearchOffersHandler) Handle(ctx context.Context, query SearchOffersQuery) ([]domain.OfferSummary, error) {
	var location *domain.GeoLocation
	if query.Longitude != nil && query.Latitude != nil {
		loc, err := domain.NewGeoLocation(*query.Longitude, *query.Latitude)
		if err != nil {
			return nil, err
		}
		location = &loc
	}

	limit := query.Limit
	if limit == 0 {
		limit = 20
	}

	return h.readRepo.SearchOffers(ctx, query.Query, location, limit)
}

// =============================================================================
// Get Nearby Offers Query
// =============================================================================

// GetNearbyOffersQuery retrieves offers near a location.
type GetNearbyOffersQuery struct {
	Longitude float64
	Latitude  float64
	RadiusKm  float64
	Limit     int
}

// GetNearbyOffersHandler handles the get nearby offers query.
type GetNearbyOffersHandler struct {
	readRepo domain.OfferReadRepository
}

// NewGetNearbyOffersHandler creates a new GetNearbyOffersHandler.
func NewGetNearbyOffersHandler(readRepo domain.OfferReadRepository) *GetNearbyOffersHandler {
	return &GetNearbyOffersHandler{
		readRepo: readRepo,
	}
}

// Handle executes the get nearby offers query.
func (h *GetNearbyOffersHandler) Handle(ctx context.Context, query GetNearbyOffersQuery) ([]domain.OfferSummary, error) {
	location, err := domain.NewGeoLocation(query.Longitude, query.Latitude)
	if err != nil {
		return nil, err
	}

	radiusKm := query.RadiusKm
	if radiusKm == 0 {
		radiusKm = 10 // Default 10km
	}

	limit := query.Limit
	if limit == 0 {
		limit = 20
	}

	return h.readRepo.GetOffersNearLocation(ctx, location, radiusKm, limit)
}

// =============================================================================
// Get Recommended Offers Query
// =============================================================================

// GetRecommendedOffersQuery retrieves personalized offers for a user.
type GetRecommendedOffersQuery struct {
	UserID         string
	UserCategories []string
	Longitude      float64
	Latitude       float64
	Limit          int
}

// GetRecommendedOffersHandler handles the get recommended offers query.
type GetRecommendedOffersHandler struct {
	readRepo domain.OfferReadRepository
}

// NewGetRecommendedOffersHandler creates a new GetRecommendedOffersHandler.
func NewGetRecommendedOffersHandler(readRepo domain.OfferReadRepository) *GetRecommendedOffersHandler {
	return &GetRecommendedOffersHandler{
		readRepo: readRepo,
	}
}

// Handle executes the get recommended offers query.
func (h *GetRecommendedOffersHandler) Handle(ctx context.Context, query GetRecommendedOffersQuery) ([]domain.OfferSummary, error) {
	location, err := domain.NewGeoLocation(query.Longitude, query.Latitude)
	if err != nil {
		return nil, err
	}

	// Convert category IDs
	categoryIDs := make([]domain.CategoryID, len(query.UserCategories))
	for i, id := range query.UserCategories {
		categoryIDs[i] = domain.CategoryID(id)
	}

	limit := query.Limit
	if limit == 0 {
		limit = 20
	}

	return h.readRepo.GetRecommendedOffers(
		ctx,
		domain.UserID(query.UserID),
		categoryIDs,
		location,
		limit,
	)
}

// =============================================================================
// Get Trending Offers Query
// =============================================================================

// GetTrendingOffersQuery retrieves trending offers.
type GetTrendingOffersQuery struct {
	Longitude *float64
	Latitude  *float64
	Limit     int
}

// GetTrendingOffersHandler handles the get trending offers query.
type GetTrendingOffersHandler struct {
	readRepo domain.OfferReadRepository
}

// NewGetTrendingOffersHandler creates a new GetTrendingOffersHandler.
func NewGetTrendingOffersHandler(readRepo domain.OfferReadRepository) *GetTrendingOffersHandler {
	return &GetTrendingOffersHandler{
		readRepo: readRepo,
	}
}

// Handle executes the get trending offers query.
func (h *GetTrendingOffersHandler) Handle(ctx context.Context, query GetTrendingOffersQuery) ([]domain.OfferSummary, error) {
	var location *domain.GeoLocation
	if query.Longitude != nil && query.Latitude != nil {
		loc, err := domain.NewGeoLocation(*query.Longitude, *query.Latitude)
		if err != nil {
			return nil, err
		}
		location = &loc
	}

	limit := query.Limit
	if limit == 0 {
		limit = 20
	}

	return h.readRepo.GetTrendingOffers(ctx, location, limit)
}

// =============================================================================
// Get Offer Snapshot Query (for Booking Service)
// =============================================================================

// GetOfferSnapshotQuery retrieves an offer snapshot for booking.
type GetOfferSnapshotQuery struct {
	OfferID string
}

// GetOfferSnapshotHandler handles the get offer snapshot query.
type GetOfferSnapshotHandler struct {
	offerRepo    domain.OfferRepository
	categoryRepo domain.CategoryRepository
}

// NewGetOfferSnapshotHandler creates a new GetOfferSnapshotHandler.
func NewGetOfferSnapshotHandler(offerRepo domain.OfferRepository, categoryRepo domain.CategoryRepository) *GetOfferSnapshotHandler {
	return &GetOfferSnapshotHandler{
		offerRepo:    offerRepo,
		categoryRepo: categoryRepo,
	}
}

// Handle executes the get offer snapshot query.
func (h *GetOfferSnapshotHandler) Handle(ctx context.Context, query GetOfferSnapshotQuery) (*domain.OfferSnapshot, error) {
	offer, err := h.offerRepo.FindByID(ctx, domain.OfferID(query.OfferID))
	if err != nil {
		return nil, err
	}
	if offer == nil {
		return nil, domain.ErrOfferNotFound
	}

	// Get category name
	categoryName := ""
	category, err := h.categoryRepo.FindByID(ctx, offer.CategoryID())
	if err == nil && category != nil {
		categoryName = category.Name().FR
	}

	snapshot := offer.ToSnapshot()
	snapshot.CategoryName = categoryName

	return &snapshot, nil
}

// =============================================================================
// Can Book Offer Query
// =============================================================================

// CanBookOfferQuery checks if an offer can be booked.
type CanBookOfferQuery struct {
	OfferID          string
	UserBookingCount int
}

// CanBookOfferResult represents the result of can book check.
type CanBookOfferResult struct {
	CanBook bool
	Reason  string
}

// CanBookOfferHandler handles the can book offer query.
type CanBookOfferHandler struct {
	offerRepo domain.OfferRepository
}

// NewCanBookOfferHandler creates a new CanBookOfferHandler.
func NewCanBookOfferHandler(offerRepo domain.OfferRepository) *CanBookOfferHandler {
	return &CanBookOfferHandler{
		offerRepo: offerRepo,
	}
}

// Handle executes the can book offer query.
func (h *CanBookOfferHandler) Handle(ctx context.Context, query CanBookOfferQuery) (*CanBookOfferResult, error) {
	offer, err := h.offerRepo.FindByID(ctx, domain.OfferID(query.OfferID))
	if err != nil {
		return nil, err
	}
	if offer == nil {
		return &CanBookOfferResult{
			CanBook: false,
			Reason:  "Offer not found",
		}, nil
	}

	if err := offer.CanUserBook(query.UserBookingCount); err != nil {
		return &CanBookOfferResult{
			CanBook: false,
			Reason:  err.Error(),
		}, nil
	}

	return &CanBookOfferResult{
		CanBook: true,
		Reason:  "",
	}, nil
}
