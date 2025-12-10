// Package mongodb contains MongoDB repository implementations for the Discovery service.
package mongodb

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/yousoon/discovery-service/internal/domain"
)

// OfferDocument represents the MongoDB document structure for offers.
type OfferDocument struct {
	ID              primitive.ObjectID `bson:"_id,omitempty"`
	PartnerID       string             `bson:"partner_id"`
	EstablishmentID string             `bson:"establishment_id"`

	Title            string   `bson:"title"`
	Description      string   `bson:"description"`
	ShortDescription string   `bson:"short_description"`
	CategoryID       string   `bson:"category_id"`
	Tags             []string `bson:"tags"`

	Discount           DiscountDoc     `bson:"discount"`
	Conditions         []ConditionDoc  `bson:"conditions"`
	TermsAndConditions string          `bson:"terms_and_conditions"`
	Validity           ValidityDoc     `bson:"validity"`
	Schedule           ScheduleDoc     `bson:"schedule"`
	Quota              QuotaDoc        `bson:"quota"`
	Images             []OfferImageDoc `bson:"images"`

	PartnerSnapshot       PartnerSnapshotDoc       `bson:"_partner"`
	EstablishmentSnapshot EstablishmentSnapshotDoc `bson:"_establishment"`

	Stats      OfferStatsDoc `bson:"stats"`
	Status     string        `bson:"status"`
	Moderation ModerationDoc `bson:"moderation"`
	IsActive   bool          `bson:"is_active"`

	CreatedAt   time.Time  `bson:"created_at"`
	UpdatedAt   time.Time  `bson:"updated_at"`
	PublishedAt *time.Time `bson:"published_at"`
	DeletedAt   *time.Time `bson:"deleted_at"`
}

// Subdocuments
type DiscountDoc struct {
	Type          string `bson:"type"`
	Value         int    `bson:"value"`
	OriginalPrice *int64 `bson:"original_price"`
	Formula       string `bson:"formula"`
}

type ConditionDoc struct {
	Type  string      `bson:"type"`
	Value interface{} `bson:"value"`
	Label string      `bson:"label"`
}

type ValidityDoc struct {
	StartDate time.Time `bson:"start_date"`
	EndDate   time.Time `bson:"end_date"`
	Timezone  string    `bson:"timezone"`
}

type ScheduleDoc struct {
	AllDay bool          `bson:"all_day"`
	Slots  []TimeSlotDoc `bson:"slots"`
}

type TimeSlotDoc struct {
	DayOfWeek int    `bson:"day_of_week"`
	StartTime string `bson:"start_time"`
	EndTime   string `bson:"end_time"`
}

type QuotaDoc struct {
	Total   *int `bson:"total"`
	PerUser *int `bson:"per_user"`
	PerDay  *int `bson:"per_day"`
	Used    int  `bson:"used"`
}

type OfferImageDoc struct {
	URL       string `bson:"url"`
	Alt       string `bson:"alt"`
	IsPrimary bool   `bson:"is_primary"`
	Order     int    `bson:"order"`
}

type PartnerSnapshotDoc struct {
	Name     string `bson:"name"`
	Logo     string `bson:"logo"`
	Category string `bson:"category"`
}

type EstablishmentSnapshotDoc struct {
	Name     string         `bson:"name"`
	Address  string         `bson:"address"`
	City     string         `bson:"city"`
	Location GeoLocationDoc `bson:"location"`
}

type GeoLocationDoc struct {
	Type        string    `bson:"type"`
	Coordinates []float64 `bson:"coordinates"`
}

type OfferStatsDoc struct {
	Views     int `bson:"views"`
	Bookings  int `bson:"bookings"`
	Checkins  int `bson:"checkins"`
	Favorites int `bson:"favorites"`
}

type ModerationDoc struct {
	Status     string     `bson:"status"`
	ReviewerID *string    `bson:"reviewer_id"`
	ReviewedAt *time.Time `bson:"reviewed_at"`
	Comment    *string    `bson:"comment"`
}

// OfferRepository implements domain.OfferRepository using MongoDB.
type OfferRepository struct {
	collection *mongo.Collection
}

// NewOfferRepository creates a new OfferRepository.
func NewOfferRepository(db *mongo.Database) *OfferRepository {
	return &OfferRepository{
		collection: db.Collection("offers"),
	}
}

// EnsureIndexes creates necessary indexes for the offer collection.
func (r *OfferRepository) EnsureIndexes(ctx context.Context) error {
	indexes := []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "partner_id", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "establishment_id", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "category_id", Value: 1}},
		},
		{
			Keys: bson.D{
				{Key: "status", Value: 1},
				{Key: "is_active", Value: 1},
			},
		},
		{
			Keys: bson.D{{Key: "_establishment.location", Value: "2dsphere"}},
		},
		{
			Keys:    bson.D{{Key: "title", Value: "text"}, {Key: "description", Value: "text"}},
			Options: options.Index().SetDefaultLanguage("french"),
		},
		{
			Keys: bson.D{
				{Key: "validity.start_date", Value: 1},
				{Key: "validity.end_date", Value: 1},
			},
		},
		{
			Keys: bson.D{{Key: "created_at", Value: -1}},
		},
		{
			Keys: bson.D{{Key: "deleted_at", Value: 1}},
		},
	}

	_, err := r.collection.Indexes().CreateMany(ctx, indexes)
	return err
}

// Save persists an offer (create or update).
func (r *OfferRepository) Save(ctx context.Context, offer *domain.Offer) error {
	doc := r.toDocument(offer)

	if offer.ID() == "" {
		return errors.New("offer ID is required")
	}

	filter := bson.M{"_id": doc.ID}
	update := bson.M{"$set": doc}
	opts := options.Update().SetUpsert(true)

	_, err := r.collection.UpdateOne(ctx, filter, update, opts)
	return err
}

// FindByID retrieves an offer by ID.
func (r *OfferRepository) FindByID(ctx context.Context, id domain.OfferID) (*domain.Offer, error) {
	objectID, err := primitive.ObjectIDFromHex(string(id))
	if err != nil {
		return nil, nil // Invalid ID format, return nil
	}

	filter := bson.M{
		"_id":        objectID,
		"deleted_at": nil,
	}

	var doc OfferDocument
	err = r.collection.FindOne(ctx, filter).Decode(&doc)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}

	return r.toDomain(&doc), nil
}

// FindByPartnerID retrieves all offers for a partner.
func (r *OfferRepository) FindByPartnerID(ctx context.Context, partnerID domain.PartnerID) ([]*domain.Offer, error) {
	filter := bson.M{
		"partner_id": string(partnerID),
		"deleted_at": nil,
	}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	return r.cursorToOffers(ctx, cursor)
}

// FindByEstablishmentID retrieves all offers for an establishment.
func (r *OfferRepository) FindByEstablishmentID(ctx context.Context, establishmentID domain.EstablishmentID) ([]*domain.Offer, error) {
	filter := bson.M{
		"establishment_id": string(establishmentID),
		"deleted_at":       nil,
	}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	return r.cursorToOffers(ctx, cursor)
}

// FindByCategory retrieves offers in a category.
func (r *OfferRepository) FindByCategory(ctx context.Context, categoryID domain.CategoryID, offset, limit int) ([]*domain.Offer, error) {
	filter := bson.M{
		"category_id": string(categoryID),
		"status":      "active",
		"deleted_at":  nil,
	}

	opts := options.Find().
		SetSkip(int64(offset)).
		SetLimit(int64(limit)).
		SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	return r.cursorToOffers(ctx, cursor)
}

// List retrieves offers with filters.
func (r *OfferRepository) List(ctx context.Context, filter domain.OfferFilter) (*domain.OfferListResult, error) {
	mongoFilter := r.buildFilter(filter)

	// Count total
	total, err := r.collection.CountDocuments(ctx, mongoFilter)
	if err != nil {
		return nil, err
	}

	// Build options
	opts := options.Find()
	if filter.Limit > 0 {
		opts.SetLimit(int64(filter.Limit))
	}
	if filter.Offset > 0 {
		opts.SetSkip(int64(filter.Offset))
	}

	// Sorting
	sort := r.buildSort(filter)
	if sort != nil {
		opts.SetSort(sort)
	}

	// For geo queries, use aggregation
	var cursor *mongo.Cursor
	if filter.Location != nil {
		cursor, err = r.findWithLocation(ctx, filter, mongoFilter, opts)
	} else {
		cursor, err = r.collection.Find(ctx, mongoFilter, opts)
	}

	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	offers, err := r.cursorToOffers(ctx, cursor)
	if err != nil {
		return nil, err
	}

	return &domain.OfferListResult{
		Offers:     offers,
		TotalCount: total,
		Offset:     filter.Offset,
		Limit:      filter.Limit,
	}, nil
}

// Delete soft-deletes an offer.
func (r *OfferRepository) Delete(ctx context.Context, id domain.OfferID) error {
	objectID, err := primitive.ObjectIDFromHex(string(id))
	if err != nil {
		return err
	}

	update := bson.M{
		"$set": bson.M{
			"deleted_at": time.Now(),
			"updated_at": time.Now(),
		},
	}

	_, err = r.collection.UpdateByID(ctx, objectID, update)
	return err
}

// Count returns the total count of offers matching the filter.
func (r *OfferRepository) Count(ctx context.Context, filter domain.OfferFilter) (int64, error) {
	mongoFilter := r.buildFilter(filter)
	return r.collection.CountDocuments(ctx, mongoFilter)
}

// ExistsActiveForPartner checks if a partner has any active offers.
func (r *OfferRepository) ExistsActiveForPartner(ctx context.Context, partnerID domain.PartnerID) (bool, error) {
	filter := bson.M{
		"partner_id": string(partnerID),
		"status":     "active",
		"deleted_at": nil,
	}

	count, err := r.collection.CountDocuments(ctx, filter, options.Count().SetLimit(1))
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// =============================================================================
// Helper Methods
// =============================================================================

func (r *OfferRepository) buildFilter(filter domain.OfferFilter) bson.M {
	mongoFilter := bson.M{
		"deleted_at": nil,
	}

	if filter.PartnerID != nil {
		mongoFilter["partner_id"] = string(*filter.PartnerID)
	}

	if filter.EstablishmentID != nil {
		mongoFilter["establishment_id"] = string(*filter.EstablishmentID)
	}

	if filter.CategoryID != nil {
		mongoFilter["category_id"] = string(*filter.CategoryID)
	}

	if filter.Status != nil {
		mongoFilter["status"] = string(*filter.Status)
	}

	if filter.OnlyActive {
		mongoFilter["status"] = "active"
		mongoFilter["is_active"] = true
	}

	if len(filter.Tags) > 0 {
		mongoFilter["tags"] = bson.M{"$in": filter.Tags}
	}

	if filter.SearchQuery != "" {
		mongoFilter["$text"] = bson.M{"$search": filter.SearchQuery}
	}

	if filter.ModerationStatus != nil {
		mongoFilter["moderation.status"] = string(*filter.ModerationStatus)
	}

	return mongoFilter
}

func (r *OfferRepository) buildSort(filter domain.OfferFilter) bson.D {
	order := -1
	if filter.SortOrder == "asc" {
		order = 1
	}

	switch filter.SortBy {
	case "distance":
		return nil // Handled by geo query
	case "discount":
		return bson.D{{Key: "discount.value", Value: order}}
	case "popularity":
		return bson.D{{Key: "stats.views", Value: order}}
	default:
		return bson.D{{Key: "created_at", Value: order}}
	}
}

func (r *OfferRepository) findWithLocation(ctx context.Context, filter domain.OfferFilter, mongoFilter bson.M, opts *options.FindOptions) (*mongo.Cursor, error) {
	pipeline := mongo.Pipeline{
		{{Key: "$geoNear", Value: bson.M{
			"near":          filter.Location,
			"distanceField": "distance",
			"maxDistance":   filter.RadiusKm * 1000, // km to meters
			"spherical":     true,
			"query":         mongoFilter,
		}}},
	}

	if filter.Offset > 0 {
		pipeline = append(pipeline, bson.D{{Key: "$skip", Value: filter.Offset}})
	}

	if filter.Limit > 0 {
		pipeline = append(pipeline, bson.D{{Key: "$limit", Value: filter.Limit}})
	}

	return r.collection.Aggregate(ctx, pipeline)
}

func (r *OfferRepository) cursorToOffers(ctx context.Context, cursor *mongo.Cursor) ([]*domain.Offer, error) {
	var offers []*domain.Offer

	for cursor.Next(ctx) {
		var doc OfferDocument
		if err := cursor.Decode(&doc); err != nil {
			return nil, err
		}
		offers = append(offers, r.toDomain(&doc))
	}

	return offers, cursor.Err()
}

// =============================================================================
// Mapping Methods
// =============================================================================

func (r *OfferRepository) toDocument(offer *domain.Offer) *OfferDocument {
	objectID, _ := primitive.ObjectIDFromHex(string(offer.ID()))

	// Map images
	images := make([]OfferImageDoc, len(offer.Images()))
	for i, img := range offer.Images() {
		images[i] = OfferImageDoc{
			URL:       img.URL,
			Alt:       img.Alt,
			IsPrimary: img.IsPrimary,
			Order:     img.Order,
		}
	}

	// Map conditions
	conditions := make([]ConditionDoc, len(offer.Conditions()))
	for i, c := range offer.Conditions() {
		conditions[i] = ConditionDoc{
			Type:  string(c.Type),
			Value: c.Value,
			Label: c.Label,
		}
	}

	// Map time slots
	slots := make([]TimeSlotDoc, len(offer.Schedule().Slots))
	for i, s := range offer.Schedule().Slots {
		slots[i] = TimeSlotDoc{
			DayOfWeek: s.DayOfWeek,
			StartTime: s.StartTime,
			EndTime:   s.EndTime,
		}
	}

	return &OfferDocument{
		ID:               objectID,
		PartnerID:        string(offer.PartnerID()),
		EstablishmentID:  string(offer.EstablishmentID()),
		Title:            offer.Title(),
		Description:      offer.Description(),
		ShortDescription: offer.ShortDescription(),
		CategoryID:       string(offer.CategoryID()),
		Tags:             offer.Tags(),
		Discount: DiscountDoc{
			Type:          string(offer.Discount().Type),
			Value:         offer.Discount().Value,
			OriginalPrice: offer.Discount().OriginalPrice,
			Formula:       offer.Discount().Formula,
		},
		Conditions:         conditions,
		TermsAndConditions: offer.TermsAndConditions(),
		Validity: ValidityDoc{
			StartDate: offer.Validity().StartDate,
			EndDate:   offer.Validity().EndDate,
			Timezone:  offer.Validity().Timezone,
		},
		Schedule: ScheduleDoc{
			AllDay: offer.Schedule().AllDay,
			Slots:  slots,
		},
		Quota: QuotaDoc{
			Total:   offer.Quota().Total,
			PerUser: offer.Quota().PerUser,
			PerDay:  offer.Quota().PerDay,
			Used:    offer.Quota().Used,
		},
		Images: images,
		PartnerSnapshot: PartnerSnapshotDoc{
			Name:     offer.PartnerSnapshot().Name,
			Logo:     offer.PartnerSnapshot().Logo,
			Category: offer.PartnerSnapshot().Category,
		},
		EstablishmentSnapshot: EstablishmentSnapshotDoc{
			Name:    offer.EstablishmentSnapshot().Name,
			Address: offer.EstablishmentSnapshot().Address,
			City:    offer.EstablishmentSnapshot().City,
			Location: GeoLocationDoc{
				Type:        offer.EstablishmentSnapshot().Location.Type,
				Coordinates: offer.EstablishmentSnapshot().Location.Coordinates,
			},
		},
		Stats: OfferStatsDoc{
			Views:     offer.Stats().Views,
			Bookings:  offer.Stats().Bookings,
			Checkins:  offer.Stats().Checkins,
			Favorites: offer.Stats().Favorites,
		},
		Status: string(offer.Status()),
		Moderation: ModerationDoc{
			Status:     string(offer.Moderation().Status),
			ReviewerID: offer.Moderation().ReviewerID,
			ReviewedAt: offer.Moderation().ReviewedAt,
			Comment:    offer.Moderation().Comment,
		},
		IsActive:    offer.IsActive(),
		CreatedAt:   offer.CreatedAt(),
		UpdatedAt:   offer.UpdatedAt(),
		PublishedAt: offer.PublishedAt(),
		DeletedAt:   offer.DeletedAt(),
	}
}

func (r *OfferRepository) toDomain(doc *OfferDocument) *domain.Offer {
	// Map images
	images := make([]domain.OfferImage, len(doc.Images))
	for i, img := range doc.Images {
		images[i] = domain.OfferImage{
			URL:       img.URL,
			Alt:       img.Alt,
			IsPrimary: img.IsPrimary,
			Order:     img.Order,
		}
	}

	// Map conditions
	conditions := make([]domain.Condition, len(doc.Conditions))
	for i, c := range doc.Conditions {
		conditions[i] = domain.Condition{
			Type:  domain.ConditionType(c.Type),
			Value: c.Value,
			Label: c.Label,
		}
	}

	// Map time slots
	slots := make([]domain.TimeSlot, len(doc.Schedule.Slots))
	for i, s := range doc.Schedule.Slots {
		slots[i] = domain.TimeSlot{
			DayOfWeek: s.DayOfWeek,
			StartTime: s.StartTime,
			EndTime:   s.EndTime,
		}
	}

	return domain.ReconstructOffer(
		domain.OfferID(doc.ID.Hex()),
		domain.PartnerID(doc.PartnerID),
		domain.EstablishmentID(doc.EstablishmentID),
		doc.Title,
		doc.Description,
		doc.ShortDescription,
		domain.CategoryID(doc.CategoryID),
		doc.Tags,
		domain.Discount{
			Type:          domain.DiscountType(doc.Discount.Type),
			Value:         doc.Discount.Value,
			OriginalPrice: doc.Discount.OriginalPrice,
			Formula:       doc.Discount.Formula,
		},
		conditions,
		doc.TermsAndConditions,
		domain.Validity{
			StartDate: doc.Validity.StartDate,
			EndDate:   doc.Validity.EndDate,
			Timezone:  doc.Validity.Timezone,
		},
		domain.Schedule{
			AllDay: doc.Schedule.AllDay,
			Slots:  slots,
		},
		domain.Quota{
			Total:   doc.Quota.Total,
			PerUser: doc.Quota.PerUser,
			PerDay:  doc.Quota.PerDay,
			Used:    doc.Quota.Used,
		},
		images,
		domain.PartnerSnapshot{
			Name:     doc.PartnerSnapshot.Name,
			Logo:     doc.PartnerSnapshot.Logo,
			Category: doc.PartnerSnapshot.Category,
		},
		domain.EstablishmentSnapshot{
			Name:    doc.EstablishmentSnapshot.Name,
			Address: doc.EstablishmentSnapshot.Address,
			City:    doc.EstablishmentSnapshot.City,
			Location: domain.GeoLocation{
				Type:        doc.EstablishmentSnapshot.Location.Type,
				Coordinates: doc.EstablishmentSnapshot.Location.Coordinates,
			},
		},
		domain.OfferStats{
			Views:     doc.Stats.Views,
			Bookings:  doc.Stats.Bookings,
			Checkins:  doc.Stats.Checkins,
			Favorites: doc.Stats.Favorites,
		},
		domain.OfferStatus(doc.Status),
		domain.Moderation{
			Status:     domain.ModerationStatus(doc.Moderation.Status),
			ReviewerID: doc.Moderation.ReviewerID,
			ReviewedAt: doc.Moderation.ReviewedAt,
			Comment:    doc.Moderation.Comment,
		},
		doc.CreatedAt,
		doc.UpdatedAt,
		doc.PublishedAt,
		doc.DeletedAt,
	)
}

// =============================================================================
// OfferReadRepository Implementation
// =============================================================================

// GetOfferSummaries returns offer summaries for lists.
func (r *OfferRepository) GetOfferSummaries(ctx context.Context, filter domain.OfferFilter) ([]domain.OfferSummary, int64, error) {
	result, err := r.List(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	summaries := make([]domain.OfferSummary, len(result.Offers))
	for i, offer := range result.Offers {
		summaries[i] = r.toSummary(offer)
	}

	return summaries, result.TotalCount, nil
}

// GetOffersNearLocation returns offers near a location.
func (r *OfferRepository) GetOffersNearLocation(ctx context.Context, location domain.GeoLocation, radiusKm float64, limit int) ([]domain.OfferSummary, error) {
	filter := domain.OfferFilter{
		Location:   &location,
		RadiusKm:   radiusKm,
		Limit:      limit,
		OnlyActive: true,
	}

	result, err := r.List(ctx, filter)
	if err != nil {
		return nil, err
	}

	summaries := make([]domain.OfferSummary, len(result.Offers))
	for i, offer := range result.Offers {
		summaries[i] = r.toSummary(offer)
	}

	return summaries, nil
}

// GetOffersByIDs returns offers by IDs.
func (r *OfferRepository) GetOffersByIDs(ctx context.Context, ids []domain.OfferID) ([]*domain.Offer, error) {
	objectIDs := make([]primitive.ObjectID, 0, len(ids))
	for _, id := range ids {
		oid, err := primitive.ObjectIDFromHex(string(id))
		if err != nil {
			continue
		}
		objectIDs = append(objectIDs, oid)
	}

	filter := bson.M{
		"_id":        bson.M{"$in": objectIDs},
		"deleted_at": nil,
	}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	return r.cursorToOffers(ctx, cursor)
}

// SearchOffers performs a full-text search on offers.
func (r *OfferRepository) SearchOffers(ctx context.Context, query string, location *domain.GeoLocation, limit int) ([]domain.OfferSummary, error) {
	filter := domain.OfferFilter{
		SearchQuery: query,
		Location:    location,
		Limit:       limit,
		OnlyActive:  true,
	}

	result, err := r.List(ctx, filter)
	if err != nil {
		return nil, err
	}

	summaries := make([]domain.OfferSummary, len(result.Offers))
	for i, offer := range result.Offers {
		summaries[i] = r.toSummary(offer)
	}

	return summaries, nil
}

// GetRecommendedOffers returns personalized offers for a user.
func (r *OfferRepository) GetRecommendedOffers(ctx context.Context, userID domain.UserID, userCategories []domain.CategoryID, location domain.GeoLocation, limit int) ([]domain.OfferSummary, error) {
	// Simple implementation: get offers in user's preferred categories near their location
	filter := domain.OfferFilter{
		Location:   &location,
		RadiusKm:   10, // Default 10km
		Limit:      limit,
		OnlyActive: true,
	}

	// If user has category preferences, filter by first one (simplified)
	if len(userCategories) > 0 {
		filter.CategoryID = &userCategories[0]
	}

	result, err := r.List(ctx, filter)
	if err != nil {
		return nil, err
	}

	summaries := make([]domain.OfferSummary, len(result.Offers))
	for i, offer := range result.Offers {
		summaries[i] = r.toSummary(offer)
	}

	return summaries, nil
}

// GetTrendingOffers returns trending offers (most booked).
func (r *OfferRepository) GetTrendingOffers(ctx context.Context, location *domain.GeoLocation, limit int) ([]domain.OfferSummary, error) {
	filter := domain.OfferFilter{
		Location:   location,
		RadiusKm:   50, // Wider radius for trending
		Limit:      limit,
		OnlyActive: true,
		SortBy:     "popularity",
		SortOrder:  "desc",
	}

	result, err := r.List(ctx, filter)
	if err != nil {
		return nil, err
	}

	summaries := make([]domain.OfferSummary, len(result.Offers))
	for i, offer := range result.Offers {
		summaries[i] = r.toSummary(offer)
	}

	return summaries, nil
}

// GetNewOffers returns recently published offers.
func (r *OfferRepository) GetNewOffers(ctx context.Context, location *domain.GeoLocation, limit int) ([]domain.OfferSummary, error) {
	filter := domain.OfferFilter{
		Location:   location,
		RadiusKm:   50,
		Limit:      limit,
		OnlyActive: true,
		SortBy:     "created_at",
		SortOrder:  "desc",
	}

	result, err := r.List(ctx, filter)
	if err != nil {
		return nil, err
	}

	summaries := make([]domain.OfferSummary, len(result.Offers))
	for i, offer := range result.Offers {
		summaries[i] = r.toSummary(offer)
	}

	return summaries, nil
}

// GetExpiringOffers returns offers expiring soon.
func (r *OfferRepository) GetExpiringOffers(ctx context.Context, withinDays int, limit int) ([]domain.OfferSummary, error) {
	now := time.Now()
	expiryDate := now.AddDate(0, 0, withinDays)

	mongoFilter := bson.M{
		"status":     "active",
		"is_active":  true,
		"deleted_at": nil,
		"validity.end_date": bson.M{
			"$gte": now,
			"$lte": expiryDate,
		},
	}

	opts := options.Find().
		SetLimit(int64(limit)).
		SetSort(bson.D{{Key: "validity.end_date", Value: 1}})

	cursor, err := r.collection.Find(ctx, mongoFilter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	offers, err := r.cursorToOffers(ctx, cursor)
	if err != nil {
		return nil, err
	}

	summaries := make([]domain.OfferSummary, len(offers))
	for i, offer := range offers {
		summaries[i] = r.toSummary(offer)
	}

	return summaries, nil
}

// GetOfferCountByCategory returns offer counts per category.
func (r *OfferRepository) GetOfferCountByCategory(ctx context.Context) (map[domain.CategoryID]int, error) {
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{
			"status":     "active",
			"is_active":  true,
			"deleted_at": nil,
		}}},
		{{Key: "$group", Value: bson.M{
			"_id":   "$category_id",
			"count": bson.M{"$sum": 1},
		}}},
	}

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	result := make(map[domain.CategoryID]int)
	for cursor.Next(ctx) {
		var doc struct {
			ID    string `bson:"_id"`
			Count int    `bson:"count"`
		}
		if err := cursor.Decode(&doc); err != nil {
			continue
		}
		result[domain.CategoryID(doc.ID)] = doc.Count
	}

	return result, cursor.Err()
}

// GetPartnerOfferStats returns statistics for a partner's offers.
func (r *OfferRepository) GetPartnerOfferStats(ctx context.Context, partnerID domain.PartnerID) (*domain.PartnerOfferStats, error) {
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{
			"partner_id": string(partnerID),
			"deleted_at": nil,
		}}},
		{{Key: "$group", Value: bson.M{
			"_id":            nil,
			"total_offers":   bson.M{"$sum": 1},
			"active_offers":  bson.M{"$sum": bson.M{"$cond": bson.A{bson.M{"$eq": bson.A{"$status", "active"}}, 1, 0}}},
			"total_views":    bson.M{"$sum": "$stats.views"},
			"total_bookings": bson.M{"$sum": "$stats.bookings"},
			"total_checkins": bson.M{"$sum": "$stats.checkins"},
		}}},
	}

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var result struct {
		TotalOffers   int `bson:"total_offers"`
		ActiveOffers  int `bson:"active_offers"`
		TotalViews    int `bson:"total_views"`
		TotalBookings int `bson:"total_bookings"`
		TotalCheckins int `bson:"total_checkins"`
	}

	if cursor.Next(ctx) {
		if err := cursor.Decode(&result); err != nil {
			return nil, err
		}
	}

	return &domain.PartnerOfferStats{
		TotalOffers:   result.TotalOffers,
		ActiveOffers:  result.ActiveOffers,
		TotalViews:    result.TotalViews,
		TotalBookings: result.TotalBookings,
		TotalCheckins: result.TotalCheckins,
	}, nil
}

// toSummary converts an Offer to an OfferSummary.
func (r *OfferRepository) toSummary(offer *domain.Offer) domain.OfferSummary {
	var primaryImage string
	for _, img := range offer.Images() {
		if img.IsPrimary {
			primaryImage = img.URL
			break
		}
	}
	if primaryImage == "" && len(offer.Images()) > 0 {
		primaryImage = offer.Images()[0].URL
	}

	return domain.OfferSummary{
		ID:                offer.ID(),
		Title:             offer.Title(),
		ShortDescription:  offer.ShortDescription(),
		PrimaryImage:      primaryImage,
		Discount:          offer.Discount(),
		PartnerName:       offer.PartnerSnapshot().Name,
		EstablishmentName: offer.EstablishmentSnapshot().Name,
		City:              offer.EstablishmentSnapshot().City,
		Location:          offer.EstablishmentSnapshot().Location,
		CategoryID:        offer.CategoryID(),
	}
}
