package services

import (
	"strings"

	"layanesia-server/internal/apperrors"
	"layanesia-server/internal/models"
	"layanesia-server/internal/repositories"
	"layanesia-server/internal/utils"

	"golang.org/x/crypto/bcrypt"
)

type RegisterDTO struct {
	Name          string          `json:"name" validate:"required,min=2,max=255"`
	Email         string          `json:"email" validate:"required,email,max=255"`
	Phone         string          `json:"phone" validate:"required,id_phone"`
	Password      string          `json:"password" validate:"required,min=8,max=72"`
	Province      string          `json:"province" validate:"required,max=100"`
	City          string          `json:"city" validate:"required,max=100"`
	District      string          `json:"district" validate:"required,max=100"`
	AddressDetail *string         `json:"address_detail" validate:"omitempty,max=1000"`
	ActiveMode    models.UserMode `json:"active_mode" validate:"omitempty,oneof=SEEKER EMPLOYER"`
}

type LoginDTO struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type AuthResponse struct {
	Token string       `json:"token"`
	User  *models.User `json:"user"`
}

type AuthService interface {
	Register(req RegisterDTO) (*AuthResponse, error)
	Login(req LoginDTO) (*AuthResponse, error)
}

type authService struct {
	userRepo repositories.UserRepository
}

func NewAuthService(userRepo repositories.UserRepository) AuthService {
	return &authService{userRepo: userRepo}
}

var errBadCredentials = apperrors.Unauthorized("Email atau password tidak sesuai.")

func (s *authService) Register(req RegisterDTO) (*AuthResponse, error) {
	req.Name = strings.TrimSpace(req.Name)
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))
	req.Phone = utils.NormalizePhone(req.Phone)
	req.Province = strings.TrimSpace(req.Province)
	req.City = strings.TrimSpace(req.City)
	req.District = strings.TrimSpace(req.District)

	if err := utils.ValidateStruct(req); err != nil {
		return nil, err
	}

	if existing, err := s.userRepo.FindByEmail(req.Email); err != nil {
		return nil, apperrors.Internal(err)
	} else if existing != nil {
		return nil, apperrors.Conflict("EMAIL_TAKEN", "Email sudah terdaftar.")
	}
	if existing, err := s.userRepo.FindByPhone(req.Phone); err != nil {
		return nil, apperrors.Internal(err)
	} else if existing != nil {
		return nil, apperrors.Conflict("PHONE_TAKEN", "Nomor telepon sudah terdaftar.")
	}

	activeMode := req.ActiveMode
	if activeMode == "" {
		activeMode = models.ModeSeeker
	}

	user := &models.User{
		Name:          req.Name,
		Email:         req.Email,
		Phone:         req.Phone,
		Password:      req.Password, // di-hash oleh hook BeforeCreate
		ActiveMode:    activeMode,
		Country:       "Indonesia",
		Province:      req.Province,
		City:          req.City,
		District:      req.District,
		AddressDetail: req.AddressDetail,
		KTPStatus:     models.KTPUnverified,
	}
	if err := s.userRepo.Create(user); err != nil {
		return nil, apperrors.Internal(err)
	}

	token, err := utils.GenerateJWTToken(user.ID, user.Email, string(user.ActiveMode))
	if err != nil {
		return nil, apperrors.Internal(err)
	}
	return &AuthResponse{Token: token, User: user}, nil
}

func (s *authService) Login(req LoginDTO) (*AuthResponse, error) {
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))
	if err := utils.ValidateStruct(req); err != nil {
		return nil, err
	}

	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		return nil, apperrors.Internal(err)
	}
	if user == nil {
		// Tetap lakukan perbandingan bcrypt agar waktu respons serupa (anti user-enumeration).
		_ = bcrypt.CompareHashAndPassword([]byte("$2a$10$invalidsaltinvalidsaltinvalidsaltinvalidsaltinvalidsalt"), []byte(req.Password))
		return nil, errBadCredentials
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errBadCredentials
	}

	token, err := utils.GenerateJWTToken(user.ID, user.Email, string(user.ActiveMode))
	if err != nil {
		return nil, apperrors.Internal(err)
	}
	return &AuthResponse{Token: token, User: user}, nil
}
