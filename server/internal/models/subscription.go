package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SubscriptionPlan string

const (
	PlanApplyJob5K   SubscriptionPlan = "APPLY_JOB_5K"
	PlanPostSkill10K SubscriptionPlan = "POST_SKILL_10K"
)

func (p SubscriptionPlan) Valid() bool { return p == PlanApplyJob5K || p == PlanPostSkill10K }

// Price mengembalikan harga paket dalam Rupiah (integer, Midtrans tidak menerima desimal).
func (p SubscriptionPlan) Price() int64 {
	switch p {
	case PlanApplyJob5K:
		return 5000
	case PlanPostSkill10K:
		return 10000
	}
	return 0
}

func (p SubscriptionPlan) DisplayName() string {
	switch p {
	case PlanApplyJob5K:
		return "Paket Melamar Lowongan (30 Hari)"
	case PlanPostSkill10K:
		return "Paket Etalase Keahlian (30 Hari)"
	}
	return string(p)
}

// SubscriptionDuration adalah masa aktif satu kali pembelian paket.
const SubscriptionDuration = 30 * 24 * time.Hour

type SubscriptionStatus string

const (
	SubStatusPending   SubscriptionStatus = "PENDING" // menunggu pembayaran
	SubStatusActive    SubscriptionStatus = "ACTIVE"
	SubStatusExpired   SubscriptionStatus = "EXPIRED"
	SubStatusCancelled SubscriptionStatus = "CANCELLED"
)

type PaymentStatus string

const (
	PaymentPending  PaymentStatus = "PENDING"
	PaymentPaid     PaymentStatus = "PAID"
	PaymentFailed   PaymentStatus = "FAILED"
	PaymentExpired  PaymentStatus = "EXPIRED"
	PaymentRefunded PaymentStatus = "REFUNDED"
)

type Subscription struct {
	ID       uuid.UUID        `gorm:"type:uuid;primaryKey" json:"id"`
	UserID   uuid.UUID        `gorm:"type:uuid;not null;index" json:"user_id"`
	User     *User            `gorm:"foreignKey:UserID" json:"-"`
	PlanType SubscriptionPlan `gorm:"type:varchar(50);not null" json:"plan_type"`
	Amount   float64          `gorm:"type:numeric(12,2);not null" json:"amount"`

	Status          SubscriptionStatus `gorm:"type:varchar(30);not null;default:'PENDING'" json:"status"`
	PaymentStatus   PaymentStatus      `gorm:"type:varchar(30);not null;default:'PENDING'" json:"payment_status"`
	PaymentProvider string             `gorm:"type:varchar(30);not null;default:'MIDTRANS'" json:"payment_provider"`
	PaymentType     *string            `gorm:"type:varchar(50)" json:"payment_type,omitempty"`
	OrderID         string             `gorm:"type:varchar(64);not null;uniqueIndex" json:"order_id"`
	SnapToken       *string            `gorm:"type:varchar(255)" json:"snap_token,omitempty"`
	RedirectURL     *string            `gorm:"type:text" json:"redirect_url,omitempty"`
	PaidAt          *time.Time         `json:"paid_at,omitempty"`

	StartsAt  time.Time `gorm:"not null" json:"starts_at"`
	ExpiresAt time.Time `gorm:"not null" json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (s *Subscription) BeforeCreate(tx *gorm.DB) (err error) {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	if s.Status == "" {
		s.Status = SubStatusPending
	}
	if s.PaymentStatus == "" {
		s.PaymentStatus = PaymentPending
	}
	if s.PaymentProvider == "" {
		s.PaymentProvider = "MIDTRANS"
	}
	return nil
}

// IsActiveAt melaporkan apakah langganan sedang berlaku pada waktu t.
func (s *Subscription) IsActiveAt(t time.Time) bool {
	return s.Status == SubStatusActive && !t.Before(s.StartsAt) && t.Before(s.ExpiresAt)
}
