package models

import (
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserMode string

const (
	ModeSeeker   UserMode = "SEEKER"
	ModeEmployer UserMode = "EMPLOYER"
)

func (m UserMode) Valid() bool { return m == ModeSeeker || m == ModeEmployer }

type KTPStatus string

const (
	KTPUnverified    KTPStatus = "UNVERIFIED"
	KTPPendingReview KTPStatus = "PENDING_REVIEW"
	KTPVerified      KTPStatus = "VERIFIED"
	KTPRejected      KTPStatus = "REJECTED"
)

// User adalah entitas pengguna. Serialisasi JSON struct ini hanya boleh dikirim
// kepada PEMILIK akun (endpoint /users/me dan respons auth). Untuk ditampilkan ke
// pengguna lain gunakan User.Public() / PublicUser agar data pribadi tidak bocor.
type User struct {
	ID            uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Email         string    `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
	Phone         string    `gorm:"type:varchar(50);uniqueIndex;not null" json:"phone"`
	Name          string    `gorm:"type:varchar(255);not null" json:"name"`
	Password      string    `gorm:"type:varchar(255);column:password;not null" json:"-"`
	ActiveMode    UserMode  `gorm:"type:varchar(20);default:'SEEKER';not null" json:"active_mode"`
	IsKTPVerified bool      `gorm:"not null;default:false" json:"is_ktp_verified"`

	KTPStatus          KTPStatus  `gorm:"type:varchar(30);not null;default:'UNVERIFIED'" json:"ktp_status"`
	KTPNumber          *string    `gorm:"type:varchar(50)" json:"ktp_number,omitempty"`
	KTPImageURL        *string    `gorm:"type:varchar(500)" json:"ktp_image_url,omitempty"`
	KTPSubmittedAt     *time.Time `json:"ktp_submitted_at,omitempty"`
	KTPVerifiedAt      *time.Time `json:"ktp_verified_at,omitempty"`
	KTPRejectionReason *string    `gorm:"type:text" json:"ktp_rejection_reason,omitempty"`

	ResumeURL *string `gorm:"type:varchar(500)" json:"resume_url,omitempty"`

	Country          string   `gorm:"type:varchar(100);default:'Indonesia'" json:"country"`
	Province         string   `gorm:"type:varchar(100);not null" json:"province"`
	City             string   `gorm:"type:varchar(100);not null" json:"city"`
	District         string   `gorm:"type:varchar(100);not null" json:"district"`
	AddressDetail    *string  `gorm:"type:text" json:"address_detail,omitempty"`
	Latitude         *float64 `gorm:"type:decimal(10,8)" json:"latitude,omitempty"`
	Longitude        *float64 `gorm:"type:decimal(11,8)" json:"longitude,omitempty"`
	FormattedAddress *string  `gorm:"type:text" json:"formatted_address,omitempty"`
	OSMPlaceID       *string  `gorm:"type:varchar(255)" json:"osm_place_id,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// PublicUser adalah proyeksi User yang aman ditampilkan kepada pengguna lain
// (pencarian talent/lowongan, daftar pelamar, offer, rating). Tidak memuat email,
// nomor HP, KTP, alamat detail, maupun koordinat.
type PublicUser struct {
	ID            uuid.UUID `json:"id"`
	Name          string    `json:"name"`
	ActiveMode    UserMode  `json:"active_mode"`
	IsKTPVerified bool      `json:"is_ktp_verified"`
	Province      string    `json:"province"`
	City          string    `json:"city"`
	District      string    `json:"district"`
	CreatedAt     time.Time `json:"created_at"`
}

func (u *User) Public() *PublicUser {
	if u == nil {
		return nil
	}
	return &PublicUser{
		ID:            u.ID,
		Name:          u.Name,
		ActiveMode:    u.ActiveMode,
		IsKTPVerified: u.IsKTPVerified,
		Province:      u.Province,
		City:          u.City,
		District:      u.District,
		CreatedAt:     u.CreatedAt,
	}
}

// ApplicantView adalah proyeksi pelamar untuk pemberi kerja: profil publik + resume
// (PRD: employer meninjau resume & status KTP sebelum menerima). Kontak tetap disembunyikan
// sampai status TERHUBUNG.
type ApplicantView struct {
	PublicUser
	ResumeURL *string `json:"resume_url,omitempty"`
}

func (u *User) Applicant() *ApplicantView {
	if u == nil {
		return nil
	}
	return &ApplicantView{PublicUser: *u.Public(), ResumeURL: u.ResumeURL}
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	if u.ActiveMode == "" {
		u.ActiveMode = ModeSeeker
	}
	if u.Country == "" {
		u.Country = "Indonesia"
	}
	if u.KTPStatus == "" {
		u.KTPStatus = KTPUnverified
	}
	if u.Password != "" && !isBcryptHash(u.Password) {
		hashed, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		u.Password = string(hashed)
	}
	return nil
}

func isBcryptHash(s string) bool {
	return strings.HasPrefix(s, "$2a$") || strings.HasPrefix(s, "$2b$") || strings.HasPrefix(s, "$2y$")
}
