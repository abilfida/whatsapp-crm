package controllers

import (
	"strconv"
	"whatsapp-crm/internal/services"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ConversationController struct { db *gorm.DB; csv *services.ConversationService; ms *services.MessageService }

func NewConversationController(db *gorm.DB, csv *services.ConversationService, ms *services.MessageService) *ConversationController { return &ConversationController{db: db, csv: csv, ms: ms} }

func (cc *ConversationController) List(c *fiber.Ctx) error {
	// simplified list by recent
	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	var convs []struct{ ID string; CustomerID string; Status string; UpdatedAt string }
	cc.db.Table("conversations").Select("id, customer_id, status, updated_at").Order("updated_at desc").Limit(limit).Scan(&convs)
	return c.JSON(fiber.Map{"conversations": convs})
}

func (cc *ConversationController) Create(c *fiber.Ctx) error {
	var req struct{ CustomerID string `json:"customer_id"` }
	if err := c.BodyParser(&req); err != nil { return c.Status(400).JSON(fiber.Map{"error":"Invalid body"}) }
	cid, err := uuid.Parse(req.CustomerID); if err != nil { return c.Status(400).JSON(fiber.Map{"error":"Invalid customer_id"}) }
	conv, err := cc.csv.Create(cid); if err != nil { return c.Status(500).JSON(fiber.Map{"error":"Failed to create"}) }
	return c.Status(201).JSON(conv)
}

func (cc *ConversationController) Detail(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id")); if err != nil { return c.Status(400).JSON(fiber.Map{"error":"Invalid ID"}) }
	var conv interface{}
	if err := cc.db.Preload("Customer").Preload("Agent").Preload("Messages").First(&conv, "id = ?", id).Error; err != nil { return c.Status(404).JSON(fiber.Map{"error":"Not found"}) }
	return c.JSON(conv)
}

func (cc *ConversationController) Assign(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id")); if err != nil { return c.Status(400).JSON(fiber.Map{"error":"Invalid ID"}) }
	var req struct{ AgentID string `json:"agent_id"` }
	if err := c.BodyParser(&req); err != nil { return c.Status(400).JSON(fiber.Map{"error":"Invalid body"}) }
	aid, err := uuid.Parse(req.AgentID); if err != nil { return c.Status(400).JSON(fiber.Map{"error":"Invalid agent_id"}) }
	if err := cc.csv.Assign(id, aid); err != nil { return c.Status(500).JSON(fiber.Map{"error":"Failed to assign"}) }
	return c.JSON(fiber.Map{"message":"Assigned"})
}

func (cc *ConversationController) UpdateStatus(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id")); if err != nil { return c.Status(400).JSON(fiber.Map{"error":"Invalid ID"}) }
	var req struct{ Status string `json:"status"` }
	if err := c.BodyParser(&req); err != nil { return c.Status(400).JSON(fiber.Map{"error":"Invalid body"}) }
	if err := cc.db.Model(&struct{ }{}).Table("conversations").Where("id = ?", id).Update("status", req.Status).Error; err != nil { return c.Status(500).JSON(fiber.Map{"error":"Failed to update"}) }
	return c.JSON(fiber.Map{"message":"Updated"})
}

func (cc *ConversationController) UpdatePriority(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id")); if err != nil { return c.Status(400).JSON(fiber.Map{"error":"Invalid ID"}) }
	var req struct{ Priority string `json:"priority"` }
	if err := c.BodyParser(&req); err != nil { return c.Status(400).JSON(fiber.Map{"error":"Invalid body"}) }
	if err := cc.db.Table("conversations").Where("id = ?", id).Update("priority", req.Priority).Error; err != nil { return c.Status(500).JSON(fiber.Map{"error":"Failed to update"}) }
	return c.JSON(fiber.Map{"message":"Updated"})
}

func (cc *ConversationController) UpdateNotes(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id")); if err != nil { return c.Status(400).JSON(fiber.Map{"error":"Invalid ID"}) }
	var req struct{ Notes string `json:"notes"` }
	if err := c.BodyParser(&req); err != nil { return c.Status(400).JSON(fiber.Map{"error":"Invalid body"}) }
	if err := cc.db.Table("conversations").Where("id = ?", id).Update("notes", req.Notes).Error; err != nil { return c.Status(500).JSON(fiber.Map{"error":"Failed to update"}) }
	return c.JSON(fiber.Map{"message":"Updated"})
}
