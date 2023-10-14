package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/PainCodermax/FashionShop_Website_Backend/database"
	"github.com/PainCodermax/FashionShop_Website_Backend/model"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// get collection
var adminCollection *mongo.Collection = database.GetCollection(database.Client, "admin")

func LoginAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var Admin model.Admin
		var AdminFound model.Admin

		if err := c.BindJSON(&Admin); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err := adminCollection.FindOne(ctx, bson.M{"username": Admin.UserName, "password": Admin.Password}).Decode(&AdminFound)
		defer cancel()

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "username or password is not correct !!"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"data": AdminFound, "message": "success"})
	}
}
