package mongodb

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/yousoon/apps/services/booking-service/internal/domain"
)

// =============================================================================
// MONGODB DOCUMENT
// =============================================================================

type OutingDocument struct {
	ID     primitive.ObjectID `bson:"_id,omitempty"`
	UserID string             `bson:"user_id"`

	// Offer snapshot
	Offer OfferSnapshotDoc `bson:"offer"`

	// User snapshot
	User UserSnapshotDoc `bson:"user"`

	// QR Code
	QRCode QRCodeDoc `bson:"qr_code"`

	// Status
	Status   string             `bson:"status"`
	Timeline []TimelineEntryDoc `bson:"timeline"`

	// Check-in (optional)
	CheckIn *CheckInInfoDoc `bson:"check_in,omitempty"`

	// Cancellation (optional)
	Cancellation *CancellationInfoDoc `bson:"cancellation,omitempty"`

	// Timing
	BookedAt  time.Time `bson:"booked_at"`
	ExpiresAt time.Time `bson:"expires_at"`

	// Metadata
	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
}

type OfferSnapshotDoc struct {
	OfferID              string    `bson:"offer_id"`
	PartnerID            string    `bson:"partner_id"`
	EstablishmentID      string    `bson:"establishment_id"`
	Title                string    `bson:"title"`
	Description          string    `bson:"description"`
	DiscountType         string    `bson:"discount_type"`
	DiscountValue        int       `bson:"discount_value"`
	Category             string    `bson:"category"`
	EstablishmentName    string    `bson:"establishment_name"`
	EstablishmentAddress string    `bson:"establishment_address"`
	Latitude             float64   `bson:"latitude"`
	Longitude            float64   `bson:"longitude"`
	ImageURL             string    `bson:"image_url"`
	CapturedAt           time.Time `bson:"captured_at"`
}

type UserSnapshotDoc struct {
	UserID    string `bson:"user_id"`
	FirstName string `bson:"first_name"`
	LastName  string `bson:"last_name"`
	Email     string `bson:"email"`
}

type QRCodeDoc struct {
	Code      string    `bson:"code"`
	Signature string    `bson:"signature"`
	CreatedAt time.Time `bson:"created_at"`
	ExpiresAt time.Time `bson:"expires_at"`
}

type TimelineEntryDoc struct {
	Status    string                 `bson:"status"`
	Timestamp time.Time              `bson:"timestamp"`
	Actor     string                 `bson:"actor"`
	Metadata  map[string]interface{} `bson:"metadata,omitempty"`
}

type CheckInInfoDoc struct {
	CheckedInAt time.Time `bson:"checked_in_at"`
	CheckedInBy string    `bson:"checked_in_by"`
	Method      string    `bson:"method"`
	Latitude    *float64  `bson:"latitude,omitempty"`
	Longitude   *float64  `bson:"longitude,omitempty"`
}

type CancellationInfoDoc struct {
	CancelledAt time.Time `bson:"cancelled_at"`
	CancelledBy string    `bson:"cancelled_by"`
	Reason      string    `bson:"reason"`
}

// =============================================================================
// REPOSITORY IMPLEMENTATION
// =============================================================================

type OutingRepository struct {
	collection *mongo.Collection
}

func NewOutingRepository(db *mongo.Database) *OutingRepository {
	collection := db.Collection("outings")

	// Create indexes
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	indexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "user_id", Value: 1}, {Key: "created_at", Value: -1}},
			Options: options.Index().SetName("user_outings"),
		},
		{
			Keys:    bson.D{{Key: "offer.offer_id", Value: 1}},
			Options: options.Index().SetName("offer_outings"),
		},
		{
			Keys:    bson.D{{Key: "offer.partner_id", Value: 1}, {Key: "created_at", Value: -1}},
			Options: options.Index().SetName("partner_outings"),
		},
		{
			Keys:    bson.D{{Key: "offer.establishment_id", Value: 1}},
			Options: options.Index().SetName("establishment_outings"),
		},
		{
			Keys:    bson.D{{Key: "qr_code.code", Value: 1}},
			Options: options.Index().SetUnique(true).SetName("qr_code_unique"),
		},
		{
			Keys:    bson.D{{Key: "status", Value: 1}},
			Options: options.Index().SetName("status_idx"),
		},
		{
			Keys:    bson.D{{Key: "expires_at", Value: 1}},
			Options: options.Index().SetName("expires_at_idx"),
		},
		{
			Keys:    bson.D{{Key: "user_id", Value: 1}, {Key: "offer.offer_id", Value: 1}, {Key: "status", Value: 1}},
			Options: options.Index().SetName("user_offer_active"),
		},
	}

	_, _ = collection.Indexes().CreateMany(ctx, indexes)

	return &OutingRepository{collection: collection}
}

func (r *OutingRepository) Create(ctx context.Context, outing *domain.Outing) error {
	doc := r.toDocument(outing)
	doc.ID = primitive.NewObjectID()

	_, err := r.collection.InsertOne(ctx, doc)
	if err != nil {
		return fmt.Errorf("failed to insert outing: %w", err)
	}

	return nil
}

func (r *OutingRepository) Update(ctx context.Context, outing *domain.Outing) error {
	doc := r.toDocument(outing)

	filter := bson.D{{Key: "_id", Value: doc.ID}}
	update := bson.D{{Key: "$set", Value: doc}}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update outing: %w", err)
	}

	return nil
}

func (r *OutingRepository) GetByID(ctx context.Context, id string) (*domain.Outing, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, domain.ErrOutingNotFound
	}

	var doc OutingDocument
	err = r.collection.FindOne(ctx, bson.D{{Key: "_id", Value: oid}}).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, domain.ErrOutingNotFound
		}
		return nil, fmt.Errorf("failed to get outing: %w", err)
	}

	return r.toDomain(&doc), nil
}

func (r *OutingRepository) GetByQRCode(ctx context.Context, qrCode string) (*domain.Outing, error) {
	// Extract code part if full code is provided
	code := qrCode
	if len(qrCode) > 32 {
		code = qrCode[:32] // First 32 chars is the code
	}

	var doc OutingDocument
	err := r.collection.FindOne(ctx, bson.D{{Key: "qr_code.code", Value: code}}).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, domain.ErrOutingNotFound
		}
		return nil, fmt.Errorf("failed to get outing by QR: %w", err)
	}

	return r.toDomain(&doc), nil
}

func (r *OutingRepository) GetByUserID(ctx context.Context, userID string, filter domain.OutingFilter) ([]*domain.Outing, int64, error) {
	query := bson.D{{Key: "user_id", Value: userID}}
	return r.findWithFilter(ctx, query, filter)
}

func (r *OutingRepository) GetByOfferID(ctx context.Context, offerID string, filter domain.OutingFilter) ([]*domain.Outing, int64, error) {
	query := bson.D{{Key: "offer.offer_id", Value: offerID}}
	return r.findWithFilter(ctx, query, filter)
}

func (r *OutingRepository) GetByPartnerID(ctx context.Context, partnerID string, filter domain.OutingFilter) ([]*domain.Outing, int64, error) {
	query := bson.D{{Key: "offer.partner_id", Value: partnerID}}
	return r.findWithFilter(ctx, query, filter)
}

func (r *OutingRepository) GetByEstablishmentID(ctx context.Context, establishmentID string, filter domain.OutingFilter) ([]*domain.Outing, int64, error) {
	query := bson.D{{Key: "offer.establishment_id", Value: establishmentID}}
	return r.findWithFilter(ctx, query, filter)
}

func (r *OutingRepository) GetActiveByUserAndOffer(ctx context.Context, userID, offerID string) (*domain.Outing, error) {
	activeStatuses := []string{
		string(domain.OutingStatusPending),
		string(domain.OutingStatusConfirmed),
	}

	query := bson.D{
		{Key: "user_id", Value: userID},
		{Key: "offer.offer_id", Value: offerID},
		{Key: "status", Value: bson.D{{Key: "$in", Value: activeStatuses}}},
		{Key: "expires_at", Value: bson.D{{Key: "$gt", Value: time.Now()}}},
	}

	var doc OutingDocument
	err := r.collection.FindOne(ctx, query).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, domain.ErrOutingNotFound
		}
		return nil, fmt.Errorf("failed to get active outing: %w", err)
	}

	return r.toDomain(&doc), nil
}

func (r *OutingRepository) CountByUserAndOffer(ctx context.Context, userID, offerID string) (int64, error) {
	query := bson.D{
		{Key: "user_id", Value: userID},
		{Key: "offer.offer_id", Value: offerID},
	}

	count, err := r.collection.CountDocuments(ctx, query)
	if err != nil {
		return 0, fmt.Errorf("failed to count user offer bookings: %w", err)
	}

	return count, nil
}

func (r *OutingRepository) CountByOfferAndPeriod(ctx context.Context, offerID string, start, end time.Time) (int64, error) {
	query := bson.D{
		{Key: "offer.offer_id", Value: offerID},
		{Key: "created_at", Value: bson.D{
			{Key: "$gte", Value: start},
			{Key: "$lte", Value: end},
		}},
	}

	count, err := r.collection.CountDocuments(ctx, query)
	if err != nil {
		return 0, fmt.Errorf("failed to count offer period bookings: %w", err)
	}

	return count, nil
}

func (r *OutingRepository) GetExpiredOutings(ctx context.Context, before time.Time, limit int) ([]*domain.Outing, error) {
	activeStatuses := []string{
		string(domain.OutingStatusPending),
		string(domain.OutingStatusConfirmed),
	}

	query := bson.D{
		{Key: "status", Value: bson.D{{Key: "$in", Value: activeStatuses}}},
		{Key: "expires_at", Value: bson.D{{Key: "$lt", Value: before}}},
	}

	opts := options.Find().SetLimit(int64(limit))

	cursor, err := r.collection.Find(ctx, query, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to find expired outings: %w", err)
	}
	defer cursor.Close(ctx)

	var outings []*domain.Outing
	for cursor.Next(ctx) {
		var doc OutingDocument
		if err := cursor.Decode(&doc); err != nil {
			continue
		}
		outings = append(outings, r.toDomain(&doc))
	}

	return outings, nil
}

func (r *OutingRepository) Delete(ctx context.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return domain.ErrOutingNotFound
	}

	_, err = r.collection.DeleteOne(ctx, bson.D{{Key: "_id", Value: oid}})
	if err != nil {
		return fmt.Errorf("failed to delete outing: %w", err)
	}

	return nil
}

// =============================================================================
// HELPER METHODS
// =============================================================================

func (r *OutingRepository) findWithFilter(ctx context.Context, baseQuery bson.D, filter domain.OutingFilter) ([]*domain.Outing, int64, error) {
	query := baseQuery

	// Add status filter
	if len(filter.Status) > 0 {
		statuses := make([]string, len(filter.Status))
		for i, s := range filter.Status {
			statuses[i] = string(s)
		}
		query = append(query, bson.E{Key: "status", Value: bson.D{{Key: "$in", Value: statuses}}})
	}

	// Add date filters
	if filter.StartDate != nil {
		query = append(query, bson.E{Key: "created_at", Value: bson.D{{Key: "$gte", Value: *filter.StartDate}}})
	}
	if filter.EndDate != nil {
		query = append(query, bson.E{Key: "created_at", Value: bson.D{{Key: "$lte", Value: *filter.EndDate}}})
	}

	// Count total
	total, err := r.collection.CountDocuments(ctx, query)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count outings: %w", err)
	}

	// Build options
	opts := options.Find().
		SetSkip(int64(filter.Offset)).
		SetLimit(int64(filter.Limit))

	// Sort
	sortOrder := -1
	if filter.SortOrder == "asc" {
		sortOrder = 1
	}
	sortField := filter.SortBy
	if sortField == "" {
		sortField = "created_at"
	}
	opts.SetSort(bson.D{{Key: sortField, Value: sortOrder}})

	// Execute query
	cursor, err := r.collection.Find(ctx, query, opts)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to find outings: %w", err)
	}
	defer cursor.Close(ctx)

	var outings []*domain.Outing
	for cursor.Next(ctx) {
		var doc OutingDocument
		if err := cursor.Decode(&doc); err != nil {
			continue
		}
		outings = append(outings, r.toDomain(&doc))
	}

	return outings, total, nil
}

// =============================================================================
// MAPPERS
// =============================================================================

func (r *OutingRepository) toDocument(outing *domain.Outing) *OutingDocument {
	doc := &OutingDocument{
		UserID: outing.UserID(),
		Offer: OfferSnapshotDoc{
			OfferID:              outing.Offer().OfferID(),
			PartnerID:            outing.Offer().PartnerID(),
			EstablishmentID:      outing.Offer().EstablishmentID(),
			Title:                outing.Offer().Title(),
			Description:          outing.Offer().Description(),
			DiscountType:         outing.Offer().DiscountType(),
			DiscountValue:        outing.Offer().DiscountValue(),
			Category:             outing.Offer().Category(),
			EstablishmentName:    outing.Offer().EstablishmentName(),
			EstablishmentAddress: outing.Offer().EstablishmentAddress(),
			Latitude:             outing.Offer().Latitude(),
			Longitude:            outing.Offer().Longitude(),
			ImageURL:             outing.Offer().ImageURL(),
			CapturedAt:           outing.Offer().CapturedAt(),
		},
		User: UserSnapshotDoc{
			UserID:    outing.User().UserID(),
			FirstName: outing.User().FirstName(),
			LastName:  outing.User().LastName(),
			Email:     outing.User().Email(),
		},
		QRCode: QRCodeDoc{
			Code:      outing.QRCode().Code(),
			Signature: outing.QRCode().Signature(),
			CreatedAt: outing.QRCode().CreatedAt(),
			ExpiresAt: outing.QRCode().ExpiresAt(),
		},
		Status:    string(outing.Status()),
		BookedAt:  outing.BookedAt(),
		ExpiresAt: outing.ExpiresAt(),
		CreatedAt: outing.CreatedAt(),
		UpdatedAt: outing.UpdatedAt(),
	}

	// Parse ID if it's a valid ObjectID
	if oid, err := primitive.ObjectIDFromHex(outing.ID()); err == nil {
		doc.ID = oid
	}

	// Map timeline
	for _, entry := range outing.Timeline() {
		doc.Timeline = append(doc.Timeline, TimelineEntryDoc{
			Status:    string(entry.Status()),
			Timestamp: entry.Timestamp(),
			Actor:     entry.Actor(),
			Metadata:  entry.Metadata(),
		})
	}

	// Map check-in
	if outing.CheckIn() != nil {
		doc.CheckIn = &CheckInInfoDoc{
			CheckedInAt: outing.CheckIn().CheckedInAt(),
			CheckedInBy: outing.CheckIn().CheckedInBy(),
			Method:      string(outing.CheckIn().Method()),
			Latitude:    outing.CheckIn().Latitude(),
			Longitude:   outing.CheckIn().Longitude(),
		}
	}

	// Map cancellation
	if outing.Cancellation() != nil {
		doc.Cancellation = &CancellationInfoDoc{
			CancelledAt: outing.Cancellation().CancelledAt(),
			CancelledBy: string(outing.Cancellation().CancelledBy()),
			Reason:      outing.Cancellation().Reason(),
		}
	}

	return doc
}

func (r *OutingRepository) toDomain(doc *OutingDocument) *domain.Outing {
	// Reconstruct offer snapshot
	offer := domain.ReconstructOfferSnapshot(
		doc.Offer.OfferID,
		doc.Offer.PartnerID,
		doc.Offer.EstablishmentID,
		doc.Offer.Title,
		doc.Offer.Description,
		doc.Offer.DiscountType,
		doc.Offer.DiscountValue,
		doc.Offer.Category,
		doc.Offer.EstablishmentName,
		doc.Offer.EstablishmentAddress,
		doc.Offer.Latitude,
		doc.Offer.Longitude,
		doc.Offer.ImageURL,
		doc.Offer.CapturedAt,
	)

	// Reconstruct user snapshot
	user := domain.NewUserSnapshot(
		doc.User.UserID,
		doc.User.FirstName,
		doc.User.LastName,
		doc.User.Email,
	)

	// Reconstruct QR code
	qrCode := domain.ReconstructQRCode(
		doc.QRCode.Code,
		doc.QRCode.Signature,
		doc.QRCode.CreatedAt,
		doc.QRCode.ExpiresAt,
	)

	// Reconstruct timeline
	var timeline []domain.TimelineEntry
	for _, entry := range doc.Timeline {
		timeline = append(timeline, domain.ReconstructTimelineEntry(
			domain.OutingStatus(entry.Status),
			entry.Timestamp,
			entry.Actor,
			entry.Metadata,
		))
	}

	// Reconstruct check-in
	var checkIn *domain.CheckInInfo
	if doc.CheckIn != nil {
		ci := domain.ReconstructCheckInInfo(
			doc.CheckIn.CheckedInAt,
			doc.CheckIn.CheckedInBy,
			domain.CheckInMethod(doc.CheckIn.Method),
			doc.CheckIn.Latitude,
			doc.CheckIn.Longitude,
		)
		checkIn = &ci
	}

	// Reconstruct cancellation
	var cancellation *domain.CancellationInfo
	if doc.Cancellation != nil {
		c := domain.ReconstructCancellationInfo(
			doc.Cancellation.CancelledAt,
			domain.CancellationActor(doc.Cancellation.CancelledBy),
			doc.Cancellation.Reason,
		)
		cancellation = &c
	}

	return domain.ReconstructOuting(
		doc.ID.Hex(),
		doc.UserID,
		offer,
		user,
		qrCode,
		domain.OutingStatus(doc.Status),
		timeline,
		checkIn,
		cancellation,
		doc.BookedAt,
		doc.ExpiresAt,
		doc.CreatedAt,
		doc.UpdatedAt,
	)
}
