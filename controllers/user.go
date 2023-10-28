package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/PainCodermax/FashionShop_Website_Backend/models"
	generate "github.com/PainCodermax/FashionShop_Website_Backend/tokens"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

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
