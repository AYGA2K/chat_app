package models

import (
	"time"

	"gorm.io/gorm"
)

type Room struct {
	gorm.Model
	CreatedAt time.Time
	Name      string    `json:"name"`
	Users     []*User   `gorm:"many2many:user_rooms;"`
	Messages  []Message `json:"messages"`
	ID        uint      `json:"id" gorm:"primaryKey"`
}
