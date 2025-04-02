package database

import (
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	db  *gorm.DB
	err error
)

// Creates schema echo and 2 models: echo.User and echo.Task
func Init() {
	db, err = gorm.Open(postgres.Open(os.Getenv("POSTGRES_DSN")), &gorm.Config{})

	if err != nil {
		log.Fatal("Database failed to open: %w", err)
	}

	err = db.Exec("CREATE SCHEMA IF NOT EXISTS echo").Error
	if err != nil {
		log.Fatal("Failed to create schema: %w", err)
	}

	db.Exec("SET search_path TO echo")

	db.AutoMigrate(&UserStored{}, &TaskStored{})
}
