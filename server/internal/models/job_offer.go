package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OfferStatus string

const (
	OfferStatusPending  OfferStatus = "PENDING"
	OfferStatusAccepted OfferStatus = "ACCEPTED"
	OfferStatusRejected OfferStatus = "REJECTED"
)

func (s OfferStatus) IsResponse() bool { return s == OfferStatusAccepted || s == OfferStatusRejected }

type JobOffer struct {
	ID             uuid.UUID     `gorm:"type:uuid;primaryKey" json:"id"`
	SkillPostingID uuid.UUID     `gorm:"type:uuid;not null;index" json:"skill_posting_id"`
	SkillPosting   *SkillPosting `gorm:"foreignKey:SkillPostingID" json:"-"`
	EmployerID     uuid.UUID     `gorm:"type:uuid;not null;index" json:"employer_id"`
	Employer       *User         `gorm:"foreignKey:EmployerID" json:"-"`
	WorkerID       uuid.UUID     `gorm:"type:uuid;not null;index" json:"worker_id"`
	Worker         *User         `gorm:"foreignKey:WorkerID" json:"-"`
	OfferedBudget  float64       `gorm:"type:numeric(12,2);not null" json:"offered_budget"`
	WorkDate       time.Time     `gorm:"type:date;not null" json:"work_date"`
	Status         OfferStatus   `gorm:"type:varchar(30);default:'PENDING'" json:"status"`
	CreatedAt      time.Time     `json:"created_at"`
	UpdatedAt      time.Time     `json:"updated_at"`
}

func (jo *JobOffer) BeforeCreate(tx *gorm.DB) (err error) {
	if jo.ID == uuid.Nil {
		jo.ID = uuid.New()
	}
	if jo.Status == "" {
		jo.Status = OfferStatusPending
	}
	return
}
