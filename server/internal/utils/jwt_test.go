package utils

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"layanesia-server/config"

	"github.com/google/uuid"
)

func TestMain(m *testing.M) {
	config.AppConfig.JWTSecret = "test-jwt-secret-at-least-32-characters-long"
	config.AppConfig.JWTTTL = time.Hour
	m.Run()
}

func TestJWTRoundTrip(t *testing.T) {
	id := uuid.New()
	tok, err := GenerateJWTToken(id, "a@b.c", "SEEKER")
	if err != nil {
		t.Fatal(err)
	}
	claims, err := ValidateJWTToken(tok)
	if err != nil {
		t.Fatal(err)
	}
	if claims.UserID != id || claims.Email != "a@b.c" || claims.ActiveMode != "SEEKER" {
		t.Fatalf("claims mismatch: %+v", claims)
	}
}

func TestJWTRejectsTamperAndNoneAlg(t *testing.T) {
	id := uuid.New()
	tok, err := GenerateJWTToken(id, "a@b.c", "SEEKER")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := ValidateJWTToken(tok[:len(tok)-2] + "xx"); err == nil {
		t.Fatal("expected invalid signature")
	}

	parts := strings.Split(tok, ".")
	noneHeader := b64([]byte(`{"alg":"none","typ":"JWT"}`))
	if _, err := ValidateJWTToken(noneHeader + "." + parts[1] + "."); err == nil {
		t.Fatal("alg=none must be rejected")
	}
	if _, err := ValidateJWTToken("not-a-jwt"); err == nil {
		t.Fatal("malformed token must be rejected")
	}
}

func TestJWTExpired(t *testing.T) {
	id := uuid.New()
	headerJSON, _ := json.Marshal(map[string]string{"alg": "HS256", "typ": "JWT"})
	claims := JWTClaims{UserID: id, Email: "a@b.c", ActiveMode: "SEEKER", IssuedAt: 1, Exp: time.Now().Add(-time.Hour).Unix(), TokenID: "x"}
	claimsJSON, _ := json.Marshal(claims)
	unsigned := b64(headerJSON) + "." + b64(claimsJSON)
	tok := unsigned + "." + b64(sign(unsigned))
	if _, err := ValidateJWTToken(tok); err != ErrTokenExpired {
		t.Fatalf("want ErrTokenExpired, got %v", err)
	}
}
