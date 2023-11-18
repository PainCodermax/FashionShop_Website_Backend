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
)

var CategoryCollection *mongo.Collection = database.ProductData(database.Client, "category")

func AddCategory() gin.HandlerFunc {
	return func(c *gin.Context) {
		if value, ok := c.Get("isAdmin"); ok {
			if value == true {
				var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
				var category models.Category
				defer cancel()
				if err := c.BindJSON(&category); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
					return
				}
				category.ID = primitive.NewObjectID()
				category.CategoryId = utils.GenerateCode("CATE", 5)

				_, anyerr := CategoryCollection.InsertOne(ctx, category)
				if anyerr != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Not Created"})
					return
				}
				c.JSON(http.StatusOK, "Successfully added our category Admin!!")
			}
			if value == false {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Not have authorization"})
				return
			}
		} else {
			c.JSON(http.StatusNotFound, gin.H{"error": "Cannot add category"})
		}
	}
}

func GetCategory() gin.HandlerFunc {
	return func(c *gin.Context) {
		if value, ok := c.Get("isAdmin"); ok {
			if value == true {
				var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
				defer cancel()
				var category []models.Category
				categoryID := c.Query("id")
				filter := bson.D{{"category_id", categoryID}}
				err := CategoryCollection.FindOne(ctx, filter).Decode(&category[0])
				if err != nil {
					c.JSON(http.StatusNotFound, gin.H{"error": "cannot found"})
					return
				}
				c.JSON(http.StatusOK, models.CategoryResponse{
					Status:  200,
					Message: "get category successfully",
					Data:    category,
				})
			}

		} else {
			c.JSON(http.StatusNotFound, gin.H{"error": "Cannot add category"})
			return
		}
	}
}

func GetCategoryList() gin.HandlerFunc {
	return func(c *gin.Context) {
		if value, ok := c.Get("isAdmin"); ok {
			if value == true {
				var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
				defer cancel()
				isMen, error := utils.ParseStringToBool(c.Query("isMen"))
				filter := bson.D{}
				if error == nil {
					filter = bson.D{{"is_men", isMen}}
				}
				var categoryList []models.Category
				result, err := CategoryCollection.Find(ctx, filter)
				if err != nil {
					c.JSON(http.StatusNotFound, gin.H{"message": "cannot get category list !!"})
					return
				}

				for result.Next(ctx) {
					singleCategory := models.Category{}
					if err := result.Decode(&singleCategory); err != nil {
						c.JSON(http.StatusInternalServerError, models.CategoryResponse{
							Status:  500,
							Message: "List product is empty",
							Data:    []models.Category{},
						})
					}
					categoryList = append(categoryList, singleCategory)
				}
				c.JSON(http.StatusOK, models.CategoryResponse{
					Status:  200,
					Message: "Get list category successfully",
					Data:    categoryList,
				})

			} else {
				c.JSON(http.StatusNotFound, gin.H{"error": "you don't have permission !"})
				return
			}
		}
	}
}

func UpdateCategory() gin.HandlerFunc {
	return func(c *gin.Context) {
		if value, ok := c.Get("isAdmin"); ok {
			if value == true {
				var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
				defer cancel()
				categoryId := c.Param("categoryId")
				var category models.Category
				if err := c.BindJSON(&category); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
					return
				}
				filter := bson.D{{"category_id", categoryId}}
				update := bson.M{"$set": category}

				result, err := CategoryCollection.UpdateOne(ctx, filter, update)

				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
				c.JSON(http.StatusOK, result)
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Cannot Update category"})
			}
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "you don't have permission"})
		}
	}
}

// func VerifyGmail() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
// 		defer cancel()
// 		var user models.User
// 		if err := c.BindJSON(&user); err != nil {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 			return
// 		}
		
// 	}
// }
