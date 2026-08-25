package routes

import (
	"github.com/gin-gonic/gin"

	"shopsphere/order-service/controllers"
	"shopsphere/order-service/middleware"
)

func RegisterOrderRoutes(
	router *gin.RouterGroup,
	controller *controllers.OrderController,
	jwtSecret string,
) {

	order := router.Group("/orders")

	order.Use(
		middleware.JWTAuthMiddleware(jwtSecret),
	)

	order.POST("", controller.CreateOrder)

	order.GET("", controller.GetOrders)

	order.GET("/:id", controller.GetOrder)

	order.PUT("/:id/cancel", controller.CancelOrder)
}
