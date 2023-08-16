package main

import (
	"github.com/AYGA2K/chat_app/controllers"
	"github.com/AYGA2K/chat_app/db"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

func main() {
	db.ConnectDb()
	app := iris.New().Configure(iris.WithFireMethodNotAllowed)
	mvc.Configure(app.Party("/ws"), wsSetup)
	mvc.Configure(app.Party("/user"), userSetup)
	app.Listen(":8080")
}

func wsSetup(app *mvc.Application) {
	app.Handle(new(controllers.RoomController))
}

func userSetup(app *mvc.Application) {
	app.Handle(new(controllers.UserController))
}
