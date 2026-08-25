package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"shopsphere/product-service/models"
	"shopsphere/product-service/services"
)

type CategoryController struct {
	categoryService *services.CategoryService
}

func NewCategoryController(
	categoryService *services.CategoryService,
) *CategoryController {
	return &CategoryController{
		categoryService: categoryService,
	}
}

func (c *CategoryController) CreateCategory(ctx *gin.Context) {

	var req models.CreateCategoryRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid category data",
			"error":   err.Error(),
		})
		return
	}

	category, err := c.categoryService.Create(
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
		"status":   "success",
		"message":  "Category created successfully",
		"category": category,
	})
}

func (c *CategoryController) GetCategories(ctx *gin.Context) {

	categories, err := c.categoryService.GetAll(
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
		"status":     "success",
		"categories": categories,
		"count":      len(categories),
	})
}

func (c *CategoryController) GetCategory(ctx *gin.Context) {

	id, err := strconv.ParseInt(
		ctx.Param("id"),
		10,
		64,
	)

	if err != nil || id <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid category ID",
		})
		return
	}

	category, err := c.categoryService.GetByID(
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

	if category == nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Category not found",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":   "success",
		"category": category,
	})
}

func (c *CategoryController) UpdateCategory(ctx *gin.Context) {

	id, err := strconv.ParseInt(
		ctx.Param("id"),
		10,
		64,
	)

	if err != nil || id <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid category ID",
		})
		return
	}

	var req models.UpdateCategoryRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid category data",
			"error":   err.Error(),
		})
		return
	}

	category, err := c.categoryService.Update(
		ctx.Request.Context(),
		id,
		req,
	)

	if err != nil {

		if err.Error() == "category not found" {
			ctx.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"message": "Category not found",
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
		"status":   "success",
		"message":  "Category updated successfully",
		"category": category,
	})
}

func (c *CategoryController) DeleteCategory(ctx *gin.Context) {

	id, err := strconv.ParseInt(
		ctx.Param("id"),
		10,
		64,
	)

	if err != nil || id <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid category ID",
		})
		return
	}

	err = c.categoryService.Delete(
		ctx.Request.Context(),
		id,
	)

	if err != nil {

		if err.Error() == "category not found" {
			ctx.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"message": "Category not found",
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
		"message": "Category deleted successfully",
	})
}
