package controllers

import (
	"layanesia-server/internal/apperrors"
	"layanesia-server/internal/middlewares"
	"layanesia-server/internal/models"
	"layanesia-server/internal/services"

	"github.com/gofiber/fiber/v2"
)

type SubscriptionController struct {
	subService services.SubscriptionService
}

func NewSubscriptionController(subService services.SubscriptionService) *SubscriptionController {
	return &SubscriptionController{subService: subService}
}

func (ctrl *SubscriptionController) Checkout(c *fiber.Ctx) error {
	userID, err := middlewares.MustUserID(c)
	if err != nil {
		return err
	}
	var req services.CheckoutRequestDTO
	if err := parseBody(c, &req); err != nil {
		return err
	}
	res, err := ctrl.subService.Checkout(c.Context(), userID, req)
	if err != nil {
		return err
	}
	status := fiber.StatusCreated
	msg := "Transaksi dibuat. Selesaikan pembayaran melalui Midtrans Snap."
	if res.Reused {
		status = fiber.StatusOK
		msg = "Anda masih memiliki transaksi yang belum dibayar; lanjutkan pembayaran."
	}
	return c.Status(status).JSON(fiber.Map{"message": msg, "data": res})
}

func (ctrl *SubscriptionController) GetMySubscriptions(c *fiber.Ctx) error {
	userID, err := middlewares.MustUserID(c)
	if err != nil {
		return err
	}
	subs, err := ctrl.subService.GetUserSubscriptions(userID)
	if err != nil {
		return err
	}
	return c.JSON(fiber.Map{"subscriptions": subs, "count": len(subs)})
}

func (ctrl *SubscriptionController) CheckActiveStatus(c *fiber.Ctx) error {
	userID, err := middlewares.MustUserID(c)
	if err != nil {
		return err
	}
	plan := models.SubscriptionPlan(c.Query("plan_type"))
	if !plan.Valid() {
		return apperrors.Validation("Parameter plan_type wajib diisi.", map[string]any{"plan_type": "APPLY_JOB_5K atau POST_SKILL_10K"})
	}
	sub, err := ctrl.subService.CheckActiveSubscription(userID, plan)
	if err != nil {
		return err
	}
	resp := fiber.Map{"plan_type": plan, "is_active": sub != nil}
	if sub != nil {
		resp["expires_at"] = sub.ExpiresAt
		resp["subscription_id"] = sub.ID
	}
	return c.JSON(resp)
}

// SyncStatus menarik status transaksi dari Midtrans (fallback jika webhook terlambat).
func (ctrl *SubscriptionController) SyncStatus(c *fiber.Ctx) error {
	userID, err := middlewares.MustUserID(c)
	if err != nil {
		return err
	}
	subID, err := paramUUID(c, "id", "ID langganan")
	if err != nil {
		return err
	}
	sub, err := ctrl.subService.SyncPaymentStatus(c.Context(), userID, subID)
	if err != nil {
		return err
	}
	return c.JSON(fiber.Map{"subscription": sub})
}

// SimulateActivate hanya terdaftar jika ENABLE_PAYMENT_SIMULATION=true (non-prod).
func (ctrl *SubscriptionController) SimulateActivate(c *fiber.Ctx) error {
	userID, err := middlewares.MustUserID(c)
	if err != nil {
		return err
	}
	var req services.CheckoutRequestDTO
	if err := parseBody(c, &req); err != nil {
		return err
	}
	sub, err := ctrl.subService.SimulateActivate(userID, req.PlanType)
	if err != nil {
		return err
	}
	return c.JSON(fiber.Map{"message": "Simulasi aktivasi paket langganan berhasil (mode pengembangan).", "subscription": sub})
}

// HandleWebhook menerima HTTP notification Midtrans. Selalu balas 200 untuk
// notifikasi valid (Midtrans akan retry pada non-2xx).
func (ctrl *SubscriptionController) HandleWebhook(c *fiber.Ctx) error {
	var n services.MidtransNotification
	if err := c.BodyParser(&n); err != nil {
		return apperrors.New(fiber.StatusBadRequest, "INVALID_BODY", "Payload notifikasi tidak valid.").WithCause(err)
	}
	if err := ctrl.subService.HandleMidtransWebhook(c.Context(), n); err != nil {
		return err
	}
	return c.JSON(fiber.Map{"status": "ok"})
}
