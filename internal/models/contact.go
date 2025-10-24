package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ContactStatus string

const (
	ContactStatusValid   ContactStatus = "valid"
	ContactStatusInvalid ContactStatus = "invalid"
	ContactStatusBlocked ContactStatus = "blocked"
)

type Contact struct {
	ID            uuid.UUID     `json:"id" gorm:"type:char(36);primaryKey"`
	CustomerID    uuid.UUID     `json:"customer_id" gorm:"type:char(36);index"`
	WhatsAppID    string        `json:"whatsapp_id" gorm:"uniqueIndex;not null"`
	DisplayName   string        `json:"display_name"`
	PushName      string        `json:"push_name"`
	ProfilePic    string        `json:"profile_pic"`
	Status        ContactStatus `json:"status" gorm:"type:enum('valid','invalid','blocked');default:'valid'"`
	IsGroup       bool          `json:"is_group" gorm:"default:false"`
	GroupSubject  string        `json:"group_subject"`
	GroupDesc     string        `json:"group_desc"`
	LastSeen      *time.Time    `json:"last_seen"`
	CreatedAt     time.Time     `json:"created_at"`
	UpdatedAt     time.Time     `json:"updated_at"`

	// Relationships
	Customer Customer `json:"customer,omitempty" gorm:"foreignKey:CustomerID"`
}

func (c *Contact) BeforeCreate(tx *gorm.DB) (err error) {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return
}