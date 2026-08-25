package controllers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"shopsphere/auth-service/models"
	"shopsphere/auth-service/services"
)

type AuthController struct {
	authService *services.AuthService
}

func NewAuthController(
	authService *services.AuthService,
) *AuthController {
	return &AuthController{
		authService: authService,
	}
}

func (c *AuthController) Register(ctx *gin.Context) {

	var req models.RegisterRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid registration data",
			"error":   err.Error(),
		})
		return
	}

	user, err := c.authService.Register(
		ctx.Request.Context(),
		req,
	)

	if err != nil {

		if errors.Is(err, services.ErrEmailExists) {
			ctx.JSON(http.StatusConflict, gin.H{
				"status":  "error",
				"message": "Email already registered",
			})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to register user",
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "User registered successfully",
		"user": gin.H{
			"id":         user.ID,
			"first_name": user.FirstName,
			"last_name":  user.LastName,
			"email":      user.Email,
			"phone":      user.Phone,
			"role_id":    user.RoleID,
		},
	})
}

func (c *AuthController) Login(ctx *gin.Context) {

	var req models.LoginRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid login data",
			"error":   err.Error(),
		})
		return
	}

	user, token, err := c.authService.Login(
		ctx.Request.Context(),
		req,
	)

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":       "success",
		"message":      "Login successful",
		"access_token": token,
		"token_type":   "Bearer",
		"expires_in":   86400,
		"user": gin.H{
			"id":         user.ID,
			"first_name": user.FirstName,
			"last_name":  user.LastName,
			"email":      user.Email,
			"role_id":    user.RoleID,
		},
	})
}
func (c *AuthController) Profile(ctx *gin.Context) {

	userID, exists := ctx.Get("user_id")

	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "User information not found",
		})
		return
	}

	email, _ := ctx.Get("email")
	roleID, _ := ctx.Get("role_id")

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Protected profile accessed successfully",
		"user": gin.H{
			"id":      userID,
			"email":   email,
			"role_id": roleID,
		},
	})
}
func (c *AuthController) AdminTest(ctx *gin.Context) {

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Admin access granted",
		"role":    "ADMIN",
	})
}
