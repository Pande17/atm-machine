package rest

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

func jsonResponse(c *fiber.Ctx, statusCode int, message string, errLocate string, data any, deletedAt string) error {
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

	if deletedAt != "" {
		response["deletedAt"] = deletedAt
	}

	return c.Status(statusCode).JSON(response)
}

func OK(c *fiber.Ctx, message string, data any) error {
	return jsonResponse(c, fiber.StatusOK, message, "", data, "")
}

func BadRequest(c *fiber.Ctx, message string, errLocate string) error {
	return jsonResponse(c, fiber.StatusOK, message, errLocate, nil, "")
}

func Conflict(c *fiber.Ctx, message string, errLocate string) error {
	return jsonResponse(c, fiber.StatusConflict, message, errLocate, nil, "")
}

func Unauthorized(c *fiber.Ctx, message string, errLocate string) error {
	return jsonResponse(c, fiber.StatusUnauthorized, message, errLocate, nil, "")
}
