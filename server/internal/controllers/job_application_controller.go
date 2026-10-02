package controllers

import (
	"layanesia-server/internal/middlewares"
	"layanesia-server/internal/models"
	"layanesia-server/internal/services"

	"github.com/gofiber/fiber/v2"
)

type JobApplicationController struct {
	appService services.JobApplicationService
}

func NewJobApplicationController(appService services.JobApplicationService) *JobApplicationController {
	return &JobApplicationController{appService: appService}
}

func (ctrl *JobApplicationController) Apply(c *fiber.Ctx) error {
	applicantID, err := middlewares.MustUserID(c)
	if err != nil {
		return err
	}
	jobID, err := paramUUID(c, "id", "ID lowongan")
	if err != nil {
		return err
	}
	app, err := ctrl.appService.ApplyToJob(applicantID, jobID)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Berhasil melamar lowongan pekerjaan. Lamaran Anda sedang ditinjau oleh Pemberi Kerja.",
		"data":    app.View(),
	})
}

func (ctrl *JobApplicationController) GetApplicants(c *fiber.Ctx) error {
	employerID, err := middlewares.MustUserID(c)
	if err != nil {
		return err
	}
	jobID, err := paramUUID(c, "id", "ID lowongan")
	if err != nil {
		return err
	}
	apps, err := ctrl.appService.GetApplicantsForJob(employerID, jobID)
	if err != nil {
		return err
	}
	return c.JSON(fiber.Map{"data": models.JobApplicationViews(apps), "count": len(apps)})
}

func (ctrl *JobApplicationController) GetMyApplications(c *fiber.Ctx) error {
	applicantID, err := middlewares.MustUserID(c)
	if err != nil {
		return err
	}
	apps, err := ctrl.appService.GetMyApplications(applicantID)
	if err != nil {
		return err
	}
	return c.JSON(fiber.Map{"data": models.JobApplicationViews(apps), "count": len(apps)})
}

func (ctrl *JobApplicationController) Accept(c *fiber.Ctx) error {
	return ctrl.respond(c, models.AppStatusAccepted,
		"Selamat! Anda telah menerima pelamar ini dan status Anda kini TERHUBUNG.")
}

func (ctrl *JobApplicationController) Reject(c *fiber.Ctx) error {
	return ctrl.respond(c, models.AppStatusRejected, "Lamaran ditolak.")
}

func (ctrl *JobApplicationController) respond(c *fiber.Ctx, status models.ApplicationStatus, msg string) error {
	employerID, err := middlewares.MustUserID(c)
	if err != nil {
		return err
	}
	appID, err := paramUUID(c, "id", "ID lamaran")
	if err != nil {
		return err
	}
	app, err := ctrl.appService.UpdateApplicationStatus(employerID, appID, status)
	if err != nil {
		return err
	}
	return c.JSON(fiber.Map{"message": msg, "data": app.View()})
}
