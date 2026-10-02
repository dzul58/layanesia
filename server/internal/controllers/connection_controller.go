package controllers

import (
	"layanesia-server/internal/middlewares"
	"layanesia-server/internal/models"
	"layanesia-server/internal/services"

	"github.com/gofiber/fiber/v2"
)

type ConnectionController struct {
	connService services.ConnectionService
}

func NewConnectionController(connService services.ConnectionService) *ConnectionController {
	return &ConnectionController{connService: connService}
}

// List: GET /connections?status=ACTIVE|COMPLETED (tanpa status = semua).
func (ctrl *ConnectionController) List(c *fiber.Ctx) error {
	userID, err := middlewares.MustUserID(c)
	if err != nil {
		return err
	}
	var status *models.ConnectionStatus
	if raw := c.Query("status"); raw != "" {
		s := models.ConnectionStatus(raw)
		status = &s
	}
	conns, err := ctrl.connService.GetConnections(userID, status)
	if err != nil {
		return err
	}
	return c.JSON(fiber.Map{"data": conns, "count": len(conns)})
}

// GetActive: alias kompatibel untuk klien lama (GET /connections/active).
func (ctrl *ConnectionController) GetActive(c *fiber.Ctx) error {
	userID, err := middlewares.MustUserID(c)
	if err != nil {
		return err
	}
	status := models.ConnStatusActive
	conns, err := ctrl.connService.GetConnections(userID, &status)
	if err != nil {
		return err
	}
	return c.JSON(fiber.Map{"data": conns, "count": len(conns)})
}

func (ctrl *ConnectionController) GetByID(c *fiber.Ctx) error {
	userID, err := middlewares.MustUserID(c)
	if err != nil {
		return err
	}
	id, err := paramUUID(c, "id", "ID koneksi")
	if err != nil {
		return err
	}
	conn, err := ctrl.connService.GetConnectionByID(userID, id)
	if err != nil {
		return err
	}
	return c.JSON(fiber.Map{"data": conn})
}

func (ctrl *ConnectionController) Complete(c *fiber.Ctx) error {
	userID, err := middlewares.MustUserID(c)
	if err != nil {
		return err
	}
	id, err := paramUUID(c, "id", "ID koneksi")
	if err != nil {
		return err
	}
	conn, err := ctrl.connService.CompleteConnection(userID, id)
	if err != nil {
		return err
	}
	return c.JSON(fiber.Map{"message": "Pekerjaan ditandai selesai. Anda kini dapat memberi ulasan.", "data": conn})
}
