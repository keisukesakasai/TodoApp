package models

import (
	"context"
	"log"
	"time"

	"go.opentelemetry.io/otel"
)

type Todo struct {
	ID        int
	Content   string
	UserID    int
	CreatedAt time.Time
}

func (u *User) CreateTodo(ctx context.Context, content string) (err error) {
	tracer := otel.Tracer("CreateTodo")
	ctx, span := tracer.Start(ctx, "CreateTodo")
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

func GetTodo(ctx context.Context, id int) (todo Todo, err error) {
	tracer := otel.Tracer("GetTodo")
	ctx, span := tracer.Start(ctx, "GetTodo")
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

func GetTodos(ctx context.Context) (todos []Todo, err error) {
	tracer := otel.Tracer("GetTodos")
	ctx, span := tracer.Start(ctx, "GetTodos")
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

func (u *User) GetTodosByUser(ctx context.Context) (todos []Todo, err error) {
	tracer := otel.Tracer("GetTodosByUser")
	ctx, span := tracer.Start(ctx, "GetTodosByUser")
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

func (t *Todo) UpdateTodo(ctx context.Context) error {
	tracer := otel.Tracer("UpdateTodo")
	ctx, span := tracer.Start(ctx, "UpdateTodo")
	defer span.End()

	cmd := `update todos set content = $1, user_id = $2 
	where id = $3`
	_, err = Db.Exec(cmd, t.Content, t.UserID, t.ID)
	if err != nil {
		log.Fatalln(err)
	}
	return err
}

func (t *Todo) DeleteTodo(ctx context.Context) error {
	tracer := otel.Tracer("DeleteTodo")
	ctx, span := tracer.Start(ctx, "DeleteTodo")
	defer span.End()

	cmd := `delete from todos where id = $1`
	_, err = Db.Exec(cmd, t.ID)
	if err != nil {
		log.Fatalln(err)
	}
	return err
}
