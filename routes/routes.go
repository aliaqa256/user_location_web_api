package routes

import (
	"github.com/gofiber/fiber/v2"
	controllers "github.com/aliaqa256/user_location_web_api/controllers"
)

func SetupRoutes(app *fiber.App) {
	app.Get("/", controllers.RootHandler)
	app.Post("/user", controllers.CreateUserHandler)
	app.Post("/user/info",controllers.AddUserInfoHandler)
	app.Get("/user/lastlocation/:id", controllers.GetLastLocationHandler)
	app.Post("/user/pastlocations",controllers.GetPastLocationsHandler)

}

