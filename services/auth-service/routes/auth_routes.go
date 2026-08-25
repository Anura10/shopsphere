package routes

import (
	"github.com/gin-gonic/gin"

	"shopsphere/auth-service/controllers"
	"shopsphere/auth-service/middleware"
)

func RegisterAuthRoutes(
	router *gin.RouterGroup,
	authController *controllers.AuthController,
) {
	auth := router.Group("/auth")

	// Public routes
	auth.POST("/register", authController.Register)
	auth.POST("/login", authController.Login)

	// Authentication required
	protected := auth.Group("")
	protected.Use(middleware.JWTAuthMiddleware())

	protected.GET("/profile", authController.Profile)

	// Admin-only routes
	admin := protected.Group("/admin")
	admin.Use(middleware.RequireRole(2))

	admin.GET("/test", authController.AdminTest)
}
