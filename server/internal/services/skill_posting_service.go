package services

import (
	"strings"

	"layanesia-server/internal/apperrors"
	"layanesia-server/internal/models"
	"layanesia-server/internal/repositories"
	"layanesia-server/internal/utils"

	"github.com/google/uuid"
)

type CreateSkillRequest struct {
	Category         models.SkillCategory      `json:"category" validate:"omitempty,oneof=SERABUTAN PROFESIONAL"`
	Title            string                    `json:"title" validate:"required,min=3,max=255"`
	Description      string                    `json:"description" validate:"required,min=10,max=5000"`
	Country          string                    `json:"country" validate:"omitempty,max=100"`
	Province         string                    `json:"province" validate:"omitempty,max=100"`
	City             string                    `json:"city" validate:"omitempty,max=100"`
	District         string                    `json:"district" validate:"omitempty,max=100"`
	AddressDetail    *string                   `json:"address_detail" validate:"omitempty,max=1000"`
	Latitude         *float64                  `json:"latitude" validate:"omitempty,gte=-90,lte=90"`
	Longitude        *float64                  `json:"longitude" validate:"omitempty,gte=-180,lte=180"`
	FormattedAddress *string                   `json:"formatted_address" validate:"omitempty,max=1000"`
	OSMPlaceID       *string                   `json:"osm_place_id" validate:"omitempty,max=255"`
	RateType         models.RateType           `json:"rate_type" validate:"omitempty,oneof=PER_HOUR PER_DAY"`
	RateAmount       float64                   `json:"rate_amount" validate:"required,gt=0,lte=999999999"`
	Availability     models.AvailabilityStatus `json:"availability" validate:"omitempty,oneof=AVAILABLE BUSY"`
}

type UpdateSkillAvailabilityRequest struct {
	Availability models.AvailabilityStatus `json:"availability" validate:"required,oneof=AVAILABLE BUSY"`
}

type SkillPostingService interface {
	CreateSkillPosting(userID uuid.UUID, req CreateSkillRequest) (*models.SkillPosting, error)
	GetMySkillPostings(userID uuid.UUID) ([]models.SkillPosting, error)
	SearchSkillPostings(filter repositories.SkillFilter) ([]models.SkillPosting, int64, error)
	GetSkillByID(id uuid.UUID) (*models.SkillPosting, error)
	UpdateAvailability(userID uuid.UUID, id uuid.UUID, availability models.AvailabilityStatus) (*models.SkillPosting, error)
	DeleteSkillPosting(userID uuid.UUID, id uuid.UUID) error
}

type skillPostingService struct {
	skillRepo repositories.SkillPostingRepository
	userRepo  repositories.UserRepository
	subRepo   repositories.SubscriptionRepository
}

func NewSkillPostingService(
	skillRepo repositories.SkillPostingRepository,
	userRepo repositories.UserRepository,
	subRepo repositories.SubscriptionRepository,
) SkillPostingService {
	return &skillPostingService{skillRepo: skillRepo, userRepo: userRepo, subRepo: subRepo}
}

func (s *skillPostingService) CreateSkillPosting(userID uuid.UUID, req CreateSkillRequest) (*models.SkillPosting, error) {
	req.Title = strings.TrimSpace(req.Title)
	req.Description = strings.TrimSpace(req.Description)
	if err := utils.ValidateStruct(req); err != nil {
		return nil, err
	}

	// Paket POST_SKILL_10K wajib aktif.
	sub, err := s.subRepo.FindActiveByUserAndPlan(userID, models.PlanPostSkill10K)
	if err != nil {
		return nil, apperrors.Internal(err)
	}
	if sub == nil {
		return nil, apperrors.ErrSubscriptionRequired.WithDetails(map[string]any{
			"plan_type": models.PlanPostSkill10K,
			"message":   "Anda wajib memiliki Paket Etalase Keahlian (Rp 10.000 / 30 hari) yang aktif untuk memasang profil keahlian.",
		})
	}

	user, err := s.userRepo.FindByID(userID)
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

	posting := &models.SkillPosting{
		UserID:           userID,
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
		RateType:         orDefault(req.RateType, models.RatePerHour),
		RateAmount:       req.RateAmount,
		Availability:     orDefault(req.Availability, models.AvailAvailable),
	}
	if err := s.skillRepo.Create(posting); err != nil {
		return nil, apperrors.Internal(err)
	}
	created, err := s.skillRepo.FindByID(posting.ID)
	if err != nil {
		return nil, apperrors.Internal(err)
	}
	return created, nil
}

func (s *skillPostingService) GetMySkillPostings(userID uuid.UUID) ([]models.SkillPosting, error) {
	items, err := s.skillRepo.FindByUserID(userID)
	if err != nil {
		return nil, apperrors.Internal(err)
	}
	return items, nil
}

func (s *skillPostingService) SearchSkillPostings(filter repositories.SkillFilter) ([]models.SkillPosting, int64, error) {
	if filter.Category != "" && !models.SkillCategory(filter.Category).Valid() {
		return nil, 0, apperrors.Validation("Kategori tidak valid.", map[string]any{"category": "SERABUTAN atau PROFESIONAL"})
	}
	if filter.Availability != "" && !models.AvailabilityStatus(filter.Availability).Valid() {
		return nil, 0, apperrors.Validation("Availability tidak valid.", map[string]any{"availability": "AVAILABLE atau BUSY"})
	}
	items, total, err := s.skillRepo.FindAll(filter)
	if err != nil {
		return nil, 0, apperrors.Internal(err)
	}
	return items, total, nil
}

func (s *skillPostingService) GetSkillByID(id uuid.UUID) (*models.SkillPosting, error) {
	item, err := s.skillRepo.FindByID(id)
	if err != nil {
		return nil, apperrors.Internal(err)
	}
	if item == nil {
		return nil, apperrors.NotFound("Postingan keahlian tidak ditemukan")
	}
	return item, nil
}

func (s *skillPostingService) UpdateAvailability(userID uuid.UUID, id uuid.UUID, availability models.AvailabilityStatus) (*models.SkillPosting, error) {
	if !availability.Valid() {
		return nil, apperrors.Validation("Availability tidak valid.", map[string]any{"availability": "AVAILABLE atau BUSY"})
	}
	item, err := s.GetSkillByID(id)
	if err != nil {
		return nil, err
	}
	if item.UserID != userID {
		return nil, apperrors.Forbidden("NOT_OWNER", "Anda bukan pemilik postingan keahlian ini.")
	}
	item.Availability = availability
	if err := s.skillRepo.Update(item); err != nil {
		return nil, apperrors.Internal(err)
	}
	return item, nil
}

func (s *skillPostingService) DeleteSkillPosting(userID uuid.UUID, id uuid.UUID) error {
	n, err := s.skillRepo.Delete(id, userID)
	if err != nil {
		return apperrors.Internal(err)
	}
	if n == 0 {
		return apperrors.NotFound("Postingan keahlian tidak ditemukan atau bukan milik Anda")
	}
	return nil
}

func orDefault[T ~string](v, def T) T {
	if v == "" {
		return def
	}
	return v
}
