package domain

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
)

// =============================================================================
// User ID
// =============================================================================

// UserID is a unique identifier for a user.
type UserID string

// NewUserID creates a new UserID.
func NewUserID() UserID {
	return UserID(uuid.New().String())
}

// ParseUserID parses a string into a UserID.
func ParseUserID(s string) (UserID, error) {
	if _, err := uuid.Parse(s); err != nil {
		return "", fmt.Errorf("invalid user ID: %w", err)
	}
	return UserID(s), nil
}

// String returns the string representation.
func (id UserID) String() string {
	return string(id)
}

// IsZero returns true if the ID is empty.
func (id UserID) IsZero() bool {
	return id == ""
}

// =============================================================================
// Subscription ID
// =============================================================================

// SubscriptionID is a unique identifier for a subscription.
type SubscriptionID string

// NewSubscriptionID creates a new SubscriptionID.
func NewSubscriptionID() SubscriptionID {
	return SubscriptionID(uuid.New().String())
}

// String returns the string representation.
func (id SubscriptionID) String() string {
	return string(id)
}

// =============================================================================
// Email Value Object
// =============================================================================

// Email represents a validated email address.
type Email struct {
	value string
}

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

// NewEmail creates a new validated Email.
func NewEmail(email string) (Email, error) {
	email = strings.TrimSpace(strings.ToLower(email))
	if !emailRegex.MatchString(email) {
		return Email{}, fmt.Errorf("invalid email format: %s", email)
	}
	return Email{value: email}, nil
}

// String returns the email string.
func (e Email) String() string {
	return e.value
}

// Equals checks if two emails are equal.
func (e Email) Equals(other Email) bool {
	return e.value == other.value
}

// Domain returns the email domain.
func (e Email) Domain() string {
	parts := strings.Split(e.value, "@")
	if len(parts) == 2 {
		return parts[1]
	}
	return ""
}

// MarshalJSON implements json.Marshaler.
func (e Email) MarshalJSON() ([]byte, error) {
	return []byte(`"` + e.value + `"`), nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (e *Email) UnmarshalJSON(data []byte) error {
	s := strings.Trim(string(data), `"`)
	email, err := NewEmail(s)
	if err != nil {
		return err
	}
	*e = email
	return nil
}

// =============================================================================
// Phone Value Object
// =============================================================================

// Phone represents a validated phone number in E.164 format.
type Phone struct {
	value string
}

var phoneRegex = regexp.MustCompile(`^\+[1-9]\d{1,14}$`)

// NewPhone creates a new validated Phone.
func NewPhone(phone string) (Phone, error) {
	phone = strings.TrimSpace(phone)
	if !phoneRegex.MatchString(phone) {
		return Phone{}, fmt.Errorf("invalid phone format (E.164 required): %s", phone)
	}
	return Phone{value: phone}, nil
}

// String returns the phone string.
func (p Phone) String() string {
	return p.value
}

// Equals checks if two phones are equal.
func (p Phone) Equals(other Phone) bool {
	return p.value == other.value
}

// MarshalJSON implements json.Marshaler.
func (p Phone) MarshalJSON() ([]byte, error) {
	return []byte(`"` + p.value + `"`), nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (p *Phone) UnmarshalJSON(data []byte) error {
	s := strings.Trim(string(data), `"`)
	if s == "" {
		return nil
	}
	phone, err := NewPhone(s)
	if err != nil {
		return err
	}
	*p = phone
	return nil
}

// =============================================================================
// Password Value Object
// =============================================================================

// Password represents a hashed password.
type Password struct {
	hash string
}

// NewPassword creates a new password from plain text.
// It hashes the password using bcrypt.
func NewPassword(plaintext string) (Password, error) {
	if len(plaintext) < 8 {
		return Password{}, ErrPasswordTooShort
	}

	// Check password strength (at least one uppercase, one lowercase, one digit)
	hasUpper := false
	hasLower := false
	hasDigit := false
	for _, c := range plaintext {
		switch {
		case c >= 'A' && c <= 'Z':
			hasUpper = true
		case c >= 'a' && c <= 'z':
			hasLower = true
		case c >= '0' && c <= '9':
			hasDigit = true
		}
	}
	if !hasUpper || !hasLower || !hasDigit {
		return Password{}, ErrPasswordTooWeak
	}

	hash, err := bcryptHash(plaintext)
	if err != nil {
		return Password{}, fmt.Errorf("failed to hash password: %w", err)
	}

	return Password{hash: hash}, nil
}

// NewPasswordFromHash creates a Password from an existing hash.
// Use this when loading from database.
func NewPasswordFromHash(hash string) Password {
	return Password{hash: hash}
}

// Hash returns the password hash.
func (p Password) Hash() string {
	return p.hash
}

// Matches checks if the plain text password matches the hash.
func (p Password) Matches(plaintext string) bool {
	return bcryptCompare(p.hash, plaintext)
}

// bcryptHash hashes a password using bcrypt.
func bcryptHash(password string) (string, error) {
	// Using golang.org/x/crypto/bcrypt
	// For now, we'll use a simple implementation
	// In production, import bcrypt package
	// hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	// return string(hash), err

	// Placeholder - in real implementation use bcrypt
	return "hashed:" + password, nil
}

// bcryptCompare compares a password with a hash.
func bcryptCompare(hash, password string) bool {
	// err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	// return err == nil

	// Placeholder - in real implementation use bcrypt
	return hash == "hashed:"+password
}

// =============================================================================
// Profile Value Object
// =============================================================================

// Profile contains user profile information.
type Profile struct {
	FirstName   string     `json:"firstName" bson:"firstName"`
	LastName    string     `json:"lastName" bson:"lastName"`
	DisplayName string     `json:"displayName" bson:"displayName"`
	Avatar      *string    `json:"avatar,omitempty" bson:"avatar,omitempty"`
	BirthDate   *time.Time `json:"birthDate,omitempty" bson:"birthDate,omitempty"`
	Gender      *Gender    `json:"gender,omitempty" bson:"gender,omitempty"`
}

// NewProfile creates a new profile.
func NewProfile(firstName, lastName string) Profile {
	return Profile{
		FirstName:   firstName,
		LastName:    lastName,
		DisplayName: fmt.Sprintf("%s %s", firstName, lastName),
	}
}

// FullName returns the full name.
func (p Profile) FullName() string {
	return fmt.Sprintf("%s %s", p.FirstName, p.LastName)
}

// Age returns the age in years, or 0 if birth date is not set.
func (p Profile) Age() int {
	if p.BirthDate == nil {
		return 0
	}
	now := time.Now()
	age := now.Year() - p.BirthDate.Year()
	if now.YearDay() < p.BirthDate.YearDay() {
		age--
	}
	return age
}

// Gender represents user gender.
type Gender string

const (
	GenderMale   Gender = "male"
	GenderFemale Gender = "female"
	GenderOther  Gender = "other"
)

// =============================================================================
// Preferences Value Object
// =============================================================================

// Preferences contains user preferences.
type Preferences struct {
	Language             string               `json:"language" bson:"language"`
	NotificationSettings NotificationSettings `json:"notifications" bson:"notifications"`
	FavoriteCategories   []string             `json:"categories,omitempty" bson:"categories,omitempty"`
	MaxDistance          int                  `json:"maxDistance" bson:"maxDistance"` // in km
}

// NotificationSettings contains notification preferences.
type NotificationSettings struct {
	Push      bool `json:"push" bson:"push"`
	Email     bool `json:"email" bson:"email"`
	SMS       bool `json:"sms" bson:"sms"`
	Marketing bool `json:"marketing" bson:"marketing"`
}

// DefaultPreferences returns default preferences.
func DefaultPreferences() Preferences {
	return Preferences{
		Language: "fr",
		NotificationSettings: NotificationSettings{
			Push:      true,
			Email:     true,
			SMS:       false,
			Marketing: false,
		},
		MaxDistance: 10, // 10 km default
	}
}

// =============================================================================
// GeoLocation Value Object
// =============================================================================

// GeoLocation represents a geographic location.
type GeoLocation struct {
	Type        string    `json:"type" bson:"type"`
	Coordinates []float64 `json:"coordinates" bson:"coordinates"` // [longitude, latitude]
	UpdatedAt   time.Time `json:"updatedAt" bson:"updatedAt"`
}

// NewGeoLocation creates a new GeoLocation.
func NewGeoLocation(longitude, latitude float64) (GeoLocation, error) {
	if longitude < -180 || longitude > 180 {
		return GeoLocation{}, fmt.Errorf("invalid longitude: %f", longitude)
	}
	if latitude < -90 || latitude > 90 {
		return GeoLocation{}, fmt.Errorf("invalid latitude: %f", latitude)
	}
	return GeoLocation{
		Type:        "Point",
		Coordinates: []float64{longitude, latitude},
		UpdatedAt:   time.Now(),
	}, nil
}

// Longitude returns the longitude.
func (g GeoLocation) Longitude() float64 {
	if len(g.Coordinates) >= 2 {
		return g.Coordinates[0]
	}
	return 0
}

// Latitude returns the latitude.
func (g GeoLocation) Latitude() float64 {
	if len(g.Coordinates) >= 2 {
		return g.Coordinates[1]
	}
	return 0
}

// =============================================================================
// User Status
// =============================================================================

// UserStatus represents the status of a user account.
type UserStatus string

const (
	UserStatusActive    UserStatus = "active"
	UserStatusSuspended UserStatus = "suspended"
	UserStatusDeleted   UserStatus = "deleted"
)

// =============================================================================
// User Grade
// =============================================================================

// UserGrade represents the user's loyalty grade.
type UserGrade int

const (
	GradeExplorateur   UserGrade = 1
	GradeAventurier    UserGrade = 2
	GradeGrandVoyageur UserGrade = 3
	GradeConquerant    UserGrade = 4
)

// String returns the grade name.
func (g UserGrade) String() string {
	switch g {
	case GradeExplorateur:
		return "explorateur"
	case GradeAventurier:
		return "aventurier"
	case GradeGrandVoyageur:
		return "grand_voyageur"
	case GradeConquerant:
		return "conquerant"
	default:
		return "unknown"
	}
}

// =============================================================================
// Identity Verification
// =============================================================================

// VerificationID is a unique identifier for a verification.
type VerificationID string

// NewVerificationID creates a new VerificationID.
func NewVerificationID() VerificationID {
	return VerificationID(uuid.New().String())
}

// String returns the string representation.
func (id VerificationID) String() string {
	return string(id)
}

// IdentityVerification represents the identity verification state.
type IdentityVerification struct {
	ID              VerificationID     `json:"id" bson:"id"`
	Status          VerificationStatus `json:"status" bson:"status"`
	DocumentType    DocumentType       `json:"documentType" bson:"documentType"`
	Method          VerificationMethod `json:"method" bson:"method"`
	AttemptCount    int                `json:"attemptCount" bson:"attemptCount"`
	SubmittedAt     time.Time          `json:"submittedAt" bson:"submittedAt"`
	VerifiedAt      *time.Time         `json:"verifiedAt,omitempty" bson:"verifiedAt,omitempty"`
	RejectedAt      *time.Time         `json:"rejectedAt,omitempty" bson:"rejectedAt,omitempty"`
	RejectionReason *string            `json:"rejectionReason,omitempty" bson:"rejectionReason,omitempty"`
}

// NewIdentityVerification creates a new identity verification.
func NewIdentityVerification(docType DocumentType, method VerificationMethod) IdentityVerification {
	return IdentityVerification{
		ID:           NewVerificationID(),
		Status:       VerificationStatusPending,
		DocumentType: docType,
		Method:       method,
		AttemptCount: 1,
		SubmittedAt:  time.Now(),
	}
}

// VerificationStatus represents the status of identity verification.
type VerificationStatus string

const (
	VerificationStatusNotSubmitted VerificationStatus = "not_submitted"
	VerificationStatusPending      VerificationStatus = "pending"
	VerificationStatusVerified     VerificationStatus = "verified"
	VerificationStatusRejected     VerificationStatus = "rejected"
)

// DocumentType represents the type of identity document.
type DocumentType string

const (
	DocumentTypeCNI            DocumentType = "cni"
	DocumentTypePassport       DocumentType = "passport"
	DocumentTypeDrivingLicense DocumentType = "driving_license"
)

// VerificationMethod represents how verification is performed.
type VerificationMethod string

const (
	VerificationMethodInternalOCR VerificationMethod = "internal_ocr"
	VerificationMethodExternal    VerificationMethod = "external"
)

// =============================================================================
// Social Account
// =============================================================================

// SocialAccount represents a linked social account.
type SocialAccount struct {
	Provider   SocialProvider `json:"provider" bson:"provider"`
	ProviderID string         `json:"providerId" bson:"providerId"`
	Email      string         `json:"email,omitempty" bson:"email,omitempty"`
	LinkedAt   time.Time      `json:"linkedAt" bson:"linkedAt"`
}

// SocialProvider represents a social login provider.
type SocialProvider string

const (
	SocialProviderGoogle   SocialProvider = "google"
	SocialProviderApple    SocialProvider = "apple"
	SocialProviderFacebook SocialProvider = "facebook"
)

// =============================================================================
// FCM Token
// =============================================================================

// FCMToken represents a Firebase Cloud Messaging token.
type FCMToken struct {
	Token    string    `json:"token" bson:"token"`
	Platform Platform  `json:"platform" bson:"platform"`
	AddedAt  time.Time `json:"addedAt" bson:"addedAt"`
}

// NewFCMToken creates a new FCM token.
func NewFCMToken(token string, platform Platform) FCMToken {
	return FCMToken{
		Token:    token,
		Platform: platform,
		AddedAt:  time.Now(),
	}
}

// Platform represents the mobile platform.
type Platform string

const (
	PlatformIOS     Platform = "ios"
	PlatformAndroid Platform = "android"
	PlatformWeb     Platform = "web"
)
