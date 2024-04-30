package controllers

import (
	"context"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strconv"
	"strings"
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

var DeliveryCollection *mongo.Collection = database.ProductData(database.Client, "delivery")

func Checkout(orderUpdateChannel chan<- string) gin.HandlerFunc {
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
			CartID:        checkout.CartID,
			UserID:        utils.InterfaceToString(userID),
			OrderID:       orderID,
			Items:         checkout.CartItems,
			Price:         checkout.TotalPrice + checkout.ShipFee,
			DileveryDate:  t,
			Status:        enum.Pending,
			Address:       checkout.Address,
			Quantity:      checkout.Quantity,
			Created_At:    time.Now(),
			PaymentMethod: checkout.PaymentMethod,
		}

		_, anyerr := OrderCollection.InsertOne(ctx, newOrder)
		if anyerr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Not Created"})
			return
		}

		// orderUpdateChannel <- orderID
		time.AfterFunc(10*time.Second, func() {
			UpdateOrder(orderID)
		})

		mailErr := email.ConfirmOrder(utils.InterfaceToString(emailUser), orderID)
		if mailErr != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "email not found"})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "Checkout successfully",
			"orderId": orderID,
		})
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
				Status: enum.Cancelled,
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

		userID, ok := c.Get("uid")
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Cannot get userID"})
			return
		}

		var order models.Order
		err := OrderCollection.FindOne(ctx, bson.D{{"order_id", orderID}}).Decode(&order)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.OrderResponse{
				Message: "cannot find this order",
			})
			return
		}
		var rawItems []models.CartItem
		for _, item := range order.Items {
			var rating *models.Rating

			err := RatingCollection.FindOne(ctx, bson.M{
				"user_id":    userID,
				"product_id": item.ProductID,
			}).Decode(&rating)
			if err != nil {
				rawItems = append(rawItems, item)
				continue
			}
			item.IsRate = true

			rawItems = append(rawItems, item)
		}
		order.Items = rawItems

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

		var wg sync.WaitGroup
		wg.Add(4)

		go func() {
			defer wg.Done() // Đánh dấu hoàn thành goroutine khi kết thúc
			update := bson.M{
				"$set": models.Order{
					Status: enum.Submitted,
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

			go func() {
				defer wg.Done() // Đánh dấu hoàn thành goroutine khi kết thúc
				mailErr := email.CancelOrder(utils.ParsePoitnerToString(user.Email), orderID)
				if mailErr != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
			}()
		}()

		go func() {
			defer wg.Done()
			delivery := models.Delivery{
				ID:             primitive.NewObjectID(),
				DeliveryID:     utils.GenerateCode("DElI", 5),
				OrderID:        &orderID,
				DeliveryDate:   time.Now().Add(7 * 24 * time.Hour), // Thêm 7 ngày vào ngày hiện tại
				Created_At:     time.Now(),
				Address:        order.Address,
				DeliveryStatus: enum.Shipping,
			}

			_, err := DeliveryCollection.InsertOne(ctx, delivery)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}()

		wg.Wait()
		c.JSON(http.StatusOK, gin.H{"message": "order submission in progress"})
	}
}

func PaymentByVnPay() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		orderID := c.PostForm("orderId")
		amount, _ := strconv.Atoi(c.PostForm("amount"))
		var order models.Order

		filter := bson.D{{"order_id", orderID}}
		err := OrderCollection.FindOne(ctx, filter).Decode(&order)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if order.PaymentMethod != enum.VNPAY {
			c.JSON(
				http.StatusBadRequest,
				gin.H{
					"error": "Order is COD",
				},
			)
		}
		// Thông tin từ VNPay
		vnpUrl := os.Getenv("VNP_URL")

		// Thông tin merchant từ tài khoản của bạn
		merchantId := os.Getenv("MERCHANT_ID")

		// Tạo request
		data := url.Values{}
		data.Set("vnp_Version", "2.1.0")
		data.Set("vnp_Command", "pay")
		data.Set("vnp_TmnCode", merchantId)
		data.Set("vnp_Locale", "vn")
		data.Set("vnp_CurrCode", "VND")
		data.Set("vnp_TxnRef", orderID)
		data.Set("vnp_OrderInfo", "Thanh toan don hang")
		data.Set("vnp_OrderType", "billpayment")
		data.Set("vnp_Amount", fmt.Sprintf("%d", amount*100)) // Chuyển đổi sang đơn vị xu
		data.Set("vnp_ReturnUrl", "http://localhost:3000/users/payment/vnpay/callback")
		data.Set("vnp_IpAddr", "127.0.0.1")

		// Tính mã hash
		secureHash := "SERCURE_HASH"
		data.Set("vnp_SecureHashType", "MD5")
		query := data.Encode() + "&vnp_SecureHash=" + generateMD5Hash(data.Encode()+secureHash)

		// Trả về URL thanh toán VNPay
		c.JSON(http.StatusOK, gin.H{
			"payment_url": vnpUrl + "?" + query,
		})
	}
}

func generateMD5Hash(data string) string {
	hash := md5.New()
	hash.Write([]byte(data))
	return hex.EncodeToString(hash.Sum(nil))
}

func PaymentByVnPayCallBack() gin.HandlerFunc {
	return func(c *gin.Context) {
		var _, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		orderId := c.PostForm("vnp_TxnRef")
		amount, _ := strconv.Atoi(c.PostForm("vnp_Amount"))
		status := c.PostForm("vnp_ResponseCode")

		if status == "00" {
			// Cập nhật trạng thái đơn hàng là đã thanh toán thành công trong cơ sở dữ liệu của bạn
			// Ví dụ: updateOrderStatus(orderId, "paid")

			// Trả về trang cảm ơn
			c.JSON(http.StatusOK, gin.H{
				"order_id": orderId,
				"amount":   amount,
			})
		} else {
			// Xử lý trường hợp thanh toán không thành công
			c.String(http.StatusOK, "Payment failed")
		}
	}
}

func PaymentByVnPay2() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		loc, _ := time.LoadLocation("Asia/Ho_Chi_Minh")
		date := time.Now().In(loc)
		createDate := date.Format("20060102150405")

		ipAddr := c.Request.Header.Get("X-Forwarded-For")
		if ipAddr == "" {
			ipAddr = c.Request.RemoteAddr
		}

		tmnCode := os.Getenv("MERCHANT_ID")
		secretKey := os.Getenv("SERCURE_HASH")
		vnpURL := os.Getenv("VNP_URL")
		returnURL := os.Getenv("VNP_RETURN_URL")

		var payment models.RequestPayment
		if err := c.BindJSON(&payment); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		bankCode := c.PostForm("bankCode")

		var order models.Order

		filter := bson.D{{"order_id", payment.OrderId}}
		err := OrderCollection.FindOne(ctx, filter).Decode(&order)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if order.PaymentMethod != enum.VNPAY {
			c.JSON(
				http.StatusBadRequest,
				gin.H{
					"error": "Order is COD",
				},
			)
		}

		locale := c.PostForm("language")
		if locale == "" {
			locale = "vn"
		}
		currCode := "VND"

		vnpParams := url.Values{}
		vnpParams.Set("vnp_Version", "2.1.0")
		vnpParams.Set("vnp_Command", "pay")
		vnpParams.Set("vnp_TmnCode", tmnCode)
		vnpParams.Set("vnp_Locale", locale)
		vnpParams.Set("vnp_CurrCode", currCode)
		vnpParams.Set("vnp_TxnRef", payment.OrderId)
		vnpParams.Set("vnp_OrderInfo", "Thanh toan cho ma GD:"+payment.OrderId)
		vnpParams.Set("vnp_OrderType", "other")
		vnpParams.Set("vnp_Amount", utils.ParseIntToString(order.Price*100))
		vnpParams.Set("vnp_ReturnUrl", returnURL)
		vnpParams.Set("vnp_IpAddr", ipAddr)
		vnpParams.Set("vnp_CreateDate", createDate)
		if bankCode != "" {
			vnpParams.Set("vnp_BankCode", bankCode)
		}

		signData := strings.Join([]string{vnpParams.Encode()}, "&")
		mac := hmac.New(sha512.New, []byte(secretKey))
		mac.Write([]byte(signData))
		signature := fmt.Sprintf("%x", mac.Sum(nil))

		vnpParams.Set("vnp_SecureHash", signature)

		vnpURL += "?" + vnpParams.Encode()

		c.JSON(http.StatusOK, gin.H{
			"message": "get vnpay payment success",
			"url":     vnpURL,
		})

	}
}

func VnpayReturnHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		vnpParams := c.Request.URL.Query()
		secureHash := vnpParams.Get("vnp_SecureHash")

		delete(vnpParams, "vnp_SecureHash")
		delete(vnpParams, "vnp_SecureHashType")

		sortedParams := sortParams(vnpParams)

		secretKey := os.Getenv("SERCURE_HASH")

		signData := url.Values{}
		for _, param := range sortedParams {
			signData.Add(param[0], param[1])
		}

		hash := hmac.New(sha512.New, []byte(secretKey))
		hash.Write([]byte(signData.Encode()))
		signed := hex.EncodeToString(hash.Sum(nil))

		if secureHash == signed {
			orderId := vnpParams["vnp_TxnRef"]
			update := bson.M{
				"$set": models.Order{
					IsPaid: true, // Simplified the bool expression
				},
			}
			_, uErr := OrderCollection.UpdateOne(ctx, bson.D{{"order_id", orderId[0]}}, update)
			if uErr != nil {
				log.Println("UpdateOne error:", uErr)
				c.JSON(http.StatusInternalServerError, gin.H{"error": uErr.Error()})
				return
			}
			// workerChannel <- orderId
			c.Redirect(http.StatusSeeOther, "https://fashion-shop-client-379b8.web.app/thankyou")
			return
		}

		c.String(http.StatusOK, "Failed! Response Code: 97")
	}
}

func sortParams(params url.Values) [][2]string {
	var sortedParams [][2]string
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		for _, v := range params[k] {
			sortedParams = append(sortedParams, [2]string{k, v})
		}
	}
	return sortedParams
}

func UpdateOrder(orderId string) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()
	var status enum.OrderStatus = enum.Submitted

	var order models.Order
	filter := bson.D{{"order_id", orderId}}
	e := OrderCollection.FindOne(ctx, filter).Decode(&order)
	if e != nil {
		log.Println("UpdateOne error:", e)
		return
	}
	update := bson.M{
		"$set": models.Order{
			Status: status, // Simplified the bool expression
		},
	}

	if order.Status == enum.Pending {
		if order.PaymentMethod == enum.VNPAY && !order.IsPaid {
			status = enum.Cancelled
		}
		_, uErr := OrderCollection.UpdateOne(ctx, bson.D{{"order_id", orderId}}, update)
		if uErr != nil {
			log.Println("UpdateOne error:", uErr)
			return
		}
	}
}

func UpdateOrderStatus(order models.Order) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var status enum.OrderStatus = enum.Received

	if order.Status == enum.Submitted {
		status = enum.Processing
	}

	if order.Status == enum.Processing {
		status = enum.Delivery
	}

	if order.Status == enum.Delivery {
		status = enum.Shipping
	}

	update := bson.M{
		"$set": models.Order{
			Status: status,
		},
	}

	if order.Status == enum.Delivery {
		status = enum.Received
		update = bson.M{
			"$set": models.Order{
				Status: status,
				IsPaid: true,
			},
		}
	}

	_, uErr := OrderCollection.UpdateOne(ctx, bson.D{{"order_id", order.OrderID}}, update)
	if uErr != nil {
		log.Println("UpdateOne error:", uErr)
		return
	}
}
