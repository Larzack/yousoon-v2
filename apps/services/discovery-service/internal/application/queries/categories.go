// Package queries contains query handlers for categories.
package queries

import (
	"context"

	"github.com/yousoon/services/discovery/internal/domain"
)

// =============================================================================
// Get Category Query
// =============================================================================

// GetCategoryQuery retrieves a single category by ID.
type GetCategoryQuery struct {
	CategoryID string
}

// GetCategoryHandler handles the get category query.
type GetCategoryHandler struct {
	categoryRepo domain.CategoryRepository
}

// NewGetCategoryHandler creates a new GetCategoryHandler.
func NewGetCategoryHandler(categoryRepo domain.CategoryRepository) *GetCategoryHandler {
	return &GetCategoryHandler{
		categoryRepo: categoryRepo,
	}
}

// Handle executes the get category query.
func (h *GetCategoryHandler) Handle(ctx context.Context, query GetCategoryQuery) (*domain.Category, error) {
	category, err := h.categoryRepo.FindByID(ctx, domain.CategoryID(query.CategoryID))
	if err != nil {
		return nil, err
	}
	if category == nil {
		return nil, domain.ErrCategoryNotFound
	}
	return category, nil
}

// =============================================================================
// Get Category By Slug Query
// =============================================================================

// GetCategoryBySlugQuery retrieves a category by slug.
type GetCategoryBySlugQuery struct {
	Slug string
}

// GetCategoryBySlugHandler handles the get category by slug query.
type GetCategoryBySlugHandler struct {
	categoryRepo domain.CategoryRepository
}

// NewGetCategoryBySlugHandler creates a new GetCategoryBySlugHandler.
func NewGetCategoryBySlugHandler(categoryRepo domain.CategoryRepository) *GetCategoryBySlugHandler {
	return &GetCategoryBySlugHandler{
		categoryRepo: categoryRepo,
	}
}

// Handle executes the get category by slug query.
func (h *GetCategoryBySlugHandler) Handle(ctx context.Context, query GetCategoryBySlugQuery) (*domain.Category, error) {
	category, err := h.categoryRepo.FindBySlug(ctx, query.Slug)
	if err != nil {
		return nil, err
	}
	if category == nil {
		return nil, domain.ErrCategoryNotFound
	}
	return category, nil
}

// =============================================================================
// List Categories Query
// =============================================================================

// ListCategoriesQuery retrieves all categories.
type ListCategoriesQuery struct {
	OnlyActive    bool
	OnlyRoot      bool
	ParentID      *string
	IncludeHidden bool
}

// ListCategoriesHandler handles the list categories query.
type ListCategoriesHandler struct {
	categoryRepo domain.CategoryRepository
}

// NewListCategoriesHandler creates a new ListCategoriesHandler.
func NewListCategoriesHandler(categoryRepo domain.CategoryRepository) *ListCategoriesHandler {
	return &ListCategoriesHandler{
		categoryRepo: categoryRepo,
	}
}

// Handle executes the list categories query.
func (h *ListCategoriesHandler) Handle(ctx context.Context, query ListCategoriesQuery) ([]*domain.Category, error) {
	// Get root categories only
	if query.OnlyRoot {
		return h.categoryRepo.FindRootCategories(ctx)
	}

	// Get categories by parent
	if query.ParentID != nil {
		parentID := domain.CategoryID(*query.ParentID)
		return h.categoryRepo.FindByParentID(ctx, &parentID)
	}

	// Get active categories only
	if query.OnlyActive {
		return h.categoryRepo.FindActive(ctx)
	}

	// Get all categories
	return h.categoryRepo.FindAll(ctx)
}

// =============================================================================
// Get Category Tree Query
// =============================================================================

// GetCategoryTreeQuery retrieves the full category tree.
type GetCategoryTreeQuery struct {
	OnlyActive bool
}

// GetCategoryTreeHandler handles the get category tree query.
type GetCategoryTreeHandler struct {
	categoryRepo domain.CategoryRepository
}

// NewGetCategoryTreeHandler creates a new GetCategoryTreeHandler.
func NewGetCategoryTreeHandler(categoryRepo domain.CategoryRepository) *GetCategoryTreeHandler {
	return &GetCategoryTreeHandler{
		categoryRepo: categoryRepo,
	}
}

// Handle executes the get category tree query.
func (h *GetCategoryTreeHandler) Handle(ctx context.Context, query GetCategoryTreeQuery) ([]*domain.CategoryTree, error) {
	return h.categoryRepo.GetCategoryTree(ctx)
}

// =============================================================================
// Get Category Summaries Query (with offer counts)
// =============================================================================

// GetCategorySummariesQuery retrieves category summaries with offer counts.
type GetCategorySummariesQuery struct {
	OnlyActive bool
	Language   string // "fr" or "en"
}

// GetCategorySummariesHandler handles the get category summaries query.
type GetCategorySummariesHandler struct {
	categoryRepo domain.CategoryRepository
}

// NewGetCategorySummariesHandler creates a new GetCategorySummariesHandler.
func NewGetCategorySummariesHandler(categoryRepo domain.CategoryRepository) *GetCategorySummariesHandler {
	return &GetCategorySummariesHandler{
		categoryRepo: categoryRepo,
	}
}

// Handle executes the get category summaries query.
func (h *GetCategorySummariesHandler) Handle(ctx context.Context, query GetCategorySummariesQuery) ([]domain.CategorySummary, error) {
	return h.categoryRepo.GetCategorySummaries(ctx)
}

// =============================================================================
// Get Child Categories Query
// =============================================================================

// GetChildCategoriesQuery retrieves child categories of a parent.
type GetChildCategoriesQuery struct {
	ParentID string
}

// GetChildCategoriesHandler handles the get child categories query.
type GetChildCategoriesHandler struct {
	categoryRepo domain.CategoryRepository
}

// NewGetChildCategoriesHandler creates a new GetChildCategoriesHandler.
func NewGetChildCategoriesHandler(categoryRepo domain.CategoryRepository) *GetChildCategoriesHandler {
	return &GetChildCategoriesHandler{
		categoryRepo: categoryRepo,
	}
}

// Handle executes the get child categories query.
func (h *GetChildCategoriesHandler) Handle(ctx context.Context, query GetChildCategoriesQuery) ([]*domain.Category, error) {
	parentID := domain.CategoryID(query.ParentID)
	return h.categoryRepo.FindByParentID(ctx, &parentID)
}

// =============================================================================
// Check Category Exists Query
// =============================================================================

// CheckCategoryExistsQuery checks if a category exists.
type CheckCategoryExistsQuery struct {
	CategoryID string
}

// CheckCategoryExistsHandler handles the check category exists query.
type CheckCategoryExistsHandler struct {
	categoryRepo domain.CategoryRepository
}

// NewCheckCategoryExistsHandler creates a new CheckCategoryExistsHandler.
func NewCheckCategoryExistsHandler(categoryRepo domain.CategoryRepository) *CheckCategoryExistsHandler {
	return &CheckCategoryExistsHandler{
		categoryRepo: categoryRepo,
	}
}

// Handle executes the check category exists query.
func (h *CheckCategoryExistsHandler) Handle(ctx context.Context, query CheckCategoryExistsQuery) (bool, error) {
	category, err := h.categoryRepo.FindByID(ctx, domain.CategoryID(query.CategoryID))
	if err != nil {
		return false, err
	}
	return category != nil, nil
}
