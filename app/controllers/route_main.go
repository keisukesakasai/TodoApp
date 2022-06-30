package controllers

import (
	"TodoApp/app/models"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

func top(c *gin.Context) {
	_, span := tracer.Start(c.Request.Context(), "top")
	defer span.End()

	fmt.Println("===top===")
	fmt.Println(c.Cookie("_cookie"))
	fmt.Println("===top===")

	_, err := session(c)
	fmt.Println("---")
	fmt.Println(err)
	fmt.Println("---")
	if err != nil {
		generateHTML(c, "hello", "top", "layout", "top", "public_navbar")
	} else {
		_, span = tracer.Start(c.Request.Context(), "redirect")
		defer span.End()

		//http.Redirect(w, r, "/todos", http.StatusFound)
		// c.Redirect(http.StatusFound, "/todos")
		index(c)
	}
}

func index(c *gin.Context) {
	_, span := tracer.Start(c.Request.Context(), "index")
	defer span.End()

	sess, err := session(c)
	if err != nil {
		_, span = tracer.Start(c.Request.Context(), "redirect")
		defer span.End()

		// http.Redirect(w, r, "/", http.StatusFound)
		top(c)
	} else {
		user, err := sess.GetUserBySession(c)
		if err != nil {
			log.Println(err)
		}
		todos, _ := user.GetTodosByUser(c)
		user.Todos = todos
		generateHTML(c, user, "index", "layout", "private_navbar", "index")
	}
}

func todoNew(c *gin.Context) {
	_, span := tracer.Start(c.Request.Context(), "todoNew")
	defer span.End()

	_, err := session(c)
	if err != nil {
		_, span = tracer.Start(c.Request.Context(), "redirect")
		defer span.End()

		// http.Redirect(w, r, "/login", http.StatusFound)
		login(c)
	} else {
		generateHTML(c, nil, "todoNew", "layout", "private_navbar", "todo_new")
	}
}

func todoSave(c *gin.Context) {
	_, span := tracer.Start(c.Request.Context(), "todoSave")
	defer span.End()

	sess, err := session(c)
	if err != nil {
		_, span = tracer.Start(c.Request.Context(), "redirect")
		defer span.End()

		// http.Redirect(w, r, "/login", http.StatusFound)
		login(c)
	} else {
		err = c.Request.ParseForm()
		if err != nil {
			log.Println(err)
		}
		user, err := sess.GetUserBySession(c)
		if err != nil {
			log.Println(err)
		}
		content := c.Request.PostFormValue(("content"))
		if err := user.CreateTodo(c, content); err != nil {
			log.Println(err)
		}
		_, span = tracer.Start(c.Request.Context(), "redirect")
		defer span.End()

		// http.Redirect(w, r, "/todos", http.StatusFound)
		index(c)
	}
}

func todoEdit(c *gin.Context, id int) {
	_, span := tracer.Start(c.Request.Context(), "todoEdit")
	defer span.End()

	sess, err := session(c)
	if err != nil {
		// ctx = r.Context()
		_, span = tracer.Start(c.Request.Context(), "redirect")
		defer span.End()

		// http.Redirect(w, r, "/login", http.StatusFound)
		login(c)
	} else {
		err = c.Request.ParseForm()
		if err != nil {
			log.Println(err)
		}
		_, err := sess.GetUserBySession(c)
		if err != nil {
			log.Println(err)

		}
		t, err := models.GetTodo(c, id)
		if err != nil {
			log.Fatalln(err)
		}
		generateHTML(c, t, "todoEdit", "layout", "private_navbar", "todo_edit")
	}
}

func todoUpdate(c *gin.Context, id int) {
	_, span := tracer.Start(c.Request.Context(), "todoUpdate")
	defer span.End()

	sess, err := session(c)
	if err != nil {
		_, span = tracer.Start(c.Request.Context(), "redirect")
		defer span.End()

		// http.Redirect(w, r, "/login", http.StatusFound)
		login(c)
	} else {
		err := c.Request.ParseForm()
		if err != nil {
			log.Println(err)
		}
		user, err := sess.GetUserBySession(c)
		if err != nil {
			log.Println(err)
		}
		content := c.Request.PostFormValue("content")
		t := &models.Todo{ID: id, Content: content, UserID: user.ID}
		if err := t.UpdateTodo(c); err != nil {
			log.Println(err)
		}
		_, span = tracer.Start(c.Request.Context(), "redirect")
		defer span.End()

		//http.Redirect(w, r, "/todos", http.StatusFound)
		index(c)
	}
}

func todoDelete(c *gin.Context, id int) {
	_, span := tracer.Start(c.Request.Context(), "todoDelete")
	defer span.End()

	sess, err := session(c)
	if err != nil {
		_, span = tracer.Start(c.Request.Context(), "redirect")
		defer span.End()

		// http.Redirect(w, r, "/login", http.StatusFound)
		login(c)
	} else {
		_, err := sess.GetUserBySession(c)
		if err != nil {
			log.Println(err)
		}
		t, err := models.GetTodo(c, id)
		if err != nil {
			log.Println(err)
		}
		if err := t.DeleteTodo(c); err != nil {
			log.Println(err)
		}
		_, span = tracer.Start(c.Request.Context(), "redirect")
		defer span.End()

		// http.Redirect(w, r, "/todos", http.StatusFound)
		index(c)
	}
}
