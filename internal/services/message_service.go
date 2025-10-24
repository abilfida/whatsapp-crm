package services

import (
	"fmt"
	"time"
	"whatsapp-crm/internal/config"
	"whatsapp-crm/internal/models"
	"whatsapp-crm/pkg/whatsapp"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MessageService struct { db *gorm.DB; wa *whatsapp.Client }

func NewMessageService(db *gorm.DB, cfg *config.Config) *MessageService { return &MessageService{db: db, wa: whatsapp.NewClient(cfg)} }

func (ms *MessageService) SendText(conversationID uuid.UUID, content string) (*models.Message, error) {
	var conv models.Conversation
	if err := ms.db.Preload("Customer").First(&conv, "id = ?", conversationID).Error; err != nil { return nil, err }
	resp, err := ms.wa.SendTextMessage(conv.Customer.WhatsAppID, content)
	if err != nil { return nil, fmt.Errorf("send text: %w", err) }
	now := time.Now()
	msg := models.Message{ConversationID: conversationID, WhatsAppID: resp.ID, Type: models.MessageTypeText, Direction: models.MessageDirectionOutbound, Status: models.MessageStatusSent, Content: content, SentAt: &now}
	if err := ms.db.Create(&msg).Error; err != nil { return nil, err }
	conv.LastMessageAt = &now
	_ = ms.db.Save(&conv)
	return &msg, nil
}

func (ms *MessageService) SendMediaMessage(conversationID uuid.UUID, mediaType, mediaURL, caption, filename string) (*models.Message, error) {
	var conv models.Conversation
	if err := ms.db.Preload("Customer").First(&conv, "id = ?", conversationID).Error; err != nil { return nil, err }
	var resp *whatsapp.SendMessageResponse
	var err error
	switch mediaType {
	case "image":
		resp, err = ms.wa.SendImageMessage(conv.Customer.WhatsAppID, mediaURL, caption)
	case "document":
		resp, err = ms.wa.SendDocumentMessage(conv.Customer.WhatsAppID, mediaURL, filename, caption)
	default:
		return nil, fmt.Errorf("unsupported media type: %s", mediaType)
	}
	if err != nil { return nil, fmt.Errorf("send media: %w", err) }
	now := time.Now()
	msg := models.Message{ConversationID: conversationID, WhatsAppID: resp.ID, Type: models.MessageType(mediaType), Direction: models.MessageDirectionOutbound, Status: models.MessageStatusSent, MediaURL: mediaURL, Caption: caption, FileName: filename, SentAt: &now}
	if err := ms.db.Create(&msg).Error; err != nil { return nil, err }
	conv.LastMessageAt = &now
	_ = ms.db.Save(&conv)
	return &msg, nil
}

func (ms *MessageService) SendTemplateMessage(conversationID uuid.UUID, templateID uuid.UUID, variables map[string]string) (*models.Message, error) {
	var conv models.Conversation
	if err := ms.db.Preload("Customer").First(&conv, "id = ?", conversationID).Error; err != nil { return nil, err }
	var tpl models.Template
	if err := ms.db.First(&tpl, "id = ?", templateID).Error; err != nil { return nil, err }
	// build components simple body order per map iteration
	components := []whatsapp.TemplateComponent{}
	if len(variables) > 0 {
		params := []whatsapp.TemplateParameter{}
		for _, v := range variables {
			params = append(params, whatsapp.TemplateParameter{Type: "text", Text: v})
		}
		components = append(components, whatsapp.TemplateComponent{Type: "body", Parameters: params})
	}
	resp, err := ms.wa.SendTemplateMessage(conv.Customer.WhatsAppID, tpl.Name, tpl.Language, components)
	if err != nil { return nil, fmt.Errorf("send template: %w", err) }
	now := time.Now()
	msg := models.Message{ConversationID: conversationID, WhatsAppID: resp.ID, Type: models.MessageTypeTemplate, Direction: models.MessageDirectionOutbound, Status: models.MessageStatusSent, Content: tpl.Content, TemplateID: &templateID, SentAt: &now}
	if err := ms.db.Create(&msg).Error; err != nil { return nil, err }
	ms.db.Model(&tpl).UpdateColumn("usage_count", gorm.Expr("usage_count + 1"))
	conv.LastMessageAt = &now
	_ = ms.db.Save(&conv)
	return &msg, nil
}
