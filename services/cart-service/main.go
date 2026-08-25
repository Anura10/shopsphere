package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"shopsphere/cart-service/config"
	"shopsphere/cart-service/controllers"
	"shopsphere/cart-service/database"
	"shopsphere/cart-service/repository"
	"shopsphere/cart-service/routes"
	"shopsphere/cart-service/services"
)

func main() {

	cfg := config.LoadConfig()

	db, err := database.Connect(cfg.DatabaseURL)

	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}

	defer db.Close()

	cartRepository := repository.NewCartRepository(db)

	cartService := services.NewCartService(
		cartRepository,
	)

	cartController := controllers.NewCartController(
		cartService,
	)

	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {

		c.JSON(http.StatusOK, gin.H{
			"status":   "success",
			"service":  "cart-service",
			"message":  "ShopSphere Cart Service is running",
			"database": "connected",
		})
	})

	api := router.Group("/api/v1")

	routes.RegisterCartRoutes(
		api,
		cartController,
		cfg.JWTSecret,
	)

	log.Printf(
		"Cart Service running on port %s",
		cfg.Port,
	)

	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatal(err)
	}
}