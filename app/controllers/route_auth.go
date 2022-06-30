package controllers

import (
	"TodoApp/app/models"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
)

func getSignup(c *gin.Context) {
	tracer := otel.Tracer("getSignup")
	_, span := tracer.Start(c.Request.Context(), "getSignup")
	defer span.End()

	_, err := session(c)
	if err != nil {
		generateHTML(c, nil, "signup", "layout", "signup", "public_navbar")
	} else {
		_, span := tracer.Start(c.Request.Context(), "redirect")
		defer span.End()

		// c.Redirect(http.StatusFound, "/todos")
		//http.Redirect(w, r, "/todos", http.StatusFound)
		index(c)
	}
}

func postSignup(c *gin.Context) {
	_, span := tracer.Start(c.Request.Context(), "postSignup")
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

	_, span = tracer.Start(c.Request.Context(), "redirect")
	defer span.End()
	top(c)
}

func login(c *gin.Context) {
	_, span := tracer.Start(c.Request.Context(), "login")
	defer span.End()

	_, err := session(c)
	if err != nil {
		generateHTML(c, nil, "login", "layout", "login", "public_navbar")
	} else {
		_, span = tracer.Start(c.Request.Context(), "redirect")
		defer span.End()

		// http.Redirect(w, r, "/todos", http.StatusFound)
		index(c)
	}
}

func authenticate(c *gin.Context) {
	_, span := tracer.Start(c.Request.Context(), "authenticate")
	defer span.End()

	user, err := models.GetUserByEmail(c, c.Request.PostFormValue("email"))
	if err != nil {
		log.Println(err)
		_, span = tracer.Start(c.Request.Context(), "redirect")
		defer span.End()

		// http.Redirect(w, r, "/login", http.StatusFound)
		login(c)
	}
	if user.PassWord == models.Encrypt(c, c.Request.PostFormValue("password")) {
		session, err := user.CreateSession(c)

		log.Println(session)
		log.Println(session.UUID)
		sessionUUID := session.UUID
		log.Println(sessionUUID)

		if err != nil {
			log.Println(err)
		}

		/*
			cookie := http.Cookie{
				Name:     "_cookie",
				Value:    session.UUID,
				HttpOnly: true,
			}
		*/

		cookie := new(http.Cookie)
		cookie.Value = session.UUID
		c.SetSameSite(http.SameSiteNoneMode)

		c.SetCookie("_cookie", cookie.Value, -1, "/", "localhost", true, true)
		c.SetCookie("_cookie", cookie.Value, 3600, "/", "localhost", true, true)
		// http.SetCookie(c.Writer, &cookie)
		fmt.Println("===setcookie===")
		fmt.Println(c.Cookie("_cookie"))
		fmt.Println(c.Cookie("_cookie"))
		fmt.Println(c.Cookie("_cookie"))
		fmt.Println("===setcookie===")
		// ctx = r.Context()

		_, span = tracer.Start(c.Request.Context(), "redirect")
		defer span.End()

		//http.Redirect(w, r, "/", http.StatusFound)
		top(c)
	} else {
		// ctx = r.Context()
		_, span = tracer.Start(c.Request.Context(), "redirect")
		defer span.End()
		// http.Redirect(w, r, "/login", http.StatusFound)
		login(c)
	}
}

func logout(c *gin.Context) {
	_, span := tracer.Start(c.Request.Context(), "logout")
	defer span.End()

	cookie, err := c.Request.Cookie("_cookie")
	if err != nil {
		log.Println(err)
	}

	if err != http.ErrNoCookie {
		session := models.Session{UUID: cookie.Value}
		session.DeleteSessionByUUID(c)
	}
	_, span = tracer.Start(c.Request.Context(), "redirect")
	defer span.End()

	// http.Redirect(w, r, "/login", http.StatusFound)
	login(c)
}
