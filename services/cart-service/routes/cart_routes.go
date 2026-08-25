package routes

import (
	"github.com/gin-gonic/gin"

	"shopsphere/cart-service/controllers"
	"shopsphere/cart-service/middleware"
)

func RegisterCartRoutes(
	router *gin.RouterGroup,
	controller *controllers.CartController,
	jwtSecret string,
) {

	cart := router.Group("/cart")

	cart.Use(middleware.JWTAuthMiddleware(jwtSecret))

	cart.GET("", controller.GetCart)

	cart.POST(
		"/items",
		controller.AddItem,
	)

	cart.PUT(
		"/items/:product_id",
		controller.UpdateItem,
	)

	cart.DELETE(
		"/items/:product_id",
		controller.RemoveItem,
	)

	cart.DELETE(
		"",
		controller.ClearCart,
	)
}