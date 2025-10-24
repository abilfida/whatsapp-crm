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
