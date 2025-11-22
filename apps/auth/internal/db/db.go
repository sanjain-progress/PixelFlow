package db

import (
	"fmt"
	"log"

	"github.com/sanjain/pixelflow/apps/auth/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Handler struct {
	DB *gorm.DB
}

func Init(url string) *Handler {
	db, err := gorm.Open(postgres.Open(url), &gorm.Config{})
	if err != nil {
		log.Fatalln(err)
	}

	// Auto Migrate the User model
	db.AutoMigrate(&models.User{})

	fmt.Println("Database connected and migrated")

	return &Handler{DB: db}
}
