package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"shopsphere/order-service/config"
	"shopsphere/order-service/controllers"
	"shopsphere/order-service/database"
	"shopsphere/order-service/repository"
	"shopsphere/order-service/routes"
	"shopsphere/order-service/services"
)

func main() {
	cfg := config.LoadConfig()

	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}

	defer db.Close()

	orderRepository := repository.NewOrderRepository(db)

	cartClient := services.NewCartClient(
		cfg.CartURL,
	)

	inventoryClient := services.NewInventoryClient(
		cfg.InventoryURL,
	)

	orderService := services.NewOrderService(
		orderRepository,
		cartClient,
		inventoryClient,
	)

	orderController := controllers.NewOrderController(
		orderService,
	)

	router := gin.Default()

	// CORS middleware
	router.Use(func(c *gin.Context) {

		c.Writer.Header().Set(
			"Access-Control-Allow-Origin",
			"http://localhost:5173",
		)

		c.Writer.Header().Set(
			"Access-Control-Allow-Methods",
			"GET, POST, PUT, PATCH, DELETE, OPTIONS",
		)

		c.Writer.Header().Set(
			"Access-Control-Allow-Headers",
			"Origin, Content-Type, Accept, Authorization",
		)

		c.Writer.Header().Set(
			"Access-Control-Allow-Credentials",
			"true",
		)

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	})

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":   "success",
			"service":  "order-service",
			"message":  "ShopSphere Order Service is running",
			"database": "connected",
		})
	})

	api := router.Group("/api/v1")

	routes.RegisterOrderRoutes(
		api,
		orderController,
		cfg.JWTSecret,
	)

	log.Printf(
		"Order Service running on port %s",
		cfg.Port,
	)

	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatal(err)
	}
}