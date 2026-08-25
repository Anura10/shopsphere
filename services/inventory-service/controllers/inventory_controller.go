package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"shopsphere/inventory-service/models"
	"shopsphere/inventory-service/services"
)

type InventoryController struct {
	service *services.InventoryService
}

func NewInventoryController(
	service *services.InventoryService,
) *InventoryController {
	return &InventoryController{
		service: service,
	}
}

func (c *InventoryController) Create(ctx *gin.Context) {

	var req models.CreateInventoryRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid inventory data",
			"error":   err.Error(),
		})
		return
	}

	inventory, err := c.service.Create(
		ctx.Request.Context(),
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
		"status":    "success",
		"message":   "Inventory created successfully",
		"inventory": inventory,
	})
}

func (c *InventoryController) GetByProductID(ctx *gin.Context) {

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

	inventory, err := c.service.GetByProductID(
		ctx.Request.Context(),
		productID,
	)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	if inventory == nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Inventory not found",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":    "success",
		"inventory": inventory,
	})
}

func (c *InventoryController) Update(ctx *gin.Context) {

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

	var req models.UpdateInventoryRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid inventory data",
			"error":   err.Error(),
		})
		return
	}

	inventory, err := c.service.Update(
		ctx.Request.Context(),
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

	if inventory == nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Inventory not found",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":    "success",
		"message":   "Inventory updated successfully",
		"inventory": inventory,
	})
}

func (c *InventoryController) AddStock(ctx *gin.Context) {

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

	var req models.StockRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid stock data",
			"error":   err.Error(),
		})
		return
	}

	inventory, err := c.service.AddStock(
		ctx.Request.Context(),
		productID,
		req.Quantity,
	)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	if inventory == nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Inventory not found",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":    "success",
		"message":   "Stock added successfully",
		"inventory": inventory,
	})
}

func (c *InventoryController) RemoveStock(ctx *gin.Context) {

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

	var req models.StockRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid stock data",
			"error":   err.Error(),
		})
		return
	}

	inventory, err := c.service.RemoveStock(
		ctx.Request.Context(),
		productID,
		req.Quantity,
	)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":    "success",
		"message":   "Stock removed successfully",
		"inventory": inventory,
	})
}
