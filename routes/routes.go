package routes

import (
	"github.com/PainCodermax/FashionShop_Website_Backend/controllers"

	"github.com/gin-gonic/gin"
)

func LoginRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("/users/signup", controllers.SignUp())
	incomingRoutes.POST("/users/login", controllers.Login())
	incomingRoutes.POST("/users/verify", controllers.VerifyUser())
}

func UserRoutes(incomingRoutes *gin.Engine) {
	// incomingRoutes.POST("/users/signup", controllers.SignUp())
	// incomingRoutes.POST("/users/login", controllers.Login())
	incomingRoutes.POST("/users/get-token", controllers.GetNewToken())
	incomingRoutes.GET("/users/get-user", controllers.GetUser())
	incomingRoutes.POST("/users/forget-password", controllers.ForGotPassword())
	incomingRoutes.PUT("/users/update-password", controllers.UpdatePassWord())

	// product
	incomingRoutes.GET("/users/product/list", controllers.GetListProduct())
	incomingRoutes.GET("/users/product/:productId", controllers.GetProduct())

	//category
	incomingRoutes.GET("/users/get-category/list", controllers.GetCategoryList())

	//cart
	incomingRoutes.POST("/users/cart/add", controllers.AddToCart())
	incomingRoutes.GET("/users/cart", controllers.GetCart())
	incomingRoutes.PUT("/users/cart/update", controllers.UpdateCart())
	incomingRoutes.DELETE("/users/cart/delete/:cartItemID", controllers.DeleteCartItem())

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
}
