package utils

import "testing"

func TestNormalizePhone(t *testing.T) {
	cases := map[string]string{
		"081234567890":    "081234567890",
		"+6281234567890":  "081234567890",
		"6281234567890":   "081234567890",
		"0812-3456-7890":  "081234567890",
	}
	for in, want := range cases {
		if got := NormalizePhone(in); got != want {
			t.Errorf("NormalizePhone(%q)=%q want %q", in, got, want)
		}
	}
}

func TestValidateRegisterFields(t *testing.T) {
	type dto struct {
		Email    string `json:"email" validate:"required,email"`
		Phone    string `json:"phone" validate:"required,id_phone"`
		Password string `json:"password" validate:"required,min=8"`
		NIK      string `json:"ktp_number" validate:"required,nik"`
	}
	if err := ValidateStruct(dto{
		Email: "budi@example.com", Phone: "081234567890", Password: "Password123!", NIK: "3174012304900001",
	}); err != nil {
		t.Fatalf("valid dto rejected: %v", err)
	}
	if err := ValidateStruct(dto{Email: "not-an-email", Phone: "123", Password: "short", NIK: "00"}); err == nil {
		t.Fatal("invalid dto accepted")
	}
}
