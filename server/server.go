package main

import (
	"github.com/AYGA2K/chat_app/controllers"
	"github.com/AYGA2K/chat_app/db"
	"github.com/AYGA2K/chat_app/ws"
	"github.com/iris-contrib/middleware/cors"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

func main() {
	db.ConnectDb()

	app := iris.New().Configure(iris.WithFireMethodNotAllowed)
	crs := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
	})
	app.UseRouter(crs)
	mvc.Configure(app.Party("/user"), userSetup)

	hub := ws.NewHub()
	wsHandler := controllers.NewHandler(hub)
	go hub.Run()
	app.Any("/ws/:roomId", wsHandler.JoinRoom)
	app.Post("/ws/createRoom", wsHandler.CreateRoom)
	app.Get("/ws/getRooms", wsHandler.GetRooms)
	app.Get("/ws/getClients/:roomId", wsHandler.GetClients)
	app.Listen(":8080")
}

func userSetup(app *mvc.Application) {
	app.Handle(new(controllers.UserController))
}
