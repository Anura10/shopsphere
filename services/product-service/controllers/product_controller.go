package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"shopsphere/product-service/models"
	"shopsphere/product-service/services"
)

type ProductController struct {
	productService *services.ProductService
}

func NewProductController(
	productService *services.ProductService,
) *ProductController {
	return &ProductController{
		productService: productService,
	}
}

// CreateProduct creates a new product.
func (c *ProductController) CreateProduct(ctx *gin.Context) {

	var req models.CreateProductRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid product data",
			"error":   err.Error(),
		})
		return
	}

	product, err := c.productService.Create(
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
		"status":  "success",
		"message": "Product created successfully",
		"product": product,
	})
}

// GetProduct returns a single product.
func (c *ProductController) GetProduct(ctx *gin.Context) {

	id, err := strconv.ParseInt(
		ctx.Param("id"),
		10,
		64,
	)

	if err != nil || id <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid product ID",
		})
		return
	}

	product, err := c.productService.GetByID(
		ctx.Request.Context(),
		id,
	)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	if product == nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Product not found",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"product": product,
	})
}

// GetProducts returns all active products.
func (c *ProductController) GetProducts(ctx *gin.Context) {

	products, err := c.productService.GetAll(
		ctx.Request.Context(),
	)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":   "success",
		"products": products,
		"count":    len(products),
	})
}
func (c *ProductController) UpdateProduct(ctx *gin.Context) {

	id, err := strconv.ParseInt(
		ctx.Param("id"),
		10,
		64,
	)

	if err != nil || id <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid product ID",
		})
		return
	}

	var req models.UpdateProductRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid product data",
			"error":   err.Error(),
		})
		return
	}

	product, err := c.productService.Update(
		ctx.Request.Context(),
		id,
		req,
	)

	if err != nil {

		if err.Error() == "product not found" {
			ctx.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"message": "Product not found",
			})
			return
		}

		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Product updated successfully",
		"product": product,
	})
}

func (c *ProductController) DeleteProduct(ctx *gin.Context) {

	id, err := strconv.ParseInt(
		ctx.Param("id"),
		10,
		64,
	)

	if err != nil || id <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid product ID",
		})
		return
	}

	err = c.productService.Delete(
		ctx.Request.Context(),
		id,
	)

	if err != nil {

		if err.Error() == "product not found" {
			ctx.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"message": "Product not found",
			})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Product deleted successfully",
	})
}
func (c *ProductController) SearchProducts(ctx *gin.Context) {

	var filter models.ProductFilter

	filter.Search = ctx.Query("search")

	if categoryID := ctx.Query("category_id"); categoryID != "" {
		id, err := strconv.ParseInt(categoryID, 10, 64)

		if err != nil || id <= 0 {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Invalid category_id",
			})
			return
		}

		filter.CategoryID = id
	}

	if minPrice := ctx.Query("min_price"); minPrice != "" {
		price, err := strconv.ParseFloat(minPrice, 64)

		if err != nil || price < 0 {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Invalid min_price",
			})
			return
		}

		filter.MinPrice = price
	}

	if maxPrice := ctx.Query("max_price"); maxPrice != "" {
		price, err := strconv.ParseFloat(maxPrice, 64)

		if err != nil || price < 0 {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Invalid max_price",
			})
			return
		}

		filter.MaxPrice = price
	}

	if filter.MinPrice > 0 &&
		filter.MaxPrice > 0 &&
		filter.MinPrice > filter.MaxPrice {

		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "min_price cannot be greater than max_price",
		})
		return
	}

	products, err := c.productService.Search(
		ctx.Request.Context(),
		filter,
	)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":   "success",
		"products": products,
		"count":    len(products),
	})
}
