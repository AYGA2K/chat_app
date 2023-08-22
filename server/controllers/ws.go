package controllers

import (
	"net/http"

	"github.com/AYGA2K/chat_app/ws"
	"github.com/gorilla/websocket"
	"github.com/kataras/iris/v12"
)

type (
	Handler struct {
		hub *ws.Hub
	}
)

func NewHandler(h *ws.Hub) *Handler {
	return &Handler{
		hub: h,
	}
}

type ReqRoom struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (h *Handler) CreateRoom(ctx iris.Context) {
	var room ReqRoom
	ctx.ReadJSON(&room)
	h.hub.Rooms[room.ID] = &ws.Room{
		ID:      room.ID,
		Name:    room.Name,
		Clients: make(map[string]*ws.Client),
	}
	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(room)
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (h *Handler) JoinRoom(ctx iris.Context) {
	conn, err := upgrader.Upgrade(ctx.ResponseWriter(), ctx.Request(), nil)
	if err != nil {
		return
	}
	roomID := ctx.Params().Get("roomId")
	clientID := ctx.URLParam("userId")
	username := ctx.URLParam("username")
	cl := &ws.Client{
		Conn:     conn,
		Message:  make(chan *ws.Message, 10),
		ID:       clientID,
		RoomID:   roomID,
		Username: username,
	}

	m := &ws.Message{
		Content:  "A new user has joined the room",
		RoomID:   roomID,
		Username: username,
	}

	h.hub.Register <- cl
	h.hub.Broadcast <- m

	go cl.WriteMessage()
	cl.ReadMessage(h.hub)
}

type RoomRes struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (h *Handler) GetRooms(ctx iris.Context) {
	rooms := make([]RoomRes, 0)

	for _, r := range h.hub.Rooms {
		rooms = append(rooms, RoomRes{
			ID:   r.ID,
			Name: r.Name,
		})
	}
	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(rooms)
}

type ClientRes struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

func (h *Handler) GetClients(ctx iris.Context) {
	var clients []ClientRes
	roomId := ctx.Params().Get("roomId")
	if _, ok := h.hub.Rooms[roomId]; !ok {
		clients = make([]ClientRes, 0)
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(clients)

	}

	for _, c := range h.hub.Rooms[roomId].Clients {
		clients = append(clients, ClientRes{
			ID:       c.ID,
			Username: c.Username,
		})
	}
	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(clients)
}
