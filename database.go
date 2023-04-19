package main

import (
	"log"
	"os"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

//Database connection
func Connect() *gorm.DB {
	dsn := "user=" + os.Getenv("DB_USER") +
		" password=" + os.Getenv("DB_PASSWORD") +
		" host=" + os.Getenv("DB_URL") +
		" dbname=" + os.Getenv("DB_DATABASE") +
		" sslmode=disable"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Printf("Failed to connect to database. \n")
		os.Exit(100)
	}

	if err := createSchema(db); err != nil {
		log.Printf("Failed to create schema: %v\n", err)
	}

	log.Printf("Database connection successful! \n")

	return db
}


func createSchema(db *gorm.DB) error {
	models := []interface{}{
		&File{},
	}

	for _, model := range models {
		err := db.AutoMigrate(model)
		if err != nil {
			return err
		}
	}

	return nil
}