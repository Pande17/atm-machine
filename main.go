package main

import (
	"atm-machine/config"
	"atm-machine/database"
	"atm-machine/routes"
	"log"

	"github.com/gofiber/fiber/v2"
)

func main() {
	config.InitEnv()

	database.ConnectMongoDB()

	app := fiber.New(fiber.Config{
		Network: "tcp",
	})
	routes.Routes(app)

	log.Fatal(app.Listen(":3000"))
}
