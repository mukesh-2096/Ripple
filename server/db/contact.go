package db

import (
	"time"
)

type ContactMessage struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"not null" json:"name"`
	Email     string    `gorm:"not null" json:"email"`
	Message   string    `gorm:"not null" json:"message"`
	CreatedAt time.Time `json:"created_at"`
}
