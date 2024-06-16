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

var FlashSaleCollection *mongo.Collection = database.ProductData(database.Client, "flash_sale")

func AddFlashSale() gin.HandlerFunc {
	return func(c *gin.Context) {
		if value, ok := c.Get("isAdmin"); ok {
			if value == true {
				var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
				var flashSale models.FlashSale
				defer cancel()
				if err := c.BindJSON(&flashSale); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
					return
				}
				flashSale.ID = primitive.NewObjectID()
				flashSale.FlashSaleId = utils.GenerateCode("FLASH", 7)
				flashSale.Created_At = time.Now()

				_, anyerr := FlashSaleCollection.InsertOne(ctx, flashSale)
				if anyerr != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Not Created"})
					return
				}
				defer cancel()
				c.JSON(http.StatusOK, "Successfully added our flash sale Admin!!")
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

func GetFlashSale() gin.HandlerFunc {
	return func(c *gin.Context) {
		if value, ok := c.Get("isAdmin"); ok {
			if value == true {
				var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
				var flashSale models.FlashSale
				defer cancel()
				flashSaleId := c.Param("flashSaleId")
				filter := bson.D{{"flash_sale_id", flashSaleId}}
				anyerr := FlashSaleCollection.FindOne(ctx, filter).Decode(&flashSale)
				if anyerr != nil {
					c.JSON(http.StatusNotFound, gin.H{"error": "Not Found"})
					return
				}

				filter = bson.D{
					{"product_id", bson.D{
						{"$in", flashSale.ProductIdList},
					}},
				}

				// filter = bson.D{{"product_id", bson.D{{"$in", flashSale.ProductIdList}}}}
				rs, err := ProductCollection.Find(ctx, filter)
				if err != nil {
					c.JSON(http.StatusNotFound, gin.H{"error": "Not Found"})
					return
				}
				var productList []models.Product

				if err = rs.All(ctx, &productList); err != nil {
					c.JSON(http.StatusNotFound, gin.H{"error": "Not Found"})
					return
				}

				for i := range productList {
					productList[i].FlashSalePrice = *productList[i].Price - (*productList[i].Price/100)*flashSale.Discount
				}
				flashSale.ProductList = productList

				defer cancel()
				c.JSON(http.StatusOK, models.FlashSalehResponse{
					Status:  200,
					Message: "get flashsale successfully",
					Data:    []models.FlashSale{flashSale},
				})
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

func GetFlashSales() gin.HandlerFunc {
	return func(c *gin.Context) {
		if value, ok := c.Get("isAdmin"); ok {
			if value == true {
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
					Sort:  bson.D{{Key: "created_at", Value: -1}},
				}
				result, err := FlashSaleCollection.Find(ctx, bson.D{{}}, &opt)
				var flashSaleList []models.FlashSale
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Can Not Get List"})
					return
				}
				totalCount, _ := FlashSaleCollection.CountDocuments(ctx, bson.M{})
				for result.Next(ctx) {
					singleFlashSale := models.FlashSale{}
					if err := result.Decode(&singleFlashSale); err != nil {
						c.JSON(http.StatusInternalServerError, models.ProductResponse{
							Status:  500,
							Message: "List flash sale is empty",
							Data:    []models.Product{},
						})
						return
					}
					flashSaleList = append(flashSaleList, singleFlashSale)
				}
				c.JSON(http.StatusOK, models.FlashSalehResponse{
					Status:  200,
					Message: "Get List flash sale success",
					Data:    flashSaleList,
					Total:   int(totalCount),
				})
				if value == false {
					c.JSON(http.StatusUnauthorized, gin.H{"error": "Not have authorization"})
					return
				}
			} else {
				c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
			}
		}
	}
}

func UpdateFlashSale() gin.HandlerFunc {
	return func(c *gin.Context) {
		if value, ok := c.Get("isAdmin"); ok {
			if value == true {
				var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
				var flashSale models.FlashSale
				var newFlashSale models.FlashSale
				defer cancel()
				flashSaleId := c.Param("flashSaleId")

				if err := c.BindJSON(&flashSale); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
					return
				}

				filter := bson.D{{"flash_sale_id", flashSaleId}}

				anyerr := FlashSaleCollection.FindOne(ctx, filter).Decode(&newFlashSale)
				if anyerr != nil {
					c.JSON(http.StatusNotFound, gin.H{"error": "Not Found"})
					return
				}

				if flashSale.Discount > 0 {
					newFlashSale.Discount = flashSale.Discount
				}

				if len(flashSale.ProductIdList) > 0 {
					newFlashSale.ProductIdList = flashSale.ProductIdList
				}

				if flashSale.TimeExpired != nil {
					newFlashSale.TimeExpired = flashSale.TimeExpired
				}

				if flashSale.TimeStarted != nil {
					newFlashSale.TimeStarted = flashSale.TimeStarted
				}

				newFlashSale.Updated_At = time.Now()
				update := bson.M{
					"$set": newFlashSale,
				}
				_, err := FlashSaleCollection.UpdateOne(ctx, filter, update)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}

				defer cancel()
				c.JSON(http.StatusOK, models.FlashSalehResponse{
					Status:  200,
					Message: "update flashsale successfully",
					Data:    []models.FlashSale{newFlashSale},
				})
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
