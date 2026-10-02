package middlewares

import (
	"crypto/subtle"
	"errors"
	"strings"

	"layanesia-server/internal/apperrors"
	"layanesia-server/internal/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// Kunci c.Locals. Selalu gunakan konstanta ini (dan helper di bawah), jangan
// string literal, agar tidak terulang bug key/tipe berbeda antar controller.
const (
	localsUserID     = "auth.user_id"
	localsUserEmail  = "auth.user_email"
	localsActiveMode = "auth.active_mode"
)

// Protected memverifikasi header Authorization: Bearer <jwt> dan menyimpan
// identitas ke Locals.
func Protected() fiber.Handler {
	return func(c *fiber.Ctx) error {
		token, ok := bearerToken(c.Get(fiber.HeaderAuthorization))
		if !ok {
			return apperrors.Unauthorized("Header Authorization tidak ada atau formatnya salah. Gunakan 'Bearer <token>'.")
		}

		claims, err := utils.ValidateJWTToken(token)
		if err != nil {
			if errors.Is(err, utils.ErrTokenExpired) {
				return apperrors.New(fiber.StatusUnauthorized, "TOKEN_EXPIRED", "Sesi Anda telah berakhir. Silakan masuk kembali.")
			}
			return apperrors.Unauthorized("Token tidak valid.")
		}

		c.Locals(localsUserID, claims.UserID)
		c.Locals(localsUserEmail, claims.Email)
		c.Locals(localsActiveMode, claims.ActiveMode)
		return c.Next()
	}
}

// GetUserID mengambil ID pengguna yang sudah diautentikasi oleh Protected().
func GetUserID(c *fiber.Ctx) (uuid.UUID, bool) {
	id, ok := c.Locals(localsUserID).(uuid.UUID)
	return id, ok && id != uuid.Nil
}

// MustUserID seperti GetUserID tetapi mengembalikan AppError siap-return.
func MustUserID(c *fiber.Ctx) (uuid.UUID, error) {
	id, ok := GetUserID(c)
	if !ok {
		return uuid.Nil, apperrors.ErrSessionInvalid
	}
	return id, nil
}

// AdminOnly melindungi endpoint operasional (misal approve KTP) dengan header
// X-Admin-Key yang harus sama dengan ADMIN_API_KEY.
func AdminOnly(adminKey string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if adminKey == "" {
			return apperrors.Forbidden("ADMIN_DISABLED", "Endpoint admin tidak aktif: ADMIN_API_KEY belum dikonfigurasi.")
		}
		provided := c.Get("X-Admin-Key")
		if provided == "" || subtle.ConstantTimeCompare([]byte(provided), []byte(adminKey)) != 1 {
			return apperrors.Forbidden("ADMIN_KEY_INVALID", "Kunci admin tidak valid.")
		}
		return c.Next()
	}
}

func bearerToken(header string) (string, bool) {
	if header == "" {
		return "", false
	}
	parts := strings.Fields(header)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") || parts[1] == "" {
		return "", false
	}
	return parts[1], true
}
