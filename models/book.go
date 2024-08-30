package models

import (
	"database/sql"
	"example/web-service-gin/database"
	"fmt"
)

// Book represents a book entity with basic information.
type Book struct {
	ID     uint    `json:"id"`
	Title  string  `json:"title"`
	Author string  `json:"author"`
	Price  float64 `json:"price"`
}

// GetBooks retrieves all books from the database.
func GetBooks() ([]Book, error) {
	// Execute the query to retrieve all books.
	rows, err := database.DB.Query("SELECT * FROM Books")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var Books []Book
	// Iterate over the rows and scan each one into a Book struct.
	for rows.Next() {
		var Book Book
		if err := rows.Scan(&Book.ID, &Book.Title, &Book.Author, &Book.Price); err != nil {
			return nil, err
		}
		Books = append(Books, Book)
	}
	return Books, nil
}

// GetBooksById retrieves a single book from the database by its ID.
func GetBooksById(id string) (Book, error) {
	var Book Book
	// Prepare the query to retrieve a book with the specified ID.
	query := "SELECT * FROM Books WHERE id = ?"
	row := database.DB.QueryRow(query, id)
	// Scan the result into the Book struct.
	err := row.Scan(&Book.ID, &Book.Title, &Book.Author, &Book.Price)
	if err == sql.ErrNoRows {
		return Book, fmt.Errorf("Book not found")
	}
	return Book, err
}

// AddBook inserts a new book into the database.
func (a *Book) AddBook() error {
	// Prepare the query to insert a new book into the database.
	query := "INSERT INTO Books (title, Author, price) VALUES (?,?,?)"
	_, err := database.DB.Exec(query, a.Title, a.Author, a.Price)
	return err
}
