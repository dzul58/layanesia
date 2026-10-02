package repositories

import (
	"strings"
	"time"

	"layanesia-server/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type JobFilter struct {
	Province string
	City     string
	District string
	Category string
	Search   string
	Status   string
	Limit    int
	Offset   int
}

type JobPostingRepository interface {
	WithTx(tx *gorm.DB) JobPostingRepository
	Create(job *models.JobPosting) error
	FindByEmployerID(employerID uuid.UUID) ([]models.JobPosting, error)
	FindAll(filter JobFilter) ([]models.JobPosting, int64, error)
	FindByID(id uuid.UUID) (*models.JobPosting, error)
	Update(job *models.JobPosting) error
	UpdateStatus(id uuid.UUID, status models.JobStatus, closedAt *time.Time) error
	Delete(id uuid.UUID, employerID uuid.UUID) (int64, error)
}

type jobPostingRepository struct {
	db *gorm.DB
}

func NewJobPostingRepository(db *gorm.DB) JobPostingRepository {
	return &jobPostingRepository{db: db}
}

func (r *jobPostingRepository) WithTx(tx *gorm.DB) JobPostingRepository {
	return &jobPostingRepository{db: tx}
}

func (r *jobPostingRepository) Create(job *models.JobPosting) error {
	return r.db.Create(job).Error
}

func (r *jobPostingRepository) FindByEmployerID(employerID uuid.UUID) ([]models.JobPosting, error) {
	var jobs []models.JobPosting
	err := r.db.Preload("Employer").
		Where("employer_id = ?", employerID).
		Order("created_at DESC").
		Find(&jobs).Error
	return jobs, err
}

func (r *jobPostingRepository) FindAll(filter JobFilter) ([]models.JobPosting, int64, error) {
	base := r.db.Model(&models.JobPosting{})
	base = applyLocationFilter(base, filter.Province, filter.City, filter.District)
	if filter.Category != "" {
		base = base.Where("category = ?", filter.Category)
	}
	// Pencarian publik hanya menampilkan lowongan OPEN kecuali status diminta eksplisit.
	if filter.Status != "" {
		base = base.Where("status = ?", filter.Status)
	} else {
		base = base.Where("status = ?", models.JobStatusOpen)
	}
	if s := sanitizeSearch(filter.Search); s != "" {
		base = base.Where("(title ILIKE ? ESCAPE '\\' OR description ILIKE ? ESCAPE '\\')", s, s)
	}

	var total int64
	if err := base.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var jobs []models.JobPosting
	err := base.Preload("Employer").
		Limit(clampLimit(filter.Limit, 20, 100)).Offset(maxInt(filter.Offset, 0)).
		Order("created_at DESC").
		Find(&jobs).Error
	return jobs, total, err
}

func (r *jobPostingRepository) FindByID(id uuid.UUID) (*models.JobPosting, error) {
	var job models.JobPosting
	err := r.db.Preload("Employer").Where("id = ?", id).First(&job).Error
	return oneOrNil(&job, err)
}

func (r *jobPostingRepository) Update(job *models.JobPosting) error {
	return r.db.Save(job).Error
}

func (r *jobPostingRepository) UpdateStatus(id uuid.UUID, status models.JobStatus, closedAt *time.Time) error {
	updates := map[string]any{"status": status, "updated_at": time.Now()}
	if closedAt != nil {
		updates["closed_at"] = *closedAt
	}
	return r.db.Model(&models.JobPosting{}).Where("id = ?", id).Updates(updates).Error
}

func (r *jobPostingRepository) Delete(id uuid.UUID, employerID uuid.UUID) (int64, error) {
	res := r.db.Where("id = ? AND employer_id = ?", id, employerID).Delete(&models.JobPosting{})
	return res.RowsAffected, res.Error
}

// sanitizeSearch membatasi panjang, meng-escape wildcard LIKE, dan membungkus dengan %.
func sanitizeSearch(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return ""
	}
	if len(s) > 100 {
		s = s[:100]
	}
	s = strings.NewReplacer(`\`, `\\`, `%`, `\%`, `_`, `\_`).Replace(s)
	return "%" + s + "%"
}
