package routes

import (
	"siashish/application/controllers"

	"github.com/gofiber/fiber/v2"
)

func UserRoute(app *fiber.App) {
	//All routes related to users comes here
	app.Post("/user", controllers.CreateUser)
	app.Get("/user/:username", controllers.GetAUser)
	app.Patch("/user/:username", controllers.EditAUser)
	app.Delete("/user/:username", controllers.DeleteAUser)
	app.Get("/users", controllers.GetAllUsers)
}
