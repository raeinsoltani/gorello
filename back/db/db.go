package db

import (
	"fmt"
	"log"

	"github.com/raeinsoltani/gorello/back/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init() {
	var err error
	dsn := fmt.Sprintf(
		// "host=%s user=%s dbname=%s sslmode=disable password=%s",
		// os.Getenv("DB_HOST"),
		// os.Getenv("DB_USER"),
		// os.Getenv("DB_NAME"),
		// os.Getenv("DB_PASSWORD"),
		"postgres://postgres:password@127.0.0.1:5432/gorello",
	)
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database!", err)
	}
	fmt.Println("Database connected")

	err = DB.AutoMigrate(&models.User{}, &models.UserWorkspaceRole{}, &models.Workspace{}, &models.Task{})
	if err != nil {
		log.Fatal("Failed to migrate database!", err)
	}
	fmt.Println("Database Migrated")
}
