package config

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Database

func ConnectDB() {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal("MongoDB bağlantı hatası:", err)
	}

	DB = client.Database("havuzMavisi")

	fmt.Println("✅ MongoDB bağlantısı kuruldu.")
}

func GetCollection(collectionName string) *mongo.Collection {
	return DB.Collection(collectionName)
}