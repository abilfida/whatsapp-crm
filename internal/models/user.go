package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRole string

const (
	RoleAdmin      UserRole = "admin"
	RoleAgent      UserRole = "agent"
	RoleSupervisor UserRole = "supervisor"
)

type UserStatus string

const (
	UserStatusActive   UserStatus = "active"
	UserStatusInactive UserStatus = "inactive"
	UserStatusOnline   UserStatus = "online"
	UserStatusOffline  UserStatus = "offline"
)

type User struct {
	ID        uuid.UUID  `json:"id" gorm:"type:char(36);primaryKey"`
	Email     string     `json:"email" gorm:"uniqueIndex;not null"`
	Password  string     `json:"-" gorm:"not null"`
	Name      string     `json:"name" gorm:"not null"`
	Phone     string     `json:"phone" gorm:"index"`
	Role      UserRole   `json:"role" gorm:"type:enum('admin','agent','supervisor');default:'agent'"`
	Status    UserStatus `json:"status" gorm:"type:enum('active','inactive','online','offline');default:'active'"`
	Avatar    string     `json:"avatar"`
	LastLogin *time.Time `json:"last_login"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	Conversations []Conversation `json:"conversations,omitempty" gorm:"foreignKey:AgentID"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return
}