package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TemplateCategory string

const (
	TemplateCategoryMarketing     TemplateCategory = "marketing"
	TemplateCategoryUtility       TemplateCategory = "utility"
	TemplateCategoryAuthentication TemplateCategory = "authentication"
)

type TemplateStatus string

const (
	TemplateStatusApproved TemplateStatus = "approved"
	TemplateStatusPending  TemplateStatus = "pending"
	TemplateStatusRejected TemplateStatus = "rejected"
	TemplateStatusActive   TemplateStatus = "active"
	TemplateStatusInactive TemplateStatus = "inactive"
)

type Template struct {
	ID          uuid.UUID        `json:"id" gorm:"type:char(36);primaryKey"`
	Name        string           `json:"name" gorm:"uniqueIndex;not null"`
	Language    string           `json:"language" gorm:"default:'en'"`
	Category    TemplateCategory `json:"category" gorm:"type:enum('marketing','utility','authentication');not null"`
	Status      TemplateStatus   `json:"status" gorm:"type:enum('approved','pending','rejected','active','inactive');default:'active'"`
	Content     string           `json:"content" gorm:"type:text;not null"`
	Header      string           `json:"header" gorm:"type:text"`
	Footer      string           `json:"footer" gorm:"type:text"`
	Buttons     string           `json:"buttons" gorm:"type:json"`
	Variables   string           `json:"variables" gorm:"type:json"`
	MediaURL    string           `json:"media_url"`
	MediaType   string           `json:"media_type"`
	UsageCount  int              `json:"usage_count" gorm:"default:0"`
	CreatedBy   uuid.UUID        `json:"created_by" gorm:"type:char(36);index"`
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at"`

	// Relationships
	Creator  User      `json:"creator,omitempty" gorm:"foreignKey:CreatedBy"`
	Messages []Message `json:"messages,omitempty" gorm:"foreignKey:TemplateID"`
}

func (t *Template) BeforeCreate(tx *gorm.DB) (err error) {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	return
}