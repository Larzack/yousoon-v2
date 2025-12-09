// Package resolver contains GraphQL resolvers for the Discovery service.
package resolver

import (
	"context"
	"fmt"
	"time"

	"github.com/yousoon/discovery-service/internal/application/commands"
	"github.com/yousoon/discovery-service/internal/application/queries"
	"github.com/yousoon/discovery-service/internal/domain"
	"github.com/yousoon/discovery-service/internal/interface/graphql/model"
	"github.com/yousoon/shared/infrastructure/nats"
)

// Resolver is the root resolver.
type Resolver struct {
	// Repositories
	offerRepo    domain.OfferRepository
	categoryRepo domain.CategoryRepository
	searchRepo   domain.OfferSearchService

	// Command handlers
	createOfferHandler    *commands.CreateOfferHandler
	publishOfferHandler   *commands.PublishOfferHandler
	expireOfferHandler    *commands.ExpireOfferHandler
	archiveOfferHandler   *commands.ArchiveOfferHandler
	extendOfferHandler    *commands.ExtendOfferHandler
	createCategoryHandler *commands.CreateCategoryHandler
	updateCategoryHandler *commands.UpdateCategoryHandler
	deleteCategoryHandler *commands.DeleteCategoryHandler

	// Query handlers
	getOfferHandler           *queries.GetOfferHandler
	listOffersHandler         *queries.ListOffersHandler
	searchOffersHandler       *queries.SearchOffersHandler
	getOffersByPartnerHandler *queries.GetOffersByPartnerHandler
	getNearbyOffersHandler    *queries.GetNearbyOffersHandler
	getTrendingOffersHandler  *queries.GetTrendingOffersHandler
	getCategoryHandler        *queries.GetCategoryHandler
	listCategoriesHandler     *queries.ListCategoriesHandler
	getCategoriesTreeHandler  *queries.GetCategoriesTreeHandler

	// Event publisher
	eventPublisher *nats.Publisher
}

// NewResolver creates a new resolver with all dependencies.
func NewResolver(
	offerRepo domain.OfferRepository,
	categoryRepo domain.CategoryRepository,
	searchRepo domain.OfferSearchService,
	eventPublisher *nats.Publisher,
) *Resolver {
	return &Resolver{
		offerRepo:      offerRepo,
		categoryRepo:   categoryRepo,
		searchRepo:     searchRepo,
		eventPublisher: eventPublisher,

		// Initialize command handlers
		createOfferHandler:    commands.NewCreateOfferHandler(offerRepo, categoryRepo, eventPublisher),
		publishOfferHandler:   commands.NewPublishOfferHandler(offerRepo, searchRepo, eventPublisher),
		expireOfferHandler:    commands.NewExpireOfferHandler(offerRepo, searchRepo, eventPublisher),
		archiveOfferHandler:   commands.NewArchiveOfferHandler(offerRepo, searchRepo, eventPublisher),
		extendOfferHandler:    commands.NewExtendOfferHandler(offerRepo, searchRepo, eventPublisher),
		createCategoryHandler: commands.NewCreateCategoryHandler(categoryRepo, eventPublisher),
		updateCategoryHandler: commands.NewUpdateCategoryHandler(categoryRepo, eventPublisher),
		deleteCategoryHandler: commands.NewDeleteCategoryHandler(categoryRepo, offerRepo, eventPublisher),

		// Initialize query handlers
		getOfferHandler:           queries.NewGetOfferHandler(offerRepo),
		listOffersHandler:         queries.NewListOffersHandler(offerRepo),
		searchOffersHandler:       queries.NewSearchOffersHandler(searchRepo),
		getOffersByPartnerHandler: queries.NewGetOffersByPartnerHandler(offerRepo),
		getNearbyOffersHandler:    queries.NewGetNearbyOffersHandler(searchRepo),
		getTrendingOffersHandler:  queries.NewGetTrendingOffersHandler(offerRepo),
		getCategoryHandler:        queries.NewGetCategoryHandler(categoryRepo),
		listCategoriesHandler:     queries.NewListCategoriesHandler(categoryRepo),
		getCategoriesTreeHandler:  queries.NewGetCategoriesTreeHandler(categoryRepo),
	}
}

// =============================================================================
// Query Resolvers
// =============================================================================

// Offer returns an offer by ID.
func (r *Resolver) Offer(ctx context.Context, id string) (*model.Offer, error) {
	offerID, err := domain.ParseOfferID(id)
	if err != nil {
		return nil, fmt.Errorf("invalid offer ID: %w", err)
	}

	result, err := r.getOfferHandler.Handle(ctx, queries.GetOfferQuery{ID: offerID})
	if err != nil {
		return nil, err
	}

	return mapOfferToModel(result), nil
}

// Offers returns a paginated list of offers.
func (r *Resolver) Offers(ctx context.Context, filter *model.OfferFilterInput) (*model.OfferListResult, error) {
	query := queries.ListOffersQuery{
		Offset: 0,
		Limit:  20,
	}

	if filter != nil {
		if filter.Offset != nil {
			query.Offset = *filter.Offset
		}
		if filter.Limit != nil {
			query.Limit = *filter.Limit
		}
		if filter.PartnerID != nil {
			partnerID, _ := domain.ParsePartnerID(*filter.PartnerID)
			query.PartnerID = &partnerID
		}
		if filter.CategoryID != nil {
			categoryID, _ := domain.ParseCategoryID(*filter.CategoryID)
			query.CategoryID = &categoryID
		}
		if filter.Status != nil {
			status := mapOfferStatusToDomain(*filter.Status)
			query.Status = &status
		}
		if filter.OnlyActive != nil {
			query.OnlyActive = *filter.OnlyActive
		}
	}

	result, err := r.listOffersHandler.Handle(ctx, query)
	if err != nil {
		return nil, err
	}

	offers := make([]*model.Offer, len(result.Offers))
	for i, offer := range result.Offers {
		offers[i] = mapOfferToModel(offer)
	}

	return &model.OfferListResult{
		Offers:  offers,
		Total:   int(result.Total),
		HasMore: result.HasMore,
	}, nil
}

// SearchOffers performs a full-text search on offers.
func (r *Resolver) SearchOffers(ctx context.Context, query string, filter *model.OfferFilterInput) (*model.OfferSearchResult, error) {
	searchQuery := queries.SearchOffersQuery{
		Query:  query,
		Offset: 0,
		Limit:  20,
	}

	if filter != nil {
		if filter.Offset != nil {
			searchQuery.Offset = *filter.Offset
		}
		if filter.Limit != nil {
			searchQuery.Limit = *filter.Limit
		}
		if filter.CategoryID != nil {
			categoryID, _ := domain.ParseCategoryID(*filter.CategoryID)
			searchQuery.CategoryID = &categoryID
		}
		if filter.Latitude != nil && filter.Longitude != nil {
			searchQuery.Latitude = filter.Latitude
			searchQuery.Longitude = filter.Longitude
		}
		if filter.RadiusKm != nil {
			searchQuery.RadiusKm = filter.RadiusKm
		}
	}

	result, err := r.searchOffersHandler.Handle(ctx, searchQuery)
	if err != nil {
		return nil, err
	}

	summaries := make([]*model.OfferSummary, len(result.Offers))
	for i, summary := range result.Offers {
		summaries[i] = mapOfferSummaryToModel(summary)
	}

	return &model.OfferSearchResult{
		Offers:  summaries,
		Total:   int(result.Total),
		HasMore: result.HasMore,
	}, nil
}

// NearbyOffers returns offers near a location.
func (r *Resolver) NearbyOffers(ctx context.Context, latitude, longitude, radiusKm float64, filter *model.OfferFilterInput) (*model.OfferSearchResult, error) {
	query := queries.GetNearbyOffersQuery{
		Latitude:  latitude,
		Longitude: longitude,
		RadiusKm:  radiusKm,
		Offset:    0,
		Limit:     20,
	}

	if filter != nil {
		if filter.Offset != nil {
			query.Offset = *filter.Offset
		}
		if filter.Limit != nil {
			query.Limit = *filter.Limit
		}
		if filter.CategoryID != nil {
			categoryID, _ := domain.ParseCategoryID(*filter.CategoryID)
			query.CategoryID = &categoryID
		}
	}

	result, err := r.getNearbyOffersHandler.Handle(ctx, query)
	if err != nil {
		return nil, err
	}

	summaries := make([]*model.OfferSummary, len(result.Offers))
	for i, summary := range result.Offers {
		summaries[i] = mapOfferSummaryToModel(summary)
	}

	return &model.OfferSearchResult{
		Offers:  summaries,
		Total:   int(result.Total),
		HasMore: result.HasMore,
	}, nil
}

// TrendingOffers returns trending offers.
func (r *Resolver) TrendingOffers(ctx context.Context, limit *int) ([]*model.TrendingOffer, error) {
	l := 10
	if limit != nil {
		l = *limit
	}

	result, err := r.getTrendingOffersHandler.Handle(ctx, queries.GetTrendingOffersQuery{Limit: l})
	if err != nil {
		return nil, err
	}

	trending := make([]*model.TrendingOffer, len(result))
	for i, t := range result {
		trending[i] = &model.TrendingOffer{
			Offer:  mapOfferToModel(t.Offer),
			Score:  t.Score,
			Reason: t.Reason,
		}
	}

	return trending, nil
}

// Category returns a category by ID.
func (r *Resolver) Category(ctx context.Context, id string) (*model.Category, error) {
	categoryID, err := domain.ParseCategoryID(id)
	if err != nil {
		return nil, fmt.Errorf("invalid category ID: %w", err)
	}

	result, err := r.getCategoryHandler.Handle(ctx, queries.GetCategoryQuery{ID: categoryID})
	if err != nil {
		return nil, err
	}

	return mapCategoryToModel(result), nil
}

// Categories returns all categories.
func (r *Resolver) Categories(ctx context.Context, activeOnly *bool) ([]*model.Category, error) {
	active := true
	if activeOnly != nil {
		active = *activeOnly
	}

	result, err := r.listCategoriesHandler.Handle(ctx, queries.ListCategoriesQuery{ActiveOnly: active})
	if err != nil {
		return nil, err
	}

	categories := make([]*model.Category, len(result))
	for i, cat := range result {
		categories[i] = mapCategoryToModel(cat)
	}

	return categories, nil
}

// CategoryTree returns the category tree.
func (r *Resolver) CategoryTree(ctx context.Context, activeOnly *bool) ([]*model.CategoryTree, error) {
	active := true
	if activeOnly != nil {
		active = *activeOnly
	}

	result, err := r.getCategoriesTreeHandler.Handle(ctx, queries.GetCategoriesTreeQuery{ActiveOnly: active})
	if err != nil {
		return nil, err
	}

	return mapCategoryTreeToModel(result), nil
}

// =============================================================================
// Mutation Resolvers
// =============================================================================

// CreateOffer creates a new offer.
func (r *Resolver) CreateOffer(ctx context.Context, input model.CreateOfferInput) (*model.Offer, error) {
	partnerID, err := domain.ParsePartnerID(input.PartnerID)
	if err != nil {
		return nil, fmt.Errorf("invalid partner ID: %w", err)
	}

	establishmentID, err := domain.ParseEstablishmentID(input.EstablishmentID)
	if err != nil {
		return nil, fmt.Errorf("invalid establishment ID: %w", err)
	}

	categoryID, err := domain.ParseCategoryID(input.CategoryID)
	if err != nil {
		return nil, fmt.Errorf("invalid category ID: %w", err)
	}

	cmd := commands.CreateOfferCommand{
		PartnerID:       partnerID,
		EstablishmentID: establishmentID,
		Title:           input.Title,
		Description:     input.Description,
		CategoryID:      categoryID,
		Tags:            input.Tags,
		DiscountType:    mapDiscountTypeToDomain(input.Discount.Type),
		DiscountValue:   input.Discount.Value,
		ValidityStart:   input.Validity.StartDate,
		ValidityEnd:     input.Validity.EndDate,
	}

	if input.ShortDescription != nil {
		cmd.ShortDescription = *input.ShortDescription
	}
	if input.TermsAndConditions != nil {
		cmd.TermsAndConditions = *input.TermsAndConditions
	}
	if input.Discount.Formula != nil {
		cmd.DiscountFormula = *input.Discount.Formula
	}
	if input.Validity.Timezone != nil {
		cmd.Timezone = *input.Validity.Timezone
	}

	offer, err := r.createOfferHandler.Handle(ctx, cmd)
	if err != nil {
		return nil, err
	}

	return mapOfferToModel(offer), nil
}

// PublishOffer publishes an offer.
func (r *Resolver) PublishOffer(ctx context.Context, id string) (*model.Offer, error) {
	offerID, err := domain.ParseOfferID(id)
	if err != nil {
		return nil, fmt.Errorf("invalid offer ID: %w", err)
	}

	offer, err := r.publishOfferHandler.Handle(ctx, commands.PublishOfferCommand{OfferID: offerID})
	if err != nil {
		return nil, err
	}

	return mapOfferToModel(offer), nil
}

// ArchiveOffer archives an offer.
func (r *Resolver) ArchiveOffer(ctx context.Context, id string) (*model.Offer, error) {
	offerID, err := domain.ParseOfferID(id)
	if err != nil {
		return nil, fmt.Errorf("invalid offer ID: %w", err)
	}

	offer, err := r.archiveOfferHandler.Handle(ctx, commands.ArchiveOfferCommand{OfferID: offerID})
	if err != nil {
		return nil, err
	}

	return mapOfferToModel(offer), nil
}

// ExtendOffer extends an offer's validity.
func (r *Resolver) ExtendOffer(ctx context.Context, id string, newEndDate time.Time) (*model.Offer, error) {
	offerID, err := domain.ParseOfferID(id)
	if err != nil {
		return nil, fmt.Errorf("invalid offer ID: %w", err)
	}

	offer, err := r.extendOfferHandler.Handle(ctx, commands.ExtendOfferCommand{
		OfferID:    offerID,
		NewEndDate: newEndDate,
	})
	if err != nil {
		return nil, err
	}

	return mapOfferToModel(offer), nil
}

// CreateCategory creates a new category.
func (r *Resolver) CreateCategory(ctx context.Context, input model.CreateCategoryInput) (*model.Category, error) {
	nameEN := ""
	if input.NameEn != nil {
		nameEN = *input.NameEn
	}

	cmd := commands.CreateCategoryCommand{
		Slug:   input.Slug,
		NameFR: input.NameFr,
		NameEN: nameEN,
		Icon:   ptrToString(input.Icon),
	}

	if input.DescriptionFr != nil {
		cmd.DescriptionFR = *input.DescriptionFr
	}
	if input.DescriptionEn != nil {
		cmd.DescriptionEN = *input.DescriptionEn
	}
	if input.Color != nil {
		cmd.Color = *input.Color
	}
	if input.Image != nil {
		cmd.Image = *input.Image
	}
	if input.ParentID != nil {
		parentID, err := domain.ParseCategoryID(*input.ParentID)
		if err == nil {
			cmd.ParentID = &parentID
		}
	}

	category, err := r.createCategoryHandler.Handle(ctx, cmd)
	if err != nil {
		return nil, err
	}

	return mapCategoryToModel(category), nil
}

// DeleteCategory deletes a category.
func (r *Resolver) DeleteCategory(ctx context.Context, id string) (bool, error) {
	categoryID, err := domain.ParseCategoryID(id)
	if err != nil {
		return false, fmt.Errorf("invalid category ID: %w", err)
	}

	err = r.deleteCategoryHandler.Handle(ctx, commands.DeleteCategoryCommand{CategoryID: categoryID})
	if err != nil {
		return false, err
	}

	return true, nil
}

// =============================================================================
// Entity Resolvers (Federation)
// =============================================================================

// FindOfferByID resolves an Offer entity by ID for federation.
func (r *Resolver) FindOfferByID(ctx context.Context, id string) (*model.Offer, error) {
	return r.Offer(ctx, id)
}

// FindCategoryByID resolves a Category entity by ID for federation.
func (r *Resolver) FindCategoryByID(ctx context.Context, id string) (*model.Category, error) {
	return r.Category(ctx, id)
}

// =============================================================================
// Mapper Functions
// =============================================================================

func mapOfferToModel(offer *domain.Offer) *model.Offer {
	if offer == nil {
		return nil
	}

	m := &model.Offer{
		ID:               offer.ID().String(),
		PartnerID:        offer.PartnerID().String(),
		EstablishmentID:  offer.EstablishmentID().String(),
		Title:            offer.Title(),
		Description:      offer.Description(),
		ShortDescription: offer.ShortDescription(),
		CategoryID:       offer.CategoryID().String(),
		Tags:             offer.Tags(),
		Status:           mapOfferStatusToModel(offer.Status()),
		IsActive:         offer.IsActive(),
		IsAvailableNow:   offer.IsAvailableNow(),
		CreatedAt:        offer.CreatedAt(),
		UpdatedAt:        offer.UpdatedAt(),
		PublishedAt:      offer.PublishedAt(),
	}

	// Map discount
	discount := offer.Discount()
	m.Discount = &model.Discount{
		Type:  mapDiscountTypeToModel(discount.Type),
		Value: discount.Value,
	}
	if discount.OriginalPrice != nil {
		m.Discount.OriginalPrice = &model.Money{
			Amount:   int(discount.OriginalPrice.Amount),
			Currency: discount.OriginalPrice.Currency,
		}
	}
	if discount.DiscountedPrice != nil {
		m.Discount.DiscountedPrice = &model.Money{
			Amount:   int(discount.DiscountedPrice.Amount),
			Currency: discount.DiscountedPrice.Currency,
		}
	}
	if discount.Formula != "" {
		m.Discount.Formula = &discount.Formula
	}

	// Map validity
	validity := offer.Validity()
	m.Validity = &model.Validity{
		StartDate: validity.StartDate,
		EndDate:   validity.EndDate,
		Timezone:  validity.Timezone,
	}

	// Map schedule
	schedule := offer.Schedule()
	m.Schedule = &model.Schedule{
		AllDay: schedule.AllDay,
		Slots:  make([]*model.TimeSlot, len(schedule.Slots)),
	}
	for i, slot := range schedule.Slots {
		m.Schedule.Slots[i] = &model.TimeSlot{
			DayOfWeek: slot.DayOfWeek,
			StartTime: slot.StartTime,
			EndTime:   slot.EndTime,
		}
	}

	// Map quota
	quota := offer.Quota()
	m.Quota = &model.Quota{
		Used: quota.Used,
	}
	if quota.Total > 0 {
		m.Quota.Total = &quota.Total
		remaining := quota.Total - quota.Used
		m.Quota.Remaining = &remaining
		m.RemainingQuota = &remaining
	}
	if quota.PerUser > 0 {
		m.Quota.PerUser = &quota.PerUser
	}
	if quota.PerDay > 0 {
		m.Quota.PerDay = &quota.PerDay
	}

	// Map images
	images := offer.Images()
	m.Images = make([]*model.OfferImage, len(images))
	for i, img := range images {
		m.Images[i] = &model.OfferImage{
			URL:       img.URL,
			IsPrimary: img.IsPrimary,
			Order:     img.Order,
		}
		if img.Alt != "" {
			m.Images[i].Alt = &img.Alt
		}
	}

	// Map partner snapshot
	partner := offer.PartnerSnapshot()
	m.Partner = &model.PartnerSnapshot{
		ID:       partner.ID.String(),
		Name:     partner.Name,
		Category: partner.Category,
	}
	if partner.Logo != "" {
		m.Partner.Logo = &partner.Logo
	}

	// Map establishment snapshot
	establishment := offer.EstablishmentSnapshot()
	m.Establishment = &model.EstablishmentSnapshot{
		ID:      establishment.ID.String(),
		Name:    establishment.Name,
		Address: establishment.Address,
		City:    establishment.City,
		Location: &model.GeoLocation{
			Latitude:  establishment.Location.Latitude,
			Longitude: establishment.Location.Longitude,
		},
	}

	// Map stats
	stats := offer.Stats()
	m.Stats = &model.OfferStats{
		Views:       int(stats.Views),
		Clicks:      int(stats.Clicks),
		Bookings:    int(stats.Bookings),
		Checkins:    int(stats.Checkins),
		Favorites:   int(stats.Favorites),
		AvgRating:   stats.AvgRating,
		ReviewCount: stats.ReviewCount,
	}

	// Map moderation
	moderation := offer.Moderation()
	m.Moderation = &model.Moderation{
		Status: mapModerationStatusToModel(moderation.Status),
	}
	if moderation.ReviewedBy != "" {
		m.Moderation.ReviewedBy = &moderation.ReviewedBy
	}
	if moderation.ReviewedAt != nil {
		m.Moderation.ReviewedAt = moderation.ReviewedAt
	}
	if moderation.Comment != "" {
		m.Moderation.Comment = &moderation.Comment
	}

	// Map conditions
	conditions := offer.Conditions()
	m.Conditions = make([]*model.Condition, len(conditions))
	for i, cond := range conditions {
		m.Conditions[i] = &model.Condition{
			Type:  mapConditionTypeToModel(cond.Type),
			Value: cond.Value,
			Label: cond.Label,
		}
	}

	if offer.TermsAndConditions() != "" {
		tc := offer.TermsAndConditions()
		m.TermsAndConditions = &tc
	}

	return m
}

func mapCategoryToModel(cat *domain.Category) *model.Category {
	if cat == nil {
		return nil
	}

	m := &model.Category{
		ID:        cat.ID().String(),
		Slug:      cat.Slug(),
		Order:     cat.Order(),
		IsActive:  cat.IsActive(),
		CreatedAt: cat.CreatedAt(),
		UpdatedAt: cat.UpdatedAt(),
	}

	// Map name
	name := cat.Name()
	m.Name = &model.LocalizedString{
		FR: name.FR,
	}
	if name.EN != "" {
		m.Name.EN = &name.EN
	}

	// Map description
	desc := cat.Description()
	if desc.FR != "" || desc.EN != "" {
		m.Description = &model.LocalizedString{
			FR: desc.FR,
		}
		if desc.EN != "" {
			m.Description.EN = &desc.EN
		}
	}

	// Map optional fields
	if cat.Icon() != "" {
		icon := cat.Icon()
		m.Icon = &icon
	}
	if cat.Color() != "" {
		color := cat.Color()
		m.Color = &color
	}
	if cat.Image() != "" {
		image := cat.Image()
		m.Image = &image
	}
	if cat.ParentID() != nil {
		parentID := cat.ParentID().String()
		m.ParentID = &parentID
	}

	return m
}

func mapCategoryTreeToModel(trees []*domain.CategoryTree) []*model.CategoryTree {
	result := make([]*model.CategoryTree, len(trees))
	for i, tree := range trees {
		result[i] = &model.CategoryTree{
			Category: mapCategoryToModel(tree.Category),
			Children: mapCategoryTreeToModel(tree.Children),
		}
	}
	return result
}

func mapOfferSummaryToModel(summary *domain.OfferSummary) *model.OfferSummary {
	if summary == nil {
		return nil
	}

	return &model.OfferSummary{
		ID:                summary.ID.String(),
		PartnerID:         summary.PartnerID.String(),
		EstablishmentID:   summary.EstablishmentID.String(),
		Title:             summary.Title,
		ShortDescription:  summary.ShortDescription,
		CategoryID:        summary.CategoryID.String(),
		DiscountType:      mapDiscountTypeToModel(summary.DiscountType),
		DiscountValue:     summary.DiscountValue,
		Status:            mapOfferStatusToModel(summary.Status),
		PartnerName:       summary.PartnerName,
		EstablishmentName: summary.EstablishmentName,
		EstablishmentCity: summary.EstablishmentCity,
		Location: &model.GeoLocation{
			Latitude:  summary.Location.Latitude,
			Longitude: summary.Location.Longitude,
		},
		Views:       int(summary.Views),
		AvgRating:   summary.AvgRating,
		ReviewCount: summary.ReviewCount,
		PublishedAt: summary.PublishedAt,
	}
}

func mapOfferStatusToModel(status domain.OfferStatus) model.OfferStatus {
	switch status {
	case domain.OfferStatusDraft:
		return model.OfferStatusDraft
	case domain.OfferStatusPending:
		return model.OfferStatusPending
	case domain.OfferStatusActive:
		return model.OfferStatusActive
	case domain.OfferStatusPaused:
		return model.OfferStatusPaused
	case domain.OfferStatusExpired:
		return model.OfferStatusExpired
	case domain.OfferStatusArchived:
		return model.OfferStatusArchived
	default:
		return model.OfferStatusDraft
	}
}

func mapOfferStatusToDomain(status model.OfferStatus) domain.OfferStatus {
	switch status {
	case model.OfferStatusDraft:
		return domain.OfferStatusDraft
	case model.OfferStatusPending:
		return domain.OfferStatusPending
	case model.OfferStatusActive:
		return domain.OfferStatusActive
	case model.OfferStatusPaused:
		return domain.OfferStatusPaused
	case model.OfferStatusExpired:
		return domain.OfferStatusExpired
	case model.OfferStatusArchived:
		return domain.OfferStatusArchived
	default:
		return domain.OfferStatusDraft
	}
}

func mapDiscountTypeToModel(dt domain.DiscountType) model.DiscountType {
	switch dt {
	case domain.DiscountTypePercentage:
		return model.DiscountTypePercentage
	case domain.DiscountTypeFixed:
		return model.DiscountTypeFixed
	case domain.DiscountTypeFormula:
		return model.DiscountTypeFormula
	default:
		return model.DiscountTypePercentage
	}
}

func mapDiscountTypeToDomain(dt model.DiscountType) domain.DiscountType {
	switch dt {
	case model.DiscountTypePercentage:
		return domain.DiscountTypePercentage
	case model.DiscountTypeFixed:
		return domain.DiscountTypeFixed
	case model.DiscountTypeFormula:
		return domain.DiscountTypeFormula
	default:
		return domain.DiscountTypePercentage
	}
}

func mapModerationStatusToModel(status domain.ModerationStatus) model.ModerationStatus {
	switch status {
	case domain.ModerationStatusPending:
		return model.ModerationStatusPending
	case domain.ModerationStatusApproved:
		return model.ModerationStatusApproved
	case domain.ModerationStatusRejected:
		return model.ModerationStatusRejected
	default:
		return model.ModerationStatusPending
	}
}

func mapConditionTypeToModel(ct domain.ConditionType) model.ConditionType {
	switch ct {
	case domain.ConditionTypeMinPurchase:
		return model.ConditionTypeMinPurchase
	case domain.ConditionTypeMinPeople:
		return model.ConditionTypeMinPeople
	case domain.ConditionTypeFirstVisit:
		return model.ConditionTypeFirstVisit
	case domain.ConditionTypeSpecificDays:
		return model.ConditionTypeSpecificDays
	case domain.ConditionTypeSpecificHours:
		return model.ConditionTypeSpecificHours
	default:
		return model.ConditionTypeOther
	}
}

func ptrToString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
