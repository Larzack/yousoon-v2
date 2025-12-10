// Package commands contains command handlers for category management.
package commands

import (
	"context"
	"errors"

	"github.com/yousoon/services/discovery/internal/domain"
)

// =============================================================================
// Create Category Command
// =============================================================================

// CreateCategoryCommand represents a command to create a new category.
type CreateCategoryCommand struct {
	Slug          string
	NameFR        string
	NameEN        string
	DescriptionFR string
	DescriptionEN string
	Icon          string
	Color         string
	Image         string
	ParentID      *string
	Order         int
}

// CreateCategoryHandler handles the create category command.
type CreateCategoryHandler struct {
	categoryRepo domain.CategoryRepository
}

// NewCreateCategoryHandler creates a new CreateCategoryHandler.
func NewCreateCategoryHandler(categoryRepo domain.CategoryRepository) *CreateCategoryHandler {
	return &CreateCategoryHandler{
		categoryRepo: categoryRepo,
	}
}

// Handle executes the create category command.
func (h *CreateCategoryHandler) Handle(ctx context.Context, cmd CreateCategoryCommand) (*domain.Category, error) {
	// Check if slug already exists
	exists, err := h.categoryRepo.ExistsBySlug(ctx, cmd.Slug)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, domain.ErrCategorySlugExists
	}

	// Validate parent if provided
	if cmd.ParentID != nil {
		parent, err := h.categoryRepo.FindByID(ctx, domain.CategoryID(*cmd.ParentID))
		if err != nil {
			return nil, err
		}
		if parent == nil {
			return nil, errors.New("parent category not found")
		}
	}

	// Create category
	category, err := domain.NewCategory(cmd.Slug, cmd.NameFR, cmd.NameEN, cmd.Icon)
	if err != nil {
		return nil, err
	}

	// Set optional fields
	if cmd.DescriptionFR != "" || cmd.DescriptionEN != "" {
		category.UpdateDescription(cmd.DescriptionFR, cmd.DescriptionEN)
	}
	if cmd.Color != "" {
		category.UpdateColor(cmd.Color)
	}
	if cmd.Image != "" {
		category.UpdateImage(cmd.Image)
	}
	if cmd.ParentID != nil {
		parentID := domain.CategoryID(*cmd.ParentID)
		category.SetParent(&parentID)
	}
	if cmd.Order > 0 {
		category.SetOrder(cmd.Order)
	}

	// Save category
	if err := h.categoryRepo.Save(ctx, category); err != nil {
		return nil, err
	}

	return category, nil
}

// =============================================================================
// Update Category Command
// =============================================================================

// UpdateCategoryCommand represents a command to update a category.
type UpdateCategoryCommand struct {
	CategoryID    string
	Slug          *string
	NameFR        *string
	NameEN        *string
	DescriptionFR *string
	DescriptionEN *string
	Icon          *string
	Color         *string
	Image         *string
	ParentID      *string
	Order         *int
}

// UpdateCategoryHandler handles the update category command.
type UpdateCategoryHandler struct {
	categoryRepo domain.CategoryRepository
}

// NewUpdateCategoryHandler creates a new UpdateCategoryHandler.
func NewUpdateCategoryHandler(categoryRepo domain.CategoryRepository) *UpdateCategoryHandler {
	return &UpdateCategoryHandler{
		categoryRepo: categoryRepo,
	}
}

// Handle executes the update category command.
func (h *UpdateCategoryHandler) Handle(ctx context.Context, cmd UpdateCategoryCommand) (*domain.Category, error) {
	category, err := h.categoryRepo.FindByID(ctx, domain.CategoryID(cmd.CategoryID))
	if err != nil {
		return nil, err
	}
	if category == nil {
		return nil, domain.ErrCategoryNotFound
	}

	// Update slug
	if cmd.Slug != nil && *cmd.Slug != category.Slug() {
		exists, err := h.categoryRepo.ExistsBySlug(ctx, *cmd.Slug)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, domain.ErrCategorySlugExists
		}
		if err := category.UpdateSlug(*cmd.Slug); err != nil {
			return nil, err
		}
	}

	// Update name
	if cmd.NameFR != nil || cmd.NameEN != nil {
		nameFR := category.Name().FR
		nameEN := category.Name().EN
		if cmd.NameFR != nil {
			nameFR = *cmd.NameFR
		}
		if cmd.NameEN != nil {
			nameEN = *cmd.NameEN
		}
		if err := category.UpdateName(nameFR, nameEN); err != nil {
			return nil, err
		}
	}

	// Update description
	if cmd.DescriptionFR != nil || cmd.DescriptionEN != nil {
		descFR := category.Description().FR
		descEN := category.Description().EN
		if cmd.DescriptionFR != nil {
			descFR = *cmd.DescriptionFR
		}
		if cmd.DescriptionEN != nil {
			descEN = *cmd.DescriptionEN
		}
		category.UpdateDescription(descFR, descEN)
	}

	// Update icon
	if cmd.Icon != nil {
		category.UpdateIcon(*cmd.Icon)
	}

	// Update color
	if cmd.Color != nil {
		category.UpdateColor(*cmd.Color)
	}

	// Update image
	if cmd.Image != nil {
		category.UpdateImage(*cmd.Image)
	}

	// Update parent
	if cmd.ParentID != nil {
		if *cmd.ParentID == "" {
			category.SetParent(nil)
		} else {
			parentID := domain.CategoryID(*cmd.ParentID)
			category.SetParent(&parentID)
		}
	}

	// Update order
	if cmd.Order != nil {
		category.SetOrder(*cmd.Order)
	}

	// Save changes
	if err := h.categoryRepo.Save(ctx, category); err != nil {
		return nil, err
	}

	return category, nil
}

// =============================================================================
// Delete Category Command
// =============================================================================

// DeleteCategoryCommand represents a command to delete a category.
type DeleteCategoryCommand struct {
	CategoryID string
}

// DeleteCategoryHandler handles the delete category command.
type DeleteCategoryHandler struct {
	categoryRepo domain.CategoryRepository
	offerRepo    domain.OfferRepository
}

// NewDeleteCategoryHandler creates a new DeleteCategoryHandler.
func NewDeleteCategoryHandler(categoryRepo domain.CategoryRepository, offerRepo domain.OfferRepository) *DeleteCategoryHandler {
	return &DeleteCategoryHandler{
		categoryRepo: categoryRepo,
		offerRepo:    offerRepo,
	}
}

// Handle executes the delete category command.
func (h *DeleteCategoryHandler) Handle(ctx context.Context, cmd DeleteCategoryCommand) error {
	categoryID := domain.CategoryID(cmd.CategoryID)

	category, err := h.categoryRepo.FindByID(ctx, categoryID)
	if err != nil {
		return err
	}
	if category == nil {
		return domain.ErrCategoryNotFound
	}

	// Check for child categories
	children, err := h.categoryRepo.FindByParentID(ctx, &categoryID)
	if err != nil {
		return err
	}
	if len(children) > 0 {
		return domain.ErrCategoryHasChildren
	}

	// Check for offers in this category
	count, err := h.offerRepo.Count(ctx, domain.OfferFilter{CategoryID: &categoryID})
	if err != nil {
		return err
	}
	if count > 0 {
		return domain.ErrCategoryHasOffers
	}

	return h.categoryRepo.Delete(ctx, categoryID)
}

// =============================================================================
// Activate/Deactivate Category Commands
// =============================================================================

// ToggleCategoryCommand activates or deactivates a category.
type ToggleCategoryCommand struct {
	CategoryID string
	Activate   bool
}

// ToggleCategoryHandler handles the toggle category command.
type ToggleCategoryHandler struct {
	categoryRepo domain.CategoryRepository
}

// NewToggleCategoryHandler creates a new ToggleCategoryHandler.
func NewToggleCategoryHandler(categoryRepo domain.CategoryRepository) *ToggleCategoryHandler {
	return &ToggleCategoryHandler{
		categoryRepo: categoryRepo,
	}
}

// Handle executes the toggle category command.
func (h *ToggleCategoryHandler) Handle(ctx context.Context, cmd ToggleCategoryCommand) (*domain.Category, error) {
	category, err := h.categoryRepo.FindByID(ctx, domain.CategoryID(cmd.CategoryID))
	if err != nil {
		return nil, err
	}
	if category == nil {
		return nil, domain.ErrCategoryNotFound
	}

	if cmd.Activate {
		category.Activate()
	} else {
		category.Deactivate()
	}

	if err := h.categoryRepo.Save(ctx, category); err != nil {
		return nil, err
	}

	return category, nil
}

// =============================================================================
// Reorder Categories Command
// =============================================================================

// ReorderCategoriesCommand reorders categories.
type ReorderCategoriesCommand struct {
	CategoryOrders []CategoryOrder
}

// CategoryOrder represents the order of a category.
type CategoryOrder struct {
	CategoryID string
	Order      int
}

// ReorderCategoriesHandler handles the reorder categories command.
type ReorderCategoriesHandler struct {
	categoryRepo domain.CategoryRepository
}

// NewReorderCategoriesHandler creates a new ReorderCategoriesHandler.
func NewReorderCategoriesHandler(categoryRepo domain.CategoryRepository) *ReorderCategoriesHandler {
	return &ReorderCategoriesHandler{
		categoryRepo: categoryRepo,
	}
}

// Handle executes the reorder categories command.
func (h *ReorderCategoriesHandler) Handle(ctx context.Context, cmd ReorderCategoriesCommand) error {
	for _, order := range cmd.CategoryOrders {
		category, err := h.categoryRepo.FindByID(ctx, domain.CategoryID(order.CategoryID))
		if err != nil {
			return err
		}
		if category == nil {
			continue // Skip non-existent categories
		}

		category.SetOrder(order.Order)

		if err := h.categoryRepo.Save(ctx, category); err != nil {
			return err
		}
	}

	return nil
}
