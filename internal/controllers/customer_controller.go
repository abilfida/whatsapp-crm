package controllers

import (
	"strconv"
	"whatsapp-crm/internal/services"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CustomerController struct {
	db  *gorm.DB
	csv *services.CustomerService
}

func NewCustomerController(db *gorm.DB, csv *services.CustomerService) *CustomerController {
	return &CustomerController{db: db, csv: csv}
}

func (cc *CustomerController) List(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	search := c.Query("search")

	customers, total, err := cc.csv.GetCustomers(page, limit, search)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch customers"})
	}

	return c.JSON(fiber.Map{
		"customers": customers,
		"pagination": fiber.Map{
			"page":  page,
			"limit": limit,
			"total": total,
			"pages": (total + int64(limit) - 1) / int64(limit),
		},
	})
}

func (cc *CustomerController) Create(c *fiber.Ctx) error {
	var req struct {
		Name       string `json:"name"`
		Email      string `json:"email"`
		Phone      string `json:"phone"`
		WhatsAppID string `json:"whatsapp_id"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid body"})
	}

	cust, err := cc.csv.UpdateCustomer(uuid.Nil, map[string]interface{}{})
	_ = cust
	// Create direct via GORM to keep it simple here
	type createReq = struct{ Name, Email, Phone, WhatsAppID string }
	cr := createReq{req.Name, req.Email, req.Phone, req.WhatsAppID}
	if cr.WhatsAppID != "" {
		if _, err := cc.csv.GetOrCreateCustomer(cr.WhatsAppID, cr.Name); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create customer"})
		}
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Customer created/ensured"})
}

func (cc *CustomerController) Detail(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID"})
	}
	cust, err := cc.csv.GetCustomerByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Customer not found"})
	}
	return c.JSON(cust)
}

func (cc *CustomerController) Update(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID"})
	}
	updates := map[string]interface{}{}
	if err := c.BodyParser(&updates); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid body"})
	}
	cust, err := cc.csv.UpdateCustomer(id, updates)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update"})
	}
	return c.JSON(cust)
}

func (cc *CustomerController) Delete(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID"})
	}
	if err := cc.csv.DeleteCustomer(id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete"})
	}
	return c.JSON(fiber.Map{"message": "Deleted"})
}
