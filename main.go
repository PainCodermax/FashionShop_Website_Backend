package main

import (
	"log"
	"os"

	"github.com/PainCodermax/FashionShop_Website_Backend/client"
	"github.com/PainCodermax/FashionShop_Website_Backend/cronjob"
	"github.com/PainCodermax/FashionShop_Website_Backend/middleware"
	"github.com/PainCodermax/FashionShop_Website_Backend/routes"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/robfig/cron"
)

func main() {
	client.GOs()
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
	client.Init()
	// gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(middleware.CORSMiddleware())
	routes.LoginRoutes(router)
	// user

	router.Use(middleware.Authentication())
	routes.UserRoutes(router)
	// admin
	routes.AdminRouter(router)

	c := cron.New()
	c.AddFunc("0 0 * * *", cronjob.UpdateOrderStatusJob)
	c.Start()
	defer c.Stop()

	log.Fatal(router.Run(":" + port))
}
 