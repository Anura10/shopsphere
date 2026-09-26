package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"shopsphere/auth-service/config"
	"shopsphere/auth-service/controllers"
	"shopsphere/auth-service/database"
	"shopsphere/auth-service/repository"
	"shopsphere/auth-service/routes"
	"shopsphere/auth-service/services"
)

func main() {
	cfg := config.LoadConfig()

	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}

	defer db.Close()

	// Repository
	userRepository := repository.NewUserRepository(db)

	// Service
	authService := services.NewAuthService(userRepository)

	// Controller
	authController := controllers.NewAuthController(authService)

	// Gin
	router := gin.Default()

	// =========================
	// CORS
	// =========================

	router.Use(func(c *gin.Context) {
		origin := c.GetHeader("Origin")

		if origin == "http://localhost:5173" ||
			origin == "http://localhost:5174" {
			c.Header("Access-Control-Allow-Origin", origin)
		}

		c.Header(
			"Access-Control-Allow-Credentials",
			"true",
		)

		c.Header(
			"Access-Control-Allow-Headers",
			"Origin, Content-Type, Accept, Authorization",
		)

		c.Header(
			"Access-Control-Allow-Methods",
			"GET, POST, PUT, PATCH, DELETE, OPTIONS",
		)

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	})

	// =========================
	// Health Check
	// =========================

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":   "success",
			"service":  "auth-service",
			"database": "connected",
			"message":  "ShopSphere Auth Service is running",
		})
	})

	// =========================
	// API v1
	// =========================

	api := router.Group("/api/v1")

	routes.RegisterAuthRoutes(
		api,
		authController,
	)

	log.Printf(
		"Auth Service running on port %s",
		cfg.Port,
	)

	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatal(err)
	}
}