package models

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

type User struct {
	ID        int
	UUID      string
	Name      string
	Email     string
	PassWord  string
	CreatedAt time.Time
	Todos     []Todo
}

func (u *User) CreateUser(c *gin.Context) (err error) {
	ctx := c.Request.Context()
	_, span := tracer.Start(ctx, "CRUD : CreateUser")
	defer span.End()

	cmd := `insert into users (
		uuid,
		name,
		email,
		password,
		created_at) values ($1, $2, $3, $4, $5)`

	_, err = Db.Exec(cmd,
		createUUID(c),
		u.Name,
		u.Email,
		Encrypt(c, u.PassWord),
		time.Now())

	if err != nil {
		log.Fatalln(err)
	}
	return err
}

func GetUser(c *gin.Context, id int) (user User, err error) {
	_, span := tracer.Start(c.Request.Context(), "CRUD : GetUser")
	defer span.End()

	user = User{}
	cmd := `select id, uuid, name, email, password, created_at
	from users where id = $1`
	err = Db.QueryRow(cmd, id).Scan(
		&user.ID,
		&user.UUID,
		&user.Name,
		&user.Email,
		&user.PassWord,
		&user.CreatedAt)

	return user, err
}

func (u *User) UpdateUser(c *gin.Context) (err error) {
	_, span := tracer.Start(c.Request.Context(), "CRUD : UpdateUser")
	defer span.End()

	cmd := `update users set name = $1, email = $2 where id = $3`
	_, err = Db.Exec(cmd, u.Name, u.Email, u.ID)
	if err != nil {
		log.Fatalln(err)
	}
	return err
}

func (u *User) DeleteUser(c *gin.Context) (err error) {
	_, span := tracer.Start(c.Request.Context(), "CRUD : DeleteUser")
	defer span.End()

	cmd := `delete from users where id = $1`
	_, err = Db.Exec(cmd, u.ID)
	if err != nil {
		log.Fatalln(err)
	}
	return err
}

func GetUserByEmail(c *gin.Context, email string) (user User, err error) {
	_, span := tracer.Start(c.Request.Context(), "CRUD : GetUserByEmail")
	defer span.End()

	user = User{}
	cmd := `select id, uuid, name, email, password, created_at
	from users where email = $1`
	err = Db.QueryRow(cmd, email).Scan(
		&user.ID,
		&user.UUID,
		&user.Name,
		&user.Email,
		&user.PassWord,
		&user.CreatedAt)

	return user, err
}
