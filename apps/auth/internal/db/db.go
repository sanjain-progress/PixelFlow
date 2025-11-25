package db

import (
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/sanjain/pixelflow/apps/auth/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Handler struct {
	DB *gorm.DB
}

func Init(url string) *Handler {
	// Configure GORM logger to show SQL queries
	gormLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second,   // Slow SQL threshold
			LogLevel:                  logger.Info,   // Log level (Info shows all SQL queries)
			IgnoreRecordNotFoundError: true,          // Ignore ErrRecordNotFound error for logger
			Colorful:                  true,          // Enable color
		},
	)

	db, err := gorm.Open(postgres.Open(url), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		log.Fatalln(err)
	}

	// Auto Migrate the User model
	db.AutoMigrate(&models.User{})

	slog.Info("Database connected and migrated")

	return &Handler{DB: db}
}
