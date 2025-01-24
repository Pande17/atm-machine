package rest

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

func jsonResponse(c *fiber.Ctx, statusCode int, message string, errLocate string, data any, deletedAt any) error {
	response := fiber.Map{
		"status":    statusCode,
		"message":   message,
		"timestamp": time.Now().Format(time.RFC3339),
	}

	if data != nil {
		response["data"] = data
	}

	if errLocate != "" {
		response["error_location"] = errLocate
	}

	if deletedAt != nil {
		response["deletedAt"] = deletedAt
	}

	return c.Status(statusCode).JSON(response)
}

func OK(c *fiber.Ctx, message string, data any) error {
	return jsonResponse(c, fiber.StatusOK, message, "", data, nil)
}

func BadRequest(c *fiber.Ctx, message string, errLocate string) error {
	return jsonResponse(c, fiber.StatusOK, message, errLocate, nil, nil)
}

func Conflict(c *fiber.Ctx, message string, errLocate string) error {
	return jsonResponse(c, fiber.StatusConflict, message, errLocate, nil, nil)
}

func Unauthorized(c *fiber.Ctx, message string, errLocate string) error {
	return jsonResponse(c, fiber.StatusUnauthorized, message, errLocate, nil, nil)
}

func NotFound(c *fiber.Ctx, message string, errLocate string) error {
	return jsonResponse(c, fiber.StatusNotFound, message, errLocate, nil, nil)
}

func AlreadyDeleted(c *fiber.Ctx, message string, errLocate string, deletedAt any) error {
	return jsonResponse(c, fiber.StatusNotFound, message, errLocate, nil, deletedAt)
}
