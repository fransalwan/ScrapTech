package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"

	"scrapflow-backend/internal/domain"
	"scrapflow-backend/pkg/response"
	"scrapflow-backend/pkg/utils"
)

func AuthMiddleware(secretKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Unauthorized(c, "Authorization header is required")
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Unauthorized(c, "Authorization format must be Bearer <token>")
			c.Abort()
			return
		}

		claims, err := utils.ValidateToken(parts[1], secretKey)
		if err != nil {
			response.Unauthorized(c, "Invalid or expired authorization token")
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("email", claims.Email)
		c.Set("role", claims.Role)
		c.Set("company_id", claims.CompanyID)

		c.Next()
	}
}

func RequireRoles(allowedRoles ...domain.UserRole) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleVal, exists := c.Get("role")
		if !exists {
			response.Forbidden(c, "Access denied: missing role context")
			c.Abort()
			return
		}

		userRole := domain.UserRole(roleVal.(string))
		for _, role := range allowedRoles {
			if userRole == role {
				c.Next()
				return
			}
		}

		response.Forbidden(c, "Access denied: insufficient permission for this role")
		c.Abort()
	}
}
