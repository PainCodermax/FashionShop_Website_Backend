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
	incomingRoutes.POST("/users/gettoken", controllers.GetNewToken())
	incomingRoutes.GET("/users/get-user", controllers.GetUser())
	incomingRoutes.POST("/users/forget-password", controllers.ForGotPassword())
	// incomingRoutes.PUT("/users/update-password", controllers.UpdatePassWord())

	// product
	incomingRoutes.GET("/users/product/list", controllers.GetListProduct())
	incomingRoutes.GET("/users/product/:productId", controllers.GetProduct())

	//category
	incomingRoutes.GET("/users/getcategory/list", controllers.GetCategoryList())

	//cart
	incomingRoutes.POST("/users/cart/add", controllers.AddToCart())
	incomingRoutes.GET("/users/cart", controllers.GetCart())

}

func AdminRouter(incomingRoutes *gin.Engine) {

	//product
	incomingRoutes.POST("/admin/addproduct", controllers.AddProduct())
	incomingRoutes.GET("/admin/getlistproduct", controllers.GetListProduct())
	incomingRoutes.PUT("/admin/product/update/:productId", controllers.UpdateProduct())
	incomingRoutes.DELETE("/admin/product/delete/:productId", controllers.DeleteProduct())

	//category
	incomingRoutes.POST("/admin/addcategory", controllers.AddCategory())
	incomingRoutes.GET("/admin/getcategory", controllers.GetCategory())
	incomingRoutes.GET("/admin/getcategory/list", controllers.GetCategoryList())
	incomingRoutes.PUT("/admin/updatecategory/:categoryId", controllers.UpdateCategory())
}
