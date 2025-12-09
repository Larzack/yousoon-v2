// Package mongodb implements the persistence layer for categories.
package mongodb

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/yousoon/discovery-service/internal/domain"
)

const categoryCollection = "categories"

// CategoryRepository implements domain.CategoryRepository using MongoDB.
type CategoryRepository struct {
	collection *mongo.Collection
}

// NewCategoryRepository creates a new MongoDB category repository.
func NewCategoryRepository(db *mongo.Database) *CategoryRepository {
	return &CategoryRepository{
		collection: db.Collection(categoryCollection),
	}
}

// categoryDocument represents a category in MongoDB.
type categoryDocument struct {
	ID          string            `bson:"_id"`
	Name        map[string]string `bson:"name"`
	Slug        string            `bson:"slug"`
	Description map[string]string `bson:"description,omitempty"`
	Icon        string            `bson:"icon,omitempty"`
	Color       string            `bson:"color,omitempty"`
	Image       string            `bson:"image,omitempty"`
	ParentID    *string           `bson:"parent_id,omitempty"`
	Position    int               `bson:"position"`
	IsActive    bool              `bson:"is_active"`
	CreatedAt   time.Time         `bson:"created_at"`
	UpdatedAt   time.Time         `bson:"updated_at"`
}

// EnsureIndexes creates necessary indexes for the category collection.
func (r *CategoryRepository) EnsureIndexes(ctx context.Context) error {
	indexes := []mongo.IndexModel{
		{
			Keys:    bson.D{bson.E{Key: "slug", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{bson.E{Key: "parent_id", Value: 1}},
		},
		{
			Keys: bson.D{
				bson.E{Key: "is_active", Value: 1},
				bson.E{Key: "position", Value: 1},
			},
		},
		{
			Keys: bson.D{bson.E{Key: "position", Value: 1}},
		},
	}

	_, err := r.collection.Indexes().CreateMany(ctx, indexes)
	if err != nil {
		return fmt.Errorf("failed to create indexes: %w", err)
	}

	return nil
}

// Save persists a category to MongoDB.
func (r *CategoryRepository) Save(ctx context.Context, category *domain.Category) error {
	doc := r.toDocument(category)

	filter := bson.M{"_id": doc.ID}
	update := bson.M{"$set": doc}
	opts := options.Update().SetUpsert(true)

	_, err := r.collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return fmt.Errorf("failed to save category: %w", err)
	}

	return nil
}

// FindByID retrieves a category by its ID.
func (r *CategoryRepository) FindByID(ctx context.Context, id domain.CategoryID) (*domain.Category, error) {
	filter := bson.M{"_id": id.String()}

	var doc categoryDocument
	err := r.collection.FindOne(ctx, filter).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, domain.ErrCategoryNotFound
		}
		return nil, fmt.Errorf("failed to find category: %w", err)
	}

	return r.toDomain(&doc)
}

// FindBySlug retrieves a category by its slug.
func (r *CategoryRepository) FindBySlug(ctx context.Context, slug string) (*domain.Category, error) {
	filter := bson.M{"slug": slug}

	var doc categoryDocument
	err := r.collection.FindOne(ctx, filter).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, domain.ErrCategoryNotFound
		}
		return nil, fmt.Errorf("failed to find category by slug: %w", err)
	}

	return r.toDomain(&doc)
}

// FindAll retrieves all categories, optionally filtered by active status.
func (r *CategoryRepository) FindAll(ctx context.Context, activeOnly bool) ([]*domain.Category, error) {
	filter := bson.M{}
	if activeOnly {
		filter["is_active"] = true
	}

	opts := options.Find().SetSort(bson.D{bson.E{Key: "position", Value: 1}})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to find categories: %w", err)
	}
	defer cursor.Close(ctx)

	var docs []categoryDocument
	if err := cursor.All(ctx, &docs); err != nil {
		return nil, fmt.Errorf("failed to decode categories: %w", err)
	}

	categories := make([]*domain.Category, 0, len(docs))
	for _, doc := range docs {
		category, err := r.toDomain(&doc)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}

	return categories, nil
}

// FindByParentID retrieves all child categories of a parent.
func (r *CategoryRepository) FindByParentID(ctx context.Context, parentID *domain.CategoryID, activeOnly bool) ([]*domain.Category, error) {
	filter := bson.M{}

	if parentID != nil {
		filter["parent_id"] = parentID.String()
	} else {
		filter["parent_id"] = nil
	}

	if activeOnly {
		filter["is_active"] = true
	}

	opts := options.Find().SetSort(bson.D{bson.E{Key: "position", Value: 1}})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to find categories: %w", err)
	}
	defer cursor.Close(ctx)

	var docs []categoryDocument
	if err := cursor.All(ctx, &docs); err != nil {
		return nil, fmt.Errorf("failed to decode categories: %w", err)
	}

	categories := make([]*domain.Category, 0, len(docs))
	for _, doc := range docs {
		category, err := r.toDomain(&doc)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}

	return categories, nil
}

// FindRootCategories retrieves all root categories (no parent).
func (r *CategoryRepository) FindRootCategories(ctx context.Context, activeOnly bool) ([]*domain.Category, error) {
	return r.FindByParentID(ctx, nil, activeOnly)
}

// Delete removes a category from the database.
func (r *CategoryRepository) Delete(ctx context.Context, id domain.CategoryID) error {
	filter := bson.M{"_id": id.String()}

	result, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("failed to delete category: %w", err)
	}

	if result.DeletedCount == 0 {
		return domain.ErrCategoryNotFound
	}

	return nil
}

// HasChildren checks if a category has child categories.
func (r *CategoryRepository) HasChildren(ctx context.Context, id domain.CategoryID) (bool, error) {
	filter := bson.M{"parent_id": id.String()}

	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return false, fmt.Errorf("failed to count children: %w", err)
	}

	return count > 0, nil
}

// CountOffers counts the number of offers in a category.
func (r *CategoryRepository) CountOffers(ctx context.Context, id domain.CategoryID) (int64, error) {
	// This would typically query the offers collection
	// For now, return 0 as this is cross-aggregate
	return 0, nil
}

// GetCategoryTree retrieves the full category tree.
func (r *CategoryRepository) GetCategoryTree(ctx context.Context, activeOnly bool) ([]*domain.CategoryTree, error) {
	categories, err := r.FindAll(ctx, activeOnly)
	if err != nil {
		return nil, err
	}

	// Use the domain function to build the tree
	return domain.BuildCategoryTree(categories), nil
}

func ptrCategoryID(id domain.CategoryID) *domain.CategoryID {
	return &id
}

// UpdatePositions updates the positions of multiple categories.
func (r *CategoryRepository) UpdatePositions(ctx context.Context, positions map[domain.CategoryID]int) error {
	for id, position := range positions {
		filter := bson.M{"_id": id.String()}
		update := bson.M{
			"$set": bson.M{
				"position":   position,
				"updated_at": time.Now(),
			},
		}

		_, err := r.collection.UpdateOne(ctx, filter, update)
		if err != nil {
			return fmt.Errorf("failed to update position for category %s: %w", id, err)
		}
	}

	return nil
}

// ExistsBySlug checks if a category with the given slug exists.
func (r *CategoryRepository) ExistsBySlug(ctx context.Context, slug string) (bool, error) {
	filter := bson.M{"slug": slug}

	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return false, fmt.Errorf("failed to check slug existence: %w", err)
	}

	return count > 0, nil
}

// GetMaxPosition returns the maximum position value.
func (r *CategoryRepository) GetMaxPosition(ctx context.Context, parentID *domain.CategoryID) (int, error) {
	filter := bson.M{}
	if parentID != nil {
		filter["parent_id"] = parentID.String()
	} else {
		filter["parent_id"] = nil
	}

	opts := options.FindOne().SetSort(bson.D{bson.E{Key: "position", Value: -1}})

	var doc categoryDocument
	err := r.collection.FindOne(ctx, filter, opts).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return 0, nil
		}
		return 0, fmt.Errorf("failed to get max position: %w", err)
	}

	return doc.Position, nil
}

// toDocument converts a domain category to a MongoDB document.
func (r *CategoryRepository) toDocument(category *domain.Category) *categoryDocument {
	doc := &categoryDocument{
		ID:          category.ID().String(),
		Name:        map[string]string{"fr": category.Name().FR, "en": category.Name().EN},
		Slug:        category.Slug(),
		Description: map[string]string{"fr": category.Description().FR, "en": category.Description().EN},
		Icon:        category.Icon(),
		Color:       category.Color(),
		Image:       category.Image(),
		Position:    category.Order(),
		IsActive:    category.IsActive(),
		CreatedAt:   category.CreatedAt(),
		UpdatedAt:   category.UpdatedAt(),
	}

	if category.ParentID() != nil {
		parentIDStr := category.ParentID().String()
		doc.ParentID = &parentIDStr
	}

	return doc
}

// toDomain converts a MongoDB document to a domain category.
func (r *CategoryRepository) toDomain(doc *categoryDocument) (*domain.Category, error) {
	id, err := domain.ParseCategoryID(doc.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to parse category ID: %w", err)
	}

	var parentID *domain.CategoryID
	if doc.ParentID != nil {
		pid, err := domain.ParseCategoryID(*doc.ParentID)
		if err != nil {
			return nil, fmt.Errorf("failed to parse parent ID: %w", err)
		}
		parentID = &pid
	}

	// Build LocalizedString from map
	name := domain.LocalizedString{
		FR: doc.Name["fr"],
		EN: doc.Name["en"],
	}
	description := domain.LocalizedString{
		FR: doc.Description["fr"],
		EN: doc.Description["en"],
	}

	return domain.ReconstructCategory(
		id,
		doc.Slug,
		name,
		description,
		doc.Icon,
		doc.Color,
		doc.Image,
		parentID,
		doc.Position,
		doc.IsActive,
		doc.CreatedAt,
		doc.UpdatedAt,
	), nil
}
