package routes

import (
	"github.com/gin-gonic/gin"

	"shopsphere/inventory-service/controllers"
	"shopsphere/inventory-service/middleware"
)

func RegisterInventoryRoutes(
	router *gin.RouterGroup,
	controller *controllers.InventoryController,
	jwtSecret string,
) {
	inventory := router.Group("/inventory")

	// Public stock lookup
	inventory.GET(
		"/:product_id",
		controller.GetByProductID,
	)

	// ADMIN operations
	admin := inventory.Group("")
	admin.Use(middleware.JWTAuthMiddleware(jwtSecret))
	admin.Use(middleware.RequireRole(2))

	admin.POST("", controller.Create)
	admin.PUT("/:product_id", controller.Update)
	admin.POST("/:product_id/add", controller.AddStock)
	admin.POST("/:product_id/remove", controller.RemoveStock)
}
