package mongodb

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// BaseDocument represents common MongoDB document fields.
type BaseDocument struct {
	ID        primitive.ObjectID  `bson:"_id,omitempty"`
	CreatedAt primitive.DateTime  `bson:"createdAt"`
	UpdatedAt primitive.DateTime  `bson:"updatedAt"`
	DeletedAt *primitive.DateTime `bson:"deletedAt,omitempty"`
}

// NewBaseDocument creates a new base document with current timestamps.
func NewBaseDocument() BaseDocument {
	now := primitive.NewDateTimeFromTime(time.Now())
	return BaseDocument{
		ID:        primitive.NewObjectID(),
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// NewBaseDocumentWithID creates a new base document with a specific ID.
func NewBaseDocumentWithID(id primitive.ObjectID) BaseDocument {
	now := primitive.NewDateTimeFromTime(time.Now())
	return BaseDocument{
		ID:        id,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// Touch updates the UpdatedAt timestamp.
func (bd *BaseDocument) Touch() {
	bd.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())
}

// MarkDeleted sets the DeletedAt timestamp for soft delete.
func (bd *BaseDocument) MarkDeleted() {
	now := primitive.NewDateTimeFromTime(time.Now())
	bd.DeletedAt = &now
}

// IsDeleted returns true if the document has been soft deleted.
func (bd *BaseDocument) IsDeleted() bool {
	return bd.DeletedAt != nil
}

// ToTime converts a primitive.DateTime to time.Time.
func ToTime(dt primitive.DateTime) time.Time {
	return dt.Time()
}

// FromTime converts a time.Time to primitive.DateTime.
func FromTime(t time.Time) primitive.DateTime {
	return primitive.NewDateTimeFromTime(t)
}

// ToTimePtr converts a *primitive.DateTime to *time.Time.
func ToTimePtr(dt *primitive.DateTime) *time.Time {
	if dt == nil {
		return nil
	}
	t := dt.Time()
	return &t
}

// FromTimePtr converts a *time.Time to *primitive.DateTime.
func FromTimePtr(t *time.Time) *primitive.DateTime {
	if t == nil {
		return nil
	}
	dt := primitive.NewDateTimeFromTime(*t)
	return &dt
}

// GeoJSONPoint represents a GeoJSON point for geospatial queries.
type GeoJSONPoint struct {
	Type        string    `bson:"type"`
	Coordinates []float64 `bson:"coordinates"` // [longitude, latitude]
}

// NewGeoJSONPoint creates a new GeoJSON point.
func NewGeoJSONPoint(longitude, latitude float64) GeoJSONPoint {
	return GeoJSONPoint{
		Type:        "Point",
		Coordinates: []float64{longitude, latitude},
	}
}

// Longitude returns the longitude coordinate.
func (p GeoJSONPoint) Longitude() float64 {
	if len(p.Coordinates) >= 1 {
		return p.Coordinates[0]
	}
	return 0
}

// Latitude returns the latitude coordinate.
func (p GeoJSONPoint) Latitude() float64 {
	if len(p.Coordinates) >= 2 {
		return p.Coordinates[1]
	}
	return 0
}

// NearQuery creates a $near query for geospatial searches.
func NearQuery(longitude, latitude float64, maxDistanceMeters int) bson.M {
	return bson.M{
		"$near": bson.M{
			"$geometry": bson.M{
				"type":        "Point",
				"coordinates": []float64{longitude, latitude},
			},
			"$maxDistance": maxDistanceMeters,
		},
	}
}

// GeoWithinQuery creates a $geoWithin query for radius search.
func GeoWithinQuery(longitude, latitude, radiusMeters float64) bson.M {
	return bson.M{
		"$geoWithin": bson.M{
			"$centerSphere": []interface{}{
				[]float64{longitude, latitude},
				radiusMeters / 6378100, // Convert meters to radians
			},
		},
	}
}

// StringArrayToObjectIDs converts a string array to ObjectID array.
func StringArrayToObjectIDs(ids []string) ([]primitive.ObjectID, error) {
	objectIDs := make([]primitive.ObjectID, len(ids))
	for i, id := range ids {
		oid, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return nil, err
		}
		objectIDs[i] = oid
	}
	return objectIDs, nil
}

// ObjectIDsToStrings converts an ObjectID array to string array.
func ObjectIDsToStrings(ids []primitive.ObjectID) []string {
	strings := make([]string, len(ids))
	for i, id := range ids {
		strings[i] = id.Hex()
	}
	return strings
}

// NotDeletedFilter returns a filter for non-deleted documents.
func NotDeletedFilter() bson.M {
	return bson.M{"deletedAt": nil}
}

// CombineFilters combines multiple filters with $and.
func CombineFilters(filters ...bson.M) bson.M {
	if len(filters) == 0 {
		return bson.M{}
	}
	if len(filters) == 1 {
		return filters[0]
	}
	return bson.M{"$and": filters}
}

// InFilter creates an $in filter.
func InFilter(field string, values interface{}) bson.M {
	return bson.M{field: bson.M{"$in": values}}
}

// RegexFilter creates a case-insensitive regex filter.
func RegexFilter(field, pattern string) bson.M {
	return bson.M{field: bson.M{"$regex": pattern, "$options": "i"}}
}

// DateRangeFilter creates a date range filter.
func DateRangeFilter(field string, from, to *time.Time) bson.M {
	filter := bson.M{}
	if from != nil {
		filter["$gte"] = primitive.NewDateTimeFromTime(*from)
	}
	if to != nil {
		filter["$lte"] = primitive.NewDateTimeFromTime(*to)
	}
	if len(filter) == 0 {
		return bson.M{}
	}
	return bson.M{field: filter}
}

// SortOrder represents sort direction.
type SortOrder int

const (
	// SortAsc is ascending order.
	SortAsc SortOrder = 1
	// SortDesc is descending order.
	SortDesc SortOrder = -1
)

// BuildSort creates a sort document from field-order pairs.
func BuildSort(fields map[string]SortOrder) bson.D {
	sort := bson.D{}
	for field, order := range fields {
		sort = append(sort, bson.E{Key: field, Value: int(order)})
	}
	return sort
}

// LocalizedString represents a multi-language string.
type LocalizedString struct {
	FR string `bson:"fr"`
	EN string `bson:"en"`
}

// Get returns the string for the given language code.
func (ls LocalizedString) Get(lang string) string {
	switch lang {
	case "en":
		return ls.EN
	default:
		return ls.FR
	}
}

// MoneyDocument represents money in MongoDB.
type MoneyDocument struct {
	Amount   int64  `bson:"amount"`   // In cents
	Currency string `bson:"currency"` // ISO 4217
}

// DiscountDocument represents a discount in MongoDB.
type DiscountDocument struct {
	Type        string         `bson:"type"` // percentage, fixed
	Value       int            `bson:"value"`
	MinPurchase *MoneyDocument `bson:"minPurchase,omitempty"`
	MaxDiscount *MoneyDocument `bson:"maxDiscount,omitempty"`
}

// AddressDocument represents an address in MongoDB.
type AddressDocument struct {
	Street     string `bson:"street"`
	City       string `bson:"city"`
	PostalCode string `bson:"postalCode"`
	Country    string `bson:"country"`
	Formatted  string `bson:"formatted"`
}

// ScheduleSlotDocument represents a time slot in MongoDB.
type ScheduleSlotDocument struct {
	DayOfWeek int    `bson:"dayOfWeek"` // 0 = Sunday
	Open      string `bson:"open"`      // HH:MM
	Close     string `bson:"close"`     // HH:MM
	IsClosed  bool   `bson:"isClosed"`
}
