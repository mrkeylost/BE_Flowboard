package seed

import (
	"log"

	"github.com/google/uuid"
	"github.com/mrkeylost/Flowboard_Backend/config"
	"github.com/mrkeylost/Flowboard_Backend/model"
	"github.com/mrkeylost/Flowboard_Backend/utils"
)

func SeedAdmin() {
	password, _ := utils.Encrypt("admin123")

	admin := model.User{
		Name:     "Superadmin",
		Email:    "admin@mail.com",
		Password: password,
		Role:     "admin",
		PublicID: uuid.New(),
	}

	err := config.DBConn.FirstOrCreate(&admin, model.User{Email: admin.Email}).Error
	if err != nil {
		log.Println("Data admin seed failed")
	} else {
		log.Println("User admin seeded")
	}
}
