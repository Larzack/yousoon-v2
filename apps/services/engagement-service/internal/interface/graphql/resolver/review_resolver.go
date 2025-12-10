package resolver

import (
	"context"
	"time"

	"github.com/yousoon/apps/services/engagement-service/internal/domain"
	"github.com/yousoon/apps/services/engagement-service/internal/interface/graphql/model"
)

// CreateReview crée un nouvel avis
func (r *Resolver) CreateReview(ctx context.Context, input model.CreateReviewInput) (*model.Review, error) {
	userID := ctx.Value("userID").(string)

	// Vérifier si l'utilisateur a déjà laissé un avis
	exists, err := r.reviewRepo.Exists(ctx, userID, input.OfferID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, domain.ErrReviewAlreadyExists
	}

	now := time.Now()
	review := &domain.Review{
		UserID:  userID,
		OfferID: input.OfferID,
		Rating:  input.Rating,
		Title:   derefString(input.Title),
		Content: input.Content,
		Images:  input.Images,
		Moderation: domain.Moderation{
			Status: domain.ModerationStatusPending,
		},
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := r.reviewRepo.Create(ctx, review); err != nil {
		return nil, err
	}

	return r.toModelReview(review), nil
}

// UpdateReview met à jour un avis
func (r *Resolver) UpdateReview(ctx context.Context, id string, input model.UpdateReviewInput) (*model.Review, error) {
	userID := ctx.Value("userID").(string)

	review, err := r.reviewRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if review.UserID != userID {
		return nil, domain.ErrUnauthorized
	}

	if input.Rating != nil {
		review.Rating = *input.Rating
	}
	if input.Title != nil {
		review.Title = *input.Title
	}
	if input.Content != nil {
		review.Content = *input.Content
	}
	if input.Images != nil {
		review.Images = input.Images
	}

	review.UpdatedAt = time.Now()

	if err := r.reviewRepo.Update(ctx, review); err != nil {
		return nil, err
	}

	return r.toModelReview(review), nil
}

// DeleteReview supprime un avis
func (r *Resolver) DeleteReview(ctx context.Context, id string) (bool, error) {
	userID := ctx.Value("userID").(string)

	review, err := r.reviewRepo.FindByID(ctx, id)
	if err != nil {
		return false, err
	}

	if review.UserID != userID {
		return false, domain.ErrUnauthorized
	}

	if err := r.reviewRepo.Delete(ctx, id); err != nil {
		return false, err
	}

	return true, nil
}

// ReportReview signale un avis
func (r *Resolver) ReportReview(ctx context.Context, id, reason string) (bool, error) {
	userID := ctx.Value("userID").(string)

	if err := r.reviewRepo.AddReport(ctx, id, userID, reason); err != nil {
		return false, err
	}

	return true, nil
}

// OfferReviews retourne les avis d'une offre
func (r *Resolver) OfferReviews(ctx context.Context, offerID string, first *int, after *string) (*model.ReviewsConnection, error) {
	limit := 20
	if first != nil {
		limit = *first
	}

	offset := 0

	reviews, total, err := r.reviewRepo.FindByOfferID(ctx, offerID, limit, offset)
	if err != nil {
		return nil, err
	}

	avg, _, _ := r.reviewRepo.GetAverageRating(ctx, offerID)

	edges := make([]*model.ReviewEdge, len(reviews))
	for i, rev := range reviews {
		edges[i] = &model.ReviewEdge{
			Node:   r.toModelReview(rev),
			Cursor: rev.ID,
		}
	}

	return &model.ReviewsConnection{
		Edges: edges,
		PageInfo: &model.PageInfo{
			HasNextPage:     len(reviews) == limit && offset+limit < total,
			HasPreviousPage: offset > 0,
		},
		TotalCount:    total,
		AverageRating: &avg,
	}, nil
}

// MyReviews retourne les avis de l'utilisateur
func (r *Resolver) MyReviews(ctx context.Context, first *int, after *string) (*model.ReviewsConnection, error) {
	userID := ctx.Value("userID").(string)

	limit := 20
	if first != nil {
		limit = *first
	}

	offset := 0

	reviews, total, err := r.reviewRepo.FindByUserID(ctx, userID, limit, offset)
	if err != nil {
		return nil, err
	}

	edges := make([]*model.ReviewEdge, len(reviews))
	for i, rev := range reviews {
		edges[i] = &model.ReviewEdge{
			Node:   r.toModelReview(rev),
			Cursor: rev.ID,
		}
	}

	return &model.ReviewsConnection{
		Edges: edges,
		PageInfo: &model.PageInfo{
			HasNextPage:     len(reviews) == limit && offset+limit < total,
			HasPreviousPage: offset > 0,
		},
		TotalCount: total,
	}, nil
}

// AverageRating retourne la note moyenne d'une offre
func (r *Resolver) AverageRating(ctx context.Context, offerID string) (*float64, error) {
	avg, _, err := r.reviewRepo.GetAverageRating(ctx, offerID)
	if err != nil {
		return nil, err
	}
	if avg == 0 {
		return nil, nil
	}
	return &avg, nil
}

// ReviewCount retourne le nombre d'avis d'une offre
func (r *Resolver) ReviewCount(ctx context.Context, offerID string) (int, error) {
	_, count, err := r.reviewRepo.GetAverageRating(ctx, offerID)
	return count, err
}

func (r *Resolver) toModelReview(review *domain.Review) *model.Review {
	return &model.Review{
		ID:                 review.ID,
		UserID:             review.UserID,
		OfferID:            review.OfferID,
		Rating:             review.Rating,
		Title:              &review.Title,
		Content:            review.Content,
		Images:             review.Images,
		HelpfulCount:       review.HelpfulCount,
		IsVerifiedPurchase: review.IsVerifiedPurchase,
		CreatedAt:          review.CreatedAt,
		UpdatedAt:          review.UpdatedAt,
	}
}

func derefString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
