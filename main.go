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

	// Auth Instance
	userRepo := repository.NewUserRepository()
	userService := services.NewUserService(userRepo)
	userController := controller.NewUserController(userService)

	// Board Member Instance
	boardMemberRepo := repository.NewBoardMemberRepository()

	// Board Instance
	boardRepo := repository.NewBoardRepository()
	boardService := services.NewBoardService(boardRepo, userRepo, boardMemberRepo)
	boardController := controller.NewBoardController(boardService)

	route.Setup(app, userController, boardController)

	PORT := config.AppConfig.Port
	log.Println("Listening server on port :", PORT)

	log.Fatal(app.Listen(":" + PORT))
}
