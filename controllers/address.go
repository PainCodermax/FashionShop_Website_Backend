package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/PainCodermax/FashionShop_Website_Backend/client"
	"github.com/PainCodermax/FashionShop_Website_Backend/database"
	"github.com/PainCodermax/FashionShop_Website_Backend/models"
	"github.com/PainCodermax/FashionShop_Website_Backend/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
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
			address.AddressID = utils.GenerateCode("ADD", 6)
			address.FullAddress = address.Street + ", " + client.GetAddressString(address.ProvinceID, address.DistrictID, address.WardID)
			_, addErr := UserAddressCollection.InsertOne(ctx, address)
			if addErr != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "not created"})
				return
			}
			c.JSON(http.StatusOK, "Successfully add address !")
		}
	}
}

func GetAddressUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var address models.UserAddress
		addressId := c.Query("addressId")

		filter := bson.M{"address_id": addressId}
		err := OrderCollection.FindOne(ctx, filter).Decode(&address)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.OrderResponse{
				Message: "cannot find this order",
			})
			return
		}
		c.JSON(http.StatusOK, address)
	}
}

func GetAddressUserList() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		if userID, ok := c.Get("uid"); ok {
			filter := bson.D{{"user_id", userID}}

			rs, addErr := UserAddressCollection.Find(ctx, filter)
			if addErr != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "not found"})
				return
			}
			var address []models.UserAddress
			for rs.Next(ctx) {
				var singleAddress models.UserAddress
				if err := rs.Decode(&singleAddress); err != nil {
					c.JSON(http.StatusInternalServerError, models.AddressResponse{
						Status:  500,
						Message: "List product is empty",
						Data:    []models.UserAddress{},
					})
				}
				address = append(address, singleAddress)
			}
			c.JSON(http.StatusOK, models.AddressResponse{
				Status:  200,
				Message: "get product list success",
				Data:    address,
				Total:   len(address),
			})
		}
	}
}
