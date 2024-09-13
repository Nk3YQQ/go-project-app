package database

import (
	"backend/internal/models"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func LoadDotenv() {
	err := godotenv.Load()

	if err != nil {
		log.Fatalf("Error during loding .env")
	}
}

func InitDB() {
	LoadDotenv()

	var err error

	dsn := fmt.Sprintf(
		"user=%s password=%s host=%s port=%s dbname=%s sslmode=%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
		os.Getenv("SSL_MODE"),
	)

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Failed connection to database", err)
	}

	if err := DB.AutoMigrate(&models.User{}, &models.Project{}, &models.Task{}); err != nil {
		log.Fatal("Failed to automigrate models: ", err)
	}
}
