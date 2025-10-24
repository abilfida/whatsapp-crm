package controllers

import (
	"context"
	"fmt"
	"whatsapp-crm/internal/models"
	"whatsapp-crm/internal/services"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UploadController struct{ db *gorm.DB; mu *services.MediaUploader }

func NewUploadController(db *gorm.DB, mu *services.MediaUploader) *UploadController { return &UploadController{db: db, mu: mu} }

// POST /api/v1/messages/conversation/:id/upload
func (uc *UploadController) UploadAndSend(c *fiber.Ctx) error {
	cid, err := uuid.Parse(c.Params("id"))
	if err != nil { return c.Status(400).JSON(fiber.Map{"error":"invalid conversation id"}) }
	file, err := c.FormFile("file")
	if err != nil { return c.Status(400).JSON(fiber.Map{"error":"file is required"}) }
	caption := c.FormValue("caption")
	declaredType := c.FormValue("type")

	var conv models.Conversation
	if err := uc.db.Preload("Customer").First(&conv, "id = ?", cid).Error; err != nil { return c.Status(404).JSON(fiber.Map{"error":"conversation not found"}) }

	msg, err := uc.mu.UploadAndSend(context.Background(), &conv, file, declaredType, caption)
	if err != nil { return c.Status(502).JSON(fiber.Map{"error": fmt.Sprintf("upload/send failed: %v", err)}) }

	// persist message
	if err := uc.db.Create(msg).Error; err != nil { return c.Status(500).JSON(fiber.Map{"error":"failed to save message"}) }
	return c.Status(201).JSON(msg)
}
