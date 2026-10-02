package services

import (
	"strings"
	"time"

	"layanesia-server/internal/apperrors"
	"layanesia-server/internal/models"
	"layanesia-server/internal/repositories"
	"layanesia-server/internal/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CreateJobRequest struct {
	Category         models.SkillCategory `json:"category" validate:"omitempty,oneof=SERABUTAN PROFESIONAL"`
	Title            string               `json:"title" validate:"required,min=3,max=255"`
	Description      string               `json:"description" validate:"required,min=10,max=5000"`
	Country          string               `json:"country" validate:"omitempty,max=100"`
	Province         string               `json:"province" validate:"omitempty,max=100"`
	City             string               `json:"city" validate:"omitempty,max=100"`
	District         string               `json:"district" validate:"omitempty,max=100"`
	AddressDetail    *string              `json:"address_detail" validate:"omitempty,max=1000"`
	Latitude         *float64             `json:"latitude" validate:"omitempty,gte=-90,lte=90"`
	Longitude        *float64             `json:"longitude" validate:"omitempty,gte=-180,lte=180"`
	FormattedAddress *string              `json:"formatted_address" validate:"omitempty,max=1000"`
	OSMPlaceID       *string              `json:"osm_place_id" validate:"omitempty,max=255"`
	WorkDate         string               `json:"work_date" validate:"omitempty"`
	DurationType     models.DurationType  `json:"duration_type" validate:"omitempty,oneof=HOURS DAYS"`
	DurationValue    int                  `json:"duration_value" validate:"omitempty,gte=1,lte=365"`
	Budget           float64              `json:"budget" validate:"required,gt=0,lte=999999999"`
}

type JobPostingService interface {
	CreateJobPosting(employerID uuid.UUID, req CreateJobRequest) (*models.JobPosting, error)
	GetMyJobPostings(employerID uuid.UUID) ([]models.JobPosting, error)
	SearchJobPostings(filter repositories.JobFilter) ([]models.JobPosting, int64, error)
	GetJobByID(id uuid.UUID) (*models.JobPosting, error)
	CloseJob(employerID uuid.UUID, id uuid.UUID, status models.JobStatus) (*models.JobPosting, error)
}

type jobPostingService struct {
	db       *gorm.DB
	jobRepo  repositories.JobPostingRepository
	appRepo  repositories.JobApplicationRepository
	userRepo repositories.UserRepository
}

func NewJobPostingService(
	db *gorm.DB,
	jobRepo repositories.JobPostingRepository,
	appRepo repositories.JobApplicationRepository,
	userRepo repositories.UserRepository,
) JobPostingService {
	return &jobPostingService{db: db, jobRepo: jobRepo, appRepo: appRepo, userRepo: userRepo}
}

func (s *jobPostingService) CreateJobPosting(employerID uuid.UUID, req CreateJobRequest) (*models.JobPosting, error) {
	req.Title = strings.TrimSpace(req.Title)
	req.Description = strings.TrimSpace(req.Description)
	if err := utils.ValidateStruct(req); err != nil {
		return nil, err
	}

	user, err := s.userRepo.FindByID(employerID)
	if err != nil {
		return nil, apperrors.Internal(err)
	}
	if user == nil {
		return nil, apperrors.ErrUserNotFound
	}

	loc, err := resolveLocation(locationInput{
		Country: req.Country, Province: req.Province, City: req.City, District: req.District,
		AddressDetail: req.AddressDetail, FormattedAddress: req.FormattedAddress, OSMPlaceID: req.OSMPlaceID,
		Latitude: req.Latitude, Longitude: req.Longitude,
	}, user)
	if err != nil {
		return nil, err
	}

	workDate, err := parseWorkDate(req.WorkDate, time.Now())
	if err != nil {
		return nil, err
	}

	durationValue := req.DurationValue
	if durationValue <= 0 {
		durationValue = 1
	}

	job := &models.JobPosting{
		EmployerID:       employerID,
		Category:         orDefault(req.Category, models.CategorySerabutan),
		Title:            req.Title,
		Description:      req.Description,
		Country:          loc.Country,
		Province:         loc.Province,
		City:             loc.City,
		District:         loc.District,
		AddressDetail:    loc.AddressDetail,
		Latitude:         loc.Latitude,
		Longitude:        loc.Longitude,
		FormattedAddress: loc.FormattedAddress,
		OSMPlaceID:       loc.OSMPlaceID,
		WorkDate:         workDate,
		DurationType:     orDefault(req.DurationType, models.DurationHours),
		DurationValue:    durationValue,
		Budget:           req.Budget,
		Status:           models.JobStatusOpen,
	}
	if err := s.jobRepo.Create(job); err != nil {
		return nil, apperrors.Internal(err)
	}
	created, err := s.jobRepo.FindByID(job.ID)
	if err != nil {
		return nil, apperrors.Internal(err)
	}
	return created, nil
}

func (s *jobPostingService) GetMyJobPostings(employerID uuid.UUID) ([]models.JobPosting, error) {
	items, err := s.jobRepo.FindByEmployerID(employerID)
	if err != nil {
		return nil, apperrors.Internal(err)
	}
	return items, nil
}

func (s *jobPostingService) SearchJobPostings(filter repositories.JobFilter) ([]models.JobPosting, int64, error) {
	if filter.Category != "" && !models.SkillCategory(filter.Category).Valid() {
		return nil, 0, apperrors.Validation("Kategori tidak valid.", map[string]any{"category": "SERABUTAN atau PROFESIONAL"})
	}
	if filter.Status != "" && !models.JobStatus(filter.Status).Valid() {
		return nil, 0, apperrors.Validation("Status tidak valid.", map[string]any{"status": "OPEN, IN_PROGRESS, DONE, atau CANCELLED"})
	}
	items, total, err := s.jobRepo.FindAll(filter)
	if err != nil {
		return nil, 0, apperrors.Internal(err)
	}
	return items, total, nil
}

func (s *jobPostingService) GetJobByID(id uuid.UUID) (*models.JobPosting, error) {
	item, err := s.jobRepo.FindByID(id)
	if err != nil {
		return nil, apperrors.Internal(err)
	}
	if item == nil {
		return nil, apperrors.NotFound("Lowongan pekerjaan tidak ditemukan")
	}
	return item, nil
}

// CloseJob menutup lowongan (DONE atau CANCELLED). Lamaran PENDING otomatis ditolak.
func (s *jobPostingService) CloseJob(employerID uuid.UUID, id uuid.UUID, status models.JobStatus) (*models.JobPosting, error) {
	if status != models.JobStatusDone && status != models.JobStatusCancelled {
		return nil, apperrors.Validation("Status penutupan tidak valid.", map[string]any{"status": "DONE atau CANCELLED"})
	}
	job, err := s.GetJobByID(id)
	if err != nil {
		return nil, err
	}
	if job.EmployerID != employerID {
		return nil, apperrors.Forbidden("NOT_OWNER", "Anda bukan pemilik lowongan ini.")
	}
	if job.Status == models.JobStatusDone || job.Status == models.JobStatusCancelled {
		return nil, apperrors.Conflict("JOB_ALREADY_CLOSED", "Lowongan ini sudah ditutup.")
	}

	now := time.Now()
	err = s.db.Transaction(func(tx *gorm.DB) error {
		if err := s.jobRepo.WithTx(tx).UpdateStatus(id, status, &now); err != nil {
			return err
		}
		_, err := s.appRepo.WithTx(tx).RejectOtherPending(id, uuid.Nil)
		return err
	})
	if err != nil {
		return nil, apperrors.Internal(err)
	}
	return s.GetJobByID(id)
}
