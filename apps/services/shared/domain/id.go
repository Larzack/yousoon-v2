package domain

import (
	"encoding/json"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ErrInvalidID is returned when an ID string is invalid.
var ErrInvalidID = errors.New("invalid ID format")

// ID is the base interface for all strongly-typed IDs.
type ID interface {
	String() string
	IsZero() bool
	ObjectID() primitive.ObjectID
}

// BaseID provides a common implementation for strongly-typed IDs.
type BaseID struct {
	value primitive.ObjectID
}

// NewBaseID creates a new BaseID with a generated ObjectID.
func NewBaseID() BaseID {
	return BaseID{value: primitive.NewObjectID()}
}

// BaseIDFromString creates a BaseID from a hex string.
func BaseIDFromString(s string) (BaseID, error) {
	if s == "" {
		return BaseID{}, nil
	}
	oid, err := primitive.ObjectIDFromHex(s)
	if err != nil {
		return BaseID{}, fmt.Errorf("%w: %s", ErrInvalidID, s)
	}
	return BaseID{value: oid}, nil
}

// BaseIDFromObjectID creates a BaseID from a primitive.ObjectID.
func BaseIDFromObjectID(oid primitive.ObjectID) BaseID {
	return BaseID{value: oid}
}

// String returns the hex string representation.
func (id BaseID) String() string {
	return id.value.Hex()
}

// IsZero returns true if the ID is zero/empty.
func (id BaseID) IsZero() bool {
	return id.value.IsZero()
}

// ObjectID returns the underlying MongoDB ObjectID.
func (id BaseID) ObjectID() primitive.ObjectID {
	return id.value
}

// MarshalJSON implements json.Marshaler.
func (id BaseID) MarshalJSON() ([]byte, error) {
	return json.Marshal(id.value.Hex())
}

// UnmarshalJSON implements json.Unmarshaler.
func (id *BaseID) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	parsed, err := BaseIDFromString(s)
	if err != nil {
		return err
	}
	*id = parsed
	return nil
}

// =============================================================================
// Strongly-typed IDs for each aggregate
// =============================================================================

// UserID identifies a User aggregate.
type UserID struct{ BaseID }

// NewUserID creates a new UserID.
func NewUserID() UserID {
	return UserID{NewBaseID()}
}

// UserIDFromString creates a UserID from a string.
func UserIDFromString(s string) (UserID, error) {
	base, err := BaseIDFromString(s)
	if err != nil {
		return UserID{}, err
	}
	return UserID{base}, nil
}

// MustUserID creates a UserID from string, panics on error.
func MustUserID(s string) UserID {
	id, err := UserIDFromString(s)
	if err != nil {
		panic(err)
	}
	return id
}

// PartnerID identifies a Partner aggregate.
type PartnerID struct{ BaseID }

// NewPartnerID creates a new PartnerID.
func NewPartnerID() PartnerID {
	return PartnerID{NewBaseID()}
}

// PartnerIDFromString creates a PartnerID from a string.
func PartnerIDFromString(s string) (PartnerID, error) {
	base, err := BaseIDFromString(s)
	if err != nil {
		return PartnerID{}, err
	}
	return PartnerID{base}, nil
}

// EstablishmentID identifies an Establishment entity.
type EstablishmentID struct{ BaseID }

// NewEstablishmentID creates a new EstablishmentID.
func NewEstablishmentID() EstablishmentID {
	return EstablishmentID{NewBaseID()}
}

// EstablishmentIDFromString creates an EstablishmentID from a string.
func EstablishmentIDFromString(s string) (EstablishmentID, error) {
	base, err := BaseIDFromString(s)
	if err != nil {
		return EstablishmentID{}, err
	}
	return EstablishmentID{base}, nil
}

// OfferID identifies an Offer aggregate.
type OfferID struct{ BaseID }

// NewOfferID creates a new OfferID.
func NewOfferID() OfferID {
	return OfferID{NewBaseID()}
}

// OfferIDFromString creates an OfferID from a string.
func OfferIDFromString(s string) (OfferID, error) {
	base, err := BaseIDFromString(s)
	if err != nil {
		return OfferID{}, err
	}
	return OfferID{base}, nil
}

// CategoryID identifies a Category aggregate.
type CategoryID struct{ BaseID }

// NewCategoryID creates a new CategoryID.
func NewCategoryID() CategoryID {
	return CategoryID{NewBaseID()}
}

// CategoryIDFromString creates a CategoryID from a string.
func CategoryIDFromString(s string) (CategoryID, error) {
	base, err := BaseIDFromString(s)
	if err != nil {
		return CategoryID{}, err
	}
	return CategoryID{base}, nil
}

// OutingID identifies an Outing (booking) aggregate.
type OutingID struct{ BaseID }

// NewOutingID creates a new OutingID.
func NewOutingID() OutingID {
	return OutingID{NewBaseID()}
}

// OutingIDFromString creates an OutingID from a string.
func OutingIDFromString(s string) (OutingID, error) {
	base, err := BaseIDFromString(s)
	if err != nil {
		return OutingID{}, err
	}
	return OutingID{base}, nil
}

// ReviewID identifies a Review aggregate.
type ReviewID struct{ BaseID }

// NewReviewID creates a new ReviewID.
func NewReviewID() ReviewID {
	return ReviewID{NewBaseID()}
}

// ReviewIDFromString creates a ReviewID from a string.
func ReviewIDFromString(s string) (ReviewID, error) {
	base, err := BaseIDFromString(s)
	if err != nil {
		return ReviewID{}, err
	}
	return ReviewID{base}, nil
}

// ConversationID identifies a Conversation aggregate.
type ConversationID struct{ BaseID }

// NewConversationID creates a new ConversationID.
func NewConversationID() ConversationID {
	return ConversationID{NewBaseID()}
}

// ConversationIDFromString creates a ConversationID from a string.
func ConversationIDFromString(s string) (ConversationID, error) {
	base, err := BaseIDFromString(s)
	if err != nil {
		return ConversationID{}, err
	}
	return ConversationID{base}, nil
}

// MessageID identifies a Message entity.
type MessageID struct{ BaseID }

// NewMessageID creates a new MessageID.
func NewMessageID() MessageID {
	return MessageID{NewBaseID()}
}

// MessageIDFromString creates a MessageID from a string.
func MessageIDFromString(s string) (MessageID, error) {
	base, err := BaseIDFromString(s)
	if err != nil {
		return MessageID{}, err
	}
	return MessageID{base}, nil
}

// NotificationID identifies a Notification aggregate.
type NotificationID struct{ BaseID }

// NewNotificationID creates a new NotificationID.
func NewNotificationID() NotificationID {
	return NotificationID{NewBaseID()}
}

// NotificationIDFromString creates a NotificationID from a string.
func NotificationIDFromString(s string) (NotificationID, error) {
	base, err := BaseIDFromString(s)
	if err != nil {
		return NotificationID{}, err
	}
	return NotificationID{base}, nil
}

// SubscriptionID identifies a Subscription entity.
type SubscriptionID struct{ BaseID }

// NewSubscriptionID creates a new SubscriptionID.
func NewSubscriptionID() SubscriptionID {
	return SubscriptionID{NewBaseID()}
}

// SubscriptionIDFromString creates a SubscriptionID from a string.
func SubscriptionIDFromString(s string) (SubscriptionID, error) {
	base, err := BaseIDFromString(s)
	if err != nil {
		return SubscriptionID{}, err
	}
	return SubscriptionID{base}, nil
}

// VerificationID identifies an IdentityVerification entity.
type VerificationID struct{ BaseID }

// NewVerificationID creates a new VerificationID.
func NewVerificationID() VerificationID {
	return VerificationID{NewBaseID()}
}

// VerificationIDFromString creates a VerificationID from a string.
func VerificationIDFromString(s string) (VerificationID, error) {
	base, err := BaseIDFromString(s)
	if err != nil {
		return VerificationID{}, err
	}
	return VerificationID{base}, nil
}

// TeamMemberID identifies a TeamMember entity.
type TeamMemberID struct{ BaseID }

// NewTeamMemberID creates a new TeamMemberID.
func NewTeamMemberID() TeamMemberID {
	return TeamMemberID{NewBaseID()}
}

// TeamMemberIDFromString creates a TeamMemberID from a string.
func TeamMemberIDFromString(s string) (TeamMemberID, error) {
	base, err := BaseIDFromString(s)
	if err != nil {
		return TeamMemberID{}, err
	}
	return TeamMemberID{base}, nil
}
