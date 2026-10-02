package services

import (
	"strings"
	"time"

	"layanesia-server/internal/apperrors"
	"layanesia-server/internal/models"
	"layanesia-server/internal/repositories"
	"layanesia-server/internal/utils"

	"github.com/google/uuid"
)

type UpdateProfileDTO struct {
	Name             *string  `json:"name" validate:"omitempty,min=2,max=255"`
	Phone            *string  `json:"phone" validate:"omitempty,id_phone"`
	Province         *string  `json:"province" validate:"omitempty,max=100"`
	City             *string  `json:"city" validate:"omitempty,max=100"`
	District         *string  `json:"district" validate:"omitempty,max=100"`
	AddressDetail    *string  `json:"address_detail" validate:"omitempty,max=1000"`
	Latitude         *float64 `json:"latitude" validate:"omitempty,gte=-90,lte=90"`
	Longitude        *float64 `json:"longitude" validate:"omitempty,gte=-180,lte=180"`
	FormattedAddress *string  `json:"formatted_address" validate:"omitempty,max=1000"`
	OSMPlaceID       *string  `json:"osm_place_id" validate:"omitempty,max=255"`
}

type VerifyKTPDTO struct {
	KTPNumber   string `json:"ktp_number" validate:"required,nik"`
	KTPImageURL string `json:"ktp_image_url" validate:"required,max=500"`
}

type UploadResumeDTO struct {
	ResumeURL string `json:"resume_url" validate:"required,max=500"`
}

type SwitchModeDTO struct {
	Mode models.UserMode `json:"mode" validate:"required,oneof=SEEKER EMPLOYER"`
}

type ReviewKTPDTO struct {
	Approve bool    `json:"approve"`
	Reason  *string `json:"reason" validate:"omitempty,max=500"`
}

type UserService interface {
	GetProfile(userID uuid.UUID) (*models.User, error)
	UpdateProfile(userID uuid.UUID, req UpdateProfileDTO) (*models.User, error)
	SwitchMode(userID uuid.UUID, mode models.UserMode) (*models.User, string, error)
	SubmitKTP(userID uuid.UUID, req VerifyKTPDTO) (*models.User, error)
	UploadResume(userID uuid.UUID, req UploadResumeDTO) (*models.User, error)
	ReviewKTP(userID uuid.UUID, req ReviewKTPDTO) (*models.User, error)
	ListPendingKTP(limit, offset int) ([]models.User, error)
}

type userService struct {
	userRepo      repositories.UserRepository
	ktpAutoVerify bool
}

func NewUserService(userRepo repositories.UserRepository, ktpAutoVerify bool) UserService {
	return &userService{userRepo: userRepo, ktpAutoVerify: ktpAutoVerify}
}

func (s *userService) mustGet(userID uuid.UUID) (*models.User, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, apperrors.Internal(err)
	}
	if user == nil {
		return nil, apperrors.ErrUserNotFound
	}
	return user, nil
}

func (s *userService) GetProfile(userID uuid.UUID) (*models.User, error) {
	return s.mustGet(userID)
}

// blankToNil memperlakukan string kosong sebagai "tidak diubah" (kompatibel dengan
// klien yang mengirim seluruh form).
func blankToNil(p *string) *string {
	if p == nil || strings.TrimSpace(*p) == "" {
		return nil
	}
	return p
}

func (s *userService) UpdateProfile(userID uuid.UUID, req UpdateProfileDTO) (*models.User, error) {
	req.Name = blankToNil(req.Name)
	req.Phone = blankToNil(req.Phone)
	req.Province = blankToNil(req.Province)
	req.City = blankToNil(req.City)
	req.District = blankToNil(req.District)
	req.OSMPlaceID = blankToNil(req.OSMPlaceID)
	if err := utils.ValidateStruct(req); err != nil {
		return nil, err
	}
	user, err := s.mustGet(userID)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		user.Name = strings.TrimSpace(*req.Name)
	}
	if req.Phone != nil {
		phone := utils.NormalizePhone(*req.Phone)
		if phone != user.Phone {
			if other, err := s.userRepo.FindByPhone(phone); err != nil {
				return nil, apperrors.Internal(err)
			} else if other != nil && other.ID != user.ID {
				return nil, apperrors.Conflict("PHONE_TAKEN", "Nomor telepon sudah dipakai akun lain.")
			}
			user.Phone = phone
		}
	}
	if req.Province != nil {
		user.Province = strings.TrimSpace(*req.Province)
	}
	if req.City != nil {
		user.City = strings.TrimSpace(*req.City)
	}
	if req.District != nil {
		user.District = strings.TrimSpace(*req.District)
	}
	if req.AddressDetail != nil {
		user.AddressDetail = req.AddressDetail
	}
	if req.Latitude != nil {
		user.Latitude = req.Latitude
	}
	if req.Longitude != nil {
		user.Longitude = req.Longitude
	}
	if req.FormattedAddress != nil {
		user.FormattedAddress = req.FormattedAddress
	}
	if req.OSMPlaceID != nil {
		user.OSMPlaceID = req.OSMPlaceID
	}
	if user.Province == "" || user.City == "" || user.District == "" {
		return nil, apperrors.Validation("Provinsi, kota, dan kecamatan tidak boleh kosong.", nil)
	}

	if err := s.userRepo.Update(user); err != nil {
		return nil, apperrors.Internal(err)
	}
	return user, nil
}

func (s *userService) SwitchMode(userID uuid.UUID, mode models.UserMode) (*models.User, string, error) {
	if !mode.Valid() {
		return nil, "", apperrors.Validation("Mode tidak valid.", map[string]any{"mode": "harus SEEKER atau EMPLOYER"})
	}
	user, err := s.userRepo.UpdateActiveMode(userID, mode)
	if err != nil {
		return nil, "", apperrors.Internal(err)
	}
	if user == nil {
		return nil, "", apperrors.ErrUserNotFound
	}
	token, err := utils.GenerateJWTToken(user.ID, user.Email, string(user.ActiveMode))
	if err != nil {
		return nil, "", apperrors.Internal(err)
	}
	return user, token, nil
}

// SubmitKTP menerima pengajuan verifikasi KTP. Status menjadi PENDING_REVIEW sampai
// disetujui admin (ReviewKTP). Di dev, KTP_AUTO_VERIFY=true langsung VERIFIED.
func (s *userService) SubmitKTP(userID uuid.UUID, req VerifyKTPDTO) (*models.User, error) {
	req.KTPNumber = strings.TrimSpace(req.KTPNumber)
	req.KTPImageURL = strings.TrimSpace(req.KTPImageURL)
	if err := utils.ValidateStruct(req); err != nil {
		return nil, err
	}
	user, err := s.mustGet(userID)
	if err != nil {
		return nil, err
	}
	if user.KTPStatus == models.KTPVerified {
		return nil, apperrors.Conflict("KTP_ALREADY_VERIFIED", "KTP Anda sudah terverifikasi.")
	}

	now := time.Now()
	user.KTPNumber = &req.KTPNumber
	user.KTPImageURL = &req.KTPImageURL
	user.KTPSubmittedAt = &now
	user.KTPRejectionReason = nil

	if s.ktpAutoVerify {
		user.KTPStatus = models.KTPVerified
		user.IsKTPVerified = true
		user.KTPVerifiedAt = &now
	} else {
		user.KTPStatus = models.KTPPendingReview
		user.IsKTPVerified = false
		user.KTPVerifiedAt = nil
	}

	if err := s.userRepo.Update(user); err != nil {
		return nil, apperrors.Internal(err)
	}
	return user, nil
}

// ReviewKTP dipakai admin untuk menyetujui/menolak pengajuan KTP.
func (s *userService) ReviewKTP(userID uuid.UUID, req ReviewKTPDTO) (*models.User, error) {
	if err := utils.ValidateStruct(req); err != nil {
		return nil, err
	}
	user, err := s.mustGet(userID)
	if err != nil {
		return nil, err
	}
	if user.KTPStatus != models.KTPPendingReview {
		return nil, apperrors.Conflict("KTP_NOT_PENDING", "Pengguna ini tidak memiliki pengajuan KTP yang menunggu review.")
	}

	now := time.Now()
	if req.Approve {
		user.KTPStatus = models.KTPVerified
		user.IsKTPVerified = true
		user.KTPVerifiedAt = &now
		user.KTPRejectionReason = nil
	} else {
		user.KTPStatus = models.KTPRejected
		user.IsKTPVerified = false
		user.KTPVerifiedAt = nil
		user.KTPRejectionReason = req.Reason
	}
	if err := s.userRepo.Update(user); err != nil {
		return nil, apperrors.Internal(err)
	}
	return user, nil
}

func (s *userService) ListPendingKTP(limit, offset int) ([]models.User, error) {
	users, err := s.userRepo.FindByKTPStatus(models.KTPPendingReview, limit, offset)
	if err != nil {
		return nil, apperrors.Internal(err)
	}
	return users, nil
}

func (s *userService) UploadResume(userID uuid.UUID, req UploadResumeDTO) (*models.User, error) {
	req.ResumeURL = strings.TrimSpace(req.ResumeURL)
	if err := utils.ValidateStruct(req); err != nil {
		return nil, err
	}
	user, err := s.mustGet(userID)
	if err != nil {
		return nil, err
	}
	user.ResumeURL = &req.ResumeURL
	if err := s.userRepo.Update(user); err != nil {
		return nil, apperrors.Internal(err)
	}
	return user, nil
}
