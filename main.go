// Build a simple User Management API using Gin with basic HTTP routing and request handling.
package main

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// User represents a user in our system
type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}

// Response represents a standard API response
type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
	Error   string      `json:"error,omitempty"`
	Code    int         `json:"code,omitempty"`
}

// In-memory storage
var users = []User{
	{ID: 1, Name: "John Doe", Email: "john@example.com", Age: 30},
	{ID: 2, Name: "Jane Smith", Email: "jane@example.com", Age: 25},
	{ID: 3, Name: "Bob Wilson", Email: "bob@example.com", Age: 35},
}
var nextID = 4

func main() {
	router := gin.Default()

	router.GET("/users", getAllUsers)
	router.GET("/users/:id", getUserByID)
	router.POST("/users", createUser)
	router.PUT("/users/:id", updateUser)
	router.DELETE("/users/:id", deleteUser)
	router.GET("/users/search", searchUsers)

	router.Run(":8080")
}

// getAllUsers handles GET /users
func getAllUsers(c *gin.Context) {
	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    users,
		Code:    http.StatusOK,
	})
}

// getUserByID handles GET /users/:id
func getUserByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Invalid user ID format. Must be an Integer.",
			Code:    http.StatusBadRequest,
		})
		return
	}
	user, _ := findUserByID(id)
	if user == nil {
		c.JSON(http.StatusNotFound, Response{
			Success: false,
			Error:   "User not found.",
			Code:    http.StatusNotFound,
		})
		return
	}
	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    user,
		Code:    http.StatusOK,
	})
}

// createUser handles POST /users
func createUser(c *gin.Context) {
	var newUser User
	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Invalid JSON format or missing fields.",
			Code:    http.StatusBadRequest,
		})
		return
	}
	if err := validateUser(newUser); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	newUser.ID = nextID
	users = append(users, newUser)
	nextID++

	c.JSON(http.StatusCreated, Response{
		Success: true,
		Message: "User created successfully.",
		Data:    newUser,
		Code:    http.StatusCreated,
	})
}

// updateUser handles PUT /users/:id
func updateUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Invalid ID format.",
			Code:    http.StatusBadRequest,
		})
		return
	}

	existingUser, idx := findUserByID(id)
	if existingUser == nil {
		c.JSON(http.StatusNotFound, Response{
			Success: false,
			Error:   "User not found.",
			Code:    http.StatusNotFound,
		})
		return
	}

	var updatedData User
	if err := c.ShouldBindJSON(&updatedData); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Invalid JSON format.",
			Code:    http.StatusBadRequest,
		})
		return
	}

	if err := validateUser(updatedData); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	updatedData.ID = id
	users[idx] = updatedData
	c.JSON(http.StatusOK, Response{
		Success: true,
		Message: "User updated successfully.",
		Data:    users[idx],
		Code:    http.StatusOK,
	})
}

// deleteUser handles DELETE /users/:id
func deleteUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Invalid ID format.",
			Code:    http.StatusBadRequest,
		})
		return
	}

	user, idx := findUserByID(id)
	if user == nil {
		c.JSON(http.StatusNotFound, Response{
			Success: false,
			Error:   "User not found.",
			Code:    http.StatusNotFound,
		})
		return
	}

	users = append(users[:idx], users[idx+1:]...)
	c.JSON(http.StatusOK, Response{
		Success: true,
		Message: "User deleted successfully.",
		Code:    http.StatusOK,
	})
}

// searchUsers handles GET /users/search?name=value
func searchUsers(c *gin.Context) {
	searchName := c.Query("name")
	if searchName == "" {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Missing 'name' query parameter for search.",
			Code:    http.StatusBadRequest,
		})
		return
	}

	var matchingUsers []User
	searchNameLower := strings.ToLower(searchName)
	for _, user := range users {
		if strings.Contains(strings.ToLower(user.Name), searchNameLower) {
			matchingUsers = append(matchingUsers, user)
		}
	}

	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    matchingUsers,
		Code:    http.StatusOK,
	})
}

// Helper function to find user by ID
func findUserByID(id int) (*User, int) {
	for i, u := range users {
		if u.ID == id {
			return &users[i], i
		}
	}
	return nil, -1
}

// Helper function to validate user data
func validateUser(user User) error {
	if strings.TrimSpace(user.Name) == "" {
		return errors.New("name is required")
	}
	if strings.TrimSpace(user.Email) == "" {
		return errors.New("email is required")
	}
	if !strings.Contains(user.Email, "@") {
		return errors.New("invalid email format")
	}
	if user.Age <= 0 || user.Age > 150 {
		return errors.New("age must be a positive integer (1-150)")
	}
	return nil
}