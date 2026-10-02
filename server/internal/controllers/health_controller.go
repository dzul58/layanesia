package controllers

import (
	"context"
	"time"

	"layanesia-server/config"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type HealthController struct {
	db      *gorm.DB
	version string
}

func NewHealthController(db *gorm.DB, version string) *HealthController {
	return &HealthController{db: db, version: version}
}

// Live: proses hidup (untuk liveness probe).
func (h *HealthController) Live(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"status": "ok"})
}

// Ready/CheckHealth: proses hidup DAN database bisa dijangkau (readiness probe).
func (h *HealthController) CheckHealth(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), 2*time.Second)
	defer cancel()

	dbStatus := "ok"
	status := fiber.StatusOK
	if err := config.PingDB(ctx, h.db); err != nil {
		dbStatus = "unreachable"
		status = fiber.StatusServiceUnavailable
	}
	return c.Status(status).JSON(fiber.Map{
		"status":    map[bool]string{true: "ok", false: "degraded"}[status == fiber.StatusOK],
		"timestamp": time.Now().UTC().Format(time.RFC3339),
		"version":   h.version,
		"services":  fiber.Map{"database": dbStatus},
	})
}
