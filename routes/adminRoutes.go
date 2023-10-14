package routes

import (
	"github.com/PainCodermax/FashionShop_Website_Backend/controllers"
	"github.com/gin-gonic/gin"
)

func AdminRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("admin/login", controllers.LoginAdmin())
}