package database

import (
	"context"
	"log"
	"time"

	"github.com/fatih/color"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func DBinstance() *mongo.Client {
	MongoDb := "mongodb+srv://<username>:<password>@mflix.0uncx.mongodb.net/test"
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
	color.Green("Succesfully Connected to MongoDB")
	return client
}

func OpenCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	var collection *mongo.Collection = client.Database("restaurant").Collection(collectionName)
	return collection
}
