package controllers

import (
	"os"
	"strconv"
	"time"

	"github.com/AYGA2K/chat_app/db"
	"github.com/AYGA2K/chat_app/models"
	"github.com/golang-jwt/jwt/v4"
	"github.com/joho/godotenv"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"golang.org/x/crypto/bcrypt"
)

type (
	UserController struct {
		Ctx iris.Context
	}
	User struct {
		Name     string `json:"name,omitempty"`
		Email    string `json:"email"`
		Password string `json:"password,omitempty"`
		ID       uint   `json:"id,omitempty"`
	}
)

func (u *UserController) BeforeActivation(b mvc.BeforeActivation) {
	b.Handle("GET", "/{id:string}", "GetById")
}

func (u *UserController) Get() {
	users := []models.User{}
	if err := db.Database.Db.Find(&users).Error; err != nil {
		u.Ctx.StatusCode(iris.StatusBadRequest)
		u.Ctx.JSON(iris.Map{
			"error": err.Error(),
		})
		return
	}
	u.Ctx.StatusCode(iris.StatusOK)
	u.Ctx.JSON(users)
}

func (u *UserController) GetById() {
	userID := u.Ctx.Params().Get("id")
	user := models.User{}
	if err := db.Database.Db.Find(&user, "id=?", userID).Error; err != nil {
		u.Ctx.StatusCode(iris.StatusBadRequest)
		u.Ctx.JSON(iris.Map{
			"error": err.Error(),
		})
		return
	}
	u.Ctx.StatusCode(iris.StatusOK)
	u.Ctx.JSON(user)
}

func (u *UserController) PostSignup() {
	user := models.User{}
	u.Ctx.ReadJSON(&user)
	password := []byte(user.Password)
	hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		u.Ctx.StatusCode(iris.StatusInternalServerError)
		u.Ctx.JSON(iris.Map{
			"error": err.Error(),
		})
		return
	}

	user.Password = string(hashedPassword)
	if err := db.Database.Db.Create(&user).Error; err != nil {
		u.Ctx.StatusCode(iris.StatusBadRequest)
		u.Ctx.JSON(iris.Map{
			"error": err.Error(),
		})
		return
	}

	responseUser := User{
		Email: user.Email,
		Name:  user.Name,
		ID:    user.ID,
	}

	u.Ctx.StatusCode(iris.StatusOK)
	u.Ctx.JSON(responseUser)
}

func (u *UserController) PostLogin() {
	reqUser := User{}
	user := models.User{}
	u.Ctx.ReadJSON(&reqUser)
	if err := db.Database.Db.Find(&user, "email=?", reqUser.Email).Error; err != nil {
		u.Ctx.StatusCode(iris.StatusBadRequest)
		u.Ctx.JSON(iris.Map{
			"error": err.Error(),
		})
		return
	}
	if user.Email == "" {
		u.Ctx.StatusCode(iris.StatusBadRequest)
		u.Ctx.JSON(iris.Map{
			"error": "user not found",
		})
		return

	}
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(reqUser.Password))
	if err != nil {
		u.Ctx.StatusCode(iris.StatusUnauthorized)
		u.Ctx.JSON(iris.Map{
			"error": err.Error(),
		})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    strconv.Itoa(int(user.ID)),
		"email": user.Email,
		"exp":   time.Now().Add(24 * time.Hour).Unix(),
	})

	if err := godotenv.Load(); err != nil {
		u.Ctx.StatusCode(iris.StatusInternalServerError)
		u.Ctx.JSON(iris.Map{
			"error": err,
		})

	}

	access := os.Getenv("ACCESS_KEY")
	tokenString, err := token.SignedString([]byte(access))
	if err != nil {
		u.Ctx.StatusCode(iris.StatusUnauthorized)
		u.Ctx.JSON(iris.Map{
			"error": err,
		})
	}
	u.Ctx.SetCookieKV("jwt", tokenString, iris.CookieExpires(time.Duration(24*time.Hour)), iris.CookiePath("/"), iris.CookieHTTPOnly(true))

	u.Ctx.StatusCode(iris.StatusOK)
	u.Ctx.JSON(iris.Map{
		"userID":   user.ID,
		"userName": user.Name,
		"message":  "login successfull",
	})
}

func (u *UserController) GetLogout() {
	u.Ctx.SetCookieKV("jwt", "", iris.CookieExpires(time.Duration(-1*time.Hour)), iris.CookiePath("/"), iris.CookieHTTPOnly(true))
	u.Ctx.StatusCode(iris.StatusOK)
	u.Ctx.JSON(iris.Map{
		"message": "logout successfull",
	})
}
