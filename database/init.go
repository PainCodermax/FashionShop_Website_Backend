package database

import (
	"context"
    "fmt"
    "log"
    "os"
    "time"

    "github.com/joho/godotenv"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)


//DBinstance func
func DBinstance() *mongo.Client {
    err := godotenv.Load(".env")

    if err != nil {
        log.Fatal("Error loading .env file")
    }

    MongoDb := os.Getenv("URL")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

    client, err := mongo.Connect(ctx, options.Client().ApplyURI(MongoDb))
    if err != nil {
        log.Fatal(err)
    }


    defer cancel()
    err = client.Ping(context.Background(), nil)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("Connected to MongoDB!")

    return client
}

var Client *mongo.Client = DBinstance()

func GetCollection (client *mongo.Client, collectionName string) *mongo.Collection { 
	var collection *mongo.Collection = client.Database("Fashion_Shop").Collection(collectionName)
	return collection
}

