package controllers

import (
	"github.com/AYGA2K/chat_app/db"
	"github.com/AYGA2K/chat_app/models"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

type RoomController struct {
	Ctx iris.Context
}

func (r *RoomController) BeforeActivation(b mvc.BeforeActivation) {
	b.Handle("POST", "/create_room/{id:string}", "CreateRoom")
	b.Handle("POST", "/join_room/{id:string}", "JoinRoom")
}

func (r *RoomController) CreateRoom() {
	userID := r.Ctx.Params().Get("id")
	room := models.Room{}
	r.Ctx.ReadJSON(&room)
	if err := db.Database.Db.Create(&room).Error; err != nil {
		r.Ctx.StatusCode(iris.StatusInternalServerError)
		r.Ctx.JSON(iris.Map{
			"error": err.Error(),
		})
		return
	}
	user := models.User{}
	if err := db.Database.Db.Find(&user, "id=?", userID).Error; err != nil {
		r.Ctx.StatusCode(iris.StatusInternalServerError)
		r.Ctx.JSON(iris.Map{
			"error": err.Error(),
		})
		return
	}
	if err := db.Database.Db.Model(&user).Association("Rooms").Append(&room); err != nil {
		r.Ctx.StatusCode(iris.StatusInternalServerError)
		r.Ctx.JSON(iris.Map{
			"error": err.Error(),
		})
		return
	}
	r.Ctx.StatusCode(iris.StatusOK)
	r.Ctx.JSON(iris.Map{
		"message": "room created successfully ",
	})
}

func (r *RoomController) JoinRoom() {
	room := models.Room{}
	r.Ctx.ReadJSON(&room)
	if err := db.Database.Db.Find(&room).Error; err != nil {
		r.Ctx.StatusCode(iris.StatusInternalServerError)
		r.Ctx.JSON(iris.Map{
			"error": err.Error(),
		})
	}
	userID := r.Ctx.Params().Get("id")
	user := models.User{}
	if err := db.Database.Db.Find(&user, "id=?", userID).Error; err != nil {
		r.Ctx.StatusCode(iris.StatusInternalServerError)
		r.Ctx.JSON(iris.Map{
			"error": err.Error(),
		})
		return
	}
	if err := db.Database.Db.Model(&user).Association("Rooms").Append(&room); err != nil {
		r.Ctx.StatusCode(iris.StatusInternalServerError)
		r.Ctx.JSON(iris.Map{
			"error": err.Error(),
		})
		return
	}

	r.Ctx.StatusCode(iris.StatusOK)
	r.Ctx.JSON(iris.Map{
		"message": "joined room successfully ",
	})
}
