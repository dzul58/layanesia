package repositories

import (
	"time"

	"layanesia-server/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ConnectionRepository interface {
	WithTx(tx *gorm.DB) ConnectionRepository
	Create(conn *models.Connection) error
	FindByUserID(userID uuid.UUID, status *models.ConnectionStatus) ([]models.Connection, error)
	FindByID(id uuid.UUID) (*models.Connection, error)
	FindByIDForUpdate(id uuid.UUID) (*models.Connection, error)
	FindBySource(sourceType models.SourceType, sourceID uuid.UUID) (*models.Connection, error)
	MarkCompleted(id uuid.UUID, at time.Time) error
}

type connectionRepository struct {
	db *gorm.DB
}

func NewConnectionRepository(db *gorm.DB) ConnectionRepository {
	return &connectionRepository{db: db}
}

func (r *connectionRepository) WithTx(tx *gorm.DB) ConnectionRepository {
	return &connectionRepository{db: tx}
}

// Create membuat koneksi; duplikasi per sumber dicegah oleh unique index
// uq_connections_source dan dianggap sukses (idempoten).
func (r *connectionRepository) Create(conn *models.Connection) error {
	return r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "source_type"}, {Name: "source_id"}},
		DoNothing: true,
	}).Create(conn).Error
}

func (r *connectionRepository) withRelations() *gorm.DB {
	return r.db.Preload("Employer").Preload("Worker")
}

func (r *connectionRepository) FindByUserID(userID uuid.UUID, status *models.ConnectionStatus) ([]models.Connection, error) {
	var conns []models.Connection
	q := r.withRelations().Where("employer_id = ? OR worker_id = ?", userID, userID)
	if status != nil {
		q = q.Where("status = ?", *status)
	}
	err := q.Order("created_at DESC").Find(&conns).Error
	return conns, err
}

func (r *connectionRepository) FindByID(id uuid.UUID) (*models.Connection, error) {
	var conn models.Connection
	err := r.withRelations().Where("id = ?", id).First(&conn).Error
	return oneOrNil(&conn, err)
}

func (r *connectionRepository) FindByIDForUpdate(id uuid.UUID) (*models.Connection, error) {
	var conn models.Connection
	err := r.db.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", id).First(&conn).Error
	return oneOrNil(&conn, err)
}

func (r *connectionRepository) FindBySource(sourceType models.SourceType, sourceID uuid.UUID) (*models.Connection, error) {
	var conn models.Connection
	err := r.withRelations().Where("source_type = ? AND source_id = ?", sourceType, sourceID).First(&conn).Error
	return oneOrNil(&conn, err)
}

func (r *connectionRepository) MarkCompleted(id uuid.UUID, at time.Time) error {
	return r.db.Model(&models.Connection{}).Where("id = ?", id).
		Updates(map[string]any{"status": models.ConnStatusCompleted, "completed_at": at, "updated_at": at}).Error
}
