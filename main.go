package main

import (
	"context"
	"log"
	"os"

	"github.com/PainCodermax/FashionShop_Website_Backend/client"
	"github.com/PainCodermax/FashionShop_Website_Backend/middleware"
	"github.com/PainCodermax/FashionShop_Website_Backend/routes"
	"github.com/PainCodermax/FashionShop_Website_Backend/worker"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
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
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	router.Use(gin.Logger())
	router.Use(middleware.CORSMiddleware())
	routes.LoginRoutes(router)
	// user

	workerChannel := make(chan string)
	go worker.Worker(ctx, workerChannel)

	router.Use(middleware.Authentication())
	routes.UserRoutes(router)
	// admin
	routes.AdminRouter(router)
	log.Fatal(router.Run(":" + port))
}
