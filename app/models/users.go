package models

import (
	"fmt"
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

type Session struct {
	ID        int
	UUID      string
	Email     string
	UserID    int
	CreatedAt time.Time
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
	_, span := tracer.Start(c.Request.Context(), "GetUser")
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
	_, span := tracer.Start(c.Request.Context(), "UpdateUser")
	defer span.End()

	cmd := `update users set name = $1, email = $2 where id = $3`
	_, err = Db.Exec(cmd, u.Name, u.Email, u.ID)
	if err != nil {
		log.Fatalln(err)
	}
	return err
}

func (u *User) DeleteUser(c *gin.Context) (err error) {
	_, span := tracer.Start(c.Request.Context(), "DeleteUser")
	defer span.End()

	cmd := `delete from users where id = $1`
	_, err = Db.Exec(cmd, u.ID)
	if err != nil {
		log.Fatalln(err)
	}
	return err
}

func GetUserByEmail(c *gin.Context, email string) (user User, err error) {
	_, span := tracer.Start(c.Request.Context(), "GetUserByEmail")
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

func (u *User) CreateSession(c *gin.Context) (session Session, err error) {
	_, span := tracer.Start(c.Request.Context(), "CRUD : CreateSession")
	defer span.End()

	session = Session{}
	cmd1 := `insert into sessions (
		uuid, 
		email, 
		user_id, 
		created_at) values ($1, $2, $3, $4)`

	_, err = Db.Exec(cmd1, createUUID(c), u.Email, u.ID, time.Now())
	if err != nil {
		log.Println(err)
	}

	_, span = tracer.Start(c.Request.Context(), "CRUD : QuerySession")
	defer span.End()

	cmd2 := `select id, uuid, email, user_id, created_at
	 from sessions where user_id = $1 and email = $2`

	err = Db.QueryRow(cmd2, u.ID, u.Email).Scan(
		&session.ID,
		&session.UUID,
		&session.Email,
		&session.UserID,
		&session.CreatedAt)

	return session, err
}

func (sess *Session) CheckSession(c *gin.Context) (valid bool, err error) {
	_, span := tracer.Start(c.Request.Context(), "CRUD : CheckSession")
	defer span.End()

	cmd := `select id, uuid, email, user_id, created_at
	from sessions where uuid = $1`

	fmt.Println("===checksession===")
	fmt.Println(sess.UUID)
	fmt.Println("===checksession===")

	err = Db.QueryRow(cmd, sess.UUID).Scan(
		&sess.ID,
		&sess.UUID,
		&sess.Email,
		&sess.UserID,
		&sess.CreatedAt)

	fmt.Println("===checksession===")
	fmt.Println(err)
	fmt.Println(sess.ID)
	fmt.Println("===checksession===")

	if err != nil {
		valid = false
		return
	}
	if sess.ID != 0 {
		valid = true
		return
	}
	return valid, err
}

func (sess *Session) DeleteSessionByUUID(c *gin.Context) (err error) {
	_, span := tracer.Start(c.Request.Context(), "DeleteSessionByUUID")
	defer span.End()

	cmd := `delete from sessions where uuid = $1`
	_, err = Db.Exec(cmd, sess.UUID)
	if err != nil {
		log.Fatalln(err)
	}
	return err
}

func (sess *Session) GetUserBySession(c *gin.Context) (user User, err error) {
	_, span := tracer.Start(c.Request.Context(), "GetUserBySession")
	defer span.End()

	user = User{}
	cmd := `select id, uuid, name, email, created_at FROM users
	where id = $1`
	err = Db.QueryRow(cmd, sess.UserID).Scan(
		&user.ID,
		&user.UUID,
		&user.Name,
		&user.Email,
		&user.CreatedAt)

	return user, err
}
