package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"shopsphere/order-service/models"
	"shopsphere/order-service/services"
)

type OrderController struct {
	service *services.OrderService
}

func NewOrderController(
	service *services.OrderService,
) *OrderController {
	return &OrderController{
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

func getToken(ctx *gin.Context) (string, bool) {
	token, exists := ctx.Get("token")

	if !exists {
		return "", false
	}

	tokenString, ok := token.(string)

	return tokenString, ok
}

func (c *OrderController) CreateOrder(ctx *gin.Context) {

	userID, ok := getUserID(ctx)

	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "User ID not found",
		})
		return
	}

	token, ok := getToken(ctx)

	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "Authentication token not found",
		})
		return
	}

	var req models.CreateOrderRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid order data",
			"error":   err.Error(),
		})
		return
	}

	order, err := c.service.CreateOrder(
		ctx.Request.Context(),
		userID,
		token,
		req,
	)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "Order created successfully",
		"order":   order,
	})
}

func (c *OrderController) GetOrders(ctx *gin.Context) {

	userID, ok := getUserID(ctx)

	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "User ID not found",
		})
		return
	}

	orders, err := c.service.GetUserOrders(
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
		"orders": orders,
	})
}

func (c *OrderController) GetOrder(ctx *gin.Context) {

	userID, ok := getUserID(ctx)

	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "User ID not found",
		})
		return
	}

	orderID, err := strconv.ParseInt(
		ctx.Param("id"),
		10,
		64,
	)

	if err != nil || orderID <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid order ID",
		})
		return
	}

	order, err := c.service.GetOrder(
		ctx.Request.Context(),
		userID,
		orderID,
	)

	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"order":  order,
	})
}

func (c *OrderController) CancelOrder(ctx *gin.Context) {

	userID, ok := getUserID(ctx)

	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "User ID not found",
		})
		return
	}

	orderID, err := strconv.ParseInt(
		ctx.Param("id"),
		10,
		64,
	)

	if err != nil || orderID <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid order ID",
		})
		return
	}

	order, err := c.service.CancelOrder(
		ctx.Request.Context(),
		userID,
		orderID,
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
		"message": "Order cancelled successfully",
		"order":   order,
	})
}
