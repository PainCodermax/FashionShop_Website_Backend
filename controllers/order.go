package controllers

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/PainCodermax/FashionShop_Website_Backend/client"
	"github.com/PainCodermax/FashionShop_Website_Backend/database"
	"github.com/PainCodermax/FashionShop_Website_Backend/email"
	"github.com/PainCodermax/FashionShop_Website_Backend/enum"
	"github.com/PainCodermax/FashionShop_Website_Backend/models"
	"github.com/PainCodermax/FashionShop_Website_Backend/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var OrderCollection *mongo.Collection = database.ProductData(database.Client, "order")
var DeliveryCollection *mongo.Collection = database.ProductData(database.Client, "delivery")

func Checkout() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		userID, ok := c.Get("uid")
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Cannot get userID"})
			return
		}

		emailUser, ok := c.Get("email")
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Cannot get email"})
			return
		}
		var checkout models.RequestOrder
		if err := c.BindJSON(&checkout); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		filter := bson.D{{"cart_id", checkout.CartID}}
		// _, err := CartCollection.DeleteOne(ctx, filter)
		// if err != nil {
		// 	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		// 	return
		// }

		_, err := CartItemCollection.DeleteMany(ctx, filter, nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		for _, cartItem := range checkout.CartItems {
			var product models.Product
			updateFilter := bson.D{{"product_id", cartItem.ProductID}}
			e := ProductCollection.FindOne(ctx, updateFilter).Decode(&product)
			if e != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": e.Error()})
				return
			}
			quantity := *product.Quantity - cartItem.ItemQuantity
			productUpdate := models.Product{
				Quantity: &quantity,
			}
			_, uErr := ProductCollection.UpdateOne(ctx, updateFilter, bson.M{
				"$set": productUpdate,
			})
			if uErr != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}
		t := time.Now().Add(7 * 24 * time.Hour)

		orderID := utils.GenerateCode("ORD", 6)

		newOrder := models.Order{
			CartID:       checkout.CartID,
			UserID:       utils.InterfaceToString(userID),
			OrderID:      orderID,
			Items:        checkout.CartItems,
			Price:        checkout.TotalPrice + checkout.ShipFee,
			DileveryDate: t,
			Status:       "SUCCESS",
			Address:      checkout.Address,
			Quantity:     checkout.Quantity,
			Created_At:   time.Now(),
		}

		_, anyerr := OrderCollection.InsertOne(ctx, newOrder)
		if anyerr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Not Created"})
			return
		}
		mailErr := email.ConfirmOrder(utils.InterfaceToString(emailUser), orderID)
		if mailErr != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "email not found"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Checkout successfully"})
	}
}

func GetListOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		userID, ok := c.Get("uid")
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Cannot get userID"})
			return
		}

		rs, err := OrderCollection.Find(ctx, bson.D{{
			"user_id", userID,
		}})

		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "cannot found"})
			return
		}
		var orders []models.Order
		for rs.Next(ctx) {
			order := models.Order{}
			if err := rs.Decode(&order); err != nil {
				c.JSON(http.StatusInternalServerError, models.OrderResponse{
					Status:  500,
					Message: "List order is empty",
					Data:    []models.Order{},
				})
				return
			}
			orders = append(orders, order)
		}
		c.JSON(http.StatusOK, models.OrderResponse{
			Status:  200,
			Message: "Get List product success",
			Data:    orders,
		})

	}
}

func CancelOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		orderID := c.Param("orderId")
		emailUser, ok := c.Get("email")
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Cannot get email"})
			return
		}

		update := bson.M{
			"$set": models.Order{
				Status: string(enum.Cancelled),
			},
		}

		rs, err := OrderCollection.UpdateOne(ctx, bson.D{{
			"order_id", orderID,
		}}, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.OrderResponse{
				Message: "cannot cancel this order",
			})
			return
		}
		mailErr := email.CancelOrder(utils.InterfaceToString(emailUser), orderID)
		if mailErr != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "email not found"})
			return
		}
		c.JSON(http.StatusOK, rs)
	}
}

func GetOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		orderID := c.Query("orderId")
		var order models.Order
		err := OrderCollection.FindOne(ctx, bson.D{{"order_id", orderID}}).Decode(&order)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.OrderResponse{
				Message: "cannot find this order",
			})
			return
		}
		c.JSON(http.StatusOK, models.OrderResponse{
			Status:  200,
			Message: "Get List product success",
			Data:    []models.Order{order},
		})
	}
}

func GetAllOrder() gin.HandlerFunc {
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

		opt := options.FindOptions{
			Limit: utils.ParseIn64ToPointer(limit),
			Skip:  utils.ParseIn64ToPointer(offset * limit),
		}

		rs, err := OrderCollection.Find(ctx, bson.M{}, &opt)
		var orders []models.Order
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Can Not get list"})
			return
		}
		for rs.Next(ctx) {
			order := models.Order{}
			if err := rs.Decode(&order); err != nil {
				c.JSON(http.StatusInternalServerError, models.OrderResponse{
					Status:  500,
					Message: "List order is empty",
					Data:    []models.Order{},
				})
				return
			}
			orders = append(orders, order)
		}
		totalCount, _ := OrderCollection.CountDocuments(ctx, bson.M{})
		c.JSON(http.StatusOK, models.OrderResponse{
			Status:  200,
			Message: "Get List product success",
			Data:    orders,
			Total:   int(totalCount),
		})

	}
}

func GetRawOrder() gin.HandlerFunc {
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
		var user models.User
		err := CartCollection.FindOne(ctx, bson.M{"user_id": userID}).Decode(&cart)

		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"message": "cart not found !"})
			return
		}

		userErr := UserCollection.FindOne(ctx, bson.M{"user_id": userID}).Decode(&user)
		if userErr != nil {
			c.JSON(http.StatusNotFound, gin.H{"message": "user not found !"})
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
		fee := client.CheckShipingFee("", user.District, user.Ward)

		c.JSON(http.StatusOK, models.CartItemResponse{
			Status:     200,
			Message:    "Get cart successfully",
			Total:      total,
			TotalPrice: totalPrice,
			Data:       items,
			ShipFee:    fee,
			Province:   user.Province,
			District:   user.District,
			Ward:       user.Ward,
		})

	}
}

func GetOneOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		if value, ok := c.Get("isAdmin"); ok {
			if value != true {
				c.JSON(http.StatusNotFound, gin.H{"error": "no permision"})
				return
			}
		}
		var order models.Order
		orderID := c.Param("orderID")
		if orderID == "" {
			c.JSON(http.StatusNotFound, gin.H{"Error": "Wrong id not provided"})
			return
		}

		filter := bson.D{{"order_id", orderID}}
		err := OrderCollection.FindOne(ctx, filter).Decode(&order)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, order)
	}
}

func SubmitOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		if value, ok := c.Get("isAdmin"); ok {
			if value != true {
				c.JSON(http.StatusNotFound, gin.H{"error": "no permission"})
				return
			}
		}

		orderID := c.Param("orderID")
		var order models.Order

		filter := bson.D{{"order_id", orderID}}
		err := OrderCollection.FindOne(ctx, filter).Decode(&order)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		var wg sync.WaitGroup // Khởi tạo WaitGroup

		// Thêm số lượng goroutine cần đợi vào WaitGroup
		wg.Add(4)

		// Thực hiện cập nhật trạng thái đơn hàng trong một goroutine riêng
		go func() {
			defer wg.Done() // Đánh dấu hoàn thành goroutine khi kết thúc
			update := bson.M{
				"$set": models.Order{
					Status: string(enum.Processing),
				},
			}
			_, err := OrderCollection.UpdateOne(ctx, bson.D{{"order_id", orderID}}, update)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}()

		// Lấy thông tin user trong một goroutine riêng
		go func() {
			defer wg.Done() // Đánh dấu hoàn thành goroutine khi kết thúc
			var user models.User
			err := UserCollection.FindOne(ctx, bson.D{{"user_id", order.UserID}}).Decode(&user)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			// Gửi email trong một goroutine riêng
			go func() {
				defer wg.Done() // Đánh dấu hoàn thành goroutine khi kết thúc
				mailErr := email.CancelOrder(utils.ParsePoitnerToString(user.Email), orderID)
				if mailErr != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
			}()
		}()

		// Thêm dữ liệu vào collection Delivery trong một goroutine riêng
		go func() {
			defer wg.Done() // Đánh dấu hoàn thành goroutine khi kết thúc
			delivery := models.Delivery{
				ID:             primitive.NewObjectID(),
				DeliveryID:     utils.GenerateCode("DElI", 5),
				OrderID:        &orderID,
				DeliveryDate:   time.Now().Add(7 * 24 * time.Hour), // Thêm 7 ngày vào ngày hiện tại
				Created_At:     time.Now(),
				DeliveryStatus: enum.Shipping,
			}

			_, err := DeliveryCollection.InsertOne(ctx, delivery)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}()

		// Chờ cho tất cả các goroutine hoàn thành trước khi trả về phản hồi cho client
		wg.Wait()
		c.JSON(http.StatusOK, gin.H{"message": "order submission in progress"})
	}
}
