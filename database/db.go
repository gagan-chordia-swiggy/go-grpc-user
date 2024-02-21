package database

import (
	"fmt"
	"go-project/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

func DatabaseConnection() *gorm.DB {
	host := "localhost"
	port := 5432
	user := "postgres"
	password := "root1234"
	dbName := "postgres"

	connectionString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbName)

	db, err := gorm.Open(postgres.Open(connectionString), &gorm.Config{})

	if err != nil {
		log.Fatalln("Error connecting to the database")
	}

	log.Print("Connected to the database")

	err = db.AutoMigrate(&models.User{})

	if err != nil {
		log.Fatalln("Error migrating the database")
	}

	return db
}

func GetUserByUsername(db *gorm.DB, username string) (*models.User, error) {
	var user models.User
	err := db.Where("username = ?", username).Find(&user).Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}
