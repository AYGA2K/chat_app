package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name     string        `json:"name"`
	Email    string        `json:"email" gorm:"unique"`
	Password string        `json:"password"`
	Rooms    []*Room       `gorm:"many2many:user_rooms;"`
	Message  chan *Message `gorm:"-"` // this will be ignored during migration
	Messages []Message     `json:"messages"`
	ID       uint          `json:"id" gorm:"primaryKey"`
}
