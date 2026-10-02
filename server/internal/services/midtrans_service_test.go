package services

import (
	"crypto/sha512"
	"encoding/hex"
	"testing"

	"layanesia-server/internal/models"
)

func TestMidtransSignatureAndOutcome(t *testing.T) {
	svc := NewMidtransService(MidtransConfig{ServerKey: "supersecret"})
	n := MidtransNotification{
		OrderID:           "SUB-1",
		StatusCode:        "200",
		GrossAmount:       "5000.00",
		TransactionStatus: "settlement",
	}
	sum := sha512.Sum512([]byte(n.OrderID + n.StatusCode + n.GrossAmount + "supersecret"))
	n.SignatureKey = hex.EncodeToString(sum[:])

	if !svc.VerifySignature(n) {
		t.Fatal("valid signature rejected")
	}
	n.SignatureKey = "deadbeef"
	if svc.VerifySignature(n) {
		t.Fatal("invalid signature accepted")
	}

	paid, final := MidtransNotification{TransactionStatus: "settlement"}.Outcome()
	if paid != models.PaymentPaid || !final {
		t.Fatalf("settlement -> %+v %v", paid, final)
	}
	pending, final := MidtransNotification{TransactionStatus: "pending"}.Outcome()
	if pending != models.PaymentPending || final {
		t.Fatalf("pending -> %+v %v", pending, final)
	}
	challenge, final := MidtransNotification{TransactionStatus: "capture", FraudStatus: "challenge"}.Outcome()
	if challenge != models.PaymentPending || final {
		t.Fatalf("challenge capture should stay pending, got %s final=%v", challenge, final)
	}
}

func TestMidtransDisabledWithoutKey(t *testing.T) {
	svc := NewMidtransService(MidtransConfig{})
	if svc.Enabled() {
		t.Fatal("empty server key must be disabled")
	}
	if svc.VerifySignature(MidtransNotification{SignatureKey: "x"}) {
		t.Fatal("disabled service must reject signatures")
	}
}
