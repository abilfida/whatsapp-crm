package services

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
	"whatsapp-crm/internal/config"
	"whatsapp-crm/internal/models"
	"whatsapp-crm/pkg/whatsapp"

	"github.com/google/uuid"
)

// UploadAndSendLocal handles local file uploads then sends appropriate media message
func (ms *MessageService) UploadAndSendLocal(conversationID uuid.UUID, filePath, mediaType, caption string) (*models.Message, error) {
	var conv models.Conversation
	if err := ms.db.Preload("Customer").First(&conv, "id = ?", conversationID).Error; err != nil { return nil, err }

	// ensure file exists
	if _, err := os.Stat(filePath); err != nil { return nil, fmt.Errorf("file not found: %w", err) }

	// upload to WA gateway
	uploadURL, err := ms.wa.UploadLocalMedia(filePath)
	if err != nil { return nil, fmt.Errorf("upload media: %w", err) }

	// send by media type
	var resp *whatsapp.SendMessageResponse
	switch mediaType {
	case "image":
		resp, err = ms.wa.SendImageMessage(conv.Customer.WhatsAppID, uploadURL, caption)
	case "document":
		resp, err = ms.wa.SendDocumentMessage(conv.Customer.WhatsAppID, uploadURL, filepath.Base(filePath), caption)
	case "audio":
		resp, err = ms.wa.SendAudioMessage(conv.Customer.WhatsAppID, uploadURL)
	case "video":
		resp, err = ms.wa.SendVideoMessage(conv.Customer.WhatsAppID, uploadURL, caption)
	default:
		return nil, fmt.Errorf("unsupported media type: %s", mediaType)
	}
	if err != nil { return nil, err }

	now := time.Now()
	msg := models.Message{ConversationID: conversationID, WhatsAppID: resp.ID, Type: models.MessageType(mediaType), Direction: models.MessageDirectionOutbound, Status: models.MessageStatusSent, MediaURL: uploadURL, Caption: caption, FileName: filepath.Base(filePath), SentAt: &now}
	if err := ms.db.Create(&msg).Error; err != nil { return nil, err }
	conv.LastMessageAt = &now
	_ = ms.db.Save(&conv)
	return &msg, nil
}
