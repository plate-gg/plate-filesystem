package main

import (
    "context"
    "log"
    "os"
	"time"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

func connect() (*mongo.Client, *mongo.Database) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGODB_URI")))

	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	log.Printf("Database connection successful! \n")
	return client, client.Database(os.Getenv("MONGODB_DATABASE_Name"))
}

