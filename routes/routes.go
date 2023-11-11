package routes

import (
	"github.com/PainCodermax/FashionShop_Website_Backend/controllers"

	"github.com/gin-gonic/gin"
)

func UserRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("/users/signup", controllers.SignUp())
	incomingRoutes.POST("/users/login", controllers.Login())
	incomingRoutes.POST("/users/gettoken", controllers.GetNewToken())
	// incomingRoutes.POST("/admin/addproduct", controllers.ProductViewerAdmin())
	incomingRoutes.GET("/users/productview", controllers.SearchProduct())
	incomingRoutes.GET("/users/search", controllers.SearchProductByQuery())
}

func AdminRouter(incomingRoutes *gin.Engine) {

	//product
	incomingRoutes.POST("/admin/addproduct", controllers.AddProduct())
	incomingRoutes.GET("/admin/getlistproduct", controllers.GetListProduct())
	incomingRoutes.PUT("/admin/product/update", controllers.UpdateProduct())
	incomingRoutes.DELETE("/admin/product/delete", controllers.DeleteProduct())

	//category
	incomingRoutes.POST("/admin/addcategory", controllers.AddCategory())
	incomingRoutes.GET("/admin/getcategory", controllers.GetCategory())
	incomingRoutes.GET("/admin/getcategory/list", controllers.GetCategoryList())
}
