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
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var WishItemCollection *mongo.Collection = database.DB(database.Client, "wish_item")

func LikeItem() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		userID, ok := c.Get("uid")
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Cannot get userID"})
			return
		}
		var wishItem models.WishItem
		if err := c.BindJSON(&wishItem); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var founduser *models.User
		err := UserCollection.FindOne(ctx, bson.M{"user_id": userID}).Decode(&founduser)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error":  "not found user",
				"status": 404,
			})
			return
		}
		wishItem.UserId = utils.InterfaceToString(userID)
		var newWish models.WishItem
		if err := WishItemCollection.FindOne(ctx, bson.M{
			"product_id": wishItem.ProductId,
			"user_id":    userID,
		}).Decode(&newWish); err == nil {
			c.JSON(http.StatusOK, "Successfully add wishlist!!")
			return
		}
		wishItem.WishItemId = utils.GenerateCode("WISH", 7)
		if wishItem.ProductId != "" {
			f := bson.M{"user_id": userID, "product_id": wishItem.ProductId}
			errExist := WishItemCollection.FindOne(ctx, f)
			if errExist.Err() == nil {
				c.JSON(http.StatusNotFound, gin.H{
					"error":  "item already exist",
					"status": 400,
				})
				return
			}
			var foundProduct models.Product
			filter := bson.D{{"product_id", wishItem.ProductId}}
			err = ProductCollection.FindOne(ctx, filter).Decode(&foundProduct)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			wishItem.Product = foundProduct
		}
		wishItem.ID = primitive.NewObjectID()
		wishItem.Created_At = time.Now()
		_, anyerr := WishItemCollection.InsertOne(ctx, wishItem)
		if anyerr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Not Created"})
			return
		}
		defer cancel()
		c.JSON(http.StatusOK, "Successfully add wishlist!!")
	}
}

func UnLikeItem() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		product_id := c.Param("id")
		userID, ok := c.Get("uid")
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Cannot get userID"})
			return
		}
		if product_id == "" {
			c.JSON(http.StatusNotFound, gin.H{"Error": "Wrong id not provided"})
			c.Abort()
			return
		}
		filter := bson.M{
			"product_id": product_id,
			"user_id":    userID,
		}
		_, err := WishItemCollection.DeleteOne(ctx, filter)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, "delete wish item complete")
	}
}

func GetWishList() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		userID, ok := c.Get("uid")
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Cannot get userID"})
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
		var listProduct []models.Product
		skip := utils.ParseIn64ToPointer(offset * limit)
		opt := options.FindOptions{
			Limit: utils.ParseIn64ToPointer(limit),
			Skip:  skip,
			Sort:  bson.D{{Key: "created_at", Value: -1}},
		}
		totalCount, _ := ProductCollection.CountDocuments(ctx, bson.D{{"user_id", userID}})

		result, err := WishItemCollection.Find(ctx, bson.D{{"user_id", userID}}, &opt)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Can Not Get List"})
			return
		}
		for result.Next(ctx) {
			wishItem := models.WishItem{}
			if err := result.Decode(&wishItem); err != nil {
				c.JSON(http.StatusBadRequest, models.ProductResponse{
					Status:  400,
					Message: "List product is empty",
					Data:    []models.Product{},
				})
				return
			}
			listProduct = append(listProduct, wishItem.Product)
		}
		c.JSON(http.StatusOK, models.ProductResponse{
			Status:  200,
			Message: "Get List product success",
			Data:    listProduct,
			Total:   int(totalCount),
		})
	}
}
