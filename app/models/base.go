package models

import (
	"context"
	"crypto/sha1"
	"database/sql"
	"fmt"
	"log"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"go.opentelemetry.io/otel"
	// _ "github.com/mattn/go-sqlite3"
)

var Db *sql.DB

var err error

var tracer = otel.Tracer("models")

const (
	tableNameUser    = "users"
	tableNameTodo    = "todos"
	tableNameSession = "sessions"
)

func init() {
	fmt.Println("Now migration...")
	Db, err = sql.Open("postgres", "host=postgresql.prod.svc.cluster.local port=5432 user=postgres dbname=postgres password=postgres sslmode=disable")
	if err != nil {
		log.Println(err)
	}

	cmdU := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s(
		id serial PRIMARY KEY,
		uuid text NOT NULL UNIQUE,
		name text,
		email text,
		password text,
		created_at timestamp)`, "users")

	Db.Exec(cmdU)
	fmt.Println("Now migration...DONE!!")

	cmdT := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s(
		id serial PRIMARY KEY,
		content text,
		user_id integer,
		created_at timestamp)`, "todos")

	Db.Exec(cmdT)

	cmdS := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s(
		id serial PRIMARY KEY,
		uuid text NOT NULL UNIQUE,
		email text,
		user_id integer,
		created_at timestamp)`, "sessions")

	Db.Exec(cmdS)
	fmt.Println("initializing...DONE!!!!")
}

/*
func init() {
	fmt.Println("initializing...")
	Db, err = sql.Open("sqlite3", config.Config.DbName)
	if err != nil {
		log.Fatalln(err)
	}

	cmdU := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		uuid STRING NOT NULL UNIQUE,
		name STRING,
		email STRING,
		password STRING,
		created_at DATETIME)`, tableNameUser)

	Db.Exec(cmdU)

	cmdT := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		content TEXT,
		user_id INTEGER,
		created_at DATETIME)`, tableNameTodo)

	Db.Exec(cmdT)

	cmdS := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		uuid STRING NOT NULL UNIQUE,
		email STRING,
		user_id INTEGER,
		created_at DATETIME)`, tableNameSession)

	Db.Exec(cmdS)
	fmt.Println("initializing...DONE!!!!")
}
*/

func createUUID(ctx context.Context) (uuidobj uuid.UUID) {
	tracer := otel.Tracer("createUUID")
	ctx, span := tracer.Start(ctx, "createUUID")
	defer span.End()

	uuidobj, _ = uuid.NewUUID()
	return uuidobj
}

func Encrypt(ctx context.Context, plaintext string) (cryptext string) {
	tracer := otel.Tracer("Encrypt")
	ctx, span := tracer.Start(ctx, "Encrypt")
	defer span.End()

	cryptext = fmt.Sprintf("%x", sha1.Sum([]byte(plaintext)))
	return cryptext
}
