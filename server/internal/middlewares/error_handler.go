package middlewares

import (
	"errors"
	"log"

	"layanesia-server/internal/apperrors"

	"github.com/gofiber/fiber/v2"
)

// ErrorHandler adalah fiber.Config.ErrorHandler terpusat.
//   - *apperrors.AppError  -> status & body sesuai definisi (aman untuk klien)
//   - *fiber.Error         -> status bawaan Fiber (404 route, 413 body too large, dll)
//   - error lain           -> dicatat ke log, klien menerima 500 generik
func ErrorHandler(c *fiber.Ctx, err error) error {
	if err == nil {
		return nil
	}

	if ae, ok := apperrors.As(err); ok {
		if ae.Status >= 500 {
			log.Printf("[error] %s %s -> %d %s: %v", c.Method(), c.OriginalURL(), ae.Status, ae.Code, err)
		}
		return c.Status(ae.Status).JSON(ae)
	}

	var fe *fiber.Error
	if errors.As(err, &fe) {
		code := "REQUEST_ERROR"
		switch fe.Code {
		case fiber.StatusNotFound:
			code = "ROUTE_NOT_FOUND"
		case fiber.StatusRequestEntityTooLarge:
			code = "PAYLOAD_TOO_LARGE"
		case fiber.StatusMethodNotAllowed:
			code = "METHOD_NOT_ALLOWED"
		}
		return c.Status(fe.Code).JSON(fiber.Map{"code": code, "error": fe.Message})
	}

	log.Printf("[error] %s %s -> 500 unhandled: %v", c.Method(), c.OriginalURL(), err)
	internal := apperrors.Internal(err)
	return c.Status(internal.Status).JSON(internal)
}
