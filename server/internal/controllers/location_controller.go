package controllers

import (
	"layanesia-server/internal/apperrors"
	"layanesia-server/internal/services"

	"github.com/gofiber/fiber/v2"
)

type LocationController struct {
	locationService services.LocationService
}

func NewLocationController(locationService services.LocationService) *LocationController {
	return &LocationController{locationService: locationService}
}

func (ctrl *LocationController) GetProvinces(c *fiber.Ctx) error {
	provinces, err := ctrl.locationService.GetProvinces()
	if err != nil {
		return err
	}
	return c.JSON(fiber.Map{"status": "success", "data": provinces})
}

func (ctrl *LocationController) GetCities(c *fiber.Ctx) error {
	province := c.Query("province")
	if province == "" {
		return apperrors.Validation("Query parameter 'province' wajib diisi.", nil)
	}
	cities, err := ctrl.locationService.GetCitiesByProvince(province)
	if err != nil {
		return err
	}
	return c.JSON(fiber.Map{"status": "success", "data": cities})
}

func (ctrl *LocationController) GetDistricts(c *fiber.Ctx) error {
	city := c.Query("city")
	if city == "" {
		return apperrors.Validation("Query parameter 'city' wajib diisi.", nil)
	}
	districts, err := ctrl.locationService.GetDistrictsByCity(city)
	if err != nil {
		return err
	}
	return c.JSON(fiber.Map{"status": "success", "data": districts})
}

func (ctrl *LocationController) SearchLocations(c *fiber.Ctx) error {
	locations, err := ctrl.locationService.SearchLocations(c.Query("q"))
	if err != nil {
		return err
	}
	return c.JSON(fiber.Map{"status": "success", "data": locations})
}

func (ctrl *LocationController) GeocodeNominatim(c *fiber.Ctx) error {
	results, err := ctrl.locationService.GeocodeNominatim(c.Context(), c.Query("q"))
	if err != nil {
		return err
	}
	return c.JSON(fiber.Map{"status": "success", "data": results})
}
