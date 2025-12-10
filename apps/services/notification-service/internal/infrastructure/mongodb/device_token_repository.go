package mongodb

import (
	"context"
	"time"

	"github.com/yousoon/apps/services/notification-service/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// DeviceTokenDocument représente un token de device en MongoDB
type DeviceTokenDocument struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	UserID    string             `bson:"userId"`
	Token     string             `bson:"token"`
	Platform  string             `bson:"platform"`
	IsActive  bool               `bson:"isActive"`
	CreatedAt time.Time          `bson:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt"`
}

// DeviceTokenRepository implémentation MongoDB
type DeviceTokenRepository struct {
	collection *mongo.Collection
}

// NewDeviceTokenRepository crée un nouveau repository
func NewDeviceTokenRepository(db *mongo.Database) *DeviceTokenRepository {
	coll := db.Collection("device_tokens")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	indexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "token", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{{Key: "userId", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "userId", Value: 1}, {Key: "isActive", Value: 1}},
		},
	}

	coll.Indexes().CreateMany(ctx, indexes)

	return &DeviceTokenRepository{collection: coll}
}

// Create crée un nouveau token
func (r *DeviceTokenRepository) Create(ctx context.Context, token *domain.DeviceToken) error {
	now := time.Now()
	doc := &DeviceTokenDocument{
		UserID:    token.UserID,
		Token:     token.Token,
		Platform:  string(token.Platform),
		IsActive:  true,
		CreatedAt: now,
		UpdatedAt: now,
	}

	result, err := r.collection.InsertOne(ctx, doc)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			// Si le token existe déjà, le mettre à jour avec le nouvel utilisateur
			return r.UpdateUserID(ctx, token.Token, token.UserID)
		}
		return err
	}

	token.ID = result.InsertedID.(primitive.ObjectID).Hex()
	return nil
}

// UpdateUserID met à jour l'utilisateur d'un token
func (r *DeviceTokenRepository) UpdateUserID(ctx context.Context, token, userID string) error {
	_, err := r.collection.UpdateOne(ctx,
		bson.M{"token": token},
		bson.M{"$set": bson.M{
			"userId":    userID,
			"isActive":  true,
			"updatedAt": time.Now(),
		}},
	)
	return err
}

// Delete supprime un token
func (r *DeviceTokenRepository) Delete(ctx context.Context, token string) error {
	result, err := r.collection.DeleteOne(ctx, bson.M{"token": token})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return domain.ErrDeviceTokenNotFound
	}
	return nil
}

// Deactivate désactive un token
func (r *DeviceTokenRepository) Deactivate(ctx context.Context, token string) error {
	_, err := r.collection.UpdateOne(ctx,
		bson.M{"token": token},
		bson.M{"$set": bson.M{"isActive": false, "updatedAt": time.Now()}},
	)
	return err
}

// FindByUserID trouve les tokens d'un utilisateur
func (r *DeviceTokenRepository) FindByUserID(ctx context.Context, userID string) ([]*domain.DeviceToken, error) {
	cursor, err := r.collection.Find(ctx, bson.M{
		"userId":   userID,
		"isActive": true,
	})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var tokens []*domain.DeviceToken
	for cursor.Next(ctx) {
		var doc DeviceTokenDocument
		if err := cursor.Decode(&doc); err != nil {
			return nil, err
		}
		tokens = append(tokens, r.toDomain(&doc))
	}

	return tokens, nil
}

// FindByToken trouve un token par sa valeur
func (r *DeviceTokenRepository) FindByToken(ctx context.Context, token string) (*domain.DeviceToken, error) {
	var doc DeviceTokenDocument
	err := r.collection.FindOne(ctx, bson.M{"token": token}).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, domain.ErrDeviceTokenNotFound
		}
		return nil, err
	}
	return r.toDomain(&doc), nil
}

func (r *DeviceTokenRepository) toDomain(doc *DeviceTokenDocument) *domain.DeviceToken {
	return &domain.DeviceToken{
		ID:        doc.ID.Hex(),
		UserID:    doc.UserID,
		Token:     doc.Token,
		Platform:  domain.Platform(doc.Platform),
		IsActive:  doc.IsActive,
		CreatedAt: doc.CreatedAt,
		UpdatedAt: doc.UpdatedAt,
	}
}
