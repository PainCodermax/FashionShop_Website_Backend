package main

import (
	"log"
	"os"

	"github.com/PainCodermax/FashionShop_Website_Backend/controllers"
	"github.com/PainCodermax/FashionShop_Website_Backend/middleware"
	"github.com/PainCodermax/FashionShop_Website_Backend/routes"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(".env")
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	// app := controllers.NewApplication(
	// 	database.ProductData(database.Client, "product"),
	// 	database.UserData(database.Client, "user"),
	// 	database.UserData(database.Client, "category"),
	// )

	// gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(middleware.CORSMiddleware())
	routes.UserRoutes(router)
	router.Use(middleware.Authentication())

	// admin
	routes.AdminRouter(router)

	router.POST("/addaddress", controllers.AddAddress())
	router.PUT("/edithomeaddress", controllers.EditHomeAddress())
	router.PUT("/editworkaddress", controllers.EditWorkAddress())
	router.GET("/deleteaddresses", controllers.DeleteAddress())
	log.Fatal(router.Run(":" + port))
}
