package database

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/fatih/color"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		color.Red("Failing in Loading Godotenv")
		log.Fatal("Error loading .env file")
	}
}

func DBinstance() *mongo.Client {
	MongoDb := os.Getenv("MONGO_DB_URI")
	client, err := mongo.NewClient(options.Client().ApplyURI(MongoDb))
	if err != nil {
		color.Red("MongoDB Connection Failed")
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		color.Red("MongoDB Connection Failed")
		log.Fatal(err)
	}
	color.Green("Successfully Connected to MongoDB")
	return client
}

func OpenCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	var collection *mongo.Collection = client.Database("restaurant").Collection(collectionName)
	return collection
}
