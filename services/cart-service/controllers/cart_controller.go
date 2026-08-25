package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"shopsphere/cart-service/models"
	"shopsphere/cart-service/services"
)

type CartController struct {
	service *services.CartService
}

func NewCartController(
	service *services.CartService,
) *CartController {
	return &CartController{
		service: service,
	}
}

func getUserID(ctx *gin.Context) (int64, bool) {

	value, exists := ctx.Get("user_id")

	if !exists {
		return 0, false
	}

	userID, ok := value.(int64)

	return userID, ok
}

func (c *CartController) GetCart(ctx *gin.Context) {

	userID, ok := getUserID(ctx)

	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "User ID not found",
		})
		return
	}

	cart, err := c.service.GetOrCreateCart(
		ctx.Request.Context(),
		userID,
	)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"cart":   cart,
	})
}

func (c *CartController) AddItem(ctx *gin.Context) {

	userID, ok := getUserID(ctx)

	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "User ID not found",
		})
		return
	}

	var req models.AddCartItemRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid cart item data",
			"error":   err.Error(),
		})
		return
	}

	cart, err := c.service.AddItem(
		ctx.Request.Context(),
		userID,
		req,
	)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Product added to cart successfully",
		"cart":    cart,
	})
}

func (c *CartController) UpdateItem(ctx *gin.Context) {

	userID, ok := getUserID(ctx)

	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "User ID not found",
		})
		return
	}

	productID, err := strconv.ParseInt(
		ctx.Param("product_id"),
		10,
		64,
	)

	if err != nil || productID <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid product ID",
		})
		return
	}

	var req models.UpdateCartItemRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid quantity",
			"error":   err.Error(),
		})
		return
	}

	cart, err := c.service.UpdateItem(
		ctx.Request.Context(),
		userID,
		productID,
		req,
	)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Cart item updated successfully",
		"cart":    cart,
	})
}

func (c *CartController) RemoveItem(ctx *gin.Context) {

	userID, ok := getUserID(ctx)

	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "User ID not found",
		})
		return
	}

	productID, err := strconv.ParseInt(
		ctx.Param("product_id"),
		10,
		64,
	)

	if err != nil || productID <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid product ID",
		})
		return
	}

	cart, err := c.service.RemoveItem(
		ctx.Request.Context(),
		userID,
		productID,
	)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Product removed from cart successfully",
		"cart":    cart,
	})
}

func (c *CartController) ClearCart(ctx *gin.Context) {

	userID, ok := getUserID(ctx)

	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "User ID not found",
		})
		return
	}

	cart, err := c.service.ClearCart(
		ctx.Request.Context(),
		userID,
	)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Cart cleared successfully",
		"cart":    cart,
	})
}