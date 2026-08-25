package routes

import (
	"github.com/gin-gonic/gin"

	"shopsphere/product-service/controllers"
	"shopsphere/product-service/middleware"
)

func RegisterProductRoutes(
	router *gin.RouterGroup,
	productController *controllers.ProductController,
	jwtSecret string,
) {
	products := router.Group("/products")

	// Public routes for now.
	products.GET("", productController.GetProducts)
	products.GET("/search", productController.SearchProducts)
	products.GET("/:id", productController.GetProduct)

	// Authentication required.
	protected := products.Group("")
	protected.Use(middleware.JWTAuthMiddleware(jwtSecret))

	// ADMIN only.
	admin := protected.Group("")
	admin.Use(middleware.RequireRole(2))

	admin.POST("", productController.CreateProduct)
	admin.PUT("/:id", productController.UpdateProduct)
	admin.DELETE("/:id", productController.DeleteProduct)
}
