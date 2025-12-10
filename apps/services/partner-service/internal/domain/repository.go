package domain

import "context"

// =============================================================================
// Repository Interfaces
// =============================================================================

// PartnerRepository is the repository interface for Partner aggregate.
type PartnerRepository interface {
	// Save saves a partner (insert or update).
	Save(ctx context.Context, partner *Partner) error

	// FindByID finds a partner by ID.
	FindByID(ctx context.Context, id PartnerID) (*Partner, error)

	// FindByOwnerUserID finds a partner by owner user ID.
	FindByOwnerUserID(ctx context.Context, userID UserID) (*Partner, error)

	// FindBySIRET finds a partner by SIRET.
	FindBySIRET(ctx context.Context, siret string) (*Partner, error)

	// FindByTeamMemberEmail finds partners where the email is a team member.
	FindByTeamMemberEmail(ctx context.Context, email Email) ([]*Partner, error)

	// Delete deletes a partner.
	Delete(ctx context.Context, id PartnerID) error

	// List lists partners with pagination and filters.
	List(ctx context.Context, filter PartnerFilter) ([]*Partner, int64, error)

	// Count counts partners matching a filter.
	Count(ctx context.Context, filter PartnerFilter) (int64, error)
}

// PartnerFilter contains filtering options for partner queries.
type PartnerFilter struct {
	// Status filter
	Status *PartnerStatus

	// Category filter
	Category string

	// Search term (matches company name, trade name)
	Search string

	// Pagination
	Offset int
	Limit  int
}

// NewPartnerFilter creates a new partner filter with defaults.
func NewPartnerFilter() PartnerFilter {
	return PartnerFilter{
		Offset: 0,
		Limit:  20,
	}
}

// WithStatus adds a status filter.
func (f PartnerFilter) WithStatus(status PartnerStatus) PartnerFilter {
	f.Status = &status
	return f
}

// WithCategory adds a category filter.
func (f PartnerFilter) WithCategory(category string) PartnerFilter {
	f.Category = category
	return f
}

// WithSearch adds a search term.
func (f PartnerFilter) WithSearch(search string) PartnerFilter {
	f.Search = search
	return f
}

// WithPagination adds pagination.
func (f PartnerFilter) WithPagination(offset, limit int) PartnerFilter {
	f.Offset = offset
	if limit > 0 && limit <= 100 {
		f.Limit = limit
	}
	return f
}

// =============================================================================
// Read Models (Queries)
// =============================================================================

// PartnerReadRepository is the read-side repository for partner queries.
type PartnerReadRepository interface {
	// GetPartnerSummary gets a partner summary for display.
	GetPartnerSummary(ctx context.Context, id PartnerID) (*PartnerSummary, error)

	// GetEstablishmentsNearLocation gets establishments near a location.
	GetEstablishmentsNearLocation(ctx context.Context, location GeoLocation, radiusKm float64, limit int) ([]EstablishmentSummary, error)

	// GetPartnerByEstablishmentID gets a partner by establishment ID.
	GetPartnerByEstablishmentID(ctx context.Context, estID EstablishmentID) (*Partner, error)

	// SearchPartners searches partners by name.
	SearchPartners(ctx context.Context, query string, limit int) ([]*PartnerSummary, error)
}

// =============================================================================
// Read Model DTOs
// =============================================================================

// PartnerSummary is a summary view of a partner.
type PartnerSummary struct {
	ID                 PartnerID     `json:"id"`
	CompanyName        string        `json:"companyName"`
	TradeName          string        `json:"tradeName"`
	Category           string        `json:"category"`
	Logo               string        `json:"logo"`
	Status             PartnerStatus `json:"status"`
	EstablishmentCount int           `json:"establishmentCount"`
	ActiveOfferCount   int           `json:"activeOfferCount"`
	AvgRating          float64       `json:"avgRating"`
	ReviewCount        int           `json:"reviewCount"`
}

// EstablishmentSummary is a summary view of an establishment.
type EstablishmentSummary struct {
	ID           EstablishmentID `json:"id"`
	PartnerID    PartnerID       `json:"partnerId"`
	Name         string          `json:"name"`
	Type         string          `json:"type"`
	Address      string          `json:"address"`
	City         string          `json:"city"`
	Location     GeoLocation     `json:"location"`
	PrimaryImage string          `json:"primaryImage"`
	DistanceKm   float64         `json:"distanceKm"`
}
