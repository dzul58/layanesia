package app_test

import (
	"bytes"
	"context"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"layanesia-server/config"
	dbmigrate "layanesia-server/db"
	"layanesia-server/internal/app"
	"layanesia-server/internal/models"
	"layanesia-server/internal/services"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type stubMidtrans struct{ key string }

func (m stubMidtrans) Enabled() bool { return true }

func (m stubMidtrans) CreateSnapTransaction(_ context.Context, orderID string, _ int64, _ string, _ *models.User) (*services.MidtransSnapResponse, error) {
	return &services.MidtransSnapResponse{Token: "snap-test", RedirectURL: "https://pay.test/" + orderID}, nil
}

func (m stubMidtrans) GetTransactionStatus(_ context.Context, orderID string) (*services.MidtransNotification, error) {
	return &services.MidtransNotification{OrderID: orderID, TransactionStatus: "settlement", StatusCode: "200", GrossAmount: "5000.00"}, nil
}

func (m stubMidtrans) VerifySignature(n services.MidtransNotification) bool {
	if n.SignatureKey == "" {
		return false
	}
	sum := sha512.Sum512([]byte(n.OrderID + n.StatusCode + n.GrossAmount + m.key))
	return hex.EncodeToString(sum[:]) == strings.ToLower(n.SignatureKey)
}

type stubUpload struct{}

func (stubUpload) Upload(_ context.Context, _ uuid.UUID, _ services.UploadPurpose, _ *multipart.FileHeader) (*services.UploadResult, error) {
	return &services.UploadResult{URL: "https://cdn.test/file.pdf", MIMEType: "application/pdf", Storage: "memory"}, nil
}
func (stubUpload) LocalDir() string       { return "" }
func (stubUpload) UsesLocalStorage() bool { return false }

func TestProtectedRoutesAndPII(t *testing.T) {
	fiberApp, gdb, cleanup := newTestApp(t)
	defer cleanup()

	email := fmt.Sprintf("qa.%s@test.local", uuid.NewString()[:8])
	phone := fmt.Sprintf("0818%08d", time.Now().UnixNano()%100000000)
	body := fmt.Sprintf(`{"name":"QA User","email":%q,"phone":%q,"password":"Password123!","province":"DKI Jakarta","city":"Jakarta Selatan","district":"Tebet"}`, email, phone)

	code, raw := do(t, fiberApp, "POST", "/api/v1/auth/register", "", body)
	if code != http.StatusCreated {
		t.Fatalf("register: %d %s", code, raw)
	}
	token := jsonStr(raw, "token")
	if token == "" {
		t.Fatalf("missing token: %s", raw)
	}
	t.Cleanup(func() {
		gdb.Exec("DELETE FROM users WHERE email = ?", email)
	})

	code, raw = do(t, fiberApp, "GET", "/api/v1/users/me", token, "")
	if code != http.StatusOK {
		t.Fatalf("me: %d %s", code, raw)
	}

	// Regression P0-1: Locals key mismatch used to 401 every protected skill/job/offer route.
	for _, path := range []string{
		"/api/v1/skills/my-postings",
		"/api/v1/jobs/my-postings",
		"/api/v1/offers/received",
		"/api/v1/connections/active",
	} {
		code, raw = do(t, fiberApp, "GET", path, token, "")
		if code != http.StatusOK {
			t.Fatalf("%s: %d %s", path, code, raw)
		}
	}

	code, raw = do(t, fiberApp, "POST", "/api/v1/skills", token, `{"title":"Tukang Las Harian","description":"Pengalaman las listrik lebih dari lima tahun.","rate_amount":150000}`)
	if code != http.StatusForbidden || !strings.Contains(string(raw), "SUBSCRIPTION_REQUIRED") {
		t.Fatalf("post skill without sub: %d %s", code, raw)
	}

	code, raw = do(t, fiberApp, "POST", "/api/v1/subscriptions/simulate-activate", token, `{"plan_type":"POST_SKILL_10K"}`)
	if code != http.StatusOK {
		t.Fatalf("simulate: %d %s", code, raw)
	}

	code, raw = do(t, fiberApp, "POST", "/api/v1/skills", token, `{"title":"Tukang Las Harian","description":"Pengalaman las listrik lebih dari lima tahun.","rate_amount":150000}`)
	if code != http.StatusCreated {
		t.Fatalf("post skill: %d %s", code, raw)
	}

	code, raw = do(t, fiberApp, "GET", "/api/v1/skills?limit=50", "", "")
	if code != http.StatusOK {
		t.Fatalf("public skills: %d %s", code, raw)
	}
	public := string(raw)
	if strings.Contains(public, email) || strings.Contains(public, phone) {
		t.Fatalf("public skill search leaked PII: %s", public)
	}

	code, raw = do(t, fiberApp, "POST", "/api/v1/subscriptions/webhook", "", `{"order_id":"FAKE","transaction_status":"settlement","signature_key":"nope"}`)
	if code != http.StatusForbidden {
		t.Fatalf("webhook without signature: %d %s", code, raw)
	}

	code, raw = do(t, fiberApp, "POST", "/api/v1/subscriptions/checkout", token, `{"plan_type":"APPLY_JOB_5K"}`)
	if code != http.StatusCreated && code != http.StatusOK {
		t.Fatalf("checkout: %d %s", code, raw)
	}
	if !strings.Contains(string(raw), `"status":"PENDING"`) && !strings.Contains(string(raw), `"payment_status":"PENDING"`) {
		t.Fatalf("checkout must stay PENDING until paid: %s", raw)
	}
}

func newTestApp(t *testing.T) (*fiber.App, *gorm.DB, func()) {
	t.Helper()
	t.Setenv("ENV", "test")
	t.Setenv("DB_HOST", getenvDefault("DB_HOST", "localhost"))
	t.Setenv("DB_PORT", getenvDefault("DB_PORT", "5432"))
	t.Setenv("DB_USER", getenvDefault("DB_USER", "postgres"))
	t.Setenv("DB_PASSWORD", getenvDefault("DB_PASSWORD", "postgres"))
	t.Setenv("DB_NAME", getenvDefault("DB_NAME", "layanesia_db"))
	t.Setenv("DB_SSLMODE", "disable")
	t.Setenv("DB_RUN_MIGRATIONS", "true")
	t.Setenv("JWT_SECRET", "test-jwt-secret-at-least-32-characters-long!!")
	t.Setenv("ENABLE_PAYMENT_SIMULATION", "true")
	t.Setenv("KTP_AUTO_VERIFY", "true")
	t.Setenv("CORS_ALLOW_ORIGINS", "http://localhost:5173")

	cfg, err := config.LoadConfig()
	if err != nil {
		t.Fatal(err)
	}
	gdb, err := config.ConnectDB(cfg)
	if err != nil {
		t.Skipf("database tidak tersedia: %v", err)
	}
	if err := dbmigrate.MigrateUp(cfg.DSN()); err != nil {
		t.Fatalf("migrasi: %v", err)
	}

	fiberApp := app.New(cfg, gdb, app.Options{
		Midtrans: stubMidtrans{key: "test-midtrans-key"},
		Upload:   stubUpload{},
	})
	return fiberApp, gdb, func() {
		sqlDB, _ := gdb.DB()
		_ = sqlDB.Close()
	}
}

func do(t *testing.T, a *fiber.App, method, path, token, body string) (int, []byte) {
	t.Helper()
	var rdr io.Reader
	if body != "" {
		rdr = bytes.NewBufferString(body)
	}
	req := httptest.NewRequest(method, path, rdr)
	if body != "" {
		req.Header.Set("Content-Type", "application/json")
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	resp, err := a.Test(req, 15000)
	if err != nil {
		t.Fatal(err)
	}
	raw, _ := io.ReadAll(resp.Body)
	return resp.StatusCode, raw
}

func jsonStr(raw []byte, key string) string {
	var m map[string]any
	if err := json.Unmarshal(raw, &m); err != nil {
		return ""
	}
	s, _ := m[key].(string)
	return s
}

func getenvDefault(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
