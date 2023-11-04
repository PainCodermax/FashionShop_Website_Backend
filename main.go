package main

import (
	"log"
	"os"

	"github.com/PainCodermax/FashionShop_Website_Backend/controllers"
	"github.com/PainCodermax/FashionShop_Website_Backend/database"
	"github.com/PainCodermax/FashionShop_Website_Backend/middleware"
	"github.com/PainCodermax/FashionShop_Website_Backend/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	app := controllers.NewApplication(
		database.ProductData(database.Client, "product"),
		database.UserData(database.Client, "user"),
		database.UserData(database.Client, "category"),
	)

	// gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(middleware.CORSMiddleware())
	routes.UserRoutes(router)
	router.Use(middleware.Authentication())

	// admin
	routes.AdminRouter(router)

	router.GET("/addtocart", app.AddToCart())
	router.GET("/removeitem", app.RemoveItem())
	router.GET("/listcart", controllers.GetItemFromCart())
	router.POST("/addaddress", controllers.AddAddress())
	router.PUT("/edithomeaddress", controllers.EditHomeAddress())
	router.PUT("/editworkaddress", controllers.EditWorkAddress())
	router.GET("/deleteaddresses", controllers.DeleteAddress())
	router.GET("/cartcheckout", app.BuyFromCart())
	router.GET("/instantbuy", app.InstantBuy())
	log.Fatal(router.Run(":" + port))
}


