package controllers

import (
	"TodoApp/app/models"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func top(c *gin.Context) {
	_, span := tracer.Start(c.Request.Context(), "TOP画面取得")
	defer span.End()

	log.Println("TOP画面取得")
	generateHTML(c, "hello", "top", "layout", "top", "public_navbar")
}

func index(c *gin.Context) {
	_, span := tracer.Start(c.Request.Context(), "TODO画面取得")
	defer span.End()

	UserId, _ := c.Get("UserId")
	user, err := models.GetUserByEmail(c, UserId.(string))
	if err != nil {
		log.Println(err)
	}

	// ユーザの Todo を取得
	todos, _ := user.GetTodosByUser(c)
	user.Todos = todos

	log.Println("TODO画面取得")
	generateHTML(c, user, "index", "layout", "private_navbar", "index")
}

func todoNew(c *gin.Context) {
	_, span := tracer.Start(c.Request.Context(), "TODO作成画面取得")
	defer span.End()

	log.Println("TODO作成画面取得")
	generateHTML(c, nil, "todoNew", "layout", "private_navbar", "todo_new")
}

func todoSave(c *gin.Context) {
	_, span := tracer.Start(c.Request.Context(), "TODO保存")
	defer span.End()

	UserId, _ := c.Get("UserId")
	user, err := models.GetUserByEmail(c, UserId.(string))
	if err != nil {
		log.Println(err)
	}

	content := c.Request.PostFormValue("content")
	if err := user.CreateTodo(c, content); err != nil {
		log.Println(err)
	}
	log.Println("TODO保存")

	_, span = tracer.Start(c.Request.Context(), "TODO画面にリダイレクト")
	defer span.End()

	log.Println("TODO画面にリダイレクト")
	c.Redirect(http.StatusFound, "/menu/todos")
}

func todoEdit(c *gin.Context, id int) {
	_, span := tracer.Start(c.Request.Context(), "TODO編集画面取得")
	defer span.End()

	err := c.Request.ParseForm()
	if err != nil {
		log.Println(err)
	}

	UserId, _ := c.Get("UserId")
	_, err = models.GetUserByEmail(c, UserId.(string))
	if err != nil {
		log.Println(err)
	}

	t, err := models.GetTodo(c, id)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("TODO編集画面取得")
	generateHTML(c, t, "todoEdit", "layout", "private_navbar", "todo_edit")
}

func todoUpdate(c *gin.Context, id int) {
	_, span := tracer.Start(c.Request.Context(), "TODO更新")
	defer span.End()

	err := c.Request.ParseForm()
	if err != nil {
		log.Println(err)
	}

	UserId, _ := c.Get("UserId")
	user, err := models.GetUserByEmail(c, UserId.(string))
	if err != nil {
		log.Println(err)
	}

	content := c.Request.PostFormValue("content")
	t := &models.Todo{ID: id, Content: content, UserID: user.ID}
	if err := t.UpdateTodo(c); err != nil {
		log.Println(err)
	}
	log.Println("TODO更新")

	_, span = tracer.Start(c.Request.Context(), "TODO画面にリダイレクト")
	defer span.End()

	log.Println("TODO画面にリダイレクト")
	c.Redirect(http.StatusFound, "/menu/todos")
}

func todoDelete(c *gin.Context, id int) {
	_, span := tracer.Start(c.Request.Context(), "TODO削除")
	defer span.End()

	t, err := models.GetTodo(c, id)
	if err != nil {
		log.Println(err)
	}

	if err := t.DeleteTodo(c); err != nil {
		log.Println(err)
	}
	log.Println("TODO削除")

	_, span = tracer.Start(c.Request.Context(), "TODO画面にリダイレクト")
	defer span.End()

	log.Println("TODO画面にリダイレクト")
	c.Redirect(http.StatusFound, "/menu/todos")
}
