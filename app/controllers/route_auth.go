package controllers

import (
	"TodoApp/app/models"
	"log"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func getSignup(c *gin.Context) {
	_, span := tracer.Start(c.Request.Context(), "ユーザ登録画面取得")
	defer span.End()

	log.Println("ユーザ登録画面取得")
	generateHTML(c, nil, "signup", "layout", "signup", "public_navbar")
}

func postSignup(c *gin.Context) {
	_, span := tracer.Start(c.Request.Context(), "ユーザ登録")
	defer span.End()

	err := c.Request.ParseForm()
	if err != nil {
		log.Println(err)
	}
	user := models.User{
		Name:     c.Request.PostFormValue("name"),
		Email:    c.Request.PostFormValue("email"),
		PassWord: c.Request.PostFormValue("password"),
	}
	if err := user.CreateUser(c); err != nil {
		log.Println(err)
	}

	UserId := c.PostForm("email")
	log.Println("ログイン処理")
	login(c, UserId)

	_, span = tracer.Start(c.Request.Context(), "TODO画面にリダイレクト")
	defer span.End()

	log.Println("TODO画面にリダイレクト")
	c.Redirect(http.StatusMovedPermanently, "/menu/todos")
}

func login(c *gin.Context, UserId string) {
	_, span := tracer.Start(c.Request.Context(), "ログイン")
	defer span.End()

	session := sessions.Default(c)
	session.Set("UserId", UserId)
	session.Save()
	log.Println("ログイン")
}

func logout(c *gin.Context) {
	_, span := tracer.Start(c.Request.Context(), "ログアウト")
	defer span.End()

	session := sessions.Default(c)
	session.Clear()
	session.Save()
	log.Println("ログアウト")
}

func getLogin(c *gin.Context) {
	_, span := tracer.Start(c.Request.Context(), "ログイン画面取得")
	defer span.End()

	log.Println("ログイン画面取得")
	generateHTML(c, nil, "login", "layout", "login", "public_navbar")
}

func postLogin(c *gin.Context) {
	user, err := models.GetUserByEmail(c, c.Request.PostFormValue("email"))
	if err != nil {
		log.Println(err)
		log.Println("ユーザがいません")
		c.Redirect(http.StatusFound, "/login")
	} else if user.PassWord == models.Encrypt(c, c.Request.PostFormValue("password")) {
		UserId := c.PostForm("email")
		log.Println("ログイン処理")
		login(c, UserId)
		c.SetCookie("UserId", user.Email, 60, "/", "localhost", false, true)
		// index(c)
		c.Redirect(http.StatusMovedPermanently, "/menu/todos")
	} else {
		log.Println("PW が間違っています")
		c.Redirect(http.StatusFound, "/login")
	}
}

func getLogout(c *gin.Context) {
	_, span := tracer.Start(c.Request.Context(), "ログアウト")
	defer span.End()

	logout(c)

	_, span = tracer.Start(c.Request.Context(), "TOP画面にリダイレクト")
	defer span.End()

	log.Println("TOP画面にリダイレクト")
	c.Redirect(http.StatusMovedPermanently, "/")
}
