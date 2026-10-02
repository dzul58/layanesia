package controllers

import (
	"layanesia-server/internal/services"

	"github.com/gofiber/fiber/v2"
)

type AuthController struct {
	authService services.AuthService
}

func NewAuthController(authService services.AuthService) *AuthController {
	return &AuthController{authService: authService}
}

func (ctrl *AuthController) Register(c *fiber.Ctx) error {
	var req services.RegisterDTO
	if err := parseBody(c, &req); err != nil {
		return err
	}
	res, err := ctrl.authService.Register(req)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Registrasi berhasil! Selamat datang di Layanesia.",
		"token":   res.Token,
		"user":    res.User,
	})
}

func (ctrl *AuthController) Login(c *fiber.Ctx) error {
	var req services.LoginDTO
	if err := parseBody(c, &req); err != nil {
		return err
	}
	res, err := ctrl.authService.Login(req)
	if err != nil {
		return err
	}
	return c.JSON(fiber.Map{
		"message": "Login berhasil! Selamat datang kembali.",
		"token":   res.Token,
		"user":    res.User,
	})
}
