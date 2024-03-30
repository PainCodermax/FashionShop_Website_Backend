package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/PainCodermax/FashionShop_Website_Backend/database"
	"github.com/PainCodermax/FashionShop_Website_Backend/models"
	"github.com/PainCodermax/FashionShop_Website_Backend/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

var UserAddressCollection *mongo.Collection = database.ProductData(database.Client, "user_address")

func AddAdressUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var address models.UserAddress
		if userID, ok := c.Get("uid"); ok {
			if err := c.BindJSON(&address); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			address.UserID = utils.InterfaceToString(userID)
			_, addErr := UserAddressCollection.InsertOne(ctx, address)
			if addErr != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "not created"})
				return
			}

			c.JSON(http.StatusOK, "Successfully add address !")
		}
	}
}
