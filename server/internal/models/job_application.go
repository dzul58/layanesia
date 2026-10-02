package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ApplicationStatus string

const (
	AppStatusPending  ApplicationStatus = "PENDING"
	AppStatusAccepted ApplicationStatus = "ACCEPTED"
	AppStatusRejected ApplicationStatus = "REJECTED"
)

type JobApplication struct {
	ID           uuid.UUID         `gorm:"type:uuid;primaryKey" json:"id"`
	JobPostingID uuid.UUID         `gorm:"type:uuid;not null;uniqueIndex:idx_job_applicant" json:"job_posting_id"`
	JobPosting   *JobPosting       `gorm:"foreignKey:JobPostingID" json:"-"`
	ApplicantID  uuid.UUID         `gorm:"type:uuid;not null;uniqueIndex:idx_job_applicant" json:"applicant_id"`
	Applicant    *User             `gorm:"foreignKey:ApplicantID" json:"-"`
	Status       ApplicationStatus `gorm:"type:varchar(30);default:'PENDING'" json:"status"`
	CreatedAt    time.Time         `json:"created_at"`
	UpdatedAt    time.Time         `json:"updated_at"`
}

func (ja *JobApplication) BeforeCreate(tx *gorm.DB) (err error) {
	if ja.ID == uuid.Nil {
		ja.ID = uuid.New()
	}
	if ja.Status == "" {
		ja.Status = AppStatusPending
	}
	return
}
