package models

import (
	"time"

	"gorm.io/gorm"
)

type Message struct {
	gorm.Model
	CreatedAt time.Time
	Content   string `json:"content"`
	UserId    uint   `json:"userId"`
	RoomID    uint   `json:"roomId"`
	ID        uint   `json:"id" gorm:"primaryKey"`
}
