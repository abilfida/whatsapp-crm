package controllers

import (
	"strconv"
	"whatsapp-crm/internal/models"
	"whatsapp-crm/pkg/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserController struct {
	db *gorm.DB
}

func NewUserController(db *gorm.DB) *UserController {
	return &UserController{db: db}
}

// GetUsers returns list of users with pagination
func (uc *UserController) GetUsers(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	role := c.Query("role")
	status := c.Query("status")
	search := c.Query("search")

	offset := (page - 1) * limit

	query := uc.db.Model(&models.User{})

	// Apply filters
	if role != "" {
		query = query.Where("role = ?", role)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if search != "" {
		query = query.Where("name LIKE ? OR email LIKE ?", "%"+search+"%", "%"+search+"%")
	}

	// Get total count
	var total int64
	query.Count(&total)

	// Get users
	var users []models.User
	if err := query.Select("id, name, email, phone, role, status, avatar, last_login, created_at, updated_at").Offset(offset).Limit(limit).Find(&users).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch users",
		})
	}

	return c.JSON(fiber.Map{
		"users": users,
		"pagination": fiber.Map{
			"page":  page,
			"limit": limit,
			"total": total,
			"pages": (total + int64(limit) - 1) / int64(limit),
		},
	})
}

// GetUser returns single user by ID
func (uc *UserController) GetUser(c *fiber.Ctx) error {
	userID := c.Params("id")

	var user models.User
	if err := uc.db.Select("id, name, email, phone, role, status, avatar, last_login, created_at, updated_at").First(&user, "id = ?", userID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "User not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch user",
		})
	}

	return c.JSON(user)
}

// CreateUser creates new user (admin only)
func (uc *UserController) CreateUser(c *fiber.Ctx) error {
	var req struct {
		Name     string `json:"name" validate:"required"`
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,min=6"`
		Phone    string `json:"phone"`
		Role     string `json:"role" validate:"required"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Check if email already exists
	var existingUser models.User
	if err := uc.db.First(&existingUser, "email = ?", req.Email).Error; err == nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": "Email already exists",
		})
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to hash password",
		})
	}

	// Create user
	user := models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: hashedPassword,
		Phone:    req.Phone,
		Role:     models.UserRole(req.Role),
		Status:   models.UserStatusActive,
	}

	if err := uc.db.Create(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create user",
		})
	}

	user.Password = ""
	return c.Status(fiber.StatusCreated).JSON(user)
}

// UpdateUser updates user by ID
func (uc *UserController) UpdateUser(c *fiber.Ctx) error {
	userID := c.Params("id")

	var req struct {
		Name   string `json:"name"`
		Phone  string `json:"phone"`
		Role   string `json:"role"`
		Status string `json:"status"`
		Avatar string `json:"avatar"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	var user models.User
	if err := uc.db.First(&user, "id = ?", userID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "User not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch user",
		})
	}

	// Update fields
	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Phone != "" {
		user.Phone = req.Phone
	}
	if req.Role != "" {
		user.Role = models.UserRole(req.Role)
	}
	if req.Status != "" {
		user.Status = models.UserStatus(req.Status)
	}
	if req.Avatar != "" {
		user.Avatar = req.Avatar
	}

	if err := uc.db.Save(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update user",
		})
	}

	user.Password = ""
	return c.JSON(user)
}

// DeleteUser deletes user by ID (soft delete)
func (uc *UserController) DeleteUser(c *fiber.Ctx) error {
	userID := c.Params("id")
	currentUser := c.Locals("user").(*models.User)

	// Prevent self-deletion
	if currentUser.ID.String() == userID {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot delete your own account",
		})
	}

	var user models.User
	if err := uc.db.First(&user, "id = ?", userID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "User not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch user",
		})
	}

	if err := uc.db.Delete(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete user",
		})
	}

	return c.JSON(fiber.Map{
		"message": "User deleted successfully",
	})
}