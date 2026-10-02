package controllers

import (
	"layanesia-server/internal/middlewares"
	"layanesia-server/internal/models"
	"layanesia-server/internal/repositories"
	"layanesia-server/internal/services"

	"github.com/gofiber/fiber/v2"
)

type SkillPostingController struct {
	skillService services.SkillPostingService
}

func NewSkillPostingController(skillService services.SkillPostingService) *SkillPostingController {
	return &SkillPostingController{skillService: skillService}
}

func (ctrl *SkillPostingController) Create(c *fiber.Ctx) error {
	userID, err := middlewares.MustUserID(c)
	if err != nil {
		return err
	}
	var req services.CreateSkillRequest
	if err := parseBody(c, &req); err != nil {
		return err
	}
	posting, err := ctrl.skillService.CreateSkillPosting(userID, req)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Profil keahlian berhasil dipasang di etalase!",
		"data":    posting.View(),
	})
}

func (ctrl *SkillPostingController) GetMyPostings(c *fiber.Ctx) error {
	userID, err := middlewares.MustUserID(c)
	if err != nil {
		return err
	}
	postings, err := ctrl.skillService.GetMySkillPostings(userID)
	if err != nil {
		return err
	}
	return c.JSON(fiber.Map{"data": models.SkillPostingViews(postings), "count": len(postings)})
}

func (ctrl *SkillPostingController) Search(c *fiber.Ctx) error {
	limit, offset := pagination(c, 50, 100)
	filter := repositories.SkillFilter{
		Province:     c.Query("province"),
		City:         c.Query("city"),
		District:     c.Query("district"),
		Category:     c.Query("category"),
		Search:       c.Query("search"),
		Availability: c.Query("availability"),
		Limit:        limit,
		Offset:       offset,
	}
	postings, total, err := ctrl.skillService.SearchSkillPostings(filter)
	if err != nil {
		return err
	}
	return c.JSON(fiber.Map{
		"data":   models.SkillPostingViews(postings),
		"count":  len(postings),
		"total":  total,
		"limit":  limit,
		"offset": offset,
	})
}

func (ctrl *SkillPostingController) GetByID(c *fiber.Ctx) error {
	id, err := paramUUID(c, "id", "ID postingan")
	if err != nil {
		return err
	}
	posting, err := ctrl.skillService.GetSkillByID(id)
	if err != nil {
		return err
	}
	return c.JSON(fiber.Map{"data": posting.View()})
}

func (ctrl *SkillPostingController) UpdateAvailability(c *fiber.Ctx) error {
	userID, err := middlewares.MustUserID(c)
	if err != nil {
		return err
	}
	id, err := paramUUID(c, "id", "ID postingan")
	if err != nil {
		return err
	}
	var req services.UpdateSkillAvailabilityRequest
	if err := parseBody(c, &req); err != nil {
		return err
	}
	posting, err := ctrl.skillService.UpdateAvailability(userID, id, req.Availability)
	if err != nil {
		return err
	}
	return c.JSON(fiber.Map{"message": "Status ketersediaan diperbarui.", "data": posting.View()})
}

func (ctrl *SkillPostingController) Delete(c *fiber.Ctx) error {
	userID, err := middlewares.MustUserID(c)
	if err != nil {
		return err
	}
	id, err := paramUUID(c, "id", "ID postingan")
	if err != nil {
		return err
	}
	if err := ctrl.skillService.DeleteSkillPosting(userID, id); err != nil {
		return err
	}
	return c.JSON(fiber.Map{"message": "Postingan keahlian dihapus."})
}
