package routes

import (
	"github.com/PainCodermax/FashionShop_Website_Backend/controllers"

	"github.com/gin-gonic/gin"
)

func LoginRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("/users/signup", controllers.SignUp())
	incomingRoutes.POST("/users/login", controllers.Login())
	incomingRoutes.POST("/users/verify", controllers.VerifyUser())
	incomingRoutes.POST("/users/forget-password", controllers.ForGotPassword())
	incomingRoutes.PUT("/users/update-password", controllers.UpdatePassWord())
}

func UserRoutes(incomingRoutes *gin.Engine) {
	//users
	incomingRoutes.POST("/users/get-token", controllers.GetNewToken())
	incomingRoutes.GET("/users/get-user", controllers.GetUser())
	incomingRoutes.PUT("/users/update-user", controllers.UpdateUser())

	// product
	incomingRoutes.GET("/users/product/list", controllers.GetListProduct())
	incomingRoutes.GET("/users/product/:productId", controllers.GetProduct())
	incomingRoutes.GET("/users/product", controllers.SearchProduct())
	incomingRoutes.GET("/users/product-by-category", controllers.GetProductByCategory())

	//category
	incomingRoutes.GET("/users/get-category/list", controllers.GetCategoryList())

	//cart
	incomingRoutes.POST("/users/cart/add", controllers.AddToCart())
	incomingRoutes.GET("/users/cart", controllers.GetCart())
	incomingRoutes.PUT("/users/cart/update", controllers.UpdateCart())
	incomingRoutes.DELETE("/users/cart/delete/:cartItemID", controllers.DeleteCartItem())

	//order
	incomingRoutes.POST("/users/checkout", controllers.Checkout())
	incomingRoutes.GET("/users/orders", controllers.GetListOrder())
	incomingRoutes.PUT("/users/orders/cancel/:orderId", controllers.CancelOrder())
	incomingRoutes.GET("/users/orders/get-single", controllers.GetOrder())
	incomingRoutes.GET("/users/orders/get-raw-order", controllers.GetRawOrder())
	incomingRoutes.PUT("/users/order/update", controllers.GetOneOrder())
	incomingRoutes.PATCH("/")

	//address
	// incomingRoutes.POST("/users/address", controllers.AddAddress())

	//rating
	incomingRoutes.POST("/users/rating", controllers.CreateRating())
	incomingRoutes.GET("/users/rating/:productId", controllers.GetRating())

}

func AdminRouter(incomingRoutes *gin.Engine) {

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
}
