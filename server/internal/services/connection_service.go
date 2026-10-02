package services

import (
	"time"

	"layanesia-server/internal/apperrors"
	"layanesia-server/internal/models"
	"layanesia-server/internal/repositories"
	"layanesia-server/internal/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UserContact adalah kontak pihak lawan; hanya dibuka kepada peserta koneksi
// (setelah lamaran/tawaran diterima) — ini satu-satunya tempat nomor telepon
// pihak lain terekspos.
type UserContact struct {
	ID               uuid.UUID `json:"id"`
	Name             string    `json:"name"`
	Phone            string    `json:"phone"`
	WhatsAppLink     string    `json:"whatsapp_link"`
	Province         string    `json:"province"`
	City             string    `json:"city"`
	District         string    `json:"district"`
	AddressDetail    *string   `json:"address_detail,omitempty"`
	FormattedAddress *string   `json:"formatted_address,omitempty"`
	IsKTPVerified    bool      `json:"is_ktp_verified"`
}

type ConnectionDetailResponse struct {
	Connection models.Connection `json:"connection"`
	Role       models.ReviewerRole `json:"my_role"` // EMPLOYER | WORKER
	Contact    *UserContact        `json:"contact"`
}

type ConnectionService interface {
	GetConnections(userID uuid.UUID, status *models.ConnectionStatus) ([]ConnectionDetailResponse, error)
	GetConnectionByID(userID uuid.UUID, id uuid.UUID) (*ConnectionDetailResponse, error)
	CompleteConnection(userID uuid.UUID, id uuid.UUID) (*ConnectionDetailResponse, error)
}

type connectionService struct {
	db       *gorm.DB
	connRepo repositories.ConnectionRepository
	jobRepo  repositories.JobPostingRepository
	appRepo  repositories.JobApplicationRepository
}

func NewConnectionService(
	db *gorm.DB,
	connRepo repositories.ConnectionRepository,
	jobRepo repositories.JobPostingRepository,
	appRepo repositories.JobApplicationRepository,
) ConnectionService {
	return &connectionService{db: db, connRepo: connRepo, jobRepo: jobRepo, appRepo: appRepo}
}

func toContact(u *models.User) *UserContact {
	if u == nil {
		return nil
	}
	return &UserContact{
		ID:               u.ID,
		Name:             u.Name,
		Phone:            u.Phone,
		WhatsAppLink:     utils.PhoneToWhatsApp(u.Phone),
		Province:         u.Province,
		City:             u.City,
		District:         u.District,
		AddressDetail:    u.AddressDetail,
		FormattedAddress: u.FormattedAddress,
		IsKTPVerified:    u.IsKTPVerified,
	}
}

func toDetail(conn *models.Connection, userID uuid.UUID) ConnectionDetailResponse {
	role := models.RoleWorker
	if conn.EmployerID == userID {
		role = models.RoleEmployer
	}
	return ConnectionDetailResponse{
		Connection: *conn,
		Role:       role,
		Contact:    toContact(conn.Counterpart(userID)),
	}
}

func (s *connectionService) GetConnections(userID uuid.UUID, status *models.ConnectionStatus) ([]ConnectionDetailResponse, error) {
	if status != nil && *status != models.ConnStatusActive && *status != models.ConnStatusCompleted {
		return nil, apperrors.Validation("Status koneksi tidak valid.", map[string]any{"status": "ACTIVE atau COMPLETED"})
	}
	conns, err := s.connRepo.FindByUserID(userID, status)
	if err != nil {
		return nil, apperrors.Internal(err)
	}
	out := make([]ConnectionDetailResponse, 0, len(conns))
	for i := range conns {
		out = append(out, toDetail(&conns[i], userID))
	}
	return out, nil
}

func (s *connectionService) GetConnectionByID(userID uuid.UUID, id uuid.UUID) (*ConnectionDetailResponse, error) {
	conn, err := s.connRepo.FindByID(id)
	if err != nil {
		return nil, apperrors.Internal(err)
	}
	// 404 juga untuk non-peserta agar keberadaan koneksi tidak bisa di-enumerate.
	if conn == nil || !conn.Involves(userID) {
		return nil, apperrors.NotFound("Koneksi tidak ditemukan")
	}
	d := toDetail(conn, userID)
	return &d, nil
}

// CompleteConnection menandai pekerjaan selesai. Boleh dilakukan salah satu pihak.
// Jika sumbernya lamaran, lowongan terkait ikut ditandai DONE.
func (s *connectionService) CompleteConnection(userID uuid.UUID, id uuid.UUID) (*ConnectionDetailResponse, error) {
	err := s.db.Transaction(func(tx *gorm.DB) error {
		connRepo := s.connRepo.WithTx(tx)
		conn, err := connRepo.FindByIDForUpdate(id)
		if err != nil {
			return apperrors.Internal(err)
		}
		if conn == nil || !conn.Involves(userID) {
			return apperrors.NotFound("Koneksi tidak ditemukan")
		}
		if conn.Status == models.ConnStatusCompleted {
			return apperrors.Conflict("CONNECTION_ALREADY_COMPLETED", "Koneksi ini sudah ditandai selesai.")
		}
		now := time.Now()
		if err := connRepo.MarkCompleted(id, now); err != nil {
			return apperrors.Internal(err)
		}
		if conn.SourceType == models.SourceJobApplication {
			app, err := s.appRepo.WithTx(tx).FindByID(conn.SourceID)
			if err != nil {
				return apperrors.Internal(err)
			}
			if app != nil {
				if err := s.jobRepo.WithTx(tx).UpdateStatus(app.JobPostingID, models.JobStatusDone, &now); err != nil {
					return apperrors.Internal(err)
				}
			}
		}
		return nil
	})
	if err != nil {
		if ae, ok := apperrors.As(err); ok {
			return nil, ae
		}
		return nil, apperrors.Internal(err)
	}
	return s.GetConnectionByID(userID, id)
}
