package models

import (
	"context"
	"log"
	"time"

	"go.opentelemetry.io/otel"
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

func (u *User) CreateUser(ctx context.Context) (err error) {
	tracer := otel.Tracer("CreateUser")
	ctx, span := tracer.Start(ctx, "CreateUser")
	defer span.End()

	cmd := `insert into users (
		uuid,
		name,
		email,
		password,
		created_at) values ($1, $2, $3, $4, $5)`

	_, err = Db.Exec(cmd,
		createUUID(ctx),
		u.Name,
		u.Email,
		Encrypt(ctx, u.PassWord),
		time.Now())

	if err != nil {
		log.Fatalln(err)
	}
	return err
}

func GetUser(ctx context.Context, id int) (user User, err error) {
	tracer := otel.Tracer("GetUser")
	ctx, span := tracer.Start(ctx, "GetUser")
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

func (u *User) UpdateUser(ctx context.Context) (err error) {
	tracer := otel.Tracer("UpdateUser")
	ctx, span := tracer.Start(ctx, "UpdateUser")
	defer span.End()

	cmd := `update users set name = $1, email = $2 where id = $3`
	_, err = Db.Exec(cmd, u.Name, u.Email, u.ID)
	if err != nil {
		log.Fatalln(err)
	}
	return err
}

func (u *User) DeleteUser(ctx context.Context) (err error) {
	tracer := otel.Tracer("DeleteUser")
	ctx, span := tracer.Start(ctx, "DeleteUser")
	defer span.End()

	cmd := `delete from users where id = $1`
	_, err = Db.Exec(cmd, u.ID)
	if err != nil {
		log.Fatalln(err)
	}
	return err
}

func GetUserByEmail(ctx context.Context, email string) (user User, err error) {
	tracer := otel.Tracer("GetUserByEmail")
	ctx, span := tracer.Start(ctx, "GetUserByEmail")
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

func (u *User) CreateSession(ctx context.Context) (session Session, err error) {
	tracer := otel.Tracer("CreateSession")

	ctx, span := tracer.Start(ctx, "createsession")
	defer span.End()

	session = Session{}
	cmd1 := `insert into sessions (
		uuid, 
		email, 
		user_id, 
		created_at) values ($1, $2, $3, $4)`

	_, err = Db.Exec(cmd1, createUUID(ctx), u.Email, u.ID, time.Now())
	if err != nil {
		log.Println(err)
	}

	_, span = tracer.Start(ctx, "querysession")
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

func (sess *Session) CheckSession(ctx context.Context) (valid bool, err error) {
	tracer := otel.Tracer("CheckSession")
	ctx, span := tracer.Start(ctx, "CheckSession")
	defer span.End()

	cmd := `select id, uuid, email, user_id, created_at
	from sessions where uuid = $1`

	err = Db.QueryRow(cmd, sess.UUID).Scan(
		&sess.ID,
		&sess.UUID,
		&sess.Email,
		&sess.UserID,
		&sess.CreatedAt)

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

func (sess *Session) DeleteSessionByUUID(ctx context.Context) (err error) {
	tracer := otel.Tracer("DeleteSessionByUUID")
	ctx, span := tracer.Start(ctx, "DeleteSessionByUUID")
	defer span.End()

	cmd := `delete from sessions where uuid = $1`
	_, err = Db.Exec(cmd, sess.UUID)
	if err != nil {
		log.Fatalln(err)
	}
	return err
}

func (sess *Session) GetUserBySession(ctx context.Context) (user User, err error) {
	tracer := otel.Tracer("GetUserBySession")
	ctx, span := tracer.Start(ctx, "GetUserBySession")
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
