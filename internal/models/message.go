package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MessageType string

const (
	MessageTypeText     MessageType = "text"
	MessageTypeImage    MessageType = "image"
	MessageTypeDocument MessageType = "document"
	MessageTypeAudio    MessageType = "audio"
	MessageTypeVideo    MessageType = "video"
	MessageTypeSticker  MessageType = "sticker"
	MessageTypeLocation MessageType = "location"
	MessageTypeContact  MessageType = "contact"
	MessageTypeTemplate MessageType = "template"
)

type MessageDirection string

const (
	MessageDirectionInbound  MessageDirection = "inbound"
	MessageDirectionOutbound MessageDirection = "outbound"
)

type MessageStatus string

const (
	MessageStatusSent      MessageStatus = "sent"
	MessageStatusDelivered MessageStatus = "delivered"
	MessageStatusRead      MessageStatus = "read"
	MessageStatusFailed    MessageStatus = "failed"
	MessageStatusPending   MessageStatus = "pending"
)

type Message struct {
	ID             uuid.UUID        `json:"id" gorm:"type:char(36);primaryKey"`
	ConversationID uuid.UUID        `json:"conversation_id" gorm:"type:char(36);index;not null"`
	WhatsAppID     string           `json:"whatsapp_id" gorm:"uniqueIndex"`
	Type           MessageType      `json:"type" gorm:"type:enum('text','image','document','audio','video','sticker','location','contact','template');not null"`
	Direction      MessageDirection `json:"direction" gorm:"type:enum('inbound','outbound');not null"`
	Status         MessageStatus    `json:"status" gorm:"type:enum('sent','delivered','read','failed','pending');default:'pending'"`
	Content        string           `json:"content" gorm:"type:text"`
	MediaURL       string           `json:"media_url"`
	MediaMimeType  string           `json:"media_mime_type"`
	MediaSize      int64            `json:"media_size"`
	FileName       string           `json:"file_name"`
	Caption        string           `json:"caption"`
	Latitude       float64          `json:"latitude"`
	Longitude      float64          `json:"longitude"`
	LocationName   string           `json:"location_name"`
	ContactName    string           `json:"contact_name"`
	ContactPhone   string           `json:"contact_phone"`
	QuotedID       *uuid.UUID       `json:"quoted_id" gorm:"type:char(36);index"`
	TemplateID     *uuid.UUID       `json:"template_id" gorm:"type:char(36);index"`
	SentAt         *time.Time       `json:"sent_at"`
	DeliveredAt    *time.Time       `json:"delivered_at"`
	ReadAt         *time.Time       `json:"read_at"`
	CreatedAt      time.Time        `json:"created_at"`
	UpdatedAt      time.Time        `json:"updated_at"`

	// Relationships
	Conversation Conversation `json:"conversation,omitempty" gorm:"foreignKey:ConversationID"`
	QuotedMessage *Message    `json:"quoted_message,omitempty" gorm:"foreignKey:QuotedID"`
	Template     *Template   `json:"template,omitempty" gorm:"foreignKey:TemplateID"`
}

func (m *Message) BeforeCreate(tx *gorm.DB) (err error) {
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}
	return
}