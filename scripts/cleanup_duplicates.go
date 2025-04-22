package main

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// MongoDB'ye bağlan
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27018"))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(context.Background())

	db := client.Database("havuzMavisi")
	usersCollection := db.Collection("users")

	// Tüm kullanıcıları email'e göre grupla
	pipeline := []bson.M{
		{
			"$group": bson.M{
				"_id": "$email",
				"users": bson.M{
					"$push": bson.M{
						"id":        "$_id",
						"name":      "$name",
						"email":     "$email",
						"password":  "$password",
						"role":      "$role",
						"createdAt": "$createdAt",
					},
				},
			},
		},
		{
			"$match": bson.M{
				"users.1": bson.M{"$exists": true}, // Sadece birden fazla kaydı olanları al
			},
		},
	}

	cursor, err := usersCollection.Aggregate(context.Background(), pipeline)
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(context.Background())

	var results []bson.M
	if err = cursor.All(context.Background(), &results); err != nil {
		log.Fatal(err)
	}

	// Her duplicate grup için en son kaydı tut, diğerlerini sil
	for _, result := range results {
		email := result["_id"].(string)
		users := result["users"].(bson.A)

		// En son kaydı bul
		var latestUser bson.M
		var latestTime time.Time
		for _, user := range users {
			userMap := user.(bson.M)
			createdAt := userMap["createdAt"].(time.Time)
			if createdAt.After(latestTime) {
				latestTime = createdAt
				latestUser = userMap
			}
		}

		// En son kayıt dışındaki tüm kayıtları sil
		_, err = usersCollection.DeleteMany(
			context.Background(),
			bson.M{
				"email": email,
				"_id": bson.M{
					"$ne": latestUser["id"],
				},
			},
		)
		if err != nil {
			log.Printf("Error deleting duplicates for email %s: %v\n", email, err)
			continue
		}

		log.Printf("Cleaned up duplicates for email: %s\n", email)
	}

	log.Println("Duplicate cleanup completed!")
} 