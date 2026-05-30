package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	DBConn    *gorm.DB
	AppConfig *Config
)

type Config struct {
	Port             string
	DBHost           string
	DBPort           string
	DBUser           string
	DBPassword       string
	DBName           string
	JWTSecret        string
	JWTExpireMinutes string
	JWTRefreshToken  string
	JWTExpire        string
}

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}

	AppConfig = &Config{
		Port:            getEnv("PORT", "3001"),
		DBHost:          getEnv("DB_HOST", "localhost"),
		DBPort:          getEnv("DB_PORT", "5432"),
		DBUser:          getEnv("DB_USER", "postgres"),
		DBPassword:      getEnv("DB_PASSWORD", "secretpassword"),
		DBName:          getEnv("DB_NAME", "Flowboard_DB"),
		JWTSecret:       getEnv("JWT_SECRET", "secret"),
		JWTExpire:       getEnv("JWT_EXPIRED", "1h"),
		JWTRefreshToken: getEnv("REFRESH_TOKEN_EXPIRED", "24h"),
	}
}

func getEnv(key string, fallback string) string {
	value, exist := os.LookupEnv(key)
	if exist {
		return value
	} else {
		return fallback
	}
}

func ConnectDB() {
	conf := AppConfig
	dataSourceName := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", conf.DBHost, conf.DBPort, conf.DBUser, conf.DBPassword, conf.DBName)

	db, err := gorm.Open(postgres.Open(dataSourceName), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect database", err)
	}

	sql, err := db.DB()
	if err != nil {
		log.Fatal("Failed to get database instance", err)
	}

	sql.SetMaxIdleConns(10)
	sql.SetMaxOpenConns(100)
	sql.SetConnMaxLifetime(time.Hour)

	DBConn = db
}
