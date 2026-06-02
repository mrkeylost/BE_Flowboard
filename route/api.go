package route

import (
	"log"

	jwtware "github.com/gofiber/contrib/v3/jwt"
	"github.com/gofiber/fiber/v3"
	"github.com/joho/godotenv"
	"github.com/mrkeylost/Flowboard_Backend/config"
	"github.com/mrkeylost/Flowboard_Backend/controller"
	"github.com/mrkeylost/Flowboard_Backend/utils"
)

func Setup(
	app *fiber.App,
	authController *controller.AuthController,
	boardController *controller.BoardController,
	listController *controller.ListController,
) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error load env value")
	}

	public := app.Group("/api")

	// Authentication Public Route
	publicUser := public.Group("/auth")
	publicUser.Post("/register", authController.Register)
	publicUser.Post("/login", authController.Login)

	protected := app.Group("/api", jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{
			Key: []byte(config.AppConfig.JWTSecret),
		},
		ErrorHandler: func(ctx fiber.Ctx, err error) error {
			return utils.Unauthorized(ctx, "Unauthorized user", err.Error())
		},
	}))

	// Authentication Protected Route
	protectedUser := protected.Group("/auth")
	protectedUser.Get("/", authController.GetAllUser)
	protectedUser.Get("/:id", authController.GetUserDetail)
	protectedUser.Put("/:id", authController.UpdateUser)
	protectedUser.Delete("/:id", authController.DeleteUser)

	// Board Protected Route
	protectedBoard := protected.Group("/board")
	protectedBoard.Get("/my-board", boardController.GetAllBoardByUserID)
	protectedBoard.Post("/", boardController.CreateBoard)
	protectedBoard.Put("/:id", boardController.UpdateBoard)
	protectedBoard.Post("/:id/members", boardController.AddBoardMembers)
	protectedBoard.Delete("/:id/members", boardController.RemoveBoardMembers)

	// List Protected Route
	protectedList := protected.Group("/list")
	protectedList.Post("/", listController.CreateList)
}
