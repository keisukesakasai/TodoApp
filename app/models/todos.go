package models

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

type Todo struct {
	ID        int
	Content   string
	UserID    int
	CreatedAt time.Time
}

func (u *User) CreateTodo(c *gin.Context, content string) (err error) {
	_, span := tracer.Start(c.Request.Context(), "CRUD : CreateTodo")
	defer span.End()

	cmd := `insert into todos (
		content, 
		user_id, 
		created_at) values ($1, $2, $3)`

	_, err = Db.Exec(cmd, content, u.ID, time.Now())
	if err != nil {
		log.Fatalln(err)
	}

	return err
}

func GetTodo(c *gin.Context, id int) (todo Todo, err error) {
	_, span := tracer.Start(c.Request.Context(), "CRUD : GetTodo")
	defer span.End()

	cmd := `select id, content, user_id, created_at from todos
	where id = $1`
	todo = Todo{}

	err = Db.QueryRow(cmd, id).Scan(
		&todo.ID,
		&todo.Content,
		&todo.UserID,
		&todo.CreatedAt)

	return todo, err
}

func GetTodos(c *gin.Context) (todos []Todo, err error) {
	_, span := tracer.Start(c.Request.Context(), "CRUD : GetTodos")
	defer span.End()

	cmd := `select id, content, user_id, created_at from todos`
	rows, err := Db.Query(cmd)
	if err != nil {
		log.Fatalln(err)
	}
	for rows.Next() {
		var todo Todo
		err = rows.Scan(&todo.ID,
			&todo.Content,
			&todo.UserID,
			&todo.CreatedAt)
		if err != nil {
			log.Fatalln(err)
		}
		todos = append(todos, todo)
	}
	rows.Close()

	return todos, err
}

func (u *User) GetTodosByUser(c *gin.Context) (todos []Todo, err error) {
	_, span := tracer.Start(c.Request.Context(), "CRUD : GetTodosByUser")
	defer span.End()

	cmd := `select id, content, user_id, created_at from todos
	where user_id = $1`

	rows, err := Db.Query(cmd, u.ID)
	if err != nil {
		log.Fatalln(err)
	}

	for rows.Next() {
		var todo Todo
		err = rows.Scan(&todo.ID,
			&todo.Content,
			&todo.UserID,
			&todo.CreatedAt)
		if err != nil {
			log.Fatalln(err)
		}
		todos = append(todos, todo)
	}
	rows.Close()

	return todos, err
}

func (t *Todo) UpdateTodo(c *gin.Context) error {
	_, span := tracer.Start(c.Request.Context(), "CRUD : UpdateTodo")
	defer span.End()

	cmd := `update todos set content = $1, user_id = $2 
	where id = $3`
	_, err = Db.Exec(cmd, t.Content, t.UserID, t.ID)
	if err != nil {
		log.Fatalln(err)
	}
	return err
}

func (t *Todo) DeleteTodo(c *gin.Context) error {
	_, span := tracer.Start(c.Request.Context(), "CRUD : DeleteTodo")
	defer span.End()

	cmd := `delete from todos where id = $1`
	_, err = Db.Exec(cmd, t.ID)
	if err != nil {
		log.Fatalln(err)
	}
	return err
}
