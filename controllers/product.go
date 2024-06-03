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

var ProductCollection *mongo.Collection = database.ProductData(database.Client, "product")
var RecommendationCollection *mongo.Collection = database.DB(database.Client, "recommendation")

func GetListProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		// if value, ok := c.Get("isAdmin"); ok {
		// 	if value == true {
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
		var listProduct []models.Product
		skip := utils.ParseIn64ToPointer(offset * limit)
		opt := options.FindOptions{
			Limit: utils.ParseIn64ToPointer(limit),
			Skip:  skip,
		}
		result, err := ProductCollection.Find(ctx, bson.M{}, &opt)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Can Not Get List"})
			return
		}
		totalCount, _ := ProductCollection.CountDocuments(ctx, bson.M{})
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
			Total:   int(totalCount),
		})
	}
	// 	if value == false {
	// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not have authorization"})
	// 		return
	// 	}
	// }
	// }
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
				products.Product_ID = utils.GenerateCode("PRO", 5)
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
				productId := c.Param("productId")
				if productId == "" {
					c.JSON(http.StatusNotFound, gin.H{"Error": "Wrong id not provided"})
					return
				}
				var editProduct models.Product
				if err := c.BindJSON(&editProduct); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"Error": "cannot format input"})
				}
				filter := bson.D{{"product_id", productId}}
				update := bson.M{
					"$set": editProduct,
				}
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
				productId := c.Param("productId")
				if productId == "" {
					c.JSON(http.StatusNotFound, gin.H{"Error": "Wrong id not provided"})
					c.Abort()
					return
				}
				filter := bson.D{{"product_id", productId}}
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

func GetProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var foundProduct models.Product
		productId := c.Param("productId")
		if productId == "" {
			c.JSON(http.StatusNotFound, gin.H{"Error": "Wrong id not provided"})
			c.Abort()
			return
		}
		filter := bson.D{{"product_id", productId}}
		err := ProductCollection.FindOne(ctx, filter).Decode(&foundProduct)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		cateFilter := bson.D{{"category_id", foundProduct.CategoryID}}
		category := make([]models.Category, 1)
		cateErr := CategoryCollection.FindOne(ctx, cateFilter).Decode(&category[0])
		if cateErr != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "cannot found"})
			return
		}
		if len(category) > 0 {
			foundProduct.CategoryMame = utils.ParsePoitnerToString(category[0].Name)
			foundProduct.Gender = category[0].Gender
		}
		c.JSON(http.StatusOK, foundProduct)
	}
}

func SearchProduct() gin.HandlerFunc {
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
		text := c.Query("text")

		opt := options.FindOptions{
			Limit: utils.ParseIn64ToPointer(limit),
			Skip:  utils.ParseIn64ToPointer(offset * limit),
		}

		filter := bson.M{}
		if text != "" {
			filter["product_name"] = bson.M{"$regex": primitive.Regex{Pattern: text, Options: "i"}}
		}

		result, err := ProductCollection.Find(ctx, filter, &opt)
		var listProduct []models.Product
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Can Not Get List"})
			return
		}
		totalCount, _ := ProductCollection.CountDocuments(ctx, bson.M{})
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
			Message: "search products success",
			Data:    listProduct,
			Total:   int(totalCount),
		})
	}
}

func GetProductByCategory() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		categoryID := c.Query("categoryId")
		limit, _ := utils.ParseStringToIn64(c.Query("limit"))
		offset, _ := utils.ParseStringToIn64(c.Query("offset"))
		if categoryID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"Error": "missing category id"})
			return
		}
		opt := options.FindOptions{
			Limit: utils.ParseIn64ToPointer(limit),
			Skip:  utils.ParseIn64ToPointer(offset * limit),
		}
		filter := bson.D{{"category_id", categoryID}}

		result, err := ProductCollection.Find(ctx, filter, &opt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Can Not Get List"})
			return
		}
		var listProduct []models.Product
		totalCount, _ := ProductCollection.CountDocuments(ctx, filter)
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
			Total:   int(totalCount),
		})
	}
}

func GetRecommendList() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
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
		skip := offset * limit

		opt := options.FindOptions{
			Limit: utils.ParseIn64ToPointer(limit * 5),
			Skip:  utils.ParseIn64ToPointer(skip),
			Sort:  bson.D{{Key: "weight", Value: 1}},
		}
		filter := bson.D{{"product_id_1", productId}}

		var listProductId []string
		result, err := RecommendationCollection.Find(ctx, filter, &opt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Can Not Get List"})
			return
		}
		defer result.Close(ctx)

		for result.Next(ctx) {
			singleRecommend := models.Recommendation{}
			if err := result.Decode(&singleRecommend); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error decoding data"})
				return
			}
			listProductId = append(listProductId, singleRecommend.Product_id_2)
		}

		if len(listProductId) == 0 {
			c.JSON(http.StatusNotFound, gin.H{"Message": "not found product"})
			return
		}

		// Lấy danh sách sản phẩm từ ProductCollection
		productOpt := options.FindOptions{
			Limit: utils.ParseIn64ToPointer(limit),
			Skip:  utils.ParseIn64ToPointer(skip),
		}
		products, err := ProductCollection.Find(ctx, bson.M{"product_id": bson.M{"$in": listProductId}}, &productOpt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Can Not Get List"})
			return
		}
		defer products.Close(ctx)

		var listProduct []models.Product
		for products.Next(ctx) {
			singleProduct := models.Product{}
			if err := products.Decode(&singleProduct); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error decoding data"})
				return
			}
			listProduct = append(listProduct, singleProduct)
		}

		totalCount, _ := RecommendationCollection.CountDocuments(ctx, filter)

		c.JSON(http.StatusOK, models.ProductResponse{
			Status:  200,
			Message: "Get List product success",
			Data:    listProduct,
			Total:   int(totalCount),
		})
	}
}
