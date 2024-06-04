package main

import (
	"database/sql"
	"errors"
	_ "modernc.org/sqlite"
	"net/http"
	"strconv"

	"example/go-project/db"
	"github.com/gin-gonic/gin"
)

type book struct {
	ID          int64  `json:"id"`
	BookTitle   string `json:"booktitle"`
	ISBN        int    `json:"isbn"`
	BookAuthor  string `json:"bookauthor"`
	ReleaseDate int    `json:"releasedate"`
}

func getBooks(c *gin.Context) {
	var books []book
	rows, err := db.GetDB().Query("SELECT id, booktitle, isbn, bookauthor, releasedate FROM books")
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "No database found."})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var b book
		if err := rows.Scan(&b.ID, &b.BookTitle, &b.ISBN, &b.BookAuthor, &b.ReleaseDate); err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "No database found."})
			return
		}
		books = append(books, b)
	}
	c.IndentedJSON(http.StatusOK, books)
}

func bookById(c *gin.Context) {
	id := c.Param("id")
	book, err := getBookById(id)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "No book found."})
		return
	}

	c.IndentedJSON(http.StatusOK, book)
}

func getBookById(id string) (*book, error) {
	bookID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, err
	}

	var b book
	err = db.GetDB().QueryRow("SELECT id, booktitle, isbn, bookauthor, releasedate FROM books WHERE id = ?", bookID).Scan(&b.ID, &b.BookTitle, &b.ISBN, &b.BookAuthor, &b.ReleaseDate)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("No book found.")
		}
		return nil, err
	}
	return &b, nil
}

func createBook(c *gin.Context) {
	var newBook book

	if err := c.BindJSON(&newBook); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid request."})
		return
	}

	res, err := db.GetDB().Exec("INSERT INTO books (booktitle, isbn, bookauthor, releasedate) VALUES (?, ?, ?, ?)", newBook.BookTitle, newBook.ISBN, newBook.BookAuthor, newBook.ReleaseDate)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Database error."})
		return
	}
	id, err := res.LastInsertId()
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Database error."})
		return
	}
	newBook.ID = id
	c.IndentedJSON(http.StatusCreated, newBook)
}

func updateBookById(c *gin.Context) {
	id := c.Param("id")
	bookID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Wrong book ID."})
		return
	}

	var updatedBook book
	if err := c.BindJSON(&updatedBook); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Wrong book ID."})
		return
	}

	updatedBook.ID = bookID

	_, err = db.GetDB().Exec("UPDATE books SET booktitle = ?, isbn = ?, bookauthor = ?, releasedate = ? WHERE id = ?", updatedBook.BookTitle, updatedBook.ISBN, updatedBook.BookAuthor, updatedBook.ReleaseDate, bookID)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Database error."})
		return
	}
	c.IndentedJSON(http.StatusOK, updatedBook)
}

func updateBook(c *gin.Context) {
	var updatedBook book

	if err := c.BindJSON(&updatedBook); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Wrong request."})
		return
	}

	bookID := updatedBook.ID

	_, err := db.GetDB().Exec("UPDATE books SET booktitle = ?, isbn = ?, bookauthor = ?, releasedate = ? WHERE id = ?", updatedBook.BookTitle, updatedBook.ISBN, updatedBook.BookAuthor, updatedBook.ReleaseDate, bookID)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Database error."})
		return
	}
	c.IndentedJSON(http.StatusOK, updatedBook)
}

func deleteBook(c *gin.Context) {
	id := c.Param("id")
	bookID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Wrong book ID."})
		return
	}

	_, err = db.GetDB().Exec("DELETE FROM books WHERE id = ?", bookID)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Database error."})
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"message": "Book was deleted successfully."})
}

func main() {
	db.Init()
	router := gin.Default()
	router.GET("/books", getBooks)
	router.GET("/books/:id", bookById)
	router.POST("/books", createBook)
	router.PUT("/books/:id", updateBookById)
	router.PUT("/books", updateBook)
	router.DELETE("/books/:id", deleteBook)
	router.Run("localhost:8080")
}
