package handlers

import (
	"example/web-service-gin/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetBooks(c *gin.Context) {
	Books, err := models.GetBooks()
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, Books)
}

func GetBooksById(c *gin.Context) {
	id := c.Param("id")

	Book, err := models.GetBooksById(id)
	if err != nil {
		if err.Error() == "Book not found" {
			c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book not found"})
		} else {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		}
		return
	}
	c.IndentedJSON(http.StatusOK, Book)
}

func PostBooks(c *gin.Context) {
	var newBook models.Book

	if err := c.BindJSON(&newBook); err != nil {
		return
	}

	if err := newBook.AddBook(); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusCreated, newBook)
}
