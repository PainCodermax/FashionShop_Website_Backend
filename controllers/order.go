package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/PainCodermax/FashionShop_Website_Backend/database"
	"github.com/PainCodermax/FashionShop_Website_Backend/models"
	"github.com/PainCodermax/FashionShop_Website_Backend/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var OrderCollection *mongo.Collection = database.ProductData(database.Client, "order")

func Checkout() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		userID, ok := c.Get("uid")
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Cannot get userID"})
			return
		}

		var checkout models.RequestOrder
		if err := c.BindJSON(&checkout); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		filter := bson.D{{"cart_id", checkout.CartID}}
		// _, err := CartCollection.DeleteOne(ctx, filter)
		// if err != nil {
		// 	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		// 	return
		// }

		_, err := CartItemCollection.DeleteMany(ctx, filter, nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		for _, cartItem := range checkout.CartItems {
			var product models.Product
			updateFilter := bson.D{{"product_id", cartItem.ProductID}}
			e := ProductCollection.FindOne(ctx, updateFilter).Decode(&product)
			if e != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": e.Error()})
				return
			}
			quantity := *product.Quantity - cartItem.ItemQuantity
			productUpdate := models.Product{
				Quantity: &quantity,
			}
			_, uErr := ProductCollection.UpdateOne(ctx, updateFilter, bson.M{
				"$set": productUpdate,
			})
			if uErr != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}

		newOrder := models.Order{
			CartID:   checkout.CartID,
			UserID:   utils.InterfaceToString(userID),
			OrderID:  utils.GenerateCode("ORD", 6),
			Items:    checkout.CartItems,
			Price:    checkout.TotalPrice,
			Status:   "SUCCESS",
			Address:  checkout.Address,
			Quantity: checkout.Quantity,
		}

		_, anyerr := OrderCollection.InsertOne(ctx, newOrder)
		if anyerr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Not Created"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Checkout successfully"})
	}
}

func GetListOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		userID, ok := c.Get("uid")
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Cannot get userID"})
			return
		}

		rs, err := OrderCollection.Find(ctx, bson.D{{
			"user_id", userID,
		}})

		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "cannot found"})
			return
		}
		var orders []models.Order
		for rs.Next(ctx) {
			order := models.Order{}
			if err := rs.Decode(&order); err != nil {
				c.JSON(http.StatusInternalServerError, models.OrderResponse{
					Status:  500,
					Message: "List order is empty",
					Data:    []models.Order{},
				})
				return
			}
			orders = append(orders, order)
		}
		c.JSON(http.StatusOK, models.OrderResponse{
			Status:  200,
			Message: "Get List product success",
			Data:    orders,
		})

	}
}
