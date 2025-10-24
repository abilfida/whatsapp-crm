package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ConversationStatus string

const (
	ConversationStatusOpen     ConversationStatus = "open"
	ConversationStatusAssigned ConversationStatus = "assigned"
	ConversationStatusClosed   ConversationStatus = "closed"
	ConversationStatusPending  ConversationStatus = "pending"
)

type ConversationPriority string

const (
	PriorityLow    ConversationPriority = "low"
	PriorityMedium ConversationPriority = "medium"
	PriorityHigh   ConversationPriority = "high"
	PriorityUrgent ConversationPriority = "urgent"
)

type Conversation struct {
	ID               uuid.UUID            `json:"id" gorm:"type:char(36);primaryKey"`
	CustomerID       uuid.UUID            `json:"customer_id" gorm:"type:char(36);index;not null"`
	AgentID          *uuid.UUID           `json:"agent_id" gorm:"type:char(36);index"`
	Status           ConversationStatus   `json:"status" gorm:"type:enum('open','assigned','closed','pending');default:'open'"`
	Priority         ConversationPriority `json:"priority" gorm:"type:enum('low','medium','high','urgent');default:'medium'"`
	Subject          string               `json:"subject"`
	Tags             string               `json:"tags" gorm:"type:text"`
	Notes            string               `json:"notes" gorm:"type:text"`
	LastMessageAt    *time.Time           `json:"last_message_at"`
	AssignedAt       *time.Time           `json:"assigned_at"`
	ClosedAt         *time.Time           `json:"closed_at"`
	ResponseTime     int                  `json:"response_time" gorm:"comment:'Response time in seconds'"`
	ResolutionTime   int                  `json:"resolution_time" gorm:"comment:'Resolution time in seconds'"`
	UnreadCount      int                  `json:"unread_count" gorm:"default:0"`
	CreatedAt        time.Time            `json:"created_at"`
	UpdatedAt        time.Time            `json:"updated_at"`

	// Relationships
	Customer Customer `json:"customer,omitempty" gorm:"foreignKey:CustomerID"`
	Agent    *User    `json:"agent,omitempty" gorm:"foreignKey:AgentID"`
	Messages []Message `json:"messages,omitempty" gorm:"foreignKey:ConversationID"`
}

func (c *Conversation) BeforeCreate(tx *gorm.DB) (err error) {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return
}