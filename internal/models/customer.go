package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Customer struct {
	ID          uuid.UUID `json:"id" gorm:"type:char(36);primaryKey"`
	Name        string    `json:"name" gorm:"not null"`
	Email       string    `json:"email" gorm:"index"`
	Phone       string    `json:"phone" gorm:"index"`
	WhatsAppID  string    `json:"whatsapp_id" gorm:"uniqueIndex"`
	ProfilePic  string    `json:"profile_pic"`
	Company     string    `json:"company"`
	Address     string    `json:"address"`
	City        string    `json:"city"`
	Country     string    `json:"country"`
	Notes       string    `json:"notes" gorm:"type:text"`
	Tags        string    `json:"tags" gorm:"type:text"`
	LastSeen    *time.Time `json:"last_seen"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	Contact       *Contact       `json:"contact,omitempty" gorm:"foreignKey:CustomerID"`
	Conversations []Conversation `json:"conversations,omitempty" gorm:"foreignKey:CustomerID"`
}

func (c *Customer) BeforeCreate(tx *gorm.DB) (err error) {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return
}