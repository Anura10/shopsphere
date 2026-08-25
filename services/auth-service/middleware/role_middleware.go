package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RequireRole(allowedRoles ...int64) gin.HandlerFunc {
	return func(c *gin.Context) {

		roleIDValue, exists := c.Get("role_id")

		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  "error",
				"message": "User role not found",
			})
			c.Abort()
			return
		}

		roleID, ok := roleIDValue.(int64)

		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  "error",
				"message": "Invalid user role",
			})
			c.Abort()
			return
		}

		for _, allowedRole := range allowedRoles {
			if roleID == allowedRole {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, gin.H{
			"status":  "error",
			"message": "You do not have permission to access this resource",
		})

		c.Abort()
	}
}