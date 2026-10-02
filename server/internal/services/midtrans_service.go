package services

import (
	"bytes"
	"context"
	"crypto/sha512"
	"crypto/subtle"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"layanesia-server/internal/apperrors"
	"layanesia-server/internal/models"
)

// MidtransConfig adalah konfigurasi minimal yang dibutuhkan klien Midtrans.
type MidtransConfig struct {
	ServerKey    string
	IsProduction bool
	HTTPClient   *http.Client
}

type MidtransSnapRequest struct {
	TransactionDetails MidtransTransactionDetails `json:"transaction_details"`
	CustomerDetails    MidtransCustomerDetails    `json:"customer_details"`
	ItemDetails        []MidtransItemDetail       `json:"item_details"`
	Expiry             *MidtransExpiry            `json:"expiry,omitempty"`
}

type MidtransTransactionDetails struct {
	OrderID  string `json:"order_id"`
	GrossAmt int64  `json:"gross_amount"`
}

type MidtransCustomerDetails struct {
	FirstName string `json:"first_name"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
}

type MidtransItemDetail struct {
	ID       string `json:"id"`
	Price    int64  `json:"price"`
	Quantity int    `json:"quantity"`
	Name     string `json:"name"`
}

type MidtransExpiry struct {
	Unit     string `json:"unit"`
	Duration int    `json:"duration"`
}

type MidtransSnapResponse struct {
	Token       string `json:"token"`
	RedirectURL string `json:"redirect_url"`
}

// MidtransNotification adalah payload webhook / status API yang relevan.
// Lihat https://docs.midtrans.com/reference/http-notification-webhooks
type MidtransNotification struct {
	OrderID           string `json:"order_id"`
	StatusCode        string `json:"status_code"`
	GrossAmount       string `json:"gross_amount"`
	SignatureKey      string `json:"signature_key"`
	TransactionStatus string `json:"transaction_status"`
	FraudStatus       string `json:"fraud_status"`
	PaymentType       string `json:"payment_type"`
	TransactionID     string `json:"transaction_id"`
	TransactionTime   string `json:"transaction_time"`
	SettlementTime    string `json:"settlement_time"`
}

// Outcome menerjemahkan transaction_status + fraud_status menjadi status pembayaran internal.
// Mengembalikan (status, final). final=false berarti belum ada keputusan (pending).
func (n MidtransNotification) Outcome() (models.PaymentStatus, bool) {
	switch strings.ToLower(n.TransactionStatus) {
	case "capture":
		if strings.ToLower(n.FraudStatus) == "challenge" {
			return models.PaymentPending, false
		}
		return models.PaymentPaid, true
	case "settlement":
		return models.PaymentPaid, true
	case "deny", "cancel", "failure":
		return models.PaymentFailed, true
	case "expire":
		return models.PaymentExpired, true
	case "refund", "partial_refund", "chargeback", "partial_chargeback":
		return models.PaymentRefunded, true
	default: // pending, authorize, dll
		return models.PaymentPending, false
	}
}

type MidtransService interface {
	Enabled() bool
	CreateSnapTransaction(ctx context.Context, orderID string, amount int64, itemName string, user *models.User) (*MidtransSnapResponse, error)
	GetTransactionStatus(ctx context.Context, orderID string) (*MidtransNotification, error)
	VerifySignature(n MidtransNotification) bool
}

type midtransService struct {
	cfg MidtransConfig
}

func NewMidtransService(cfg MidtransConfig) MidtransService {
	if cfg.HTTPClient == nil {
		cfg.HTTPClient = &http.Client{Timeout: 15 * time.Second}
	}
	return &midtransService{cfg: cfg}
}

func (s *midtransService) Enabled() bool { return s.cfg.ServerKey != "" }

func (s *midtransService) snapBaseURL() string {
	if s.cfg.IsProduction {
		return "https://app.midtrans.com/snap/v1"
	}
	return "https://app.sandbox.midtrans.com/snap/v1"
}

func (s *midtransService) apiBaseURL() string {
	if s.cfg.IsProduction {
		return "https://api.midtrans.com/v2"
	}
	return "https://api.sandbox.midtrans.com/v2"
}

var errMidtransDisabled = apperrors.Unavailable(
	"Layanan pembayaran belum dikonfigurasi. Hubungi administrator.",
	errors.New("MIDTRANS_SERVER_KEY kosong"),
)

// CreateSnapTransaction membuat transaksi Snap. Tidak ada fallback token palsu:
// jika Midtrans gagal, error dikembalikan dan checkout dibatalkan.
func (s *midtransService) CreateSnapTransaction(ctx context.Context, orderID string, amount int64, itemName string, user *models.User) (*MidtransSnapResponse, error) {
	if !s.Enabled() {
		return nil, errMidtransDisabled
	}

	reqBody := MidtransSnapRequest{
		TransactionDetails: MidtransTransactionDetails{OrderID: orderID, GrossAmt: amount},
		CustomerDetails: MidtransCustomerDetails{
			FirstName: user.Name,
			Email:     user.Email,
			Phone:     user.Phone,
		},
		ItemDetails: []MidtransItemDetail{{ID: orderID, Price: amount, Quantity: 1, Name: itemName}},
		Expiry:      &MidtransExpiry{Unit: "hour", Duration: 24},
	}

	var snapResp MidtransSnapResponse
	if err := s.do(ctx, http.MethodPost, s.snapBaseURL()+"/transactions", reqBody, &snapResp); err != nil {
		return nil, err
	}
	if snapResp.Token == "" {
		return nil, apperrors.Unavailable("Midtrans tidak mengembalikan token pembayaran.", errors.New("empty snap token"))
	}
	return &snapResp, nil
}

// GetTransactionStatus menanyakan status transaksi ke Midtrans (dipakai untuk
// sinkronisasi manual bila webhook terlambat/tidak sampai).
func (s *midtransService) GetTransactionStatus(ctx context.Context, orderID string) (*MidtransNotification, error) {
	if !s.Enabled() {
		return nil, errMidtransDisabled
	}
	var n MidtransNotification
	if err := s.do(ctx, http.MethodGet, s.apiBaseURL()+"/"+orderID+"/status", nil, &n); err != nil {
		return nil, err
	}
	return &n, nil
}

// VerifySignature memvalidasi signature_key webhook:
// SHA512(order_id + status_code + gross_amount + ServerKey).
func (s *midtransService) VerifySignature(n MidtransNotification) bool {
	if !s.Enabled() || n.SignatureKey == "" {
		return false
	}
	sum := sha512.Sum512([]byte(n.OrderID + n.StatusCode + n.GrossAmount + s.cfg.ServerKey))
	expected := hex.EncodeToString(sum[:])
	return subtle.ConstantTimeCompare([]byte(strings.ToLower(n.SignatureKey)), []byte(expected)) == 1
}

func (s *midtransService) do(ctx context.Context, method, url string, body any, out any) error {
	var reader io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return apperrors.Internal(err)
		}
		reader = bytes.NewReader(b)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, reader)
	if err != nil {
		return apperrors.Internal(err)
	}
	req.SetBasicAuth(s.cfg.ServerKey, "")
	req.Header.Set("Accept", "application/json")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := s.cfg.HTTPClient.Do(req)
	if err != nil {
		return apperrors.Unavailable("Tidak dapat menghubungi Midtrans. Coba lagi beberapa saat.", err)
	}
	defer resp.Body.Close()

	raw, _ := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return apperrors.Unavailable(
			"Midtrans menolak permintaan pembayaran.",
			fmt.Errorf("midtrans %s %s -> %d: %s", method, url, resp.StatusCode, truncate(string(raw), 500)),
		)
	}
	if out != nil {
		if err := json.Unmarshal(raw, out); err != nil {
			return apperrors.Unavailable("Respons Midtrans tidak dapat dibaca.", err)
		}
	}
	return nil
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}
