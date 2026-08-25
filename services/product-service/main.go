package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"shopsphere/product-service/config"
	"shopsphere/product-service/controllers"
	"shopsphere/product-service/database"
	"shopsphere/product-service/repository"
	"shopsphere/product-service/routes"
	"shopsphere/product-service/services"
)

func main() {

	cfg := config.LoadConfig()

	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}

	defer db.Close()

	// Repository
	productRepository := repository.NewProductRepository(db)

	// Service
	productService := services.NewProductService(
		productRepository,
	)

	// Controller
	productController := controllers.NewProductController(
		productService,
	)
	categoryRepository := repository.NewCategoryRepository(db)

	categoryService := services.NewCategoryService(
		categoryRepository,
	)

	categoryController := controllers.NewCategoryController(
		categoryService,
	)

	// Router
	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":   "success",
			"service":  "product-service",
			"database": "connected",
			"message":  "ShopSphere Product Service is running",
		})
	})

	api := router.Group("/api/v1")

	routes.RegisterProductRoutes(
		api,
		productController,
		cfg.JWTSecret,
	)
	routes.RegisterCategoryRoutes(
		api,
		categoryController,
		cfg.JWTSecret,
	)

	log.Printf(
		"Product Service running on port %s",
		cfg.Port,
	)

	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatal(err)
	}
}
