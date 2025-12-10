package mongodb

import (
	"context"
	"time"

	"github.com/yousoon/apps/services/engagement-service/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ReviewDocument représente un avis en MongoDB
type ReviewDocument struct {
	ID                 primitive.ObjectID `bson:"_id,omitempty"`
	UserID             string             `bson:"userId"`
	OfferID            string             `bson:"offerId"`
	PartnerID          string             `bson:"partnerId"`
	EstablishmentID    string             `bson:"establishmentId"`
	OutingID           string             `bson:"outingId,omitempty"`
	Rating             int                `bson:"rating"`
	Title              string             `bson:"title,omitempty"`
	Content            string             `bson:"content"`
	Images             []string           `bson:"images,omitempty"`
	HelpfulCount       int                `bson:"helpfulCount"`
	IsVerifiedPurchase bool               `bson:"isVerifiedPurchase"`
	Moderation         ModerationDocument `bson:"moderation"`
	User               UserSnapshotDoc    `bson:"_user"`
	Offer              OfferSnapshotDoc   `bson:"_offer"`
	Partner            PartnerSnapshotDoc `bson:"_partner"`
	CreatedAt          time.Time          `bson:"createdAt"`
	UpdatedAt          time.Time          `bson:"updatedAt"`
}

type ModerationDocument struct {
	Status       string           `bson:"status"`
	Reports      []ReportDocument `bson:"reports,omitempty"`
	ReviewedBy   string           `bson:"reviewedBy,omitempty"`
	ReviewedAt   *time.Time       `bson:"reviewedAt,omitempty"`
	RejectReason string           `bson:"rejectReason,omitempty"`
}

type ReportDocument struct {
	UserID     string    `bson:"userId"`
	Reason     string    `bson:"reason"`
	ReportedAt time.Time `bson:"reportedAt"`
}

type UserSnapshotDoc struct {
	FirstName string `bson:"firstName"`
	Avatar    string `bson:"avatar,omitempty"`
}

type OfferSnapshotDoc struct {
	Title string `bson:"title"`
}

type PartnerSnapshotDoc struct {
	Name string `bson:"name"`
}

// ReviewRepository implémentation MongoDB
type ReviewRepository struct {
	collection *mongo.Collection
}

// NewReviewRepository crée un nouveau repository
func NewReviewRepository(db *mongo.Database) *ReviewRepository {
	coll := db.Collection("reviews")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	indexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "userId", Value: 1}, {Key: "offerId", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{{Key: "offerId", Value: 1}, {Key: "createdAt", Value: -1}},
		},
		{
			Keys: bson.D{{Key: "partnerId", Value: 1}, {Key: "createdAt", Value: -1}},
		},
		{
			Keys: bson.D{{Key: "userId", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "moderation.status", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "rating", Value: 1}},
		},
	}

	coll.Indexes().CreateMany(ctx, indexes)

	return &ReviewRepository{collection: coll}
}

// Create crée un avis
func (r *ReviewRepository) Create(ctx context.Context, review *domain.Review) error {
	doc := r.toDocument(review)

	result, err := r.collection.InsertOne(ctx, doc)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return domain.ErrReviewAlreadyExists
		}
		return err
	}

	review.ID = result.InsertedID.(primitive.ObjectID).Hex()
	return nil
}

// Update met à jour un avis
func (r *ReviewRepository) Update(ctx context.Context, review *domain.Review) error {
	id, err := primitive.ObjectIDFromHex(review.ID)
	if err != nil {
		return domain.ErrReviewNotFound
	}

	doc := r.toDocument(review)
	doc.ID = id

	result, err := r.collection.ReplaceOne(ctx, bson.M{"_id": id}, doc)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return domain.ErrReviewNotFound
	}
	return nil
}

// Delete supprime un avis
func (r *ReviewRepository) Delete(ctx context.Context, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return domain.ErrReviewNotFound
	}

	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return domain.ErrReviewNotFound
	}
	return nil
}

// FindByID trouve un avis par ID
func (r *ReviewRepository) FindByID(ctx context.Context, id string) (*domain.Review, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, domain.ErrReviewNotFound
	}

	var doc ReviewDocument
	err = r.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, domain.ErrReviewNotFound
		}
		return nil, err
	}

	return r.toDomain(&doc), nil
}

// FindByOfferID trouve les avis d'une offre
func (r *ReviewRepository) FindByOfferID(ctx context.Context, offerID string, limit, offset int) ([]*domain.Review, int, error) {
	filter := bson.M{
		"offerId":           offerID,
		"moderation.status": "approved",
	}

	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	opts := options.Find().
		SetSort(bson.D{{Key: "createdAt", Value: -1}}).
		SetSkip(int64(offset)).
		SetLimit(int64(limit))

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var reviews []*domain.Review
	for cursor.Next(ctx) {
		var doc ReviewDocument
		if err := cursor.Decode(&doc); err != nil {
			return nil, 0, err
		}
		reviews = append(reviews, r.toDomain(&doc))
	}

	return reviews, int(total), nil
}

// FindByUserID trouve les avis d'un utilisateur
func (r *ReviewRepository) FindByUserID(ctx context.Context, userID string, limit, offset int) ([]*domain.Review, int, error) {
	filter := bson.M{"userId": userID}

	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	opts := options.Find().
		SetSort(bson.D{{Key: "createdAt", Value: -1}}).
		SetSkip(int64(offset)).
		SetLimit(int64(limit))

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var reviews []*domain.Review
	for cursor.Next(ctx) {
		var doc ReviewDocument
		if err := cursor.Decode(&doc); err != nil {
			return nil, 0, err
		}
		reviews = append(reviews, r.toDomain(&doc))
	}

	return reviews, int(total), nil
}

// GetAverageRating calcule la note moyenne d'une offre
func (r *ReviewRepository) GetAverageRating(ctx context.Context, offerID string) (float64, int, error) {
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{
			"offerId":           offerID,
			"moderation.status": "approved",
		}}},
		{{Key: "$group", Value: bson.M{
			"_id":   nil,
			"avg":   bson.M{"$avg": "$rating"},
			"count": bson.M{"$sum": 1},
		}}},
	}

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return 0, 0, err
	}
	defer cursor.Close(ctx)

	if cursor.Next(ctx) {
		var result struct {
			Avg   float64 `bson:"avg"`
			Count int     `bson:"count"`
		}
		if err := cursor.Decode(&result); err != nil {
			return 0, 0, err
		}
		return result.Avg, result.Count, nil
	}

	return 0, 0, nil
}

// AddReport ajoute un signalement
func (r *ReviewRepository) AddReport(ctx context.Context, reviewID, userID, reason string) error {
	objID, err := primitive.ObjectIDFromHex(reviewID)
	if err != nil {
		return domain.ErrReviewNotFound
	}

	report := ReportDocument{
		UserID:     userID,
		Reason:     reason,
		ReportedAt: time.Now(),
	}

	_, err = r.collection.UpdateOne(ctx,
		bson.M{"_id": objID},
		bson.M{
			"$push": bson.M{"moderation.reports": report},
			"$set":  bson.M{"moderation.status": "reported"},
		},
	)
	return err
}

// Exists vérifie si un avis existe
func (r *ReviewRepository) Exists(ctx context.Context, userID, offerID string) (bool, error) {
	count, err := r.collection.CountDocuments(ctx, bson.M{
		"userId":  userID,
		"offerId": offerID,
	})
	return count > 0, err
}

func (r *ReviewRepository) toDocument(review *domain.Review) *ReviewDocument {
	doc := &ReviewDocument{
		UserID:             review.UserID,
		OfferID:            review.OfferID,
		PartnerID:          review.PartnerID,
		EstablishmentID:    review.EstablishmentID,
		OutingID:           review.OutingID,
		Rating:             review.Rating,
		Title:              review.Title,
		Content:            review.Content,
		Images:             review.Images,
		HelpfulCount:       review.HelpfulCount,
		IsVerifiedPurchase: review.IsVerifiedPurchase,
		Moderation: ModerationDocument{
			Status: string(review.Moderation.Status),
		},
		User: UserSnapshotDoc{
			FirstName: review.User.FirstName,
			Avatar:    review.User.Avatar,
		},
		Offer: OfferSnapshotDoc{
			Title: review.Offer.Title,
		},
		Partner: PartnerSnapshotDoc{
			Name: review.Partner.Name,
		},
		CreatedAt: review.CreatedAt,
		UpdatedAt: review.UpdatedAt,
	}

	if review.ID != "" {
		if id, err := primitive.ObjectIDFromHex(review.ID); err == nil {
			doc.ID = id
		}
	}

	return doc
}

func (r *ReviewRepository) toDomain(doc *ReviewDocument) *domain.Review {
	return &domain.Review{
		ID:                 doc.ID.Hex(),
		UserID:             doc.UserID,
		OfferID:            doc.OfferID,
		PartnerID:          doc.PartnerID,
		EstablishmentID:    doc.EstablishmentID,
		OutingID:           doc.OutingID,
		Rating:             doc.Rating,
		Title:              doc.Title,
		Content:            doc.Content,
		Images:             doc.Images,
		HelpfulCount:       doc.HelpfulCount,
		IsVerifiedPurchase: doc.IsVerifiedPurchase,
		Moderation: domain.Moderation{
			Status: domain.ModerationStatus(doc.Moderation.Status),
		},
		User: domain.UserSnapshot{
			FirstName: doc.User.FirstName,
			Avatar:    doc.User.Avatar,
		},
		Offer: domain.OfferSnapshot{
			Title: doc.Offer.Title,
		},
		Partner: domain.PartnerSnapshot{
			Name: doc.Partner.Name,
		},
		CreatedAt: doc.CreatedAt,
		UpdatedAt: doc.UpdatedAt,
	}
}
