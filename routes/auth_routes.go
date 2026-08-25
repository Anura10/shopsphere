package routes

import (
	"github.com/gin-gonic/gin"

	"shopsphere/auth-service/controllers"
)

func RegisterAuthRoutes(
	router *gin.RouterGroup,
	authController *controllers.AuthController,
) {
	auth := router.Group("/auth")

	auth.POST("/register", authController.Register)
}