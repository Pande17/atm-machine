package routes

import (
	"atm-machine/handler/rest"

	"github.com/gofiber/fiber/v2"
)

func Routes(r *fiber.App) {
	api := r.Group("/api")

	api.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "TES",
		})
	})

	api.Post("/account/register", rest.RegisterAccount)
	api.Delete("/account/:id", rest.DeleteAccount)
	api.Get("/account", rest.GetAllAccount)
}
