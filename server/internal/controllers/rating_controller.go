package controllers

import (
	"layanesia-server/internal/middlewares"
	"layanesia-server/internal/services"

	"github.com/gofiber/fiber/v2"
)

type RatingController struct {
	ratingService services.RatingService
}

func NewRatingController(ratingService services.RatingService) *RatingController {
	return &RatingController{ratingService: ratingService}
}

func (ctrl *RatingController) Create(c *fiber.Ctx) error {
	reviewerID, err := middlewares.MustUserID(c)
	if err != nil {
		return err
	}
	var req services.CreateRatingRequest
	if err := parseBody(c, &req); err != nil {
		return err
	}
	rating, err := ctrl.ratingService.CreateRating(reviewerID, req)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Terima kasih! Ulasan dan rating Anda telah berhasil disimpan.",
		"data":    rating,
	})
}

func (ctrl *RatingController) GetUserRatings(c *fiber.Ctx) error {
	userID, err := paramUUID(c, "userID", "ID pengguna")
	if err != nil {
		return err
	}
	limit, offset := pagination(c, 20, 100)
	summary, err := ctrl.ratingService.GetUserRatings(userID, limit, offset)
	if err != nil {
		return err
	}
	return c.JSON(fiber.Map{
		"data":           summary.Ratings,
		"average_rating": summary.AverageStars,
		"total_reviews":  summary.Total,
		"limit":          limit,
		"offset":         offset,
	})
}
