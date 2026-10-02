package controllers

import (
	"layanesia-server/internal/middlewares"
	"layanesia-server/internal/models"
	"layanesia-server/internal/repositories"
	"layanesia-server/internal/services"

	"github.com/gofiber/fiber/v2"
)

type JobPostingController struct {
	jobService services.JobPostingService
}

func NewJobPostingController(jobService services.JobPostingService) *JobPostingController {
	return &JobPostingController{jobService: jobService}
}

func (ctrl *JobPostingController) Create(c *fiber.Ctx) error {
	employerID, err := middlewares.MustUserID(c)
	if err != nil {
		return err
	}
	var req services.CreateJobRequest
	if err := parseBody(c, &req); err != nil {
		return err
	}
	job, err := ctrl.jobService.CreateJobPosting(employerID, req)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Lowongan pekerjaan berhasil dipublikasikan!",
		"data":    job.View(),
	})
}

func (ctrl *JobPostingController) GetMyPostings(c *fiber.Ctx) error {
	employerID, err := middlewares.MustUserID(c)
	if err != nil {
		return err
	}
	jobs, err := ctrl.jobService.GetMyJobPostings(employerID)
	if err != nil {
		return err
	}
	return c.JSON(fiber.Map{"data": models.JobPostingViews(jobs), "count": len(jobs)})
}

func (ctrl *JobPostingController) Search(c *fiber.Ctx) error {
	limit, offset := pagination(c, 50, 100)
	filter := repositories.JobFilter{
		Province: c.Query("province"),
		City:     c.Query("city"),
		District: c.Query("district"),
		Category: c.Query("category"),
		Search:   c.Query("search"),
		Status:   c.Query("status"),
		Limit:    limit,
		Offset:   offset,
	}
	jobs, total, err := ctrl.jobService.SearchJobPostings(filter)
	if err != nil {
		return err
	}
	return c.JSON(fiber.Map{
		"data":   models.JobPostingViews(jobs),
		"count":  len(jobs),
		"total":  total,
		"limit":  limit,
		"offset": offset,
	})
}

func (ctrl *JobPostingController) GetByID(c *fiber.Ctx) error {
	id, err := paramUUID(c, "id", "ID lowongan")
	if err != nil {
		return err
	}
	job, err := ctrl.jobService.GetJobByID(id)
	if err != nil {
		return err
	}
	return c.JSON(fiber.Map{"data": job.View()})
}

type closeJobRequest struct {
	Status models.JobStatus `json:"status"` // DONE (default) | CANCELLED
}

func (ctrl *JobPostingController) Close(c *fiber.Ctx) error {
	employerID, err := middlewares.MustUserID(c)
	if err != nil {
		return err
	}
	id, err := paramUUID(c, "id", "ID lowongan")
	if err != nil {
		return err
	}
	req := closeJobRequest{Status: models.JobStatusDone}
	if len(c.Body()) > 0 {
		if err := parseBody(c, &req); err != nil {
			return err
		}
	}
	if req.Status == "" {
		req.Status = models.JobStatusDone
	}
	job, err := ctrl.jobService.CloseJob(employerID, id, req.Status)
	if err != nil {
		return err
	}
	return c.JSON(fiber.Map{"message": "Lowongan ditutup.", "data": job.View()})
}
