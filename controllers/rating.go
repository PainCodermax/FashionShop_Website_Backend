package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/PainCodermax/FashionShop_Website_Backend/models"
	"github.com/PainCodermax/FashionShop_Website_Backend/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func CreateRating() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		userID, ok := c.Get("uid")
		var rating models.Rating
		if err := c.BindJSON(&rating); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var founduser *models.User
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Cannot get userID"})
			return
		}
		err := UserCollection.FindOne(ctx, bson.M{"user_id": userID}).Decode(&founduser)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error":  "not found user",
				"status": 404,
			})
			return
		}
		rating.User = founduser

		if rating.OrderID != nil {
			var order models.Order
			filter := bson.D{{"order_id", rating.OrderID}}
			err = OrderCollection.FindOne(ctx, filter).Decode(&order)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			rating.Order = &order
		}

		if rating.Product_ID != "" {
			var foundProduct models.Product
			filter := bson.D{{"product_id", rating.Product_ID}}
			err = ProductCollection.FindOne(ctx, filter).Decode(&foundProduct)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			rating.Product = &foundProduct
		}
		rating.ID = primitive.NewObjectID()

		_, anyerr := RatingCollection.InsertOne(ctx, rating)
		if anyerr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Not Created"})
			return
		}
		defer cancel()
		c.JSON(http.StatusOK, "Successfully add our rating!!")
	}

}

func GetRating() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		productId := c.Param("productId")
		if productId == "" {
			c.JSON(http.StatusNotFound, gin.H{"Error": "Wrong id not provided"})
			return
		}

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
		result, err := RatingCollection.Find(ctx, bson.D{{"product_id", productId}}, &opt)
		var listRating []models.Rating
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Can Not Get List"})
			return
		}
		totalCount, _ := RatingCollection.CountDocuments(ctx, bson.M{})
		for result.Next(ctx) {
			singleRating := models.Rating{}
			if err := result.Decode(&singleRating); err != nil {
				c.JSON(http.StatusInternalServerError, models.ProductResponse{
					Status:  500,
					Message: "List product is empty",
					Data:    []models.Product{},
				})
				return
			}
			listRating = append(listRating, singleRating)
		}
		c.JSON(http.StatusOK, models.RatingResponse{
			Status:  200,
			Message: "Get List product success",
			Data:    listRating,
			Total:   int(totalCount),
		})
	}
}
