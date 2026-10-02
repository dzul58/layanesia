package services

import (
	"layanesia-server/internal/apperrors"
	"layanesia-server/internal/models"
	"layanesia-server/internal/repositories"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UpdateApplicationStatusRequest struct {
	Status models.ApplicationStatus `json:"status" validate:"required,oneof=ACCEPTED REJECTED"`
}

type JobApplicationService interface {
	ApplyToJob(applicantID uuid.UUID, jobPostingID uuid.UUID) (*models.JobApplication, error)
	GetApplicantsForJob(employerID uuid.UUID, jobPostingID uuid.UUID) ([]models.JobApplication, error)
	GetMyApplications(applicantID uuid.UUID) ([]models.JobApplication, error)
	UpdateApplicationStatus(employerID uuid.UUID, applicationID uuid.UUID, status models.ApplicationStatus) (*models.JobApplication, error)
}

type jobApplicationService struct {
	db       *gorm.DB
	appRepo  repositories.JobApplicationRepository
	jobRepo  repositories.JobPostingRepository
	userRepo repositories.UserRepository
	subRepo  repositories.SubscriptionRepository
	connRepo repositories.ConnectionRepository
}

func NewJobApplicationService(
	db *gorm.DB,
	appRepo repositories.JobApplicationRepository,
	jobRepo repositories.JobPostingRepository,
	userRepo repositories.UserRepository,
	subRepo repositories.SubscriptionRepository,
	connRepo repositories.ConnectionRepository,
) JobApplicationService {
	return &jobApplicationService{
		db: db, appRepo: appRepo, jobRepo: jobRepo, userRepo: userRepo, subRepo: subRepo, connRepo: connRepo,
	}
}

// ApplyToJob: pelamar wajib KTP terverifikasi & punya resume; lowongan harus OPEN.
func (s *jobApplicationService) ApplyToJob(applicantID uuid.UUID, jobPostingID uuid.UUID) (*models.JobApplication, error) {
	user, err := s.userRepo.FindByID(applicantID)
	if err != nil {
		return nil, apperrors.Internal(err)
	}
	if user == nil {
		return nil, apperrors.ErrUserNotFound
	}
	if !user.IsKTPVerified {
		return nil, apperrors.ErrKTPNotVerified
	}
	if user.ResumeURL == nil || *user.ResumeURL == "" {
		return nil, apperrors.ErrResumeNotUploaded
	}

	sub, err := s.subRepo.FindActiveByUserAndPlan(applicantID, models.PlanApplyJob5K)
	if err != nil {
		return nil, apperrors.Internal(err)
	}
	if sub == nil {
		return nil, apperrors.ErrSubscriptionRequired.WithDetails(map[string]any{
			"plan_type": models.PlanApplyJob5K,
			"message":   "Anda wajib memiliki Paket Melamar (Rp 5.000 / 30 hari) yang aktif untuk melamar lowongan.",
		})
	}

	job, err := s.jobRepo.FindByID(jobPostingID)
	if err != nil {
		return nil, apperrors.Internal(err)
	}
	if job == nil {
		return nil, apperrors.NotFound("Lowongan pekerjaan tidak ditemukan")
	}
	if job.EmployerID == applicantID {
		return nil, apperrors.Forbidden("SELF_APPLY", "Anda tidak dapat melamar ke lowongan Anda sendiri.")
	}
	if job.Status != models.JobStatusOpen {
		return nil, apperrors.Conflict("JOB_NOT_OPEN", "Lowongan ini sudah tidak menerima lamaran.")
	}

	existing, err := s.appRepo.FindExisting(jobPostingID, applicantID)
	if err != nil {
		return nil, apperrors.Internal(err)
	}
	if existing != nil {
		return nil, apperrors.Conflict("ALREADY_APPLIED", "Anda sudah melamar ke lowongan ini.")
	}

	app := &models.JobApplication{
		JobPostingID: jobPostingID,
		ApplicantID:  applicantID,
		Status:       models.AppStatusPending,
	}
	if err := s.appRepo.Create(app); err != nil {
		// Race dengan unique index (job_posting_id, applicant_id).
		return nil, apperrors.Conflict("ALREADY_APPLIED", "Anda sudah melamar ke lowongan ini.").WithCause(err)
	}
	created, err := s.appRepo.FindByID(app.ID)
	if err != nil {
		return nil, apperrors.Internal(err)
	}
	return created, nil
}

func (s *jobApplicationService) GetApplicantsForJob(employerID uuid.UUID, jobPostingID uuid.UUID) ([]models.JobApplication, error) {
	job, err := s.jobRepo.FindByID(jobPostingID)
	if err != nil {
		return nil, apperrors.Internal(err)
	}
	if job == nil {
		return nil, apperrors.NotFound("Lowongan pekerjaan tidak ditemukan")
	}
	if job.EmployerID != employerID {
		return nil, apperrors.Forbidden("NOT_OWNER", "Anda bukan pemilik lowongan ini.")
	}
	apps, err := s.appRepo.FindByJobPostingID(jobPostingID)
	if err != nil {
		return nil, apperrors.Internal(err)
	}
	return apps, nil
}

func (s *jobApplicationService) GetMyApplications(applicantID uuid.UUID) ([]models.JobApplication, error) {
	apps, err := s.appRepo.FindByApplicantID(applicantID)
	if err != nil {
		return nil, apperrors.Internal(err)
	}
	return apps, nil
}

// UpdateApplicationStatus dijalankan dalam satu transaksi:
// ACCEPTED -> lamaran ACCEPTED, lamaran PENDING lain REJECTED, lowongan IN_PROGRESS, koneksi dibuat.
func (s *jobApplicationService) UpdateApplicationStatus(employerID uuid.UUID, applicationID uuid.UUID, status models.ApplicationStatus) (*models.JobApplication, error) {
	if status != models.AppStatusAccepted && status != models.AppStatusRejected {
		return nil, apperrors.Validation("Status tidak valid.", map[string]any{"status": "ACCEPTED atau REJECTED"})
	}

	err := s.db.Transaction(func(tx *gorm.DB) error {
		appRepo := s.appRepo.WithTx(tx)
		jobRepo := s.jobRepo.WithTx(tx)
		connRepo := s.connRepo.WithTx(tx)

		app, err := appRepo.FindByIDForUpdate(applicationID)
		if err != nil {
			return apperrors.Internal(err)
		}
		if app == nil {
			return apperrors.NotFound("Lamaran tidak ditemukan")
		}
		job, err := jobRepo.FindByID(app.JobPostingID)
		if err != nil {
			return apperrors.Internal(err)
		}
		if job == nil {
			return apperrors.NotFound("Lowongan pekerjaan tidak ditemukan")
		}
		if job.EmployerID != employerID {
			return apperrors.Forbidden("NOT_OWNER", "Anda bukan pemilik lowongan ini.")
		}
		if app.Status != models.AppStatusPending {
			return apperrors.Conflict("APPLICATION_ALREADY_RESPONDED", "Lamaran ini sudah diproses sebelumnya.")
		}
		if status == models.AppStatusAccepted && job.Status != models.JobStatusOpen {
			return apperrors.Conflict("JOB_NOT_OPEN", "Lowongan sudah tidak terbuka; tidak dapat menerima pelamar baru.")
		}

		if err := appRepo.UpdateStatus(applicationID, status); err != nil {
			return apperrors.Internal(err)
		}
		if status != models.AppStatusAccepted {
			return nil
		}

		if _, err := appRepo.RejectOtherPending(job.ID, applicationID); err != nil {
			return apperrors.Internal(err)
		}
		if err := jobRepo.UpdateStatus(job.ID, models.JobStatusInProgress, nil); err != nil {
			return apperrors.Internal(err)
		}
		conn := &models.Connection{
			SourceType: models.SourceJobApplication,
			SourceID:   app.ID,
			EmployerID: job.EmployerID,
			WorkerID:   app.ApplicantID,
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

	updated, err := s.appRepo.FindByID(applicationID)
	if err != nil {
		return nil, apperrors.Internal(err)
	}
	return updated, nil
}
