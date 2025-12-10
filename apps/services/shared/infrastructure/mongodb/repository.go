package mongodb

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Repository provides a generic base for MongoDB repositories.
type Repository[T any] struct {
	collection *mongo.Collection
	mapper     Mapper[T]
}

// Mapper defines the interface for mapping between domain entities and MongoDB documents.
type Mapper[T any] interface {
	// ToDocument converts a domain entity to a MongoDB document.
	ToDocument(entity T) (bson.M, error)
	// FromDocument converts a MongoDB document to a domain entity.
	FromDocument(doc bson.M) (T, error)
}

// NewRepository creates a new generic repository.
func NewRepository[T any](collection *mongo.Collection, mapper Mapper[T]) *Repository[T] {
	return &Repository[T]{
		collection: collection,
		mapper:     mapper,
	}
}

// Collection returns the underlying MongoDB collection.
func (r *Repository[T]) Collection() *mongo.Collection {
	return r.collection
}

// FindByID finds a document by its ObjectID.
func (r *Repository[T]) FindByID(ctx context.Context, id primitive.ObjectID) (T, error) {
	var zero T

	filter := bson.M{"_id": id}
	result := r.collection.FindOne(ctx, filter)
	if result.Err() != nil {
		if result.Err() == mongo.ErrNoDocuments {
			return zero, ErrNotFound
		}
		return zero, fmt.Errorf("failed to find document: %w", result.Err())
	}

	var doc bson.M
	if err := result.Decode(&doc); err != nil {
		return zero, fmt.Errorf("failed to decode document: %w", err)
	}

	entity, err := r.mapper.FromDocument(doc)
	if err != nil {
		return zero, fmt.Errorf("failed to map document: %w", err)
	}

	return entity, nil
}

// FindOne finds a single document matching the filter.
func (r *Repository[T]) FindOne(ctx context.Context, filter bson.M) (T, error) {
	var zero T

	result := r.collection.FindOne(ctx, filter)
	if result.Err() != nil {
		if result.Err() == mongo.ErrNoDocuments {
			return zero, ErrNotFound
		}
		return zero, fmt.Errorf("failed to find document: %w", result.Err())
	}

	var doc bson.M
	if err := result.Decode(&doc); err != nil {
		return zero, fmt.Errorf("failed to decode document: %w", err)
	}

	entity, err := r.mapper.FromDocument(doc)
	if err != nil {
		return zero, fmt.Errorf("failed to map document: %w", err)
	}

	return entity, nil
}

// FindMany finds all documents matching the filter.
func (r *Repository[T]) FindMany(ctx context.Context, filter bson.M, opts ...*options.FindOptions) ([]T, error) {
	cursor, err := r.collection.Find(ctx, filter, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to find documents: %w", err)
	}
	defer cursor.Close(ctx)

	var entities []T
	for cursor.Next(ctx) {
		var doc bson.M
		if err := cursor.Decode(&doc); err != nil {
			return nil, fmt.Errorf("failed to decode document: %w", err)
		}

		entity, err := r.mapper.FromDocument(doc)
		if err != nil {
			return nil, fmt.Errorf("failed to map document: %w", err)
		}

		entities = append(entities, entity)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %w", err)
	}

	return entities, nil
}

// FindWithPagination finds documents with pagination support.
func (r *Repository[T]) FindWithPagination(ctx context.Context, filter bson.M, page, pageSize int64, sort bson.M) (*PaginatedResult[T], error) {
	skip := (page - 1) * pageSize

	opts := options.Find().
		SetSkip(skip).
		SetLimit(pageSize)

	if sort != nil {
		opts.SetSort(sort)
	}

	// Get total count
	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to count documents: %w", err)
	}

	// Get documents
	entities, err := r.FindMany(ctx, filter, opts)
	if err != nil {
		return nil, err
	}

	totalPages := (total + pageSize - 1) / pageSize

	return &PaginatedResult[T]{
		Items:       entities,
		Total:       total,
		Page:        page,
		PageSize:    pageSize,
		TotalPages:  totalPages,
		HasNext:     page < totalPages,
		HasPrevious: page > 1,
	}, nil
}

// PaginatedResult holds paginated query results.
type PaginatedResult[T any] struct {
	Items       []T
	Total       int64
	Page        int64
	PageSize    int64
	TotalPages  int64
	HasNext     bool
	HasPrevious bool
}

// Insert inserts a new document.
func (r *Repository[T]) Insert(ctx context.Context, entity T) (primitive.ObjectID, error) {
	doc, err := r.mapper.ToDocument(entity)
	if err != nil {
		return primitive.NilObjectID, fmt.Errorf("failed to map entity: %w", err)
	}

	result, err := r.collection.InsertOne(ctx, doc)
	if err != nil {
		return primitive.NilObjectID, fmt.Errorf("failed to insert document: %w", err)
	}

	id, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return primitive.NilObjectID, fmt.Errorf("unexpected ID type")
	}

	return id, nil
}

// InsertMany inserts multiple documents.
func (r *Repository[T]) InsertMany(ctx context.Context, entities []T) ([]primitive.ObjectID, error) {
	docs := make([]interface{}, len(entities))
	for i, entity := range entities {
		doc, err := r.mapper.ToDocument(entity)
		if err != nil {
			return nil, fmt.Errorf("failed to map entity %d: %w", i, err)
		}
		docs[i] = doc
	}

	result, err := r.collection.InsertMany(ctx, docs)
	if err != nil {
		return nil, fmt.Errorf("failed to insert documents: %w", err)
	}

	ids := make([]primitive.ObjectID, len(result.InsertedIDs))
	for i, insertedID := range result.InsertedIDs {
		id, ok := insertedID.(primitive.ObjectID)
		if !ok {
			return nil, fmt.Errorf("unexpected ID type at index %d", i)
		}
		ids[i] = id
	}

	return ids, nil
}

// Update updates a document by ID.
func (r *Repository[T]) Update(ctx context.Context, id primitive.ObjectID, entity T) error {
	doc, err := r.mapper.ToDocument(entity)
	if err != nil {
		return fmt.Errorf("failed to map entity: %w", err)
	}

	// Remove _id from update to avoid immutable field error
	delete(doc, "_id")

	filter := bson.M{"_id": id}
	update := bson.M{"$set": doc}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update document: %w", err)
	}

	if result.MatchedCount == 0 {
		return ErrNotFound
	}

	return nil
}

// UpdatePartial performs a partial update using the provided update document.
func (r *Repository[T]) UpdatePartial(ctx context.Context, id primitive.ObjectID, update bson.M) error {
	filter := bson.M{"_id": id}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update document: %w", err)
	}

	if result.MatchedCount == 0 {
		return ErrNotFound
	}

	return nil
}

// UpdateMany updates multiple documents matching the filter.
func (r *Repository[T]) UpdateMany(ctx context.Context, filter bson.M, update bson.M) (int64, error) {
	result, err := r.collection.UpdateMany(ctx, filter, update)
	if err != nil {
		return 0, fmt.Errorf("failed to update documents: %w", err)
	}

	return result.ModifiedCount, nil
}

// Delete deletes a document by ID.
func (r *Repository[T]) Delete(ctx context.Context, id primitive.ObjectID) error {
	filter := bson.M{"_id": id}

	result, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("failed to delete document: %w", err)
	}

	if result.DeletedCount == 0 {
		return ErrNotFound
	}

	return nil
}

// SoftDelete performs a soft delete by setting deletedAt timestamp.
func (r *Repository[T]) SoftDelete(ctx context.Context, id primitive.ObjectID) error {
	now := primitive.NewDateTimeFromTime(time.Now())
	return r.UpdatePartial(ctx, id, bson.M{
		"$set": bson.M{
			"updatedAt": now,
			"deletedAt": now,
		},
	})
}

// DeleteMany deletes multiple documents matching the filter.
func (r *Repository[T]) DeleteMany(ctx context.Context, filter bson.M) (int64, error) {
	result, err := r.collection.DeleteMany(ctx, filter)
	if err != nil {
		return 0, fmt.Errorf("failed to delete documents: %w", err)
	}

	return result.DeletedCount, nil
}

// Count counts documents matching the filter.
func (r *Repository[T]) Count(ctx context.Context, filter bson.M) (int64, error) {
	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, fmt.Errorf("failed to count documents: %w", err)
	}
	return count, nil
}

// Exists checks if a document matching the filter exists.
func (r *Repository[T]) Exists(ctx context.Context, filter bson.M) (bool, error) {
	count, err := r.Count(ctx, filter)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// ExistsByID checks if a document with the given ID exists.
func (r *Repository[T]) ExistsByID(ctx context.Context, id primitive.ObjectID) (bool, error) {
	return r.Exists(ctx, bson.M{"_id": id})
}

// Aggregate performs an aggregation pipeline.
func (r *Repository[T]) Aggregate(ctx context.Context, pipeline mongo.Pipeline, opts ...*options.AggregateOptions) (*mongo.Cursor, error) {
	cursor, err := r.collection.Aggregate(ctx, pipeline, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to aggregate: %w", err)
	}
	return cursor, nil
}

// BulkWrite performs bulk write operations.
func (r *Repository[T]) BulkWrite(ctx context.Context, models []mongo.WriteModel, opts ...*options.BulkWriteOptions) (*mongo.BulkWriteResult, error) {
	result, err := r.collection.BulkWrite(ctx, models, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to bulk write: %w", err)
	}
	return result, nil
}

// CreateIndex creates an index on the collection.
func (r *Repository[T]) CreateIndex(ctx context.Context, keys bson.D, opts ...*options.IndexOptions) (string, error) {
	indexModel := mongo.IndexModel{
		Keys:    keys,
		Options: nil,
	}
	if len(opts) > 0 {
		indexModel.Options = opts[0]
	}

	name, err := r.collection.Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		return "", fmt.Errorf("failed to create index: %w", err)
	}
	return name, nil
}

// ErrNotFound is returned when a document is not found.
var ErrNotFound = fmt.Errorf("document not found")
