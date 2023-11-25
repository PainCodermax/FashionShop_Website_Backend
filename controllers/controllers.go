package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/PainCodermax/FashionShop_Website_Backend/models"
	generate "github.com/PainCodermax/FashionShop_Website_Backend/tokens"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

var Validate = validator.New()

func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}
	return string(bytes)
}

func VerifyPassword(userpassword string, givenpassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(givenpassword), []byte(userpassword))
	valid := true
	msg := ""
	if err != nil {
		msg = "Login Or Passowrd is Incorerct"
		valid = false
	}
	return valid, msg
}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var user models.User
		var founduser models.User
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}
		err := UserCollection.FindOne(ctx,bson.M{"email": user.Email}).Decode(&founduser)
		defer cancel() 
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":  "login or password incorrect",
				"status": 401,
			})
			return
		}
		PasswordIsValid, msg := VerifyPassword(*user.Password, *founduser.Password)
		defer cancel()
		if !PasswordIsValid {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":  msg,
				"status": 401})
			fmt.Println(msg)
			return
		}
		token, refreshToken, _ := generate.TokenGenerator(*founduser.Email, *founduser.First_Name, *founduser.Last_Name, founduser.User_ID, founduser.IsAdmin)
		defer cancel()
		generate.UpdateAllTokens(token, refreshToken, founduser.User_ID)
		founduser.Refresh_Token = &refreshToken
		c.JSON(http.StatusOK, founduser)
	}
}

// func ProductViewerAdmin() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
// 		var products models.Product
// 		defer cancel()
// 		if err := c.BindJSON(&products); err != nil {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 			return
// 		}
// 		products.Product_ID = primitive.NewObjectID()
// 		_, anyerr := ProductCollection.InsertOne(ctx, products)
// 		if anyerr != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "Not Created"})
// 			return
// 		}
// 		defer cancel()
// 		c.JSON(http.StatusOK, "Successfully added our Product Admin!!")
// 	}
// }

// func SearchProduct() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		var productlist []models.Product
// 		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
// 		defer cancel()
// 		cursor, err := ProductCollection.Find(ctx, bson.D{{}})
// 		if err != nil {
// 			c.IndentedJSON(http.StatusInternalServerError, "Someting Went Wrong Please Try After Some Time")
// 			return
// 		}
// 		err = cursor.All(ctx, &productlist)
// 		if err != nil {
// 			log.Println(err)
// 			c.AbortWithStatus(http.StatusInternalServerError)
// 			return
// 		}
// 		defer cursor.Close(ctx)
// 		if err := cursor.Err(); err != nil {
// 			// Don't forget to log errors. I log them really simple here just
// 			// to get the point across.
// 			log.Println(err)
// 			c.IndentedJSON(400, "invalid")
// 			return
// 		}
// 		defer cancel()
// 		c.IndentedJSON(200, productlist)

// 	}
// }

// func SearchProductByQuery() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		var searchproducts []models.Product
// 		queryParam := c.Query("name")
// 		if queryParam == "" {
// 			log.Println("query is empty")
// 			c.Header("Content-Type", "application/json")
// 			c.JSON(http.StatusNotFound, gin.H{"Error": "Invalid Search Index"})
// 			c.Abort()
// 			return
// 		}
// 		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
// 		defer cancel()
// 		searchquerydb, err := ProductCollection.Find(ctx, bson.M{"product_name": bson.M{"$regex": queryParam}})
// 		if err != nil {
// 			c.IndentedJSON(404, "something went wrong in fetching the dbquery")
// 			return
// 		}
// 		err = searchquerydb.All(ctx, &searchproducts)
// 		if err != nil {
// 			log.Println(err)
// 			c.IndentedJSON(400, "invalid")
// 			return
// 		}
// 		defer searchquerydb.Close(ctx)
// 		if err := searchquerydb.Err(); err != nil {
// 			log.Println(err)
// 			c.IndentedJSON(400, "invalid request")
// 			return
// 		}
// 		defer cancel()
// 		c.IndentedJSON(200, searchproducts)
// 	}
// }
