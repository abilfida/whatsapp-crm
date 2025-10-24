package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type WebhookEventType string

const (
	WebhookEventMessage       WebhookEventType = "message"
	WebhookEventMessageStatus WebhookEventType = "message_status"
	WebhookEventPresence      WebhookEventType = "presence"
	WebhookEventTyping        WebhookEventType = "typing"
	WebhookEventContact       WebhookEventType = "contact"
)

type WebhookStatus string

const (
	WebhookStatusReceived  WebhookStatus = "received"
	WebhookStatusProcessed WebhookStatus = "processed"
	WebhookStatusFailed    WebhookStatus = "failed"
	WebhookStatusIgnored   WebhookStatus = "ignored"
)

type WebhookLog struct {
	ID           uuid.UUID        `json:"id" gorm:"type:char(36);primaryKey"`
	EventType    WebhookEventType `json:"event_type" gorm:"type:enum('message','message_status','presence','typing','contact');not null"`
	Status       WebhookStatus    `json:"status" gorm:"type:enum('received','processed','failed','ignored');default:'received'"`
	Payload      string           `json:"payload" gorm:"type:longtext;not null"`
	Response     string           `json:"response" gorm:"type:text"`
	ErrorMessage string           `json:"error_message" gorm:"type:text"`
	ProcessedAt  *time.Time       `json:"processed_at"`
	CreatedAt    time.Time        `json:"created_at"`
	UpdatedAt    time.Time        `json:"updated_at"`
}

func (w *WebhookLog) BeforeCreate(tx *gorm.DB) (err error) {
	if w.ID == uuid.Nil {
		w.ID = uuid.New()
	}
	return
}