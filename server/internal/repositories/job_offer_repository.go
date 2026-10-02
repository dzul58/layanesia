package repositories

import (
	"time"

	"layanesia-server/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type JobOfferRepository interface {
	WithTx(tx *gorm.DB) JobOfferRepository
	Create(offer *models.JobOffer) error
	FindByWorkerID(workerID uuid.UUID) ([]models.JobOffer, error)
	FindByEmployerID(employerID uuid.UUID) ([]models.JobOffer, error)
	FindByID(id uuid.UUID) (*models.JobOffer, error)
	FindByIDForUpdate(id uuid.UUID) (*models.JobOffer, error)
	FindPending(skillPostingID uuid.UUID, employerID uuid.UUID) (*models.JobOffer, error)
	UpdateStatus(id uuid.UUID, status models.OfferStatus) error
}

type jobOfferRepository struct {
	db *gorm.DB
}

func NewJobOfferRepository(db *gorm.DB) JobOfferRepository {
	return &jobOfferRepository{db: db}
}

func (r *jobOfferRepository) WithTx(tx *gorm.DB) JobOfferRepository { return &jobOfferRepository{db: tx} }

func (r *jobOfferRepository) Create(offer *models.JobOffer) error {
	return r.db.Create(offer).Error
}

func (r *jobOfferRepository) withRelations() *gorm.DB {
	return r.db.Preload("SkillPosting").Preload("SkillPosting.User").Preload("Employer").Preload("Worker")
}

func (r *jobOfferRepository) FindByWorkerID(workerID uuid.UUID) ([]models.JobOffer, error) {
	var offers []models.JobOffer
	err := r.withRelations().Where("worker_id = ?", workerID).Order("created_at DESC").Find(&offers).Error
	return offers, err
}

func (r *jobOfferRepository) FindByEmployerID(employerID uuid.UUID) ([]models.JobOffer, error) {
	var offers []models.JobOffer
	err := r.withRelations().Where("employer_id = ?", employerID).Order("created_at DESC").Find(&offers).Error
	return offers, err
}

func (r *jobOfferRepository) FindByID(id uuid.UUID) (*models.JobOffer, error) {
	var offer models.JobOffer
	err := r.withRelations().Where("id = ?", id).First(&offer).Error
	return oneOrNil(&offer, err)
}

func (r *jobOfferRepository) FindByIDForUpdate(id uuid.UUID) (*models.JobOffer, error) {
	var offer models.JobOffer
	err := r.db.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", id).First(&offer).Error
	return oneOrNil(&offer, err)
}

func (r *jobOfferRepository) FindPending(skillPostingID uuid.UUID, employerID uuid.UUID) (*models.JobOffer, error) {
	var offer models.JobOffer
	err := r.db.Where("skill_posting_id = ? AND employer_id = ? AND status = ?",
		skillPostingID, employerID, models.OfferStatusPending).First(&offer).Error
	return oneOrNil(&offer, err)
}

func (r *jobOfferRepository) UpdateStatus(id uuid.UUID, status models.OfferStatus) error {
	return r.db.Model(&models.JobOffer{}).Where("id = ?", id).
		Updates(map[string]any{"status": status, "updated_at": time.Now()}).Error
}
