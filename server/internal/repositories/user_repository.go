package repositories

import (
	"errors"

	"layanesia-server/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository interface {
	WithTx(tx *gorm.DB) UserRepository
	Create(user *models.User) error
	FindByEmail(email string) (*models.User, error)
	FindByPhone(phone string) (*models.User, error)
	FindByID(id uuid.UUID) (*models.User, error)
	Update(user *models.User) error
	UpdateActiveMode(id uuid.UUID, mode models.UserMode) (*models.User, error)
	FindByKTPStatus(status models.KTPStatus, limit, offset int) ([]models.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) WithTx(tx *gorm.DB) UserRepository { return &userRepository{db: tx} }

func (r *userRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Where("LOWER(email) = LOWER(?)", email).First(&user).Error
	return oneOrNil(&user, err)
}

func (r *userRepository) FindByPhone(phone string) (*models.User, error) {
	var user models.User
	err := r.db.Where("phone = ?", phone).First(&user).Error
	return oneOrNil(&user, err)
}

func (r *userRepository) FindByID(id uuid.UUID) (*models.User, error) {
	var user models.User
	err := r.db.Where("id = ?", id).First(&user).Error
	return oneOrNil(&user, err)
}

func (r *userRepository) Update(user *models.User) error {
	return r.db.Save(user).Error
}

func (r *userRepository) UpdateActiveMode(id uuid.UUID, mode models.UserMode) (*models.User, error) {
	res := r.db.Model(&models.User{}).Where("id = ?", id).Update("active_mode", mode)
	if res.Error != nil {
		return nil, res.Error
	}
	if res.RowsAffected == 0 {
		return nil, nil
	}
	return r.FindByID(id)
}

func (r *userRepository) FindByKTPStatus(status models.KTPStatus, limit, offset int) ([]models.User, error) {
	var users []models.User
	err := r.db.Where("ktp_status = ?", status).
		Order("ktp_submitted_at ASC NULLS LAST").
		Limit(clampLimit(limit, 50, 200)).Offset(maxInt(offset, 0)).
		Find(&users).Error
	return users, err
}

// oneOrNil mengubah gorm.ErrRecordNotFound menjadi (nil, nil).
func oneOrNil[T any](v *T, err error) (*T, error) {
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return v, nil
}

func clampLimit(v, def, max int) int {
	if v <= 0 {
		return def
	}
	if v > max {
		return max
	}
	return v
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
