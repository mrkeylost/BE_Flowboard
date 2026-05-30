package route

import (
	"log"

	"github.com/gofiber/fiber/v3"
	"github.com/joho/godotenv"
	"github.com/mrkeylost/Flowboard_Backend/controller"
)

func Setup(app *fiber.App, authController *controller.AuthController) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error load env value")
	}

	app.Post("/auth/register", authController.Register)
	app.Post("/auth/login", authController.Login)
}
