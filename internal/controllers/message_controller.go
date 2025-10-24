package controllers

import (
	"strconv"
	"whatsapp-crm/internal/services"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MessageController struct { db *gorm.DB; ms *services.MessageService }

func NewMessageController(db *gorm.DB, ms *services.MessageService) *MessageController { return &MessageController{db: db, ms: ms} }

func (mc *MessageController) ListByConversation(c *fiber.Ctx) error {
	cid, err := uuid.Parse(c.Params("id")); if err != nil { return c.Status(400).JSON(fiber.Map{"error":"Invalid conversation id"}) }
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "50"))
	var total int64
	mc.db.Table("messages").Where("conversation_id = ?", cid).Count(&total)
	var msgs []map[string]any
	mc.db.Table("messages").Where("conversation_id = ?", cid).Order("created_at asc").Limit(limit).Offset((page-1)*limit).Find(&msgs)
	return c.JSON(fiber.Map{"messages": msgs, "pagination": fiber.Map{"page":page, "limit":limit, "total":total}})
}

func (mc *MessageController) SendText(c *fiber.Ctx) error {
	cid, err := uuid.Parse(c.Params("id")); if err != nil { return c.Status(400).JSON(fiber.Map{"error":"Invalid conversation id"}) }
	var req struct{ Content string `json:"content"` }
	if err := c.BodyParser(&req); err != nil { return c.Status(400).JSON(fiber.Map{"error":"Invalid body"}) }
	msg, err := mc.ms.SendText(cid, req.Content); if err != nil { return c.Status(500).JSON(fiber.Map{"error": err.Error()}) }
	return c.Status(201).JSON(msg)
}

func (mc *MessageController) SendMedia(c *fiber.Ctx) error { return c.Status(501).JSON(fiber.Map{"error":"Not implemented yet"}) }
func (mc *MessageController) SendTemplate(c *fiber.Ctx) error { return c.Status(501).JSON(fiber.Map{"error":"Not implemented yet"}) }
