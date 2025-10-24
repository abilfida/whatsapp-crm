package controllers

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"whatsapp-crm/internal/config"
	"whatsapp-crm/internal/models"
	"whatsapp-crm/internal/services"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UploadController struct{ db *gorm.DB; mu *services.MediaUploader; cfg *config.Config }

func NewUploadController(db *gorm.DB, mu *services.MediaUploader, cfg *config.Config) *UploadController { return &UploadController{db: db, mu: mu, cfg: cfg} }

func (uc *UploadController) withinSize(size int64) bool { return size <= uc.cfg.MaxFileSize }

func (uc *UploadController) validateType(declared, filename string) (string, bool) {
	ext := strings.TrimPrefix(strings.ToLower(filepath.Ext(filename)), ".")
	if ext == "" { return "", false }
	in := func(list []string, v string) bool { for _, x := range list { if x == v { return true } }; return false }

	// declared type precedence; if not set, detect by ext
	mediaType := strings.ToLower(strings.TrimSpace(declared))
	if mediaType == "" {
		switch {
		case in(uc.cfg.AllowedImageTypes, ext): mediaType = "image"
		case in(uc.cfg.AllowedDocumentTypes, ext): mediaType = "document"
		case in(uc.cfg.AllowedAudioTypes, ext): mediaType = "audio"
		case in(uc.cfg.AllowedVideoTypes, ext): mediaType = "video"
		default: return "", false
		}
	} else {
		// verify ext is allowed for declared type
		switch mediaType {
		case "image": if !in(uc.cfg.AllowedImageTypes, ext) { return "", false }
		case "document": if !in(uc.cfg.AllowedDocumentTypes, ext) { return "", false }
		case "audio": if !in(uc.cfg.AllowedAudioTypes, ext) { return "", false }
		case "video": if !in(uc.cfg.AllowedVideoTypes, ext) { return "", false }
		default: return "", false
		}
	}
	return mediaType, true
}

// POST /api/v1/messages/conversation/:id/upload
func (uc *UploadController) UploadAndSend(c *fiber.Ctx) error {
	cid, err := uuid.Parse(c.Params("id"))
	if err != nil { return c.Status(400).JSON(fiber.Map{"error":"invalid conversation id"}) }
	file, err := c.FormFile("file")
	if err != nil { return c.Status(400).JSON(fiber.Map{"error":"file is required"}) }

	if !uc.withinSize(file.Size) {
		return c.Status(400).JSON(fiber.Map{"error": fmt.Sprintf("file too large (max %d bytes)", uc.cfg.MaxFileSize)})
	}

	caption := c.FormValue("caption")
	declaredType := c.FormValue("type")
	mediaType, ok := uc.validateType(declaredType, file.Filename)
	if !ok { return c.Status(400).JSON(fiber.Map{"error":"unsupported file type/extension"}) }

	var conv models.Conversation
	if err := uc.db.Preload("Customer").First(&conv, "id = ?", cid).Error; err != nil { return c.Status(404).JSON(fiber.Map{"error":"conversation not found"}) }

	// pass mediaType as declared type so uploader uses intended WA API
	msg, err := uc.mu.UploadAndSend(context.Background(), &conv, file, mediaType, caption)
	if err != nil { return c.Status(502).JSON(fiber.Map{"error": fmt.Sprintf("upload/send failed: %v", err)}) }

	if err := uc.db.Create(msg).Error; err != nil { return c.Status(500).JSON(fiber.Map{"error":"failed to save message"}) }
	return c.Status(201).JSON(msg)
}
