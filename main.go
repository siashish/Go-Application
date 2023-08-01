package main

import (
	"siashish/application/configs"
	"siashish/application/routes"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(&fiber.Map{"data": "Hello from Fiber & mongoDB By Ashish"})
	})

	//run database
	configs.ConnectDB()

	//routes
	routes.UserRoute(app)

	app.Listen(":6000")
}
