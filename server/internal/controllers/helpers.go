package controllers

import (
	"strconv"

	"layanesia-server/internal/apperrors"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// parseBody membaca JSON body; error parsing dikembalikan sebagai 400 dengan kode
// INVALID_BODY tanpa membocorkan detail internal parser.
func parseBody(c *fiber.Ctx, dst any) error {
	if err := c.BodyParser(dst); err != nil {
		return apperrors.New(fiber.StatusBadRequest, "INVALID_BODY", "Format request tidak valid (JSON tidak dapat dibaca).").WithCause(err)
	}
	return nil
}

// paramUUID mem-parse path param sebagai UUID.
func paramUUID(c *fiber.Ctx, name, label string) (uuid.UUID, error) {
	id, err := uuid.Parse(c.Params(name))
	if err != nil {
		return uuid.Nil, apperrors.New(fiber.StatusBadRequest, "INVALID_ID", label+" tidak valid.")
	}
	return id, nil
}

// pagination membaca limit/offset query dengan batas aman.
func pagination(c *fiber.Ctx, defLimit, maxLimit int) (limit, offset int) {
	limit, _ = strconv.Atoi(c.Query("limit", strconv.Itoa(defLimit)))
	offset, _ = strconv.Atoi(c.Query("offset", "0"))
	if limit <= 0 {
		limit = defLimit
	}
	if limit > maxLimit {
		limit = maxLimit
	}
	if offset < 0 {
		offset = 0
	}
	return limit, offset
}
