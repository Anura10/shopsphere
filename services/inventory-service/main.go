package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"shopsphere/inventory-service/config"
	"shopsphere/inventory-service/controllers"
	"shopsphere/inventory-service/database"
	"shopsphere/inventory-service/repository"
	"shopsphere/inventory-service/routes"
	"shopsphere/inventory-service/services"
)

func main() {

	cfg := config.LoadConfig()

	db, err := database.Connect(cfg.DatabaseURL)

	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}

	defer db.Close()

	inventoryRepository := repository.NewInventoryRepository(db)

	inventoryService := services.NewInventoryService(
		inventoryRepository,
	)

	inventoryController := controllers.NewInventoryController(
		inventoryService,
	)

	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {

		c.JSON(http.StatusOK, gin.H{
			"status":   "success",
			"service":  "inventory-service",
			"message":  "ShopSphere Inventory Service is running",
			"database": "connected",
		})
	})

	api := router.Group("/api/v1")

	routes.RegisterInventoryRoutes(
		api,
		inventoryController,
		cfg.JWTSecret,
	)

	log.Printf(
		"Inventory Service running on port %s",
		cfg.Port,
	)

	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatal(err)
	}
}
