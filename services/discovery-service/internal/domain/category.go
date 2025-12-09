// Package domain contains the Category aggregate root.
package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// Category represents a category of offers.
type Category struct {
	id          CategoryID
	slug        string
	name        LocalizedString
	description LocalizedString
	icon        string
	color       string
	image       string
	parentID    *CategoryID
	order       int
	isActive    bool
	createdAt   time.Time
	updatedAt   time.Time

	// Domain events
	events []interface{}
}

// LocalizedString represents a string with translations.
type LocalizedString struct {
	FR string `json:"fr" bson:"fr"`
	EN string `json:"en" bson:"en"`
}

// Get returns the value for the given language code.
func (ls LocalizedString) Get(lang string) string {
	switch lang {
	case "en":
		return ls.EN
	case "fr":
		fallthrough
	default:
		return ls.FR
	}
}

// NewCategory creates a new Category.
func NewCategory(
	slug string,
	nameFR, nameEN string,
	icon string,
) (*Category, error) {
	if slug == "" {
		return nil, errors.New("slug is required")
	}
	if nameFR == "" {
		return nil, errors.New("name (FR) is required")
	}

	now := time.Now()
	return &Category{
		id:   CategoryID(uuid.New().String()),
		slug: slug,
		name: LocalizedString{
			FR: nameFR,
			EN: nameEN,
		},
		icon:      icon,
		isActive:  true,
		order:     0,
		createdAt: now,
		updatedAt: now,
		events:    make([]interface{}, 0),
	}, nil
}

// Getters

func (c *Category) ID() CategoryID               { return c.id }
func (c *Category) Slug() string                 { return c.slug }
func (c *Category) Name() LocalizedString        { return c.name }
func (c *Category) Description() LocalizedString { return c.description }
func (c *Category) Icon() string                 { return c.icon }
func (c *Category) Color() string                { return c.color }
func (c *Category) Image() string                { return c.image }
func (c *Category) ParentID() *CategoryID        { return c.parentID }
func (c *Category) Order() int                   { return c.order }
func (c *Category) IsActive() bool               { return c.isActive }
func (c *Category) CreatedAt() time.Time         { return c.createdAt }
func (c *Category) UpdatedAt() time.Time         { return c.updatedAt }
func (c *Category) Events() []interface{}        { return c.events }
func (c *Category) ClearEvents()                 { c.events = make([]interface{}, 0) }

// IsRoot checks if the category is a root category (no parent).
func (c *Category) IsRoot() bool {
	return c.parentID == nil
}

// =============================================================================
// Commands
// =============================================================================

// UpdateName updates the category name.
func (c *Category) UpdateName(fr, en string) error {
	if fr == "" {
		return errors.New("name (FR) is required")
	}
	c.name = LocalizedString{FR: fr, EN: en}
	c.updatedAt = time.Now()
	return nil
}

// UpdateDescription updates the category description.
func (c *Category) UpdateDescription(fr, en string) {
	c.description = LocalizedString{FR: fr, EN: en}
	c.updatedAt = time.Now()
}

// UpdateSlug updates the category slug.
func (c *Category) UpdateSlug(slug string) error {
	if slug == "" {
		return errors.New("slug is required")
	}
	c.slug = slug
	c.updatedAt = time.Now()
	return nil
}

// UpdateIcon updates the category icon.
func (c *Category) UpdateIcon(icon string) {
	c.icon = icon
	c.updatedAt = time.Now()
}

// UpdateColor updates the category color.
func (c *Category) UpdateColor(color string) {
	c.color = color
	c.updatedAt = time.Now()
}

// UpdateImage updates the category image.
func (c *Category) UpdateImage(image string) {
	c.image = image
	c.updatedAt = time.Now()
}

// SetParent sets the parent category.
func (c *Category) SetParent(parentID *CategoryID) {
	c.parentID = parentID
	c.updatedAt = time.Now()
}

// SetOrder sets the display order.
func (c *Category) SetOrder(order int) {
	c.order = order
	c.updatedAt = time.Now()
}

// Activate activates the category.
func (c *Category) Activate() {
	c.isActive = true
	c.updatedAt = time.Now()
}

// Deactivate deactivates the category.
func (c *Category) Deactivate() {
	c.isActive = false
	c.updatedAt = time.Now()
}

// =============================================================================
// Reconstruction
// =============================================================================

// ReconstructCategory reconstructs a category from persistence.
func ReconstructCategory(
	id CategoryID,
	slug string,
	name LocalizedString,
	description LocalizedString,
	icon string,
	color string,
	image string,
	parentID *CategoryID,
	order int,
	isActive bool,
	createdAt time.Time,
	updatedAt time.Time,
) *Category {
	return &Category{
		id:          id,
		slug:        slug,
		name:        name,
		description: description,
		icon:        icon,
		color:       color,
		image:       image,
		parentID:    parentID,
		order:       order,
		isActive:    isActive,
		createdAt:   createdAt,
		updatedAt:   updatedAt,
		events:      make([]interface{}, 0),
	}
}

// =============================================================================
// CategoryTree for hierarchical display
// =============================================================================

// CategoryTree represents a category with its children.
type CategoryTree struct {
	Category *Category
	Children []*CategoryTree
}

// BuildCategoryTree builds a hierarchical tree from a flat list of categories.
func BuildCategoryTree(categories []*Category) []*CategoryTree {
	// Create a map for quick lookup
	categoryMap := make(map[CategoryID]*CategoryTree)
	for _, cat := range categories {
		categoryMap[cat.ID()] = &CategoryTree{
			Category: cat,
			Children: make([]*CategoryTree, 0),
		}
	}

	// Build tree structure
	roots := make([]*CategoryTree, 0)
	for _, cat := range categories {
		node := categoryMap[cat.ID()]
		if cat.ParentID() == nil {
			roots = append(roots, node)
		} else {
			if parent, ok := categoryMap[*cat.ParentID()]; ok {
				parent.Children = append(parent.Children, node)
			}
		}
	}

	return roots
}
