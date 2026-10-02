package services

import (
	"strings"

	"layanesia-server/internal/apperrors"
	"layanesia-server/internal/models"
	"layanesia-server/internal/repositories"
	"layanesia-server/internal/utils"

	"github.com/google/uuid"
)

type CreateRatingRequest struct {
	ConnectionID uuid.UUID `json:"connection_id" validate:"required"`
	RatingStars  int       `json:"rating_stars" validate:"required,min=1,max=5"`
	Comment      *string   `json:"comment" validate:"omitempty,max=2000"`
}

type RatingSummary struct {
	Ratings      []models.RatingView `json:"ratings"`
	AverageStars float64             `json:"average_rating"`
	Total        int                 `json:"total_ratings"`
}

type RatingService interface {
	CreateRating(reviewerID uuid.UUID, req CreateRatingRequest) (*models.Rating, error)
	GetUserRatings(userID uuid.UUID, limit, offset int) (*RatingSummary, error)
}

type ratingService struct {
	ratingRepo repositories.RatingRepository
	connRepo   repositories.ConnectionRepository
}

func NewRatingService(ratingRepo repositories.RatingRepository, connRepo repositories.ConnectionRepository) RatingService {
	return &ratingService{ratingRepo: ratingRepo, connRepo: connRepo}
}

// CreateRating hanya untuk peserta koneksi yang sudah COMPLETED; satu ulasan per
// reviewer per koneksi.
func (s *ratingService) CreateRating(reviewerID uuid.UUID, req CreateRatingRequest) (*models.Rating, error) {
	if req.Comment != nil {
		c := strings.TrimSpace(*req.Comment)
		if c == "" {
			req.Comment = nil
		} else {
			req.Comment = &c
		}
	}
	if err := utils.ValidateStruct(req); err != nil {
		return nil, err
	}

	conn, err := s.connRepo.FindByID(req.ConnectionID)
	if err != nil {
		return nil, apperrors.Internal(err)
	}
	if conn == nil || !conn.Involves(reviewerID) {
		return nil, apperrors.NotFound("Koneksi tidak ditemukan")
	}
	if conn.Status != models.ConnStatusCompleted {
		return nil, apperrors.Conflict("CONNECTION_NOT_COMPLETED", "Ulasan hanya dapat diberikan setelah pekerjaan ditandai selesai.")
	}

	exists, err := s.ratingRepo.Exists(req.ConnectionID, reviewerID)
	if err != nil {
		return nil, apperrors.Internal(err)
	}
	if exists {
		return nil, apperrors.Conflict("RATING_EXISTS", "Anda sudah memberikan ulasan untuk koneksi ini.")
	}

	revieweeID, role := conn.WorkerID, models.RoleEmployer
	if conn.WorkerID == reviewerID {
		revieweeID, role = conn.EmployerID, models.RoleWorker
	}

	rating := &models.Rating{
		ConnectionID: req.ConnectionID,
		ReviewerID:   reviewerID,
		RevieweeID:   revieweeID,
		ReviewerRole: role,
		RatingStars:  req.RatingStars,
		Comment:      req.Comment,
	}
	if err := s.ratingRepo.Create(rating); err != nil {
		return nil, apperrors.Conflict("RATING_EXISTS", "Anda sudah memberikan ulasan untuk koneksi ini.").WithCause(err)
	}
	return rating, nil
}

func (s *ratingService) GetUserRatings(userID uuid.UUID, limit, offset int) (*RatingSummary, error) {
	ratings, err := s.ratingRepo.FindByRevieweeID(userID, limit, offset)
	if err != nil {
		return nil, apperrors.Internal(err)
	}
	avg, count, err := s.ratingRepo.GetAverageRating(userID)
	if err != nil {
		return nil, apperrors.Internal(err)
	}
	return &RatingSummary{
		Ratings:      models.RatingViews(ratings),
		AverageStars: avg,
		Total:        count,
	}, nil
}
