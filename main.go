package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type book struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Author   string `json:"author"`
	Quantity int    `json:"quantity"`
}

var books = []book{
	{ID: "1", Title: "Jack and the Bean Stalk", Author: "James Gadenya", Quantity: 2},
	{ID: "2", Title: "Santos and the Bean Stalk", Author: "Sant Gadenya", Quantity: 3},
	{ID: "3", Title: "Patrick and the Bean Stalk", Author: "Pat Gadenya", Quantity: 4},
}

func main() {
	router := gin.Default()
	router.GET("/books", getAllBooks)
	router.POST("/books", addNewBook)
	router.PUT("/books/:id", updateBook)
	router.PUT("/books/checkout", checkoutBook)
	router.PUT("/books/topUp/:id/:topUp", topupBooks)
	router.GET("books/:id", bookById) // :id is for path param
	router.Run("localhost:8080")
}

func getAllBooks(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, books)
}

func addNewBook(c *gin.Context) {
	var newBook book
	if err := c.BindJSON(&newBook); err != nil {
		return
	}
	books = append(books, newBook)
	c.IndentedJSON(http.StatusCreated, newBook)
}

func getBookById(id string) (*book, error) {
	for _, v := range books {
		if v.ID == id {
			return &v, nil
		}
	}
	return nil, errors.New("book not found")
}

func bookById(c *gin.Context) {
	var id = c.Param("id")
	book, err := getBookById(id)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book not found"})
		return
	}
	c.IndentedJSON(http.StatusOK, book)
}

func updateBook(c *gin.Context) {
	id := c.Param("id")
	var updateBook book
	if err := c.BindJSON(&updateBook); err != nil {
		return
	}
	for i, v := range books {
		if v.ID == id {
			books[i] = updateBook
			c.IndentedJSON(http.StatusOK, gin.H{"message": "Book updated successfully"})
		}
	}

}

func checkoutBook(c *gin.Context) {
	id, ok := c.GetQuery("id")
	if !ok {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Missing parameter id"})

	}
	for i, v := range books {
		if v.ID == id {
			q := v.Quantity
			q--
			v.Quantity = q
			if q == -1 {
				c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Unable to checkout more books"})
			} else {
				books[i] = v
				message := fmt.Sprintf("Books left %d", q)
				c.IndentedJSON(http.StatusOK, gin.H{"message": message})

			}
		}
	}
}

func topupBooks(c *gin.Context) {
	id := c.Param("id")
	topUpAmnt, err := strconv.Atoi(c.Param("topUp"))

	if err != nil {
		return
	}

	for i, v := range books {
		if v.ID == id {
			q := v.Quantity
			q += topUpAmnt
			v.Quantity = q
			books[i] = v
			message := fmt.Sprintf("Books left :%d", q)
			c.IndentedJSON(http.StatusOK, gin.H{"message": message})
		}
	}
}
