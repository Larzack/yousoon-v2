package mongodb

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/yousoon/services/identity/internal/domain"
)

const usersCollection = "users"

// UserRepository implements domain.UserRepository using MongoDB.
type UserRepository struct {
	db *mongo.Database
}

// NewUserRepository creates a new UserRepository.
func NewUserRepository(db *mongo.Database) *UserRepository {
	return &UserRepository{db: db}
}

// userDocument represents the MongoDB document structure for a user.
type userDocument struct {
	ID             primitive.ObjectID `bson:"_id,omitempty"`
	Email          string             `bson:"email"`
	PasswordHash   string             `bson:"passwordHash"`
	Phone          *string            `bson:"phone,omitempty"`
	Profile        profileDocument    `bson:"profile"`
	Preferences    prefsDocument      `bson:"preferences"`
	Identity       *identityDocument  `bson:"identity,omitempty"`
	LastLocation   *locationDocument  `bson:"lastLocation,omitempty"`
	FCMTokens      []fcmTokenDocument `bson:"fcmTokens"`
	SocialAccounts []socialDocument   `bson:"socialAccounts"`
	SubscriptionID *string            `bson:"subscriptionId,omitempty"`
	Status         string             `bson:"status"`
	Grade          int                `bson:"grade"`
	EmailVerified  bool               `bson:"emailVerified"`
	PhoneVerified  bool               `bson:"phoneVerified"`
	CreatedAt      time.Time          `bson:"createdAt"`
	UpdatedAt      time.Time          `bson:"updatedAt"`
	LastLoginAt    *time.Time         `bson:"lastLoginAt,omitempty"`
	DeletedAt      *time.Time         `bson:"deletedAt,omitempty"`
}

type profileDocument struct {
	FirstName   string     `bson:"firstName"`
	LastName    string     `bson:"lastName"`
	DisplayName string     `bson:"displayName"`
	Avatar      *string    `bson:"avatar,omitempty"`
	BirthDate   *time.Time `bson:"birthDate,omitempty"`
	Gender      *string    `bson:"gender,omitempty"`
}

type prefsDocument struct {
	Language      string   `bson:"language"`
	Notifications notifDoc `bson:"notifications"`
	Categories    []string `bson:"categories"`
	MaxDistance   int      `bson:"maxDistance"`
}

type notifDoc struct {
	Push      bool `bson:"push"`
	Email     bool `bson:"email"`
	SMS       bool `bson:"sms"`
	Marketing bool `bson:"marketing"`
}

type identityDocument struct {
	ID              string     `bson:"id"`
	Status          string     `bson:"status"`
	DocumentType    string     `bson:"documentType"`
	Method          string     `bson:"method"`
	AttemptCount    int        `bson:"attemptCount"`
	SubmittedAt     time.Time  `bson:"submittedAt"`
	VerifiedAt      *time.Time `bson:"verifiedAt,omitempty"`
	RejectedAt      *time.Time `bson:"rejectedAt,omitempty"`
	RejectionReason *string    `bson:"rejectionReason,omitempty"`
}

type locationDocument struct {
	Type        string    `bson:"type"`
	Coordinates []float64 `bson:"coordinates"`
	UpdatedAt   time.Time `bson:"updatedAt"`
}

type fcmTokenDocument struct {
	Token    string    `bson:"token"`
	Platform string    `bson:"platform"`
	AddedAt  time.Time `bson:"addedAt"`
}

type socialDocument struct {
	Provider   string    `bson:"provider"`
	ProviderID string    `bson:"providerId"`
	Email      string    `bson:"email,omitempty"`
	LinkedAt   time.Time `bson:"linkedAt"`
}

// collection returns the users collection.
func (r *UserRepository) collection() *mongo.Collection {
	return r.db.Collection(usersCollection)
}

// Create creates a new user.
func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
	doc := toUserDocument(user)
	doc.ID = primitive.NewObjectID()

	_, err := r.collection().InsertOne(ctx, doc)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return domain.ErrEmailAlreadyExists
		}
		return err
	}

	return nil
}

// Update updates an existing user.
func (r *UserRepository) Update(ctx context.Context, user *domain.User) error {
	doc := toUserDocument(user)
	doc.UpdatedAt = time.Now()

	filter := bson.M{"email": user.Email.String()}
	update := bson.M{"$set": doc}

	result, err := r.collection().UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return domain.ErrUserNotFound
	}

	return nil
}

// Delete permanently deletes a user.
func (r *UserRepository) Delete(ctx context.Context, id domain.UserID) error {
	oid, err := primitive.ObjectIDFromHex(id.String())
	if err != nil {
		return domain.ErrUserNotFound
	}

	filter := bson.M{"_id": oid}
	result, err := r.collection().DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return domain.ErrUserNotFound
	}

	return nil
}

// FindByID finds a user by ID.
func (r *UserRepository) FindByID(ctx context.Context, id domain.UserID) (*domain.User, error) {
	oid, err := primitive.ObjectIDFromHex(id.String())
	if err != nil {
		return nil, domain.ErrUserNotFound
	}

	filter := bson.M{"_id": oid}
	var doc userDocument
	if err := r.collection().FindOne(ctx, filter).Decode(&doc); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}

	return toUserDomain(&doc), nil
}

// FindByEmail finds a user by email.
func (r *UserRepository) FindByEmail(ctx context.Context, email domain.Email) (*domain.User, error) {
	filter := bson.M{"email": email.String()}
	var doc userDocument
	if err := r.collection().FindOne(ctx, filter).Decode(&doc); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}

	return toUserDomain(&doc), nil
}

// FindByPhone finds a user by phone.
func (r *UserRepository) FindByPhone(ctx context.Context, phone domain.Phone) (*domain.User, error) {
	filter := bson.M{"phone": phone.String()}
	var doc userDocument
	if err := r.collection().FindOne(ctx, filter).Decode(&doc); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}

	return toUserDomain(&doc), nil
}

// FindBySocialAccount finds a user by social account.
func (r *UserRepository) FindBySocialAccount(ctx context.Context, provider domain.SocialProvider, providerID string) (*domain.User, error) {
	filter := bson.M{
		"socialAccounts": bson.M{
			"$elemMatch": bson.M{
				"provider":   string(provider),
				"providerId": providerID,
			},
		},
	}
	var doc userDocument
	if err := r.collection().FindOne(ctx, filter).Decode(&doc); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}

	return toUserDomain(&doc), nil
}

// ExistsByEmail checks if a user with the email exists.
func (r *UserRepository) ExistsByEmail(ctx context.Context, email domain.Email) (bool, error) {
	filter := bson.M{"email": email.String()}
	count, err := r.collection().CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// ExistsByPhone checks if a user with the phone exists.
func (r *UserRepository) ExistsByPhone(ctx context.Context, phone domain.Phone) (bool, error) {
	filter := bson.M{"phone": phone.String()}
	count, err := r.collection().CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// FindUsersForDeletion finds users scheduled for deletion.
func (r *UserRepository) FindUsersForDeletion(ctx context.Context, beforeTimestamp int64) ([]*domain.User, error) {
	before := time.Unix(beforeTimestamp, 0)
	filter := bson.M{
		"deletedAt": bson.M{"$lte": before},
		"status":    string(domain.UserStatusDeleted),
	}

	cursor, err := r.collection().Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []*domain.User
	for cursor.Next(ctx) {
		var doc userDocument
		if err := cursor.Decode(&doc); err != nil {
			return nil, err
		}
		users = append(users, toUserDomain(&doc))
	}

	return users, cursor.Err()
}

// Count returns the total number of users.
func (r *UserRepository) Count(ctx context.Context) (int64, error) {
	return r.collection().CountDocuments(ctx, bson.M{})
}

// CountActive returns the number of active users.
func (r *UserRepository) CountActive(ctx context.Context) (int64, error) {
	filter := bson.M{"status": string(domain.UserStatusActive)}
	return r.collection().CountDocuments(ctx, filter)
}

// EnsureIndexes creates the necessary indexes.
func (r *UserRepository) EnsureIndexes(ctx context.Context) error {
	indexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "email", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys:    bson.D{{Key: "phone", Value: 1}},
			Options: options.Index().SetSparse(true),
		},
		{
			Keys: bson.D{{Key: "lastLocation", Value: "2dsphere"}},
		},
		{
			Keys: bson.D{{Key: "status", Value: 1}},
		},
		{
			Keys: bson.D{
				{Key: "socialAccounts.provider", Value: 1},
				{Key: "socialAccounts.providerId", Value: 1},
			},
		},
	}

	_, err := r.collection().Indexes().CreateMany(ctx, indexes)
	return err
}

// toUserDocument converts domain user to MongoDB document.
func toUserDocument(user *domain.User) *userDocument {
	doc := &userDocument{
		Email:         user.Email.String(),
		PasswordHash:  user.PasswordHash,
		Status:        string(user.Status),
		Grade:         int(user.Grade),
		EmailVerified: user.EmailVerified,
		PhoneVerified: user.PhoneVerified,
		CreatedAt:     user.AggregateRoot.CreatedAt,
		UpdatedAt:     user.AggregateRoot.UpdatedAt,
		LastLoginAt:   user.LastLoginAt,
		DeletedAt:     user.DeletedAt,
		Profile: profileDocument{
			FirstName:   user.Profile.FirstName,
			LastName:    user.Profile.LastName,
			DisplayName: user.Profile.DisplayName,
			Avatar:      user.Profile.Avatar,
			BirthDate:   user.Profile.BirthDate,
		},
		Preferences: prefsDocument{
			Language:    user.Preferences.Language,
			MaxDistance: user.Preferences.MaxDistance,
			Notifications: notifDoc{
				Push:      user.Preferences.NotificationSettings.Push,
				Email:     user.Preferences.NotificationSettings.Email,
				SMS:       user.Preferences.NotificationSettings.SMS,
				Marketing: user.Preferences.NotificationSettings.Marketing,
			},
			Categories: user.Preferences.FavoriteCategories,
		},
	}

	if user.Phone != nil {
		phone := user.Phone.String()
		doc.Phone = &phone
	}

	if user.Profile.Gender != nil {
		gender := string(*user.Profile.Gender)
		doc.Profile.Gender = &gender
	}

	if user.SubscriptionID != nil {
		subID := user.SubscriptionID.String()
		doc.SubscriptionID = &subID
	}

	if user.Identity != nil {
		doc.Identity = &identityDocument{
			ID:              user.Identity.ID.String(),
			Status:          string(user.Identity.Status),
			DocumentType:    string(user.Identity.DocumentType),
			Method:          string(user.Identity.Method),
			AttemptCount:    user.Identity.AttemptCount,
			SubmittedAt:     user.Identity.SubmittedAt,
			VerifiedAt:      user.Identity.VerifiedAt,
			RejectedAt:      user.Identity.RejectedAt,
			RejectionReason: user.Identity.RejectionReason,
		}
	}

	if user.LastLocation != nil {
		doc.LastLocation = &locationDocument{
			Type:        user.LastLocation.Type,
			Coordinates: user.LastLocation.Coordinates,
			UpdatedAt:   user.LastLocation.UpdatedAt,
		}
	}

	for _, token := range user.FCMTokens {
		doc.FCMTokens = append(doc.FCMTokens, fcmTokenDocument{
			Token:    token.Token,
			Platform: string(token.Platform),
			AddedAt:  token.AddedAt,
		})
	}

	for _, social := range user.SocialAccounts {
		doc.SocialAccounts = append(doc.SocialAccounts, socialDocument{
			Provider:   string(social.Provider),
			ProviderID: social.ProviderID,
			Email:      social.Email,
			LinkedAt:   social.LinkedAt,
		})
	}

	return doc
}

// toUserDomain converts MongoDB document to domain user.
func toUserDomain(doc *userDocument) *domain.User {
	email, _ := domain.NewEmail(doc.Email)

	user := &domain.User{
		ID:            domain.UserID(doc.ID.Hex()),
		Email:         email,
		PasswordHash:  doc.PasswordHash,
		Status:        domain.UserStatus(doc.Status),
		Grade:         domain.UserGrade(doc.Grade),
		EmailVerified: doc.EmailVerified,
		PhoneVerified: doc.PhoneVerified,
		LastLoginAt:   doc.LastLoginAt,
		DeletedAt:     doc.DeletedAt,
		Profile: domain.Profile{
			FirstName:   doc.Profile.FirstName,
			LastName:    doc.Profile.LastName,
			DisplayName: doc.Profile.DisplayName,
			Avatar:      doc.Profile.Avatar,
			BirthDate:   doc.Profile.BirthDate,
		},
		Preferences: domain.Preferences{
			Language:    doc.Preferences.Language,
			MaxDistance: doc.Preferences.MaxDistance,
			NotificationSettings: domain.NotificationSettings{
				Push:      doc.Preferences.Notifications.Push,
				Email:     doc.Preferences.Notifications.Email,
				SMS:       doc.Preferences.Notifications.SMS,
				Marketing: doc.Preferences.Notifications.Marketing,
			},
			FavoriteCategories: doc.Preferences.Categories,
		},
	}

	// Set AggregateRoot timestamps
	user.AggregateRoot.CreatedAt = doc.CreatedAt
	user.AggregateRoot.UpdatedAt = doc.UpdatedAt

	if doc.Phone != nil {
		phone, _ := domain.NewPhone(*doc.Phone)
		user.Phone = &phone
	}

	if doc.Profile.Gender != nil {
		gender := domain.Gender(*doc.Profile.Gender)
		user.Profile.Gender = &gender
	}

	if doc.SubscriptionID != nil {
		subID := domain.SubscriptionID(*doc.SubscriptionID)
		user.SubscriptionID = &subID
	}

	if doc.Identity != nil {
		user.Identity = &domain.IdentityVerification{
			ID:              domain.VerificationID(doc.Identity.ID),
			Status:          domain.VerificationStatus(doc.Identity.Status),
			DocumentType:    domain.DocumentType(doc.Identity.DocumentType),
			Method:          domain.VerificationMethod(doc.Identity.Method),
			AttemptCount:    doc.Identity.AttemptCount,
			SubmittedAt:     doc.Identity.SubmittedAt,
			VerifiedAt:      doc.Identity.VerifiedAt,
			RejectedAt:      doc.Identity.RejectedAt,
			RejectionReason: doc.Identity.RejectionReason,
		}
	}

	if doc.LastLocation != nil {
		user.LastLocation = &domain.GeoLocation{
			Type:        doc.LastLocation.Type,
			Coordinates: doc.LastLocation.Coordinates,
			UpdatedAt:   doc.LastLocation.UpdatedAt,
		}
	}

	for _, token := range doc.FCMTokens {
		user.FCMTokens = append(user.FCMTokens, domain.FCMToken{
			Token:    token.Token,
			Platform: domain.Platform(token.Platform),
			AddedAt:  token.AddedAt,
		})
	}

	for _, social := range doc.SocialAccounts {
		user.SocialAccounts = append(user.SocialAccounts, domain.SocialAccount{
			Provider:   domain.SocialProvider(social.Provider),
			ProviderID: social.ProviderID,
			Email:      social.Email,
			LinkedAt:   social.LinkedAt,
		})
	}

	return user
}
