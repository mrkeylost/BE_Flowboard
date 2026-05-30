package main

import (
	"log"

	"github.com/gofiber/fiber/v3"
	"github.com/mrkeylost/Flowboard_Backend/config"
	"github.com/mrkeylost/Flowboard_Backend/controller"
	"github.com/mrkeylost/Flowboard_Backend/db/seed"
	"github.com/mrkeylost/Flowboard_Backend/repository"
	"github.com/mrkeylost/Flowboard_Backend/route"
	"github.com/mrkeylost/Flowboard_Backend/services"
)

func main() {
	config.LoadEnv()
	config.ConnectDB()

	seed.SeedAdmin()

	app := fiber.New()

	userRepo := repository.NewUserRepository()
	userService := services.NewUserService(userRepo)
	userController := controller.NewUserController(userService)

	route.Setup(app, userController)

	PORT := config.AppConfig.Port
	log.Println("Listening server on port :", PORT)

	log.Fatal(app.Listen(":" + PORT))
}
