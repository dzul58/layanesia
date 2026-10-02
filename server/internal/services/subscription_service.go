package services

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"layanesia-server/internal/apperrors"
	"layanesia-server/internal/models"
	"layanesia-server/internal/repositories"
	"layanesia-server/internal/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CheckoutRequestDTO struct {
	PlanType models.SubscriptionPlan `json:"plan_type" validate:"required,oneof=APPLY_JOB_5K POST_SKILL_10K"`
}

type CheckoutResponseDTO struct {
	SubscriptionID string                    `json:"subscription_id"`
	OrderID        string                    `json:"order_id"`
	PlanType       models.SubscriptionPlan   `json:"plan_type"`
	Amount         int64                     `json:"amount"`
	SnapToken      string                    `json:"snap_token"`
	RedirectURL    string                    `json:"redirect_url"`
	Status         models.SubscriptionStatus `json:"status"`
	PaymentStatus  models.PaymentStatus      `json:"payment_status"`
	Reused         bool                      `json:"reused"`
}

type SubscriptionService interface {
	Checkout(ctx context.Context, userID uuid.UUID, req CheckoutRequestDTO) (*CheckoutResponseDTO, error)
	HandleMidtransWebhook(ctx context.Context, n MidtransNotification) error
	SyncPaymentStatus(ctx context.Context, userID uuid.UUID, subscriptionID uuid.UUID) (*models.Subscription, error)
	GetUserSubscriptions(userID uuid.UUID) ([]models.Subscription, error)
	CheckActiveSubscription(userID uuid.UUID, planType models.SubscriptionPlan) (*models.Subscription, error)
	SimulateActivate(userID uuid.UUID, planType models.SubscriptionPlan) (*models.Subscription, error)
}

type subscriptionService struct {
	db          *gorm.DB
	subRepo     repositories.SubscriptionRepository
	userRepo    repositories.UserRepository
	midtransSvc MidtransService
	allowSim    bool
}

func NewSubscriptionService(
	db *gorm.DB,
	subRepo repositories.SubscriptionRepository,
	userRepo repositories.UserRepository,
	midtransSvc MidtransService,
	allowSimulation bool,
) SubscriptionService {
	return &subscriptionService{
		db:          db,
		subRepo:     subRepo,
		userRepo:    userRepo,
		midtransSvc: midtransSvc,
		allowSim:    allowSimulation,
	}
}

// pendingReuseWindow: checkout PENDING yang lebih muda dari ini dipakai ulang
// (token Snap Midtrans berlaku 24 jam, kita sisakan margin).
const pendingReuseWindow = 20 * time.Hour

func (s *subscriptionService) Checkout(ctx context.Context, userID uuid.UUID, req CheckoutRequestDTO) (*CheckoutResponseDTO, error) {
	if err := utils.ValidateStruct(req); err != nil {
		return nil, err
	}
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, apperrors.Internal(err)
	}
	if user == nil {
		return nil, apperrors.ErrUserNotFound
	}

	// Tolak jika paket yang sama masih aktif.
	if active, err := s.subRepo.FindActiveByUserAndPlan(userID, req.PlanType); err != nil {
		return nil, apperrors.Internal(err)
	} else if active != nil {
		return nil, apperrors.Conflict("SUBSCRIPTION_ALREADY_ACTIVE",
			fmt.Sprintf("Anda sudah memiliki %s yang aktif hingga %s.", req.PlanType.DisplayName(), active.ExpiresAt.Format("02 Jan 2006 15:04")))
	}

	// Pakai ulang checkout PENDING yang masih segar; kadaluarsakan yang basi.
	if pending, err := s.subRepo.FindPendingByUserAndPlan(userID, req.PlanType); err != nil {
		return nil, apperrors.Internal(err)
	} else if pending != nil {
		if time.Since(pending.CreatedAt) < pendingReuseWindow && pending.SnapToken != nil && *pending.SnapToken != "" {
			return toCheckoutResponse(pending, true), nil
		}
		pending.Status = models.SubStatusCancelled
		pending.PaymentStatus = models.PaymentExpired
		if err := s.subRepo.Update(pending); err != nil {
			return nil, apperrors.Internal(err)
		}
	}

	orderID, err := newOrderID()
	if err != nil {
		return nil, apperrors.Internal(err)
	}
	amount := req.PlanType.Price()

	snap, err := s.midtransSvc.CreateSnapTransaction(ctx, orderID, amount, req.PlanType.DisplayName(), user)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	sub := &models.Subscription{
		UserID:        userID,
		PlanType:      req.PlanType,
		Amount:        float64(amount),
		Status:        models.SubStatusPending,
		PaymentStatus: models.PaymentPending,
		OrderID:       orderID,
		SnapToken:     &snap.Token,
		RedirectURL:   &snap.RedirectURL,
		// Periode definitif ditetapkan saat pembayaran terkonfirmasi (activate).
		StartsAt:  now,
		ExpiresAt: now.Add(models.SubscriptionDuration),
	}
	if err := s.subRepo.Create(sub); err != nil {
		return nil, apperrors.Internal(err)
	}
	return toCheckoutResponse(sub, false), nil
}

// HandleMidtransWebhook memproses notifikasi Midtrans. Signature diverifikasi di
// sini (bukan hanya di controller) agar tidak bisa terlewat.
func (s *subscriptionService) HandleMidtransWebhook(ctx context.Context, n MidtransNotification) error {
	if n.OrderID == "" {
		return apperrors.BadRequest("order_id kosong")
	}
	if !s.midtransSvc.VerifySignature(n) {
		log.Printf("[webhook] signature tidak valid untuk order %s", n.OrderID)
		return apperrors.Forbidden("INVALID_SIGNATURE", "Signature webhook tidak valid.")
	}
	return s.applyNotification(n)
}

// SyncPaymentStatus menarik status terbaru dari Midtrans untuk satu langganan milik user.
func (s *subscriptionService) SyncPaymentStatus(ctx context.Context, userID uuid.UUID, subscriptionID uuid.UUID) (*models.Subscription, error) {
	sub, err := s.subRepo.FindByID(subscriptionID)
	if err != nil {
		return nil, apperrors.Internal(err)
	}
	if sub == nil || sub.UserID != userID {
		return nil, apperrors.NotFound("Langganan tidak ditemukan")
	}
	if sub.PaymentStatus != models.PaymentPending {
		return sub, nil
	}

	n, err := s.midtransSvc.GetTransactionStatus(ctx, sub.OrderID)
	if err != nil {
		return nil, err
	}
	// Status API tidak selalu menyertakan signature; respons datang lewat kanal
	// terautentikasi (basic auth server key) sehingga dipercaya langsung.
	if err := s.applyNotification(*n); err != nil {
		return nil, err
	}
	return s.subRepo.FindByID(subscriptionID)
}

// applyNotification menerapkan hasil transaksi ke record langganan secara idempoten
// dan atomik (row lock FOR UPDATE).
func (s *subscriptionService) applyNotification(n MidtransNotification) error {
	outcome, final := n.Outcome()

	return s.db.Transaction(func(tx *gorm.DB) error {
		repo := s.subRepo.WithTx(tx)

		sub, err := repo.FindByOrderIDForUpdate(n.OrderID)
		if err != nil {
			return apperrors.Internal(err)
		}
		if sub == nil {
			// Order tidak dikenal: bisa dari environment lain. Balas 200 agar Midtrans
			// tidak retry terus, tetapi catat.
			log.Printf("[webhook] order %s tidak ditemukan, diabaikan", n.OrderID)
			return nil
		}

		// Validasi nominal untuk mencegah pembayaran dengan jumlah yang dimanipulasi.
		if n.GrossAmount != "" {
			paid, perr := strconv.ParseFloat(n.GrossAmount, 64)
			if perr == nil && outcome == models.PaymentPaid && paid+0.001 < sub.Amount {
				log.Printf("[webhook] order %s: gross_amount %s < amount %.2f", n.OrderID, n.GrossAmount, sub.Amount)
				return apperrors.BadRequest("Nominal pembayaran tidak sesuai.")
			}
		}

		// Idempoten: sudah final, abaikan notifikasi ulang.
		if sub.PaymentStatus == models.PaymentPaid && outcome == models.PaymentPaid {
			return nil
		}
		if !final {
			return nil
		}

		now := time.Now()
		if n.PaymentType != "" {
			pt := n.PaymentType
			sub.PaymentType = &pt
		}
		sub.PaymentStatus = outcome

		switch outcome {
		case models.PaymentPaid:
			// Perpanjangan: mulai setelah langganan aktif yang ada berakhir.
			startsAt := now
			if active, err := repo.FindActiveByUserAndPlan(sub.UserID, sub.PlanType); err != nil {
				return apperrors.Internal(err)
			} else if active != nil && active.ID != sub.ID && active.ExpiresAt.After(now) {
				startsAt = active.ExpiresAt
			}
			sub.Status = models.SubStatusActive
			sub.PaidAt = &now
			sub.StartsAt = startsAt
			sub.ExpiresAt = startsAt.Add(models.SubscriptionDuration)
		case models.PaymentRefunded:
			sub.Status = models.SubStatusCancelled
		default: // FAILED / EXPIRED
			sub.Status = models.SubStatusCancelled
		}

		if err := repo.Update(sub); err != nil {
			return apperrors.Internal(err)
		}
		log.Printf("[payment] order %s -> %s (%s)", sub.OrderID, sub.PaymentStatus, sub.Status)
		return nil
	})
}

func (s *subscriptionService) GetUserSubscriptions(userID uuid.UUID) ([]models.Subscription, error) {
	if _, err := s.subRepo.ExpireStale(&userID); err != nil {
		return nil, apperrors.Internal(err)
	}
	subs, err := s.subRepo.FindByUserID(userID)
	if err != nil {
		return nil, apperrors.Internal(err)
	}
	return subs, nil
}

func (s *subscriptionService) CheckActiveSubscription(userID uuid.UUID, planType models.SubscriptionPlan) (*models.Subscription, error) {
	sub, err := s.subRepo.FindActiveByUserAndPlan(userID, planType)
	if err != nil {
		return nil, apperrors.Internal(err)
	}
	return sub, nil
}

// SimulateActivate mengaktifkan paket tanpa pembayaran. Hanya tersedia bila
// ENABLE_PAYMENT_SIMULATION=true (dev/staging); di production route-nya tidak didaftarkan
// dan method ini menolak.
func (s *subscriptionService) SimulateActivate(userID uuid.UUID, planType models.SubscriptionPlan) (*models.Subscription, error) {
	if !s.allowSim {
		return nil, apperrors.Forbidden("SIMULATION_DISABLED", "Simulasi pembayaran tidak diizinkan pada environment ini.")
	}
	if !planType.Valid() {
		return nil, apperrors.BadRequest("Paket langganan tidak valid. Pilih APPLY_JOB_5K atau POST_SKILL_10K.")
	}
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, apperrors.Internal(err)
	}
	if user == nil {
		return nil, apperrors.ErrUserNotFound
	}
	if active, err := s.subRepo.FindActiveByUserAndPlan(userID, planType); err != nil {
		return nil, apperrors.Internal(err)
	} else if active != nil {
		return active, nil
	}

	orderID, err := newOrderID()
	if err != nil {
		return nil, apperrors.Internal(err)
	}
	now := time.Now()
	sim := "SIMULATION"
	sub := &models.Subscription{
		UserID:          userID,
		PlanType:        planType,
		Amount:          float64(planType.Price()),
		Status:          models.SubStatusActive,
		PaymentStatus:   models.PaymentPaid,
		PaymentProvider: sim,
		PaymentType:     &sim,
		OrderID:         "SIM-" + orderID,
		PaidAt:          &now,
		StartsAt:        now,
		ExpiresAt:       now.Add(models.SubscriptionDuration),
	}
	if err := s.subRepo.Create(sub); err != nil {
		return nil, apperrors.Internal(err)
	}
	return sub, nil
}

func toCheckoutResponse(sub *models.Subscription, reused bool) *CheckoutResponseDTO {
	res := &CheckoutResponseDTO{
		SubscriptionID: sub.ID.String(),
		OrderID:        sub.OrderID,
		PlanType:       sub.PlanType,
		Amount:         int64(sub.Amount),
		Status:         sub.Status,
		PaymentStatus:  sub.PaymentStatus,
		Reused:         reused,
	}
	if sub.SnapToken != nil {
		res.SnapToken = *sub.SnapToken
	}
	if sub.RedirectURL != nil {
		res.RedirectURL = *sub.RedirectURL
	}
	return res
}

// newOrderID menghasilkan order_id unik yang memenuhi batasan Midtrans (<= 50 char, alnum & '-').
func newOrderID() (string, error) {
	b := make([]byte, 6)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return strings.ToUpper(fmt.Sprintf("SUB-%s-%s", time.Now().Format("20060102150405"), hex.EncodeToString(b))), nil
}
