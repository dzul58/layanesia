package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"layanesia-server/config"

	"github.com/google/uuid"
)

// JWTClaims adalah payload token akses. Token ditandatangani HMAC-SHA256 dengan
// JWT_SECRET; verifikasi SELALU memakai HMAC (header `alg` diabaikan) sehingga
// serangan alg=none tidak berlaku.
type JWTClaims struct {
	UserID     uuid.UUID `json:"user_id"`
	Email      string    `json:"email"`
	ActiveMode string    `json:"active_mode"`
	IssuedAt   int64     `json:"iat"`
	Exp        int64     `json:"exp"`
	TokenID    string    `json:"jti"`
}

var (
	ErrTokenInvalid = errors.New("token tidak valid")
	ErrTokenExpired = errors.New("token sudah kedaluwarsa")
)

func GenerateJWTToken(userID uuid.UUID, email string, activeMode string) (string, error) {
	ttl := config.AppConfig.JWTTTL
	if ttl <= 0 {
		ttl = 7 * 24 * time.Hour
	}
	now := time.Now()
	claims := JWTClaims{
		UserID:     userID,
		Email:      email,
		ActiveMode: activeMode,
		IssuedAt:   now.Unix(),
		Exp:        now.Add(ttl).Unix(),
		TokenID:    uuid.NewString(),
	}

	headerJSON, err := json.Marshal(map[string]string{"alg": "HS256", "typ": "JWT"})
	if err != nil {
		return "", err
	}
	claimsJSON, err := json.Marshal(claims)
	if err != nil {
		return "", err
	}

	unsigned := b64(headerJSON) + "." + b64(claimsJSON)
	return unsigned + "." + b64(sign(unsigned)), nil
}

func ValidateJWTToken(tokenString string) (*JWTClaims, error) {
	parts := strings.Split(tokenString, ".")
	if len(parts) != 3 {
		return nil, ErrTokenInvalid
	}

	provided, err := b64decode(parts[2])
	if err != nil {
		return nil, ErrTokenInvalid
	}
	if !hmac.Equal(sign(parts[0]+"."+parts[1]), provided) {
		return nil, ErrTokenInvalid
	}

	claimsJSON, err := b64decode(parts[1])
	if err != nil {
		return nil, ErrTokenInvalid
	}
	var claims JWTClaims
	if err := json.Unmarshal(claimsJSON, &claims); err != nil {
		return nil, ErrTokenInvalid
	}
	if claims.UserID == uuid.Nil {
		return nil, ErrTokenInvalid
	}
	if time.Now().Unix() >= claims.Exp {
		return nil, ErrTokenExpired
	}
	return &claims, nil
}

func sign(data string) []byte {
	mac := hmac.New(sha256.New, []byte(config.AppConfig.JWTSecret))
	mac.Write([]byte(data))
	return mac.Sum(nil)
}

func b64(data []byte) string { return base64.RawURLEncoding.EncodeToString(data) }

func b64decode(s string) ([]byte, error) { return base64.RawURLEncoding.DecodeString(s) }
