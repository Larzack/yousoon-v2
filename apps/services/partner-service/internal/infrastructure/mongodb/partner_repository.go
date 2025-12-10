// Package mongodb implements the MongoDB repository for the Partner bounded context.
package mongodb

import (
	"context"
	"time"

	partnerdomain "github.com/yousoon/services/partner/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// =============================================================================
// Partner Repository Implementation
// =============================================================================

// PartnerRepository is the MongoDB implementation of the partner repository.
type PartnerRepository struct {
	collection *mongo.Collection
}

// NewPartnerRepository creates a new partner repository.
func NewPartnerRepository(db *mongo.Database) *PartnerRepository {
	collection := db.Collection("partners")

	// Create indexes
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	indexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "company.siret", Value: 1}},
			Options: options.Index().SetUnique(true).SetSparse(true),
		},
		{
			Keys: bson.D{{Key: "ownerUserId", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "status", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "category", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "teamMembers.email.value", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "establishments.location", Value: "2dsphere"}},
		},
		{
			Keys: bson.D{
				{Key: "company.name", Value: "text"},
				{Key: "company.tradeName", Value: "text"},
			},
		},
	}

	_, _ = collection.Indexes().CreateMany(ctx, indexes)

	return &PartnerRepository{collection: collection}
}

// Save saves a partner (insert or update).
func (r *PartnerRepository) Save(ctx context.Context, partner *partnerdomain.Partner) error {
	partner.MarkUpdated()

	filter := bson.M{"_id": partner.ID}
	update := bson.M{"$set": partner}
	opts := options.Update().SetUpsert(true)

	_, err := r.collection.UpdateOne(ctx, filter, update, opts)
	return err
}

// FindByID finds a partner by ID.
func (r *PartnerRepository) FindByID(ctx context.Context, id partnerdomain.PartnerID) (*partnerdomain.Partner, error) {
	var partner partnerdomain.Partner

	filter := bson.M{
		"_id":       id,
		"deletedAt": bson.M{"$eq": nil},
	}

	err := r.collection.FindOne(ctx, filter).Decode(&partner)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, partnerdomain.ErrPartnerNotFound
		}
		return nil, err
	}

	return &partner, nil
}

// FindByOwnerUserID finds a partner by owner user ID.
func (r *PartnerRepository) FindByOwnerUserID(ctx context.Context, userID partnerdomain.UserID) (*partnerdomain.Partner, error) {
	var partner partnerdomain.Partner

	filter := bson.M{
		"ownerUserId": userID,
		"deletedAt":   bson.M{"$eq": nil},
	}

	err := r.collection.FindOne(ctx, filter).Decode(&partner)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, partnerdomain.ErrPartnerNotFound
		}
		return nil, err
	}

	return &partner, nil
}

// FindBySIRET finds a partner by SIRET.
func (r *PartnerRepository) FindBySIRET(ctx context.Context, siret string) (*partnerdomain.Partner, error) {
	var partner partnerdomain.Partner

	filter := bson.M{
		"company.siret": siret,
		"deletedAt":     bson.M{"$eq": nil},
	}

	err := r.collection.FindOne(ctx, filter).Decode(&partner)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, partnerdomain.ErrPartnerNotFound
		}
		return nil, err
	}

	return &partner, nil
}

// FindByTeamMemberEmail finds partners where the email is a team member.
func (r *PartnerRepository) FindByTeamMemberEmail(ctx context.Context, email partnerdomain.Email) ([]*partnerdomain.Partner, error) {
	filter := bson.M{
		"teamMembers.email.value": email.Value,
		"deletedAt":               bson.M{"$eq": nil},
	}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var partners []*partnerdomain.Partner
	if err := cursor.All(ctx, &partners); err != nil {
		return nil, err
	}

	return partners, nil
}

// Delete soft-deletes a partner.
func (r *PartnerRepository) Delete(ctx context.Context, id partnerdomain.PartnerID) error {
	filter := bson.M{"_id": id}
	update := bson.M{
		"$set": bson.M{
			"deletedAt": time.Now(),
			"status":    partnerdomain.PartnerStatusSuspended,
			"updatedAt": time.Now(),
		},
	}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

// List lists partners with pagination and filters.
func (r *PartnerRepository) List(ctx context.Context, filter partnerdomain.PartnerFilter) ([]*partnerdomain.Partner, int64, error) {
	mongoFilter := r.buildFilter(filter)

	// Count total
	total, err := r.collection.CountDocuments(ctx, mongoFilter)
	if err != nil {
		return nil, 0, err
	}

	// Find with pagination
	opts := options.Find().
		SetSkip(int64(filter.Offset)).
		SetLimit(int64(filter.Limit)).
		SetSort(bson.D{{Key: "createdAt", Value: -1}})

	cursor, err := r.collection.Find(ctx, mongoFilter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var partners []*partnerdomain.Partner
	if err := cursor.All(ctx, &partners); err != nil {
		return nil, 0, err
	}

	return partners, total, nil
}

// Count counts partners matching a filter.
func (r *PartnerRepository) Count(ctx context.Context, filter partnerdomain.PartnerFilter) (int64, error) {
	mongoFilter := r.buildFilter(filter)
	return r.collection.CountDocuments(ctx, mongoFilter)
}

// buildFilter builds a MongoDB filter from a domain filter.
func (r *PartnerRepository) buildFilter(filter partnerdomain.PartnerFilter) bson.M {
	mongoFilter := bson.M{
		"deletedAt": bson.M{"$eq": nil},
	}

	if filter.Status != nil {
		mongoFilter["status"] = *filter.Status
	}

	if filter.Category != "" {
		mongoFilter["category"] = filter.Category
	}

	if filter.Search != "" {
		mongoFilter["$text"] = bson.M{"$search": filter.Search}
	}

	return mongoFilter
}

// =============================================================================
// Read Repository Implementation
// =============================================================================

// PartnerReadRepository is the MongoDB implementation of the partner read repository.
type PartnerReadRepository struct {
	collection *mongo.Collection
}

// NewPartnerReadRepository creates a new partner read repository.
func NewPartnerReadRepository(db *mongo.Database) *PartnerReadRepository {
	return &PartnerReadRepository{
		collection: db.Collection("partners"),
	}
}

// GetPartnerSummary gets a partner summary for display.
func (r *PartnerReadRepository) GetPartnerSummary(ctx context.Context, id partnerdomain.PartnerID) (*partnerdomain.PartnerSummary, error) {
	var partner partnerdomain.Partner

	filter := bson.M{
		"_id":       id,
		"deletedAt": bson.M{"$eq": nil},
	}

	err := r.collection.FindOne(ctx, filter).Decode(&partner)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, partnerdomain.ErrPartnerNotFound
		}
		return nil, err
	}

	return &partnerdomain.PartnerSummary{
		ID:                 partner.ID,
		CompanyName:        partner.Company.Name,
		TradeName:          partner.Company.TradeName,
		Category:           partner.Category,
		Logo:               partner.Branding.Logo,
		Status:             partner.Status,
		EstablishmentCount: len(partner.Establishments),
		ActiveOfferCount:   partner.Stats.ActiveOffers,
		AvgRating:          partner.Stats.AvgRating,
		ReviewCount:        partner.Stats.ReviewCount,
	}, nil
}

// GetEstablishmentsNearLocation gets establishments near a location.
func (r *PartnerReadRepository) GetEstablishmentsNearLocation(ctx context.Context, location partnerdomain.GeoLocation, radiusKm float64, limit int) ([]partnerdomain.EstablishmentSummary, error) {
	// Use aggregation pipeline for geospatial query
	pipeline := mongo.Pipeline{
		// Match active partners
		{{Key: "$match", Value: bson.M{
			"status":    partnerdomain.PartnerStatusActive,
			"deletedAt": bson.M{"$eq": nil},
		}}},
		// Unwind establishments
		{{Key: "$unwind", Value: "$establishments"}},
		// Match active establishments near location
		{{Key: "$match", Value: bson.M{
			"establishments.isActive": true,
			"establishments.location": bson.M{
				"$nearSphere": bson.M{
					"$geometry":    location,
					"$maxDistance": radiusKm * 1000, // Convert to meters
				},
			},
		}}},
		// Limit results
		{{Key: "$limit", Value: limit}},
		// Project summary fields
		{{Key: "$project", Value: bson.M{
			"establishmentId": "$establishments.id",
			"partnerId":       "$_id",
			"name":            "$establishments.name",
			"type":            "$establishments.type",
			"address":         "$establishments.address.formatted",
			"city":            "$establishments.address.city",
			"location":        "$establishments.location",
			"primaryImage": bson.M{
				"$arrayElemAt": bson.A{
					bson.M{"$filter": bson.M{
						"input": "$establishments.images",
						"as":    "img",
						"cond":  bson.M{"$eq": bson.A{"$$img.isPrimary", true}},
					}},
					0,
				},
			},
		}}},
	}

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []partnerdomain.EstablishmentSummary
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, nil
}

// GetPartnerByEstablishmentID gets a partner by establishment ID.
func (r *PartnerReadRepository) GetPartnerByEstablishmentID(ctx context.Context, estID partnerdomain.EstablishmentID) (*partnerdomain.Partner, error) {
	var partner partnerdomain.Partner

	filter := bson.M{
		"establishments.id": estID,
		"deletedAt":         bson.M{"$eq": nil},
	}

	err := r.collection.FindOne(ctx, filter).Decode(&partner)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, partnerdomain.ErrPartnerNotFound
		}
		return nil, err
	}

	return &partner, nil
}

// SearchPartners searches partners by name.
func (r *PartnerReadRepository) SearchPartners(ctx context.Context, query string, limit int) ([]*partnerdomain.PartnerSummary, error) {
	filter := bson.M{
		"$text":     bson.M{"$search": query},
		"status":    partnerdomain.PartnerStatusActive,
		"deletedAt": bson.M{"$eq": nil},
	}

	opts := options.Find().
		SetLimit(int64(limit)).
		SetProjection(bson.M{"score": bson.M{"$meta": "textScore"}}).
		SetSort(bson.M{"score": bson.M{"$meta": "textScore"}})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var partners []*partnerdomain.Partner
	if err := cursor.All(ctx, &partners); err != nil {
		return nil, err
	}

	// Map to summaries
	summaries := make([]*partnerdomain.PartnerSummary, len(partners))
	for i, p := range partners {
		summaries[i] = &partnerdomain.PartnerSummary{
			ID:                 p.ID,
			CompanyName:        p.Company.Name,
			TradeName:          p.Company.TradeName,
			Category:           p.Category,
			Logo:               p.Branding.Logo,
			Status:             p.Status,
			EstablishmentCount: len(p.Establishments),
			ActiveOfferCount:   p.Stats.ActiveOffers,
			AvgRating:          p.Stats.AvgRating,
			ReviewCount:        p.Stats.ReviewCount,
		}
	}

	return summaries, nil
}
