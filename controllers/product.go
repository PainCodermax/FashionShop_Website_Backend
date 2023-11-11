package controllers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/PainCodermax/FashionShop_Website_Backend/database"
	"github.com/PainCodermax/FashionShop_Website_Backend/models"
	"github.com/PainCodermax/FashionShop_Website_Backend/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var ProductCollection *mongo.Collection = database.ProductData(database.Client, "product")
var UserCollection *mongo.Collection = database.UserData(database.Client, "user")

func GetListProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		if value, ok := c.Get("isAdmin"); ok {
			if value == true {
				var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
				defer cancel()

				// limit := c.Query("limit")
				// offset := c.Query("query")
				// opt := options.FindOptions{
				// 	Limit: ,
				// 	offset: ,
				// }
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
						return
					}
					println(singleProduct.CategoryID)
					filter := bson.D{{"category_id", singleProduct.CategoryID}}
					category := make([]models.Category, 1)
					err := CategoryCollection.FindOne(ctx, filter).Decode(&category[0])
					if err != nil {
						c.JSON(http.StatusNotFound, gin.H{"error": "cannot found"})
						return
					}
					if len(category) > 0 {
						singleProduct.CategoryMame = utils.ParsePoitnerToString(category[0].Name)
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

				imgArr := make([]string, 0, len(products.ListImage))
				products.Product_ID = primitive.NewObjectID()
				cld, _ := utils.Credentials()
				for idx, img := range products.ListImage {
					imageString := utils.UploadImage(cld, img, idx, ctx)
					imgArr = append(imgArr, imageString)
				}
				products.ListImage = imgArr
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
		} else {
			c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
		}
	}
}

func UpdateProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		if value, ok := c.Get("isAdmin"); ok {
			if value == true {
				var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
				defer cancel()
				productId := c.Query("productId")
				if productId == "" {
					c.JSON(http.StatusNotFound, gin.H{"Error": "Wrong id not provided"})
					c.Abort()
					return
				}
				var editProduct models.Product
				if err := c.BindJSON(&editProduct); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"Error": "cannot format input"})
				}
				oid, err := primitive.ObjectIDFromHex(productId)
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
					return
				}
				filter := bson.D{primitive.E{Key: "_id", Value: oid}}
				update := bson.M{
					"$set": editProduct,
				}
				fmt.Println(update)
				result, err := ProductCollection.UpdateOne(ctx, filter, update)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
				c.JSON(http.StatusOK, result)
				return
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Cannot Update product"})
			}
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Cannot Update product"})
		}
	}
}

func DeleteProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		if value, ok := c.Get("isAdmin"); ok {
			if value == true {
				var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
				defer cancel()
				productId := c.Query("productId")
				if productId == "" {
					c.JSON(http.StatusNotFound, gin.H{"Error": "Wrong id not provided"})
					c.Abort()
					return
				}
				oid, err := primitive.ObjectIDFromHex(productId)
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
					return
				}
				filter := bson.D{primitive.E{Key: "_id", Value: oid}}
				result, err := ProductCollection.DeleteOne(ctx, filter)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
				c.JSON(http.StatusOK, result)
				return

			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Cannot Update product"})
			}
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Cannot Update product"})
		}
	}
}
