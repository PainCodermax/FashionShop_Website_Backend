package main

import (
	"os"

	"github.com/PainCodermax/FashionShop_Website_Backend/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	router.Use(gin.Logger())
	routes.AdminRoutes(router)

	router.Run(":" + port)
}
