// Package apperrors mendefinisikan error domain yang aman dikirim ke klien.
// Error lain (SQL, I/O, panic) dianggap internal: dicatat ke log dan dibalas
// dengan pesan generik 500 oleh ErrorHandler terpusat.
package apperrors

import (
	"errors"
	"fmt"
	"net/http"
)

type AppError struct {
	Status  int            `json:"-"`
	Code    string         `json:"code"`
	Message string         `json:"error"`
	Details map[string]any `json:"details,omitempty"`
	cause   error
}

func (e *AppError) Error() string {
	if e.cause != nil {
		return fmt.Sprintf("%s: %s (%v)", e.Code, e.Message, e.cause)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func (e *AppError) Unwrap() error { return e.cause }

// WithDetails mengembalikan salinan error dengan detail tambahan (misal plan_type).
func (e *AppError) WithDetails(details map[string]any) *AppError {
	cp := *e
	cp.Details = details
	return &cp
}

// WithCause menyimpan error asal untuk logging tanpa membocorkannya ke klien.
func (e *AppError) WithCause(err error) *AppError {
	cp := *e
	cp.cause = err
	return &cp
}

func New(status int, code, message string) *AppError {
	return &AppError{Status: status, Code: code, Message: message}
}

func BadRequest(message string) *AppError {
	return New(http.StatusBadRequest, "BAD_REQUEST", message)
}

func Validation(message string, fields map[string]any) *AppError {
	e := New(http.StatusBadRequest, "VALIDATION_ERROR", message)
	if len(fields) > 0 {
		e.Details = map[string]any{"fields": fields}
	}
	return e
}

func Unauthorized(message string) *AppError {
	return New(http.StatusUnauthorized, "UNAUTHORIZED", message)
}

func Forbidden(code, message string) *AppError {
	return New(http.StatusForbidden, code, message)
}

func NotFound(message string) *AppError {
	return New(http.StatusNotFound, "NOT_FOUND", message)
}

func Conflict(code, message string) *AppError {
	return New(http.StatusConflict, code, message)
}

func TooManyRequests(message string) *AppError {
	return New(http.StatusTooManyRequests, "RATE_LIMITED", message)
}

func Internal(err error) *AppError {
	return New(http.StatusInternalServerError, "INTERNAL_ERROR", "Terjadi kesalahan pada server. Silakan coba lagi.").WithCause(err)
}

func Unavailable(message string, err error) *AppError {
	return New(http.StatusServiceUnavailable, "SERVICE_UNAVAILABLE", message).WithCause(err)
}

// As mengekstrak *AppError dari rantai error.
func As(err error) (*AppError, bool) {
	var ae *AppError
	if errors.As(err, &ae) {
		return ae, true
	}
	return nil, false
}

// ---- Error domain yang dipakai lintas service ----

var (
	ErrSessionInvalid = Unauthorized("Sesi tidak valid. Silakan masuk kembali.")
	ErrUserNotFound   = NotFound("Pengguna tidak ditemukan")

	ErrKTPNotVerified = Forbidden("KTP_NOT_VERIFIED",
		"Anda wajib menyelesaikan verifikasi KTP di profil sebelum melanjutkan.")
	ErrResumeNotUploaded = Forbidden("RESUME_NOT_UPLOADED",
		"Anda wajib mengunggah Resume/CV di profil sebelum melamar lowongan.")
	ErrSubscriptionRequired = Forbidden("SUBSCRIPTION_REQUIRED",
		"Anda wajib memiliki paket langganan aktif untuk fitur ini.")
)
