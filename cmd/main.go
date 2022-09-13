package main

import (
	"os"
	models "github.com/aliaqa256/user_location_web_api/models"
	routes "github.com/aliaqa256/user_location_web_api/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		panic("Error loading .env file")
	}
	PORT:= os.Getenv("SERVER_PORT")
	modelApp := models.NewApplication()
	modelApp.ConnectDB()
	modelApp.Migrate()
	defer modelApp.CloseDB()


	app := fiber.New()
	routes.SetupRoutes(app)
	app.Listen(PORT)

}