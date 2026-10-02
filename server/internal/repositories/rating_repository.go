package repositories

import (
	"layanesia-server/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RatingRepository interface {
	Create(rating *models.Rating) error
	Exists(connectionID uuid.UUID, reviewerID uuid.UUID) (bool, error)
	FindByConnectionID(connectionID uuid.UUID) ([]models.Rating, error)
	FindByRevieweeID(revieweeID uuid.UUID, limit, offset int) ([]models.Rating, error)
	GetAverageRating(revieweeID uuid.UUID) (float64, int, error)
}

type ratingRepository struct {
	db *gorm.DB
}

func NewRatingRepository(db *gorm.DB) RatingRepository {
	return &ratingRepository{db: db}
}

func (r *ratingRepository) Create(rating *models.Rating) error {
	return r.db.Create(rating).Error
}

func (r *ratingRepository) Exists(connectionID uuid.UUID, reviewerID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.Model(&models.Rating{}).
		Where("connection_id = ? AND reviewer_id = ?", connectionID, reviewerID).
		Count(&count).Error
	return count > 0, err
}

func (r *ratingRepository) FindByConnectionID(connectionID uuid.UUID) ([]models.Rating, error) {
	var ratings []models.Rating
	err := r.db.Preload("Reviewer").Where("connection_id = ?", connectionID).Find(&ratings).Error
	return ratings, err
}

func (r *ratingRepository) FindByRevieweeID(revieweeID uuid.UUID, limit, offset int) ([]models.Rating, error) {
	var ratings []models.Rating
	err := r.db.Preload("Reviewer").
		Where("reviewee_id = ?", revieweeID).
		Order("created_at DESC").
		Limit(clampLimit(limit, 20, 100)).Offset(maxInt(offset, 0)).
		Find(&ratings).Error
	return ratings, err
}

func (r *ratingRepository) GetAverageRating(revieweeID uuid.UUID) (float64, int, error) {
	var res struct {
		AvgRating float64
		Count     int
	}
	err := r.db.Model(&models.Rating{}).
		Select("COALESCE(AVG(rating_stars), 0) AS avg_rating, COUNT(id) AS count").
		Where("reviewee_id = ?", revieweeID).
		Scan(&res).Error
	if err != nil {
		return 0, 0, err
	}
	return res.AvgRating, res.Count, nil
}
