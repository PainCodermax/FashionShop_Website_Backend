package routes

import (
	"github.com/PainCodermax/FashionShop_Website_Backend/controllers"

	"github.com/gin-gonic/gin"
)

func UserRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("/users/signup", controllers.SignUp())
	incomingRoutes.POST("/users/login", controllers.Login())
	incomingRoutes.POST("/users/gettoken", controllers.GetNewToken())
	// incomingRoutes.POST("/users/verifygmail", controllers.VerifyGmail())

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
}
