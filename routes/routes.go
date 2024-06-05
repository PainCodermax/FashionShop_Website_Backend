package routes

import (
	"context"

	"github.com/PainCodermax/FashionShop_Website_Backend/controllers"
	"github.com/PainCodermax/FashionShop_Website_Backend/worker"

	"github.com/gin-gonic/gin"
)

func LoginRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.GET("users/payment/vnpay/callback", controllers.VnpayReturnHandler())
	incomingRoutes.POST("/users/signup", controllers.SignUp())
	incomingRoutes.POST("/users/login", controllers.Login())
	incomingRoutes.POST("/users/verify", controllers.VerifyUser())
	incomingRoutes.POST("/users/forget-password", controllers.ForGotPassword())
	incomingRoutes.PUT("/users/update-password", controllers.UpdatePassWord())
}

func UserRoutes(incomingRoutes *gin.Engine) {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	workerChannel := make(chan string)

	WorkerChannel := make(chan string)
	go worker.Worker(ctx, WorkerChannel)
	//users
	incomingRoutes.POST("/users/get-token", controllers.GetNewToken())
	incomingRoutes.GET("/users/get-user", controllers.GetUser())
	incomingRoutes.PUT("/users/update-user", controllers.UpdateUser())

	// product
	incomingRoutes.GET("/users/product/list", controllers.GetListProduct())
	incomingRoutes.GET("/users/product/:productId", controllers.GetProduct())
	incomingRoutes.GET("/users/product", controllers.SearchProduct())
	incomingRoutes.GET("/users/product-by-category", controllers.GetProductByCategory())
	incomingRoutes.GET("/users/product/recommend/:productId", controllers.GetRecommendList())

	//category
	incomingRoutes.GET("/users/get-category/list", controllers.GetCategoryList())

	//cart
	incomingRoutes.POST("/users/cart/add", controllers.AddToCart())
	incomingRoutes.GET("/users/cart", controllers.GetCart())
	incomingRoutes.PUT("/users/cart/update", controllers.UpdateCart())
	incomingRoutes.DELETE("/users/cart/delete/:cartItemID", controllers.DeleteCartItem())

	//order
	incomingRoutes.POST("/users/checkout", controllers.Checkout(workerChannel))
	incomingRoutes.GET("/users/orders", controllers.GetListOrder())
	incomingRoutes.PUT("/users/orders/cancel/:orderId", controllers.CancelOrder())
	incomingRoutes.GET("/users/orders/get-single", controllers.GetOrder())
	incomingRoutes.GET("/users/orders/get-raw-order", controllers.GetRawOrder())
	incomingRoutes.PUT("/users/order/update", controllers.GetOneOrder())
	//address
	// incomingRoutes.POST("/users/address", controllers.AddAddress())

	//rating
	incomingRoutes.POST("/users/rating", controllers.CreateRating())
	incomingRoutes.GET("/users/rating/:productId", controllers.GetRating())

	incomingRoutes.POST("users/payment/vnpay", controllers.PaymentByVnPay2())

	//address
	incomingRoutes.POST("/users/address/add", controllers.AddAdressUser())
	incomingRoutes.GET("/users/address/list", controllers.GetAddressUserList())

}

func AdminRouter(incomingRoutes *gin.Engine) {
	//user
	incomingRoutes.GET("/amdin/user/list", controllers.GetUserList())
	incomingRoutes.GET("/admin/user/:userId", controllers.GetSingleUser())

	//product
	incomingRoutes.POST("/admin/product/add", controllers.AddProduct())
	incomingRoutes.GET("/admin/product/list", controllers.GetListProduct())
	incomingRoutes.PUT("/admin/product/update/:productId", controllers.UpdateProduct())
	incomingRoutes.DELETE("/admin/product/delete/:productId", controllers.DeleteProduct())

	//category
	incomingRoutes.POST("/admin/add-category", controllers.AddCategory())
	incomingRoutes.GET("/admin/get-category", controllers.GetCategory())
	incomingRoutes.GET("/admin/get-category/list", controllers.GetCategoryList())
	incomingRoutes.PUT("/admin/update-category/:categoryId", controllers.UpdateCategory())

	//order
	incomingRoutes.GET("/admin/orders/get-all-order", controllers.GetAllOrder())
	incomingRoutes.GET("/admin/order/:orderID", controllers.GetOneOrder())
	incomingRoutes.PUT("/admin/order/confirm/:orderID", controllers.SubmitOrder())

	//delivery
	incomingRoutes.GET("/admin/order/delivery/:orderID", controllers.GetDelivery())

	//report
	incomingRoutes.GET("/admin/report", controllers.GetReport())

	//flashSale
	incomingRoutes.GET("/admin/flashSale/list", controllers.GetReport())
	incomingRoutes.POST("/admin/flashSale/get", controllers.GetReport())
}
