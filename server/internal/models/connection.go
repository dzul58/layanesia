package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SourceType string

const (
	SourceJobApplication SourceType = "JOB_APPLICATION"
	SourceJobOffer       SourceType = "JOB_OFFER"
)

type ConnectionStatus string

const (
	ConnStatusActive    ConnectionStatus = "ACTIVE"
	ConnStatusCompleted ConnectionStatus = "COMPLETED"
)

type Connection struct {
	ID         uuid.UUID        `gorm:"type:uuid;primaryKey" json:"id"`
	SourceType SourceType       `gorm:"type:varchar(30);not null" json:"source_type"`
	SourceID   uuid.UUID        `gorm:"type:uuid;not null;uniqueIndex:idx_conn_empl_work_src" json:"source_id"`
	EmployerID uuid.UUID        `gorm:"type:uuid;not null;uniqueIndex:idx_conn_empl_work_src" json:"employer_id"`
	Employer   *User            `gorm:"foreignKey:EmployerID" json:"-"`
	WorkerID   uuid.UUID        `gorm:"type:uuid;not null;uniqueIndex:idx_conn_empl_work_src" json:"worker_id"`
	Worker     *User            `gorm:"foreignKey:WorkerID" json:"-"`
	Status      ConnectionStatus `gorm:"type:varchar(30);default:'ACTIVE'" json:"status"`
	CompletedAt *time.Time       `json:"completed_at,omitempty"`
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at"`
}

// Involves melaporkan apakah userID adalah salah satu pihak dalam koneksi.
func (c *Connection) Involves(userID uuid.UUID) bool {
	return c.EmployerID == userID || c.WorkerID == userID
}

// Counterpart mengembalikan pihak lawan dari userID (nil jika userID bukan peserta).
func (c *Connection) Counterpart(userID uuid.UUID) *User {
	switch userID {
	case c.EmployerID:
		return c.Worker
	case c.WorkerID:
		return c.Employer
	}
	return nil
}

func (c *Connection) BeforeCreate(tx *gorm.DB) (err error) {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	if c.Status == "" {
		c.Status = ConnStatusActive
	}
	return
}
