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
