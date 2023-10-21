package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/PainCodermax/FashionShop_Website_Backend/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetListProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		if value, ok := c.Get("isAdmin"); ok {
			if value == true {
				var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
				defer cancel()
				result, err := ProductCollection.Find(ctx, bson.M{})
				var listProduct []models.Product
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Can Not Get List"})
					return
				}

				for result.Next(ctx) {
					singleProduct := models.Product{}
					if err := result.Decode(&singleProduct); err != nil {
						c.JSON(http.StatusInternalServerError, models.ProductResponse{
							Status:  500,
							Message: "List product is empty",
							Data:    []models.Product{},
						})
						println(err.Error())
						return
					}
					listProduct = append(listProduct, singleProduct)
				}
				c.JSON(http.StatusOK, models.ProductResponse{
					Status:  200,
					Message: "Get List product success",
					Data:    listProduct,
				})
			}
			if value == false {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Not have authorization"})
				return
			}
		}
	}
}

func AddProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		if value, ok := c.Get("isAdmin"); ok {
			if value == true {
				var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
				var products models.Product
				defer cancel()
				if err := c.BindJSON(&products); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
					return
				}
				products.Product_ID = primitive.NewObjectID()
				_, anyerr := ProductCollection.InsertOne(ctx, products)
				if anyerr != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Not Created"})
					return
				}
				defer cancel()
				c.JSON(http.StatusOK, "Successfully added our Product Admin!!")
			}
			if value == false {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Not have authorization"})
				return
			}
		}
	}
}
