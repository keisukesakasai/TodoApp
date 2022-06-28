package controllers

import (
	"TodoApp/app/models"
	"log"
	"net/http"

	"go.opentelemetry.io/otel"
)

func signup(w http.ResponseWriter, r *http.Request) {
	tracer := otel.Tracer("signup")
	if r.Method == "GET" {
		ctx := r.Context()
		ctx, span := tracer.Start(ctx, "signup")
		defer span.End()

		_, err := session(ctx, w, r)
		if err != nil {
			generateHTML(ctx, w, nil, "signup", "layout", "signup", "public_navbar")
		} else {
			ctx := r.Context()
			_, span := tracer.Start(ctx, "redirect")
			defer span.End()
			http.Redirect(w, r, "/todos", http.StatusFound)
		}
	} else if r.Method == "POST" {
		ctx := r.Context()
		ctx, span := tracer.Start(ctx, "signup")
		defer span.End()

		err := r.ParseForm()
		if err != nil {
			log.Println(err)
		}
		user := models.User{
			Name:     r.PostFormValue("name"),
			Email:    r.PostFormValue("email"),
			PassWord: r.PostFormValue("password"),
		}
		if err := user.CreateUser(ctx); err != nil {
			log.Println(err)
		}
		// ctx = r.Context()
		ctx, span = tracer.Start(ctx, "redirect")
		defer span.End()

		r = r.WithContext(ctx)
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

func login(w http.ResponseWriter, r *http.Request) {
	// tracer := otel.Tracer("login")
	ctx := r.Context()
	ctx, span := tracer.Start(ctx, "login")
	defer span.End()

	_, err := session(ctx, w, r)
	if err != nil {
		generateHTML(ctx, w, nil, "login", "layout", "login", "public_navbar")
	} else {
		ctx := r.Context()
		_, span = tracer.Start(ctx, "redirect")
		defer span.End()
		http.Redirect(w, r, "/todos", http.StatusFound)
	}
}

func authenticate(w http.ResponseWriter, r *http.Request) {
	tracer := otel.Tracer("authenticate")
	ctx := r.Context()
	ctx, span := tracer.Start(ctx, "authenticate")
	defer span.End()

	//err := r.ParseForm()
	user, err := models.GetUserByEmail(ctx, r.PostFormValue("email"))
	if err != nil {
		log.Println(err)
		ctx = r.Context()
		ctx, span = tracer.Start(ctx, "redirect")
		defer span.End()
		http.Redirect(w, r, "/login", http.StatusFound)
	}
	if user.PassWord == models.Encrypt(ctx, r.PostFormValue("password")) {
		session, err := user.CreateSession(ctx)
		if err != nil {
			log.Println(err)
		}

		cookie := http.Cookie{
			Name:     "_cookie",
			Value:    session.UUID,
			HttpOnly: true,
		}
		http.SetCookie(w, &cookie)
		ctx = r.Context()
		_, span = tracer.Start(ctx, "redirect")
		defer span.End()
		http.Redirect(w, r, "/", http.StatusFound)
	} else {
		ctx = r.Context()
		_, span = tracer.Start(ctx, "redirect")
		defer span.End()
		http.Redirect(w, r, "/login", http.StatusFound)
	}
}

func logout(w http.ResponseWriter, r *http.Request) {
	tracer := otel.Tracer("logout")
	ctx := r.Context()
	ctx, span := tracer.Start(ctx, "logout")
	defer span.End()

	cookie, err := r.Cookie("_cookie")
	if err != nil {
		log.Println(err)
	}

	if err != http.ErrNoCookie {
		session := models.Session{UUID: cookie.Value}
		session.DeleteSessionByUUID(ctx)
	}
	ctx = r.Context()
	_, span = tracer.Start(ctx, "redirect")
	defer span.End()
	http.Redirect(w, r, "/login", http.StatusFound)
}
