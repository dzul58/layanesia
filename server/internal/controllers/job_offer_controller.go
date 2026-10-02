package controllers

import (
	"layanesia-server/internal/middlewares"
	"layanesia-server/internal/models"
	"layanesia-server/internal/services"

	"github.com/gofiber/fiber/v2"
)

type JobOfferController struct {
	offerService services.JobOfferService
}

func NewJobOfferController(offerService services.JobOfferService) *JobOfferController {
	return &JobOfferController{offerService: offerService}
}

func (ctrl *JobOfferController) Create(c *fiber.Ctx) error {
	employerID, err := middlewares.MustUserID(c)
	if err != nil {
		return err
	}
	var req services.CreateOfferRequest
	if err := parseBody(c, &req); err != nil {
		return err
	}
	offer, err := ctrl.offerService.CreateOffer(employerID, req)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Berhasil mengirimkan tawaran pekerjaan kepada pekerja.",
		"data":    offer.View(),
	})
}

func (ctrl *JobOfferController) GetReceived(c *fiber.Ctx) error {
	workerID, err := middlewares.MustUserID(c)
	if err != nil {
		return err
	}
	offers, err := ctrl.offerService.GetReceivedOffers(workerID)
	if err != nil {
		return err
	}
	return c.JSON(fiber.Map{"data": models.JobOfferViews(offers), "count": len(offers)})
}

func (ctrl *JobOfferController) GetSent(c *fiber.Ctx) error {
	employerID, err := middlewares.MustUserID(c)
	if err != nil {
		return err
	}
	offers, err := ctrl.offerService.GetSentOffers(employerID)
	if err != nil {
		return err
	}
	return c.JSON(fiber.Map{"data": models.JobOfferViews(offers), "count": len(offers)})
}

func (ctrl *JobOfferController) Respond(c *fiber.Ctx) error {
	workerID, err := middlewares.MustUserID(c)
	if err != nil {
		return err
	}
	offerID, err := paramUUID(c, "id", "ID tawaran")
	if err != nil {
		return err
	}
	var req services.RespondOfferRequest
	if err := parseBody(c, &req); err != nil {
		return err
	}
	offer, err := ctrl.offerService.RespondToOffer(workerID, offerID, req.Status)
	if err != nil {
		return err
	}
	msg := "Tawaran ditolak."
	if offer.Status == models.OfferStatusAccepted {
		msg = "Tawaran diterima! Anda kini TERHUBUNG dengan pemberi kerja."
	}
	return c.JSON(fiber.Map{"message": msg, "data": offer.View()})
}
