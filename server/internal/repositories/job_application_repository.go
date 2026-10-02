package repositories

import (
	"time"

	"layanesia-server/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type JobApplicationRepository interface {
	WithTx(tx *gorm.DB) JobApplicationRepository
	Create(app *models.JobApplication) error
	FindByJobPostingID(jobPostingID uuid.UUID) ([]models.JobApplication, error)
	FindByApplicantID(applicantID uuid.UUID) ([]models.JobApplication, error)
	FindByID(id uuid.UUID) (*models.JobApplication, error)
	FindByIDForUpdate(id uuid.UUID) (*models.JobApplication, error)
	FindExisting(jobPostingID uuid.UUID, applicantID uuid.UUID) (*models.JobApplication, error)
	UpdateStatus(id uuid.UUID, status models.ApplicationStatus) error
	RejectOtherPending(jobPostingID uuid.UUID, exceptID uuid.UUID) (int64, error)
}

type jobApplicationRepository struct {
	db *gorm.DB
}

func NewJobApplicationRepository(db *gorm.DB) JobApplicationRepository {
	return &jobApplicationRepository{db: db}
}

func (r *jobApplicationRepository) WithTx(tx *gorm.DB) JobApplicationRepository {
	return &jobApplicationRepository{db: tx}
}

func (r *jobApplicationRepository) Create(app *models.JobApplication) error {
	return r.db.Create(app).Error
}

func (r *jobApplicationRepository) FindByJobPostingID(jobPostingID uuid.UUID) ([]models.JobApplication, error) {
	var apps []models.JobApplication
	err := r.db.Preload("Applicant").
		Preload("JobPosting").
		Where("job_posting_id = ?", jobPostingID).
		Order("created_at DESC").
		Find(&apps).Error
	return apps, err
}

func (r *jobApplicationRepository) FindByApplicantID(applicantID uuid.UUID) ([]models.JobApplication, error) {
	var apps []models.JobApplication
	err := r.db.Preload("JobPosting").
		Preload("JobPosting.Employer").
		Where("applicant_id = ?", applicantID).
		Order("created_at DESC").
		Find(&apps).Error
	return apps, err
}

func (r *jobApplicationRepository) FindByID(id uuid.UUID) (*models.JobApplication, error) {
	var app models.JobApplication
	err := r.db.Preload("Applicant").
		Preload("JobPosting").
		Preload("JobPosting.Employer").
		Where("id = ?", id).First(&app).Error
	return oneOrNil(&app, err)
}

func (r *jobApplicationRepository) FindByIDForUpdate(id uuid.UUID) (*models.JobApplication, error) {
	var app models.JobApplication
	err := r.db.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", id).First(&app).Error
	return oneOrNil(&app, err)
}

func (r *jobApplicationRepository) FindExisting(jobPostingID uuid.UUID, applicantID uuid.UUID) (*models.JobApplication, error) {
	var app models.JobApplication
	err := r.db.Where("job_posting_id = ? AND applicant_id = ?", jobPostingID, applicantID).First(&app).Error
	return oneOrNil(&app, err)
}

func (r *jobApplicationRepository) UpdateStatus(id uuid.UUID, status models.ApplicationStatus) error {
	return r.db.Model(&models.JobApplication{}).Where("id = ?", id).
		Updates(map[string]any{"status": status, "updated_at": time.Now()}).Error
}

// RejectOtherPending menolak otomatis lamaran PENDING lain pada lowongan yang sama
// setelah satu pelamar diterima.
func (r *jobApplicationRepository) RejectOtherPending(jobPostingID uuid.UUID, exceptID uuid.UUID) (int64, error) {
	res := r.db.Model(&models.JobApplication{}).
		Where("job_posting_id = ? AND id <> ? AND status = ?", jobPostingID, exceptID, models.AppStatusPending).
		Updates(map[string]any{"status": models.AppStatusRejected, "updated_at": time.Now()})
	return res.RowsAffected, res.Error
}
