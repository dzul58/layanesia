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

type CreateOfferRequest struct {
	SkillPostingID uuid.UUID `json:"skill_posting_id" validate:"required"`
	OfferedBudget  float64   `json:"offered_budget" validate:"required,gt=0,lte=999999999"`
	WorkDate       string    `json:"work_date" validate:"omitempty"` // YYYY-MM-DD atau RFC3339
	Message        *string   `json:"message" validate:"omitempty,max=1000"`
}

type RespondOfferRequest struct {
	Status models.OfferStatus `json:"status" validate:"required,oneof=ACCEPTED REJECTED"`
}

type JobOfferService interface {
	CreateOffer(employerID uuid.UUID, req CreateOfferRequest) (*models.JobOffer, error)
	GetReceivedOffers(workerID uuid.UUID) ([]models.JobOffer, error)
	GetSentOffers(employerID uuid.UUID) ([]models.JobOffer, error)
	RespondToOffer(workerID uuid.UUID, offerID uuid.UUID, status models.OfferStatus) (*models.JobOffer, error)
}

type jobOfferService struct {
	db        *gorm.DB
	offerRepo repositories.JobOfferRepository
	skillRepo repositories.SkillPostingRepository
	userRepo  repositories.UserRepository
	connRepo  repositories.ConnectionRepository
}

func NewJobOfferService(
	db *gorm.DB,
	offerRepo repositories.JobOfferRepository,
	skillRepo repositories.SkillPostingRepository,
	userRepo repositories.UserRepository,
	connRepo repositories.ConnectionRepository,
) JobOfferService {
	return &jobOfferService{db: db, offerRepo: offerRepo, skillRepo: skillRepo, userRepo: userRepo, connRepo: connRepo}
}

// CreateOffer: pemberi kerja wajib KTP terverifikasi (§7.1 ERD), tidak boleh ke
// postingan sendiri, dan hanya satu tawaran PENDING per (skill_posting, employer).
func (s *jobOfferService) CreateOffer(employerID uuid.UUID, req CreateOfferRequest) (*models.JobOffer, error) {
	if err := utils.ValidateStruct(req); err != nil {
		return nil, err
	}

	employer, err := s.userRepo.FindByID(employerID)
	if err != nil {
		return nil, apperrors.Internal(err)
	}
	if employer == nil {
		return nil, apperrors.ErrUserNotFound
	}
	if !employer.IsKTPVerified {
		return nil, apperrors.ErrKTPNotVerified
	}

	skill, err := s.skillRepo.FindByID(req.SkillPostingID)
	if err != nil {
		return nil, apperrors.Internal(err)
	}
	if skill == nil {
		return nil, apperrors.NotFound("Postingan keahlian tidak ditemukan")
	}
	if skill.UserID == employerID {
		return nil, apperrors.Forbidden("SELF_OFFER", "Anda tidak dapat mengirim tawaran ke postingan keahlian Anda sendiri.")
	}
	if skill.Availability != models.AvailAvailable {
		return nil, apperrors.Conflict("WORKER_BUSY", "Pekerja ini sedang tidak tersedia.")
	}

	workDate, err := parseWorkDate(req.WorkDate, time.Now())
	if err != nil {
		return nil, err
	}

	pending, err := s.offerRepo.FindPending(req.SkillPostingID, employerID)
	if err != nil {
		return nil, apperrors.Internal(err)
	}
	if pending != nil {
		return nil, apperrors.Conflict("OFFER_ALREADY_PENDING", "Anda sudah memiliki tawaran yang menunggu respon pada postingan ini.")
	}

	offer := &models.JobOffer{
		SkillPostingID: req.SkillPostingID,
		EmployerID:     employerID,
		WorkerID:       skill.UserID,
		OfferedBudget:  req.OfferedBudget,
		WorkDate:       workDate,
		Status:         models.OfferStatusPending,
	}
	if err := s.offerRepo.Create(offer); err != nil {
		// Race pada partial unique index uq_job_offers_pending_per_employer.
		return nil, apperrors.Conflict("OFFER_ALREADY_PENDING", "Anda sudah memiliki tawaran yang menunggu respon pada postingan ini.").WithCause(err)
	}
	created, err := s.offerRepo.FindByID(offer.ID)
	if err != nil {
		return nil, apperrors.Internal(err)
	}
	return created, nil
}

func (s *jobOfferService) GetReceivedOffers(workerID uuid.UUID) ([]models.JobOffer, error) {
	items, err := s.offerRepo.FindByWorkerID(workerID)
	if err != nil {
		return nil, apperrors.Internal(err)
	}
	return items, nil
}

func (s *jobOfferService) GetSentOffers(employerID uuid.UUID) ([]models.JobOffer, error) {
	items, err := s.offerRepo.FindByEmployerID(employerID)
	if err != nil {
		return nil, apperrors.Internal(err)
	}
	return items, nil
}

// RespondToOffer dijalankan dalam transaksi; ACCEPTED membuat Connection.
func (s *jobOfferService) RespondToOffer(workerID uuid.UUID, offerID uuid.UUID, status models.OfferStatus) (*models.JobOffer, error) {
	if !status.IsResponse() {
		return nil, apperrors.Validation("Status respon tidak valid.", map[string]any{"status": "ACCEPTED atau REJECTED"})
	}

	err := s.db.Transaction(func(tx *gorm.DB) error {
		offerRepo := s.offerRepo.WithTx(tx)
		connRepo := s.connRepo.WithTx(tx)

		offer, err := offerRepo.FindByIDForUpdate(offerID)
		if err != nil {
			return apperrors.Internal(err)
		}
		if offer == nil {
			return apperrors.NotFound("Tawaran pekerjaan tidak ditemukan")
		}
		if offer.WorkerID != workerID {
			return apperrors.Forbidden("NOT_RECIPIENT", "Anda tidak memiliki akses untuk merespon tawaran ini.")
		}
		if offer.Status != models.OfferStatusPending {
			return apperrors.Conflict("OFFER_ALREADY_RESPONDED", "Tawaran ini telah direspon sebelumnya.")
		}
		if err := offerRepo.UpdateStatus(offerID, status); err != nil {
			return apperrors.Internal(err)
		}
		if status != models.OfferStatusAccepted {
			return nil
		}
		conn := &models.Connection{
			SourceType: models.SourceJobOffer,
			SourceID:   offer.ID,
			EmployerID: offer.EmployerID,
			WorkerID:   offer.WorkerID,
			Status:     models.ConnStatusActive,
		}
		if err := connRepo.Create(conn); err != nil {
			return apperrors.Internal(err)
		}
		return nil
	})
	if err != nil {
		if ae, ok := apperrors.As(err); ok {
			return nil, ae
		}
		return nil, apperrors.Internal(err)
	}

	updated, err := s.offerRepo.FindByID(offerID)
	if err != nil {
		return nil, apperrors.Internal(err)
	}
	return updated, nil
}
