package routes

import "github.com/gofiber/fiber/v2"

func Routes(r *fiber.App) {
	api := r.Group("/api")

	api.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "TES",
		})
	}) 

	
}
