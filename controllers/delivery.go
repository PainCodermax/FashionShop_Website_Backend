package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/PainCodermax/FashionShop_Website_Backend/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func GetDelivery() gin.HandlerFunc {
	return func(c *gin.Context) {
		if value, ok := c.Get("isAdmin"); ok {
			if value == true {
				var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
				defer cancel()
				var delivery models.Delivery
				oderID := c.Param("orderID")
				filter := bson.D{{"order_id", oderID}}
				err := DeliveryCollection.FindOne(ctx, filter).Decode(&delivery)
				if err != nil {
					c.JSON(http.StatusNotFound, gin.H{"error": "cannot found"})
					return
				}
				c.JSON(http.StatusOK, models.DeliveryResponse{
					Status:  200,
					Message: "get category successfully",
					Data:    delivery,
				})
			}

		} else {
			c.JSON(http.StatusNotFound, gin.H{"error": "Cannot add category"})
			return
		}
	}
}
