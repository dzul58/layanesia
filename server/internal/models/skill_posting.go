package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SkillCategory string

const (
	CategorySerabutan   SkillCategory = "SERABUTAN"
	CategoryProfesional SkillCategory = "PROFESIONAL"
)

type RateType string

const (
	RatePerHour RateType = "PER_HOUR"
	RatePerDay  RateType = "PER_DAY"
)

type AvailabilityStatus string

const (
	AvailAvailable AvailabilityStatus = "AVAILABLE"
	AvailBusy      AvailabilityStatus = "BUSY"
)

func (r RateType) Valid() bool           { return r == RatePerHour || r == RatePerDay }
func (a AvailabilityStatus) Valid() bool { return a == AvailAvailable || a == AvailBusy }

type SkillPosting struct {
	ID               uuid.UUID          `gorm:"type:uuid;primaryKey" json:"id"`
	UserID           uuid.UUID          `gorm:"type:uuid;not null;index" json:"user_id"`
	User             *User              `gorm:"foreignKey:UserID" json:"-"`
	Category         SkillCategory      `gorm:"type:varchar(50);not null" json:"category"`
	Title            string             `gorm:"type:varchar(255);not null" json:"title"`
	Description      string             `gorm:"type:text;not null" json:"description"`
	Country          string             `gorm:"type:varchar(100);default:'Indonesia'" json:"country"`
	Province         string             `gorm:"type:varchar(100);not null" json:"province"`
	City             string             `gorm:"type:varchar(100);not null" json:"city"`
	District         string             `gorm:"type:varchar(100);not null" json:"district"`
	AddressDetail    *string            `gorm:"type:text" json:"address_detail,omitempty"`
	Latitude         *float64           `gorm:"type:decimal(10,8)" json:"latitude,omitempty"`
	Longitude        *float64           `gorm:"type:decimal(11,8)" json:"longitude,omitempty"`
	FormattedAddress *string            `gorm:"type:text" json:"formatted_address,omitempty"`
	OSMPlaceID       *string            `gorm:"type:varchar(255)" json:"osm_place_id,omitempty"`
	RateType         RateType           `gorm:"type:varchar(30);not null" json:"rate_type"`
	RateAmount       float64            `gorm:"type:numeric(12,2);not null" json:"rate_amount"`
	Availability     AvailabilityStatus `gorm:"type:varchar(30);default:'AVAILABLE'" json:"availability"`
	CreatedAt        time.Time          `json:"created_at"`
	UpdatedAt        time.Time          `json:"updated_at"`
}

func (sp *SkillPosting) BeforeCreate(tx *gorm.DB) (err error) {
	if sp.ID == uuid.Nil {
		sp.ID = uuid.New()
	}
	if sp.Country == "" {
		sp.Country = "Indonesia"
	}
	if sp.Availability == "" {
		sp.Availability = AvailAvailable
	}
	return
}
