package main

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func getUsers(c *gin.Context) {
	file, err := os.Open("users.json")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cannot open users.json"})
		return
	}
	defer file.Close()
	var users []User
	if err := json.NewDecoder(file).Decode(&users); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cannot decode users.json"})
		return
	}
	c.JSON(http.StatusOK, users)
}

func getUserByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user id"})
		return
	}
	file, err := os.Open("users.json")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cannot open users.json"})
		return
	}
	defer file.Close()
	var users []User
	if err := json.NewDecoder(file).Decode(&users); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cannot decode users.json"})
		return
	}
	for _, u := range users {
		if u.ID == id {
			c.JSON(http.StatusOK, u)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
}

func addUser(c *gin.Context) {
	var newUser User
	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	file, err := os.Open("users.json")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cannot open users.json"})
		return
	}
	var users []User
	if err := json.NewDecoder(file).Decode(&users); err != nil {
		file.Close()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cannot decode users.json"})
		return
	}
	file.Close()
	newUser.ID = getNextUserID(users)
	users = append(users, newUser)
	data, err := json.MarshalIndent(users, "", "  ")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cannot encode users"})
		return
	}
	if err := os.WriteFile("users.json", data, 0644); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cannot write users.json"})
		return
	}
	c.JSON(http.StatusCreated, newUser)
}

func getNextUserID(users []User) int {
	maxID := 0
	for _, u := range users {
		if u.ID > maxID {
			maxID = u.ID
		}
	}
	return maxID + 1
}

func main() {
	r := gin.Default()
	r.GET("/users", getUsers)
	r.GET("/users/:id", getUserByID)
	r.POST("/users", addUser)
	r.Run(":3001")
}
