package models

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestPublicUserOmitsPII(t *testing.T) {
	ktp := "3174012304900001"
	phone := "081234567890"
	email := "secret@example.com"
	addr := "Jl. Rahasia No. 1"
	u := &User{
		ID:            uuid.New(),
		Email:         email,
		Phone:         phone,
		Name:          "Budi",
		Password:      "hash",
		ActiveMode:    ModeSeeker,
		IsKTPVerified: true,
		KTPNumber:     &ktp,
		KTPImageURL:   &addr,
		AddressDetail: &addr,
		Province:      "DKI Jakarta",
		City:          "Jakarta Selatan",
		District:      "Tebet",
		CreatedAt:     time.Now(),
	}

	raw, err := json.Marshal(u.Public())
	if err != nil {
		t.Fatal(err)
	}
	s := string(raw)
	for _, secret := range []string{email, phone, ktp, addr, "hash"} {
		if strings.Contains(s, secret) {
			t.Errorf("public JSON leaked %q: %s", secret, s)
		}
	}
}

func TestOwnerJSONHidesPassword(t *testing.T) {
	u := User{Email: "a@b.c", Password: "super-secret-hash"}
	raw, _ := json.Marshal(u)
	if strings.Contains(string(raw), "super-secret-hash") {
		t.Fatal("password leaked in owner JSON")
	}
}
