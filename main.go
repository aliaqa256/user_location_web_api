package main

import (
	"github.com/gofiber/fiber/v2"
	routes "github.com/aliaqa256/user_location_web_api/routes"
	models "github.com/aliaqa256/user_location_web_api/models"
)

func main() {

	modelApp := models.NewApplication()
	modelApp.ConnectDB()
	modelApp.Migrate()
	defer modelApp.CloseDB()


	app := fiber.New()
	routes.SetupRoutes(app)
	app.Listen(":4001")

}