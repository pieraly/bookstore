package handlers

import (
	"bytes"
	"encoding/json"
	"example/web-service-gin/database"
	"example/web-service-gin/models"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// SetupRouter initializes the Gin engine and registers the routes for the book handlers.
// This is used for testing the various endpoints related to books.
func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/books", GetBooks)         // Route to get all books
	r.POST("/books", PostBooks)       // Route to create a new book
	r.GET("/books/:id", GetBooksById) // Route to get a book by its ID
	return r
}

// TestGetBooks tests the GET /books endpoint to ensure it returns the correct status and data.
func TestGetBooks(t *testing.T) {
	database.Connect()      // Connect to the test database
	router := SetupRouter() // Set up the router with the book routes

	// Create a new GET request for the /books endpoint.
	req, err := http.NewRequest(http.MethodGet, "/books", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Record the response using httptest.
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	// Assert that the status code is 200 OK.
	assert.Equal(t, http.StatusOK, rr.Code, "Expected status OK, got %v", rr.Code)

	// Unmarshal the response body into a slice of Book objects.
	var books []models.Book
	err = json.Unmarshal(rr.Body.Bytes(), &books)
	if err != nil {
		t.Fatalf("Failed to unmarshal response body: %v", err)
	}

	// Check that the list of books is not nil and has at least 0 books (empty is acceptable).
	assert.NotNil(t, books, "Expected books to be non-nil")
	assert.GreaterOrEqual(t, len(books), 0, "Expected at least 0 books in response")

	// Ensure each book in the response has the required fields (ID, Title, Author, Price).
	for _, book := range books {
		assert.NotNil(t, book.ID, "Expected book ID to be non-nil")
		assert.NotNil(t, book.Title, "Expected book title to be non-nil")
		assert.NotNil(t, book.Author, "Expected book author to be non-nil")
		assert.NotNil(t, book.Price, "Expected book price to be non-nil")
	}
}

// TestPostBooks tests the POST /books endpoint to ensure a book can be created successfully.
func TestPostBooks(t *testing.T) {
	database.Connect()      // Connect to the test database
	router := SetupRouter() // Set up the router with the book routes

	// Create a new book object to be sent in the POST request.
	book := models.Book{
		Title:  "Test Book",
		Author: "Test Author",
		Price:  19.99,
	}

	// Marshal the book object into JSON format.
	jsonBook, err := json.Marshal(book)
	if err != nil {
		t.Fatalf("Failed to marshal book: %v", err)
	}

	// Create a new POST request with the book JSON as the body.
	req, err := http.NewRequest(http.MethodPost, "/books", bytes.NewBuffer(jsonBook))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json") // Set the content type header to application/json

	// Record the response using httptest.
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	// Assert that the status code is 201 Created.
	assert.Equal(t, http.StatusCreated, rr.Code, "Expected status Created, got %v", rr.Code)

	// Unmarshal the response body into a Book object.
	var createdBook models.Book
	err = json.Unmarshal(rr.Body.Bytes(), &createdBook)
	if err != nil {
		t.Fatalf("Failed to unmarshal response body: %v", err)
	}

	// Assert that the created book matches the book sent in the request.
	assert.NotNil(t, createdBook, "Expected createdBook to be non-nil")
	assert.Equal(t, book.Title, createdBook.Title, "Expected title to match")
	assert.Equal(t, book.Author, createdBook.Author, "Expected author to match")
	assert.Equal(t, book.Price, createdBook.Price, "Expected price to match")
}

// TestGetBooksById tests the GET /books/:id endpoint to ensure a book can be retrieved by its ID.
func TestGetBooksById(t *testing.T) {
	database.Connect()      // Connect to the test database
	router := SetupRouter() // Set up the router with the book routes

	bookID := 1 // Use a specific book ID for the test

	// Create a new GET request for the /books/:id endpoint.
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/books/%d", bookID), nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Record the response using httptest.
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	// Assert that the status code is 200 OK.
	assert.Equal(t, http.StatusOK, rr.Code, "Expected status OK, got %v", rr.Code)

	// Unmarshal the response body into a Book object.
	var book models.Book
	err = json.Unmarshal(rr.Body.Bytes(), &book)
	if err != nil {
		t.Fatalf("Failed to unmarshal response body: %v", err)
	}

	// Assert that the retrieved book's ID, title, author, and price are not nil and match the expected values.
	assert.NotNil(t, book, "Expected book to be non-nil")
	assert.Equal(t, uint(bookID), book.ID, "Expected book ID to match")
	assert.NotNil(t, book.Title, "Expected book title to be non-nil")
	assert.NotNil(t, book.Author, "Expected book author to be non-nil")
	assert.NotNil(t, book.Price, "Expected book price to be non-nil")
}
