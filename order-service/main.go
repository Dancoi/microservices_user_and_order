package main

import (
	"encoding/json"
	"os"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
)

type Order struct {
	ID     int    `json:"id"`
	UserID int    `json:"userId"`
	Item   string `json:"item"`
}

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func getOrders(c *gin.Context) {
	file, err := os.Open("orders.json")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cannot open orders.json"})
		return
	}
	defer file.Close()
	var orders []Order
	if err := json.NewDecoder(file).Decode(&orders); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cannot decode orders.json"})
		return
	}
	c.JSON(http.StatusOK, orders)
}

func getOrderByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order id"})
		return
	}
	file, err := os.Open("orders.json")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cannot open orders.json"})
		return
	}
	defer file.Close()
	var orders []Order
	if err := json.NewDecoder(file).Decode(&orders); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cannot decode orders.json"})
		return
	}
	for _, o := range orders {
		if o.ID == id {
			user, err := fetchUser(o.UserID)
			if err != nil {
				c.JSON(http.StatusOK, gin.H{"order": o, "user": nil, "user_error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, gin.H{"order": o, "user": user})
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
}

func addOrder(c *gin.Context) {
	var newOrder Order
	if err := c.ShouldBindJSON(&newOrder); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	file, err := os.Open("orders.json")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cannot open orders.json"})
		return
	}
	var orders []Order
	if err := json.NewDecoder(file).Decode(&orders); err != nil {
		file.Close()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cannot decode orders.json"})
		return
	}
	file.Close()
	newOrder.ID = getNextOrderID(orders)
	orders = append(orders, newOrder)
	data, err := json.MarshalIndent(orders, "", "  ")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cannot encode orders"})
		return
	}
	if err := os.WriteFile("orders.json", data, 0644); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cannot write orders.json"})
		return
	}
	c.JSON(http.StatusCreated, newOrder)
}

func getNextOrderID(orders []Order) int {
	maxID := 0
	for _, o := range orders {
		if o.ID > maxID {
			maxID = o.ID
		}
	}
	return maxID + 1
}

func fetchUser(userID int) (*User, error) {
	client := resty.New()
	resp, err := client.R().SetResult(&User{}).Get("http://user-service:3001/users/" + strconv.Itoa(userID))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != 200 {
		return nil, err
	}
	return resp.Result().(*User), nil
}

func main() {
	r := gin.Default()
	r.GET("/orders", getOrders)
	r.GET("/orders/:id", getOrderByID)
	r.POST("/orders", addOrder)
	r.Run(":3002")
}
