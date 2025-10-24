package services

import (
	"time"
	"whatsapp-crm/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ConversationService struct { db *gorm.DB }

func NewConversationService(db *gorm.DB) *ConversationService { return &ConversationService{db: db} }

func (cs *ConversationService) Create(customerID uuid.UUID) (*models.Conversation, error) {
	conv := models.Conversation{CustomerID: customerID, Status: models.ConversationStatusOpen, Priority: models.PriorityMedium}
	if err := cs.db.Create(&conv).Error; err != nil { return nil, err }
	return &conv, nil
}

func (cs *ConversationService) Assign(conversationID, agentID uuid.UUID) error {
	now := time.Now()
	return cs.db.Model(&models.Conversation{}).Where("id = ?", conversationID).Updates(map[string]interface{}{"agent_id": agentID, "status": models.ConversationStatusAssigned, "assigned_at": &now}).Error
}
