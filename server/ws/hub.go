package ws

import "github.com/AYGA2K/chat_app/models"

type Hub struct {
	Rooms      map[string]*models.Room
	Register   chan *models.User
	Unregister chan *models.User
	Broadcast  chan *models.Message
}

func NewHub() *Hub {
	return &Hub{
		Rooms:      make(map[string]*models.Room),
		Register:   make(chan *models.User),
		Unregister: make(chan *models.User),
		Broadcast:  make(chan *models.Message, 5),
	}
}
