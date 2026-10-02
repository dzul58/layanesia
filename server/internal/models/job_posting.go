package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DurationType string

const (
	DurationHours DurationType = "HOURS"
	DurationDays  DurationType = "DAYS"
)

type JobStatus string

const (
	JobStatusOpen       JobStatus = "OPEN"
	JobStatusInProgress JobStatus = "IN_PROGRESS"
	JobStatusDone       JobStatus = "DONE"
	JobStatusCancelled  JobStatus = "CANCELLED"
)

func (s JobStatus) Valid() bool {
	switch s {
	case JobStatusOpen, JobStatusInProgress, JobStatusDone, JobStatusCancelled:
		return true
	}
	return false
}

func (c SkillCategory) Valid() bool { return c == CategorySerabutan || c == CategoryProfesional }
func (d DurationType) Valid() bool  { return d == DurationHours || d == DurationDays }

type JobPosting struct {
	ID               uuid.UUID     `gorm:"type:uuid;primaryKey" json:"id"`
	EmployerID       uuid.UUID     `gorm:"type:uuid;not null;index" json:"employer_id"`
	Employer         *User         `gorm:"foreignKey:EmployerID" json:"-"`
	Category         SkillCategory `gorm:"type:varchar(50);not null" json:"category"`
	Title            string        `gorm:"type:varchar(255);not null" json:"title"`
	Description      string        `gorm:"type:text;not null" json:"description"`
	Country          string        `gorm:"type:varchar(100);default:'Indonesia'" json:"country"`
	Province         string        `gorm:"type:varchar(100);not null" json:"province"`
	City             string        `gorm:"type:varchar(100);not null" json:"city"`
	District         string        `gorm:"type:varchar(100);not null" json:"district"`
	AddressDetail    *string       `gorm:"type:text" json:"address_detail,omitempty"`
	Latitude         *float64      `gorm:"type:decimal(10,8)" json:"latitude,omitempty"`
	Longitude        *float64      `gorm:"type:decimal(11,8)" json:"longitude,omitempty"`
	FormattedAddress *string       `gorm:"type:text" json:"formatted_address,omitempty"`
	OSMPlaceID       *string       `gorm:"type:varchar(255)" json:"osm_place_id,omitempty"`
	WorkDate         time.Time     `gorm:"type:date;not null" json:"work_date"`
	DurationType     DurationType  `gorm:"type:varchar(30);not null" json:"duration_type"`
	DurationValue    int           `gorm:"type:int;not null" json:"duration_value"`
	Budget           float64       `gorm:"type:numeric(12,2);not null" json:"budget"`
	Status           JobStatus     `gorm:"type:varchar(30);default:'OPEN'" json:"status"`
	ClosedAt         *time.Time    `json:"closed_at,omitempty"`
	CreatedAt        time.Time     `json:"created_at"`
	UpdatedAt        time.Time     `json:"updated_at"`
}

func (jp *JobPosting) BeforeCreate(tx *gorm.DB) (err error) {
	if jp.ID == uuid.Nil {
		jp.ID = uuid.New()
	}
	if jp.Country == "" {
		jp.Country = "Indonesia"
	}
	if jp.Status == "" {
		jp.Status = JobStatusOpen
	}
	return
}
