package controllers

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type Application struct {
	prodCollection        *mongo.Collection
	userCollection        *mongo.Collection
	categoryCollection    *mongo.Collection
	addressUserCollection *mongo.Collection
}

func NewApplication(prodCollection, userCollection, categoryCollection *mongo.Collection) *Application {
	return &Application{
		prodCollection:     prodCollection,
		userCollection:     userCollection,
		categoryCollection: CategoryCollection,
	}
}
