package resolver

import (
	"context"
	"errors"

	"github.com/yousoon/apps/services/engagement-service/internal/domain"
	"github.com/yousoon/apps/services/engagement-service/internal/interface/graphql/model"
)

var ErrUnauthorized = errors.New("unauthorized")

// CreateReview crée un nouvel avis
func (r *Resolver) CreateReview(ctx context.Context, input model.CreateReviewInput) (*model.Review, error) {
	userID := ctx.Value("userID").(string)

	// Vérifier si l'utilisateur a déjà laissé un avis
	existing, _ := r.reviewRepo.GetByUserAndOffer(ctx, userID, input.OfferID)
	if existing != nil {
		return nil, domain.ErrReviewAlreadyExists
	}

	// TODO: Récupérer les infos de l'offre et de l'utilisateur via gRPC
	// Pour l'instant, on utilise des valeurs vides
	title := ""
	if input.Title != nil {
		title = *input.Title
	}

	review, err := domain.NewReview(
		userID,
		input.OfferID,
		"",  // partnerID - TODO: fetch via gRPC
		"",  // establishmentID - TODO: fetch via gRPC
		nil, // bookingID - TODO: check via gRPC
		input.Rating,
		title,
		input.Content,
		input.Images,
		"",    // userFirstName - TODO: fetch via gRPC
		"",    // userAvatar - TODO: fetch via gRPC
		"",    // offerTitle - TODO: fetch via gRPC
		"",    // partnerName - TODO: fetch via gRPC
		false, // isVerifiedPurchase - TODO: check via gRPC
	)
	if err != nil {
		return nil, err
	}

	if err := r.reviewRepo.Create(ctx, review); err != nil {
		return nil, err
	}

	return r.toModelReview(review), nil
}

// UpdateReview met à jour un avis
func (r *Resolver) UpdateReview(ctx context.Context, id string, input model.UpdateReviewInput) (*model.Review, error) {
	userID := ctx.Value("userID").(string)

	review, err := r.reviewRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if review.UserID() != userID {
		return nil, ErrUnauthorized
	}

	rating := review.Rating()
	if input.Rating != nil {
		rating = *input.Rating
	}

	title := review.Title()
	if input.Title != nil {
		title = *input.Title
	}

	content := review.Content()
	if input.Content != nil {
		content = *input.Content
	}

	images := review.Images()
	if input.Images != nil {
		images = input.Images
	}

	if err := review.Update(rating, title, content, images); err != nil {
		return nil, err
	}

	if err := r.reviewRepo.Update(ctx, review); err != nil {
		return nil, err
	}

	return r.toModelReview(review), nil
}

// DeleteReview supprime un avis
func (r *Resolver) DeleteReview(ctx context.Context, id string) (bool, error) {
	userID := ctx.Value("userID").(string)

	review, err := r.reviewRepo.GetByID(ctx, id)
	if err != nil {
		return false, err
	}

	if review.UserID() != userID {
		return false, ErrUnauthorized
	}

	if err := r.reviewRepo.Delete(ctx, id); err != nil {
		return false, err
	}

	return true, nil
}

// ReportReview signale un avis
func (r *Resolver) ReportReview(ctx context.Context, id, reason string) (bool, error) {
	userID := ctx.Value("userID").(string)

	review, err := r.reviewRepo.GetByID(ctx, id)
	if err != nil {
		return false, err
	}

	if err := review.Report(userID, reason); err != nil {
		return false, err
	}

	if err := r.reviewRepo.Update(ctx, review); err != nil {
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

	filter := domain.ReviewFilter{
		Offset: offset,
		Limit:  limit,
	}

	reviews, total, err := r.reviewRepo.GetByOfferID(ctx, offerID, filter)
	if err != nil {
		return nil, err
	}

	avg, _, _ := r.reviewRepo.GetAverageRating(ctx, offerID)

	edges := make([]*model.ReviewEdge, len(reviews))
	for i, rev := range reviews {
		edges[i] = &model.ReviewEdge{
			Node:   r.toModelReview(rev),
			Cursor: rev.ID(),
		}
	}

	return &model.ReviewsConnection{
		Edges: edges,
		PageInfo: &model.PageInfo{
			HasNextPage:     len(reviews) == limit && offset+limit < int(total),
			HasPreviousPage: offset > 0,
		},
		TotalCount:    int(total),
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

	filter := domain.ReviewFilter{
		Offset: offset,
		Limit:  limit,
	}

	reviews, total, err := r.reviewRepo.GetByUserID(ctx, userID, filter)
	if err != nil {
		return nil, err
	}

	edges := make([]*model.ReviewEdge, len(reviews))
	for i, rev := range reviews {
		edges[i] = &model.ReviewEdge{
			Node:   r.toModelReview(rev),
			Cursor: rev.ID(),
		}
	}

	return &model.ReviewsConnection{
		Edges: edges,
		PageInfo: &model.PageInfo{
			HasNextPage:     len(reviews) == limit && offset+limit < int(total),
			HasPreviousPage: offset > 0,
		},
		TotalCount: int(total),
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
	return int(count), err
}

func (r *Resolver) toModelReview(review *domain.Review) *model.Review {
	title := review.Title()
	return &model.Review{
		ID:                 review.ID(),
		UserID:             review.UserID(),
		OfferID:            review.OfferID(),
		Rating:             review.Rating(),
		Title:              &title,
		Content:            review.Content(),
		Images:             review.Images(),
		HelpfulCount:       review.HelpfulCount(),
		IsVerifiedPurchase: review.IsVerifiedPurchase(),
		CreatedAt:          review.CreatedAt(),
		UpdatedAt:          review.UpdatedAt(),
	}
}

func derefString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
