package repositories

import (
	"layanesia-server/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SkillFilter struct {
	Province     string
	City         string
	District     string
	Category     string
	Search       string
	Availability string
	Limit        int
	Offset       int
}

type SkillPostingRepository interface {
	WithTx(tx *gorm.DB) SkillPostingRepository
	Create(skill *models.SkillPosting) error
	FindByUserID(userID uuid.UUID) ([]models.SkillPosting, error)
	FindAll(filter SkillFilter) ([]models.SkillPosting, int64, error)
	FindByID(id uuid.UUID) (*models.SkillPosting, error)
	Update(skill *models.SkillPosting) error
	Delete(id uuid.UUID, userID uuid.UUID) (int64, error)
}

type skillPostingRepository struct {
	db *gorm.DB
}

func NewSkillPostingRepository(db *gorm.DB) SkillPostingRepository {
	return &skillPostingRepository{db: db}
}

func (r *skillPostingRepository) WithTx(tx *gorm.DB) SkillPostingRepository {
	return &skillPostingRepository{db: tx}
}

func (r *skillPostingRepository) Create(skill *models.SkillPosting) error {
	return r.db.Create(skill).Error
}

func (r *skillPostingRepository) FindByUserID(userID uuid.UUID) ([]models.SkillPosting, error) {
	var skills []models.SkillPosting
	err := r.db.Preload("User").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&skills).Error
	return skills, err
}

// FindAll mengembalikan hasil + total (untuk pagination). Filter lokasi memakai
// equality case-insensitive karena nilainya berasal dari master data lokasi.
func (r *skillPostingRepository) FindAll(filter SkillFilter) ([]models.SkillPosting, int64, error) {
	base := r.db.Model(&models.SkillPosting{})
	base = applyLocationFilter(base, filter.Province, filter.City, filter.District)
	if filter.Category != "" {
		base = base.Where("category = ?", filter.Category)
	}
	if filter.Availability != "" {
		base = base.Where("availability = ?", filter.Availability)
	}
	if s := sanitizeSearch(filter.Search); s != "" {
		base = base.Where("(title ILIKE ? ESCAPE '\\' OR description ILIKE ? ESCAPE '\\')", s, s)
	}

	var total int64
	if err := base.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var skills []models.SkillPosting
	err := base.Preload("User").
		Limit(clampLimit(filter.Limit, 20, 100)).Offset(maxInt(filter.Offset, 0)).
		Order("created_at DESC").
		Find(&skills).Error
	return skills, total, err
}

func (r *skillPostingRepository) FindByID(id uuid.UUID) (*models.SkillPosting, error) {
	var skill models.SkillPosting
	err := r.db.Preload("User").Where("id = ?", id).First(&skill).Error
	return oneOrNil(&skill, err)
}

func (r *skillPostingRepository) Update(skill *models.SkillPosting) error {
	return r.db.Save(skill).Error
}

func (r *skillPostingRepository) Delete(id uuid.UUID, userID uuid.UUID) (int64, error) {
	res := r.db.Where("id = ? AND user_id = ?", id, userID).Delete(&models.SkillPosting{})
	return res.RowsAffected, res.Error
}

// applyLocationFilter menambahkan filter provinsi/kota/kecamatan (case-insensitive equality).
func applyLocationFilter(q *gorm.DB, province, city, district string) *gorm.DB {
	if province != "" {
		q = q.Where("LOWER(province) = LOWER(?)", province)
	}
	if city != "" {
		q = q.Where("LOWER(city) = LOWER(?)", city)
	}
	if district != "" {
		q = q.Where("LOWER(district) = LOWER(?)", district)
	}
	return q
}
