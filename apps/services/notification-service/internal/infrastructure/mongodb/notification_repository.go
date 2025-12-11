package mongodb

import (
	"context"
	"encoding/json"
	"time"

	"github.com/yousoon/apps/services/notification-service/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// NotificationDocument représente une notification en MongoDB
type NotificationDocument struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	UserID      string             `bson:"userId"`
	Type        string             `bson:"type"`
	Channel     string             `bson:"channel"`
	Title       string             `bson:"title"`
	Body        string             `bson:"body"`
	Image       string             `bson:"image,omitempty"`
	Data        string             `bson:"data,omitempty"`
	RelatedType string             `bson:"relatedType,omitempty"`
	RelatedID   string             `bson:"relatedId,omitempty"`
	Status      string             `bson:"status"`
	SentAt      *time.Time         `bson:"sentAt,omitempty"`
	DeliveredAt *time.Time         `bson:"deliveredAt,omitempty"`
	ReadAt      *time.Time         `bson:"readAt,omitempty"`
	Error       string             `bson:"error,omitempty"`
	CreatedAt   time.Time          `bson:"createdAt"`
}

// NotificationRepository implémentation MongoDB
type NotificationRepository struct {
	collection *mongo.Collection
}

// NewNotificationRepository crée un nouveau repository
func NewNotificationRepository(db *mongo.Database) *NotificationRepository {
	coll := db.Collection("notifications")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	indexes := []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "userId", Value: 1}, {Key: "createdAt", Value: -1}},
		},
		{
			Keys: bson.D{{Key: "userId", Value: 1}, {Key: "status", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "status", Value: 1}},
		},
		{
			Keys:    bson.D{{Key: "createdAt", Value: 1}},
			Options: options.Index().SetExpireAfterSeconds(7776000), // 90 jours TTL
		},
	}

	coll.Indexes().CreateMany(ctx, indexes)

	return &NotificationRepository{collection: coll}
}

// Create crée une notification
func (r *NotificationRepository) Create(ctx context.Context, notification *domain.Notification) error {
	doc := r.toDocument(notification)

	_, err := r.collection.InsertOne(ctx, doc)
	if err != nil {
		return err
	}

	return nil
}

// Update met à jour une notification
func (r *NotificationRepository) Update(ctx context.Context, notification *domain.Notification) error {
	id, err := primitive.ObjectIDFromHex(notification.ID())
	if err != nil {
		return domain.ErrNotificationNotFound
	}

	doc := r.toDocument(notification)
	doc.ID = id

	result, err := r.collection.ReplaceOne(ctx, bson.M{"_id": id}, doc)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return domain.ErrNotificationNotFound
	}
	return nil
}

// FindByID trouve une notification par ID
func (r *NotificationRepository) FindByID(ctx context.Context, id string) (*domain.Notification, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, domain.ErrNotificationNotFound
	}

	var doc NotificationDocument
	err = r.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, domain.ErrNotificationNotFound
		}
		return nil, err
	}

	return r.toDomain(&doc), nil
}

// FindByUserID trouve les notifications d'un utilisateur
func (r *NotificationRepository) FindByUserID(ctx context.Context, userID string, unreadOnly bool, limit, offset int) ([]*domain.Notification, int, error) {
	filter := bson.M{"userId": userID}
	if unreadOnly {
		filter["readAt"] = nil
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

	var notifications []*domain.Notification
	for cursor.Next(ctx) {
		var doc NotificationDocument
		if err := cursor.Decode(&doc); err != nil {
			return nil, 0, err
		}
		notifications = append(notifications, r.toDomain(&doc))
	}

	return notifications, int(total), nil
}

// MarkAsRead marque une notification comme lue
func (r *NotificationRepository) MarkAsRead(ctx context.Context, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return domain.ErrNotificationNotFound
	}

	now := time.Now()
	result, err := r.collection.UpdateOne(ctx,
		bson.M{"_id": objID},
		bson.M{"$set": bson.M{"readAt": now, "status": "read"}},
	)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return domain.ErrNotificationNotFound
	}
	return nil
}

// MarkAllAsRead marque toutes les notifications d'un utilisateur comme lues
func (r *NotificationRepository) MarkAllAsRead(ctx context.Context, userID string) error {
	now := time.Now()
	_, err := r.collection.UpdateMany(ctx,
		bson.M{"userId": userID, "readAt": nil},
		bson.M{"$set": bson.M{"readAt": now, "status": "read"}},
	)
	return err
}

// CountUnread compte les notifications non lues
func (r *NotificationRepository) CountUnread(ctx context.Context, userID string) (int, error) {
	count, err := r.collection.CountDocuments(ctx, bson.M{
		"userId": userID,
		"readAt": nil,
	})
	return int(count), err
}

// FindPending trouve les notifications en attente d'envoi
func (r *NotificationRepository) FindPending(ctx context.Context, limit int) ([]*domain.Notification, error) {
	opts := options.Find().
		SetSort(bson.D{{Key: "createdAt", Value: 1}}).
		SetLimit(int64(limit))

	cursor, err := r.collection.Find(ctx, bson.M{"status": "pending"}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var notifications []*domain.Notification
	for cursor.Next(ctx) {
		var doc NotificationDocument
		if err := cursor.Decode(&doc); err != nil {
			return nil, err
		}
		notifications = append(notifications, r.toDomain(&doc))
	}

	return notifications, nil
}

func (r *NotificationRepository) toDocument(n *domain.Notification) *NotificationDocument {
	// Serialize Data map to JSON string
	var dataStr string
	if n.Data() != nil {
		if bytes, err := json.Marshal(n.Data()); err == nil {
			dataStr = string(bytes)
		}
	}

	// Handle optional string pointers
	var imageURL string
	if n.ImageURL() != nil {
		imageURL = *n.ImageURL()
	}

	var relatedID string
	if n.RelatedID() != nil {
		relatedID = *n.RelatedID()
	}

	var errMsg string
	if n.Error() != nil {
		errMsg = *n.Error()
	}

	doc := &NotificationDocument{
		UserID:      n.UserID(),
		Type:        string(n.Type()),
		Channel:     string(n.Channel()),
		Title:       n.Title(),
		Body:        n.Body(),
		Image:       imageURL,
		Data:        dataStr,
		RelatedType: n.RelatedType(),
		RelatedID:   relatedID,
		Status:      string(n.Status()),
		SentAt:      n.SentAt(),
		DeliveredAt: n.DeliveredAt(),
		ReadAt:      n.ReadAt(),
		Error:       errMsg,
		CreatedAt:   n.CreatedAt(),
	}

	if n.ID() != "" {
		if id, err := primitive.ObjectIDFromHex(n.ID()); err == nil {
			doc.ID = id
		}
	}

	return doc
}

func (r *NotificationRepository) toDomain(doc *NotificationDocument) *domain.Notification {
	// Deserialize Data JSON string to map
	var data map[string]interface{}
	if doc.Data != "" {
		json.Unmarshal([]byte(doc.Data), &data)
	}

	// Handle optional string pointers
	var imageURL *string
	if doc.Image != "" {
		imageURL = &doc.Image
	}

	var relatedID *string
	if doc.RelatedID != "" {
		relatedID = &doc.RelatedID
	}

	var errMsg *string
	if doc.Error != "" {
		errMsg = &doc.Error
	}

	return domain.ReconstructNotification(
		doc.ID.Hex(),
		doc.UserID,
		domain.NotificationChannel(doc.Channel),
		domain.NotificationType(doc.Type),
		doc.Title,
		doc.Body,
		imageURL,
		data,
		domain.NotificationStatus(doc.Status),
		doc.SentAt,
		doc.DeliveredAt,
		doc.ReadAt,
		errMsg,
		doc.RelatedType,
		relatedID,
		doc.CreatedAt,
	)
}
