package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ReviewerRole string

const (
	RoleEmployer ReviewerRole = "EMPLOYER"
	RoleWorker   ReviewerRole = "WORKER"
)

type Rating struct {
	ID           uuid.UUID    `gorm:"type:uuid;primaryKey" json:"id"`
	ConnectionID uuid.UUID    `gorm:"type:uuid;not null;uniqueIndex:idx_rating_conn_rev" json:"connection_id"`
	Connection   *Connection  `gorm:"foreignKey:ConnectionID" json:"-"`
	ReviewerID   uuid.UUID    `gorm:"type:uuid;not null;uniqueIndex:idx_rating_conn_rev" json:"reviewer_id"`
	Reviewer     *User        `gorm:"foreignKey:ReviewerID" json:"-"`
	RevieweeID   uuid.UUID    `gorm:"type:uuid;not null;index" json:"reviewee_id"`
	Reviewee     *User        `gorm:"foreignKey:RevieweeID" json:"-"`
	ReviewerRole ReviewerRole `gorm:"type:varchar(30);not null" json:"reviewer_role"`
	RatingStars  int          `gorm:"type:int;not null" json:"rating_stars"`
	Comment      *string      `gorm:"type:text" json:"comment,omitempty"`
	CreatedAt    time.Time    `json:"created_at"`
}

func (r *Rating) BeforeCreate(tx *gorm.DB) (err error) {
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}
	return
}
