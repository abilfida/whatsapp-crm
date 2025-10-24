package middlewares

import (
	"strings"
	"whatsapp-crm/internal/models"
	"whatsapp-crm/pkg/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm"
)

type AuthMiddleware struct {
	db *gorm.DB
}

func NewAuthMiddleware(db *gorm.DB) *AuthMiddleware {
	return &AuthMiddleware{db: db}
}

func (a *AuthMiddleware) RequireAuth(c *fiber.Ctx) error {
	// Get token from Authorization header
	authorization := c.Get("Authorization")
	if authorization == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Authorization header is required",
		})
	}

	// Extract token from "Bearer TOKEN"
	tokenParts := strings.Split(authorization, " ")
	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid authorization format",
		})
	}

	tokenString := tokenParts[1]

	// Parse and validate token
	claims, err := utils.ValidateJWT(tokenString)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token",
		})
	}

	// Get user from database
	var user models.User
	if err := a.db.First(&user, "id = ?", claims.UserID).Error; err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	// Check user status
	if user.Status == models.UserStatusInactive {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Account is inactive",
		})
	}

	// Set user in context
	c.Locals("user", &user)
	c.Locals("user_id", user.ID.String())
	c.Locals("user_role", string(user.Role))

	return c.Next()
}

func (a *AuthMiddleware) RequireRole(roles ...models.UserRole) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user := c.Locals("user").(*models.User)
		
		for _, role := range roles {
			if user.Role == role {
				return c.Next()
			}
		}

		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Insufficient privileges",
		})
	}
}