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

// FavoriteDocument représente un favori en MongoDB
type FavoriteDocument struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	UserID    string             `bson:"userId"`
	OfferID   string             `bson:"offerId"`
	CreatedAt time.Time          `bson:"createdAt"`
}

// FavoriteRepository implémentation MongoDB
type FavoriteRepository struct {
	collection *mongo.Collection
}

// NewFavoriteRepository crée un nouveau repository
func NewFavoriteRepository(db *mongo.Database) *FavoriteRepository {
	coll := db.Collection("favorites")

	// Créer les index
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	indexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "userId", Value: 1}, {Key: "offerId", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{{Key: "userId", Value: 1}, {Key: "createdAt", Value: -1}},
		},
		{
			Keys: bson.D{{Key: "offerId", Value: 1}},
		},
	}

	coll.Indexes().CreateMany(ctx, indexes)

	return &FavoriteRepository{collection: coll}
}

// Create ajoute un favori
func (r *FavoriteRepository) Create(ctx context.Context, favorite *domain.Favorite) error {
	doc := &FavoriteDocument{
		UserID:    favorite.UserID,
		OfferID:   favorite.OfferID,
		CreatedAt: favorite.CreatedAt,
	}

	result, err := r.collection.InsertOne(ctx, doc)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return domain.ErrFavoriteAlreadyExists
		}
		return err
	}

	favorite.ID = result.InsertedID.(primitive.ObjectID).Hex()
	return nil
}

// Delete supprime un favori
func (r *FavoriteRepository) Delete(ctx context.Context, userID, offerID string) error {
	result, err := r.collection.DeleteOne(ctx, bson.M{
		"userId":  userID,
		"offerId": offerID,
	})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return domain.ErrFavoriteNotFound
	}
	return nil
}

// FindByUserID trouve les favoris d'un utilisateur
func (r *FavoriteRepository) FindByUserID(ctx context.Context, userID string, limit, offset int) ([]*domain.Favorite, int, error) {
	filter := bson.M{"userId": userID}

	// Compter le total
	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	// Récupérer les favoris
	opts := options.Find().
		SetSort(bson.D{{Key: "createdAt", Value: -1}}).
		SetSkip(int64(offset)).
		SetLimit(int64(limit))

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var favorites []*domain.Favorite
	for cursor.Next(ctx) {
		var doc FavoriteDocument
		if err := cursor.Decode(&doc); err != nil {
			return nil, 0, err
		}
		favorites = append(favorites, r.toDomain(&doc))
	}

	return favorites, int(total), nil
}

// Exists vérifie si un favori existe
func (r *FavoriteRepository) Exists(ctx context.Context, userID, offerID string) (bool, error) {
	count, err := r.collection.CountDocuments(ctx, bson.M{
		"userId":  userID,
		"offerId": offerID,
	})
	return count > 0, err
}

// CountByOfferID compte les favoris d'une offre
func (r *FavoriteRepository) CountByOfferID(ctx context.Context, offerID string) (int, error) {
	count, err := r.collection.CountDocuments(ctx, bson.M{"offerId": offerID})
	return int(count), err
}

func (r *FavoriteRepository) toDomain(doc *FavoriteDocument) *domain.Favorite {
	return &domain.Favorite{
		ID:        doc.ID.Hex(),
		UserID:    doc.UserID,
		OfferID:   doc.OfferID,
		CreatedAt: doc.CreatedAt,
	}
}
