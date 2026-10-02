package repositories

import (
	"time"

	"layanesia-server/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type SubscriptionRepository interface {
	WithTx(tx *gorm.DB) SubscriptionRepository
	Create(sub *models.Subscription) error
	Update(sub *models.Subscription) error
	FindByID(id uuid.UUID) (*models.Subscription, error)
	FindByOrderID(orderID string) (*models.Subscription, error)
	FindByOrderIDForUpdate(orderID string) (*models.Subscription, error)
	FindActiveByUserAndPlan(userID uuid.UUID, planType models.SubscriptionPlan) (*models.Subscription, error)
	FindPendingByUserAndPlan(userID uuid.UUID, planType models.SubscriptionPlan) (*models.Subscription, error)
	FindByUserID(userID uuid.UUID) ([]models.Subscription, error)
	ExpireStale(userID *uuid.UUID) (int64, error)
}

type subscriptionRepository struct {
	db *gorm.DB
}

func NewSubscriptionRepository(db *gorm.DB) SubscriptionRepository {
	return &subscriptionRepository{db: db}
}

func (r *subscriptionRepository) WithTx(tx *gorm.DB) SubscriptionRepository {
	return &subscriptionRepository{db: tx}
}

func (r *subscriptionRepository) Create(sub *models.Subscription) error {
	return r.db.Create(sub).Error
}

func (r *subscriptionRepository) Update(sub *models.Subscription) error {
	return r.db.Save(sub).Error
}

func (r *subscriptionRepository) FindByID(id uuid.UUID) (*models.Subscription, error) {
	var sub models.Subscription
	err := r.db.Where("id = ?", id).First(&sub).Error
	return oneOrNil(&sub, err)
}

func (r *subscriptionRepository) FindByOrderID(orderID string) (*models.Subscription, error) {
	var sub models.Subscription
	err := r.db.Where("order_id = ?", orderID).First(&sub).Error
	return oneOrNil(&sub, err)
}

// FindByOrderIDForUpdate mengunci baris (SELECT ... FOR UPDATE) agar webhook yang
// datang bersamaan untuk order yang sama diproses berurutan.
func (r *subscriptionRepository) FindByOrderIDForUpdate(orderID string) (*models.Subscription, error) {
	var sub models.Subscription
	err := r.db.Clauses(clause.Locking{Strength: "UPDATE"}).Where("order_id = ?", orderID).First(&sub).Error
	return oneOrNil(&sub, err)
}

func (r *subscriptionRepository) FindActiveByUserAndPlan(userID uuid.UUID, planType models.SubscriptionPlan) (*models.Subscription, error) {
	var sub models.Subscription
	now := time.Now()
	err := r.db.
		Where("user_id = ? AND plan_type = ? AND status = ? AND starts_at <= ? AND expires_at > ?",
			userID, planType, models.SubStatusActive, now, now).
		Order("expires_at DESC").
		First(&sub).Error
	return oneOrNil(&sub, err)
}

func (r *subscriptionRepository) FindPendingByUserAndPlan(userID uuid.UUID, planType models.SubscriptionPlan) (*models.Subscription, error) {
	var sub models.Subscription
	err := r.db.
		Where("user_id = ? AND plan_type = ? AND status = ?", userID, planType, models.SubStatusPending).
		Order("created_at DESC").
		First(&sub).Error
	return oneOrNil(&sub, err)
}

func (r *subscriptionRepository) FindByUserID(userID uuid.UUID) ([]models.Subscription, error) {
	var subs []models.Subscription
	err := r.db.Where("user_id = ?", userID).Order("created_at DESC").Find(&subs).Error
	return subs, err
}

// ExpireStale menandai langganan ACTIVE yang sudah lewat masa berlakunya sebagai EXPIRED.
// Jika userID nil, berlaku untuk semua user (dipakai job berkala).
func (r *subscriptionRepository) ExpireStale(userID *uuid.UUID) (int64, error) {
	q := r.db.Model(&models.Subscription{}).
		Where("status = ? AND expires_at <= ?", models.SubStatusActive, time.Now())
	if userID != nil {
		q = q.Where("user_id = ?", *userID)
	}
	res := q.Updates(map[string]any{"status": models.SubStatusExpired, "updated_at": time.Now()})
	return res.RowsAffected, res.Error
}
