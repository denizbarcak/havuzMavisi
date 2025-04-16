package config

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Database

func ConnectDB() {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}

	// Ping the database
	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	DB = client.Database("havuzMavisi")

	// Create unique index for email field in users collection
	usersCollection := DB.Collection("users")
	_, err = usersCollection.Indexes().CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys: bson.D{{Key: "email", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
	)
	if err != nil {
		log.Fatal("Error creating email index:", err)
	}

	fmt.Println("✅ MongoDB bağlantısı kuruldu.")
}

func GetCollection(collectionName string) *mongo.Collection {
	return DB.Collection(collectionName)
}