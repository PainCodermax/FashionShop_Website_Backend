package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/PainCodermax/FashionShop_Website_Backend/database"
	"github.com/PainCodermax/FashionShop_Website_Backend/email"
	"github.com/PainCodermax/FashionShop_Website_Backend/models"
	"github.com/PainCodermax/FashionShop_Website_Backend/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

		emailUser, ok := c.Get("email")
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Cannot get email"})
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
		orderID := utils.GenerateCode("ORD", 6)

		newOrder := models.Order{
			CartID:   checkout.CartID,
			UserID:   utils.InterfaceToString(userID),
			OrderID:  orderID,
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
		mailErr := email.ConfirmOrder(utils.InterfaceToString(emailUser), orderID)
		if mailErr != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "email not found"})
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

func CancelOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		orderID := c.Param("orderId")
		emailUser, ok := c.Get("email")
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Cannot get email"})
			return
		}
		update := bson.M{
			"$set": models.Order{
				Status: "CANCELED",
			},
		}

		rs, err := OrderCollection.UpdateOne(ctx, bson.D{{
			"order_id", orderID,
		}}, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.OrderResponse{
				Message: "cannot cancel this order",
			})
			return
		}
		mailErr := email.CancelOrder(utils.InterfaceToString(emailUser), orderID)
		if mailErr != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "email not found"})
			return
		}
		c.JSON(http.StatusOK, rs)
	}
}

func GetOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		orderID := c.Query("orderId")
		var order models.Order
		err := OrderCollection.FindOne(ctx, bson.D{{"order_id", orderID}}).Decode(&order)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.OrderResponse{
				Message: "cannot find this order",
			})
			return
		}
		c.JSON(http.StatusOK, models.OrderResponse{
			Status:  200,
			Message: "Get List product success",
			Data:    []models.Order{order},
		})
	}
}

func GetAllOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		limit, _ := utils.ParseStringToIn64(c.Query("limit"))
		offset, _ := utils.ParseStringToIn64(c.Query("offset"))

		if limit == 0 {
			limit = 20
		}
		if offset == 0 {
			offset = 0
		}

		opt := options.FindOptions{
			Limit: utils.ParseIn64ToPointer(limit),
			Skip:  utils.ParseIn64ToPointer(offset * limit),
		}

		rs, err := OrderCollection.Find(ctx, bson.M{}, &opt)
		var orders []models.Order
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Can Not get list"})
			return
		}
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
		totalCount, _ := OrderCollection.CountDocuments(ctx, bson.M{})
		c.JSON(http.StatusOK, models.OrderResponse{
			Status:  200,
			Message: "Get List product success",
			Data:    orders,
			Total:   int(totalCount),
		})

	}
}
