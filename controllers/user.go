package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/PainCodermax/FashionShop_Website_Backend/email"
	"github.com/PainCodermax/FashionShop_Website_Backend/models"
	"github.com/PainCodermax/FashionShop_Website_Backend/models/query"
	generate "github.com/PainCodermax/FashionShop_Website_Backend/tokens"
	"github.com/PainCodermax/FashionShop_Website_Backend/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
		err := UserCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&founduser)
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

func GetNewToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var user models.User
		var founduser models.User
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}
		err := UserCollection.FindOne(ctx, bson.M{"refresh_token": *user.Refresh_Token}).Decode(&founduser)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "can not find"})
			return
		}
		token, _ := generate.AccessTokenGenerator(*founduser.Email, *founduser.First_Name, *founduser.Last_Name, founduser.User_ID, founduser.IsAdmin)
		defer cancel()
		generate.UpdateAccessToken(token, founduser.User_ID)
		founduser.Token = &token
		c.JSON(http.StatusOK, founduser)
	}
}

func SignUp() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var user models.User
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		validationErr := Validate.Struct(user)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr})
			return
		}

		count, err := UserCollection.CountDocuments(ctx, bson.M{"email": user.Email})
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
		if count > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User already exists"})
		}
		count, err = UserCollection.CountDocuments(ctx, bson.M{"phone": user.Phone})
		defer cancel()
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
		if count > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Phone is already in use"})
			return
		}
		password := HashPassword(*user.Password)
		user.Password = &password

		user.UserCode = utils.GenerateCode("USER", 5)
		user.Created_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.Updated_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.ID = primitive.NewObjectID()
		user.User_ID = user.ID.Hex()
		token, refreshtoken, _ := generate.TokenGenerator(*user.Email, *user.First_Name, *user.Last_Name, user.User_ID, user.IsAdmin)
		user.Token = &token
		user.Refresh_Token = &refreshtoken
		user.Address_Details = make([]models.Address, 0)
		user.Order_Status = make([]models.Order, 0)
		user.VerifyCode = utils.GenerateCode("VRF", 6)

		mailErr := email.SendOPTMail(*user.Email, user.VerifyCode, true)
		if mailErr != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "email not found"})
			return
		}
		_, inserterr := UserCollection.InsertOne(ctx, user)
		if inserterr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "not created"})
			return
		}
		defer cancel()
		c.JSON(http.StatusCreated, gin.H{
			"message": "Successfully Signed Up!!",
			"data":    user,
		})
	}
}

func VerifyUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var user models.User
		var request query.VerifyRequest
		if err := c.BindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		fifter := bson.M{"user_code": request.User_Code}
		err := UserCollection.FindOne(ctx, fifter).Decode(&user)
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}

		if request.VerifyCode == user.VerifyCode {
			user.IsVerified = true
			userUpdate := bson.M{
				"$set": user,
			}

			_, err := UserCollection.UpdateOne(ctx, fifter, userUpdate)
			if err == nil {
				c.JSON(http.StatusOK, gin.H{"message": "verify successfully"})
				return
			}
		}

		c.JSON(http.StatusInternalServerError, gin.H{"message": "Cannot verify"})
	}
}

func ForGotPassword() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var user models.User
		var foundUser models.User
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}
		err := UserCollection.FindOne(ctx, bson.M{"email": *user.Email}).Decode(&foundUser)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "can not find"})
			return
		}
		mailErr := email.SendOPTMail(*user.Email, user.VerifyCode, false)
		if mailErr != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "email not found"})
			return
		}
		defer cancel()
		c.JSON(http.StatusOK, gin.H{"message": "verify successfully"})
	}
}

func UpdatePassWord() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var user models.User
		var foundUser models.User
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}

		err := UserCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&foundUser)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "can not find"})
			return
		}
		filter := bson.D{{"email", user.Email}}
		update := bson.M{"$set": user}
		result, err := UserCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, result)
	}
}
