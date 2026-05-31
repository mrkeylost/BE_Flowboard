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

func Setup(app *fiber.App, authController *controller.AuthController) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error load env value")
	}

	public := app.Group("/api/auth")
	public.Post("/register", authController.Register)
	public.Post("/login", authController.Login)

	protected := app.Group("/api/auth", jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{
			Key: []byte(config.AppConfig.JWTSecret),
		},
		ErrorHandler: func(ctx fiber.Ctx, err error) error {
			return utils.Unauthorized(ctx, "Unauthorized user", err.Error())
		},
	}))

	protected.Get("/", authController.GetAllUser)
	protected.Get("/:id", authController.GetUserDetail)
	protected.Put("/:id", authController.UpdateUser)
	protected.Delete("/:id", authController.DeleteUser)
}
