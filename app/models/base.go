package models

import (
	"crypto/sha1"
	"database/sql"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	_ "github.com/lib/pq"

	// _ "github.com/mattn/go-sqlite3"
	"go.opentelemetry.io/otel"
)

var Db *sql.DB

var err error
var tracer = otel.Tracer("models")

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

	log.Println("initializing...DONE!!!!")
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
		created_at DATETIME)`, "users")

	Db.Exec(cmdU)

	cmdT := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		content TEXT,
		user_id INTEGER,
		created_at DATETIME)`, "todos")

	Db.Exec(cmdT)

	log.Println("initializing...DONE!!!!")
}
*/

func createUUID(c *gin.Context) (uuidobj uuid.UUID) {
	_, span := tracer.Start(c.Request.Context(), "createUUID")
	defer span.End()

	uuidobj, _ = uuid.NewUUID()
	return uuidobj
}

func Encrypt(c *gin.Context, plaintext string) (cryptext string) {
	_, span := tracer.Start(c.Request.Context(), "Encrypt")
	defer span.End()

	cryptext = fmt.Sprintf("%x", sha1.Sum([]byte(plaintext)))
	return cryptext
}
