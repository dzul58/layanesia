package utils

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"layanesia-server/internal/apperrors"

	"github.com/go-playground/validator/v10"
)

var (
	validate    = newValidator()
	phoneRegexp = regexp.MustCompile(`^08[1-9][0-9]{6,11}$`)
	nikRegexp   = regexp.MustCompile(`^[0-9]{16}$`)
	dateRegexp  = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)
)

func newValidator() *validator.Validate {
	v := validator.New(validator.WithRequiredStructEnabled())
	// Gunakan nama field JSON pada pesan error agar mudah dipetakan oleh frontend.
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" || name == "" {
			return fld.Name
		}
		return name
	})
	_ = v.RegisterValidation("id_phone", func(fl validator.FieldLevel) bool {
		return phoneRegexp.MatchString(NormalizePhone(fl.Field().String()))
	})
	_ = v.RegisterValidation("nik", func(fl validator.FieldLevel) bool {
		return nikRegexp.MatchString(fl.Field().String())
	})
	_ = v.RegisterValidation("ymd", func(fl validator.FieldLevel) bool {
		return dateRegexp.MatchString(fl.Field().String())
	})
	return v
}

// ValidateStruct memvalidasi DTO berdasarkan tag `validate` dan mengubah hasilnya
// menjadi AppError VALIDATION_ERROR dengan pesan per-field dalam Bahasa Indonesia.
func ValidateStruct(s any) error {
	err := validate.Struct(s)
	if err == nil {
		return nil
	}
	verrs, ok := err.(validator.ValidationErrors)
	if !ok {
		return apperrors.BadRequest("Format request tidak valid")
	}
	fields := make(map[string]any, len(verrs))
	for _, fe := range verrs {
		fields[fe.Field()] = messageFor(fe)
	}
	return apperrors.Validation("Data yang dikirim tidak valid", fields)
}

func messageFor(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "wajib diisi"
	case "email":
		return "format email tidak valid"
	case "min":
		if fe.Kind() == reflect.String {
			return fmt.Sprintf("minimal %s karakter", fe.Param())
		}
		return fmt.Sprintf("minimal %s", fe.Param())
	case "max":
		if fe.Kind() == reflect.String {
			return fmt.Sprintf("maksimal %s karakter", fe.Param())
		}
		return fmt.Sprintf("maksimal %s", fe.Param())
	case "gt":
		return fmt.Sprintf("harus lebih besar dari %s", fe.Param())
	case "gte":
		return fmt.Sprintf("harus lebih besar atau sama dengan %s", fe.Param())
	case "oneof":
		return fmt.Sprintf("harus salah satu dari: %s", strings.ReplaceAll(fe.Param(), " ", ", "))
	case "id_phone":
		return "format nomor HP Indonesia tidak valid (contoh: 081234567890)"
	case "nik":
		return "NIK harus 16 digit angka"
	case "ymd":
		return "format tanggal harus YYYY-MM-DD"
	case "url", "http_url":
		return "harus berupa URL yang valid"
	case "uuid", "uuid4":
		return "harus berupa UUID yang valid"
	}
	return "tidak valid"
}

// NormalizePhone menyamakan format nomor HP menjadi 08xxxxxxxxxx.
func NormalizePhone(p string) string {
	p = strings.TrimSpace(p)
	p = strings.NewReplacer(" ", "", "-", "", ".", "", "(", "", ")", "").Replace(p)
	switch {
	case strings.HasPrefix(p, "+62"):
		return "0" + p[3:]
	case strings.HasPrefix(p, "62") && len(p) > 10:
		return "0" + p[2:]
	}
	return p
}

// PhoneToWhatsApp mengubah 08xxxx menjadi 628xxxx untuk tautan wa.me.
func PhoneToWhatsApp(p string) string {
	p = NormalizePhone(p)
	if strings.HasPrefix(p, "0") {
		return "62" + p[1:]
	}
	return p
}
