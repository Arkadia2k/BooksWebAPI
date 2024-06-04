package db

import (
	"database/sql"
	"fmt"
	_ "modernc.org/sqlite"
)

var db *sql.DB

func GetDB() *sql.DB {
	return db
}

func Init() {
	var err error
	db, err = sql.Open("sqlite", "books.db")
	if err != nil {
		fmt.Println(err)
		fmt.Println("===============")
		panic("No database connection")
	}
	db.SetMaxOpenConns(10)
	prepDatabaseTable()
}

func prepDatabaseTable() {
	createBook := ` 
	 CREATE TABLE IF NOT EXISTS books(
	 id INTEGER PRIMARY KEY AUTOINCREMENT,
	 booktitle TEXT,
	 isbn INTEGER,
	 bookauthor TEXT,
	 releasedate INTEGER
	)`

	_, err := db.Exec(createBook)
	if err != nil {
		fmt.Println(err)
		fmt.Println("==================")
		panic("Cannot create a book")
	}
}
