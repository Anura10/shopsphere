package routes

import (
	"github.com/gin-gonic/gin"

	"shopsphere/product-service/controllers"
	"shopsphere/product-service/middleware"
)

func RegisterCategoryRoutes(
	router *gin.RouterGroup,
	categoryController *controllers.CategoryController,
	jwtSecret string,
) {
	categories := router.Group("/categories")

	// Public
	categories.GET("", categoryController.GetCategories)
	categories.GET("/:id", categoryController.GetCategory)

	// ADMIN only
	admin := categories.Group("")
	admin.Use(middleware.JWTAuthMiddleware(jwtSecret))
	admin.Use(middleware.RequireRole(2))

	admin.POST("", categoryController.CreateCategory)
	admin.PUT("/:id", categoryController.UpdateCategory)
	admin.DELETE("/:id", categoryController.DeleteCategory)
}
