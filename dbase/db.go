package dbase

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	db *sql.DB
}

func CheckDB() *sql.DB {
	_, err := os.Stat("dbase/database-sqlite.db")
	if os.IsNotExist(err) {
		createFile()
	}
	var d Database
	d.open("dbase/database-sqlite.db")
	d.createTable()
	return d.db
}

func createFile() {
	file, err := os.Create("dbase/database-sqlite.db")
	if err != nil {
		log.Fatalf("file doesn't create %v", err)
	}
	defer file.Close()
}

func (d *Database) open(file string) {
	var err error
	d.db, err = sql.Open("sqlite3", file)
	if err != nil {
		log.Fatalf("this error is in dbase/open() %v", err)
	}
}

func (d *Database) createTable() {
	_, err := d.db.Exec(`CREATE TABLE IF NOT EXISTS users (
        "id"    INTEGER NOT NULL UNIQUE,
        "login"    TEXT NOT NULL UNIQUE,
        "password"    TEXT NOT NULL,
        "email"    TEXT NOT NULL UNIQUE,
		"user_uuid" TEXT NOT NULL UNIQUE,
        PRIMARY KEY("id" AUTOINCREMENT)
    );`)
	if err != nil {
		log.Println("CAN NOT CREATE TABLE users")
	}

	_, err = d.db.Exec(`CREATE TABLE IF NOT EXISTS posts 
    (
        "id"    INTEGER NOT NULL UNIQUE,
        "title"	TEXT NOT NULL,
        "body" TEXT NOT NULL,
        "date"    DATATIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		"user_uuid" TEXT NOT NULL,
        FOREIGN KEY("user_uuid") REFERENCES "users"("user_uuid"),
        PRIMARY KEY("id" AUTOINCREMENT)
    );`)
	if err != nil {
		log.Println("CAN NOT CREATE TABLE posts")
	}

	_, err = d.db.Exec(`CREATE TABLE IF NOT EXISTS categories (
        "id"    INTEGER NOT NULL UNIQUE,
        "name"    TEXT NOT NULL UNIQUE,
        PRIMARY KEY("id" AUTOINCREMENT)
    );`)
	if err != nil {
		log.Println("CAN NOT CREATE TABLE categories")
	}
	_, err = d.db.Exec(`CREATE TABLE IF NOT EXISTS post_category (
        "id"    INTEGER NOT NULL UNIQUE,
        "post_id"    INT NOT NULL,
		"category_id"    INT NOT NULL,
		FOREIGN KEY("category_id") REFERENCES "categories"("id"),
		FOREIGN KEY("post_id") REFERENCES "posts"("id"),
        PRIMARY KEY("id" AUTOINCREMENT)
    );`)
	if err != nil {
		log.Println("CAN NOT CREATE TABLE post_category")
	}

	_, err = d.db.Exec(`INSERT INTO categories (name)
		VALUES('golang'),
			('python'),
			('java'),
			('c++'),
			('git'),
			('alem school');`)
	if err != nil {
		// log.Printf("%v\n", err)
	}

	_, err = d.db.Exec(`CREATE TABLE  IF NOT EXISTS comments (
        "id"    INTEGER NOT NULL UNIQUE,
        "user_uuid"  TEXT NOT NULL,
        "post_id"    INTEGER NOT NULL,
        "date"    DATATIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
        "comment"    TEXT NOT NULL,
        FOREIGN KEY("user_uuid") REFERENCES "users"("user_uuid"),
        FOREIGN KEY("post_id") REFERENCES "posts"("id"),
        PRIMARY KEY("id" AUTOINCREMENT)
    );`)
	if err != nil {
		log.Println("CAN NOT CREATE TABLE comments")
	}

	_, err = d.db.Exec(`CREATE TABLE IF NOT EXISTS post_likes (
        "id"    INTEGER NOT NULL UNIQUE,
        "reaction"  TEXT NOT NULL,
		"user_uuid" TEXT NOT NULL,
		"post_id" INTEGER NOT NULL,
		FOREIGN KEY("user_uuid") REFERENCES "users"("user_uuid"),
		FOREIGN KEY("post_id") REFERENCES "posts"("post_id"),
        PRIMARY KEY("id" AUTOINCREMENT)
    );`)
	if err != nil {
		log.Println("CAN NOT CREATE TABLE likes")
	}

	_, err = d.db.Exec(`CREATE TABLE IF NOT EXISTS comment_likes (
        "id"    INTEGER NOT NULL UNIQUE,
        "reaction"  TEXT NOT NULL,
        "user_uuid"   TEXT NOT NULL,
        "comment_id"    INTEGER NOT NULL,
        FOREIGN KEY("user_uuid") REFERENCES "users"("user_uuid"),
        FOREIGN KEY("comment_id") REFERENCES "comments"("id"),
        PRIMARY KEY("id" AUTOINCREMENT)
    );`)
	if err != nil {
		log.Println("CAN NOT CREATE TABLE commentslikes")
	}

	_, err = d.db.Exec(`CREATE TABLE IF NOT EXISTS sessions (
        "id"    INTEGER NOT NULL UNIQUE,
        "user_uuid" TEXT NOT NULL,
        "key"        TEXT NOT NULL,
        "date"        DATATIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY("user_uuid") REFERENCES "users"("user_uuid"),
        PRIMARY KEY("id" AUTOINCREMENT)
    );`)
	if err != nil {
		log.Println("CAN NOT CREATE TABLE sessions")
	}
}
