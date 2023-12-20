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

		err := CartCollection.FindOne(ctx, bson.M{"user_id": userID}).Decode(&cart)
		if err == nil {
			var cartSearchItem models.CartItem
			errItem := CartItemCollection.FindOne(ctx, bson.M{"product_id": cartItem.ProductID, "cart_id": cart.CartID}).Decode(&cartSearchItem)
			if errItem != nil {
				cartItem.CartID = cart.CartID
				cartItem.CartItemID = utils.GenerateCode("CARTITEM", 9)
				_, inserterr := CartItemCollection.InsertOne(ctx, cartItem)
				if inserterr != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "not created"})
					return
				}
				c.JSON(http.StatusCreated, gin.H{
					"message": "Successfully add to card!!",
				})
			} else {
				cartItem.Price = cartSearchItem.Price + cartItem.Price*cartItem.ItemQuantity
				cartItem.ItemQuantity = cartSearchItem.ItemQuantity + cartItem.ItemQuantity
				filter := bson.D{{"cart_item_id", cartSearchItem.CartItemID}}
				update := bson.M{
					"$set": cartItem,
				}

				_, err := CartItemCollection.UpdateOne(ctx, filter, update)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "not created"})
					return
				}
				c.JSON(http.StatusCreated, gin.H{
					"message": "Successfully add to card!!",
				})
			}
		} else {
			if err, cartID := createCart(ctx, utils.InterfaceToString(userID)); err != nil {
				cartItem.CartID = cartID
				cartItem.Price = cartItem.Price * cartItem.ItemQuantity
				cartItem.CartItemID = utils.GenerateCode("CARTITEM", 9)
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

func GetCart() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		userID, ok := c.Get("uid")
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Cannot get userID"})
			return
		}
		var cart models.Cart
		var items []models.CartItem
		err := CartCollection.FindOne(ctx, bson.M{"user_id": userID}).Decode(&cart)

		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"message": "cart not found !"})
			return
		}

		filter := bson.D{{"cart_id", cart.CartID}}
		rs, findErr := CartItemCollection.Find(ctx, filter)
		if findErr != nil {
			c.JSON(http.StatusNotFound, gin.H{"message": "Cannot get userID"})
			return
		}
		total := 0
		totalPrice := 0
		for rs.Next(ctx) {
			var singleCartItem models.CartItem
			if err := rs.Decode(&singleCartItem); err != nil {
				c.JSON(http.StatusOK, models.CartItemResponse{
					Status:  200,
					Message: "Cart is empty",
					Data:    []models.CartItem{},
				})
				return
			}
			total = total + singleCartItem.ItemQuantity
			totalPrice = totalPrice + singleCartItem.Price

			items = append(items, singleCartItem)
		}
		c.JSON(http.StatusOK, models.CartItemResponse{
			Status:     200,
			Message:    "Get cart successfully",
			Total:      total,
			TotalPrice: totalPrice,
			Data:       items,
		})

	}
}

func DeleteCartItem() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		cartItemID := c.Param("cartItemID")
		if cartItemID == "" {
			c.JSON(http.StatusNotFound, gin.H{"Error": "Wrong id not provided"})
			c.Abort()
			return
		}
		filter := bson.D{{"cart_item_id", cartItemID}}
		result, err := CartItemCollection.DeleteOne(ctx, filter)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, result)
	}
}

func UpdateCart() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var cartItem models.CartItem
		if err := c.BindJSON(&cartItem); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		cartItem.Price = cartItem.Price * cartItem.ItemQuantity

		filter := bson.D{{"cart_item_id", cartItem.CartItemID}}
		update := bson.M{
			"$set": cartItem,
		}
		_, err := CartItemCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "not created"})
			return
		}
		c.JSON(http.StatusCreated, gin.H{
			"message": "Successfully add to card!!",
		})
	}
}
