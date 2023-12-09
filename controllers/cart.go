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
	"go.mongodb.org/mongo-driver/mongo"
)

var CartCollection *mongo.Collection = database.DB(database.Client, "cart")
var CartItemCollection *mongo.Collection = database.DB(database.Client, "cart_item")

func AddToCart() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		userID, ok := c.Get("uid")
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Cannot get userID"})
			return
		}

		var cartItem models.CartItem
		var cart models.Cart

		if err := c.BindJSON(&cartItem); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}	

		err := CartCollection.FindOne(ctx, bson.M{"user_id": userID}).Decode(cart)
		if err == nil {
			
		} else {
			if err, cartID := createCart(ctx, utils.InterfaceToString(userID)); err != nil {
				cartItem.CartID = cartID
				_, inserterr := CartItemCollection.InsertOne(ctx, cartItem)
				if inserterr != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "not created"})
					return
				}
				c.JSON(http.StatusCreated, gin.H{
					"message": "Successfully add to card!!",
				})
			}
		}
	}
}

func createCart(ctx context.Context, userID string) (error, string) {
	cartID := utils.GenerateCode("CART", 5)
	cart := models.Cart{
		UserID: userID,
		CartID: cartID,
	}

	_, err := CartCollection.InsertOne(ctx, cart)
	if err != nil {
		return err, ""
	}
	return nil, cartID
}

// func (app *Application) AddToCart() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		productQueryID := c.Query("id")
// 		if productQueryID == "" {
// 			log.Println("product id is empty")
// 			_ = c.AbortWithError(http.StatusBadRequest, errors.New("product id is empty"))
// 			return
// 		}
// 		userQueryID := c.Query("userID")
// 		if userQueryID == "" {
// 			log.Println("user id is empty")
// 			_ = c.AbortWithError(http.StatusBadRequest, errors.New("user id is empty"))
// 			return
// 		}
// 		productID, err := primitive.ObjectIDFromHex(productQueryID)
// 		if err != nil {
// 			log.Println(err)
// 			c.AbortWithStatus(http.StatusInternalServerError)
// 			return
// 		}
// 		var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
// 		defer cancel()

// 		err = database.AddProductToCart(ctx, app.prodCollection, app.userCollection, productID, userQueryID)
// 		if err != nil {
// 			c.IndentedJSON(http.StatusInternalServerError, err)
// 		}
// 		c.IndentedJSON(200, "Successfully Added to the cart")
// 	}
// }

// func (app *Application) RemoveItem() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		productQueryID := c.Query("id")
// 		if productQueryID == "" {
// 			log.Println("product id is inavalid")
// 			_ = c.AbortWithError(http.StatusBadRequest, errors.New("product id is empty"))
// 			return
// 		}

// 		userQueryID := c.Query("userID")
// 		if userQueryID == "" {
// 			log.Println("user id is empty")
// 			_ = c.AbortWithError(http.StatusBadRequest, errors.New("UserID is empty"))
// 		}

// 		ProductID, err := primitive.ObjectIDFromHex(productQueryID)
// 		if err != nil {
// 			log.Println(err)
// 			c.AbortWithStatus(http.StatusInternalServerError)
// 			return
// 		}

// 		var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
// 		defer cancel()
// 		err = database.RemoveCartItem(ctx, app.prodCollection, app.userCollection, ProductID, userQueryID)
// 		if err != nil {
// 			c.IndentedJSON(http.StatusInternalServerError, err)
// 			return
// 		}
// 		c.IndentedJSON(200, "Successfully removed from cart")
// 	}
// }

// func GetItemFromCart() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		user_id := c.Query("id")
// 		if user_id == "" {
// 			c.Header("Content-Type", "application/json")
// 			c.JSON(http.StatusNotFound, gin.H{"error": "invalid id"})
// 			c.Abort()
// 			return
// 		}

// 		usert_id, _ := primitive.ObjectIDFromHex(user_id)

// 		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
// 		defer cancel()

// 		var filledcart models.User
// 		err := UserCollection.FindOne(ctx, bson.D{primitive.E{Key: "_id", Value: usert_id}}).Decode(&filledcart)
// 		if err != nil {
// 			log.Println(err)
// 			c.IndentedJSON(500, "not id found")
// 			return
// 		}

// 		filter_match := bson.D{{Key: "$match", Value: bson.D{primitive.E{Key: "_id", Value: usert_id}}}}
// 		unwind := bson.D{{Key: "$unwind", Value: bson.D{primitive.E{Key: "path", Value: "$usercart"}}}}
// 		grouping := bson.D{{Key: "$group", Value: bson.D{primitive.E{Key: "_id", Value: "$_id"}, {Key: "total", Value: bson.D{primitive.E{Key: "$sum", Value: "$usercart.price"}}}}}}
// 		pointcursor, err := UserCollection.Aggregate(ctx, mongo.Pipeline{filter_match, unwind, grouping})
// 		if err != nil {
// 			log.Println(err)
// 		}
// 		var listing []bson.M
// 		if err = pointcursor.All(ctx, &listing); err != nil {
// 			log.Println(err)
// 			c.AbortWithStatus(http.StatusInternalServerError)
// 		}
// 		for _, json := range listing {
// 			c.IndentedJSON(200, json["total"])
// 			// c.IndentedJSON(200, filledcart.UserCart)
// 		}
// 		ctx.Done()
// 	}
// }

// func (app *Application) BuyFromCart() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		userQueryID := c.Query("id")
// 		if userQueryID == "" {
// 			log.Panicln("user id is empty")
// 			_ = c.AbortWithError(http.StatusBadRequest, errors.New("UserID is empty"))
// 		}
// 		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
// 		defer cancel()
// 		err := database.BuyItemFromCart(ctx, app.userCollection, userQueryID)
// 		if err != nil {
// 			c.IndentedJSON(http.StatusInternalServerError, err)
// 		}
// 		c.IndentedJSON(200, "Successfully Placed the order")
// 	}
// }

// func (app *Application) InstantBuy() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		UserQueryID := c.Query("userid")
// 		if UserQueryID == "" {
// 			log.Println("UserID is empty")
// 			_ = c.AbortWithError(http.StatusBadRequest, errors.New("UserID is empty"))
// 		}
// 		ProductQueryID := c.Query("pid")
// 		if ProductQueryID == "" {
// 			log.Println("Product_ID id is empty")
// 			_ = c.AbortWithError(http.StatusBadRequest, errors.New("product_id is empty"))
// 		}
// 		productID, err := primitive.ObjectIDFromHex(ProductQueryID)
// 		if err != nil {
// 			log.Println(err)
// 			c.AbortWithStatus(http.StatusInternalServerError)
// 			return
// 		}

// 		var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
// 		defer cancel()
// 		err = database.InstantBuyer(ctx, app.prodCollection, app.userCollection, productID, UserQueryID)
// 		if err != nil {
// 			c.IndentedJSON(http.StatusInternalServerError, err)
// 		}
// 		c.IndentedJSON(200, "Successully placed the order")
// 	}
// }
